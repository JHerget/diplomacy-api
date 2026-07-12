package models

type Unit struct {
	ID           string            `json:"id" bson:"id"`
	ControlledBy string            `json:"controlledBy" bson:"controlledBy"`
	Type         UnitType          `json:"type" bson:"type"`
	Location     LocationReference `json:"location" bson:"location"`
}
