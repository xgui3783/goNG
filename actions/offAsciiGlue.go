// +build off_ascii all

package actions

import "gong/offAscii"

const (
	OFF_ASCII = "OFF_ASCII"
)

func init() {
	registerParser(OFF_ASCII, offAscii.Import, offAscii.Export)
}
