package detType

import "path/filepath"

func InferTypeFromFilename(filename string) string {
	ext := filepath.Ext(filename)
	if ext == "" {
		return NG_MESH
	}
	if ext == ".vtk" {
		return VTK
	}
	if ext == ".stl" {
		return STL_ASCII
	}
	panic("ext not recognised")
}

const (
	NG_MESH    = "NG_MESH"
	VTK        = "VTK"
	STL_ASCII  = "STL_ASCII"
	STL_BINARY = "STL_BINARY"
)
