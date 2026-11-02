package models

const (
	ProvinceOcean   ProvinceType = "ocean"
	ProvinceCoastal ProvinceType = "coastal"
	ProvinceInland  ProvinceType = "inland"
)

const (
	UnitArmy  UnitType = "army"
	UnitFleet UnitType = "fleet"
)

const (
	NorthCoast   Coast = "nc"
	SouthCoast   Coast = "sc"
	EastCoast    Coast = "ec"
	WestCoast    Coast = "wc"
	DefaultCoast Coast = "default"
)

type ProvinceType string
type UnitType string
type Coast string
