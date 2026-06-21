package orders

import (
	"diplomacy-api/internal/game"
	"diplomacy-api/internal/utils"
	"slices"
)

type validator struct {
	BoardMap map[string]game.Providence
}

func ValidateCommands(commands game.Commands, board []game.Providence) game.Commands {
	v := validator{
		BoardMap: make(map[string]game.Providence),
	}
	for _, p := range board {
		v.BoardMap[p.Id] = p
	}

	return game.Commands{
		Hold: utils.Filter(commands.Hold, v.IsValidHold),
		Move: utils.Filter(commands.Move, func(c *game.MoveCommand) bool {
			return v.IsValidMove(c, true)
		}),
		Retreat:   utils.Filter(commands.Retreat, v.IsValidRetreat),
		Support:   utils.Filter(commands.Support, v.IsValidSupport),
		Convoy:    utils.Filter(commands.Convoy, v.IsValidConvoy),
		Reinforce: utils.Filter(commands.Reinforce, v.IsValidReinforce),
		Disband:   utils.Filter(commands.Disband, v.IsValidDisband),
	}
}

func (v validator) IsValidHold(hold *game.MoveCommand) bool {
	loc, ok := v.BoardMap[hold.Location.Id]
	if !ok {
		return false
	}

	if !v.IsValidLocation(hold.Location) {
		return false
	}

	if loc.Unit == nil {
		return false
	}
	if loc.Unit.Type != hold.UnitType {
		return false
	}
	if loc.Unit.ControlledBy != hold.PlayerName {
		return false
	}
	if hold.Location.Id != hold.Destination.Id {
		return false
	}
	if hold.Location.Coast != hold.Destination.Coast {
		return false
	}

	return true
}

func (v validator) IsValidMove(move *game.MoveCommand, validatePlayer bool) bool {
	loc, ok := v.BoardMap[move.Location.Id]
	if !ok {
		return false
	}

	dest, ok := v.BoardMap[move.Destination.Id]
	if !ok {
		return false
	}

	if !v.IsValidLocation(move.Location) {
		return false
	}
	if !v.IsValidLocation(move.Destination) {
		return false
	}

	if loc.Unit == nil {
		return false
	}
	if loc.Unit.Type != move.UnitType {
		return false
	}
	if validatePlayer && loc.Unit.ControlledBy != move.PlayerName {
		return false
	}

	if !v.IsValidRoute(loc, dest, move.Destination) {
		return false
	}

	return true
}

func (v validator) IsValidRetreat(retreat *game.MoveCommand) bool {
	loc, ok := v.BoardMap[retreat.Location.Id]
	if !ok {
		return false
	}

	dest, ok := v.BoardMap[retreat.Destination.Id]
	if !ok {
		return false
	}

	if !v.IsValidLocation(retreat.Location) {
		return false
	}
	if !v.IsValidLocation(retreat.Destination) {
		return false
	}

	if loc.Unit == nil {
		return false
	}
	if loc.Unit.Type != retreat.UnitType {
		return false
	}
	if loc.Unit.ControlledBy != retreat.PlayerName {
		return false
	}

	if !v.IsValidRoute(loc, dest, retreat.Destination) {
		return false
	}

	return true
}

func (v validator) IsValidSupport(support *game.SupportCommand) bool {
	loc, ok := v.BoardMap[support.Location.Id]
	if !ok {
		return false
	}

	dest, ok := v.BoardMap[support.Move.Destination.Id]
	if !ok {
		return false
	}

	if !v.IsValidMove(&support.Move, false) {
		return false
	}
	if !v.IsValidLocation(support.Location) {
		return false
	}

	if loc.Unit == nil {
		return false
	}
	if loc.Unit.Type != support.UnitType {
		return false
	}
	if loc.Unit.ControlledBy != support.PlayerName {
		return false
	}

	if !v.IsValidRoute(loc, dest, support.Move.Destination) {
		return false
	}

	return true
}

