package tests

import "testing"

func Test101_ArmyHoldsOnExplicitOrder(t *testing.T) {
	checkMovement(t, 101, []expectedCommand{
		{"A par - H", accepted},
	})
}

func Test102_ArmyHoldsWithoutOrders(t *testing.T) {
	checkMovement(t, 102, nil)
}

func Test103_UnopposedArmyMovesToBurgundy(t *testing.T) {
	checkMovement(t, 103, []expectedCommand{
		{"A par - bur", accepted},
	})
}

func Test104_MoveRequiresUnitType(t *testing.T) {
	checkMovement(t, 104, []expectedCommand{
		{"par - bur", invalid},
	})
}

func Test105_MoveRejectsAmbiguousDestination(t *testing.T) {
	checkMovement(t, 105, []expectedCommand{
		{"A par - pic or bur", invalid},
	})
}

func Test106_HoldOverridesLaterMove(t *testing.T) {
	checkMovement(t, 106, []expectedCommand{
		{"A par - H", accepted},
		{"A par - bur", duplicate},
	})
}

func Test107_HoldOverridesEarlierMove(t *testing.T) {
	checkMovement(t, 107, []expectedCommand{
		{"A par - bur", duplicate},
		{"A par - H", accepted},
	})
}

func Test108_ValidMoveFollowsInvalidMove(t *testing.T) {
	checkMovement(t, 108, []expectedCommand{
		{"A par - lon", invalid},
		{"A par - bur", accepted},
	})
}

func Test109_FirstOfTwoValidMovesWins(t *testing.T) {
	checkMovement(t, 109, []expectedCommand{
		{"A par - pic", accepted},
		{"A par - bur", duplicate},
	})
}

func Test110_RejectUnknownProvinceAndForeignUnit(t *testing.T) {
	checkMovement(t, 110, []expectedCommand{
		{"A bur - swi", invalid},
		{"A mun - ruh", invalid},
	})
}

func Test111_RejectWrongTerrainAndNonadjacentMove(t *testing.T) {
	checkMovement(t, 111, []expectedCommand{
		{"A lon - nth", invalid},
		{"F bre - par", invalid},
		{"A mun - bel", invalid},
	})
}

func Test112_FleetCannotCrossFromRomeToApulia(t *testing.T) {
	checkMovement(t, 112, []expectedCommand{
		{"F rom - apu", invalid},
	})
}

func Test113_FleetFollowsCoastFromRomeToNaples(t *testing.T) {
	checkMovement(t, 113, []expectedCommand{
		{"F rom - nap", accepted},
	})
}

func Test114_FleetMustChooseBetweenReachableCoasts(t *testing.T) {
	checkMovement(t, 114, []expectedCommand{
		{"F con - bul", invalid},
	})
}

func Test115_FleetEntersBulgariasEastCoast(t *testing.T) {
	checkMovement(t, 115, []expectedCommand{
		{"F con - bul-ec", accepted},
	})
}

func Test116_FleetMustNameEvenTheOnlyReachableCoast(t *testing.T) {
	checkMovement(t, 116, []expectedCommand{
		{"F gas - spa", invalid},
	})
}

func Test117_GasconyCannotReachSpainsSouthCoast(t *testing.T) {
	checkMovement(t, 117, []expectedCommand{
		{"F gas - spa-sc", invalid},
	})
}

func Test118_SpainNorthCoastCannotReachWesternMed(t *testing.T) {
	checkMovement(t, 118, []expectedCommand{
		{"F spa-nc - wes", invalid},
	})
}

func Test119_SpainSouthCoastCanReachWesternMed(t *testing.T) {
	checkMovement(t, 119, []expectedCommand{
		{"F spa-sc - wes", accepted},
	})
}

func Test120_FleetMovesToDifferentSpanishCoastsBounce(t *testing.T) {
	checkMovement(t, 120, []expectedCommand{
		{"F gas - spa-nc", accepted},
		{"F gol - spa-sc", accepted},
	})
}

func Test120_ArmyCanRetreatToUncontestedSpain(t *testing.T) {
	f, players := loadFixture(t, 120)
	f.Orders = nil // Both fleets hold.
	addMarseillesDislodgement(t, &f)
	placeArmy(t, f.OutputBoard, "England", "mar", "spa")

	got, _ := resolveSpring(t, f, players)
	checkBoard(t, got, f.OutputBoard)
}

func Test120_SpanishStandoffBlocksArmyRetreat(t *testing.T) {
	f, players := loadFixture(t, 120)
	addMarseillesDislodgement(t, &f)

	got, _ := resolveSpring(t, f, players)
	checkBoard(t, got, f.OutputBoard)
}

func Test121_ArmyEntersBulgariaWithoutNamingACoast(t *testing.T) {
	checkMovement(t, 121, []expectedCommand{
		{"A ser - bul", accepted},
	})
}

func Test122_FleetsEnterKielAndConstantinople(t *testing.T) {
	checkMovement(t, 122, []expectedCommand{
		{"F hol - kie", accepted},
		{"F bla - con", accepted},
	})
}

func Test123_FleetsLeaveKielAndConstantinople(t *testing.T) {
	checkMovement(t, 123, []expectedCommand{
		{"F kie - ber", accepted},
		{"F con - aeg", accepted},
	})
}

func Test124_UnitsCannotSkipAcrossWaterways(t *testing.T) {
	checkMovement(t, 124, []expectedCommand{
		{"F hol - ber", invalid},
		{"F bla - aeg", invalid},
		{"F bal - ska", invalid},
		{"A den - ber", invalid},
	})
}

func Test125_ArmyCrossesToSwedenAsFleetEntersDenmark(t *testing.T) {
	checkMovement(t, 125, []expectedCommand{
		{"A den - swe", accepted},
		{"F bal - den", accepted},
	})
}

func Test126_RulebookRoutesForEnglishChannelAndParis(t *testing.T) {
	checkMovement(t, 126, []expectedCommand{
		{"F eng - nth", accepted},
		{"A par - mar", invalid},
	})
}

func Test127_FleetsEnterSwedenAndDenmarkWithoutCoasts(t *testing.T) {
	checkMovement(t, 127, []expectedCommand{
		{"F den - swe", accepted},
		{"F bal - den", accepted},
	})
}

func Test128_ArmiesLeaveKielAndConstantinople(t *testing.T) {
	checkMovement(t, 128, []expectedCommand{
		{"A con - smy", accepted},
		{"A kie - ber", accepted},
	})
}

func Test129_MoveTakesPrecedenceOverEarlierSupport(t *testing.T) {
	checkMovement(t, 129, []expectedCommand{
		{"A bur S A pic - bel", duplicate},
		{"A pic - bel", accepted},
		{"A bur - gas", accepted},
	})
}

func Test130_MoveOverridesMoversLaterSupport(t *testing.T) {
	checkMovement(t, 130, []expectedCommand{
		{"A bur - gas", accepted},
		{"A pic - bel", accepted},
		{"A bur S A pic - bel", duplicate},
	})
}

func Test131_BouncedMoveStillOverridesLaterMove(t *testing.T) {
	checkMovement(t, 131, []expectedCommand{
		{"A pic - bel", accepted},
		{"A pic - par", duplicate},
	})
}

func Test132_RejectUnknownIslandAndUnknownOrigin(t *testing.T) {
	checkMovement(t, 132, []expectedCommand{
		{"F tyn - sic", invalid},
		{"A xyz - rom", invalid},
	})
}
