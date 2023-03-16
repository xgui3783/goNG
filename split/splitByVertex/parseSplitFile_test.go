package splitByVertex

import (
	"testing"
)

func TestParseSplitMeshByVertexLine(t *testing.T) {
	testLines := []string{
		"1 a",
		"2 a",
		"3 b",
		"",
		"5",
	}

	meshMap := MeshMap{}
	vertexMap := VertexMap{}
	for _, line := range testLines {
		parseSplitMeshByVertexLine(line, &meshMap, &vertexMap)
	}

	// testing meshMap a
	meshMapA, ok := meshMap["a"]
	if ok == false {
		t.Errorf("meshmap a is not defined")
	}

	expectedAMap := [2]uint32{1, 2}
	for idx, val := range expectedAMap {
		if val != meshMapA[idx] {
			t.Errorf("meshmap a idx %v comparison failed, expected %v, got %v", idx, val, meshMap["a"][idx])
		}
	}

	expectedBMap := [1]uint32{3}
	for idx, val := range expectedBMap {
		if val != expectedBMap[idx] {
			t.Errorf("meshmap b idx %v comparison failed, expected %v, got %v", idx, val, meshMap["a"][idx])
		}
	}

	expectedVertexMap := VertexMap{1: "a", 2: "a", 3: "b"}
	for key, val := range expectedVertexMap {
		if val != vertexMap[key] {
			t.Errorf("vertex map not as expeted, for accessor %v, expected %v, got %v", key, val, vertexMap[key])
		}
	}
}
