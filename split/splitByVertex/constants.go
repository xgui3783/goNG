package splitByVertex

import "fmt"

var splitMeshByVertexHelperText = fmt.Sprintf(
	`%v only. Ignored otherwise.

Path to a text file outlining how the mesh should be split by vertex.
The text file should follow the format:

0 label_a
1 label_a
2 label_b
.
.
.
[index] [label]

where [index] can be parsed as uint32 and [label] should satisfy the pattern [a-zA-Z0-9-_]+
`, MethodName)

var splitMeshAmbiguousStrategyHelperTxt = fmt.Sprintf(
	`%v only. Ignored otherwise.

Strategy when there are ambiguous triangles %v`,
	MethodName,
	[]string{
		EMPTY_LABEL,
		MAJORITY_OR_FIRST_INDEX,
	},
)
