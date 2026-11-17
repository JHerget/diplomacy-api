package orders

import (
	"diplomacy-api/internal/models"
	"diplomacy-api/internal/utils"
	"fmt"
	"slices"
)

var unitTypeMismatchError = fmt.Errorf("Unit type mismatch.")

func locationError(locationID string) error {
	return fmt.Errorf("Province not found: %s.", locationID)
}
func noUnitError(locationID string) error {
	return fmt.Errorf("No unit at province '%s'.", locationID)
}
func wrongUnit(playerName string) error {
	return fmt.Errorf("The unit is not controlled by '%s'.", playerName)
}
func noRoute(locationID string, destinationID string) error {
	return fmt.Errorf("No route from '%s' to '%s'.", locationID, destinationID)
}

type validator struct {
	BoardMap map[string]models.Province
}

func ValidateCommands(commands models.Commands, board []models.Province) (models.Commands, []models.CommandValidationResult) {
	v := validator{
		BoardMap: make(map[string]models.Province),
	}
	for _, p := range board {
		v.BoardMap[p.ID] = p
	}
	results := make([]models.CommandValidationResult, 0)

	return models.Commands{
		Hold: utils.Filter(commands.Hold, func(c *models.HoldCommand) bool {
			ok, err := v.IsValidHold(c)
			results = append(results, models.CommandValidationResult{
				PlayerName: c.PlayerName,
				Value:      c.Value,
				Valid:      ok,
				Error:      errorMessage(err),
			})

			return ok
		}),
		Move: utils.Filter(commands.Move, func(c *models.MoveCommand) bool {
			ok, err := v.IsValidMove(c, true)
			results = append(results, models.CommandValidationResult{
				PlayerName: c.PlayerName,
				Value:      c.Value,
				Valid:      ok,
				Error:      errorMessage(err),
			})

			return ok
		}),
		Retreat: utils.Filter(commands.Retreat, func(c *models.RetreatCommand) bool {
			ok, err := v.IsValidRetreat(c)
			results = append(results, models.CommandValidationResult{
				PlayerName: c.PlayerName,
				Value:      c.Value,
				Valid:      ok,
				Error:      errorMessage(err),
			})

			return ok
		}),
		Support: utils.Filter(commands.Support, func(c *models.SupportCommand) bool {
			ok, err := v.IsValidSupport(c)
			results = append(results, models.CommandValidationResult{
				PlayerName: c.PlayerName,
				Value:      c.Value,
				Valid:      ok,
				Error:      errorMessage(err),
			})

			return ok
		}),
		Convoy: utils.Filter(commands.Convoy, func(c *models.ConvoyCommand) bool {
			ok, err := v.IsValidConvoy(c)
			results = append(results, models.CommandValidationResult{
				PlayerName: c.PlayerName,
				Value:      c.Value,
				Valid:      ok,
				Error:      errorMessage(err),
			})

			return ok
		}),
		Reinforce: utils.Filter(commands.Reinforce, func(c *models.ReinforceCommand) bool {
			ok, err := v.IsValidReinforce(c)
			results = append(results, models.CommandValidationResult{
				PlayerName: c.PlayerName,
				Value:      c.Value,
				Valid:      ok,
				Error:      errorMessage(err),
			})

			return ok
		}),
		Disband: utils.Filter(commands.Disband, func(c *models.DisbandCommand) bool {
			ok, err := v.IsValidDisband(c)
			results = append(results, models.CommandValidationResult{
				PlayerName: c.PlayerName,
				Value:      c.Value,
				Valid:      ok,
				Error:      errorMessage(err),
			})

			return ok
		}),
	}, results
}

func (v *validator) IsValidHold(hold *models.HoldCommand) (bool, error) {
	loc, ok := v.BoardMap[hold.Location.ID]
	if !ok {
		return false, locationError(hold.Location.ID)
	}

	if ok, err := v.isValidLocation(hold.Location); !ok {
		return false, err
	}

	if loc.Unit == nil {
		return false, noUnitError(loc.ID)
	}
	if loc.Unit.Type != hold.UnitType {
		return false, unitTypeMismatchError
	}
	if loc.Unit.ControlledBy != hold.PlayerName {
		return false, wrongUnit(hold.PlayerName)
	}
	if hold.Location.ID != hold.Destination.ID {
		return false, fmt.Errorf("Hold location and destination must match.")
	}
	if hold.Location.Coast != hold.Destination.Coast {
		return false, fmt.Errorf("Hold location coast and destination coast must match.")
	}

	return true, nil
}

