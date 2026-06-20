package board

import (
	"diplomacy-api/internal/game"
	"diplomacy-api/internal/orders"
	"diplomacy-api/internal/utils/list"
)

type attack struct {
	Unit string
	Location string
	Support []string
}

type node struct {
	Providence game.Providence
	Attacks []attack
	UnresolvedCommand *game.MoveCommand
}

func ApplyTurn(mapBuf []byte, turn *game.Turn) ([]byte, error) {
	allOrders := list.Filter(turn.Orders, func(o *game.Order) bool {
		return o.PhaseId == turn.PhaseId
	})
	allCommands := orders.GetCommands(allOrders)
	validCommands := orders.ValidateCommands(allCommands)

	return mapBuf, nil
}
