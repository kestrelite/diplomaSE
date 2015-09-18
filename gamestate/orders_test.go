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

var testOrderSetBuilt = false
var testOrderSet *OrderSet

func TestCleanInvalidOrders(t *testing.T) {
	testSet := getTestOrderSet()
	testSet.ValidateOrders()
	out := testSet.CleanInvalidOrders()
	for _, v := range out {
		if len(v.orderComment) == 0 {
			t.Error("Order comment missing for order", v.orderRegion.RegionID)
		}
	}
}

func TestValidateOrders(t *testing.T) {
	testSet := getTestOrderSet()
	testSet.ValidateOrders()

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

	//If a convoy order moves a non-army unit, or a nonexistent unit, mark it invalid.
	if !(*testSet)[RegionCode("AEG")].orderInvalid {
		t.Error("The convoy of a nonexistent unit AEG C BUL -> GRE was not marked invalid")
	}
	if !(*testSet)[RegionCode("TYN")].orderInvalid {
		t.Error("The convoy for a fleet in ROM (ROM F) (TYN C ROM -> TUN) was not marked invalid")
	}
	if !(*testSet)[RegionCode("ROM")].orderInvalid {
		t.Error("The convoy of a fleet in ROM was not marked invalid")
	}

	//If a convoy order moves a unit that doesn't also have a convoy order, mark it invalid.
	if !(*testSet)[RegionCode("BOT")].orderInvalid {
		t.Error("The convoy of a holding army BOT C LVN -> SWE was not marked invalid")
	}

	//If a convoy move order is moving a unit to somewhere it can't be, mark it invalid.
	if !(*testSet)[RegionCode("NRG")].orderInvalid {
		t.Error("The convoy of an army into NAT (NRG C EDI -> NAT) was not marked invalid")
	}
	if !(*testSet)[RegionCode("EDI")].orderInvalid {
		t.Error("The convoy move of EDI -> NAT by NRG was not marked invalid")
	}

	//Stray convoy order invalidation
	if !(*testSet)[RegionCode("HEL")].orderInvalid {
		t.Error("The stray conovy HEL C DEN -> KIE was not marked invalid")
	}

	//Test a legal convoy chain
	if (*testSet)[RegionCode("TRI")].orderInvalid {
		t.Error("A legal convoy move order from TRI -> APU was marked invalid")
	}
	if (*testSet)[RegionCode("ADR")].orderInvalid {
		t.Error("A legal convoy order ADR C TRI -> ION was marked invalid")
	}
	if (*testSet)[RegionCode("ION")].orderInvalid {
		t.Error("A legal convoy order ION C ADR -> APU was marked invalid")
	}

	//Test an illegal convoy chain
	if !(*testSet)[RegionCode("STP")].orderInvalid {
		t.Error("A convoy out of STP that has no route was not marked invalid")
	}
	if !(*testSet)[RegionCode("BAR")].orderInvalid {
		t.Error("A convoy using BAR that has no route was not marked invalid")
	}
	if !(*testSet)[RegionCode("BAL")].orderInvalid {
		t.Error("A convoy using BAL that has no route was not marked invalid")
	}

	//Legal convoy - control test
	if (*testSet)[RegionCode("ANK")].orderInvalid {
		t.Error("The convoyed move ANK -> SEV was marked invalid")
	}
	if (*testSet)[RegionCode("BLA")].orderInvalid {
		t.Error("The convoy order BLA C ANK -> SEV was marked invalid")
	}
}

