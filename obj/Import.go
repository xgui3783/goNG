package obj

import (
	"fmt"
	"gong/common"
)

func Import(manyFiles [][]byte) (returnMesh []common.Mesh) {
	for _, singleFile := range manyFiles {
		returnMesh = append(returnMesh, ParseObjBytes(singleFile)...)
	}
	return
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