func (v *validator) IsValidMove(move *models.MoveCommand, validatePlayer bool) (bool, error) {
	loc, ok := v.BoardMap[move.Location.ID]
	if !ok {
		return false, locationError(move.Location.ID)
	}

	dest, ok := v.BoardMap[move.Destination.ID]
	if !ok {
		return false, locationError(move.Destination.ID)
	}

	if ok, err := v.isValidLocation(move.Location); !ok {
		return false, err
	}
	if ok, err := v.isValidLocation(move.Destination); !ok {
		return false, err
	}

	if loc.Unit == nil {
		return false, noUnitError(loc.ID)
	}
	if loc.Unit.Type != move.UnitType {
		return false, unitTypeMismatchError
	}
	if validatePlayer && loc.Unit.ControlledBy != move.PlayerName {
		return false, wrongUnit(move.PlayerName)
	}

	validRoute, routeErr := v.isValidRoute(loc, dest, move.Destination)
	possibleConvoy := loc.Unit.IsArmy() && v.hasConvoyRoute(loc, dest)
	if !validRoute && !possibleConvoy {
		return false, routeErr
	}

	return true, nil
}

func (v *validator) IsValidRetreat(retreat *models.RetreatCommand) (bool, error) {
	loc, ok := v.BoardMap[retreat.Location.ID]
	if !ok {
		return false, locationError(retreat.Location.ID)
	}

	dest, ok := v.BoardMap[retreat.Destination.ID]
	if !ok {
		return false, locationError(retreat.Destination.ID)
	}

	if ok, err := v.isValidLocation(retreat.Location); !ok {
		return false, err
	}
	if ok, err := v.isValidLocation(retreat.Destination); !ok {
		return false, err
	}

	if loc.Unit == nil {
		return false, noUnitError(loc.ID)
	}
	if loc.Unit.Type != retreat.UnitType {
		return false, unitTypeMismatchError
	}
	if loc.Unit.ControlledBy != retreat.PlayerName {
		return false, wrongUnit(retreat.PlayerName)
	}

	if ok, err := v.isValidRoute(loc, dest, retreat.Destination); !ok {
		return false, err
	}

	return true, nil
}

func (v *validator) IsValidSupport(support *models.SupportCommand) (bool, error) {
	loc, ok := v.BoardMap[support.Location.ID]
	if !ok {
		return false, locationError(support.Location.ID)
	}

	dest, ok := v.BoardMap[support.Move.Destination.ID]
	if !ok {
		return false, locationError(support.Move.Destination.ID)
	}

	if ok, err := v.IsValidMove(&support.Move, false); !ok {
		return false, err
	}
	if ok, err := v.isValidLocation(support.Location); !ok {
		return false, err
	}

	if loc.Unit == nil {
		return false, noUnitError(loc.ID)
	}
	if loc.Unit.Type != support.UnitType {
		return false, unitTypeMismatchError
	}
	if loc.Unit.ControlledBy != support.PlayerName {
		return false, wrongUnit(support.PlayerName)
	}

	if ok, err := v.isValidRoute(loc, dest, support.Move.Destination); !ok {
		return false, err
	}

	return true, nil
}

func (v *validator) IsValidConvoy(convoy *models.ConvoyCommand) (bool, error) {
	loc, ok := v.BoardMap[convoy.Location.ID]
	if !ok {
		return false, locationError(convoy.Location.ID)
	}

	if ok, err := v.IsValidMove(&convoy.Move, false); !ok {
		return false, err
	}
	if ok, err := v.isValidLocation(convoy.Location); !ok {
		return false, err
	}

	if !loc.IsOcean() {
		return false, fmt.Errorf("Only fleets in ocean provinces can convoy.")
	}
	if loc.Unit == nil {
		return false, noUnitError(loc.ID)
	}
	if loc.Unit.Type != convoy.UnitType {
		return false, unitTypeMismatchError
	}
	if loc.Unit.ControlledBy != convoy.PlayerName {
		return false, wrongUnit(convoy.PlayerName)
	}

	return true, nil
}

func (v *validator) IsValidReinforce(reinforce *models.ReinforceCommand) (bool, error) {
	loc, ok := v.BoardMap[reinforce.Location.ID]
	if !ok {
		return false, locationError(reinforce.Location.ID)
	}

	if ok, err := v.isValidLocation(reinforce.Location); !ok {
		return false, err
	}

	if loc.Unit != nil {
		return false, noUnitError(loc.ID)
	}
	if loc.SupplyCenter == nil {
		return false, fmt.Errorf("This province does not have a supply center: %s.", loc.ID)
	}
	if loc.SupplyCenter.ControlledBy == nil {
		return false, fmt.Errorf("You do not control the supply center.")
	}
	if *loc.SupplyCenter.ControlledBy != reinforce.PlayerName {
		return false, fmt.Errorf("The supply center is not controlled by '%s'.", reinforce.PlayerName)
	}
	if reinforce.UnitType == models.UnitArmy && loc.IsOcean() {
		return false, fmt.Errorf("You cannot reinforce armies to ocean provinces.")
	}
	if reinforce.UnitType == models.UnitFleet && loc.IsInland() {
		return false, fmt.Errorf("You cannot reinforce fleets to inland provinces.")
	}

	return true, nil
}

