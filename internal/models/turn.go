package models

type Turn struct {
	ID         string  `json:"id" bson:"id"`
	PhaseID    string  `json:"phaseID" bson:"phaseID"`
	Orders     []Order `json:"orders" bson:"orders"`
	TurnNumber int     `json:"turnNumber" bson:"turnNumber"`
	StartDate  int     `json:"startDate" bson:"startDate"`
	EndDate    int     `json:"endDate" bson:"endDate"`
}

func (t *Turn) Valid() error {
	for _, o := range t.Orders {
		if err := o.Valid(); err != nil {
			return err
		}
	}

	return nil
}
