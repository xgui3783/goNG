package actions

import (
	"fmt"
	"gong/common"
	"gong/detType"
	"gong/gii"
	"gong/ngPrecomputed"
	"gong/obj"
	"gong/offAscii"
	"gong/stlAscii"
	"gong/stlBinary"
	"gong/vtk"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
)

func Convert(inputFormat string, inputSource string, outputFormat string, outputDest string, xformMatrixString string, flipTriangle bool, forceTriangleFlag bool, splitMeshConfig common.SplitMeshConfig) {

	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintf(os.Stderr, "error: %v", r)
		}
	}()

	pathToMeshSplitVertexFile := splitMeshConfig.SplitMeshByVerticesPath

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
	case detType.OFF_ASCII:
		meshes = offAscii.Import(inputSource)
	default:
		panic("incoming file type other than NG_MESH is currently not supported\n")
	}

	for mIdx, _ := range meshes {
		for vIdx, _ := range meshes[mIdx].Vertices {
			meshes[mIdx].Vertices[vIdx].Transform(xformMatrix)
		}

		if forceTriangleFlag {
			if flipTriangle {
				meshes[mIdx].FlipTriangleOrder()
			}
		} else {
			if xformMatrix.Det() < 0 {
				meshes[mIdx].FlipTriangleOrder()
			}
		}
	}

	submeshNameToMeshMap := map[string]common.Mesh{}
	if pathToMeshSplitVertexFile != "" {
		meshMap, vertexMap := common.ProcessSplitMeshByVertexfile(pathToMeshSplitVertexFile)

		for meshIdx, mesh := range meshes {
			splitMeshMap := common.SplitMesh(&mesh, &meshMap, &vertexMap, &splitMeshConfig)
			for key, val := range *splitMeshMap {
				submeshName := fmt.Sprintf("%v_%v", meshIdx, key)
				submeshNameToMeshMap[submeshName] = val
			}
		}
	}

	var outBuffer []([]byte)
	var outFileType string
	fragmentOutBuffer := map[string][]byte{}
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
		if len(submeshNameToMeshMap) > 0 {
			fmt.Printf("submeshes not yet implemented for STL_ASCII")
		}
	case detType.STL_BINARY:
		for idx, mesh := range meshes {
			outBuffer = append(outBuffer, stlBinary.WriteBinaryStlFromMesh(mesh, common.MeshMetadata{Index: idx}))
		}
		if len(submeshNameToMeshMap) > 0 {
			fmt.Printf("submeshes not yet implemented for STL_BINARY")
		}
	case detType.OBJ:
		outBuffer = obj.Export(meshes)
		for key, mesh := range submeshNameToMeshMap {
			fragmentOutBuffer[key] = obj.Export([]common.Mesh{mesh})[0]
		}
	case detType.GII:
		outBuffer = gii.Export(meshes)
		for key, mesh := range submeshNameToMeshMap {
			fragmentOutBuffer[key] = gii.Export([]common.Mesh{mesh})[0]
		}
	case detType.VTK:
		outBuffer = vtk.Export(meshes)
		for key, mesh := range submeshNameToMeshMap {
			fragmentOutBuffer[key] = vtk.Export([]common.Mesh{mesh})[0]
		}
	case detType.NG_MESH:
		outBuffer = ngPrecomputed.Export(meshes)
		for key, mesh := range submeshNameToMeshMap {
			fragmentOutBuffer[key] = ngPrecomputed.Export([]common.Mesh{mesh})[0]
		}
	case detType.OFF_ASCII:
		fallthrough
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

	if len(fragmentOutBuffer) > 0 {
		fragmentDir := fmt.Sprintf("%v_fragments/", outputDest)
		if fi, err := os.Stat(fragmentDir); os.IsNotExist(err) {
			if err := os.Mkdir(fragmentDir, 0755); err != nil {
				panic(err)
			}
		} else if !(fi.Mode().IsDir()) {
			panicText := fmt.Sprintf("%v path already exist", fragmentDir)
			panic(panicText)
		}

		for key, bytes := range fragmentOutBuffer {

			fragmentFilename := path.Join(fragmentDir, key)

			ext := filepath.Ext(outputDest)
			filename := fmt.Sprintf("%v%v", fragmentFilename, ext)

			writeBytesToFile(filename, bytes)
		}
	}

	fmt.Fprintf(os.Stderr, "done")
}

func writeBytesToFile(filename string, buf []byte) {
	err := ioutil.WriteFile(filename, buf, 0644)
	if err != nil {
		panic(err)
	}
}