func (v *validator) IsValidDisband(disband *models.DisbandCommand) (bool, error) {
	loc, ok := v.BoardMap[disband.Location.ID]
	if !ok {
		return false, locationError(disband.Location.ID)
	}

	if ok, err := v.isValidLocation(disband.Location); !ok {
		return false, err
	}

	if loc.Unit == nil {
		return false, noUnitError(loc.ID)
	}
	if loc.Unit.Type != disband.UnitType {
		return false, unitTypeMismatchError
	}
	if loc.Unit.ControlledBy != disband.PlayerName {
		return false, wrongUnit(disband.PlayerName)
	}

	return true, nil
}

func (v *validator) isValidLocation(ref models.LocationReference) (bool, error) {
	loc, ok := v.BoardMap[ref.ID]
	if !ok {
		return false, fmt.Errorf("Province not found: %s", ref.ID)
	}

	coastProvided := ref.Coast != nil
	hasCoasts := len(loc.Coasts) > 0

	if coastProvided && !hasCoasts {
		return false, fmt.Errorf("A coast was provided when not expected.")
	}
	if !coastProvided {
		return true, nil
	}

	_, routesExist := loc.Coasts[*ref.Coast]
	if !routesExist {
		return false, fmt.Errorf("The coast provided does not exist.")
	}

	return true, nil
}

func (v *validator) isValidRoute(location models.Province, destination models.Province, destRef models.LocationReference) (bool, error) {
	if location.Unit == nil {
		return false, noUnitError(location.ID)
	}

	if location.IsOcean() && destination.IsInland() {
		return false, fmt.Errorf("Units cannot move from ocean to inland.")
	}
	if location.IsInland() && destination.IsOcean() {
		return false, fmt.Errorf("Units cannot move from inland to ocean.")
	}
	if location.IsInland() && location.Unit.IsFleet() {
		return false, fmt.Errorf("Fleets cannot be inland.")
	}
	if location.IsOcean() && location.Unit.IsArmy() {
		return false, fmt.Errorf("Armies cannot be in the ocean.")
	}

	requireLocCost := location.Unit != nil && location.Unit.IsFleet() && len(location.Coasts) > 1
	requireDestCoast := len(destination.Coasts) > 1

	if destRef.Coast != nil {
		if !requireDestCoast {
			return false, fmt.Errorf("Coast provided when not expected.")
		}

		coast, ok := destination.Coasts[*destRef.Coast]
		if !ok {
			return false, fmt.Errorf("No routes exist for coast '%s'.", *destRef.Coast)
		}

		return hasRoute(coast.Routes, location.ID, destination.ID)
	}
	if destRef.Coast == nil && requireDestCoast {
		return false, fmt.Errorf("Expecting a coast, none specified.")
	}

	isCoastalRoute := location.Unit.IsFleet() && location.IsCoastal() && destination.IsCoastal()
	if isCoastalRoute {
		unitCoast := location.Unit.Location.Coast

		if unitCoast == nil && requireLocCost {
			return false, fmt.Errorf("The unit at '%s' does not have a coast specified.", location.ID)
		}

		if unitCoast != nil && location.IsCoastal() {
			coast := location.Coasts[*unitCoast]
			return hasRoute(coast.Routes, location.ID, destination.ID)
		}

		coast, ok := destination.Coasts[models.DefaultCoast]
		if !ok {
			return false, fmt.Errorf("No default coastal routes found for '%s'.", destination.ID)
		}

		return hasRoute(coast.Routes, location.ID, destination.ID)
	}

	return hasRoute(destination.Routes, location.ID, destination.ID)
}

func (v *validator) hasConvoyRoute(location, destination models.Province) bool {
	if !location.IsCoastal() || !destination.IsCoastal() {
		return false
	}

	return v.hasConvoyRouteRecursiveHelper(location, destination, make([]string, 0))
}

func (v *validator) hasConvoyRouteRecursiveHelper(location, destination models.Province, visited []string) bool {
	if slices.Contains(visited, location.ID) {
		return false
	}
	if location.ID == destination.ID {
		return false
	}

	for _, province := range location.Routes {
		if province == destination.ID {
			return true
		}

		newLocation, ok := v.BoardMap[province]
		if !ok {
			continue
		}

		if !newLocation.IsOcean() {
			continue
		}
		if v.hasConvoyRouteRecursiveHelper(newLocation, destination, append(visited, location.ID)) {
			return true
		}
	}

	return false
}

func errorMessage(err error) string {
	if err != nil {
		return err.Error()
	}

	return ""
}

func hasRoute(routes []string, locationID string, destinationID string) (bool, error) {
	if slices.Contains(routes, locationID) {
		return true, nil
	}

	return false, noRoute(locationID, destinationID)
}
