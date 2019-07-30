package obj

import (
	"fmt"
	"gong/common"
)

func Import(rootPath string) []common.Mesh {
	fileBytes := common.GetResource(rootPath)
	return ParseObjBytes(fileBytes)
}

func Export(meshes []common.Mesh) [][]byte {

	// header
	headerString := fmt.Sprintf("# %v\n", common.HEADER)
	header := []byte(headerString)

	returnBytes := [][]byte{}
	for _, mesh := range meshes {
		appendBytes := append(header, WriteMeshToBytes(mesh)...)
		returnBytes = append(returnBytes, appendBytes)
	}
	return returnBytes
}
