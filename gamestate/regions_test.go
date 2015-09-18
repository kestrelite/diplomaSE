package gamestate

import "testing"

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func joinRegionCode(r []RegionCode, del string) string {
	out := ""
	for _, v := range r {
		out += string(v) + del
	}
	return out
}

func regionCodeContained(r []RegionCode, rt RegionCode) bool {
	for _, v := range r {
		if v == rt {
			return true
		}
	}
	return false
}

func TestBuildMap(t *testing.T) {
	rMap := BuildMap()
	missingAdjs := make(map[RegionCode]bool)

	for code, reg := range rMap {
		// Map's code must equal region's code
		if reg.RegionID != code {
			t.Error("Map Construction Error: Map code (" + code + ") did not equal region code (" + reg.RegionID + ").")
		}

		// If a region is part of a linked set, then the other items
		// in that set must both exist and have identical linked sets
		if len(reg.LinkedWith) != 0 {
			for _, link := range reg.LinkedWith {
				if _, ok := rMap[link]; !ok {
					t.Error("Linking Failure: A link in " + code + " does not exist: " + link)
					t.FailNow()
				}
				if joinRegionCode(reg.LinkedWith, ",") != joinRegionCode(rMap[link].LinkedWith, ",") {
					t.Error("Linking Failure: A linked country set failed to have identical linked sets: [" +
						joinRegionCode(reg.LinkedWith, ",") + "]; [" + joinRegionCode(rMap[link].LinkedWith, ",") + "]")
				}
			}
		}

		// If a region lists another region as an adjacency, that region must exist.
		// Additionally, that region must share a connection back.
		// Also, that region can't be an adjacency to itself.
		for _, adj := range append(reg.AdjacentLand, reg.AdjacentWater...) {
			if adjReg, ok := rMap[adj]; !ok {
				t.Error("Nonexistent Adjacency: " + reg.RegionID + " lists " + adj + " as an adjacency, but it DNE.")
				missingAdjs[adj] = true
			} else {
				if !regionCodeContained(append(adjReg.AdjacentLand, adjReg.AdjacentWater...), reg.RegionID) {
					t.Error("Inconsistent Map: " + reg.RegionID + "'s adjacency " + adj + " does not list " + reg.RegionID + " in kind.")
				}
				if adjReg.RegionID == reg.RegionID {
					t.Error("Self-Adjacency: " + reg.RegionID + " is adjacent to itself.")
				}
			}
		}
	}

	if len(missingAdjs) > 0 {
		keys := make([]RegionCode, 0, len(missingAdjs))
		for k := range missingAdjs {
			keys = append(keys, k)
		}
		t.Error("Reduced list of DNEs: " + joinRegionCode(keys, ","))
	}
}
