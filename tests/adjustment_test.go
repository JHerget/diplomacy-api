package tests

import (
	"testing"

	"diplomacy-api/internal/models"
)

func Test601_SpringOccupation(t *testing.T) {
	checkTurnCase(t, 601, 1, []expectedCommand{
		{"A pic - bel", accepted},
	})
}

func Test602_FallOccupationTransfersControl(t *testing.T) {
	checkTurnCase(t, 602, 2, []expectedCommand{
		{"A pic - bel", accepted},
	})
}

func Test603_SpringVisitorLeavesByFall(t *testing.T) {
	checkTurnCase(t, 603, 2, []expectedCommand{
		{"A bel - pic", accepted},
	})
}

func Test604_CaptureNeutralByStaying(t *testing.T) {
	checkTurnCase(t, 604, 2, []expectedCommand{
		{"F bel - H", accepted},
	})
}

func Test605_VacantCenterRetainsControl(t *testing.T) {
	checkTurnCase(t, 605, 2, []expectedCommand{
		{"A bel - ruh", accepted},
	})
}

func Test606_RetreatCapturesCenterAndFundsBuild(t *testing.T) {
	checkTurnCase(t, 606, 2, []expectedCommand{
		{"A pic - bel", accepted},
		{"A bur S A pic - bel", accepted},
		{"A R bel - hol", accepted},
		{"A mun", accepted},
	})
}

func Test607_DisbandBeforeCensus(t *testing.T) {
	checkTurnCase(t, 607, 2, []expectedCommand{
		{"A pic - bel", accepted},
		{"A bur S A pic - bel", accepted},
		{"A mun", accepted},
	})
}

func Test608_SpringBuildIgnored(t *testing.T) {
	checkTurnCase(t, 608, 1, []expectedCommand{
		{"F bre", springAdjustment},
	})
}

func Test609_FallCoastalFleetBuild(t *testing.T) {
	checkTurnCase(t, 609, 2, []expectedCommand{
		{"F bre", accepted},
	})
}

func Test610_SpringRemovalIgnored(t *testing.T) {
	checkTurnCase(t, 610, 1, []expectedCommand{
		{"D A pic", springAdjustment},
	})
}

func Test611_FallOwnerSelectsRemoval(t *testing.T) {
	checkTurnCase(t, 611, 2, []expectedCommand{
		{"D A pic", accepted},
	})
}

func Test612_MultipleRequiredRemovals(t *testing.T) {
	checkTurnCase(t, 612, 2, []expectedCommand{
		{"D A bur", accepted},
		{"D F eng", accepted},
	})
}

func Test613_BalancedCount(t *testing.T) {
	checkTurnCase(t, 613, 2, []expectedCommand{
		{"A par", noBuildDue},
		{"D A pic", noRemovalDue},
	})
}

func Test614_FleetBuildInland(t *testing.T) {
	checkTurnCase(t, 614, 2, []expectedCommand{
		{"F par", invalid},
	})
}

func Test615_ArmyBuildInland(t *testing.T) {
	checkTurnCase(t, 615, 2, []expectedCommand{
		{"A par", accepted},
	})
}

func Test616_ArmyBuildOnCoast(t *testing.T) {
	checkTurnCase(t, 616, 2, []expectedCommand{
		{"A bre", accepted},
	})
}

func Test617_FleetBuildAtNamedSplitCoast(t *testing.T) {
	checkTurnCase(t, 617, 2, []expectedCommand{
		{"F stp-sc", accepted},
	})
}

func Test618_FleetBuildMissingCoast(t *testing.T) {
	checkTurnCase(t, 618, 2, []expectedCommand{
		{"F stp", invalidBecause("A coast was not provided when expected.")},
	})
}

func Test619_CaptureCannotCreateAHomeBuildSite(t *testing.T) {
	checkTurnCase(t, 619, 2, []expectedCommand{
		{"A bel", invalid},
	})
}

func Test620_HomeSiteUnderForeignControl(t *testing.T) {
	checkTurnCase(t, 620, 2, []expectedCommand{
		{"A par", invalid},
	})
}

func Test621_HomeSiteOccupiedByBuilder(t *testing.T) {
	checkTurnCase(t, 621, 2, []expectedCommand{
		{"F bre", invalid},
	})
}

func Test622_HomeSiteOccupiedByEnemy(t *testing.T) {
	checkTurnCase(t, 622, 2, []expectedCommand{
		{"F bre", invalid},
	})
}

