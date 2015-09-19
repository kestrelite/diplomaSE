package gamestate

import "testing"

func TestAdjudicate(t *testing.T) {
	BuildMap()
	tSet := getTestOrderSet()
	tSet.Adjudicate()
}
