package tests

import (
	"slices"
	"strings"
	"testing"

	"diplomacy-api/internal/models"
)

func Test501_ConditionFalse(t *testing.T) {
	checkTurnCase(t, 501, 1, []expectedCommand{
		{"A pic - bel", accepted},
		{"A R bel - hol", dormant},
	})
}

func Test502_ConditionTrue(t *testing.T) {
	checkTurnCase(t, 502, 1, []expectedCommand{
		{"A pic - bel", accepted},
		{"A bur S A pic - bel", accepted},
		{"A R bel - hol", accepted},
	})
}

func Test503_NoRetreatSubmitted(t *testing.T) {
	checkTurnCase(t, 503, 1, []expectedCommand{
		{"A pic - bel", accepted},
		{"A bur S A pic - bel", accepted},
	})
}

func Test504_MovementSubmittedBeforeRetreat(t *testing.T) {
	checkTurnCase(t, 504, 1, []expectedCommand{
		{"A bel - ruh", accepted},
		{"A R bel - hol", accepted},
		{"A pic - bel", accepted},
		{"A bur S A pic - bel", accepted},
	})
}

func Test505_RetreatSubmittedBeforeMovement(t *testing.T) {
	checkTurnCase(t, 505, 1, []expectedCommand{
		{"A R bel - hol", accepted},
		{"A bel - ruh", accepted},
		{"A pic - bel", accepted},
		{"A bur S A pic - bel", accepted},
	})
}

func Test506_RetreatSubmittedFirstAndUnitEscapesByMoving(t *testing.T) {
	checkTurnCase(t, 506, 1, []expectedCommand{
		{"A R bel - hol", dormant},
		{"A bel - ruh", accepted},
		{"A pic - bel", accepted},
		{"A bur S A pic - bel", accepted},
	})
}

func Test507_OccupiedRefuge(t *testing.T) {
	checkTurnCase(t, 507, 1, []expectedCommand{
		{"A pic - bel", accepted},
		{"A bur S A pic - bel", accepted},
		{"A R bel - hol", occupiedRefuge},
	})
}

func Test508_NewlyVacatedRefuge(t *testing.T) {
	checkTurnCase(t, 508, 1, []expectedCommand{
		{"A pic - bel", accepted},
		{"A bur S A pic - bel", accepted},
		{"A hol - kie", accepted},
		{"A R bel - hol", accepted},
	})
}

func Test509_NewlyOccupiedRefuge(t *testing.T) {
	checkTurnCase(t, 509, 1, []expectedCommand{
		{"A pic - bel", accepted},
		{"A bur S A pic - bel", accepted},
		{"F hel - hol", accepted},
		{"A R bel - hol", occupiedRefuge},
	})
}

func Test510_EmptyDislodgerOriginForbidden(t *testing.T) {
	checkTurnCase(t, 510, 1, []expectedCommand{
		{"A pic - bel", accepted},
		{"A bur S A pic - bel", accepted},
		{"A R bel - pic", dislodgerOrigin},
	})
}

func Test511_EmptyStandoffProvinceForbidden(t *testing.T) {
	checkTurnCase(t, 511, 1, []expectedCommand{
		{"A pic - bel", accepted},
		{"A bur S A pic - bel", accepted},
		{"F nth - hol", accepted},
		{"A kie - hol", accepted},
		{"A R bel - hol", standoff},
	})
}

func Test512_NonadjacentRefuge(t *testing.T) {
	checkTurnCase(t, 512, 1, []expectedCommand{
		{"A pic - bel", accepted},
		{"A bur S A pic - bel", accepted},
		{"A R bel - ber", invalidBecause("No route from 'bel' to 'ber'.")},
	})
}

func Test513_RetreatCannotUseAConvoy(t *testing.T) {
	checkTurnCase(t, 513, 1, []expectedCommand{
		{"A pic - bel", accepted},
		{"A bur S A pic - bel", accepted},
		{"F nth C A bel - lon", accepted},
		{"A R bel - lon", invalid},
	})
}

func Test514_FleetRetreatsToReachableCoast(t *testing.T) {
	checkTurnCase(t, 514, 1, []expectedCommand{
		{"F mid - gas", accepted},
		{"A par S F mid - gas", accepted},
		{"F R gas - spa-nc", accepted},
	})
}

func Test515_FleetRetreatsToUnreachableCoast(t *testing.T) {
	checkTurnCase(t, 515, 1, []expectedCommand{
		{"F mid - gas", accepted},
		{"A par S F mid - gas", accepted},
		{"F R gas - spa-sc", invalid},
	})
}

func Test516_NoLegalEscape(t *testing.T) {
	checkTurnCase(t, 516, 1, []expectedCommand{
		{"F mid - por", accepted},
		{"F spa-nc S F mid - por", accepted},
	})
}

