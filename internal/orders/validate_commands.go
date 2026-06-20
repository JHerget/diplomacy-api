package orders

import (
	"diplomacy-api/internal/game"
	"diplomacy-api/internal/utils/list"
)

type providenceMap map[string]game.Providence

func ValidateCommands(commands game.Commands, board []game.Providence) game.Commands {
	boardMap := providenceMap{}

	for _, p := range board {
		boardMap[p.Id] = p
	}

	return game.Commands{
		Hold: list.Filter(commands.Hold, isValidHold),
		Move: list.Filter(commands.Move, isValidMove),
		Retreat: list.Filter(commands.Retreat, isValidRetreat),
		Support: list.Filter(commands.Support, isValidSupport),
		Convoy: list.Filter(commands.Convoy, isValidConvoy),
		Reinforce: list.Filter(commands.Reinforce, isValidReinforce),
		Disband: list.Filter(commands.Disband, isValidDisband),
	}
}

func isValidHold(hold *game.MoveCommand) bool {
	return true
}

func isValidMove(move *game.MoveCommand) bool {
	return true
}

func isValidRetreat(retreat *game.MoveCommand) bool {
	return true
}

func isValidSupport(support *game.SupportCommand) bool {
	return true
}

func isValidConvoy(convoy *game.SupportCommand) bool {
	return true
}

func isValidReinforce(reinforce *game.AdjustCommand) bool {
	return true
}

func isValidDisband(disband *game.AdjustCommand) bool {
	return true
}
