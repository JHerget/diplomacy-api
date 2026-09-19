package models

import (
	"diplomacy-api/internal/utils"
	"time"
)

type Turn struct {
	ID         string  `json:"id" bson:"id"`
	PhaseID    string  `json:"phaseID" bson:"phaseID"`
	Orders     []Order `json:"orders" bson:"orders"`
	TurnNumber int     `json:"turnNumber" bson:"turnNumber"`
	StartDate  int     `json:"startDate" bson:"startDate"`
	EndDate    int     `json:"endDate" bson:"endDate"`
}

func (t *Turn) Valid() error {
	for _, o := range t.Orders {
		if err := o.Valid(); err != nil {
			return err
		}
	}

	return nil
}

func (t *Turn) FindOrder(orderID string) (*Order, bool) {
	return utils.Find(t.Orders, func(o *Order) bool {
		return o.ID == orderID
	})
}

func (t *Turn) FindPlayerOrder(playerName string) (*Order, bool) {
	return utils.Find(t.Orders, func(o *Order) bool {
		return o.PlayerName == playerName
	})
}

func (t *Turn) IsFinished(players []Player) bool {
	if time.Now().UTC().Unix() >= int64(t.EndDate) {
		return true
	}

	submitted := map[string]bool{}
	for _, order := range t.Orders {
		submitted[order.PlayerName] = true
	}

	for _, player := range players {
		if player.IsPlaying && !submitted[player.Name] {
			return false
		}
	}

	return true
}
