package gii

import "common"

func Import(manyFiles [][]byte) (returnMesh []common.Mesh) {
	for _, singleFile := range manyFiles {
		returnMesh = append(returnMesh, ParseGii(singleFile)...)
	}
	return
}

func Export(meshes []common.Mesh) [][]byte {
	outBytes := [][]byte{}
	for _, mesh := range meshes {
		outBytes = append(outBytes, WriteGii(mesh))
	}
	return outBytes
}
