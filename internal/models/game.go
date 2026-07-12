package models

type Game struct {
	ID            string       `json:"id" bson:"_id,omitempty"`
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
