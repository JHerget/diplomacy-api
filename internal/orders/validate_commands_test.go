package orders

import (
	"testing"

	"diplomacy-api/internal/models"
)

func TestArmyCannotEnterOceanFromCoast(t *testing.T) {
	checkRoute(t, models.UnitArmy, models.ProvinceCoastal, models.ProvinceOcean, false)
}

func TestFleetCannotEnterInlandFromCoast(t *testing.T) {
	checkRoute(t, models.UnitFleet, models.ProvinceCoastal, models.ProvinceInland, false)
}

func TestArmyCanEnterInlandFromCoast(t *testing.T) {
	checkRoute(t, models.UnitArmy, models.ProvinceCoastal, models.ProvinceInland, true)
}

func TestArmyCanEnterCoastFromInland(t *testing.T) {
	checkRoute(t, models.UnitArmy, models.ProvinceInland, models.ProvinceCoastal, true)
}

func TestFleetCanEnterCoastFromCoast(t *testing.T) {
	checkRoute(t, models.UnitFleet, models.ProvinceCoastal, models.ProvinceCoastal, true)
}

func TestFleetCanEnterOceanFromCoast(t *testing.T) {
	checkRoute(t, models.UnitFleet, models.ProvinceCoastal, models.ProvinceOcean, true)
}

func TestArmyCanEnterSplitCoastByLand(t *testing.T) {
	origin := models.Province{
		ID:   "mar",
		Type: models.ProvinceCoastal,
		Unit: &models.Unit{
			Type:     models.UnitArmy,
			Location: models.LocationReference{ID: "mar"},
		},
	}
	destination := models.Province{
		ID:     "spa",
		Type:   models.ProvinceCoastal,
		Routes: []string{"mar"},
		Coasts: map[models.Coast]models.ProvinceCoast{
			models.NorthCoast: {Routes: []string{"gas", "mid"}},
			models.SouthCoast: {Routes: []string{"mar", "wes"}},
		},
	}
	v := validator{}
	ref := models.LocationReference{ID: "spa"}
	valid, err := v.isValidRoute(origin, destination, ref)
	if !valid || err != nil {
		t.Fatalf("adjacent army route = (%v, %v), want (true, nil)", valid, err)
	}

	// Coastal adjacency must not substitute for land adjacency for an army.
	destination.Routes = nil
	valid, err = v.isValidRoute(origin, destination, ref)
	wantError := noRoute("mar", "spa").Error()
	if valid || err == nil || err.Error() != wantError {
		t.Fatalf("nonadjacent army route = (%v, %v), want (false, %q)", valid, err, wantError)
	}
}

func TestSplitCoastFleetDepartureToOcean(t *testing.T) {
	tests := []struct {
		name        string
		coast       models.Coast
		destination string
		wantError   string
	}{
		{"north to western mediterranean", models.NorthCoast, "wes", "No route from 'spa' to 'wes'."},
		{"south to western mediterranean", models.SouthCoast, "wes", ""},
		{"north to mid atlantic", models.NorthCoast, "mid", ""},
		{"south to mid atlantic", models.SouthCoast, "mid", ""},
		{"missing coast", "", "wes", "The unit at 'spa' does not have a coast specified."},
		{"unknown coast", models.EastCoast, "wes", "No routes exist for coast 'ec'."},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var unitCoast *models.Coast
			if tt.coast != "" {
				coast := tt.coast
				unitCoast = &coast
			}
			origin := models.Province{
				ID:   "spa",
				Type: models.ProvinceCoastal,
				Unit: &models.Unit{
					Type:     models.UnitFleet,
					Location: models.LocationReference{ID: "spa", Coast: unitCoast},
				},
				Coasts: map[models.Coast]models.ProvinceCoast{
					models.NorthCoast: {Routes: []string{"por", "mid", "gas"}},
					models.SouthCoast: {Routes: []string{"por", "mid", "wes", "gol", "mar"}},
				},
			}
			// General adjacency must not override the occupied coast's routes.
			destination := models.Province{
				ID: tt.destination, Type: models.ProvinceOcean, Routes: []string{"spa"},
			}
			v := validator{}
			valid, err := v.isValidRoute(origin, destination, models.LocationReference{ID: destination.ID})
			wantValid := tt.wantError == ""
			if valid != wantValid || (err == nil) != wantValid {
				t.Fatalf("isValidRoute() = (%v, %v), want valid=%v, error=%q", valid, err, wantValid, tt.wantError)
			}
			if err != nil && err.Error() != tt.wantError {
				t.Fatalf("error = %q, want %q", err.Error(), tt.wantError)
			}
		})
	}
}

