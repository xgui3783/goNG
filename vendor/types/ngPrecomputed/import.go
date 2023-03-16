package ngPrecomputed

import (
	"common"
)

func Import(manyFiles [][]byte) (returnMesh []common.Mesh) {
	for _, singleFile := range manyFiles {
		returnMesh = append(returnMesh, ParseFragmentBuffer(singleFile))
	}
	return
}

func Export(meshes []common.Mesh) [][]byte {
	returnBytes := [][]byte{}
	for _, mesh := range meshes {
		returnBytes = append(returnBytes, WriteFragmentFromMesh(mesh))
	}
	return returnBytes
}
