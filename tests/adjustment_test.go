package tests

import (
	"testing"

	"diplomacy-api/internal/models"
)

func Test701_SpringOccupation(t *testing.T) {
	checkTurnCase(t, 701, 1, []expectedCommand{
		{"A pic - bel", accepted},
	})
}

func Test702_FallOccupationTransfersControl(t *testing.T) {
	checkTurnCase(t, 702, 2, []expectedCommand{
		{"A pic - bel", accepted},
	})
}

func Test703_SpringVisitorLeavesByFall(t *testing.T) {
	checkTurnCase(t, 703, 2, []expectedCommand{
		{"A bel - pic", accepted},
	})
}

func Test704_CaptureNeutralByStaying(t *testing.T) {
	checkTurnCase(t, 704, 2, []expectedCommand{
		{"F bel - H", accepted},
	})
}

func Test705_VacantCenterRetainsControl(t *testing.T) {
	checkTurnCase(t, 705, 2, []expectedCommand{
		{"A bel - ruh", accepted},
	})
}

func Test706_RetreatCapturesCenterAndFundsBuild(t *testing.T) {
	checkTurnCase(t, 706, 2, []expectedCommand{
		{"A pic - bel", accepted},
		{"A bur S A pic - bel", accepted},
		{"A R bel - hol", accepted},
		{"A mun", accepted},
	})
}

func Test707_DisbandBeforeCensus(t *testing.T) {
	checkTurnCase(t, 707, 2, []expectedCommand{
		{"A pic - bel", accepted},
		{"A bur S A pic - bel", accepted},
		{"A mun", accepted},
	})
}

func Test708_SpringBuildIgnored(t *testing.T) {
	checkTurnCase(t, 708, 1, []expectedCommand{
		{"F bre", springAdjustment},
	})
}

func Test709_FallCoastalFleetBuild(t *testing.T) {
	checkTurnCase(t, 709, 2, []expectedCommand{
		{"F bre", accepted},
	})
}

func Test710_SpringRemovalIgnored(t *testing.T) {
	checkTurnCase(t, 710, 1, []expectedCommand{
		{"D A pic", springAdjustment},
	})
}

func Test711_FallOwnerSelectsRemoval(t *testing.T) {
	checkTurnCase(t, 711, 2, []expectedCommand{
		{"D A pic", accepted},
	})
}

func Test712_MultipleRequiredRemovals(t *testing.T) {
	checkTurnCase(t, 712, 2, []expectedCommand{
		{"D A bur", accepted},
		{"D F eng", accepted},
	})
}

func Test713_BalancedCount(t *testing.T) {
	checkTurnCase(t, 713, 2, []expectedCommand{
		{"A par", noBuildDue},
		{"D A pic", noRemovalDue},
	})
}

func Test714_FleetBuildInland(t *testing.T) {
	checkTurnCase(t, 714, 2, []expectedCommand{
		{"F par", invalid},
	})
}

func Test715_ArmyBuildInland(t *testing.T) {
	checkTurnCase(t, 715, 2, []expectedCommand{
		{"A par", accepted},
	})
}

func Test716_ArmyBuildOnCoast(t *testing.T) {
	checkTurnCase(t, 716, 2, []expectedCommand{
		{"A bre", accepted},
	})
}

func Test717_FleetBuildAtNamedSplitCoast(t *testing.T) {
	checkTurnCase(t, 717, 2, []expectedCommand{
		{"F stp-sc", accepted},
	})
}

func Test718_FleetBuildMissingCoast(t *testing.T) {
	checkTurnCase(t, 718, 2, []expectedCommand{
		{"F stp", invalidBecause("A coast was not provided when expected.")},
	})
}

func Test719_CaptureCannotCreateAHomeBuildSite(t *testing.T) {
	checkTurnCase(t, 719, 2, []expectedCommand{
		{"A bel", invalid},
	})
}

func Test720_HomeSiteUnderForeignControl(t *testing.T) {
	checkTurnCase(t, 720, 2, []expectedCommand{
		{"A par", invalid},
	})
}

