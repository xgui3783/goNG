package vtk

import "common"

func Import(manyFiles [][]byte) (returnMesh []common.Mesh) {
	for _, singleFile := range manyFiles {
		returnMesh = append(returnMesh, parseVtk(singleFile))
	}
	return
}

func Export(meshes []common.Mesh) [][]byte {
	returnBytes := [][]byte{}
	for _, mesh := range meshes {
		returnBytes = append(returnBytes, writeVtk(mesh))
	}
	return returnBytes
}
