package stlBinary

import "gong/common"

func Import(manyFiles [][]byte) (returnMesh []common.Mesh) {
	for _, singleFile := range manyFiles {
		meshPtr := parseBinaryStl(singleFile)
		returnMesh = append(returnMesh, *meshPtr)
	}
	return
}

func Export(meshes []common.Mesh) (returnVal [][]byte) {
	for _, mesh := range meshes {
		returnVal = append(returnVal, WriteBinaryStlFromMesh(mesh))
	}
	return
}
