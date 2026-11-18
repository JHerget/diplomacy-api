package board

import (
	"diplomacy-api/internal/models"
	"diplomacy-api/internal/orders"
	"diplomacy-api/internal/utils"
	"slices"
)

type applicator struct {
	BoardMap         map[string]*node
	ReceivedCommands map[string]bool
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
	}

	return a.FinalizeState(), validationResults, nil
}

func (a *applicator) ApplyConvoy(convoy models.ConvoyCommand) {
	loc, ok := a.BoardMap[convoy.Location.ID]
	if !ok {
		return
	}

	if a.ReceivedCommands[loc.Province.ID] {
		return
	}
	a.ReceivedCommands[loc.Province.ID] = true

	loc.Convoy = &convoy
}

func (a *applicator) ApplyMove(move models.MoveCommand) {
	loc, ok := a.BoardMap[move.Location.ID]
	if !ok {
		return
	}

	dest, ok := a.BoardMap[move.Destination.ID]
	if !ok {
		return
	}

	if a.ReceivedCommands[loc.Province.ID] {
		return
	}
	a.ReceivedCommands[loc.Province.ID] = true

	needsConvoy := !slices.Contains(dest.Province.Routes, loc.Province.ID)
	if needsConvoy && !a.hasConvoyRoute(loc, dest, move) {
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
		return
	}

	dest, ok := a.BoardMap[retreat.Destination.ID]
	if !ok {
		return
	}

	if a.ReceivedCommands[loc.Province.ID] {
		return
	}
	if dest.Province.Unit != nil {
		return
	}
	if loc.AttackResult.LocationID != nil && *loc.AttackResult.LocationID == retreat.Destination.ID {
		return
	}
	if dest.AttackResult.Standoff {
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
		return
	}

	dest, ok := a.BoardMap[support.Move.Destination.ID]
	if !ok {
		return
	}

	if a.ReceivedCommands[loc.Province.ID] {
		return
	}
	a.ReceivedCommands[loc.Province.ID] = true

	for _, att := range loc.Attacks {
		attLoc, ok := a.BoardMap[att.LocationID]
		if !ok {
			return
		}

		if attLoc.Province.Unit != nil && attLoc.Province.Unit.ControlledBy != support.PlayerName {
			return
		}
	}

	att, ok := utils.Find(dest.Attacks, func(a *attack) bool {
		return a.LocationID == support.Move.Location.ID
	})
	if !ok {
		return
	}

	att.Support = append(att.Support, support.Location.ID)
}

func (a *applicator) ApplyReinforce(reinforce models.ReinforceCommand) {
	loc, ok := a.BoardMap[reinforce.Location.ID]
	if !ok {
		return
	}

	if a.ReceivedCommands[loc.Province.ID] {
		return
	}

	c := a.PlayerCensus[reinforce.PlayerName]
	if c == nil {
		return
	}
	if c.UnitCount >= c.SupplyCenterCount {
		return
	}

	if loc.Province.Unit != nil {
		return
	}
	if loc.Province.SupplyCenter == nil {
		return
	}

	homeFor := loc.Province.SupplyCenter.HomeFor
	if homeFor == nil || *homeFor != reinforce.PlayerName {
		return
	}

	controlledBy := loc.Province.SupplyCenter.ControlledBy
	if controlledBy == nil || controlledBy != &reinforce.PlayerName {
		return
	}

	id, err := utils.RandomID()
	if err != nil {
		return
	}

	loc.Province.Unit = &models.Unit{
		ID: id,
		ControlledBy: reinforce.PlayerName,
		Type: reinforce.UnitType,
		Location: reinforce.Location,
	}
	a.ReceivedCommands[loc.Province.ID] = true
}

func (a *applicator) ApplyDisband(disband models.DisbandCommand) {
	loc, ok := a.BoardMap[disband.Location.ID]
	if !ok {
		return
	}

	if a.ReceivedCommands[loc.Province.ID] {
		return
	}
	a.ReceivedCommands[loc.Province.ID] = true
}

func (a *applicator) ResolveAttacks() {
	for _, n := range a.BoardMap {
		if n.UnresolvedCommand == nil {
			a.resolveAttack(n)
			continue
		}

		stack := []*node{n}

		for len(stack) > 0 {
			top := stack[len(stack)-1]
			dest := a.BoardMap[top.UnresolvedCommand.Destination.ID]

			if dest.UnresolvedCommand != nil {
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

	maxStrength := 0
	for _, att := range n.Attacks {
		strength := len(att.Support) + 1

		if strength > maxStrength {
			maxStrength = strength
		}
	}

	strongestAttacks := utils.Filter(n.Attacks, func(a *attack) bool {
		return len(a.Support)+1 >= maxStrength
	})

	if len(strongestAttacks) == 1 && maxStrength > minStrength {
		att := strongestAttacks[0]
		attLoc := a.BoardMap[att.LocationID]

		if n.Province.ID != attLoc.Province.ID {
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
