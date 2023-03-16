package splitByVertex

import "common"

type SplitByVertexConfig struct {
	FilePath         string
	ResolveAmbiguity string
}

type MeshObj struct {
	mesh              common.Mesh
	vertexPtrToIdxMap map[*common.Vertex]uint32
}
type MeshMap map[string][]uint32
type VertexMap map[uint32]string