func TestSupportIntoSplitCoast(t *testing.T) {
	tests := []struct {
		name      string
		change    func(*validator, *models.SupportCommand)
		wantError string
	}{
		{name: "different reachable coasts"},
		{
			name: "unknown supported origin",
			change: func(v *validator, support *models.SupportCommand) {
				support.Move.Location.ID = "xxx"
			},
			wantError: "Province not found: xxx.",
		},
		{
			name: "missing supported unit",
			change: func(v *validator, support *models.SupportCommand) {
				origin := v.BoardMap["gas"]
				origin.Unit = nil
				v.BoardMap["gas"] = origin
			},
			wantError: "No unit at province 'gas'.",
		},
		{
			name: "wrong supported unit type",
			change: func(v *validator, support *models.SupportCommand) {
				v.BoardMap["gas"].Unit.Type = models.UnitArmy
			},
			wantError: "Unit type mismatch.",
		},
		{
			name: "supported fleet cannot reach any coast",
			change: func(v *validator, support *models.SupportCommand) {
				v.BoardMap["spa"].Coasts[models.NorthCoast] = models.ProvinceCoast{}
			},
			wantError: "No route from 'gas' to 'spa'.",
		},
		{
			name: "supporter cannot reach any coast",
			change: func(v *validator, support *models.SupportCommand) {
				v.BoardMap["spa"].Coasts[models.SouthCoast] = models.ProvinceCoast{}
			},
			wantError: "No route from 'gol' to 'spa'.",
		},
		{
			name: "supported fleet occupies unreachable coast",
			change: func(v *validator, support *models.SupportCommand) {
				origin := v.BoardMap["gas"]
				coast := models.SouthCoast
				origin.Unit.Location.Coast = &coast
				origin.Coasts = map[models.Coast]models.ProvinceCoast{
					models.NorthCoast: {Routes: []string{"spa"}},
					models.SouthCoast: {Routes: []string{"other"}},
				}
				v.BoardMap["gas"] = origin
				support.Move.Location.Coast = &coast
			},
			wantError: "No route from 'gas' to 'spa'.",
		},
		{
			name: "supporter occupies unreachable coast",
			change: func(v *validator, support *models.SupportCommand) {
				origin := v.BoardMap["gol"]
				coast := models.SouthCoast
				origin.Type = models.ProvinceCoastal
				origin.Unit.Location.Coast = &coast
				origin.Coasts = map[models.Coast]models.ProvinceCoast{
					models.NorthCoast: {Routes: []string{"spa"}},
					models.SouthCoast: {Routes: []string{"other"}},
				}
				v.BoardMap["gol"] = origin
				support.Location.Coast = &coast
			},
			wantError: "No route from 'gol' to 'spa'.",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := validator{BoardMap: map[string]models.Province{
				"gas": {
					ID: "gas", Type: models.ProvinceCoastal,
					Unit: &models.Unit{
						Type: models.UnitFleet, ControlledBy: "France",
						Location: models.LocationReference{ID: "gas"},
					},
				},
				"gol": {
					ID: "gol", Type: models.ProvinceOcean,
					Unit: &models.Unit{
						Type: models.UnitFleet, ControlledBy: "Italy",
						Location: models.LocationReference{ID: "gol"},
					},
				},
				"spa": {
					ID: "spa", Type: models.ProvinceCoastal,
					Routes: []string{"gas", "gol"},
					Coasts: map[models.Coast]models.ProvinceCoast{
						models.NorthCoast: {Routes: []string{"gas"}},
						models.SouthCoast: {Routes: []string{"gol"}},
					},
				},
			}}
			support := models.SupportCommand{
				PlayerName: "Italy", UnitType: models.UnitFleet,
				Location: models.LocationReference{ID: "gol"},
				Move: models.MoveCommand{
					UnitType:    models.UnitFleet,
					Location:    models.LocationReference{ID: "gas"},
					Destination: models.LocationReference{ID: "spa"},
				},
			}
			if tt.change != nil {
				tt.change(&v, &support)
			}
			valid, err := v.IsValidSupport(&support)
			if valid != (tt.wantError == "") || errorMessage(err) != tt.wantError {
				t.Fatalf("IsValidSupport() = (%v, %v), want valid=%v, error=%q", valid, err, tt.wantError == "", tt.wantError)
			}
			if support.Move.Destination.ID != "spa" || support.Move.Destination.Coast != nil {
				t.Fatalf("support destination was modified: %+v", support.Move.Destination)
			}
		})
	}
}

func checkRoute(t *testing.T, unit models.UnitType, from, to models.ProvinceType, wantValid bool) {
	t.Helper()
	origin := models.Province{
		ID:   "origin",
		Type: from,
		Unit: &models.Unit{Type: unit},
	}
	// Connect the provinces so terrain alone determines whether the route is legal.
	destination := models.Province{
		ID:     "destination",
		Type:   to,
		Routes: []string{origin.ID},
	}
	if destination.IsCoastal() {
		destination.Coasts = map[models.Coast]models.ProvinceCoast{
			models.DefaultCoast: {Routes: []string{origin.ID}},
		}
	}
	v := validator{}
	valid, err := v.isValidRoute(origin, destination, models.LocationReference{ID: destination.ID})
	if valid != wantValid || (err == nil) != wantValid {
		t.Fatalf("isValidRoute() = (%v, %v), want valid=%v", valid, err, wantValid)
	}
}
