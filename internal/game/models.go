package game

import "regexp"

var CommandsRegex = struct {
	Hold      *regexp.Regexp
	Move      *regexp.Regexp
	Retreat   *regexp.Regexp
	Support   *regexp.Regexp
	Convoy    *regexp.Regexp
	Reinforce *regexp.Regexp
	Disband   *regexp.Regexp
}{
	Hold:      regexp.MustCompile(`^([af])\s*([a-z]{3}|[a-z]{3}-[nsew]c)\s*-\s*(?:hold|h)$`),
	Move:      regexp.MustCompile(`^([af])\s*([a-z]{3}|[a-z]{3}-[nsew]c)\s*-\s*([a-z]{3}|[a-z]{3}-[nsew]c)$`),
	Retreat:   regexp.MustCompile(`^([af])\s*r\s*([a-z]{3}|[a-z]{3}-[nsew]c)\s*-\s*([a-z]{3}|[a-z]{3}-[nsew]c)$`),
	Support:   regexp.MustCompile(`^([af])\s*([a-z]{3}|[a-z]{3}-[nsew]c)\s*s\s*([af])\s*([a-z]{3}|[a-z]{3}-[nsew]c)\s*-\s*(hold|h|[a-z]{3})$`),
	Convoy:    regexp.MustCompile(`^f\s*([a-z]{3}|[a-z]{3}-[nsew]c)\s*c\s*a\s*([a-z]{3}|[a-z]{3}-[nsew]c)\s*-\s*([a-z]{3}|[a-z]{3}-[nsew]c)$`),
	Reinforce: regexp.MustCompile(`^([af])\s*([a-z]{3}|[a-z]{3}-[nsew]c)$`),
	Disband:   regexp.MustCompile(`^d\s*([af])\s*([a-z]{3}|[a-z]{3}-[nsew]c)$`),
}

const (
	ProvidenceOcean   ProvidenceType = "ocean"
	ProvidenceCoastal ProvidenceType = "coastal"
	ProvidenceInland  ProvidenceType = "inland"
)

const (
	UnitArmy  UnitType = "army"
	UnitFleet UnitType = "fleet"
)

const (
	NorthCoast Coast = "nc"
	SouthCoast Coast = "sc"
	EastCoast  Coast = "ec"
	WestCoast  Coast = "wc"
)

type ProvidenceType string
type UnitType string
type Coast string

type Providence struct {
	Id            string             `json:"id" bson:"id"`
	Name          string             `json:"name" bson:"name"`
	SupplyCenter  *SupplyCenter      `json:"supplyCenter" bson:"supplyCenter"`
	Unit          *Unit              `json:"unit" bson:"unit"`
	Coordinates   Coordinates        `json:"coordinates" bson:"coordinates"`
	Type          ProvidenceType     `json:"type" bson:"type"`
	Routes        []string           `json:"routes" bson:"routes"`
	CoastalRoutes map[Coast][]string `json:"coastalRoutes" bson:"coastalRoutes"`
}

type SupplyCenter struct {
	ControlledBy *string     `json:"controlledBy" bson:"controlledBy"`
	Coordinates  Coordinates `json:"coordinates" bson:"coordinates"`
}

type Unit struct {
	Id           string            `json:"id" bson:"id"`
	ControlledBy string            `json:"controlledBy" bson:"controlledBy"`
	Type         UnitType          `json:"type" bson:"type"`
	Location     LocationReference `json:"location" bson:"location"`
}

type Coordinates struct {
	X float64 `json:"x" bson:"x"`
	Y float64 `json:"y" bson:"y"`
}

type LocationReference struct {
	Id    string `json:"id" bson:"id"`
	Coast *Coast `json:"coast" bson:"coast"`
}

type MoveCommand struct {
	PlayerName  string            `json:"playerName" bson:"playerName"`
	UnitType    UnitType          `json:"unitType" bson:"unitType"`
	Location    LocationReference `json:"location" bson:"location"`
	Destination LocationReference `json:"destination" bson:"destination"`
}

