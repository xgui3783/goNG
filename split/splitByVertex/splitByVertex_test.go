package splitByVertex

import (
	"common"
	"testing"
)

func TestSplitMesh(t *testing.T) {
	mesh := common.Mesh{
		Vertices: []common.Vertex{
			common.Vertex{0, 0, 0},
			common.Vertex{0, 1, 0},
			common.Vertex{0, 1, 1},

			common.Vertex{1, 0, 0},
			common.Vertex{1, 1, 0},
			common.Vertex{1, 1, 1},
		},
		Faces: []common.Face{
			common.Face{0, 1, 2},
			common.Face{3, 4, 5},
		},
	}

	tables := []struct {
		meshPtr             *common.Mesh
		meshMap             MeshMap
		vertexMap           VertexMap
		SplitByVertexConfig SplitByVertexConfig
		expectedResult      map[string]int
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
			SplitByVertexConfig{
				FilePath:         "",
				ResolveAmbiguity: EMPTY_LABEL,
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
			SplitByVertexConfig{
				FilePath:         "",
				ResolveAmbiguity: MAJORITY_OR_FIRST_INDEX,
			},
			map[string]int{
				"a": 2,
				"b": 0,
				"c": 0,
			},
		},
	}

	for _, table := range tables {
		splitMeshesPtr := splitMesh(table.meshPtr, &(table.meshMap), &(table.vertexMap), &(table.SplitByVertexConfig))
		for key, numFaces := range table.expectedResult {
			mesh := (*splitMeshesPtr)[key]

			if len(mesh.Faces) != numFaces {
				t.Errorf("label %v expected to have %v faces, but got %v faces", key, table.expectedResult[key], len(mesh.Faces))
			}
		}
	}
}
