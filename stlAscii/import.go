package stlAscii

import "gong/common"

func Import(rootPath string) []common.Mesh {
	fileBytes := common.GetResource(rootPath)
	return []common.Mesh{ParseAsciiStl(fileBytes)}
}
