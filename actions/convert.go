package actions

import (
	"fmt"
	"gong/common"
	"gong/detType"
	"gong/ngPrecomputed"
	"gong/stl"
	"io/ioutil"
)

func Convert(inputFormat string, inputSource string, outputFormat string, outputDest string) {

	if inputSource == "" {
		panic("inputSource is empty\n")
	}

	if outputDest == "" {
		panic("outputSource is empty\n")
	}

	var mesh common.Mesh
	var incFileType string
	if inputFormat == "" {
		incFileType = detType.InferTypeFromFilename(inputSource)
	}
	switch incFileType {
	case detType.NG_MESH:
		mesh = ngPrecomputed.Import(inputSource, nil)
	default:
		panic("incoming file type other than NG_MESH is currently not supported\n")
	}

	var outBuffer []byte
	var outFileType string
	if outputFormat == "" {
		outFileType = detType.InferTypeFromFilename(outputDest)
	}
	switch outFileType {
	case detType.STL_ASCII:
		outBuffer = stl.WriteAsciiStlFromMesh(mesh)
	case detType.STL_BINARY:
		outBuffer = stl.WriteBinaryStlFromMesh(mesh)
	default:
		panic("out filename other than STL is not current supported\n")
	}

	err := ioutil.WriteFile(outputDest, outBuffer, 0644)
	if err != nil {
		panic(err)
	}
	fmt.Print("done\n")
}
