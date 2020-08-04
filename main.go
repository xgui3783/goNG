package main

import (
	"flag"
	"fmt"
	"gong/actions"
	"gong/common"
	"gong/detType"
)

func main() {
	validSrcFormats := []string{
		detType.NG_MESH,
		detType.STL_ASCII,
		detType.STL_BINARY,
		detType.GII,
		detType.OBJ,
		detType.OFF_ASCII,
	}
	srcFormatHelperText := fmt.Sprintf("Format of the input file. If left empty, the program will try to deduce it by parsing the file extension. The possible values are %v", validSrcFormats)
	srcFormatPtr := flag.String("srcFormat", "", srcFormatHelperText)

	srcHelperText := "Source of input. May start with http:// , in which case, the program will first fetch the file, then parse it." // If left empty it will accept STDIN
	srcPtr := flag.String("src", "", srcHelperText)

	validOutputFormats := []string{
		detType.NG_MESH,
		detType.STL_BINARY,
		detType.STL_ASCII,
		detType.GII,
		detType.OBJ,
		detType.OFF_ASCII,
	}
	outputFormatHelperText := fmt.Sprintf("Format of the output file. If left empty, the program will try to deduce it by parsing the file extension. The possible values are %v", validOutputFormats)
	outputFormatPtr := flag.String("outputFormat", "", outputFormatHelperText)

	dstHelperText := "Dest of output." //  If left empty it output to STDOUT
	dstPtr := flag.String("dst", "", dstHelperText)

	xformMatrixHelperText := "4x3, organised row major, comma separated. 1,0,0,0,0,1,0,0,0,0,1,0 == identity, 1,0,0,10,0,1,0,11,0,0,1,12 === same scale, but translated by 10, 11, 12. Last row assumed to be 0,0,0,1"
	xformMatrix := flag.String("xformMatrix", "1,0,0,0,0,1,0,0,0,0,1,0", xformMatrixHelperText)

	flipTriangleHelperText := `Forces flip of triangle order.
By default, triangles will be flipped if xformMatrix determinant is less than 0.
This option will overwrite the default behaviour

Usage: -flipTriangle -flipTriangle=false
`
	flipTriangle := flag.Bool("forceTriangleFlip", false, flipTriangleHelperText)

	forceTriangleFlag := false

	splitMeshByVertexHelperText := `Path to a text file outlining how the mesh should be split by vertex
The text file should follow the format:

0 label_a
1 label_a
2 label_b

where index_# can be parsed as uint32 and label_# should be [a-zA-Z0-9-_]+

In the above example, vertices with index 0 and 1 belong to label_a, but vertex with index 2 belong to label_b 
`

	splitMeshByVertex := flag.String("splitByVertexPath", "", splitMeshByVertexHelperText)

	validAmbiguousStrategies := []string{
		common.EMPTY_LABEL,
	}
	splitMeshAmbiguousStrategyHelperTxt := fmt.Sprintf(`Strategy when there are ambiguous triangles %v`, validAmbiguousStrategies)
	splitMeshAmbiguousStrategy := flag.String("splitMeshAmbiguousStrategy", common.EMPTY_LABEL, splitMeshAmbiguousStrategyHelperTxt)

	flag.Parse()

	flag.Visit(func(f *flag.Flag) {
		if f.Name == "forceTriangleFlip" {
			forceTriangleFlag = true
		}
	})

	splitMeshConfig := common.SplitMeshConfig{
		UntangleAmbiguityMethod: *splitMeshAmbiguousStrategy,
		SplitMeshByVerticesPath: *splitMeshByVertex,
	}
	actions.Convert(*srcFormatPtr, *srcPtr, *outputFormatPtr, *dstPtr, *xformMatrix, *flipTriangle, forceTriangleFlag, splitMeshConfig)
}
