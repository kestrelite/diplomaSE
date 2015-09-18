package gamestate

import (
	"diplomaSE/data"
	"log"
	"strings"
)

// RegionIndex is the full list of possible region codes and their respective regions.
var RegionIndex = make(map[RegionCode]*Region)

// Region holds information on connectedness properties, and contained units.
type Region struct {
	RegionID       RegionCode
	IsSC           bool
	LinkedWith     []RegionCode
	AdjacentWater  []RegionCode
	AdjacentLand   []RegionCode
	OccupiedBy     UnitType
	OccupiedNation NationCode
	OwnedNation    NationCode
}

const mapLocation = "data/map.txt"

var isMapInitialized = false

func regionCodeContained(r []RegionCode, rt RegionCode) bool {
	for _, v := range r {
		if v == rt {
			return true
		}
	}
	return false
}

// IsAdjacent indicates whether a given region is adjacent to another region
func (r *Region) IsAdjacent(adj RegionCode, unit UnitType) bool {
	if unit == UnitTypeArmy {
		return regionCodeContained(r.AdjacentLand, adj)
	}
	if unit == UnitTypeFleet {
		return regionCodeContained(r.AdjacentWater, adj)
	}
	panic("IsAdjacent passed a non-unit UnitType")
}

// BuildMap Builds the map from file data
func BuildMap() (out map[RegionCode]*Region) {
	out = make(map[RegionCode]*Region)
	RegionIndex = make(map[RegionCode]*Region)
	byteData, err := data.Asset(mapLocation)
	if err != nil {
		log.Fatal(err)
	}
	linkedLocations := map[RegionCode][]RegionCode{
		"STP": []RegionCode{"STP", "STPsc", "STPnc"},
		"BUL": []RegionCode{"BUL", "BULsc", "BULnc"},
		"SPN": []RegionCode{"SPN", "SPNsc", "SPNnc"},
	}
	s := string(byteData[:])
	lines := strings.Split(s, "\n")
	lines = lines[:len(lines)-1]

	for _, line := range lines {
		if strings.HasPrefix(line, "#") || len(line) <= 1 {
			continue
		}
		token := strings.Split(line, "\\")
		isWater := (token[0] == "W")
		regCode := RegionCode(strings.TrimSpace(token[2]))

		reg := new(Region)
		if ereg, ok := RegionIndex[regCode]; ok {
			reg = ereg
		} else {
			reg.IsSC = (token[1] == "Y")
			reg.RegionID = regCode
			if link, ok := linkedLocations[regCode[:3]]; ok {
				reg.LinkedWith = link
			}
		}

		for _, adj := range strings.Split(token[3], ",") {
			adjCode := RegionCode(strings.TrimSpace(adj))
			if isWater {
				reg.AdjacentWater = append(reg.AdjacentWater, adjCode)
			} else {
				reg.AdjacentLand = append(reg.AdjacentLand, adjCode)
			}
		}

		RegionIndex[reg.RegionID] = reg
		out[reg.RegionID] = reg
	}
	isMapInitialized = true
	return
}

// IsMapBuilt returns whether the map has been constructed yet
func IsMapBuilt() bool {
	return isMapInitialized
}
