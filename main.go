package main

import (
	"diplomaSE/gamestate"
	"fmt"
	"strings"
)

func main() {
	gamestate.BuildMap()
	for k, v := range gamestate.RegionIndex {
		fmt.Printf("Region: %s\n", k)
		fmt.Printf("\tIs an SC: %t\n", v.IsSC)
		fmt.Printf("\tLinked with: " + strings.Join(v.LinkedWith, ", ") + "\n")
	}
}
