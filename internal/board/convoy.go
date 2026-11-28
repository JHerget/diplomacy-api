package board

import (
	"fmt"
	"maps"
	"slices"
	"strings"

	"diplomacy-api/internal/models"
	"diplomacy-api/internal/utils"
)

// A landing can cut support for a fleet attack, and a dislodged carrier can
// suppress the landing. Replay movement from the same position until fleet
// survival agrees with the routes used in that pass.
func resolveMovement(board []models.Province, commands models.Commands) (*applicator, error) {
	disrupted := map[string]bool{}
	tried := map[string]bool{}
	for {
		key := strings.Join(slices.Sorted(maps.Keys(disrupted)), ",")
		if tried[key] {
			return nil, fmt.Errorf("convoy dependencies did not settle for carriers %q", key)
		}
		tried[key] = true

		a := &applicator{
			BoardMap:         make(map[string]*node, len(board)),
			ConvoyCarriers:   map[string]*node{},
			ReceivedCommands: map[string]bool{},
			IgnoredCommands:  map[int]string{},
			PlayerCensus:     map[string]*census{},
			Disrupted:        disrupted,
		}
		for _, province := range board {
			if province.Unit != nil {
				unit := *province.Unit
				province.Unit = &unit
			}
			a.BoardMap[province.ID] = &node{Province: province}
		}
		utils.ForEach(commands.Convoy, a.ApplyConvoy)
		utils.ForEach(commands.Hold, a.ApplyHold)
		utils.ForEach(commands.Move, a.ApplyMove)
		utils.ForEach(commands.Support, a.ApplySupport)
		a.ResolveAttacks()

		next := map[string]bool{}
		for id, n := range a.BoardMap {
			if n.Convoy != nil && n.AttackResult.DislodgedUnit != nil {
				next[id] = true
			}
		}
		if maps.Equal(disrupted, next) {
			return a, nil
		}
		disrupted = next
	}
}

// convoyMove returns the selected passenger move from its carrier's order.
// Only armies that selected a convoy are indexed; an unsolicited foreign
// convoy does not change an adjacent army move into a convoyed move.
func (a *applicator) convoyMove(origin *node) *models.MoveCommand {
	if carrier := a.ConvoyCarriers[origin.Province.ID]; carrier != nil {
		return &carrier.Convoy.Move
	}
	return nil
}

// Search simple paths through matching sea fleets. Accept the destination
// only after traversing at least one carrier, even for adjacent land provinces.
func (a *applicator) convoyRoute(move models.MoveCommand, disrupted map[string]bool, exclude string, requireOwn bool) bool {
	origin, dest := a.BoardMap[move.Location.ID], a.BoardMap[move.Destination.ID]
	if !origin.Province.IsCoastal() || !dest.Province.IsCoastal() {
		return false
	}
	visited := map[string]bool{}
	var walk func(*node, bool) bool
	walk = func(n *node, hasOwn bool) bool {
		for _, id := range n.Province.Routes {
			if id == dest.Province.ID && n != origin && (!requireOwn || hasOwn) {
				return true
			}
			fleet := a.BoardMap[id]
			if visited[id] || disrupted[id] || id == exclude || fleet == nil ||
				fleet.Convoy == nil || !fleet.Convoy.Move.Equals(move) {
				continue
			}
			visited[id] = true
			if walk(fleet, hasOwn || fleet.Convoy.PlayerName == move.PlayerName) {
				return true
			}
			delete(visited, id)
		}
		return false
	}
	return walk(origin, false)
}
