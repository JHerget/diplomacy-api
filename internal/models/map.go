package models

import (
	"errors"
	"strings"
)

type Map struct {
	ID          string       `json:"id" bson:"_id,omitempty"`
	Name        string       `json:"name" bson:"name"`
	Players     []Player     `json:"players" bson:"players"`
	Providences []Providence `json:"providences" bson:"providences"`
	IsDeleted   bool         `json:"isDeleted" bson:"isDeleted"`
}

type MapSummary struct {
	ID   string `json:"id" bson:"_id"`
	Name string `json:"name" bson:"name"`
}

func (m *Map) Summary() MapSummary {
	return MapSummary{
		ID:   m.ID,
		Name: m.Name,
	}
}

func (m *Map) Valid() error {
	if strings.TrimSpace(m.Name) == "" {
		return errors.New("Name is required.")
	}

	for _, p := range m.Providences {
		if err := p.Valid(); err != nil {
			return err
		}
	}

	for _, p := range m.Players {
		if err := p.Valid(); err != nil {
			return err
		}
	}

	return nil
}
