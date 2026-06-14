package board

import (
	"diplomacy-api/internal/game"
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
	orders := list.Filter(turn.Orders, func(o *game.Order) bool {
		return o.PhaseId == turn.PhaseId
	})

	return mapBuf, nil
}
