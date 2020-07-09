package offAscii

import "gong/common"

func Import(rootPath string) []common.Mesh {
	fileBytes := common.GetResource(rootPath)
	return ParseOffAscii(fileBytes)
}
