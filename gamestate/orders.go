package gamestate

// OrderSet is a map[RegionCode]*Order that holds a map of orders by code
type OrderSet map[RegionCode]*Order

// Order acts as an interface to get the details about a generic order
type Order struct {
	orderRegion  *Region
	orderSup     *RegionCode
	orderTo      *RegionCode
	orderType    OrderType
	depOrders    []*Order
	orderInvalid bool
	orderComment string
}

func newHoldOrder(regID RegionCode) *Order {
	return &Order{RegionIndex[regID], nil, nil, OrderTypeHold, *new([]*Order), false, ""}
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

// ValidateOrders removes all illegal orders
func (o *OrderSet) validateOrders() {
	var convoyList, moveConvoyList []*Order

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
			v.orderComment = orderRejectBadDestination
			v.orderInvalid = true
		}

		if v.orderType == OrderTypeSupport {
			oSupported := (*o)[*v.orderSup]
			//If a support order supports a move order to hold, or vice versa, mark it invalid.
			if v.orderTo == nil && oSupported.orderType == OrderTypeMove {
				v.orderComment = orderRejectWrongSupportType
				v.orderInvalid = true
				continue //Needed because of proceeding null dereference
			}
			if v.orderTo != nil && oSupported.orderType == OrderTypeConvoy ||
				oSupported.orderType == OrderTypeHold ||
				oSupported.orderType == OrderTypeSupport {
				v.orderComment = orderRejectWrongSupportType
				v.orderInvalid = true
			}

			//If a support order's destination doesn't match the supported destination, mark it invalid.
			if oSupported.orderType == OrderTypeMove && *oSupported.orderTo != *v.orderTo {
				v.orderComment = orderRejectSupportNotConsistent
				v.orderInvalid = true
			}

			//If a support order supports a move order with a source with an illegal destination, mark it invalid.
			if oSupported.orderType == OrderTypeMove && !oSupported.orderRegion.IsAdjacent(*oSupported.orderTo, oSupported.orderRegion.OccupiedBy) {
				v.orderComment = orderRejectSupportBadDest
				v.orderInvalid = true
			}
		}

		//Index the convoy and convoy move orders for later
		if v.orderType == OrderTypeConvoy {
			convoyList = append(convoyList, v)
		}

		if v.orderType == OrderTypeMoveConvoy {
			moveConvoyList = append(moveConvoyList, v)
		}
	}

	//PROCESS CONVOY ORDERS

	//Start by assuming all convoy orders invalid
	for _, v := range append(convoyList, moveConvoyList...) {
		v.orderInvalid = true
	}

	//For each item in the moveConvoyList, check to see if it has a valid path
	//to its land territory. If it does, mark all those convoy orders valid.
	//If it doesn't, leave them marked invalid.
	for _, army := range moveConvoyList {
		//fmt.Println(army.orderRegion.RegionID, "moving to", *army.orderTo)
		//If the army isn't moving into a land-based territory, stappit.
		if len(RegionIndex[*army.orderTo].AdjacentLand) == 0 {
			//fmt.Println("\t not moving into a land territory; aborting")
			army.orderComment = orderRejectConvoyMovingOntoWater
			army.orderInvalid = true
			continue
		}
		//If the army isn't actually an army, stappit.
		if army.orderRegion.OccupiedBy != UnitTypeArmy {
			//fmt.Println("\t this isn't actually an army; aborting")
			army.orderComment = orderRejectConvoyNotAnArmy
			army.orderInvalid = true
			continue
		}

		//Get a list of possible starting transfers
		var possInitialConvoys []*Order
		for _, p := range convoyList {
			if army.orderRegion.RegionID == *(p.orderSup) &&
				p.orderRegion.IsAdjacent(*p.orderTo, UnitTypeFleet) {
				possInitialConvoys = append(possInitialConvoys, p)
			}
		}

		if len(possInitialConvoys) == 0 {
			army.orderComment = orderRejectConvoyNoStartingPath
			army.orderInvalid = true
			//fmt.Println("\t welp, looks like there's no solution to this one")
			continue
		}

		//For each item in the list of possible starting nodes:
		// - Check to see if its exit is the same as the convoy's exit
		//   - If it is, mark all orders in the chain as valid.
		// - If it's not, try to add another node to the chain
		//   - If no new node could be added, exit without a path.
		for _, startConvoy := range possInitialConvoys {
			//fmt.Println("\tPossible initial convoy:", startConvoy.orderRegion.RegionID, "convoys", *startConvoy.orderSup, "to", *startConvoy.orderTo)

			var path []*Order
			path = append(path, army, startConvoy)
			army.orderInvalid = false
			startConvoy.orderInvalid = false
			pathFound := false

			for !pathFound {
				currConvoy := path[len(path)-1]
				//fmt.Println("\t\tFocusing on convoy:", currConvoy.orderRegion.RegionID, "convoys", *currConvoy.orderSup, "to", *currConvoy.orderTo)
				if *currConvoy.orderTo == *army.orderTo {
					//fmt.Println("\t\t\tConvoy satisfies the convoy chain.")
					pathFound = true
					break
				}
				for _, nextConvoy := range convoyList {
					//fmt.Println("\t\t\tTesting possible next convoy:", nextConvoy.orderRegion.RegionID, "convoys", *nextConvoy.orderSup, "to", *nextConvoy.orderTo)
					if nextConvoy.orderInvalid == false {
						//fmt.Println("\t\t\t\t convoy not legal for use at this point")
						continue //This is either already part of a valid chain, or part of our chain
					}
					if len(nextConvoy.orderRegion.AdjacentLand) != 0 {
						//fmt.Println("\t\t\t\t ship is on water; order invalid")
						nextConvoy.orderComment = orderRejectConvoyUsingDockedShip
						continue //You can't convoy with a docked ship
					}

					if *currConvoy.orderTo == nextConvoy.orderRegion.RegionID &&
						*nextConvoy.orderSup == currConvoy.orderRegion.RegionID &&
						nextConvoy.orderRegion.IsAdjacent(*nextConvoy.orderSup, UnitTypeFleet) &&
						nextConvoy.orderRegion.IsAdjacent(*nextConvoy.orderTo, UnitTypeFleet) {
						//fmt.Println("\t\t\t\t hey, this one might work!")
						path = append(path, nextConvoy)
						nextConvoy.orderInvalid = false
						break
					}

					//fmt.Println("\t\t\t\t order does not connect to", *currConvoy.orderTo)
				}
				passConvoy := path[len(path)-1]
				if currConvoy == passConvoy {
					army.orderComment = orderRejectConvoyNoCompletePath
					//fmt.Println("\t It doesn't look like a legal path for this convoy exists.")
					break
				}
			}

			//Revert everything we've modified back to orderInvalid = true assumption
			if !pathFound {
				for _, v := range path {
					v.orderInvalid = true
				}
			}
		}
	}

	for _, v := range convoyList {
		if v.orderInvalid && len(v.orderComment) == 0 {
			v.orderComment = orderRejectConvoyNotPartOfChain
		}
	}

	//If no order was given but a unit exists, add a hold order.
	for k, v := range RegionIndex {
		if _, ok := (*o)[k]; !ok && v.OccupiedBy != UnitTypeNone {
			(*o)[k] = newHoldOrder(k)
		}
	}
}

// CleanInvalidOrders strips an OrderSet of its invalid orders, replacing them
// with hold orders.
func (o *OrderSet) cleanupOrders() (mk []*Order) {
	for k, v := range *o {
		if v.orderInvalid {
			mk = append(mk, v)
			(*o)[k] = newHoldOrder(v.orderRegion.RegionID)
		}
	}

	for k := range RegionIndex {
		if _, ok := (*o)[k]; !ok {
			(*o)[k] = newHoldOrder(k)
		}
	}
	return
}
