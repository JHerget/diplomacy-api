package orders

import (
	"diplomacy-api/internal/models"
	"diplomacy-api/internal/utils"
	"fmt"
	"slices"
)

var unitTypeMismatchError = fmt.Errorf("Unit type mismatch.")

func locationError(locationID string) error {
	return fmt.Errorf("Providence not found: %s.", locationID)
}
func noUnitError(locationID string) error {
	return fmt.Errorf("No unit at providence '%s'.", locationID)
}
func wrongUnit(playerName string) error {
	return fmt.Errorf("The unit is not controlled by '%s'.", playerName)
}

type validator struct {
	BoardMap map[string]models.Providence
}

func ValidateCommands(commands models.Commands, board []models.Providence) (models.Commands, []models.CommandValidationResult) {
	v := validator{
		BoardMap: make(map[string]models.Providence),
	}
	for _, p := range board {
		v.BoardMap[p.ID] = p
	}
	results := make([]models.CommandValidationResult, 0)

	return models.Commands{
		Hold: utils.Filter(commands.Hold, func(c *models.MoveCommand) bool {
			ok, err := v.IsValidHold(c)
			results = append(results, models.CommandValidationResult{
				PlayerName: c.PlayerName,
				Value:      c.Value,
				Valid:      ok,
				Error:      err.Error(),
			})

			return ok
		}),
		Move: utils.Filter(commands.Move, func(c *models.MoveCommand) bool {
			ok, err := v.IsValidMove(c, true)
			results = append(results, models.CommandValidationResult{
				PlayerName: c.PlayerName,
				Value:      c.Value,
				Valid:      ok,
				Error:      err.Error(),
			})

			return ok
		}),
		Retreat: utils.Filter(commands.Retreat, func(c *models.MoveCommand) bool {
			ok, err := v.IsValidRetreat(c)
			results = append(results, models.CommandValidationResult{
				PlayerName: c.PlayerName,
				Value:      c.Value,
				Valid:      ok,
				Error:      err.Error(),
			})

			return ok
		}),
		Support: utils.Filter(commands.Support, func(c *models.SupportCommand) bool {
			ok, err := v.IsValidSupport(c)
			results = append(results, models.CommandValidationResult{
				PlayerName: c.PlayerName,
				Value:      c.Value,
				Valid:      ok,
				Error:      err.Error(),
			})

			return ok
		}),
		Convoy: utils.Filter(commands.Convoy, func(c *models.SupportCommand) bool {
			ok, err := v.IsValidConvoy(c)
			results = append(results, models.CommandValidationResult{
				PlayerName: c.PlayerName,
				Value:      c.Value,
				Valid:      ok,
				Error:      err.Error(),
			})

			return ok
		}),
		Reinforce: utils.Filter(commands.Reinforce, func(c *models.AdjustCommand) bool {
			ok, err := v.IsValidReinforce(c)
			results = append(results, models.CommandValidationResult{
				PlayerName: c.PlayerName,
				Value:      c.Value,
				Valid:      ok,
				Error:      err.Error(),
			})

			return ok
		}),
		Disband: utils.Filter(commands.Disband, func(c *models.AdjustCommand) bool {
			ok, err := v.IsValidDisband(c)
			results = append(results, models.CommandValidationResult{
				PlayerName: c.PlayerName,
				Value:      c.Value,
				Valid:      ok,
				Error:      err.Error(),
			})

			return ok
		}),
	}, results
}