func (v validator) IsValidConvoy(convoy *game.SupportCommand) bool {
	loc, ok := v.BoardMap[convoy.Location.Id]
	if !ok {
		return false
	}

	dest, ok := v.BoardMap[convoy.Move.Destination.Id]
	if !ok {
		return false
	}

	if !v.IsValidMove(&convoy.Move, false) {
		return false
	}
	if !v.IsValidLocation(convoy.Location) {
		return false
	}

	if loc.Unit == nil {
		return false
	}
	if loc.Unit.Type != convoy.UnitType {
		return false
	}
	if loc.Unit.ControlledBy != convoy.PlayerName {
		return false
	}

	if !v.IsValidRoute(loc, dest, convoy.Move.Destination) {
		return false
	}

	return true
}

func (v validator) IsValidReinforce(reinforce *game.AdjustCommand) bool {
	loc, ok := v.BoardMap[reinforce.Location.Id]
	if !ok {
		return false
	}

	if !v.IsValidLocation(reinforce.Location) {
		return false
	}

	if loc.Unit != nil {
		return false
	}
	if loc.SupplyCenter == nil {
		return false
	}
	if loc.SupplyCenter.ControlledBy == nil {
		return false
	}
	if *loc.SupplyCenter.ControlledBy != reinforce.PlayerName {
		return false
	}
	if reinforce.UnitType == game.UnitArmy && loc.Type == game.ProvidenceOcean {
		return false
	}
	if reinforce.UnitType == game.UnitFleet && loc.Type == game.ProvidenceInland {
		return false
	}

	return true
}

func (v validator) IsValidDisband(disband *game.AdjustCommand) bool {
	loc, ok := v.BoardMap[disband.Location.Id]
	if !ok {
		return false
	}

	if !v.IsValidLocation(disband.Location) {
		return false
	}

	if loc.Unit == nil {
		return false
	}
	if loc.Unit.Type != disband.UnitType {
		return false
	}
	if loc.Unit.ControlledBy != disband.PlayerName {
		return false
	}
	if loc.SupplyCenter == nil {
		return false
	}
	if loc.SupplyCenter.ControlledBy == nil {
		return false
	}
	if *loc.SupplyCenter.ControlledBy != disband.PlayerName {
		return false
	}

	return true
}

func (v validator) IsValidLocation(ref game.LocationReference) bool {
	loc, ok := v.BoardMap[ref.Id]
	if !ok {
		return false
	}

	coastProvided := ref.Coast != nil
	hasCoastalRoutes := len(loc.CoastalRoutes) > 0

	if !coastProvided && hasCoastalRoutes {
		return false
	}
	if coastProvided && !hasCoastalRoutes {
		return false
	}

	if !coastProvided {
		return true
	}

	_, routesExist := loc.CoastalRoutes[*ref.Coast]
	if !routesExist {
		return false
	}

	return true
}

func (v validator) IsValidRoute(location game.Providence, destination game.Providence, destRef game.LocationReference) bool {
	if location.Unit == nil {
		return false
	}

	if location.Type == game.ProvidenceOcean && destination.Type == game.ProvidenceInland {
		return false
	}
	if location.Type == game.ProvidenceInland && destination.Type == game.ProvidenceOcean {
		return false
	}
	if location.Type == game.ProvidenceInland && location.Unit.Type == game.UnitFleet {
		return false
	}
	if location.Type == game.ProvidenceOcean && location.Unit.Type == game.UnitArmy {
		return false
	}

	hasCoastalRoutes := len(destination.CoastalRoutes) > 0

	if destRef.Coast != nil {
		if !hasCoastalRoutes {
			return slices.Contains(destination.Routes, location.Id)
		}

		routes, ok := destination.CoastalRoutes[*destRef.Coast]
		if !ok {
			return false
		}

		return slices.Contains(routes, location.Id)
	}
	if hasCoastalRoutes {
		return false
	}

	return slices.Contains(destination.Routes, location.Id)
}