func Test623_VacateHomeSiteAndRebuildThere(t *testing.T) {
	checkTurnCase(t, 623, 2, []expectedCommand{
		{"A bre - pic", accepted},
		{"F bre", accepted},
	})
}

func Test624_RecapturedHomeIsStillOccupied(t *testing.T) {
	checkTurnCase(t, 624, 2, []expectedCommand{
		{"A pic - par", accepted},
		{"A par", invalid},
	})
}

func Test625_TwoBuildsOneProvince(t *testing.T) {
	checkTurnCase(t, 625, 2, []expectedCommand{
		{"A bre", accepted},
		{"A bre", duplicate},
	})
}

func Test626_ExcessBuildRequestsForDistinctSites(t *testing.T) {
	checkTurnCase(t, 626, 2, []expectedCommand{
		{"A par", accepted},
		{"F bre", noBuildDue},
	})
}

func Test627_PartialWaiver(t *testing.T) {
	checkTurnCase(t, 627, 2, []expectedCommand{
		{"A par", accepted},
	})
}

func Test628_EntireSurplusWaived(t *testing.T) {
	checkTurnCase(t, 628, 2, nil)
}

func Test629_WaiverDoesNotBankABuild(t *testing.T) {
	f, players := loadFixture(t, 629)
	g := fixtureGame(f, players, 2)
	for _, number := range []int{2, 3} {
		turn := g.CurrentTurn()
		if turn.TurnNumber != number {
			t.Fatalf("want turn %d, got %d", number, turn.TurnNumber)
		}
		turn.Orders = []models.Order{
			{PlayerName: "France", Value: "A pic - H"}, {PlayerName: "Germany", Value: "F gas - H"},
		}
		turn.EndDate = 1
		closeGameTurn(t, &g)
		checkBoard(t, g.Board, f.InputBoard)
	}
	turn := g.CurrentTurn()
	if turn.TurnNumber != 4 {
		t.Fatalf("want Fall 1902, got turn %d", turn.TurnNumber)
	}
	turn.Orders = f.Orders
	turn.EndDate = 1
	closeGameTurn(t, &g)
	checkBoard(t, g.Board, f.OutputBoard)
}

func Test630_InvalidRemovalNamesEmptySite(t *testing.T) {
	checkTurnCase(t, 630, 2, []expectedCommand{
		{"D A gas", invalidBecause("No unit at province 'gas'.")},
		{"D A bur", accepted},
	})
}

func Test631_ActivePlayerSubmitsNoRequiredRemoval(t *testing.T) {
	checkTurnCase(t, 631, 2, []expectedCommand{
		{"A pic - H", accepted},
	})
}

func Test632_SurplusRemovalsRequested(t *testing.T) {
	checkTurnCase(t, 632, 2, []expectedCommand{
		{"D A pic", accepted},
		{"D A bur", noRemovalDue},
	})
}

func Test633_SplitCoastHomeArmy(t *testing.T) {
	checkTurnCase(t, 633, 2, []expectedCommand{
		{"A stp", accepted},
	})
}

func Test634_UnusableCoastalBuildWithoutType(t *testing.T) {
	checkTurnCase(t, 634, 2, []expectedCommand{
		{"bre", invalidBecause("Command format not recognized.")},
	})
}

func Test635_BuildAfterReleasingRecoveredCenter(t *testing.T) {
	checkTurnCase(t, 635, 4, []expectedCommand{
		{"A par - pic", accepted},
		{"A par", accepted},
	})
}

func Test636_CaptureEliminatesNeedForSubmittedRemoval(t *testing.T) {
	checkTurnCase(t, 636, 2, []expectedCommand{
		{"A pic - bre", accepted},
		{"D A pic", noRemovalDue},
	})
}

func Test637_CaptureFundsBuildInTheSameTurn(t *testing.T) {
	checkTurnCase(t, 637, 2, []expectedCommand{
		{"A pic - bel", accepted},
		{"A par", accepted},
	})
}

func Test638_EnemyRetreatCancelsApparentBuild(t *testing.T) {
	checkTurnCase(t, 638, 2, []expectedCommand{
		{"A pic - bel", accepted},
		{"A bur S A pic - bel", accepted},
		{"A R bel - hol", accepted},
		{"A par", noBuildDue},
	})
}

func Test639_ActivePlayerSubmitsOneOfTwoRequiredRemovals(t *testing.T) {
	checkTurnCase(t, 639, 2, []expectedCommand{
		{"D F eng", accepted},
	})
}
