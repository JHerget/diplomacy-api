package models

type Map struct {
	ID          string       `json:"id" bson:"_id,omitempty"`
	Name        string       `json:"name" bson:"name"`
	Players     []Player     `json:"players" bson:"players"`
	Providences []Providence `json:"providences" bson:"providences"`
	IsDeleted   bool         `json:"isDeleted" bson:"isDeleted"`
}

type MapSummary struct {
	ID   string `json:"id" bson:"id"`
	Name string `json:"name" bson:"name"`
}

func (m *Map) Summary() MapSummary {
	return MapSummary{
		ID:   m.ID,
		Name: m.Name,
	}
}
