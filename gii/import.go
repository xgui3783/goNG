package gii

import "gong/common"

func Import(rootPath string) []common.Mesh {
	fileBytes := common.GetResource(rootPath)
	return ParseGii(fileBytes)
}

func Export(meshes []common.Mesh) [][]byte {
	outBytes := [][]byte{}
	for _, mesh := range meshes {
		outBytes = append(outBytes, WriteGii(mesh))
	}
	return outBytes
}
