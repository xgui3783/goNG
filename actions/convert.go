package actions

import (
	"fmt"
	"gong/common"
	"gong/detType"
	"gong/gii"
	"gong/ngPrecomputed"
	"gong/obj"
	"gong/stlAscii"
	"gong/stlBinary"
	"gong/vtk"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
)

func Convert(inputFormat string, inputSource string, outputFormat string, outputDest string, xformMatrixString string) {

	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintf(os.Stderr, "error: %v", r)
		}
	}()

	if inputSource == "" {
		panic("inputSource is empty\n")
	}

	if outputDest == "" {
		panic("outputSource is empty\n")
	}

	var xformMatrix common.TransformationMatrix
	xformMatrix.ParseCommaDelimitedString(xformMatrixString)

	var meshes []common.Mesh
	var incFileType string
	if inputFormat == "" {
		incFileType = detType.InferTypeFromFilename(inputSource)
	} else {
		incFileType = inputFormat
	}
	switch incFileType {
	case detType.NG_MESH:
		meshes = ngPrecomputed.Import(inputSource, nil)
	case detType.STL_ASCII:
		meshes = stlAscii.Import(inputSource)
	case detType.GII:
		meshes = gii.Import(inputSource)
	case detType.OBJ:
		meshes = obj.Import(inputSource)
	case detType.VTK:
		meshes = vtk.Import(inputSource)
	default:
		panic("incoming file type other than NG_MESH is currently not supported\n")
	}

	for mIdx, _ := range meshes {
		for vIdx, _ := range meshes[mIdx].Vertices {
			meshes[mIdx].Vertices[vIdx].Transform(xformMatrix)
		}

		if xformMatrix.Det() < 0 {
			meshes[mIdx].FlipTriangleOrder()
		}
	}

	var outBuffer []([]byte)
	var outFileType string
	if outputFormat == "" {
		outFileType = detType.InferTypeFromFilename(outputDest)
	} else {
		outFileType = outputFormat
	}
	switch outFileType {
	case detType.STL_ASCII:
		for idx, mesh := range meshes {
			outBuffer = append(outBuffer, stlAscii.WriteAsciiStlFromMesh(mesh, common.MeshMetadata{Index: idx}))
		}
	case detType.STL_BINARY:
		for idx, mesh := range meshes {
			outBuffer = append(outBuffer, stlBinary.WriteBinaryStlFromMesh(mesh, common.MeshMetadata{Index: idx}))
		}
	case detType.OBJ:
		outBuffer = obj.Export(meshes)
	case detType.GII:
		outBuffer = gii.Export(meshes)
	case detType.VTK:
		outBuffer = vtk.Export(meshes)
	case detType.NG_MESH:
		outBuffer = ngPrecomputed.Export(meshes)
	default:
		panic(fmt.Sprintf("ouputFormat %v is not current supported\n", outFileType))
	}

	if len(outBuffer) == 1 {
		writeBytesToFile(outputDest, outBuffer[0])
	} else {
		for idx, bytes := range outBuffer {
			ext := filepath.Ext(outputDest)
			regexString := fmt.Sprintf("(%v)?$", ext)
			re := regexp.MustCompile(regexString)
			filename := re.ReplaceAllStringFunc(outputDest, func(ext string) string {
				return fmt.Sprintf("_%d%v", idx, ext)
			})
			writeBytesToFile(filename, bytes)
		}
	}
	fmt.Print("done\n")
}

func writeBytesToFile(filename string, buf []byte) {
	err := ioutil.WriteFile(filename, buf, 0644)
	if err != nil {
		panic(err)
	}
}
