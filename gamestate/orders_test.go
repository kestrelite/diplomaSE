package gamestate

import "testing"

func TestNewOrderSet(t *testing.T) {
	BuildMap()
	testSet := NewOrderSet()
	if len(*testSet) != 81 {
		t.Error("New order set was not of expected length!")
	}
}

func TestAddOrder(t *testing.T) {
	BuildMap()
	testSet := NewOrderSet()
	orderFrom, orderTo := RegionCode("PIC"), RegionCode("ENG")
	testSet.AddOrder(&orderFrom, nil, &orderTo, OrderTypeMove)
	ordWritten := (*testSet)[orderFrom]
	if ordWritten.orderType != OrderTypeMove || ordWritten.orderTo == &orderTo {
		t.Error(ordWritten.orderType, "=", OrderTypeMove, "; ", ordWritten.orderSup, "!=", &orderTo)
	}
}

func TestCleanOrders(t *testing.T) {
	BuildMap()
	testSet := NewOrderSet()
	var orderFrom, orderSup, orderTo RegionCode

	//If no order was given but a unit exists, add a hold order.
	RegionIndex[RegionCode("PIC")].OccupiedBy = UnitTypeArmy
	RegionIndex[RegionCode("PIC")].OccupiedNation = NationCodeAustria

	//If an order was given to a unit that doesn't exist, delete it.
	orderFrom, orderTo = RegionCode("BEL"), RegionCode("HEL")
	testSet.AddOrder(&orderFrom, nil, &orderTo, OrderTypeMove)

	//If a move order has an illegal destination, mark it invalid.
	RegionIndex[RegionCode("NWY")].OccupiedBy = UnitTypeArmy
	orderFrom, orderTo = RegionCode("NWY"), RegionCode("NTH")
	testSet.AddOrder(&orderFrom, nil, &orderTo, OrderTypeMove)

	RegionIndex[RegionCode("HOL")].OccupiedBy = UnitTypeFleet
	orderFrom, orderTo = RegionCode("HOL"), RegionCode("RUH")
	testSet.AddOrder(&orderFrom, nil, &orderTo, OrderTypeMove)

	//If a support order has an illegal destination, mark it invalid.
	RegionIndex[RegionCode("SYR")].OccupiedBy = UnitTypeArmy
	orderFrom, orderSup, orderTo = RegionCode("SYR"), RegionCode("SMY"), RegionCode("CON")
	testSet.AddOrder(&orderFrom, &orderSup, &orderTo, OrderTypeSupport)
	RegionIndex[RegionCode("SMY")].OccupiedBy = UnitTypeArmy
	orderFrom, orderTo = RegionCode("SMY"), RegionCode("CON")
	testSet.AddOrder(&orderFrom, nil, &orderTo, OrderTypeMove)

	//If a support order has a source that doesn't match the destination, mark it invalid.
	RegionIndex[RegionCode("MAO")].OccupiedBy = UnitTypeFleet
	RegionIndex[RegionCode("IRI")].OccupiedBy = UnitTypeFleet
	orderFrom, orderSup, orderTo = RegionCode("MAO"), RegionCode("IRI"), RegionCode("NAT")
	testSet.AddOrder(&orderFrom, &orderSup, &orderTo, OrderTypeSupport)
	orderFrom, orderTo = RegionCode("IRI"), RegionCode("ENG")
	testSet.AddOrder(&orderFrom, nil, &orderTo, OrderTypeMove)

	//If a support order supports a move order to hold, or vice versa, mark it invalid.
	RegionIndex[RegionCode("POR")].OccupiedBy = UnitTypeArmy
	RegionIndex[RegionCode("GOL")].OccupiedBy = UnitTypeFleet
	orderFrom, orderSup, orderTo = RegionCode("POR"), RegionCode("GOL"), RegionCode("SPNsc")
	testSet.AddOrder(&orderFrom, &orderSup, &orderTo, OrderTypeSupport)
	testSet.AddOrder(&orderSup, nil, nil, OrderTypeHold)

	RegionIndex[RegionCode("WAL")].OccupiedBy = UnitTypeArmy
	RegionIndex[RegionCode("YOR")].OccupiedBy = UnitTypeArmy
	orderFrom, orderSup, orderTo = RegionCode("WAL"), RegionCode("YOR"), RegionCode("LON")
	testSet.AddOrder(&orderFrom, &orderSup, nil, OrderTypeSupport)
	testSet.AddOrder(&orderSup, nil, &orderTo, OrderTypeMove)

	//If a support order supports a move order with a source with an illegal destination, mark it invalid.
	RegionIndex[RegionCode("GAS")].OccupiedBy = UnitTypeArmy
	RegionIndex[RegionCode("BRE")].OccupiedBy = UnitTypeFleet
	orderFrom, orderSup, orderTo = RegionCode("GAS"), RegionCode("BRE"), RegionCode("PAR")
	testSet.AddOrder(&orderFrom, &orderSup, &orderTo, OrderTypeSupport)
	testSet.AddOrder(&orderSup, nil, &orderTo, OrderTypeMove)

	//Legal support - control test
	RegionIndex[RegionCode("WES")].OccupiedBy = UnitTypeFleet
	orderFrom, orderSup, orderTo = RegionCode("WES"), RegionCode("TUN"), RegionCode("NAF")
	testSet.AddOrder(&orderFrom, &orderSup, &orderTo, OrderTypeSupport)
	RegionIndex[RegionCode("TUN")].OccupiedBy = UnitTypeArmy
	orderFrom, orderTo = RegionCode("TUN"), RegionCode("NAF")
	testSet.AddOrder(&orderFrom, nil, &orderTo, OrderTypeMove)

	//Legal convoy - control test
	// RegionIndex[RegionCode("ANK")].OccupiedBy = UnitTypeArmy
	// RegionIndex[RegionCode("BLA")].OccupiedBy = UnitTypeFleet
	// orderFrom, orderSup, orderTo = RegionCode("BLA"), RegionCode("ANK"), RegionCode("SEV")
	// testSet.AddOrder(&orderFrom, &orderSup, &orderTo, OrderTypeConvoy)
	// testSet.AddOrder(&orderSup, nil, &orderTo, OrderTypeMoveConvoy)

	//Cleanup orders for testing
	testSet.CleanOrders()

	//If no order was given but a unit exists, add a hold order.
	if _, ok := (*testSet)[RegionCode("PIC")]; !ok {
		t.Error("An auto-hold was not assigned to Picardy")
	}
	if (*testSet)[RegionCode("PIC")].orderRegion.RegionID != "PIC" {
		t.Error("Picardy's auto-hold was not assigned the correct region")
	}

	//If an order was given to a unit that doesn't exist, delete it.
	if data, ok := (*testSet)[RegionCode("BEL")]; ok {
		t.Error("Cleanup should have removed order on no unit: ", data.orderRegion)
	}

	//If a move order has an illegal destination, mark it invalid.
	if !(*testSet)[RegionCode("NWY")].orderInvalid {
		t.Error("A move from NWY -> NTH by Army was not marked invalid")
	}
	if !(*testSet)[RegionCode("HOL")].orderInvalid {
		t.Error("A move from HOL -> RUH by Fleet was not marked invalid")
	}

	//If a support order has an illegal destination, mark it invalid.
	if !(*testSet)[RegionCode("SYR")].orderInvalid {
		t.Error("The move SYR sup SMY -> CON was not marked invalid")
	}

	//If a support order supports a move order to hold, or vice versa, mark it invalid.
	if !(*testSet)[RegionCode("POR")].orderInvalid {
		t.Error("The move POR sup GOL -> SPNsc with GOL H was not marked invalid")
	}

	//If a support order has a source that doesn't match the destination, mark it invalid.
	if !(*testSet)[RegionCode("MAO")].orderInvalid {
		t.Error("The move MAO sup IRI -> NAT with IRI -> ENG was not marked invalid")
	}

	//If a support order supports a move order with a source with an illegal destination, mark it invalid.
	if !(*testSet)[RegionCode("GAS")].orderInvalid {
		t.Error("The move GAS sup BRE -> PAR with F BRE -> PAR was not marked invalid")
	}

	//Legal support - control test
	if (*testSet)[RegionCode("WES")].orderInvalid {
		t.Error("The move WES sup TUN -> NAF was marked invalid")
	}

	//Legal convoy - control test
	if (*testSet)[RegionCode("ANK")].orderInvalid {
		t.Error("The convoyed move ANK -> SEV was marked invalid")
	}
	if (*testSet)[RegionCode("BLA")].orderInvalid {
		t.Error("The convoy order BLA C ANK -> SEV was marked invalid")
	}

	//Convoy illegal destination invalid
}
