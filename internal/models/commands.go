package models

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
