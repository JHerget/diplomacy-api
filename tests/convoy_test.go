package tests

import "testing"

func Test401_SingleFleetCrossing(t *testing.T) {
	checkTurnCase(t, 401, 1, []expectedCommand{
		{"A lon - bel", accepted},
		{"F eng C A lon - bel", accepted},
	})
}

func Test402_MissingTransport(t *testing.T) {
	checkTurnCase(t, 402, 1, []expectedCommand{
		{"A lon - bel", noConvoy},
	})
}

func Test403_WrongDestinationOnConvoyOrder(t *testing.T) {
	checkTurnCase(t, 403, 1, []expectedCommand{
		{"A lon - bel", noConvoy},
		{"F eng C A lon - pic", accepted},
	})
}

func Test404_WrongArmyOnConvoyOrder(t *testing.T) {
	checkTurnCase(t, 404, 1, []expectedCommand{
		{"A lon - bel", noConvoy},
		{"F eng C A wal - bel", accepted},
	})
}

func Test405_MultinationalChain(t *testing.T) {
	checkTurnCase(t, 405, 1, []expectedCommand{
		{"A lon - tun", accepted},
		{"F eng C A lon - tun", accepted},
		{"F mid C A lon - tun", accepted},
		{"F wes C A lon - tun", accepted},
	})
}

func Test406_MissingMiddleLink(t *testing.T) {
	checkTurnCase(t, 406, 1, []expectedCommand{
		{"A lon - tun", noConvoy},
		{"F eng C A lon - tun", accepted},
		{"F wes C A lon - tun", accepted},
	})
}

func Test407_CoastalFleetCannotBridgeTwoSeas(t *testing.T) {
	checkTurnCase(t, 407, 1, []expectedCommand{
		{"A ank - gre", invalid},
		{"F bla C A ank - gre", invalid},
		{"F con C A ank - gre", invalid},
		{"F aeg C A ank - gre", invalid},
	})
}

func Test408_FleetAsPassenger(t *testing.T) {
	checkTurnCase(t, 408, 1, []expectedCommand{
		{"F lon - bel", invalid},
		{"F eng C F lon - bel", invalid},
	})
}

func Test409_InlandConvoyEndpoint(t *testing.T) {
	checkTurnCase(t, 409, 1, []expectedCommand{
		{"A lon - par", invalid},
		{"F eng C A lon - par", invalid},
	})
}

func Test410_ConvoyAddsNoStrength(t *testing.T) {
	checkTurnCase(t, 410, 1, []expectedCommand{
		{"A lon - bel", accepted},
		{"F eng C A lon - bel", accepted},
	})
}

func Test411_SupportedLanding(t *testing.T) {
	checkTurnCase(t, 411, 1, []expectedCommand{
		{"A lon - bel", accepted},
		{"F eng C A lon - bel", accepted},
		{"A pic S A lon - bel", accepted},
	})
}

func Test412_MereAttackOnTransport(t *testing.T) {
	checkTurnCase(t, 412, 1, []expectedCommand{
		{"A lon - bel", accepted},
		{"F eng C A lon - bel", accepted},
		{"A pic S A lon - bel", accepted},
		{"F bre - eng", accepted},
	})
}

func Test413_TransportDislodgedSupportedLandingFails(t *testing.T) {
	checkTurnCase(t, 413, 1, []expectedCommand{
		{"A lon - bel", accepted},
		{"F eng C A lon - bel", accepted},
		{"A pic S A lon - bel", accepted},
		{"F bre - eng", accepted},
		{"F iri S F bre - eng", accepted},
	})
}

func Test414_AlternateRouteSurvives(t *testing.T) {
	checkTurnCase(t, 414, 1, []expectedCommand{
		{"A lon - bel", accepted},
		{"F eng C A lon - bel", accepted},
		{"F nth C A lon - bel", accepted},
		{"A pic S A lon - bel", accepted},
		{"F bre - eng", accepted},
		{"F iri S F bre - eng", accepted},
	})
}

func Test415_EveryAlternateRouteBroken(t *testing.T) {
	checkTurnCase(t, 415, 1, []expectedCommand{
		{"A lon - bel", accepted},
		{"F eng C A lon - bel", accepted},
		{"F nth C A lon - bel", accepted},
		{"F bre - eng", accepted},
		{"F iri S F bre - eng", accepted},
		{"F hel - nth", accepted},
		{"F den S F hel - nth", accepted},
	})
}

func Test416_BouncedPassengerDislodgedAtHome(t *testing.T) {
	checkTurnCase(t, 416, 1, []expectedCommand{
		{"A lon - bel", accepted},
		{"F eng C A lon - bel", accepted},
		{"A yor - lon", accepted},
		{"A wal S A yor - lon", accepted},
	})
}

