// +build obj all

package actions

import "gong/obj"

const (
	OBJ = "OBJ"
)

func init() {
	registerParser(OBJ, obj.Import, obj.Export)
}
