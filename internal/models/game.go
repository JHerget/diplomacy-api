package models

import (
	"diplomacy-api/internal/utils"
	"errors"
	"fmt"
	"strings"
	"time"
)

type Game struct {
	ID            string       `json:"id" bson:"_id,omitempty"`
	ExternalID    *string      `json:"externalId" bson:"externalId"`
	OwnerID       string       `json:"ownerID" bson:"ownerID"`
	Map           MapSummary   `json:"map" bson:"map"`
	Board         []Providence `json:"board" bson:"board"`
	Players       []Player     `json:"players" bson:"players"`
	Turns         []Turn       `json:"turns" bson:"turns"`
	DaysPerTurn   int          `json:"daysPerTurn" bson:"daysPerTurn"`
	TurnStartHour int          `json:"turnStartHour" bson:"turnStartHour"`
	Timezone      int          `json:"timezone" bson:"timezone"`
	StartDate     int          `json:"startDate" bson:"startDate"`
	EndDate       int          `json:"endDate" bson:"endDate"`
	InProgress    bool         `json:"inProgress" bson:"inProgress"`
	IsDeleted     bool         `json:"isDeleted" bson:"isDeleted"`
}

func (g *Game) Valid() error {
	if strings.TrimSpace(g.OwnerID) == "" {
		return errors.New("Missing owner id.")
	}
	if strings.TrimSpace(g.Map.ID) == "" {
		return errors.New("Missing map id.")
	}

	if g.DaysPerTurn <= 0 {
		return errors.New("Days per turn must be greater than 0.")
	}

	if g.TurnStartHour < 0 || g.TurnStartHour > 23 {
		return errors.New("Turn start hour must be greater than or equal to 0 and less than or equal to 23.")
	}

	if g.StartDate < 0 {
		return errors.New("Start date must be an epoch timestamp.")
	}
	if g.EndDate < 0 {
		return errors.New("End date must be an epoch timestamp.")
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

	var latestTurn *Turn
	for _, turn := range g.Turns {
		if latestTurn == nil {
			latestTurn = &turn
			continue
		}

		if turn.TurnNumber > latestTurn.TurnNumber {
			latestTurn = &turn
		}
	}

	return latestTurn
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
		PhaseID:    "",
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
