package gamestate

type OrderSet map[RegionCode]*Order

// Order acts as an interface to get the details about a generic order
type Order struct {
	orderRegion *Region
	orderA      *Order
	orderB      *Order
	ordType     OrderType
	depOrders   []Order
}

func newHoldOrder(regID RegionCode) *Order {
	return &Order{RegionIndex[regID], nil, nil, OrderTypeHold, *new([]Order)}
}

func NewOrderSet() *OrderSet {
	oSet := make(OrderSet)
	for k := range RegionIndex {
		(oSet)[k] = newHoldOrder(k)
	}
	return &oSet
}

func (o *OrderSet) AddOrder(regID RegionCode, oType OrderType, a, b Order) {
	newOrd := new(Order)
	newOrd.orderRegion = RegionIndex[regID]
	newOrd.orderA = &a
	newOrd.orderB = &b
	newOrd.ordType = oType
	(*o)[regID] = newOrd
}
