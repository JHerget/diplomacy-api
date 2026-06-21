package board

import (
	"diplomacy-api/internal/game"
	"diplomacy-api/internal/orders"
	"diplomacy-api/internal/utils"
)

type applicator struct {
	BoardMap map[string]*node
}

type node struct {
	Providence        game.Providence
	Attacks           []attack
	UnresolvedCommand *game.MoveCommand
}

type attack struct {
	UnitType   game.UnitType
	LocationId string
	Support    []string
}

func ApplyTurn(board []game.Providence, turn *game.Turn) ([]game.Providence, error) {
	allOrders := utils.Filter(turn.Orders, func(o *game.Order) bool {
		return o.PhaseId == turn.PhaseId
	})
	allCommands := orders.GetCommands(allOrders)
	validCommands := orders.ValidateCommands(allCommands, board)

	a := applicator{
		BoardMap: make(map[string]*node),
	}
	for _, p := range board {
		a.BoardMap[p.Id] = &node{
			Providence:        p,
			Attacks:           []attack{},
			UnresolvedCommand: nil,
		}
	}

	utils.ForEach(validCommands.Move, a.ApplyMove)
	utils.ForEach(validCommands.Retreat, a.ApplyRetreat)
	utils.ForEach(validCommands.Support, a.ApplySupport)
	utils.ForEach(validCommands.Convoy, a.ApplyConvoy)

	if turn.TurnNumber%2 == 0 {
		utils.ForEach(validCommands.Reinforce, a.ApplyReinforce)
		utils.ForEach(validCommands.Disband, a.ApplyDisband)
	}

	return a.FinalizeState(), nil
}

func (a applicator) ApplyMove(move game.MoveCommand) {
	loc, ok := a.BoardMap[move.Location.Id]
	if !ok {
		return
	}

	dest, ok := a.BoardMap[move.Destination.Id]
	if !ok {
		return
	}

	loc.UnresolvedCommand = &move
	dest.Attacks = append(dest.Attacks, attack{
		UnitType:   move.UnitType,
		LocationId: move.Location.Id,
		Support:    nil,
	})
}

func (a applicator) ApplyRetreat(retreat game.MoveCommand) {}

func (a applicator) ApplySupport(support game.SupportCommand) {
	dest, ok := a.BoardMap[support.Move.Destination.Id]
	if !ok {
		return
	}

	att, ok := utils.Find(dest.Attacks, func(a *attack) bool {
		return a.LocationId == support.Move.Location.Id
	})
	if !ok {
		return
	}

	att.Support = append(att.Support, support.Location.Id)
}

func (a applicator) ApplyConvoy(convoy game.SupportCommand)      {}
func (a applicator) ApplyReinforce(reinforce game.AdjustCommand) {}
func (a applicator) ApplyDisband(disband game.AdjustCommand)     {}

func (a applicator) FinalizeState() []game.Providence {
	board := make([]game.Providence, 0)

	for _, n := range a.BoardMap {
		if n.UnresolvedCommand == nil {
			board = append(board, a.resolve(n))
			continue
		}

		stack := []*node{n}

		for len(stack) > 0 {
			top := stack[len(stack)-1]
			dest := a.BoardMap[top.UnresolvedCommand.Destination.Id]

			if dest.UnresolvedCommand != nil {
				stack = append(stack, dest)
				continue
			}

			_ = a.resolve(dest)
			board = append(board, a.resolve(top))
			stack = stack[:len(stack)-1]
		}
	}

	return board
}

func (a applicator) resolve(n *node) game.Providence {
	minStrength := 0
	if n.Providence.Unit != nil {
		minStrength = 1
	}

	maxStrength := 0
	for _, att := range n.Attacks {
		strength := len(att.Support) + 1

		if strength > maxStrength {
			maxStrength = strength
		}
	}

	strongestAttacks := utils.Filter(n.Attacks, func(a *attack) bool {
		return len(a.Support)+1 >= maxStrength
	})

	n.UnresolvedCommand = nil
	n.Attacks = nil

	if len(strongestAttacks) == 1 && maxStrength > minStrength {
		att := strongestAttacks[0]
		loc := a.BoardMap[att.LocationId]

		if n.Providence.Id != loc.Providence.Id {
			n.Providence.Unit = loc.Providence.Unit
			loc.Providence.Unit = nil
		}
	}

	return n.Providence
}
