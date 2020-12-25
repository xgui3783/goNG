// +build ng_mesh all

package actions

import "gong/ngPrecomputed"

const (
	NG_MESH = "NG_MESH"
)

func init() {
	registerParser(NG_MESH, ngPrecomputed.Import, ngPrecomputed.Export)
}
