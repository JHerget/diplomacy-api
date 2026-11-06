package models

var PlayerCombos = map[int][]string{
	2: {"Austria-Hungary/Germany/Turkey", "England/France/Russia"},
	3: {"England/Germany/Austria-Hungary", "Russia/Italy", "France/Turkey"},
	4: {"England", "Austria-Hungary/France", "Germany/Turkey", "Italy/Russia"},
	5: {"Russia", "England", "Turkey", "France", "Austria-Hungary"},
	6: {"Russia", "England", "Turkey", "France", "Austria-Hungary", "Germany"},
	7: {"Russia", "England", "Turkey", "France", "Austria-Hungary", "Germany", "Italy"},
}

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
