package board

import (
	"diplomacy-api/internal/models"
	"diplomacy-api/internal/orders"
	"diplomacy-api/internal/utils"
	"fmt"
	"maps"
	"math/rand"
	"slices"
)

func commandAlreadyProcessed(locationID string) string {
	return fmt.Sprintf("Already applied a command for the unit at '%s'.", locationID)
}
func locationNotFound(locationID string) string {
	return fmt.Sprintf("Failed to find '%s'.", locationID)
}

type applicator struct {
	Disrupted        map[string]bool
	ConvoyCarriers   map[string]*node
	BoardMap         map[string]*node
	ReceivedCommands map[string]bool
	IgnoredCommands  map[int]string
	PlayerCensus     map[string]*census
}

type census struct {
	SupplyCenterCount int
	UnitCount         int
}

type node struct {
	Province          models.Province
	Attacks           []attack
	AttackResult      attackResult
	UnresolvedCommand *models.MoveCommand
	Support           *models.SupportCommand
	Convoy            *models.ConvoyCommand
	RetreatUnit       *models.Unit
	RetreatClaims     int
}

type attack struct {
	UnitType   models.UnitType
	LocationID string
	Support    []string
}

type attackResult struct {
	DislodgedUnit *models.Unit
	LocationID    *string
	Standoff      bool
}

func ApplyTurn(board []models.Province, turn *models.Turn) ([]models.Province, []models.CommandValidationResult, error) {
	allCommands := orders.GetCommands(turn.Orders)
	validCommands, validationResults := orders.ValidateCommands(allCommands, board)

	return applyValidatedTurn(board, turn, validCommands, validationResults)
}

func ApplyGameTurn(game *models.Game, turn *models.Turn) ([]models.Province, []models.CommandValidationResult, error) {
	commands, results := orders.ValidateGameCommands(orders.GetCommands(turn.Orders), game, turn)
	return applyValidatedTurn(game.Board, turn, commands, results)
}

func applyValidatedTurn(board []models.Province, turn *models.Turn, validCommands models.Commands, validationResults []models.CommandValidationResult) ([]models.Province, []models.CommandValidationResult, error) {

	a, err := resolveMovement(board, validCommands)
	if err != nil {
		return nil, nil, err
	}
	utils.ForEach(validCommands.Retreat, a.ApplyRetreat)

	isFallTurn := turn.TurnNumber%2 == 0
	a.ResolveRetreats(isFallTurn)
	clear(a.ReceivedCommands)

	if isFallTurn {
		utils.ForEach(validCommands.Disband, a.ApplyDisband)
		utils.ForEach(validCommands.Reinforce, a.ApplyReinforce)
	} else {
		utils.ForEach(validCommands.Disband, func(c models.DisbandCommand) {
			a.IgnoredCommands[c.Ordinal] = "It is not a fall turn."
		})
		utils.ForEach(validCommands.Reinforce, func(c models.ReinforceCommand) {
			a.IgnoredCommands[c.Ordinal] = "It is not a fall turn."
		})
	}

	for i, result := range validationResults {
		if a.IgnoredCommands[result.Ordinal] == "" {
			continue
		}

		result.Valid = false
		result.Error = a.IgnoredCommands[result.Ordinal]

		validationResults[i] = result
	}

	return a.FinalizeState(isFallTurn), validationResults, nil
}

func (a *applicator) ApplyConvoy(convoy models.ConvoyCommand) {
	loc, ok := a.BoardMap[convoy.Location.ID]
	if !ok {
		a.IgnoredCommands[convoy.Ordinal] = locationNotFound(convoy.Location.ID)
		return
	}

	if a.ReceivedCommands[loc.Province.ID] {
		a.IgnoredCommands[convoy.Ordinal] = commandAlreadyProcessed(loc.Province.ID)
		return
	}
	a.ReceivedCommands[loc.Province.ID] = true

	loc.Convoy = &convoy
}

func (a *applicator) ApplyHold(hold models.HoldCommand) {
	loc, ok := a.BoardMap[hold.Location.ID]
	if !ok {
		a.IgnoredCommands[hold.Ordinal] = locationNotFound(hold.Location.ID)
		return
	}

	if a.ReceivedCommands[loc.Province.ID] {
		a.IgnoredCommands[hold.Ordinal] = commandAlreadyProcessed(loc.Province.ID)
		return
	}
	a.ReceivedCommands[loc.Province.ID] = true

	loc.Attacks = append(loc.Attacks, attack{
		UnitType:   hold.UnitType,
		LocationID: hold.Location.ID,
		Support:    nil,
	})
}

