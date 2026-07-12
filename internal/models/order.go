package models

type Order struct {
	ID          string `json:"id" bson:"id"`
	PhaseID     string `json:"phaseID" bson:"phaseID"`
	PlayerName  string `json:"playerName" bson:"playerName"`
	CreatedDate int    `json:"createdDate" bson:"createdDate"`
	Value       string `json:"value" bson:"value"`
}
