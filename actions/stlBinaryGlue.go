// +build stl_binary all

package actions

import "gong/stl/stlBinary"

const (
	STL_BINARY = "STL_BINARY"
)

func init() {
	registerParser(STL_BINARY, stlBinary.Import, stlBinary.Export)
}
