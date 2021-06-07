// +build stl_ascii all

package glue

import "types/stl/stlAscii"

const (
	STL_ASCII = "STL_ASCII"
)

func init() {
	registerParser(STL_ASCII, stlAscii.Import, stlAscii.Export)
}
