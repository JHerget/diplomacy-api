package board

type Providence struct {
    Id string `json:"id" bson:"id"`
    Name string `json:"name" bson:"name"`
    SupplyCenter *SupplyCenter `json:"supplyCenter" bson:"supplyCenter"`
    Unit *Unit `json:"unit" bson:"unit"`
    Coordinates Coordinates `json:"coordinates" bson:"coordinates"`
    Type ProvidenceType `json:"type" bson:"type"`
}

type SupplyCenter struct {}
type Unit struct {}
type Coordinates struct {}

type ProvidenceType string
const (
    ProvidenceOcean ProvidenceType = "ocean"
    ProvidenceCoastal ProvidenceType = "coastal"
    ProvidenceInland ProvidenceType = "inland"
)
