// +build stl_binary all

package glue

import "types/stl/stlBinary"

const (
	STL_BINARY = "STL_BINARY"
)

func init() {
	registerParser(STL_BINARY, stlBinary.Import, stlBinary.Export)
}
