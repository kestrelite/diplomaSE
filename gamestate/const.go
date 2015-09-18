package gamestate

type (
	OrderState int
	OrderType  int
	UnitType   int
	NationCode int
	RegionCode string
)

const (
	OrderStateUnresolved OrderState = iota
	OrderStateFailed
	OrderStateSucceeds
	OrderStatePending
)

const (
	OrderTypeHold OrderType = iota
	OrderTypeSupport
	OrderTypeMove
	OrderTypeConvoy
)

const (
	UnitTypenone UnitType = iota
	UnitTypeArmy
	UnitTypeFleet
)

const (
	NationCodeNone NationCode = iota
	NationCodeEngland
	NationCodeFrance
	NationCodeGermany
	NationCodeRussia
	NationCodeItaly
	NationCodeAustria
	NationCodeTurkey
)
