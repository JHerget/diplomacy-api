package board

import (
	"diplomacy-api/internal/models"
	"diplomacy-api/internal/orders"
	"diplomacy-api/internal/utils"
	"fmt"
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
	BoardMap         map[string]*node
	ReceivedCommands map[string]bool
	IgnoredCommands  map[string]string
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
	Convoy            *models.ConvoyCommand
	RetreatUnit       *models.Unit
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

	a := applicator{
		BoardMap: make(map[string]*node),
	}
	for _, p := range board {
		a.BoardMap[p.ID] = &node{
			Province:          p,
			Attacks:           []attack{},
			UnresolvedCommand: nil,
		}
	}
	slices.SortFunc(validCommands.Move, a.compareMove)

	utils.ForEach(validCommands.Convoy, a.ApplyConvoy)
	utils.ForEach(validCommands.Move, a.ApplyMove)
	utils.ForEach(validCommands.Support, a.ApplySupport)
	a.ResolveAttacks()
	utils.ForEach(validCommands.Retreat, a.ApplyRetreat)

	isFallTurn := turn.TurnNumber%2 == 0
	a.ResolveRetreats(isFallTurn)

	if isFallTurn {
		utils.ForEach(validCommands.Disband, a.ApplyDisband)
		utils.ForEach(validCommands.Reinforce, a.ApplyReinforce)
	} else {
		utils.ForEach(validCommands.Disband, func(c models.DisbandCommand) {
			a.IgnoredCommands[c.Value] = "It is not a fall turn."
		})
		utils.ForEach(validCommands.Reinforce, func(c models.ReinforceCommand) {
			a.IgnoredCommands[c.Value] = "It is not a fall turn."
		})
	}

	for i, result := range validationResults {
		if a.IgnoredCommands[result.Value] == "" {
			continue
		}

		result.Valid = false
		result.Error = a.IgnoredCommands[result.Value]

		validationResults[i] = result
	}

	return a.FinalizeState(), validationResults, nil
}

func (a *applicator) ApplyConvoy(convoy models.ConvoyCommand) {
	loc, ok := a.BoardMap[convoy.Location.ID]
	if !ok {
		a.IgnoredCommands[convoy.Value] = locationNotFound(convoy.Location.ID)
		return
	}

	if a.ReceivedCommands[loc.Province.ID] {
		a.IgnoredCommands[convoy.Value] = commandAlreadyProcessed(loc.Province.ID)
		return
	}
	a.ReceivedCommands[loc.Province.ID] = true

	loc.Convoy = &convoy
}

func (a *applicator) ApplyMove(move models.MoveCommand) {
	loc, ok := a.BoardMap[move.Location.ID]
	if !ok {
		a.IgnoredCommands[move.Value] = locationNotFound(move.Location.ID)
		return
	}

	dest, ok := a.BoardMap[move.Destination.ID]
	if !ok {
		a.IgnoredCommands[move.Value] = locationNotFound(move.Destination.ID)
		return
	}

	if a.ReceivedCommands[loc.Province.ID] {
		a.IgnoredCommands[move.Value] = commandAlreadyProcessed(loc.Province.ID)
		return
	}
	a.ReceivedCommands[loc.Province.ID] = true

	needsConvoy := !slices.Contains(dest.Province.Routes, loc.Province.ID)
	if needsConvoy && !a.hasConvoyRoute(loc, dest, move) {
		a.IgnoredCommands[move.Value] = fmt.Sprintf("No convoy route from '%s' to '%s'.", loc.Province.ID, dest.Province.ID)
		return
	}

	loc.UnresolvedCommand = &move
	dest.Attacks = append(dest.Attacks, attack{
		UnitType:   move.UnitType,
		LocationID: move.Location.ID,
		Support:    nil,
	})

	if dest.Convoy != nil {
		dest.Convoy = nil
	}
}

func (a *applicator) ApplyRetreat(retreat models.RetreatCommand) {
	loc, ok := a.BoardMap[retreat.Location.ID]
	if !ok {
		a.IgnoredCommands[retreat.Value] = locationNotFound(retreat.Location.ID)
		return
	}

	dest, ok := a.BoardMap[retreat.Destination.ID]
	if !ok {
		a.IgnoredCommands[retreat.Value] = locationNotFound(retreat.Destination.ID)
		return
	}

	if a.ReceivedCommands[loc.Province.ID] {
		a.IgnoredCommands[retreat.Value] = commandAlreadyProcessed(loc.Province.ID)
		return
	}
	if dest.Province.Unit != nil {
		a.IgnoredCommands[retreat.Value] = "A unit occupies the destination."
		return
	}
	if loc.AttackResult.LocationID != nil && *loc.AttackResult.LocationID == retreat.Destination.ID {
		a.IgnoredCommands[retreat.Value] = "Cannot retreat to where the attack originated."
		return
	}
	if dest.AttackResult.Standoff {
		a.IgnoredCommands[retreat.Value] = "Cannot retreat to where a standoff occurred."
		return
	}

	if dest.RetreatUnit != nil {
		dest.RetreatUnit = nil
		return
	}

	dest.RetreatUnit = loc.AttackResult.DislodgedUnit
	loc.AttackResult.DislodgedUnit = nil

	a.ReceivedCommands[loc.Province.ID] = true
}

