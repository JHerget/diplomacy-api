package tests

import (
	"encoding/json"
	"fmt"
	"maps"
	"os"
	"reflect"
	"runtime/debug"
	"slices"
	"strings"
	"testing"

	"diplomacy-api/internal/board"
	"diplomacy-api/internal/models"
	"diplomacy-api/internal/orders"
)

type fixture struct {
	ExpectedOutcome string            `json:"expectedOutcome"`
	InputBoard      []models.Province `json:"inputBoard"`
	OutputBoard     []models.Province `json:"outputBoard"`
	Orders          []models.Order    `json:"orders"`
}

type disposition string

const (
	accepted         disposition = "accepted"
	invalid          disposition = "invalid"
	duplicate        disposition = "ignored duplicate"
	unmatched        disposition = "support does not match the unit's order"
	standoff         disposition = "retreat blocked by standoff"
	dormant          disposition = "unit was not dislodged"
	occupiedRefuge   disposition = "occupied retreat destination"
	dislodgerOrigin  disposition = "dislodger origin"
	springAdjustment disposition = "out of season"
	noBuildDue       disposition = "no build due"
	noRemovalDue     disposition = "no removal due"
	noConvoy         disposition = "missing convoy route"
)

func invalidBecause(reason string) disposition {
	return disposition("invalid: " + strings.ToLower(reason))
}

func supportCutFrom(province string) disposition {
	return disposition("support was cut from '" + province + "'.")
}

type expectedCommand struct {
	value   string
	outcome disposition
}

func checkMovement(t *testing.T, id int, want []expectedCommand) {
	t.Helper()
	f, players := loadFixture(t, id)
	t.Log(f.ExpectedOutcome)

	checkSpringFixture(t, f, players, want)
}

func checkSpringFixture(t *testing.T, f fixture, players []models.Player, want []expectedCommand) {
	t.Helper()
	submitted := orders.ListCommands(f.Orders)
	got, results := resolveSpring(t, f, players)
	checkBoard(t, got, f.OutputBoard)
	checkDispositions(t, submitted, results, want)
}

func loadFixture(t *testing.T, id int) (fixture, []models.Player) {
	t.Helper()
	var f fixture
	readJSON(t, fmt.Sprintf("maps/%d.json", id), &f)

	var original models.Map
	readJSON(t, "../assets/original-diplomacy.json", &original)
	if len(f.InputBoard) != len(original.Provinces) || len(f.OutputBoard) != len(original.Provinces) {
		t.Fatalf("both fixture boards must contain all %d provinces", len(original.Provinces))
	}
	for i := range original.Players {
		original.Players[i].IsPlaying = true
	}
	return f, original.Players
}

func addMarseillesDislodgement(t *testing.T, f *fixture) {
	t.Helper()
	f.Orders = append(f.Orders,
		models.Order{PlayerName: "Germany", Value: "A bur - mar, A pie S A bur - mar"},
		models.Order{PlayerName: "England", Value: "A R mar - spa"},
	)
	placeArmy(t, f.InputBoard, "Germany", "bur", "bur")
	placeArmy(t, f.InputBoard, "Germany", "pie", "pie")
	placeArmy(t, f.InputBoard, "England", "mar", "mar")

	placeArmy(t, f.OutputBoard, "Germany", "bur", "mar")
	placeArmy(t, f.OutputBoard, "Germany", "pie", "pie")
}

func placeArmy(t *testing.T, provinces []models.Province, player, origin, destination string) {
	t.Helper()
	placeUnit(t, provinces, models.UnitArmy, player, origin, destination)
}

func placeUnit(t *testing.T, provinces []models.Province, unitType models.UnitType, player, origin, destination string) {
	t.Helper()
	for i := range provinces {
		if provinces[i].ID == destination {
			provinces[i].Unit = &models.Unit{
				ID:           "u-" + origin,
				ControlledBy: player,
				Type:         unitType,
				Location:     models.LocationReference{ID: destination},
			}
			return
		}
	}
	t.Fatalf("missing province %q", destination)
}

func readJSON(t *testing.T, path string, target any) {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(data, target); err != nil {
		t.Fatalf("%s: %v", path, err)
	}
}

func resolveSpring(t *testing.T, f fixture, players []models.Player) ([]models.Province, []models.CommandValidationResult) {
	t.Helper()
	game := models.Game{
		Board: f.InputBoard, Players: players, PlayerCount: len(players),
		Turns: []models.Turn{{TurnNumber: 1, Orders: f.Orders}},
	}
	return resolveTurn(t, &game)
}

func resolveTurn(t *testing.T, game *models.Game) ([]models.Province, []models.CommandValidationResult) {
	t.Helper()
	// A resolver panic fails this fixture while letting the other tests run.
	defer func() {
		if problem := recover(); problem != nil {
			t.Fatalf("ApplyTurn panicked: %v\n%s", problem, debug.Stack())
		}
	}()
	got, results, err := board.ApplyTurn(game.Board, &game.Turns[len(game.Turns)-1])
	if err != nil {
		t.Fatalf("ApplyTurn: %v", err)
	}
	return got, results
}

