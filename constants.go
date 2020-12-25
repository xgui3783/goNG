package main

import "fmt"

func getSrcFormatHelperText() string {
	return fmt.Sprintf(`Format of the input file.
If left empty, the program will try to deduce it by parsing the file extension.
The possible values are:

%v
`, validSrcFormats)
}

var srcHelperText = `Source of input.

May start with http , in which case, the program will first fetch the file, then parse it.
` // If left empty it will accept STDIN

func getOutputFormatHelperText() string {
	return fmt.Sprintf(`Format of the output file.
If left empty, the program will try to deduce it by parsing the file extension.
The possible values are:

%v
`, validOutputFormats)
}

var dstHelperText = `Dest of output.
` //  If left empty it output to STDOUT

var xformMatrixHelperText = `4x3 transformation matrix, organised row major, comma separated.
Last row assumed to be 0,0,0,1

e.g.
1,0,0,0,0,1,0,0,0,0,1,0 == identity
1,0,0,10,0,1,0,11,0,0,1,12 === same scale, but translated by 10, 11, 12
`

var flipTriangleHelperText = `Forces flip of triangle order.
By default, triangles will be flipped if xformMatrix determinant is less than 0.
Warning: setting this flag will overwrite the default behaviour

e.g.
[default: not setting the flag] triangles will be flipped if transform matrix determinant is < 0. Do not flip otherwise
-flipTriangle always flip triangle, after applying transformation matrix, if supplied
-flipTriangle=false never flip the triangle
`

var splitMeshByVertexHelperText = `Path to a text file outlining how the mesh should be split by vertex.
The text file should follow the format:

0 label_a
1 label_a
2 label_b
.
.
.
[index] [label]

where [index] can be parsed as uint32 and [label] should satisfy the pattern [a-zA-Z0-9-_]+
`
