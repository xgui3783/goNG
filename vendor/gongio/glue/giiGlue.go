// +build gii all

package glue

import "types/gii"

const (
	GII = "GII"
)

func init() {
	registerParser(GII, gii.Import, gii.Export)
}
