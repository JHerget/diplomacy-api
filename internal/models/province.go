package models

type Province struct {
	ID            string             `json:"id" bson:"id"`
	Name          string             `json:"name" bson:"name"`
	SupplyCenter  *SupplyCenter      `json:"supplyCenter" bson:"supplyCenter"`
	Unit          *Unit              `json:"unit" bson:"unit"`
	Coordinates   Coordinates        `json:"coordinates" bson:"coordinates"`
	Type          ProvinceType       `json:"type" bson:"type"`
	Routes        []string           `json:"routes" bson:"routes"`
	CoastalRoutes map[Coast][]string `json:"coastalRoutes" bson:"coastalRoutes"`
}

func (p *Province) Valid() error {
	return nil
}