func (v validator) IsValidHold(hold *models.MoveCommand) (bool, error) {
	loc, ok := v.BoardMap[hold.Location.ProperID()]
	if !ok {
		return false, locationError(hold.Location.ProperID())
	}

	if ok, err := v.IsValidLocation(hold.Location); !ok {
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

func (v validator) IsValidMove(move *models.MoveCommand, validatePlayer bool) (bool, error) {
	loc, ok := v.BoardMap[move.Location.ProperID()]
	if !ok {
		return false, locationError(move.Location.ProperID())
	}

	dest, ok := v.BoardMap[move.Destination.ProperID()]
	if !ok {
		return false, locationError(move.Destination.ProperID())
	}

	if ok, err := v.IsValidLocation(move.Location); !ok {
		return false, err
	}
	if ok, err := v.IsValidLocation(move.Destination); !ok {
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

	if ok, err := v.IsValidRoute(loc, dest, move.Destination); !ok {
		return false, err
	}

	return true, nil
}

func (v validator) IsValidRetreat(retreat *models.MoveCommand) (bool, error) {
	loc, ok := v.BoardMap[retreat.Location.ProperID()]
	if !ok {
		return false, locationError(retreat.Location.ProperID())
	}

	dest, ok := v.BoardMap[retreat.Destination.ProperID()]
	if !ok {
		return false, locationError(retreat.Destination.ProperID())
	}

	if ok, err := v.IsValidLocation(retreat.Location); !ok {
		return false, err
	}
	if ok, err := v.IsValidLocation(retreat.Destination); !ok {
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

	if ok, err := v.IsValidRoute(loc, dest, retreat.Destination); !ok {
		return false, err
	}

	return true, nil
}

func (v validator) IsValidSupport(support *models.SupportCommand) (bool, error) {
	loc, ok := v.BoardMap[support.Location.ProperID()]
	if !ok {
		return false, locationError(support.Location.ProperID())
	}

	dest, ok := v.BoardMap[support.Move.Destination.ProperID()]
	if !ok {
		return false, locationError(support.Move.Destination.ProperID())
	}

	if ok, err := v.IsValidMove(&support.Move, false); !ok {
		return false, err
	}
	if ok, err := v.IsValidLocation(support.Location); !ok {
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

	if ok, err := v.IsValidRoute(loc, dest, support.Move.Destination); !ok {
		return false, err
	}

	return true, nil
}

func (v validator) IsValidConvoy(convoy *models.SupportCommand) (bool, error) {
	loc, ok := v.BoardMap[convoy.Location.ProperID()]
	if !ok {
		return false, locationError(convoy.Location.ProperID())
	}

	dest, ok := v.BoardMap[convoy.Move.Destination.ProperID()]
	if !ok {
		return false, locationError(convoy.Move.Destination.ProperID())
	}

	if ok, err := v.IsValidMove(&convoy.Move, false); !ok {
		return false, err
	}
	if ok, err := v.IsValidLocation(convoy.Location); !ok {
		return false, err
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

	if ok, err := v.IsValidRoute(loc, dest, convoy.Move.Destination); !ok {
		return false, err
	}

	return true, nil
}

func (v validator) IsValidReinforce(reinforce *models.AdjustCommand) (bool, error) {
	loc, ok := v.BoardMap[reinforce.Location.ProperID()]
	if !ok {
		return false, locationError(reinforce.Location.ProperID())
	}

	if ok, err := v.IsValidLocation(reinforce.Location); !ok {
		return false, err
	}

	if loc.Unit != nil {
		return false, noUnitError(loc.ID)
	}
	if loc.SupplyCenter == nil {
		return false, fmt.Errorf("This providence does not have a supply center: %s.", loc.ID)
	}
	if loc.SupplyCenter.ControlledBy == nil {
		return false, fmt.Errorf("You do not control the supply center.")
	}
	if *loc.SupplyCenter.ControlledBy != reinforce.PlayerName {
		return false, fmt.Errorf("The supply center is not controlled by '%s'.", reinforce.PlayerName)
	}
	if reinforce.UnitType == models.UnitArmy && loc.Type == models.ProvidenceOcean {
		return false, fmt.Errorf("You cannot reinforce armies to ocean providences.")
	}
	if reinforce.UnitType == models.UnitFleet && loc.Type == models.ProvidenceInland {
		return false, fmt.Errorf("You cannot reinforce fleets to inland providences.")
	}

	return true, nil
}

func (v validator) IsValidDisband(disband *models.AdjustCommand) (bool, error) {
	loc, ok := v.BoardMap[disband.Location.ProperID()]
	if !ok {
		return false, locationError(disband.Location.ProperID())
	}

	if ok, err := v.IsValidLocation(disband.Location); !ok {
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

func (v validator) IsValidLocation(ref models.LocationReference) (bool, error) {
	loc, ok := v.BoardMap[ref.ProperID()]
	if !ok {
		return false, fmt.Errorf("Providence not found: %s", ref.ProperID())
	}

	coastProvided := ref.Coast != nil
	hasCoastalRoutes := len(loc.CoastalRoutes) > 0

	if !coastProvided && hasCoastalRoutes {
		return false, fmt.Errorf("Expecting a coast.")
	}
	if coastProvided && !hasCoastalRoutes {
		return false, fmt.Errorf("A coast was provided when not expected.")
	}

	if !coastProvided {
		return true, nil
	}

	_, routesExist := loc.CoastalRoutes[*ref.Coast]
	if !routesExist {
		return false, fmt.Errorf("The coast provided does not exist.")
	}

	return true, nil
}

func (v validator) IsValidRoute(location models.Providence, destination models.Providence, destRef models.LocationReference) (bool, error) {
	if location.Unit == nil {
		return false, noUnitError(location.ID)
	}

	if location.Type == models.ProvidenceOcean && destination.Type == models.ProvidenceInland {
		return false, fmt.Errorf("Units cannot move from ocean to inland.")
	}
	if location.Type == models.ProvidenceInland && destination.Type == models.ProvidenceOcean {
		return false, fmt.Errorf("Units cannot move from inland to ocean.")
	}
	if location.Type == models.ProvidenceInland && location.Unit.Type == models.UnitFleet {
		return false, fmt.Errorf("Fleets cannot be inland.")
	}
	if location.Type == models.ProvidenceOcean && location.Unit.Type == models.UnitArmy {
		return false, fmt.Errorf("Armies cannot be in the ocean.")
	}

	hasCoastalRoutes := len(destination.CoastalRoutes) > 0

	if destRef.Coast != nil {
		if !hasCoastalRoutes {
			return slices.Contains(destination.Routes, location.ID), nil
		}

		routes, ok := destination.CoastalRoutes[*destRef.Coast]
		if !ok {
			return false, fmt.Errorf("No routes exist for coast '%s'.", *destRef.Coast)
		}

		return slices.Contains(routes, location.ID), nil
	}
	if hasCoastalRoutes {
		return false, fmt.Errorf("Expecting a coast, none specified.")
	}

	return slices.Contains(destination.Routes, location.ID), nil
}