func (a *applicator) ApplyMove(move models.MoveCommand) {
	loc, ok := a.BoardMap[move.Location.ID]
	if !ok {
		a.IgnoredCommands[move.Ordinal] = locationNotFound(move.Location.ID)
		return
	}

	dest, ok := a.BoardMap[move.Destination.ID]
	if !ok {
		a.IgnoredCommands[move.Ordinal] = locationNotFound(move.Destination.ID)
		return
	}

	if a.ReceivedCommands[loc.Province.ID] {
		a.IgnoredCommands[move.Ordinal] = commandAlreadyProcessed(loc.Province.ID)
		return
	}
	usesConvoy := move.UnitType.IsArmy() && (move.ViaConvoy ||
		!slices.Contains(dest.Province.Routes, loc.Province.ID) ||
		a.convoyRoute(move, nil, "", true))
	if usesConvoy && !a.convoyRoute(move, nil, "", false) {
		a.IgnoredCommands[move.Ordinal] = fmt.Sprintf("No convoy route from '%s' to '%s'.", loc.Province.ID, dest.Province.ID)
		return
	}

	a.ReceivedCommands[loc.Province.ID] = true
	if usesConvoy {
		// Keep a reference to a matching fleet's order, including when the
		// route is disrupted. The army's move is already stored in that order.
		for _, carrier := range a.BoardMap {
			if carrier.Convoy != nil && carrier.Convoy.Move.Equals(move) {
				a.ConvoyCarriers[loc.Province.ID] = carrier
				break
			}
		}
		if !a.convoyRoute(move, a.Disrupted, "", false) {
			return
		}
	}

	loc.UnresolvedCommand = &move
	dest.Attacks = append(dest.Attacks, attack{
		UnitType:   move.UnitType,
		LocationID: move.Location.ID,
		Support:    nil,
	})
}

func (a *applicator) ApplyRetreat(retreat models.RetreatCommand) {
	loc, ok := a.BoardMap[retreat.Location.ID]
	if !ok {
		a.IgnoredCommands[retreat.Ordinal] = locationNotFound(retreat.Location.ID)
		return
	}

	dest, ok := a.BoardMap[retreat.Destination.ID]
	if !ok {
		a.IgnoredCommands[retreat.Ordinal] = locationNotFound(retreat.Destination.ID)
		return
	}

	if a.ReceivedCommands[loc.Province.ID] {
		a.IgnoredCommands[retreat.Ordinal] = commandAlreadyProcessed(loc.Province.ID)
		return
	}
	if loc.AttackResult.DislodgedUnit == nil {
		a.IgnoredCommands[retreat.Ordinal] = "Unit was not dislodged."
		return
	}
	if dest.Province.Unit != nil {
		a.IgnoredCommands[retreat.Ordinal] = "A unit occupies the destination."
		return
	}
	if loc.AttackResult.LocationID != nil && *loc.AttackResult.LocationID == retreat.Destination.ID {
		a.IgnoredCommands[retreat.Ordinal] = "Cannot retreat to where the attack originated."
		return
	}
	if dest.AttackResult.Standoff {
		a.IgnoredCommands[retreat.Ordinal] = "Cannot retreat to where a standoff occurred."
		return
	}

	unit := loc.AttackResult.DislodgedUnit
	loc.AttackResult.DislodgedUnit = nil
	a.ReceivedCommands[loc.Province.ID] = true
	dest.RetreatClaims++
	if dest.RetreatClaims == 1 {
		unit.Location = retreat.Destination
		dest.RetreatUnit = unit
	} else {
		dest.RetreatUnit = nil
	}
}

