package actions

import (
	"common"
)

var meshTypeToImportMap = map[string]func([][]byte) []common.Mesh{}
var meshTypeToExportMap = map[string]func([]common.Mesh) [][]byte{}

func registerParser(name string, importFnPtr func([][]byte) []common.Mesh, exportFnPtr func([]common.Mesh) [][]byte) {
	if importFnPtr != nil {
		meshTypeToImportMap[name] = importFnPtr
	}
	if exportFnPtr != nil {
		meshTypeToExportMap[name] = exportFnPtr
	}
}

func GetSupportedIncTypes() (returnVal []string) {
	for key := range meshTypeToImportMap {
		returnVal = append(returnVal, key)
	}
	return
}

func GetSupportedOutTypes() (returnVal []string) {
	for key := range meshTypeToExportMap {
		returnVal = append(returnVal, key)
	}
	return
}
