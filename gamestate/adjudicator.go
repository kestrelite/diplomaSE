package gamestate

func (o *OrderSet) setOrderDeps() {

}

func (o *OrderSet) stripOrderCycles() {

}

func resolveOrderCycle(o []*Order) {

}

func (o *OrderSet) assignStrengths() {

}

func (o *OrderSet) resolveRecursive() {

}

func (o *OrderSet) getResolvedMap(r map[RegionCode]*Region) map[RegionCode]*Region {
	return nil
}

// Adjudicate takes the set of orders and returns the new game state
func (o *OrderSet) Adjudicate() *map[RegionCode]*Region {
	o.validateOrders()
	o.cleanupOrders()
	o.setOrderDeps()
	o.stripOrderCycles()
	o.assignStrengths()
	o.resolveRecursive()
	RegionIndex = o.getResolvedMap(RegionIndex)
	return &RegionIndex
}
