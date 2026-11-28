package tests

import (
	"slices"
	"strings"
	"testing"

	"diplomacy-api/internal/models"
)

func Test601_ConditionFalse(t *testing.T) {
	checkTurnCase(t, 601, 1, []expectedCommand{
		{"A pic - bel", accepted},
		{"A R bel - hol", dormant},
	})
}

func Test602_ConditionTrue(t *testing.T) {
	checkTurnCase(t, 602, 1, []expectedCommand{
		{"A pic - bel", accepted},
		{"A bur S A pic - bel", accepted},
		{"A R bel - hol", accepted},
	})
}

func Test603_NoRetreatSubmitted(t *testing.T) {
	checkTurnCase(t, 603, 1, []expectedCommand{
		{"A pic - bel", accepted},
		{"A bur S A pic - bel", accepted},
	})
}

func Test604_MovementSubmittedBeforeRetreat(t *testing.T) {
	checkTurnCase(t, 604, 1, []expectedCommand{
		{"A bel - ruh", accepted},
		{"A R bel - hol", accepted},
		{"A pic - bel", accepted},
		{"A bur S A pic - bel", accepted},
	})
}

func Test605_RetreatSubmittedBeforeMovement(t *testing.T) {
	checkTurnCase(t, 605, 1, []expectedCommand{
		{"A R bel - hol", accepted},
		{"A bel - ruh", accepted},
		{"A pic - bel", accepted},
		{"A bur S A pic - bel", accepted},
	})
}

func Test606_RetreatSubmittedFirstAndUnitEscapesByMoving(t *testing.T) {
	checkTurnCase(t, 606, 1, []expectedCommand{
		{"A R bel - hol", dormant},
		{"A bel - ruh", accepted},
		{"A pic - bel", accepted},
		{"A bur S A pic - bel", accepted},
	})
}

func Test607_OccupiedRefuge(t *testing.T) {
	checkTurnCase(t, 607, 1, []expectedCommand{
		{"A pic - bel", accepted},
		{"A bur S A pic - bel", accepted},
		{"A R bel - hol", occupiedRefuge},
	})
}

func Test608_NewlyVacatedRefuge(t *testing.T) {
	checkTurnCase(t, 608, 1, []expectedCommand{
		{"A pic - bel", accepted},
		{"A bur S A pic - bel", accepted},
		{"A hol - kie", accepted},
		{"A R bel - hol", accepted},
	})
}

func Test609_NewlyOccupiedRefuge(t *testing.T) {
	checkTurnCase(t, 609, 1, []expectedCommand{
		{"A pic - bel", accepted},
		{"A bur S A pic - bel", accepted},
		{"F hel - hol", accepted},
		{"A R bel - hol", occupiedRefuge},
	})
}

func Test610_EmptyDislodgerOriginForbidden(t *testing.T) {
	checkTurnCase(t, 610, 1, []expectedCommand{
		{"A pic - bel", accepted},
		{"A bur S A pic - bel", accepted},
		{"A R bel - pic", dislodgerOrigin},
	})
}

func Test611_EmptyStandoffProvinceForbidden(t *testing.T) {
	checkTurnCase(t, 611, 1, []expectedCommand{
		{"A pic - bel", accepted},
		{"A bur S A pic - bel", accepted},
		{"F nth - hol", accepted},
		{"A kie - hol", accepted},
		{"A R bel - hol", standoff},
	})
}

func Test612_NonadjacentRefuge(t *testing.T) {
	checkTurnCase(t, 612, 1, []expectedCommand{
		{"A pic - bel", accepted},
		{"A bur S A pic - bel", accepted},
		{"A R bel - ber", invalidBecause("No route from 'bel' to 'ber'.")},
	})
}

func Test613_RetreatCannotUseAConvoy(t *testing.T) {
	checkTurnCase(t, 613, 1, []expectedCommand{
		{"A pic - bel", accepted},
		{"A bur S A pic - bel", accepted},
		{"F nth C A bel - lon", accepted},
		{"A R bel - lon", invalid},
	})
}

