package stlAscii

import "common"

func Import(manyFiles [][]byte) (returnMesh []common.Mesh) {
	for _, singleFile := range manyFiles {
		returnMesh = append(returnMesh, ParseAsciiStl(singleFile))
	}
	return
}

func Export(meshes []common.Mesh) (returnVal [][]byte) {
	for _, mesh := range meshes {
		returnVal = append(returnVal, WriteAsciiStlFromMesh(mesh))
	}
	return
}
