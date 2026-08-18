package models

type Player struct {
	ID        string  `json:"id" bson:"id"`
	UserID    *string `json:"userID" bson:"userID"`
	Name      string  `json:"name" bson:"name"`
	Color     string  `json:"color" bson:"color"`
	IsPlaying bool    `json:"isPlaying" bson:"isPlaying"`
}

func (p *Player) Valid() error {
	return nil
}