func (a *applicator) ApplySupport(support models.SupportCommand) {
	loc, ok := a.BoardMap[support.Location.ID]
	if !ok {
		a.IgnoredCommands[support.Ordinal] = locationNotFound(support.Location.ID)
		return
	}

	dest, ok := a.BoardMap[support.Move.Destination.ID]
	if !ok {
		a.IgnoredCommands[support.Ordinal] = locationNotFound(support.Move.Destination.ID)
		return
	}

	if a.ReceivedCommands[loc.Province.ID] {
		a.IgnoredCommands[support.Ordinal] = commandAlreadyProcessed(loc.Province.ID)
		return
	}
	a.ReceivedCommands[loc.Province.ID] = true

	for _, att := range loc.Attacks {
		attLoc, ok := a.BoardMap[att.LocationID]
		if !ok {
			a.IgnoredCommands[support.Ordinal] = fmt.Sprintf("Support was cut from an unknown province: %s.", att.LocationID)
			return
		}

		isEnemyAttack := attLoc.Province.Unit != nil && attLoc.Province.Unit.ControlledBy != support.PlayerName
		fromDestination := attLoc.Province.ID == support.Move.Destination.ID
		if passenger := a.convoyMove(attLoc); passenger != nil && dest.Convoy != nil &&
			support.Move.Location.ID != dest.Province.ID && dest.Convoy.Move.Equals(*passenger) &&
			!a.convoyRoute(*passenger, a.Disrupted, dest.Province.ID, false) {
			continue
		}
		if isEnemyAttack && !fromDestination {
			a.IgnoredCommands[support.Ordinal] = fmt.Sprintf("Support was cut from '%s'.", att.LocationID)
			return
		}
	}

	att, ok := utils.Find(dest.Attacks, func(a *attack) bool {
		return a.LocationID == support.Move.Location.ID
	})
	if !ok {
		isHold := support.Move.Location.ID == support.Move.Destination.ID
		// A disrupted convoy still selected this move. Its matching support
		// is legal, but there is no attack to strengthen.
		origin := a.BoardMap[support.Move.Location.ID]
		if passenger := a.convoyMove(origin); passenger != nil && passenger.Equals(support.Move) {
			return
		}
		if isHold && dest.UnresolvedCommand == nil && a.convoyMove(dest) == nil {
			dest.Attacks = append(dest.Attacks, attack{
				UnitType:   support.Move.UnitType,
				LocationID: support.Move.Location.ID,
				Support:    []string{support.Location.ID},
			})
			loc.Support = &support
			return
		}

		a.IgnoredCommands[support.Ordinal] = "Could not find the move to support."
		return
	}

	att.Support = append(att.Support, support.Location.ID)
	loc.Support = &support
}

func (a *applicator) withdrawDislodgedSupport(support models.SupportCommand) {
	loc := a.BoardMap[support.Location.ID]
	if loc.AttackResult.DislodgedUnit == nil || a.IgnoredCommands[support.Ordinal] != "" {
		return
	}

	dest := a.BoardMap[support.Move.Destination.ID]
	for i := range dest.Attacks {
		att := &dest.Attacks[i]
		if att.LocationID != support.Move.Location.ID {
			continue
		}

		att.Support = slices.DeleteFunc(att.Support, func(origin string) bool {
			return origin == support.Location.ID
		})
	}

	a.IgnoredCommands[support.Ordinal] = fmt.Sprintf("Support was cut from '%s'.", *loc.AttackResult.LocationID)
}

func (a *applicator) ApplyReinforce(reinforce models.ReinforceCommand) {
	loc, ok := a.BoardMap[reinforce.Location.ID]
	if !ok {
		a.IgnoredCommands[reinforce.Ordinal] = locationNotFound(reinforce.Location.ID)
		return
	}

	if a.ReceivedCommands[loc.Province.ID] {
		a.IgnoredCommands[reinforce.Ordinal] = commandAlreadyProcessed(loc.Province.ID)
		return
	}

	c := a.PlayerCensus[reinforce.PlayerName]
	if c == nil {
		a.IgnoredCommands[reinforce.Ordinal] = "Failed to determine reinforcement eligibility."
		return
	}
	if c.UnitCount >= c.SupplyCenterCount {
		a.IgnoredCommands[reinforce.Ordinal] = "Not enough supply centers to support a new unit."
		return
	}

	if loc.Province.Unit != nil {
		a.IgnoredCommands[reinforce.Ordinal] = "A unit occupies the province."
		return
	}
	if loc.Province.SupplyCenter == nil {
		a.IgnoredCommands[reinforce.Ordinal] = fmt.Sprintf("The province '%s' does not have a supply center.", loc.Province.ID)
		return
	}

	homeFor := loc.Province.SupplyCenter.HomeFor
	if homeFor == nil || *homeFor != reinforce.PlayerName {
		a.IgnoredCommands[reinforce.Ordinal] = "Reinforcements can only be added to home supply centers."
		return
	}

	controlledBy := loc.Province.SupplyCenter.ControlledBy
	if controlledBy == nil || *controlledBy != reinforce.PlayerName {
		a.IgnoredCommands[reinforce.Ordinal] = "This supply center is controlled by another player."
		return
	}

	id, err := utils.RandomID()
	if err != nil {
		a.IgnoredCommands[reinforce.Ordinal] = "Failed to generate an ID for the new unit."
		return
	}

	loc.Province.Unit = &models.Unit{
		ID:           id,
		ControlledBy: reinforce.PlayerName,
		Type:         reinforce.UnitType,
		Location:     reinforce.Location,
	}
	c.UnitCount++

	a.ReceivedCommands[loc.Province.ID] = true
}

