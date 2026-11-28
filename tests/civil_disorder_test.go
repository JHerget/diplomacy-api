package tests

import "testing"

func Test801_CivilUnitsHoldIndependently(t *testing.T) {
	checkGameVariant(t, 801, 1)
}

func Test802_OtherPowersMayDefendCivilUnits(t *testing.T) {
	checkGameVariant(t, 802, 1)
}

func Test803_NoCivilBuilds(t *testing.T) {
	checkGameVariant(t, 803, 2)
}

func Test804_RandomRemovalAtHomeOrAbroad(t *testing.T) {
	checkGameVariant(t, 804, 2)
}

func Test805_RandomRemovalFromAMixedForce(t *testing.T) {
	checkGameVariant(t, 805, 2)
}

func Test806_RandomRemovalFromTwoArmies(t *testing.T) {
	checkGameVariant(t, 806, 2)
}

func Test807_EnglishArmiesAbroadDuringCivilDisorder(t *testing.T) {
	checkGameVariant(t, 807, 2)
}

func Test808_SixPlayerNeutralDislodged(t *testing.T) {
	checkGameVariant(t, 808, 1)
}

func Test809_NeutralReceivesOutsideHoldSupport(t *testing.T) {
	checkGameVariant(t, 809, 1)
}

func Test810_FivePlayerSecondNeutralPower(t *testing.T) {
	checkGameVariant(t, 810, 1)
}

func Test811_SharedPlayerSeparateNationalCounts(t *testing.T) {
	checkGameVariant(t, 811, 2)
}

func Test812_ItalianNeutralityBarsEntry(t *testing.T) {
	checkGameVariant(t, 812, 1)
}

func Test813_ItalianNeutralityEnds(t *testing.T) {
	checkGameVariant(t, 813, 3)
}
