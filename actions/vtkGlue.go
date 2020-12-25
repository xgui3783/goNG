// +build vtk all

package actions

import "gong/vtk"

const (
	VTK = "VTK"
)

func init() {
	registerParser(VTK, vtk.Import, vtk.Export)
}
