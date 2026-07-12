package orders

import (
	"diplomacy-api/internal/models"
	"diplomacy-api/internal/utils"
	"slices"
)

type validator struct {
	BoardMap map[string]models.Providence
}

func ValidateCommands(commands models.Commands, board []models.Providence) models.Commands {
	v := validator{
		BoardMap: make(map[string]models.Providence),
	}
	for _, p := range board {
		v.BoardMap[p.ID] = p
	}

	return models.Commands{
		Hold: utils.Filter(commands.Hold, v.IsValidHold),
		Move: utils.Filter(commands.Move, func(c *models.MoveCommand) bool {
			return v.IsValidMove(c, true)
		}),
		Retreat:   utils.Filter(commands.Retreat, v.IsValidRetreat),
		Support:   utils.Filter(commands.Support, v.IsValidSupport),
		Convoy:    utils.Filter(commands.Convoy, v.IsValidConvoy),
		Reinforce: utils.Filter(commands.Reinforce, v.IsValidReinforce),
		Disband:   utils.Filter(commands.Disband, v.IsValidDisband),
	}
}

func (v validator) IsValidHold(hold *models.MoveCommand) bool {
	loc, ok := v.BoardMap[hold.Location.ID]
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
	if hold.Location.ID != hold.Destination.ID {
		return false
	}
	if hold.Location.Coast != hold.Destination.Coast {
		return false
	}

	return true
}

func (v validator) IsValidMove(move *models.MoveCommand, validatePlayer bool) bool {
	loc, ok := v.BoardMap[move.Location.ID]
	if !ok {
		return false
	}

	dest, ok := v.BoardMap[move.Destination.ID]
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

func (v validator) IsValidRetreat(retreat *models.MoveCommand) bool {
	loc, ok := v.BoardMap[retreat.Location.ID]
	if !ok {
		return false
	}

	dest, ok := v.BoardMap[retreat.Destination.ID]
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

func (v validator) IsValidSupport(support *models.SupportCommand) bool {
	loc, ok := v.BoardMap[support.Location.ID]
	if !ok {
		return false
	}

	dest, ok := v.BoardMap[support.Move.Destination.ID]
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

func (v validator) IsValidConvoy(convoy *models.SupportCommand) bool {
	loc, ok := v.BoardMap[convoy.Location.ID]
	if !ok {
		return false
	}

	dest, ok := v.BoardMap[convoy.Move.Destination.ID]
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

func (v validator) IsValidReinforce(reinforce *models.AdjustCommand) bool {
	loc, ok := v.BoardMap[reinforce.Location.ID]
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
	if reinforce.UnitType == models.UnitArmy && loc.Type == models.ProvidenceOcean {
		return false
	}
	if reinforce.UnitType == models.UnitFleet && loc.Type == models.ProvidenceInland {
		return false
	}

	return true
}

func (v validator) IsValidDisband(disband *models.AdjustCommand) bool {
	loc, ok := v.BoardMap[disband.Location.ID]
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

func (v validator) IsValidLocation(ref models.LocationReference) bool {
	loc, ok := v.BoardMap[ref.ID]
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

func (v validator) IsValidRoute(location models.Providence, destination models.Providence, destRef models.LocationReference) bool {
	if location.Unit == nil {
		return false
	}

	if location.Type == models.ProvidenceOcean && destination.Type == models.ProvidenceInland {
		return false
	}
	if location.Type == models.ProvidenceInland && destination.Type == models.ProvidenceOcean {
		return false
	}
	if location.Type == models.ProvidenceInland && location.Unit.Type == models.UnitFleet {
		return false
	}
	if location.Type == models.ProvidenceOcean && location.Unit.Type == models.UnitArmy {
		return false
	}

	hasCoastalRoutes := len(destination.CoastalRoutes) > 0

	if destRef.Coast != nil {
		if !hasCoastalRoutes {
			return slices.Contains(destination.Routes, location.ID)
		}

		routes, ok := destination.CoastalRoutes[*destRef.Coast]
		if !ok {
			return false
		}

		return slices.Contains(routes, location.ID)
	}
	if hasCoastalRoutes {
		return false
	}

	return slices.Contains(destination.Routes, location.ID)
}
