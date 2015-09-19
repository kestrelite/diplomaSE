package gamestate

import "fmt"

func (o *OrderSet) buildDFSTree() {
	for k := range *o {
		fmt.Print(k, ", ")
	}

	// for _, v := range *o {
	// 	if v.orderSup != nil {
	// 		oSup := (*o)[*v.orderSup]
	// 		fmt.Println(oSup)
	// 		oSup.depOrders = append(oSup.depOrders, v)
	// 	}
	// 	if v.orderTo != nil {
	// 		oTo := (*o)[*v.orderTo]
	// 		fmt.Println(*v.orderTo, (*o)[*v.orderTo], oTo)
	// 		oTo.depOrders = append(oTo.depOrders, v)
	// 	}
	// }
}

func (o *OrderSet) stripOrderCycles() {

}

func (o *OrderSet) resolveRecursive() {

}

func (o *OrderSet) getResolvedMap(r map[RegionCode]*Region) map[RegionCode]*Region {
	return nil
}

// Adjudicate takes the set of orders and returns the new game state
func (o *OrderSet) Adjudicate() *map[RegionCode]*Region {
	o.ValidateOrders()
	o.CleanInvalidOrders()
	o.buildDFSTree()
	o.stripOrderCycles()
	o.resolveRecursive()
	RegionIndex = o.getResolvedMap(RegionIndex)
	return &RegionIndex
}
