
import (
	"flag"
)

func t() {

	splitMeshByVertex := flag.String("splitByVertexPath", "", splitMeshByVertexHelperText)
	validAmbiguousStrategies := []string{
		common.EMPTY_LABEL,
		common.MAJORITY_OR_FIRST_INDEX,
	}

	splitMeshAmbiguousStrategy := flag.String("splitMeshAmbiguousStrategy", common.EMPTY_LABEL, splitMeshAmbiguousStrategyHelperTxt)

	splitMeshConfig := common.SplitMeshConfig{
		UntangleAmbiguityMethod: *splitMeshAmbiguousStrategy,
		SplitMeshByVerticesPath: *splitMeshByVertex,
	}
}