func Test614_FleetRetreatsToReachableCoast(t *testing.T) {
	checkTurnCase(t, 614, 1, []expectedCommand{
		{"F mid - gas", accepted},
		{"A par S F mid - gas", accepted},
		{"F R gas - spa-nc", accepted},
	})
}

func Test615_FleetRetreatsToUnreachableCoast(t *testing.T) {
	checkTurnCase(t, 615, 1, []expectedCommand{
		{"F mid - gas", accepted},
		{"A par S F mid - gas", accepted},
		{"F R gas - spa-sc", invalid},
	})
}

func Test616_NoLegalEscape(t *testing.T) {
	checkTurnCase(t, 616, 1, []expectedCommand{
		{"F mid - por", accepted},
		{"F spa-nc S F mid - por", accepted},
	})
}

func Test617_OmitRetreatAfterOrderingAHold(t *testing.T) {
	checkTurnCase(t, 617, 1, []expectedCommand{
		{"A bel - H", accepted},
		{"A pic - bel", accepted},
		{"A bur S A pic - bel", accepted},
	})
}

func Test618_TwoRetreatsCollide(t *testing.T) {
	checkTurnCase(t, 618, 1, []expectedCommand{
		{"A bre - par", accepted},
		{"A gas S A bre - par", accepted},
		{"A hol - bel", accepted},
		{"A ruh S A hol - bel", accepted},
		{"A R par - pic", accepted},
		{"A R bel - pic", accepted},
	})
}

func Test619_DistinctSimultaneousRetreats(t *testing.T) {
	checkTurnCase(t, 619, 1, []expectedCommand{
		{"A bre - par", accepted},
		{"A gas S A bre - par", accepted},
		{"A hol - bel", accepted},
		{"A ruh S A hol - bel", accepted},
		{"A R par - pic", accepted},
		{"A R bel - bur", accepted},
	})
}

func Test620_ThreeWayRetreatCollision(t *testing.T) {
	checkTurnCase(t, 620, 1, []expectedCommand{
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

func Test621_RetreatsCannotReceiveSupport(t *testing.T) {
	checkTurnCase(t, 621, 1, []expectedCommand{
		{"A bre - par", accepted},
		{"A gas S A bre - par", accepted},
		{"A hol - bel", accepted},
		{"A ruh S A hol - bel", accepted},
		{"F eng S A bel - pic", unmatched},
		{"A R par - pic", accepted},
		{"A R bel - pic", accepted},
	})
}

func Test622_RetreatsToDifferentCoastsCollide(t *testing.T) {
	checkTurnCase(t, 622, 1, []expectedCommand{
		{"A par - gas", accepted},
		{"A bre S A par - gas", accepted},
		{"A bur - mar", accepted},
		{"A pie S A bur - mar", accepted},
		{"F R gas - spa-nc", accepted},
		{"F R mar - spa-sc", accepted},
	})
}

func Test623_TwoUsableConditionalRetreats(t *testing.T) {
	checkTurnCase(t, 623, 1, []expectedCommand{
		{"A R bel - hol", accepted},
		{"A R bel - ruh", duplicate},
		{"A pic - bel", accepted},
		{"A bur S A pic - bel", accepted},
	})
}

func Test624_FirstRefugeOccupiedAfterMovement(t *testing.T) {
	checkTurnCase(t, 624, 1, []expectedCommand{
		{"A R bel - hol", occupiedRefuge},
		{"A R bel - ruh", accepted},
		{"A pic - bel", accepted},
		{"A bur S A pic - bel", accepted},
		{"F nth - hol", accepted},
	})
}

func Test625_NonadjacentFirstRefugeLegalSecondRefuge(t *testing.T) {
	checkTurnCase(t, 625, 1, []expectedCommand{
		{"A R bel - ber", invalidBecause("No route from 'bel' to 'ber'.")},
		{"A R bel - hol", accepted},
		{"A pic - bel", accepted},
		{"A bur S A pic - bel", accepted},
	})
}

func Test626_AcceptedRetreatCollidesDespiteLaterAlternative(t *testing.T) {
	checkTurnCase(t, 626, 1, []expectedCommand{
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
			f, players := loadFixture(t, 626)
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
