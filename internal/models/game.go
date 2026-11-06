package models

import (
	"crypto/rand"
	"diplomacy-api/internal/utils"
	"fmt"
	"math/big"
	"strings"
	"time"
)

type Game struct {
	ID                string     `json:"id" bson:"_id,omitempty"`
	ExternalID        *string    `json:"externalId" bson:"externalId"`
	OwnerID           string     `json:"ownerID" bson:"ownerID"`
	Map               MapSummary `json:"map" bson:"map"`
	Board             []Province `json:"board" bson:"board"`
	Players           []Player   `json:"players" bson:"players"`
	Turns             []Turn     `json:"turns" bson:"turns"`
	DaysPerTurn       int        `json:"daysPerTurn" bson:"daysPerTurn"`
	TurnStartHour     int        `json:"turnStartHour" bson:"turnStartHour"`
	Timezone          int        `json:"timezone" bson:"timezone"`
	StartDate         int        `json:"startDate" bson:"startDate"`
	EndDate           int        `json:"endDate" bson:"endDate"`
	InProgress        bool       `json:"inProgress" bson:"inProgress"`
	IsDeleted         bool       `json:"isDeleted" bson:"isDeleted"`
	RandomAssignments bool       `json:"randomAssignments" bson:"randomAssignments"`
	PlayerCount       int        `json:"playerCount" bson:"playerCount"`
}

func (g *Game) Valid() error {
	if strings.TrimSpace(g.OwnerID) == "" {
		return fmt.Errorf("Missing owner id.")
	}
	if strings.TrimSpace(g.Map.ID) == "" {
		return fmt.Errorf("Missing map id.")
	}

	if g.PlayerCount < 2 || g.PlayerCount > 7 {
		return fmt.Errorf("player count must be greater than or equal to 2 and less than or equal to 7")
	}

	if g.DaysPerTurn <= 0 {
		return fmt.Errorf("Days per turn must be greater than 0.")
	}

	if g.TurnStartHour < 0 || g.TurnStartHour > 23 {
		return fmt.Errorf("Turn start hour must be greater than or equal to 0 and less than or equal to 23.")
	}

	if g.StartDate < 0 {
		return fmt.Errorf("Start date must be an epoch timestamp.")
	}
	if g.EndDate < 0 {
		return fmt.Errorf("End date must be an epoch timestamp.")
	}
	for _, p := range g.Board {
		if err := p.Valid(); err != nil {
			return err
		}
	}
	for _, p := range g.Players {
		if err := p.Valid(); err != nil {
			return err
		}
	}
	for _, t := range g.Turns {
		if err := t.Valid(); err != nil {
			return err
		}
	}

	return nil
}

func (g *Game) FindTurn(turnID string) (*Turn, bool) {
	return utils.Find(g.Turns, func(t *Turn) bool {
		return t.ID == turnID
	})
}

func (g *Game) FindPlayer(playerID string) (*Player, bool) {
	return utils.Find(g.Players, func(p *Player) bool {
		return p.ID == playerID
	})
}

func (g *Game) NextTurnStartDate() int {
	sourceDate := g.StartDate
	if len(g.Turns) > 0 {
		sourceDate = g.Turns[len(g.Turns)-1].EndDate
	}

	location := time.FixedZone("game", g.Timezone*60*60)
	t := time.Unix(int64(sourceDate), 0).In(location)
	start := time.Date(t.Year(), t.Month(), t.Day(), g.TurnStartHour, 0, 0, 0, location)

	return int(start.Unix())
}

func (g *Game) CurrentTurn() *Turn {
	if len(g.Turns) == 0 {
		return nil
	}

	latestTurn := len(g.Turns)
	for _, turn := range g.Turns {
		if turn.TurnNumber == latestTurn {
			return &turn
		}
	}

	return nil
}