func checkBoard(t *testing.T, got, want []models.Province) {
	t.Helper()
	byID := make(map[string]models.Province, len(got))
	unitIDs := map[string]bool{}
	for _, province := range got {
		if u := province.Unit; u != nil {
			if u.ID == "" || unitIDs[u.ID] || u.Location.ID != province.ID {
				t.Errorf("board %s: invalid unit identity or nested location: %+v", province.ID, u)
			}
			unitIDs[u.ID] = true
		}
		if _, exists := byID[province.ID]; exists {
			t.Errorf("board %s: duplicate province", province.ID)
		}
		byID[province.ID] = province
	}
	// Report differences in fixture order, regardless of ApplyTurn's map iteration order.
	seen := make(map[string]bool, len(want))
	for _, expected := range want {
		if seen[expected.ID] {
			t.Fatalf("fixture contains duplicate province %q", expected.ID)
		}
		seen[expected.ID] = true
		actual, ok := byID[expected.ID]
		if !ok {
			t.Errorf("board %s: missing province", expected.ID)
			continue
		}
		delete(byID, expected.ID)
		if !reflect.DeepEqual(actual, expected) {
			reportProvinceDifferences(t, expected, actual)
		}
	}
	for _, id := range slices.Sorted(maps.Keys(byID)) {
		t.Errorf("board %s: unexpected province", id)
	}
}

func reportProvinceDifferences(t *testing.T, want, got models.Province) {
	t.Helper()
	expected, actual := reflect.ValueOf(want), reflect.ValueOf(got)
	for i := 0; i < expected.NumField(); i++ {
		wantField, gotField := expected.Field(i).Interface(), actual.Field(i).Interface()
		if reflect.DeepEqual(wantField, gotField) {
			continue
		}
		field := expected.Type().Field(i).Tag.Get("json")
		t.Errorf("board %s (%s) %s: want %s, got %s", want.Name, want.ID, field, boardValue(wantField), boardValue(gotField))
	}
}

func boardValue(value any) string {
	if unit, ok := value.(*models.Unit); ok {
		if unit == nil {
			return "empty"
		}
		location := unit.Location.ID
		if unit.Location.Coast != nil {
			location += "-" + string(*unit.Location.Coast)
		}
		return fmt.Sprintf("%s %s %s at %s", unit.ControlledBy, unit.Type, unit.ID, location)
	}
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Sprintf("%+v", value)
	}
	return string(data)
}

func checkDispositions(t *testing.T, submitted, results []models.CommandValidationResult, want []expectedCommand) {
	t.Helper()
	if len(submitted) != len(want) {
		t.Fatalf("parsed commands: want %d, got %d", len(want), len(submitted))
	}
	// Results are grouped by command kind. Identical commands have the same kind
	// and retain their submission order within the group (case 625).
	type commandKey struct{ player, value string }
	byCommand := make(map[commandKey][]models.CommandValidationResult, len(results))
	for _, result := range results {
		key := commandKey{result.PlayerName, result.Value}
		byCommand[key] = append(byCommand[key], result)
	}
	for i, raw := range submitted {
		if !strings.EqualFold(raw.Value, want[i].value) {
			t.Errorf("order %d %s: want %q, got %q", i+1, raw.PlayerName, want[i].value, raw.Value)
		}
		key := commandKey{raw.PlayerName, raw.Value}
		result := raw
		if matches := byCommand[key]; len(matches) > 0 {
			result = matches[0]
			if len(matches) == 1 {
				delete(byCommand, key)
			} else {
				byCommand[key] = matches[1:]
			}
		} else if raw.Valid {
			t.Errorf("order %d %s %q: missing resolution result", i+1, key.player, key.value)
			continue
		}
		got := accepted
		if !result.Valid {
			got = invalid
			reason := strings.ToLower(result.Error)
			if reason == "" {
				t.Errorf("order %d %s %q: missing rejection reason", i+1, key.player, key.value)
			}
			switch {
			case strings.Contains(reason, "already applied a command"), strings.Contains(reason, "duplicate"):
				got = duplicate
			case strings.HasPrefix(reason, "support was cut from "):
				got = disposition(reason)
			case reason == "could not find the move to support.":
				got = unmatched
			case reason == "unit was not dislodged.":
				got = dormant
			case reason == "a unit occupies the destination.":
				got = occupiedRefuge
			case reason == "cannot retreat to where the attack originated.":
				got = dislodgerOrigin
			case reason == "it is not a fall turn.":
				got = springAdjustment
			case reason == "not enough supply centers to support a new unit.":
				got = noBuildDue
			case reason == "the number of supply centers exceeds or matches the number of units.":
				got = noRemovalDue
			case strings.HasPrefix(reason, "no convoy route from "):
				got = noConvoy
			case reason == "cannot retreat to where a standoff occurred.":
				got = standoff
			}
		}
		if got == invalid && want[i].outcome == invalidBecause(result.Error) {
			continue
		}
		// A missing route is also a valid explanation in general illegal-move cases.
		if got != want[i].outcome && !(got == noConvoy && want[i].outcome == invalid) {
			reason := ""
			if result.Error != "" {
				reason = " (" + result.Error + ")"
			}
			t.Errorf("order %d %s %q: want %s, got %s%s", i+1, raw.PlayerName, want[i].value, want[i].outcome, got, reason)
		}
	}
	for key := range byCommand {
		t.Errorf("order %s %q: result for unsubmitted command", key.player, key.value)
	}
}
