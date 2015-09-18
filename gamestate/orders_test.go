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
	//	testSet := NewOrderSet()
	//    testSet.AddOrder("BEL", OrderTypeMove, a, b Order)
}
