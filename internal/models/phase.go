package models

type Phase struct {
	ID          string `json:"id" bson:"_id,omitempty"`
	Name        string `json:"name" bson:"name"`
	Description string `json:"description" bson:"description"`
	PhaseOrder  int8   `json:"phaseOrder" bson:"phaseOrder"`
}
