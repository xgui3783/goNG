// +build stl_ascii all

package actions

import "gong/stl/stlAscii"

const (
	STL_ASCII = "STL_ASCII"
)

func init() {
	registerParser(STL_ASCII, stlAscii.Import, stlAscii.Export)
}
