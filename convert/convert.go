package convert

import (
	"flag"
	"gongio"
	"os"
)

const SubCmd = "convert"
const HelperText = "temp convert helper text"

var validSrcFormats = []string{}
var validOutputFormats = []string{}

func SetFormats(src []string, output []string) {
	validSrcFormats = src
	validOutputFormats = output
}

func Parse(fs *flag.FlagSet) error {

	srcFormatHelperText := getSrcFormatHelperText()
	outFormatHelperText := getOutputFormatHelperText()
	ii := gongio.NewInputInterface()
	ii.SrcFormatHelperText = &srcFormatHelperText
	ii.SrcHelperText = &srcHelperText
	ii.OutFormatHelperText = &outFormatHelperText
	ii.OutHelperText = &dstHelperText

	ii.SetupFlag(fs)

	xformMatrix := (*fs).String("xformMatrix", "1,0,0,0,0,1,0,0,0,0,1,0", xformMatrixHelperText)
	flipTriangle := (*fs).Bool("forceTriangleFlip", false, flipTriangleHelperText)
	forceTriangleFlag := false

	fs.Parse(os.Args[2:])

	(*fs).Visit(func(f *flag.Flag) {
		if f.Name == "forceTriangleFlip" {
			forceTriangleFlag = true
		}
	})

	return convert(ii, *xformMatrix, *flipTriangle, forceTriangleFlag)
}
