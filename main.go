package main

import (
	"flag"
	"fmt"
	"gong/actions"
	"gong/detType"
)

func main() {
	validSrcFormats := []string{
		detType.NG_MESH,
		detType.STL_ASCII,
		detType.STL_BINARY,
		detType.GII,
		detType.OBJ,
	}
	srcFormatHelperText := fmt.Sprintf("Format of the input file. If left empty, the program will try to deduce it by parsing the file extension. The possible values are %v", validSrcFormats)
	srcFormatPtr := flag.String("srcFormat", "", srcFormatHelperText)

	srcHelperText := "Source of input. May start with http:// , in which case, the program will first fetch the file, then parse it. If left empty it will accept STDIN (NYI)"
	srcPtr := flag.String("src", "", srcHelperText)

	validOutputFormats := []string{
		detType.NG_MESH,
		detType.STL_BINARY,
		detType.STL_ASCII,
		detType.GII,
		detType.OBJ,
	}
	outputFormatHelperText := fmt.Sprintf("Format of the output file. If left empty, the program will try to deduce it by parsing the file extension. The possible values are %v", validOutputFormats)
	outputFormatPtr := flag.String("outputFormat", "", outputFormatHelperText)

	dstHelperText := "Source of output. If left empty it output to STDOUT (NYI)"
	dstPtr := flag.String("dst", "", dstHelperText)

	xformMatrixHelperText := "4x3, organised row major, comma separated. 1,0,0,0,0,1,0,0,0,0,1,0 == identity, 1,0,0,10,0,1,0,11,0,0,1,12 === same scale, but translated by 10, 11, 12. Last row assumed to be 0,0,0,1"
	xformMatrix := flag.String("xformMatrix", "1,0,0,0,0,1,0,0,0,0,1,0", xformMatrixHelperText)

	flag.Parse()

	actions.Convert(*srcFormatPtr, *srcPtr, *outputFormatPtr, *dstPtr, *xformMatrix)
}