func (a *applicator) ApplyDisband(disband models.DisbandCommand) {
	loc, ok := a.BoardMap[disband.Location.ID]
	if !ok {
		a.IgnoredCommands[disband.Ordinal] = locationNotFound(disband.Location.ID)
		return
	}

	if a.ReceivedCommands[loc.Province.ID] {
		a.IgnoredCommands[disband.Ordinal] = commandAlreadyProcessed(loc.Province.ID)
		return
	}

	c := a.PlayerCensus[disband.PlayerName]
	if c == nil {
		a.IgnoredCommands[disband.Ordinal] = "Failed to determine eligibility for disbanding."
		return
	}
	if c.UnitCount <= c.SupplyCenterCount {
		a.IgnoredCommands[disband.Ordinal] = "The number of supply centers exceeds or matches the number of units."
		return
	}

	loc.Province.Unit = nil
	c.UnitCount--

	a.ReceivedCommands[loc.Province.ID] = true
}

func (a *applicator) ResolveAttacks() {
	for _, id := range slices.Sorted(maps.Keys(a.BoardMap)) {
		n := a.BoardMap[id]
		stack := []*node{n}

		for len(stack) > 0 {
			top := stack[len(stack)-1]
			if top.UnresolvedCommand == nil {
				a.resolveAttack(top)
				stack = stack[:len(stack)-1]
				continue
			}

			dest := a.BoardMap[top.UnresolvedCommand.Destination.ID]
			hadCycle := dest.UnresolvedCommand != nil && a.resolveCycle(stack, dest)

			if hadCycle || dest.UnresolvedCommand != nil {
				stack = append(stack, dest)
				continue
			}

			a.resolveAttack(dest)
			a.resolveAttack(top)
			stack = stack[:len(stack)-1]
		}
	}

	clear(a.ReceivedCommands)
}

func (a *applicator) ResolveRetreats(fallTurn bool) {
	for _, n := range a.BoardMap {
		if n.RetreatUnit != nil && n.Province.Unit == nil {
			n.Province.Unit = n.RetreatUnit
		}
		n.AttackResult.DislodgedUnit = nil

		hasSupplyCenter := n.Province.SupplyCenter != nil
		hasUnit := n.Province.Unit != nil
		if fallTurn && hasSupplyCenter && hasUnit {
			n.Province.SupplyCenter.ControlledBy = &n.Province.Unit.ControlledBy
		}

		if hasSupplyCenter && n.Province.SupplyCenter.ControlledBy != nil {
			playerName := *n.Province.SupplyCenter.ControlledBy

			if a.PlayerCensus[playerName] == nil {
				a.PlayerCensus[playerName] = &census{}
			}
			a.PlayerCensus[playerName].SupplyCenterCount++
		}
		if hasUnit {
			playerName := n.Province.Unit.ControlledBy

			if a.PlayerCensus[playerName] == nil {
				a.PlayerCensus[playerName] = &census{}
			}
			a.PlayerCensus[playerName].UnitCount++
		}
	}
}

func (a *applicator) FinalizeState(fallTurn bool) []models.Province {
	board := make([]models.Province, 0)

	if fallTurn {
		for player, c := range a.PlayerCensus {
			if c.UnitCount > c.SupplyCenterCount {
				a.removeRandomUnits(player, c.UnitCount-c.SupplyCenterCount)
			}
		}
	}

	for _, n := range a.BoardMap {
		board = append(board, n.Province)
	}

	return board
}

func (a *applicator) resolveAttack(n *node) {
	minStrength := 0
	if n.Province.Unit != nil {
		minStrength = 1
	}

	maxStrength, strongestAttacks := a.strongestAttacks(n)

	if len(strongestAttacks) == 1 && maxStrength > minStrength {
		att := strongestAttacks[0]
		attLoc := a.BoardMap[att.LocationID]

		attackingOwnUnit := n.Province.Unit != nil && n.Province.Unit.ControlledBy == attLoc.Province.Unit.ControlledBy
		if n.Province.ID != attLoc.Province.ID && !attackingOwnUnit {
			unit := attLoc.Province.Unit
			unit.Location = attLoc.UnresolvedCommand.Destination

			n.AttackResult.DislodgedUnit = n.Province.Unit
			n.AttackResult.LocationID = &att.LocationID

			n.Province.Unit = unit
			attLoc.Province.Unit = nil
		}
	} else if len(strongestAttacks) > 0 {
		n.AttackResult.Standoff = true
	}

	if n.Support != nil {
		a.withdrawDislodgedSupport(*n.Support)
	}

	n.Attacks = nil
	n.UnresolvedCommand = nil
}

