package tests

import (
	"slices"
	"strings"
	"testing"

	"diplomacy-api/internal/models"
	"diplomacy-api/internal/orders"
)

func checkTurnCase(t *testing.T, id, turnNumber int, commands []expectedCommand) {
	t.Helper()
	f, players := loadFixture(t, id)
	t.Log(f.ExpectedOutcome)
	game := models.Game{Players: players, PlayerCount: len(players), Board: clonePosition(f.InputBoard)}

	game.Turns = append(game.Turns, models.Turn{TurnNumber: turnNumber, Orders: f.Orders})
	got, results := resolveTurn(t, &game)
	checkTurnBoard(t, id, got, f)
	checkDispositions(t, orders.ListCommands(f.Orders), results, commands)
}

// The resolver updates units and supply centers through their pointers.
func clonePosition(board []models.Province) []models.Province {
	cloned := slices.Clone(board)
	for i, province := range board {
		if province.Unit != nil {
			u := *province.Unit
			cloned[i].Unit = &u
		}
		if province.SupplyCenter != nil {
			center := *province.SupplyCenter
			cloned[i].SupplyCenter = &center
		}
	}
	return cloned
}

// Bind only the two values the rules leave open: the ID of a newly built
// unit and the identity of a randomly removed unit. All other data, including
// the original IDs and locations of survivors, is compared exactly.
func checkTurnBoard(t *testing.T, id int, got []models.Province, f fixture) {
	t.Helper()
	want := slices.Clone(f.OutputBoard)
	inputIDs := map[string]bool{}
	inputUnits := map[string]*models.Unit{}
	for _, p := range f.InputBoard {
		if p.Unit != nil {
			inputIDs[p.Unit.ID] = true
			inputUnits[p.ID] = p.Unit
		}
	}
	actualUnits := map[string]*models.Unit{}
	for _, p := range got {
		actualUnits[p.ID] = p.Unit
	}

	eligible := randomRemovalSites(id)
	for i, p := range want {
		actual := actualUnits[p.ID]
		switch {
		case slices.Contains(eligible, p.ID):
			want[i].Unit = nil
			if actual != nil {
				want[i].Unit = inputUnits[p.ID]
			}
		case p.Unit != nil && strings.HasPrefix(p.Unit.ID, "new-") && actual != nil:
			if inputIDs[actual.ID] {
				t.Errorf("build at %s reused input ID %q", p.ID, actual.ID)
			}
			unit := *p.Unit
			unit.ID = actual.ID
			want[i].Unit = &unit
		}
	}
	if len(eligible) > 0 {
		power := inputUnits[eligible[0]].ControlledBy
		units, centers := 0, 0
		for _, p := range got {
			if p.Unit != nil && p.Unit.ControlledBy == power {
				units++
			}
			if p.SupplyCenter != nil && p.SupplyCenter.ControlledBy != nil && *p.SupplyCenter.ControlledBy == power {
				centers++
			}
		}
		if units != 1 || units != centers {
			t.Errorf("%s: want one survivor for one center, got %d units for %d centers", power, units, centers)
		}
	}
	checkBoard(t, got, want)
}

func randomRemovalSites(id int) []string {
	switch id {
	case 631, 639:
		return []string{"pic", "bur"}
	case 704:
		return []string{"gal", "bre"}
	case 705:
		return []string{"mid", "pie"}
	case 706:
		return []string{"pie", "spa"}
	case 707:
		return []string{"tyr", "pic"}
	default:
		return nil
	}
}