func (a *applicator) ApplySupport(support models.SupportCommand) {
	loc, ok := a.BoardMap[support.Location.ID]
	if !ok {
		a.IgnoredCommands[support.Value] = locationNotFound(support.Location.ID)
		return
	}

	dest, ok := a.BoardMap[support.Move.Destination.ID]
	if !ok {
		a.IgnoredCommands[support.Value] = locationNotFound(support.Move.Destination.ID)
		return
	}

	if a.ReceivedCommands[loc.Province.ID] {
		a.IgnoredCommands[support.Value] = commandAlreadyProcessed(loc.Province.ID)
		return
	}
	a.ReceivedCommands[loc.Province.ID] = true

	for _, att := range loc.Attacks {
		attLoc, ok := a.BoardMap[att.LocationID]
		if !ok {
			a.IgnoredCommands[support.Value] = fmt.Sprintf("Support was cut from an unknown province: %s.", att.LocationID)
			return
		}

		enemyAttack := attLoc.Province.Unit != nil && attLoc.Province.Unit.ControlledBy != support.PlayerName
		fromDestination := attLoc.Province.ID == support.Move.Destination.ID
		if enemyAttack && !fromDestination {
			a.IgnoredCommands[support.Value] = fmt.Sprintf("Support was cut from '%s'.", att.LocationID)
			return
		}
	}

	att, ok := utils.Find(dest.Attacks, func(a *attack) bool {
		return a.LocationID == support.Move.Location.ID
	})
	if !ok {
		a.IgnoredCommands[support.Value] = "Could not find the move to support."
		return
	}

	att.Support = append(att.Support, support.Location.ID)
}

func (a *applicator) ApplyReinforce(reinforce models.ReinforceCommand) {
	loc, ok := a.BoardMap[reinforce.Location.ID]
	if !ok {
		a.IgnoredCommands[reinforce.Value] = locationNotFound(reinforce.Location.ID)
		return
	}

	if a.ReceivedCommands[loc.Province.ID] {
		a.IgnoredCommands[reinforce.Value] = commandAlreadyProcessed(loc.Province.ID)
		return
	}

	c := a.PlayerCensus[reinforce.PlayerName]
	if c == nil {
		a.IgnoredCommands[reinforce.Value] = "Failed to determine reinforcement eligibility."
		return
	}
	if c.UnitCount >= c.SupplyCenterCount {
		a.IgnoredCommands[reinforce.Value] = "Not enough supply centers to support a new unit."
		return
	}

	if loc.Province.Unit != nil {
		a.IgnoredCommands[reinforce.Value] = "A unit occupies the province."
		return
	}
	if loc.Province.SupplyCenter == nil {
		a.IgnoredCommands[reinforce.Value] = fmt.Sprintf("The province '%s' does not have a supply center.", loc.Province.ID)
		return
	}

	homeFor := loc.Province.SupplyCenter.HomeFor
	if homeFor == nil || *homeFor != reinforce.PlayerName {
		a.IgnoredCommands[reinforce.Value] = "Reinforcements can only be added to home supply centers."
		return
	}

	controlledBy := loc.Province.SupplyCenter.ControlledBy
	if controlledBy == nil || controlledBy != &reinforce.PlayerName {
		a.IgnoredCommands[reinforce.Value] = "This supply center is controlled by another player."
		return
	}

	id, err := utils.RandomID()
	if err != nil {
		a.IgnoredCommands[reinforce.Value] = "Failed to generate an ID for the new unit."
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
		a.IgnoredCommands[disband.Value] = locationNotFound(disband.Location.ID)
		return
	}

	if a.ReceivedCommands[loc.Province.ID] {
		a.IgnoredCommands[disband.Value] = commandAlreadyProcessed(loc.Province.ID)
		return
	}

	c := a.PlayerCensus[disband.PlayerName]
	if c == nil {
		a.IgnoredCommands[disband.Value] = "Failed to determine eligibility for disbanding."
		return
	}
	if c.UnitCount <= c.SupplyCenterCount {
		a.IgnoredCommands[disband.Value] = "The number of supply centers exceeds or matches the number of units."
		return
	}

	loc.Province.Unit = nil
	c.UnitCount--

	a.ReceivedCommands[loc.Province.ID] = true
}

