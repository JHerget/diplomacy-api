package board

type Providence struct {
    Id string `json:"id" bson:"id"`
    Name string `json:"name" bson:"name"`
    SupplyCenter *SupplyCenter `json:"supplyCenter" bson:"supplyCenter"`
    Unit *Unit `json:"unit" bson:"unit"`
    Coordinates Coordinates `json:"coordinates" bson:"coordinates"`
    Type ProvidenceType `json:"type" bson:"type"`
    Routes []string `json:"routes" bson:"routes"`
    CoastalRoutes map[string][]string `json:"coastalRoutes" bson:"coastalRoutes"`
}

type SupplyCenter struct {
    ControlledBy *string `json:"controlledBy" bson:"controlledBy"`
    Coordinates Coordinates `json:"coordinates" bson:"coordinates"`
}

type Unit struct {
    Id string `json:"id" bson:"id"`
    ControlledBy string `json:"controlledBy" bson:"controlledBy"`
    Type UnitType `json:"type" bson:"type"`
    Location LocationReference `json:"location" bson:"location"`
}

type Coordinates struct {
    X float64 `json:"x" bson:"x"`
    Y float64 `json:"y" bson:"y"`
}

type LocationReference struct {
    Id string `json:"id" bson:"id"`
    Coast *Coast `json:"coast" bson:"coast"`
}

type MoveCommand struct {
    PlayerName string `json:"playerName" bson:"playerName"`
    UnitType UnitType `json:"unitType" bson:"unitType"`
    Location LocationReference `json:"location" bson:"location"`
    Destination LocationReference `json:"destination" bson:"destination"`
}

type SupportCommand struct {
    PlayerName string `json:"playerName" bson:"playerName"`
    UnitType UnitType `json:"unitType" bson:"unitType"`
    Location LocationReference `json:"location" bson:"location"`
    Move MoveCommand `json:"move" bson:"move"`
}

type AdjustCommand struct {
    PlayerName string `json:"playerName" bson:"playerName"`
    UnitType UnitType `json:"unitType" bson:"unitType"`
    Location LocationReference `json:"location" bson:"location"`
}

type Commands struct {
    Hold: []MoveCommand `json:"hold" bson:"hold"`
    Move: []MoveCommand `json:"move" bson:"move"`
    Retreat: []MoveCommand `json:"retreat" bson:"retreat"`
    Support: []SupportCommand `json:"support" bson:"support"`
    Convoy: []SupportCommand `json:"convoy" bson:"convoy"`
    Reinforce: []AdjustCommand `json:"reinforce" bson:"reinforce"`
    Disband: []AdjustCommand `json:"disband" bson:"disband"`
}

type ProvidenceType string
const (
    ProvidenceOcean ProvidenceType = "ocean"
    ProvidenceCoastal ProvidenceType = "coastal"
    ProvidenceInland ProvidenceType = "inland"
)

type UnitType string
const (
    UnitArmy UnitType = "army"
    UnitFleet UnitType = "fleet"
)

type Coast string
const (
    NorthCoast Coast = "nc"
    SouthCoast Coast = "sc"
    EastCoast Coast = "ec"
    WestCoast Coast = "wc"
)
