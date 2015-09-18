package gamestate

import (
	"diplomaSE/data"
	"log"
	"strings"
)

// RegionCodeList is the full list of possible region codes and their respective regions.
var RegionIndex = make(map[string]*Region)

// Region holds information on connectedness properties, and contained units.
type Region struct {
	RegionCode    string
	IsSC          bool
	LinkedWith    []string
	AdjacentWater []string
	AdjacentLand  []string
}

const mapLocation = "data/testmap.txt"

var isMapInitialized = false

// BuildMap Builds the map from file data
func BuildMap() (out map[string]*Region) {
	out = make(map[string]*Region)
	byteData, err := data.Asset(mapLocation)
	if err != nil {
		log.Fatal(err)
	}
	linkedLocations := map[string][]string{
		"STP": []string{"STP", "STPsc", "STPnc"},
		"BUL": []string{"BUL", "BULsc", "BULnc"},
		"SPN": []string{"SPN", "SPNsc", "SPNnc"},
	}
	s := string(byteData[:])
	lines := strings.Split(s, "\n")
	lines = lines[:len(lines)-1]

	for _, line := range lines {
		token := strings.Split(line, "|")
		isWater := (token[0] == "W")
		regCode := strings.TrimSpace(token[2])

		reg := new(Region)
		if ereg, ok := RegionIndex[regCode]; ok {
			reg = ereg
		} else {
			reg.IsSC = (token[1] == "Y")
			reg.RegionCode = regCode
			if link, ok := linkedLocations[regCode[:3]]; ok {
				reg.LinkedWith = link
			}
		}

		for _, adj := range strings.Split(token[3], ",") {
			adj = strings.TrimSpace(adj)
			if isWater {
				reg.AdjacentWater = append(reg.AdjacentWater, adj)
			} else {
				reg.AdjacentLand = append(reg.AdjacentLand, adj)
			}
		}

		RegionIndex[reg.RegionCode] = reg
		out[reg.RegionCode] = reg
	}
	return
}

// IsMapBuilt returns whether the map has been constructed yet
func IsMapBuilt() bool {
	return isMapInitialized
}
