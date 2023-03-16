// +build off_ascii all

package glue

import "types/offAscii"

const (
	OFF_ASCII = "OFF_ASCII"
)

func init() {
	registerParser(OFF_ASCII, offAscii.Import, offAscii.Export)
}