func (a *applicator) strongestAttacks(n *node) (int, []*attack) {
	var maxStrength int
	strongest := make([]*attack, 0)

	for {
		maxStrength = 0
		strongest = strongest[:0]

		for i := range n.Attacks {
			att := &n.Attacks[i]
			strength := len(att.Support) + 1
			switch {
			case strength > maxStrength:
				maxStrength = strength
				strongest = []*attack{att}
			case strength == maxStrength:
				strongest = append(strongest, att)
			}
		}

		removed := false
		if len(strongest) == 1 && n.Province.Unit != nil {
			att := strongest[0]

			for index, locID := range att.Support {
				supportLoc, ok := a.BoardMap[locID]
				if !ok {
					continue
				}

				if supportLoc.Province.Unit == nil {
					continue
				}

				if supportLoc.Province.Unit.ControlledBy == n.Province.Unit.ControlledBy {
					att.Support = append(att.Support[:index], att.Support[index+1:]...)
					removed = true
					break
				}
			}

		}

		if !removed {
			break
		}
	}

	return maxStrength, strongest
}

func (a *applicator) removeRandomUnits(playerName string, amount int) {
	playerNodes := make([]*node, 0)

	for _, n := range a.BoardMap {
		unit := n.Province.Unit
		if unit != nil && unit.ControlledBy == playerName {
			playerNodes = append(playerNodes, n)
		}
	}

	for amount > 0 && len(playerNodes) > 0 {
		randIndex := rand.Intn(len(playerNodes))
		playerNodes[randIndex].Province.Unit = nil
		a.PlayerCensus[playerName].UnitCount--

		playerNodes[randIndex] = playerNodes[len(playerNodes)-1]
		playerNodes = playerNodes[:len(playerNodes)-1]
		amount--
	}
}

func (a *applicator) resolveHeadToHead(origin, dest *node) bool {
	if origin.Province.Unit.ControlledBy == dest.Province.Unit.ControlledBy {
		return false
	}

	strength, opposition := 0, 1
	for _, bid := range dest.Attacks {
		if bid.LocationID != origin.Province.ID {
			opposition = max(opposition, 1+len(bid.Support))
			continue
		}
		strength = 1
		for _, id := range bid.Support {
			unit := a.BoardMap[id].Province.Unit
			if unit != nil && unit.ControlledBy != dest.Province.Unit.ControlledBy {
				strength++
			}
		}
	}
	for _, bid := range origin.Attacks {
		if bid.LocationID == dest.Province.ID {
			opposition = max(opposition, 1+len(bid.Support))
			break
		}
	}
	if strength <= opposition {
		return false
	}

	a.resolveAttack(dest)
	origin.UnresolvedCommand = nil
	origin.Attacks = slices.DeleteFunc(origin.Attacks, func(bid attack) bool {
		return bid.LocationID == dest.Province.ID
	})
	return true
}

func rotateUnits(cycle []*node) {
	units := make([]*models.Unit, len(cycle))
	for i, origin := range cycle {
		units[i] = origin.Province.Unit
		units[i].Location = origin.UnresolvedCommand.Destination
	}
	for i, origin := range cycle {
		dest := cycle[(i+1)%len(cycle)]
		dest.Province.Unit = units[i]
		origin.Attacks = nil
		origin.UnresolvedCommand = nil
	}
}

func (a *applicator) resolveCycle(stack []*node, newNode *node) bool {
	start := slices.Index(stack, newNode)
	if start < 0 {
		return false
	}
	cycle := stack[start:]

	exchangeByConvoy := len(cycle) == 2 && (a.convoyMove(cycle[0]) != nil || a.convoyMove(cycle[1]) != nil)
	if len(cycle) == 2 && !exchangeByConvoy {
		left, right := cycle[0], cycle[1]
		if a.resolveHeadToHead(left, right) || a.resolveHeadToHead(right, left) {
			return true
		}
	}
	if len(cycle) >= 3 || exchangeByConvoy {
		for i, origin := range cycle {
			dest := cycle[(i+1)%len(cycle)]
			if len(dest.Attacks) != 1 || dest.Attacks[0].LocationID != origin.Province.ID {
				a.resolveAttack(newNode)
				return true
			}
		}
		rotateUnits(cycle)
		return true
	}

	a.resolveAttack(newNode)
	return true
}
