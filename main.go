package main

import (
	"flag"
	"fmt"
	"gong/actions"
	"gong/detType"
)

func main() {
	validInputFormats := []string{
		detType.NG_MESH,
		detType.STL_ASCII,
		detType.STL_BINARY,
		detType.GII,
		detType.OBJ,
	}
	inputFormatHelperText := fmt.Sprintf("Format of the input file. If left empty, the program will try to deduce it by parsing the file extension. The possible values are %v", validInputFormats)
	inputFormatPtr := flag.String("inputFormat", "", inputFormatHelperText)

	inputSourceHelperText := "Source of input. May start with http:// , in which case, the program will first fetch the file, then parse it. If left empty it will accept STDIN (NYI)"
	inputSourcePtr := flag.String("inputSource", "", inputSourceHelperText)

	validOutputFormats := []string{
		detType.STL_BINARY,
		detType.STL_ASCII,
		detType.GII,
		detType.OBJ,
	}
	outputFormatHelperText := fmt.Sprintf("Format of the output file. If left empty, the program will try to deduce it by parsing the file extension. The possible values are %v", validOutputFormats)
	outputFormatPtr := flag.String("outputFormat", "", outputFormatHelperText)

	outputSourceHelperText := "Source of output. If left empty it output to STDOUT (NYI)"
	outputSourcePtr := flag.String("outputSource", "", outputSourceHelperText)

	flag.Parse()

	actions.Convert(*inputFormatPtr, *inputSourcePtr, *outputFormatPtr, *outputSourcePtr)
}
