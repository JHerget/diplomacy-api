package tests

import (
	"slices"
	"testing"

	"diplomacy-api/internal/models"
)

func Test301_ForeignFleetSupportsAnAttackFromAnotherProvince(t *testing.T) {
	checkMovement(t, 301, []expectedCommand{
		{"A pic - bel", accepted},
		{"F nth S A pic - bel", accepted},
	})
}

func Test302_HoldSupportTiesASupportedAttack(t *testing.T) {
	checkMovement(t, 302, []expectedCommand{
		{"A pic - bel", accepted},
		{"F nth S A pic - bel", accepted},
		{"A bel - H", accepted},
		{"A ruh S A bel - H", accepted},
	})
}

func Test303_UnorderedDefenderReceivesHoldSupport(t *testing.T) {
	checkMovement(t, 303, []expectedCommand{
		{"A pic - bel", accepted},
		{"F nth S A pic - bel", accepted},
		{"A ruh S A bel - H", accepted},
	})
}

func Test304_MovingDefenderCannotReceiveHoldSupport(t *testing.T) {
	checkMovement(t, 304, []expectedCommand{
		{"A pic - bel", accepted},
		{"F nth S A pic - bel", accepted},
		{"A bel - hol", accepted},
		{"A ruh S A bel - H", unmatched}, // Belgium ordered a move, so there is no hold to support.
	})
}

func Test305_SupportMustMatchTheActualDestination(t *testing.T) {
	checkMovement(t, 305, []expectedCommand{
		{"A pic - bel", accepted},
		{"F eng S A pic - bre", unmatched},
	})
}

func Test306_SupportMustMatchTheActualOrigin(t *testing.T) {
	checkMovement(t, 306, []expectedCommand{
		{"A pic - bel", accepted},
		{"F nth S A bur - bel", unmatched},
	})
}

func Test307_MoveSupportDoesNotDefendAnUnorderedUnit(t *testing.T) {
	checkMovement(t, 307, []expectedCommand{
		{"A pic - bel", accepted},
		{"A bur S A pic - bel", accepted},
		{"A ruh S A bel - hol", unmatched},
	})
}

func Test308_FleetCannotSupportIntoInlandSilesia(t *testing.T) {
	checkMovement(t, 308, []expectedCommand{
		{"A pru - sil", accepted},
		{"F ber S A pru - sil", invalid},
	})
}

func Test309_ArmyCannotSupportIntoTheNorthSea(t *testing.T) {
	checkMovement(t, 309, []expectedCommand{
		{"F yor - nth", accepted},
		{"A lon S F yor - nth", invalid},
	})
}

func Test310_FleetCannotSupportAcrossAnUnreachableCoastline(t *testing.T) {
	checkMovement(t, 310, []expectedCommand{
		{"A ven - apu", accepted},
		{"F rom S A ven - apu", invalid},
	})
}

func Test311_OppositeCoastSupportHelpsFleetEnterSpain(t *testing.T) {
	checkMovement(t, 311, []expectedCommand{
		{"F gas - spa-nc", accepted},
		{"F gol S F gas - spa", accepted},
	})
}

func Test312_EquallySupportedAttacksStandOffInBelgium(t *testing.T) {
	commands := []expectedCommand{
		{"A pic - bel", accepted},
		{"A bur S A pic - bel", accepted},
		{"A ruh - bel", accepted},
		{"A hol S A ruh - bel", accepted},
	}
	checkMovement(t, 312, commands)
	checkRetreatToBelgium(t, 312, commands)
}

func Test313_ExtraSupportBreaksTheTieInBelgium(t *testing.T) {
	checkMovement(t, 313, []expectedCommand{
		{"A pic - bel", accepted},
		{"A bur S A pic - bel", accepted},
		{"A ruh - bel", accepted},
		{"A hol S A ruh - bel", accepted},
		{"F nth S A pic - bel", accepted},
	})
}

func Test314_BouncingOutsideAttackCutsSupport(t *testing.T) {
	checkMovement(t, 314, []expectedCommand{
		{"A pru - war", accepted},
		{"A sil S A pru - war", supportCutFrom("boh")},
		{"A boh - sil", accepted},
	})
}

func Test315_SamePowerAttackLeavesSupportIntact(t *testing.T) {
	checkMovement(t, 315, []expectedCommand{
		{"A pru - war", accepted},
		{"A sil S A pru - war", accepted},
		{"A boh - sil", accepted},
	})
}

func Test316_AttackFromSupportedDestinationMustDislodgeToCut(t *testing.T) {
	checkMovement(t, 316, []expectedCommand{
		{"A pru - war", accepted},
		{"A sil S A pru - war", accepted},
		{"A war - sil", accepted},
	})
}

