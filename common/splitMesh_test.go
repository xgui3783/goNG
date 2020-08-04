package common

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

func TestSplitMesh(t *testing.T) {
	mesh := Mesh{
		Vertices: []Vertex{
			Vertex{0, 0, 0},
			Vertex{0, 1, 0},
			Vertex{0, 1, 1},

			Vertex{1, 0, 0},
			Vertex{1, 1, 0},
			Vertex{1, 1, 1},
		},
		Faces: []Face{
			Face{0, 1, 2},
			Face{3, 4, 5},
		},
	}

	meshMap := MeshMap{
		"a": []uint32{
			uint32(0),
			uint32(1),
			uint32(2),
		},
		"b": []uint32{
			uint32(3),
			uint32(4),
			uint32(5),
		},
	}

	vertexMap := VertexMap{
		uint32(0): "a",
		uint32(1): "a",
		uint32(2): "a",
		uint32(3): "b",
		uint32(4): "b",
		uint32(5): "b",
	}

	splitMeshesPtr := SplitMesh(&mesh, &meshMap, &vertexMap)

	if len(*splitMeshesPtr) != 2 {
		t.Errorf("expected split mesh ptr to have length 2, but have length %v", len(*splitMeshesPtr))
	}
	for key := range *splitMeshesPtr {
		if key != "a" && key != "b" {
			t.Errorf("key other than a or b are present %v", key)
		}
	}
}
