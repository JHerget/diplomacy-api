package models

import "fmt"

type Coordinates struct {
	X float64 `json:"x" bson:"x"`
	Y float64 `json:"y" bson:"y"`
}

type LocationReference struct {
	ID    string `json:"id" bson:"id"`
	Coast *Coast `json:"coast" bson:"coast"`
}

func (l LocationReference) ProperID() string {
	if l.Coast != nil {
		return fmt.Sprintf("%s-%s", l.ID, *l.Coast)
	}

	return l.ID
}