func Test317_DislodgementByTheSupportedDestinationCutsSupport(t *testing.T) {
	commands := []expectedCommand{
		{"F ber - pru", accepted},
		{"A sil S F ber - pru", supportCutFrom("pru")},
		{"A pru - sil", accepted},
		{"A war S A pru - sil", accepted},
		{"F bal - pru", accepted},
	}
	checkMovement(t, 317, commands)
	t.Run("Reversed orders and province lists", func(t *testing.T) {
		f, players := loadFixture(t, 317)
		slices.Reverse(f.Orders)
		slices.Reverse(f.InputBoard)
		slices.Reverse(f.OutputBoard)
		want := slices.Clone(commands)
		slices.Reverse(want)
		checkSpringFixture(t, f, players, want)
	})
	t.Run("Standoff in vacated Prussia blocks a retreat from Livonia", func(t *testing.T) {
		f, players := loadFixture(t, 317)
		// Livonia can reach Prussia and its dislodger comes from Moscow.
		placeArmy(t, f.InputBoard, "Italy", "lvn", "lvn")
		placeArmy(t, f.InputBoard, "Russia", "mos", "mos")
		placeArmy(t, f.InputBoard, "Russia", "stp", "stp")
		placeArmy(t, f.OutputBoard, "Russia", "mos", "lvn")
		placeArmy(t, f.OutputBoard, "Russia", "stp", "stp")
		f.Orders = append(f.Orders,
			models.Order{PlayerName: "Russia", Value: "A mos - lvn, A stp S A mos - lvn"},
			models.Order{PlayerName: "Italy", Value: "A R lvn - pru"},
		)

		checkSpringFixture(t, f, players, append(commands,
			expectedCommand{"A mos - lvn", accepted},
			expectedCommand{"A stp S A mos - lvn", accepted},
			expectedCommand{"A R lvn - pru", standoff},
		))
	})
}

func Test318_DislodgedAttackerStillCutsSupportElsewhere(t *testing.T) {
	checkMovement(t, 318, []expectedCommand{
		{"A pru - war", accepted},
		{"A sil S A pru - war", supportCutFrom("boh")},
		{"A boh - sil", accepted},
		{"A tyr - boh", accepted},
		{"A mun S A tyr - boh", accepted},
	})
}

func Test319_BouncingAttackCutsHoldSupport(t *testing.T) {
	checkMovement(t, 319, []expectedCommand{
		{"A pic - bel", accepted},
		{"A bur S A pic - bel", accepted},
		{"A ruh S A bel - H", supportCutFrom("mun")},
		{"A mun - ruh", accepted},
	})
}

func Test320_SupportingUnitCanReceiveHoldSupport(t *testing.T) {
	checkMovement(t, 320, []expectedCommand{
		{"A pic - bel", accepted},
		{"A bur S A pic - bel", supportCutFrom("mun")},
		{"A par S A bur - H", accepted},
		{"A mun - bur", accepted},
		{"A ruh S A mun - bur", accepted},
	})
}

func Test321_ForeignSupportCannotForceSelfDislodgement(t *testing.T) {
	checkMovement(t, 321, []expectedCommand{
		{"A pic - bel", accepted},
		{"F nth S A pic - bel", accepted},
	})
}

func Test322_DefendersPowerCannotHelpDislodgeItsOwnUnit(t *testing.T) {
	checkMovement(t, 322, []expectedCommand{
		{"A pic - bel", accepted},
		{"A ruh S A pic - bel", accepted},
	})
}

func Test323_IndependentSupportIsEnoughToDislodge(t *testing.T) {
	checkMovement(t, 323, []expectedCommand{
		{"A pic - bel", accepted},
		{"A ruh S A pic - bel", accepted},
		{"A bur S A pic - bel", accepted},
	})
}

func Test324_TwoFriendlyAttacksCanBounceAForeignAttack(t *testing.T) {
	checkMovement(t, 324, []expectedCommand{
		{"A pic - bel", accepted},
		{"A bur - bel", accepted},
		{"A ruh - bel", accepted},
	})
}

func Test325_ForeignSupportBreaksASelfStandoff(t *testing.T) {
	checkMovement(t, 325, []expectedCommand{
		{"A pic - bel", accepted},
		{"A bur - bel", accepted},
		{"A ruh - bel", accepted},
		{"F nth S A pic - bel", accepted},
	})
}

func Test326_AttemptedSelfDislodgementStillOpposesARival(t *testing.T) {
	checkMovement(t, 326, []expectedCommand{
		{"A pic - bel", accepted},
		{"A bur S A pic - bel", accepted},
		{"A hol - bel", accepted},
		{"A ruh S A hol - bel", accepted},
	})
}

func Test327_FailedSupporterCanReceiveHoldSupport(t *testing.T) {
	checkMovement(t, 327, []expectedCommand{
		{"A pic - bel", accepted},
		{"A bur S A pic - bel", accepted},
		{"A bel S A par - pic", unmatched},
		{"A ruh S A bel - H", accepted},
	})
}

func Test328_IllegalMoveLeavesDefenderEligibleForHoldSupport(t *testing.T) {
	checkMovement(t, 328, []expectedCommand{
		{"A pic - bel", accepted},
		{"A bur S A pic - bel", accepted},
		{"A bel - mos", invalid},
		{"A ruh S A bel - H", accepted},
	})
}

func Test329_DefendersPowerCanSupportEntryAfterItVacates(t *testing.T) {
	checkMovement(t, 329, []expectedCommand{
		{"A pic - bel", accepted},
		{"A ruh S A pic - bel", accepted},
		{"A bel - hol", accepted},
		{"F nth - bel", accepted},
	})
}

func Test330_IllegalFleetAttackCannotCutArmySupport(t *testing.T) {
	checkMovement(t, 330, []expectedCommand{
		{"A pru - war", accepted},
		{"A sil S A pru - war", accepted},
		{"F ber - sil", invalid},
	})
}
