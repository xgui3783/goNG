// +build gii all

package actions

import "gong/gii"

const (
	GII = "GII"
)

func init() {
	registerParser(GII, gii.Import, gii.Export)
}
