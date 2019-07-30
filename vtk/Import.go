package vtk

import "gong/common"

func Import(rootPath string) []common.Mesh {
	fileBytes := common.GetResource(rootPath)
	return []common.Mesh{parseVtk(fileBytes)}
}

func Export(meshes []common.Mesh) [][]byte {
	returnBytes := [][]byte{}
	for _, mesh := range meshes {
		returnBytes = append(returnBytes, writeVtk(mesh))
	}
	return returnBytes
}