type SupportCommand struct {
	PlayerName string            `json:"playerName" bson:"playerName"`
	UnitType   UnitType          `json:"unitType" bson:"unitType"`
	Location   LocationReference `json:"location" bson:"location"`
	Move       MoveCommand       `json:"move" bson:"move"`
}

type AdjustCommand struct {
	PlayerName string            `json:"playerName" bson:"playerName"`
	UnitType   UnitType          `json:"unitType" bson:"unitType"`
	Location   LocationReference `json:"location" bson:"location"`
}

type Commands struct {
	Hold      []MoveCommand    `json:"hold" bson:"hold"`
	Move      []MoveCommand    `json:"move" bson:"move"`
	Retreat   []MoveCommand    `json:"retreat" bson:"retreat"`
	Support   []SupportCommand `json:"support" bson:"support"`
	Convoy    []SupportCommand `json:"convoy" bson:"convoy"`
	Reinforce []AdjustCommand  `json:"reinforce" bson:"reinforce"`
	Disband   []AdjustCommand  `json:"disband" bson:"disband"`
}

type Player struct {
	UserId    *string `json:"userId" bson:"userId"`
	Name      string  `json:"name" bson:"name"`
	Color     string  `json:"color" bson:"color"`
	IsPlaying bool    `json:"isPlaying" bson:"isPlaying"`
}

type Phase struct {
	Id          string `json:"id" bson:"_id"`
	Name        string `json:"name" bson:"name"`
	Description string `json:"description" bson:"description"`
	PhaseOrder  int8   `json:"phaseOrder" bson:"phaseOrder"`
}

type Turn struct {
	Id         string  `json:"id" bson:"id"`
	PhaseId    string  `json:"phaseId" bson:"phaseId"`
	Orders     []Order `json:"orders" bson:"orders"`
	TurnNumber int     `json:"turnNumber" bson:"turnNumber"`
	StartDate  int     `json:"startDate" bson:"startDate"`
	EndDate    int     `json:"endDate" bson:"endDate"`
}

type Order struct {
	Id          string `json:"id" bson:"id"`
	PhaseId     string `json:"phaseId" bson:"phaseId"`
	PlayerName  string `json:"playerName" bson:"playerName"`
	CreatedDate int    `json:"createdDate" bson:"createdDate"`
	Value       string `json:"value" bson:"value"`
}

type User struct {
	Id          string `json:"id" bson:"_id"`
	FirstName   string `json:"firstName" bson:"firstName"`
	LastName    string `json:"lastName" bson:"lastName"`
	Username    string `json:"username" bson:"username"`
	Password    string `json:"password" bson:"password"`
	Salt        []byte `json:"salt" bson:"salt"`
	CreatedDate int64  `json:"createdDate" bson:"createdDate"`
	IsDeleted   bool   `json:"isDeleted" bson:"isDeleted"`
}

type Game struct {
	Id      string `json:"id" bson:"_id"`
	OwnerId string `json:"ownerId" bson:"ownerId"`
	Map     struct {
		Id       string `json:"id" bson:"id"`
		Filename string `json:"filename" bson:"filename"`
	} `json:"map" bson:"map"`
	Board         []Providence `json:"board" bson:"board"`
	Players       []Player     `json:"players" bson:"players"`
	Turns         []Turn       `json:"turns" bson:"turns"`
	DaysPerTurn   int          `json:"daysPerTurn" bson:"daysPerTurn"`
	TurnStartHour int          `json:"turnStartHour" bson:"turnStartHour"`
	Timezone      int          `json:"timezone" bson:"timezone"`
	StartDate     int          `json:"startDate" bson:"startDate"`
	EndDate       int          `json:"endDate" bson:"endDate"`
	InProgress    bool         `json:"inProgress" bson:"inProgress"`
	IsDeleted     bool         `json:"isDeleted" bson:"isDeleted"`
}

type Map struct {
	Id          string       `json:"id" bson:"_id"`
	Filename    string       `json:"filename" bson:"filename"`
	Name        string       `json:"name" bson:"name"`
	Players     []Player     `json:"players" bson:"players"`
	Providences []Providence `json:"providences" bson:"providences"`
}