func (a *applicator) ResolveAttacks() {
	for _, n := range a.BoardMap {
		stack := []*node{n}

		for len(stack) > 0 {
			top := stack[len(stack)-1]
			if top.UnresolvedCommand == nil {
				a.resolveAttack(top)
				stack = stack[:len(stack)-1]
				continue
			}

			dest := a.BoardMap[top.UnresolvedCommand.Destination.ID]
			hadCycle := a.resolveCycle(stack, dest)

			if hadCycle || dest.UnresolvedCommand != nil {
				stack = append(stack, dest)
				continue
			}

			a.resolveAttack(dest)
			a.resolveAttack(top)
			stack = stack[:len(stack)-1]
		}
	}

	a.ReceivedCommands = nil
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

		if hasSupplyCenter {
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

func (a *applicator) FinalizeState() []models.Province {
	board := make([]models.Province, 0)

	for player, c := range a.PlayerCensus {
		if c.UnitCount > c.SupplyCenterCount {
			a.removeRandomUnits(player, c.UnitCount-c.SupplyCenterCount)
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
	} else {
		n.AttackResult.Standoff = len(strongestAttacks) > 0
	}

	n.Attacks = nil
	n.UnresolvedCommand = nil
	n.Convoy = nil
}

func (a *applicator) strongestAttacks(n *node) (int, []*attack) {
	var maxStrength int
	strongest := make([]*attack, 0)

	for {
		maxStrength = 0
		strongest = strongest[:0]

		for i := range n.Attacks {
			att := &n.Attacks[i]
			strength := len(att.Support)+1
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

func (a *applicator) compareMove(lhs, rhs models.MoveCommand) int {
	lhsDest := a.BoardMap[lhs.Destination.ID]
	lhsNeighbors := slices.Contains(lhsDest.Province.Routes, lhs.Location.ID)

	rhsDest := a.BoardMap[rhs.Destination.ID]
	rhsNeighbors := slices.Contains(rhsDest.Province.Routes, rhs.Location.ID)

	if lhsNeighbors && !rhsNeighbors {
		return -1
	}
	if !lhsNeighbors && rhsNeighbors {
		return 1
	}

	return 0
}

func (a *applicator) hasConvoyRoute(location, destination *node, move models.MoveCommand) bool {
	if !location.Province.IsCoastal() || !destination.Province.IsCoastal() {
		return false
	}

	return a.hasConvoyRouteRecursiveHelper(location, destination, move, make([]string, 0))
}

func (a *applicator) hasConvoyRouteRecursiveHelper(location, destination *node, move models.MoveCommand, visited []string) bool {
	if slices.Contains(visited, location.Province.ID) {
		return false
	}
	if location.Province.ID == destination.Province.ID {
		return false
	}

	for _, province := range location.Province.Routes {
		if province == destination.Province.ID {
			return true
		}

		newLocation, ok := a.BoardMap[province]
		if !ok {
			continue
		}

		if newLocation.Convoy == nil || !newLocation.Convoy.Move.Equals(move) {
			continue
		}
		if a.hasConvoyRouteRecursiveHelper(newLocation, destination, move, append(visited, location.Province.ID)) {
			return true
		}
	}

	return false
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

func (a *applicator) rotateUnits(nodes []*node) {
	if len(nodes) < 2 {
		return
	}

	unit := nodes[0].Province.Unit
	for _, n := range nodes[1:] {
		newUnit := n.Province.Unit
		n.Province.Unit = unit
		n.UnresolvedCommand = nil
		unit = newUnit
	}
	nodes[0].Province.Unit = unit
	nodes[0].UnresolvedCommand = nil
}

func (a *applicator) resolveCycle(stack []*node, newNode *node) bool {
	var cycle []*node
	for i, elem := range stack {
		if elem.Province.ID == newNode.Province.ID {
			cycle = stack[i:]
			break
		}
	}
	if len(cycle) == 0 {
		return false
	}

	cycle = append(cycle, newNode)
	
	var resolvable *node
	for _, elem := range cycle {
		_, strongestAttacks := a.strongestAttacks(elem)

		if len(strongestAttacks) == 1  && strongestAttacks[0].LocationID == elem.Province.ID {
			continue
		}

		resolvable = elem
		break
	}

	if resolvable != nil {
		a.resolveAttack(resolvable)
	} else {
		a.rotateUnits(cycle)
	}

	return true
}
