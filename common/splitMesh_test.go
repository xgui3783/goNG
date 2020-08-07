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

	tables := []struct {
		meshPtr         *Mesh
		meshMap         MeshMap
		vertexMap       VertexMap
		splitMeshConfig SplitMeshConfig
		expectedResult  map[string]int
	}{
		{
			&mesh,
			MeshMap{
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
			},
			VertexMap{
				uint32(0): "a",
				uint32(1): "a",
				uint32(2): "a",
				uint32(3): "b",
				uint32(4): "b",
				uint32(5): "b",
			},
			SplitMeshConfig{
				SplitMeshByVerticesPath: "",
				UntangleAmbiguityMethod: EMPTY_LABEL,
			},
			map[string]int{
				"a": 1,
				"b": 1,
			},
		},
		{
			&mesh,
			MeshMap{
				"a": []uint32{
					uint32(0),
					uint32(1),
					uint32(3),
				},
				"b": []uint32{
					uint32(2),
					uint32(4),
				},
				"c": []uint32{
					uint32(5),
				},
			},
			VertexMap{
				uint32(0): "a",
				uint32(1): "a",
				uint32(2): "b",
				uint32(3): "a",
				uint32(4): "a",
				uint32(5): "c",
			},
			SplitMeshConfig{
				SplitMeshByVerticesPath: "",
				UntangleAmbiguityMethod: MAJORITY_OR_FIRST_INDEX,
			},
			map[string]int{
				"a": 2,
				"b": 0,
				"c": 0,
			},
		},
	}

	for _, table := range tables {
		splitMeshesPtr := SplitMesh(table.meshPtr, &(table.meshMap), &(table.vertexMap), &(table.splitMeshConfig))
		for key, numFaces := range table.expectedResult {
			mesh := (*splitMeshesPtr)[key]

			if len(mesh.Faces) != numFaces {
				t.Errorf("label %v expected to have %v faces, but got %v faces", key, table.expectedResult[key], len(mesh.Faces))
			}
		}
	}
}
