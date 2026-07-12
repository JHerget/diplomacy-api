package board

import (
	"diplomacy-api/internal/models"
	"diplomacy-api/internal/utils"
)

type applicator struct {
	BoardMap map[string]*node
}

type node struct {
	Providence        models.Providence
	Attacks           []attack
	UnresolvedCommand *models.MoveCommand
}

type attack struct {
	UnitType   models.UnitType
	LocationID string
	Support    []string
}

func ApplyTurn(board []models.Providence, commands models.Commands) ([]models.Providence, error) {
	a := applicator{
		BoardMap: make(map[string]*node),
	}
	for _, p := range board {
		a.BoardMap[p.ID] = &node{
			Providence:        p,
			Attacks:           []attack{},
			UnresolvedCommand: nil,
		}
	}

	utils.ForEach(commands.Move, a.ApplyMove)
	utils.ForEach(commands.Retreat, a.ApplyRetreat)
	utils.ForEach(commands.Support, a.ApplySupport)
	utils.ForEach(commands.Convoy, a.ApplyConvoy)
	utils.ForEach(commands.Reinforce, a.ApplyReinforce)
	utils.ForEach(commands.Disband, a.ApplyDisband)

	return a.FinalizeState(), nil
}

func (a applicator) ApplyMove(move models.MoveCommand) {
	loc, ok := a.BoardMap[move.Location.ID]
	if !ok {
		return
	}

	dest, ok := a.BoardMap[move.Destination.ID]
	if !ok {
		return
	}

	loc.UnresolvedCommand = &move
	dest.Attacks = append(dest.Attacks, attack{
		UnitType:   move.UnitType,
		LocationID: move.Location.ID,
		Support:    nil,
	})
}

func (a applicator) ApplyRetreat(retreat models.MoveCommand) {}

func (a applicator) ApplySupport(support models.SupportCommand) {
	dest, ok := a.BoardMap[support.Move.Destination.ID]
	if !ok {
		return
	}

	att, ok := utils.Find(dest.Attacks, func(a *attack) bool {
		return a.LocationID == support.Move.Location.ID
	})
	if !ok {
		return
	}

	att.Support = append(att.Support, support.Location.ID)
}

func (a applicator) ApplyConvoy(convoy models.SupportCommand)      {}
func (a applicator) ApplyReinforce(reinforce models.AdjustCommand) {}
func (a applicator) ApplyDisband(disband models.AdjustCommand)     {}

func (a applicator) FinalizeState() []models.Providence {
	board := make([]models.Providence, 0)

	for _, n := range a.BoardMap {
		if n.UnresolvedCommand == nil {
			board = append(board, a.resolve(n))
			continue
		}

		stack := []*node{n}

		for len(stack) > 0 {
			top := stack[len(stack)-1]
			dest := a.BoardMap[top.UnresolvedCommand.Destination.ID]

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

func (a applicator) resolve(n *node) models.Providence {
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
		loc := a.BoardMap[att.LocationID]

		if n.Providence.ID != loc.Providence.ID {
			n.Providence.Unit = loc.Providence.Unit
			loc.Providence.Unit = nil
		}
	}

	return n.Providence
}
