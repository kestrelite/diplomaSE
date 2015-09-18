package gamestate

import (
	"strings"
	"testing"
)

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func TestBuildMap(t *testing.T) {
	rMap := BuildMap()
	missingAdjs := make(map[string]bool)

	for code, reg := range rMap {
		// Map's code must equal region's code
		if reg.RegionCode != code {
			t.Error("Map Construction Error: Map code (" + code + ") did not equal region code (" + reg.RegionCode + ").")
		}

		// If a region is part of a linked set, then the other items
		// in that set must both exist and have identical linked sets
		if len(reg.LinkedWith) != 0 {
			for _, link := range reg.LinkedWith {
				if _, ok := rMap[link]; !ok {
					t.Error("Linking Failure: A link in " + code + " does not exist: " + link)
					t.FailNow()
				}
				if strings.Join(reg.LinkedWith, "|") != strings.Join(rMap[link].LinkedWith, "|") {
					t.Error("Linking Failure: A linked country set failed to have identical linked sets: [" +
						strings.Join(reg.LinkedWith, "|") + "]; [" + strings.Join(rMap[link].LinkedWith, "|") + "]")
				}
			}
		}

		// If a region lists another region as an adjacency, that region must exist.
		// Additionally, that region must share a connection back.
		// Also, that region can't be an adjacency to itself.
		for _, adj := range append(reg.AdjacentLand, reg.AdjacentWater...) {
			if adjReg, ok := rMap[adj]; !ok {
				t.Error("Nonexistent Adjacency: " + reg.RegionCode + " lists " + adj + " as an adjacency, but it DNE.")
				missingAdjs[adj] = true
			} else {
				if !contains(append(adjReg.AdjacentLand, adjReg.AdjacentWater...), reg.RegionCode) {
					t.Error("Inconsistent Map: " + reg.RegionCode + "'s adjacency " + adj + " does not list " + reg.RegionCode + " in kind.")
				}
				if adjReg.RegionCode == reg.RegionCode {
					t.Error("Self-Adjacency: " + reg.RegionCode + " is adjacent to itself.")
				}
			}
		}
	}

	if len(missingAdjs) > 0 {
		keys := make([]string, 0, len(missingAdjs))
		for k := range missingAdjs {
			keys = append(keys, k)
		}
		t.Error("Reduced list of DNEs: " + strings.Join(keys, " "))
	}
}
