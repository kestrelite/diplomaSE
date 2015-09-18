package gamestate

// OrderSet is a map[RegionCode]*Order that holds a map of orders by code
type OrderSet map[RegionCode]*Order

// Order acts as an interface to get the details about a generic order
type Order struct {
	orderRegion  *Region
	orderSup     *RegionCode
	orderTo      *RegionCode
	orderType    OrderType
	depOrders    []Order
	orderInvalid bool
}

func newHoldOrder(regID RegionCode) *Order {
	return &Order{RegionIndex[regID], nil, nil, OrderTypeHold, *new([]Order), false}
}

// NewOrderSet builds a new order set by iterating over the keys of the RegionIndex
func NewOrderSet() *OrderSet {
	oSet := make(OrderSet)
	for k := range RegionIndex {
		(oSet)[k] = newHoldOrder(k)
	}
	return &oSet
}

// AddOrder adds an order to the set by overwriting the existing default hold order
func (o *OrderSet) AddOrder(regID, orderA, orderB *RegionCode, oType OrderType) {
	newOrd := new(Order)
	newOrd.orderRegion = RegionIndex[*regID]
	if orderA != nil {
		nOrdA := RegionCode(*orderA)
		newOrd.orderSup = &nOrdA
	}
	if orderB != nil {
		nOrdB := RegionCode(*orderB)
		newOrd.orderTo = &nOrdB
	}
	newOrd.orderType = oType
	(*o)[*regID] = newOrd
}

// CleanOrders removes all illegal orders

//If a convoy is not part of a legal convoy chain, mark all items in the chain invalid.
func (o *OrderSet) CleanOrders() {
	for k, v := range *o {
		//If an order was given to a unit that doesn't exist, delete it.
		if v.orderRegion.OccupiedBy == UnitTypeNone {
			delete(*o, k)
			continue
		}
		//If a move order has an illegal destination, mark it invalid.
		//If a support order has an illegal destination, mark it invalid.
		//If a convoy order has an illegal destination, mark it invalid.
		if (v.orderType == OrderTypeMove || v.orderType == OrderTypeSupport || v.orderType == OrderTypeConvoy) &&
			v.orderTo != nil && !v.orderRegion.IsAdjacent(*v.orderTo, v.orderRegion.OccupiedBy) {
			v.orderInvalid = true
		}

		if v.orderType == OrderTypeSupport {
			oSupported := (*o)[*v.orderSup]
			//If a support order supports a move order to hold, or vice versa, mark it invalid.
			if v.orderTo == nil && oSupported.orderType == OrderTypeMove {
				v.orderInvalid = true
				continue //Needed because of proceeding null dereference
			}
			if v.orderTo != nil && oSupported.orderType == OrderTypeConvoy ||
				oSupported.orderType == OrderTypeHold ||
				oSupported.orderType == OrderTypeSupport {
				v.orderInvalid = true
			}

			//If a support order supports a move order that has a source that doesn't match the destination, mark it invalid.
			if oSupported.orderType == OrderTypeMove && *oSupported.orderTo != *v.orderTo {
				v.orderInvalid = true
			}

			//If a support order supports a move order with a source with an illegal destination, mark it invalid.
			if oSupported.orderType == OrderTypeMove && !oSupported.orderRegion.IsAdjacent(*oSupported.orderTo, oSupported.orderRegion.OccupiedBy) {
				v.orderInvalid = true
			}
		}
	}

	//If no order was given but a unit exists, add a hold order.
	for k, v := range RegionIndex {
		if _, ok := (*o)[k]; !ok && v.OccupiedBy != UnitTypeNone {
			(*o)[k] = newHoldOrder(k)
		}
	}
}