func getTestOrderSet() (testSet *OrderSet) {
	if testOrderSetBuilt {
		testSet = testOrderSet
		return
	}

	BuildMap()
	testSet = NewOrderSet()
	var orderFrom, orderSup, orderTo RegionCode

	//MOVE EXISTENCE VALIDATION

	//If no order was given but a unit exists, add a hold order.
	RegionIndex[RegionCode("PIC")].OccupiedBy = UnitTypeArmy
	RegionIndex[RegionCode("PIC")].OccupiedNation = NationCodeAustria

	//If an order was given to a unit that doesn't exist, delete it.
	orderFrom, orderTo = RegionCode("BEL"), RegionCode("HEL")
	testSet.AddOrder(&orderFrom, nil, &orderTo, OrderTypeMove)

	//MOVE ORDER VALIDATION

	//If a move order has an illegal destination, mark it invalid.
	RegionIndex[RegionCode("NWY")].OccupiedBy = UnitTypeArmy
	orderFrom, orderTo = RegionCode("NWY"), RegionCode("NTH")
	testSet.AddOrder(&orderFrom, nil, &orderTo, OrderTypeMove)

	RegionIndex[RegionCode("HOL")].OccupiedBy = UnitTypeFleet
	orderFrom, orderTo = RegionCode("HOL"), RegionCode("RUH")
	testSet.AddOrder(&orderFrom, nil, &orderTo, OrderTypeMove)

	//SUPPORT ORDER VALIDATION

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

	//CONVOY ORDER VALIDATION

	//If a convoy order moves a non-army unit, or a nonexistent unit, mark it invalid.
	RegionIndex[RegionCode("AEG")].OccupiedBy = UnitTypeFleet
	orderFrom, orderSup, orderTo = RegionCode("AEG"), RegionCode("BUL"), RegionCode("GRE")
	testSet.AddOrder(&orderFrom, &orderSup, &orderTo, OrderTypeConvoy)

	RegionIndex[RegionCode("ROM")].OccupiedBy = UnitTypeFleet
	RegionIndex[RegionCode("TYN")].OccupiedBy = UnitTypeFleet
	orderFrom, orderSup, orderTo = RegionCode("TYN"), RegionCode("ROM"), RegionCode("TUN")
	testSet.AddOrder(&orderFrom, &orderSup, &orderTo, OrderTypeConvoy)
	testSet.AddOrder(&orderSup, nil, &orderTo, OrderTypeMoveConvoy)

	//If a convoy order moves a unit that doesn't also have a convoy order, mark it invalid.
	RegionIndex[RegionCode("BOT")].OccupiedBy = UnitTypeFleet
	RegionIndex[RegionCode("LVN")].OccupiedBy = UnitTypeArmy
	orderFrom, orderSup, orderTo = RegionCode("BOT"), RegionCode("LVN"), RegionCode("SWE")
	testSet.AddOrder(&orderFrom, &orderSup, &orderTo, OrderTypeConvoy)
	testSet.AddOrder(&orderSup, nil, nil, OrderTypeHold)

	//If a convoy move order is moving to somewhere it can't be, mark it invalid.
	RegionIndex[RegionCode("NRG")].OccupiedBy = UnitTypeFleet
	RegionIndex[RegionCode("EDI")].OccupiedBy = UnitTypeArmy
	orderFrom, orderSup, orderTo = RegionCode("NRG"), RegionCode("EDI"), RegionCode("NAT")
	testSet.AddOrder(&orderFrom, &orderSup, &orderTo, OrderTypeConvoy)
	testSet.AddOrder(&orderSup, nil, &orderTo, OrderTypeMoveConvoy)

	//Stray convoy order invalidation
	RegionIndex[RegionCode("HEL")].OccupiedBy = UnitTypeFleet
	orderFrom, orderSup, orderTo = RegionCode("HEL"), RegionCode("DEN"), RegionCode("KIE")
	testSet.AddOrder(&orderFrom, &orderSup, &orderTo, OrderTypeConvoy)

	//Test a legal convoy chain
	RegionIndex[RegionCode("ADR")].OccupiedBy = UnitTypeFleet
	RegionIndex[RegionCode("ION")].OccupiedBy = UnitTypeFleet
	RegionIndex[RegionCode("TRI")].OccupiedBy = UnitTypeArmy
	orderFrom, orderTo = RegionCode("TRI"), RegionCode("APU")
	testSet.AddOrder(&orderFrom, nil, &orderTo, OrderTypeMoveConvoy)
	orderFrom, orderSup, orderTo = RegionCode("ADR"), RegionCode("TRI"), RegionCode("ION")
	testSet.AddOrder(&orderFrom, &orderSup, &orderTo, OrderTypeConvoy)
	orderFrom, orderSup, orderTo = RegionCode("ION"), RegionCode("ADR"), RegionCode("APU")
	testSet.AddOrder(&orderFrom, &orderSup, &orderTo, OrderTypeConvoy)

	//Test an illegal convoy chain
	RegionIndex[RegionCode("STP")].OccupiedBy = UnitTypeArmy
	RegionIndex[RegionCode("BAR")].OccupiedBy = UnitTypeFleet
	RegionIndex[RegionCode("BAL")].OccupiedBy = UnitTypeFleet
	orderFrom, orderTo = RegionCode("STP"), RegionCode("PRU")
	testSet.AddOrder(&orderFrom, nil, &orderTo, OrderTypeMoveConvoy)
	orderFrom, orderSup, orderTo = RegionCode("BAR"), RegionCode("STP"), RegionCode("BAL")
	testSet.AddOrder(&orderFrom, &orderSup, &orderTo, OrderTypeConvoy)
	orderFrom, orderSup, orderTo = RegionCode("BAL"), RegionCode("BAR"), RegionCode("PRU")
	testSet.AddOrder(&orderFrom, &orderSup, &orderTo, OrderTypeConvoy)

	//Legal convoy - control test
	RegionIndex[RegionCode("ANK")].OccupiedBy = UnitTypeArmy
	RegionIndex[RegionCode("BLA")].OccupiedBy = UnitTypeFleet
	orderFrom, orderSup, orderTo = RegionCode("BLA"), RegionCode("ANK"), RegionCode("SEV")
	testSet.AddOrder(&orderFrom, &orderSup, &orderTo, OrderTypeConvoy)
	testSet.AddOrder(&orderSup, nil, &orderTo, OrderTypeMoveConvoy)
	return
}
