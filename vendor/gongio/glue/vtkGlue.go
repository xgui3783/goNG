// +build vtk all

package glue

import "types/vtk"

const (
	VTK = "VTK"
)

func init() {
	registerParser(VTK, vtk.Import, vtk.Export)
}
