package models

type Province struct {
	ID           string                  `json:"id" bson:"id"`
	Name         string                  `json:"name" bson:"name"`
	SupplyCenter *SupplyCenter           `json:"supplyCenter" bson:"supplyCenter"`
	Unit         *Unit                   `json:"unit" bson:"unit"`
	Coordinates  Coordinates             `json:"coordinates" bson:"coordinates"`
	Type         ProvinceType            `json:"type" bson:"type"`
	Routes       []string                `json:"routes" bson:"routes"`
	Coasts       map[Coast]ProvinceCoast `json:"coasts" bson:"coasts"`
}

type ProvinceCoast struct {
	Routes      []string    `json:"routes" bson:"routes"`
	Coordinates Coordinates `json:"coordinates" bson:"coordinates"`
}

func (p *Province) Valid() error {
	return nil
}

func (p *Province) IsCoastal() bool {
	return p.Type == ProvinceCoastal
}

func (p *Province) IsInland() bool {
	return p.Type == ProvinceInland
}

func (p *Province) IsOcean() bool {
	return p.Type == ProvinceOcean
}

func (p *Province) GetCoordinates() Coordinates {
	if p.Unit == nil {
		return p.Coordinates
	}

	if p.IsCoastal() {
		unitCoast := DefaultCoast
		if p.Unit.Location.Coast != nil {
			unitCoast = *p.Unit.Location.Coast
		}

		coast, ok := p.Coasts[unitCoast]
		if ok {
			return coast.Coordinates
		}
	}

	return p.Coordinates
}
