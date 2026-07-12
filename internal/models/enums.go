package models

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
