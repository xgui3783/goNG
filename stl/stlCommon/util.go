package stlCommon

import (
	"bytes"
	"encoding/binary"
	"gong/common"
	"math"
)

// STL files are quite verbose in that shared vertices are duplciated for each face
// this method is shared by both stlBinary and stlAscii, which suffer similar shortcomings
func PrepareStlAppendFace(mesh *common.Mesh) func([3]common.Vertex) {

	// Holds float32, joined and serialized into string
	mapVertexToIdx := make(map[string]uint32, 0)
	currentVertexIdx := uint32(0)

	return func(vertices [3]common.Vertex) {
		newFace := common.Face{}
		for vIdx, vertex := range vertices {
			vertexHash := serializeVertexTripletToString(vertex)
			if vertexIdx, ok := mapVertexToIdx[vertexHash]; !ok {
				mapVertexToIdx[vertexHash] = currentVertexIdx
				newFace[vIdx] = currentVertexIdx
				currentVertexIdx++
			} else {
				newFace[vIdx] = vertexIdx
			}
		}

		mesh.Faces = append(mesh.Faces, newFace)
	}
}

func AppendFloat32ToBuffer(input float32, bufPtr *bytes.Buffer) {
	binary.Write(bufPtr, binary.LittleEndian, math.Float32bits(input))
}

func serializeVertexTripletToString(vertex common.Vertex) string {
	bufPtr := new(bytes.Buffer)
	for _, coord := range vertex {
		AppendFloat32ToBuffer(coord, bufPtr)
	}
	return string(bufPtr.Bytes())
}