func (g *Game) NewTurn() (*Turn, error) {
	id, err := utils.RandomID()
	if err != nil {
		return nil, err
	}

	startDate := g.NextTurnStartDate()
	endDate := int(time.Unix(int64(startDate), 0).UTC().AddDate(0, 0, g.DaysPerTurn).Unix())
	turn := Turn{
		ID:         id,
		Orders:     []Order{},
		TurnNumber: len(g.Turns) + 1,
		StartDate:  startDate,
		EndDate:    endDate,
	}
	g.Turns = append(g.Turns, turn)

	if err := g.Valid(); err != nil {
		return nil, err
	}

	return &turn, nil
}

func (g *Game) RemoveTurn(turnID string) error {
	index := -1
	for i := range g.Turns {
		if g.Turns[i].ID == turnID {
			index = i
			break
		}
	}
	if index == -1 {
		return fmt.Errorf("invalid turn id '%s'", turnID)
	}
	g.Turns = append(g.Turns[:index], g.Turns[index+1:]...)

	if err := g.Valid(); err != nil {
		return err
	}

	return nil
}

func (g *Game) IsFinished() bool {
	winningScore := 18
	playerScores := make(map[string]int)

	if g.PlayerCount == 2 {
		winningScore = 24
	}

	for _, p := range g.Board {
		if p.SupplyCenter == nil {
			continue
		}
		if p.SupplyCenter.ControlledBy == nil {
			continue
		}

		playerScores[*p.SupplyCenter.ControlledBy]++
	}

	for _, score := range playerScores {
		if score >= winningScore {
			return true
		}
	}

	return false
}

func (g *Game) UpdatePlayerStatuses() {
	playerUnits := make(map[string]int)

	for _, p := range g.Board {
		if p.Unit == nil {
			continue
		}

		playerUnits[p.Unit.ControlledBy]++
	}

	for index, p := range g.Players {
		if playerUnits[p.Name] == 0 {
			p.IsPlaying = false
			g.Players[index] = p
		}
	}
}

func (g *Game) AssignSpecificPlayers(userID string, playerNames []string) ([]*Player, error) {
	players := make(map[string]int)
	for i, player := range g.Players {
		if player.UserID != nil && *player.UserID == userID {
			return nil, fmt.Errorf("User is already assigned to a player.")
		}

		players[player.Name] = i
	}

	results := make([]*Player, 0)
	for _, name := range playerNames {
		index := players[name]
		player := &g.Players[index]

		if player.UserID != nil {
			return nil, fmt.Errorf("%s is already assigned.", name)
		}
		player.UserID = &userID
		player.IsPlaying = true

		results = append(results, player)
	}

	return results, nil
}

func (g *Game) AssignRandomPlayers(userID string) ([]*Player, error) {
	unassignedPlayers := make(map[string]int)
	for i, player := range g.Players {
		if player.UserID != nil && *player.UserID == userID {
			return nil, fmt.Errorf("User is already assigned to a player.")
		}
		if player.UserID == nil {
			unassignedPlayers[player.Name] = i + 1
			continue
		}
	}
	if len(unassignedPlayers) == 0 {
		return nil, fmt.Errorf("The game is full.")
	}

	playerCombo := PlayerCombos[g.PlayerCount]
	unassignedCombos := make([][]string, 0)
	for _, c := range playerCombo {
		names := strings.Split(c, "/")
		skipCombo := false

		for _, name := range names {
			if unassignedPlayers[name] < 1 {
				skipCombo = true
				break
			}
		}

		if skipCombo {
			continue
		}

		unassignedCombos = append(unassignedCombos, names)
	}

	selected, err := rand.Int(rand.Reader, big.NewInt(int64(len(unassignedCombos))))
	if err != nil {
		return nil, err
	}

	results := make([]*Player, 0)
	for _, name := range unassignedCombos[selected.Int64()] {
		player := &g.Players[unassignedPlayers[name]-1]
		player.UserID = &userID
		player.IsPlaying = true

		results = append(results, player)
	}

	return results, nil
}