func Test721_HomeSiteOccupiedByBuilder(t *testing.T) {
	checkTurnCase(t, 721, 2, []expectedCommand{
		{"F bre", invalid},
	})
}

func Test722_HomeSiteOccupiedByEnemy(t *testing.T) {
	checkTurnCase(t, 722, 2, []expectedCommand{
		{"F bre", invalid},
	})
}

func Test723_VacateHomeSiteAndRebuildThere(t *testing.T) {
	checkTurnCase(t, 723, 2, []expectedCommand{
		{"A bre - pic", accepted},
		{"F bre", accepted},
	})
}

func Test724_RecapturedHomeIsStillOccupied(t *testing.T) {
	checkTurnCase(t, 724, 2, []expectedCommand{
		{"A pic - par", accepted},
		{"A par", invalid},
	})
}

func Test725_TwoBuildsOneProvince(t *testing.T) {
	checkTurnCase(t, 725, 2, []expectedCommand{
		{"A bre", accepted},
		{"A bre", duplicate},
	})
}

func Test726_ExcessBuildRequestsForDistinctSites(t *testing.T) {
	checkTurnCase(t, 726, 2, []expectedCommand{
		{"A par", accepted},
		{"F bre", noBuildDue},
	})
}

func Test727_PartialWaiver(t *testing.T) {
	checkTurnCase(t, 727, 2, []expectedCommand{
		{"A par", accepted},
	})
}

func Test728_EntireSurplusWaived(t *testing.T) {
	checkTurnCase(t, 728, 2, nil)
}

func Test729_WaiverDoesNotBankABuild(t *testing.T) {
	f, players := loadFixture(t, 729)
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

func Test730_InvalidRemovalNamesEmptySite(t *testing.T) {
	checkTurnCase(t, 730, 2, []expectedCommand{
		{"D A gas", invalidBecause("No unit at province 'gas'.")},
		{"D A bur", accepted},
	})
}

func Test731_ActivePlayerSubmitsNoRequiredRemoval(t *testing.T) {
	checkTurnCase(t, 731, 2, []expectedCommand{
		{"A pic - H", accepted},
	})
}

func Test732_SurplusRemovalsRequested(t *testing.T) {
	checkTurnCase(t, 732, 2, []expectedCommand{
		{"D A pic", accepted},
		{"D A bur", noRemovalDue},
	})
}

func Test733_SplitCoastHomeArmy(t *testing.T) {
	checkTurnCase(t, 733, 2, []expectedCommand{
		{"A stp", accepted},
	})
}

func Test734_UnusableCoastalBuildWithoutType(t *testing.T) {
	checkTurnCase(t, 734, 2, []expectedCommand{
		{"bre", invalidBecause("Command format not recognized.")},
	})
}

func Test735_BuildAfterReleasingRecoveredCenter(t *testing.T) {
	checkTurnCase(t, 735, 4, []expectedCommand{
		{"A par - pic", accepted},
		{"A par", accepted},
	})
}

func Test736_CaptureEliminatesNeedForSubmittedRemoval(t *testing.T) {
	checkTurnCase(t, 736, 2, []expectedCommand{
		{"A pic - bre", accepted},
		{"D A pic", noRemovalDue},
	})
}

func Test737_CaptureFundsBuildInTheSameTurn(t *testing.T) {
	checkTurnCase(t, 737, 2, []expectedCommand{
		{"A pic - bel", accepted},
		{"A par", accepted},
	})
}

func Test738_EnemyRetreatCancelsApparentBuild(t *testing.T) {
	checkTurnCase(t, 738, 2, []expectedCommand{
		{"A pic - bel", accepted},
		{"A bur S A pic - bel", accepted},
		{"A R bel - hol", accepted},
		{"A par", noBuildDue},
	})
}

func Test739_ActivePlayerSubmitsOneOfTwoRequiredRemovals(t *testing.T) {
	checkTurnCase(t, 739, 2, []expectedCommand{
		{"D F eng", accepted},
	})
}
