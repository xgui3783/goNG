package detType

import (
	"gong/common"
	"path/filepath"
)

func InferTypeFromFilename(filename string) string {
	ext := filepath.Ext(filename)
	if ext == "" {
		return NG_MESH
	}
	if ext == ".vtk" {
		return VTK
	}
	if ext == ".stl" {
		stlBytes := common.GetResource(filename)
		if string(stlBytes[:5]) == "solid" {
			return STL_ASCII
		} else {
			return STL_BINARY
		}
	}
	if ext == ".obj" {
		return OBJ
	}
	if ext == ".gii" {
		return GII
	}
	panic("ext not recognised")
}

const (
	NG_MESH    = "NG_MESH"
	VTK        = "VTK"
	STL_ASCII  = "STL_ASCII"
	STL_BINARY = "STL_BINARY"
	GII        = "GII"
	OBJ        = "OBJ"
)
