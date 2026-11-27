package tests

import (
	"testing"

	"diplomacy-api/internal/models"
)

func Test201_EqualAttacksLeaveBelgiumEmpty(t *testing.T) {
	commands := []expectedCommand{
		{"A pic - bel", accepted},
		{"A ruh - bel", accepted},
	}
	checkMovement(t, 201, commands)
	checkRetreatToBelgium(t, 201, commands)
}

func Test202_EqualAttacksLeaveTheBelgianDefenderInPlace(t *testing.T) {
	checkMovement(t, 202, []expectedCommand{
		{"A pic - bel", accepted},
		{"A ruh - bel", accepted},
	})
}

func Test203_ThreeEqualAttacksBounceInBelgium(t *testing.T) {
	commands := []expectedCommand{
		{"A pic - bel", accepted},
		{"A ruh - bel", accepted},
		{"F nth - bel", accepted},
	}
	checkMovement(t, 203, commands)
	checkRetreatToBelgium(t, 203, commands)
}

func Test204_ArmiesFollowSuccessfulDepartures(t *testing.T) {
	checkMovement(t, 204, []expectedCommand{
		{"A par - bur", accepted},
		{"A bur - mun", accepted},
		{"A mun - boh", accepted},
	})
}

func Test205_HoldingAtTheEndBlocksTheWholeChain(t *testing.T) {
	checkMovement(t, 205, []expectedCommand{
		{"A par - bur", accepted},
		{"A bur - mun", accepted},
		{"A mun - boh", accepted},
	})
}

func Test206_DirectExchangeBounces(t *testing.T) {
	checkMovement(t, 206, []expectedCommand{
		{"A par - bur", accepted},
		{"A bur - par", accepted},
	})
}

func Test207_StrongerHeadToHeadLetsThirdArmyFollow(t *testing.T) {
	checkMovement(t, 207, []expectedCommand{
		{"A par - bur", accepted},
		{"A bur - par", accepted},
		{"A pic S A par - bur", accepted},
		{"A gas - par", accepted},
	})
}

func Test208_ThreeArmiesRotate(t *testing.T) {
	checkMovement(t, 208, []expectedCommand{
		{"A par - bur", accepted},
		{"A bur - pic", accepted},
		{"A pic - par", accepted},
	})
}

func Test209_StandoffBlocksTheWholeRotation(t *testing.T) {
	checkMovement(t, 209, []expectedCommand{
		{"A par - bur", accepted},
		{"A bur - pic", accepted},
		{"A pic - par", accepted},
		{"A bre - pic", accepted},
	})
}

func Test210_ReorderingTheMovesPreservesTheRotation(t *testing.T) {
	checkMovement(t, 210, []expectedCommand{
		{"A pic - par", accepted},
		{"A bur - pic", accepted},
		{"A par - bur", accepted},
	})
}

func Test211_DislodgedArmyStillContestsMunich(t *testing.T) {
	commands := []expectedCommand{
		{"A par - bur", accepted},
		{"A pic S A par - bur", accepted},
		{"A bur - mun", accepted},
		{"A tyr - mun", accepted},
	}
	checkMovement(t, 211, commands)
	t.Run("Retreat to Munich is blocked by the standoff", func(t *testing.T) {
		f, players := loadFixture(t, 211)
		f.Orders = append(f.Orders, models.Order{PlayerName: "Germany", Value: "A R bur - mun"})

		checkSpringFixture(t, f, players, append(commands, expectedCommand{"A R bur - mun", standoff}))
	})

	t.Run("Burgundy can retreat to uncontested Munich", func(t *testing.T) {
		f, players := loadFixture(t, 211)
		// Only France attacks. Burgundy is still dislodged, but no one contests Munich.
		f.Orders = []models.Order{
			{PlayerName: "France", Value: "A par - bur, A pic S A par - bur"},
			{PlayerName: "Germany", Value: "A R bur - mun"},
		}
		placeArmy(t, f.OutputBoard, "Germany", "bur", "mun")

		checkSpringFixture(t, f, players, []expectedCommand{
			{"A par - bur", accepted},
			{"A pic S A par - bur", accepted},
			{"A R bur - mun", accepted},
		})
	})
}

func Test212_DislodgedArmyCannotContestItsDislodgersOrigin(t *testing.T) {
	checkMovement(t, 212, []expectedCommand{
		{"A bul - rum", accepted},
		{"F bla S A bul - rum", accepted},
		{"A rum - bul", accepted},
		{"A gre S A rum - bul", accepted},
		{"A ser S A rum - bul", accepted},
		{"A sev - rum", accepted},
	})
}

// In cases 201, 203 and 312 Belgium starts and ends empty. Dislodge a fleet
// in the Channel so its retreat reveals whether Belgium was contested.
func checkRetreatToBelgium(t *testing.T, id int, attacks []expectedCommand) {
	t.Helper()
	dislodgeFleet := func(f *fixture) {
		placeUnit(t, f.InputBoard, models.UnitFleet, "England", "eng", "eng")
		placeUnit(t, f.InputBoard, models.UnitFleet, "Italy", "bre", "bre")
		placeUnit(t, f.InputBoard, models.UnitFleet, "Italy", "mid", "mid")
		placeUnit(t, f.OutputBoard, models.UnitFleet, "Italy", "bre", "eng")
		placeUnit(t, f.OutputBoard, models.UnitFleet, "Italy", "mid", "mid")
		f.Orders = append(f.Orders,
			models.Order{PlayerName: "Italy", Value: "F bre - eng, F mid S F bre - eng"},
			models.Order{PlayerName: "England", Value: "F R eng - bel"},
		)
	}

	t.Run("Standoff in Belgium prevents the retreat", func(t *testing.T) {
		f, players := loadFixture(t, id)
		dislodgeFleet(&f)
		checkSpringFixture(t, f, players, append(attacks,
			expectedCommand{"F bre - eng", accepted},
			expectedCommand{"F mid S F bre - eng", accepted},
			expectedCommand{"F R eng - bel", standoff},
		))
	})

	t.Run("Uncontested Belgium accepts the retreat", func(t *testing.T) {
		f, players := loadFixture(t, id)
		f.Orders = nil // The original units hold.
		dislodgeFleet(&f)
		placeUnit(t, f.OutputBoard, models.UnitFleet, "England", "eng", "bel")
		checkSpringFixture(t, f, players, []expectedCommand{
			{"F bre - eng", accepted},
			{"F mid S F bre - eng", accepted},
			{"F R eng - bel", accepted},
		})
	})
}
