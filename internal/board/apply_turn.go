package board

import (
	"diplomacy-api/internal/models"
	"diplomacy-api/internal/orders"
	"diplomacy-api/internal/utils"
)

type applicator struct {
	BoardMap map[string]*node
}

type node struct {
	Province          models.Province
	Attacks           []attack
	UnresolvedCommand *models.MoveCommand
}

type attack struct {
	UnitType   models.UnitType
	LocationID string
	Support    []string
}

func ApplyTurn(board []models.Province, turn *models.Turn) ([]models.Province, []models.CommandValidationResult, error) {
	allCommands := orders.GetCommands(turn.Orders)
	validCommands, validationResults := orders.ValidateCommands(allCommands, board)

	a := applicator{
		BoardMap: make(map[string]*node),
	}
	for _, p := range board {
		a.BoardMap[p.ID] = &node{
			Province:          p,
			Attacks:           []attack{},
			UnresolvedCommand: nil,
		}
	}

	isFallTurn := turn.TurnNumber%2 == 0

	utils.ForEach(validCommands.Move, a.ApplyMove)
	utils.ForEach(validCommands.Retreat, a.ApplyRetreat)
	utils.ForEach(validCommands.Support, a.ApplySupport)
	utils.ForEach(validCommands.Convoy, a.ApplyConvoy)
	utils.ForEach(validCommands.Disband, a.ApplyDisband)
	if isFallTurn {
		utils.ForEach(validCommands.Reinforce, a.ApplyReinforce)
	}

	return a.FinalizeState(isFallTurn), validationResults, nil
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

func (a applicator) FinalizeState(fallTurn bool) []models.Province {
	board := make([]models.Province, 0)

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

	if fallTurn {
		for index, p := range board {
			if p.SupplyCenter == nil || p.Unit == nil {
				continue
			}

			p.SupplyCenter.ControlledBy = &p.Unit.ControlledBy
			board[index] = p
		}
	}

	return board
}

func (a applicator) resolve(n *node) models.Province {
	minStrength := 0
	if n.Province.Unit != nil {
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

		if n.Province.ID != loc.Province.ID {
			n.Province.Unit = loc.Province.Unit
			loc.Province.Unit = nil
		}
	}

	return n.Province
}
