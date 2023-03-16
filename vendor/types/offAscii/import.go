package offAscii

import "common"

func Import(manyFiles [][]byte) (returnMesh []common.Mesh) {
	for _, singleFile := range manyFiles {
		returnMesh = append(returnMesh, ParseOffAscii(singleFile)...)
	}
	return
}
