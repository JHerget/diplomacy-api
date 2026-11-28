package tests

import "testing"

func Test701_CivilUnitsHoldIndependently(t *testing.T) {
	checkGameVariant(t, 701, 1)
}

func Test702_OtherPowersMayDefendCivilUnits(t *testing.T) {
	checkGameVariant(t, 702, 1)
}

func Test703_NoCivilBuilds(t *testing.T) {
	checkGameVariant(t, 703, 2)
}

func Test704_RandomRemovalAtHomeOrAbroad(t *testing.T) {
	checkGameVariant(t, 704, 2)
}

func Test705_RandomRemovalFromAMixedForce(t *testing.T) {
	checkGameVariant(t, 705, 2)
}

func Test706_RandomRemovalFromTwoArmies(t *testing.T) {
	checkGameVariant(t, 706, 2)
}

func Test707_EnglishArmiesAbroadDuringCivilDisorder(t *testing.T) {
	checkGameVariant(t, 707, 2)
}

func Test708_SixPlayerNeutralDislodged(t *testing.T) {
	checkGameVariant(t, 708, 1)
}

func Test709_NeutralReceivesOutsideHoldSupport(t *testing.T) {
	checkGameVariant(t, 709, 1)
}

func Test710_FivePlayerSecondNeutralPower(t *testing.T) {
	checkGameVariant(t, 710, 1)
}

func Test711_SharedPlayerSeparateNationalCounts(t *testing.T) {
	checkGameVariant(t, 711, 2)
}

func Test712_ItalianNeutralityBarsEntry(t *testing.T) {
	checkGameVariant(t, 712, 1)
}

func Test713_ItalianNeutralityEnds(t *testing.T) {
	checkGameVariant(t, 713, 3)
}