func Test517_OmitRetreatAfterOrderingAHold(t *testing.T) {
	checkTurnCase(t, 517, 1, []expectedCommand{
		{"A bel - H", accepted},
		{"A pic - bel", accepted},
		{"A bur S A pic - bel", accepted},
	})
}

func Test518_TwoRetreatsCollide(t *testing.T) {
	checkTurnCase(t, 518, 1, []expectedCommand{
		{"A bre - par", accepted},
		{"A gas S A bre - par", accepted},
		{"A hol - bel", accepted},
		{"A ruh S A hol - bel", accepted},
		{"A R par - pic", accepted},
		{"A R bel - pic", accepted},
	})
}

func Test519_DistinctSimultaneousRetreats(t *testing.T) {
	checkTurnCase(t, 519, 1, []expectedCommand{
		{"A bre - par", accepted},
		{"A gas S A bre - par", accepted},
		{"A hol - bel", accepted},
		{"A ruh S A hol - bel", accepted},
		{"A R par - pic", accepted},
		{"A R bel - bur", accepted},
	})
}

func Test520_ThreeWayRetreatCollision(t *testing.T) {
	checkTurnCase(t, 520, 1, []expectedCommand{
		{"A bre - par", accepted},
		{"A gas S A bre - par", accepted},
		{"A hol - bel", accepted},
		{"A ruh S A hol - bel", accepted},
		{"A mun - bur", accepted},
		{"A mar S A mun - bur", accepted},
		{"A R par - pic", accepted},
		{"A R bel - pic", accepted},
		{"A R bur - pic", accepted},
	})
}

func Test521_RetreatsCannotReceiveSupport(t *testing.T) {
	checkTurnCase(t, 521, 1, []expectedCommand{
		{"A bre - par", accepted},
		{"A gas S A bre - par", accepted},
		{"A hol - bel", accepted},
		{"A ruh S A hol - bel", accepted},
		{"F eng S A bel - pic", unmatched},
		{"A R par - pic", accepted},
		{"A R bel - pic", accepted},
	})
}

func Test522_RetreatsToDifferentCoastsCollide(t *testing.T) {
	checkTurnCase(t, 522, 1, []expectedCommand{
		{"A par - gas", accepted},
		{"A bre S A par - gas", accepted},
		{"A bur - mar", accepted},
		{"A pie S A bur - mar", accepted},
		{"F R gas - spa-nc", accepted},
		{"F R mar - spa-sc", accepted},
	})
}

func Test523_TwoUsableConditionalRetreats(t *testing.T) {
	checkTurnCase(t, 523, 1, []expectedCommand{
		{"A R bel - hol", accepted},
		{"A R bel - ruh", duplicate},
		{"A pic - bel", accepted},
		{"A bur S A pic - bel", accepted},
	})
}

func Test524_FirstRefugeOccupiedAfterMovement(t *testing.T) {
	checkTurnCase(t, 524, 1, []expectedCommand{
		{"A R bel - hol", occupiedRefuge},
		{"A R bel - ruh", accepted},
		{"A pic - bel", accepted},
		{"A bur S A pic - bel", accepted},
		{"F nth - hol", accepted},
	})
}

func Test525_NonadjacentFirstRefugeLegalSecondRefuge(t *testing.T) {
	checkTurnCase(t, 525, 1, []expectedCommand{
		{"A R bel - ber", invalidBecause("No route from 'bel' to 'ber'.")},
		{"A R bel - hol", accepted},
		{"A pic - bel", accepted},
		{"A bur S A pic - bel", accepted},
	})
}

func Test526_AcceptedRetreatCollidesDespiteLaterAlternative(t *testing.T) {
	checkTurnCase(t, 526, 1, []expectedCommand{
		{"A bre - par", accepted},
		{"A gas S A bre - par", accepted},
		{"A hol - bel", accepted},
		{"A ruh S A hol - bel", accepted},
		{"A R par - pic", accepted},
		{"A R bel - pic", accepted},
		{"A R bel - bur", duplicate},
	})

	// Each army can reach Picardy when it is the only one retreating there.
	// Keep Belgium's fallback order to check that selection precedes collision.
	for _, origin := range []string{"par", "bel"} {
		t.Run("only "+origin+" retreats", func(t *testing.T) {
			f, players := loadFixture(t, 526)
			f.Orders = slices.DeleteFunc(f.Orders, func(o models.Order) bool {
				return strings.HasPrefix(o.Value, "A R ") && !strings.HasPrefix(o.Value, "A R "+origin+" ")
			})
			power := "England"
			if origin == "bel" {
				power = "Germany"
			}
			placeArmy(t, f.OutputBoard, power, origin, "pic")
			got, _ := resolveSpring(t, f, players)
			checkBoard(t, got, f.OutputBoard)
		})
	}
}
