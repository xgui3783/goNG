package info

import (
	"flag"
	"gongio"
	"os"
)

func init() {

}

const SubCmd = "info"
const HelperText = "temp helper text info"

func Parse(fs *flag.FlagSet) error {
	ii := gongio.NewInputInterface()

	ii.OptOut("OutFormat")

	outputHelperText := "outputHelperText place holder"
	ii.OutHelperText = &outputHelperText
	ii.SetupFlag(fs)

	numOfLod := fs.Int("numberOfLod", 0, "Number of LOD")
	fs.Parse(os.Args[2:])

	return writeInfo(ii, *numOfLod)
}
