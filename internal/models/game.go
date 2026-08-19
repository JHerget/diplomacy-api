package models

import (
	"diplomacy-api/internal/utils"
	"errors"
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
