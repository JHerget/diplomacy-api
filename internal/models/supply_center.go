package models

type SupplyCenter struct {
	ControlledBy *string     `json:"controlledBy" bson:"controlledBy"`
	Coordinates  Coordinates `json:"coordinates" bson:"coordinates"`
}