func Test417_FleetReceivesHoldSupportWhileConvoying(t *testing.T) {
	checkTurnCase(t, 417, 1, []expectedCommand{
		{"A lon - bel", accepted},
		{"F eng C A lon - bel", accepted},
		{"F mid S F eng - H", accepted},
		{"F bre - eng", accepted},
		{"F iri S F bre - eng", accepted},
	})
}

func Test418_OneFleetOrderedToFerryTwoArmies(t *testing.T) {
	checkTurnCase(t, 418, 1, []expectedCommand{
		{"A lon - bel", accepted},
		{"A wal - pic", noConvoy},
		{"F eng C A lon - bel", accepted},
		{"F eng C A wal - pic", duplicate},
	})
}

func Test419_BothSidesExchangeByConvoy(t *testing.T) {
	checkTurnCase(t, 419, 1, []expectedCommand{
		{"A lon - bel", accepted},
		{"A bel - lon", accepted},
		{"F nth C A lon - bel", accepted},
		{"F eng C A bel - lon", accepted},
	})
}

func Test420_AdjacentExchangeOwnTransport(t *testing.T) {
	checkTurnCase(t, 420, 1, []expectedCommand{
		{"A lon - wal", accepted},
		{"A wal - lon", accepted},
		{"F eng C A lon - wal", accepted},
	})
}

func Test421_AdjacentExchangeForeignTransport(t *testing.T) {
	checkTurnCase(t, 421, 1, []expectedCommand{
		{"A lon - wal", accepted},
		{"A wal - lon", accepted},
		{"F eng C A lon - wal", accepted},
	})
}

func Test422_ExplicitForeignConvoy(t *testing.T) {
	checkTurnCase(t, 422, 1, []expectedCommand{
		{"A lon - wal via convoy", accepted},
		{"A wal - lon", accepted},
		{"F eng C A lon - wal", accepted},
	})
}

func Test423_DisruptionCannotRestoreDiscardedLandRoute(t *testing.T) {
	checkTurnCase(t, 423, 1, []expectedCommand{
		{"A lon - wal", accepted},
		{"F eng C A lon - wal", accepted},
		{"F bre - eng", accepted},
		{"F iri S F bre - eng", accepted},
	})
}

func Test424_NecessaryRouteUnderAttackFromDestination(t *testing.T) {
	checkTurnCase(t, 424, 1, []expectedCommand{
		{"A tun - nap", accepted},
		{"F tyn C A tun - nap", accepted},
		{"F rom - tyn", accepted},
		{"F nap S F rom - tyn", accepted},
	})
}

func Test425_IndependentRouteChangesTheCut(t *testing.T) {
	checkTurnCase(t, 425, 1, []expectedCommand{
		{"A tun - nap", accepted},
		{"F tyn C A tun - nap", accepted},
		{"F ion C A tun - nap", accepted},
		{"F rom - tyn", accepted},
		{"F nap S F rom - tyn", supportCutFrom("tun")},
	})
}

func Test426_IndependentRouteWithLandingSupport(t *testing.T) {
	checkTurnCase(t, 426, 1, []expectedCommand{
		{"A tun - nap", accepted},
		{"F tyn C A tun - nap", accepted},
		{"F ion C A tun - nap", accepted},
		{"F rom - tyn", accepted},
		{"F nap S F rom - tyn", supportCutFrom("tun")},
		{"A apu S A tun - nap", accepted},
	})
}

func Test427_InlandOrigin(t *testing.T) {
	checkTurnCase(t, 427, 1, []expectedCommand{
		{"A par - lon", invalid},
		{"F eng C A par - lon", invalid},
	})
}

func Test428_ArmyCannotTransportAnotherArmy(t *testing.T) {
	checkTurnCase(t, 428, 1, []expectedCommand{
		{"A lon - bel", noConvoy},
		{"A wal C A lon - bel", invalid},
	})
}

func Test429_SupportCannotBeTransported(t *testing.T) {
	checkTurnCase(t, 429, 1, []expectedCommand{
		{"A pic - bre", accepted},
		{"A lon S A pic - bre", invalid},
		{"F eng C A lon S A pic - bre", invalid},
	})
}

func Test430_DisruptedPassengerCanBeDislodged(t *testing.T) {
	checkTurnCase(t, 430, 1, []expectedCommand{
		{"A lon - bel", accepted},
		{"F eng C A lon - bel", accepted},
		{"F bre - eng", accepted},
		{"F iri S F bre - eng", accepted},
		{"A yor - lon", accepted},
		{"A wal S A yor - lon", accepted},
	})
}

func Test431_ConvoyBeforeFleetMove(t *testing.T) {
	checkTurnCase(t, 431, 1, []expectedCommand{
		{"A lon - bel", accepted},
		{"F eng C A lon - bel", accepted},
		{"F eng - bre", duplicate},
	})
}

func Test432_FleetMoveBeforeConvoy(t *testing.T) {
	checkTurnCase(t, 432, 1, []expectedCommand{
		{"A lon - bel", accepted},
		{"F eng - bre", duplicate},
		{"F eng C A lon - bel", accepted},
	})
}
