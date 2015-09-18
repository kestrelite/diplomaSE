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
	OrderTypeMoveConvoy
	OrderTypeConvoy
)

const (
	UnitTypeNone UnitType = iota
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

const (
	orderRejectBadDestination       string = "The order has an illegal destination."
	orderRejectWrongSupportType     string = "The order supports a move to hold or vice versa."
	orderRejectSupportNotConsistent string = "The order supports a move to a different location."
	orderRejectSupportBadDest       string = "The supported order has an illegal destination."

	orderRejectConvoyMovingOntoWater string = "The convoy terminates on water."
	orderRejectConvoyNotAnArmy       string = "The unit convoyed is not an army."
	orderRejectConvoyNoStartingPath  string = "The convoyed unit has no possible starting paths."
	orderRejectConvoyNoCompletePath  string = "The convoyed unit has no legal path to its destination."
	orderRejectConvoyUsingDockedShip string = "The convoying ship is on land and is illegal."
	orderRejectConvoyNotPartOfChain  string = "The convoy was not part of a legal convoy chain."
)
