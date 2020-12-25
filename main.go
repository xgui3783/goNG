package main

import (
	"flag"
	"fmt"
	"gong/actions"
	"gong/common"
)

var validSrcFormats = []string{}
var validOutputFormats = []string{}

func init() {
	validSrcFormats = actions.GetSupportedIncTypes()
	validOutputFormats = actions.GetSupportedOutTypes()
}

func main() {
	srcFormatPtr := flag.String("srcFormat", "", getSrcFormatHelperText())
	srcPtr := flag.String("src", "", srcHelperText)
	outputFormatPtr := flag.String("outputFormat", "", getOutputFormatHelperText())
	dstPtr := flag.String("dst", "", dstHelperText)
	xformMatrix := flag.String("xformMatrix", "1,0,0,0,0,1,0,0,0,0,1,0", xformMatrixHelperText)
	flipTriangle := flag.Bool("forceTriangleFlip", false, flipTriangleHelperText)
	forceTriangleFlag := false
	splitMeshByVertex := flag.String("splitByVertexPath", "", splitMeshByVertexHelperText)
	validAmbiguousStrategies := []string{
		common.EMPTY_LABEL,
		common.MAJORITY_OR_FIRST_INDEX,
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
