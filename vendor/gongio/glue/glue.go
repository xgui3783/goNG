package glue

import (
	"common"
)

var MeshTypeToImportMap = map[string]func([][]byte) []common.Mesh{}
var MeshTypeToExportMap = map[string]func([]common.Mesh) [][]byte{}

func registerParser(name string, importFnPtr func([][]byte) []common.Mesh, exportFnPtr func([]common.Mesh) [][]byte) {
	if importFnPtr != nil {
		MeshTypeToImportMap[name] = importFnPtr
	}
	if exportFnPtr != nil {
		MeshTypeToExportMap[name] = exportFnPtr
	}
}

func GetSupportedIncTypes() (returnVal []string) {
	for key := range MeshTypeToImportMap {
		returnVal = append(returnVal, key)
	}
	return
}

func GetSupportedOutTypes() (returnVal []string) {
	for key := range MeshTypeToExportMap {
		returnVal = append(returnVal, key)
	}
	return
}
