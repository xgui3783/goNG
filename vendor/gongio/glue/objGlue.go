// +build obj all

package glue

import "types/obj"

const (
	OBJ = "OBJ"
)

func init() {
	registerParser(OBJ, obj.Import, obj.Export)
}
