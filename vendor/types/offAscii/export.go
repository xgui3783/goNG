package offAscii

import (
	"fmt"
	"common"
)

func Export(meshes []common.Mesh) [][]byte {

	// header
	headerString := fmt.Sprintf("OFF\n# %v\n", common.HEADER)
	header := []byte(headerString)

	returnBytes := [][]byte{}
	for _, mesh := range meshes {
		appendBytes := append(header, WriteMeshToBytes(mesh)...)
		returnBytes = append(returnBytes, appendBytes)
	}
	return returnBytes
}
