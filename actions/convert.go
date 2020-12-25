package actions

import (
	"fmt"
	"gong/common"
	"gong/detType"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
)

var inputBytes = make([]byte, 0)

func Convert(inputFormat string, inputSource string, outputFormat string, outputDest string, xformMatrixString string, flipTriangle bool, forceTriangleFlag bool, splitMeshConfig common.SplitMeshConfig) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintf(os.Stderr, "error: %v", r)
		}
	}()

	pathToMeshSplitVertexFile := splitMeshConfig.SplitMeshByVerticesPath

	if inputSource == "" {
		fmt.Printf("inputSource not provided, listening from stdin until EOF... \n")
		var d []byte
		for {
			_, err := fmt.Scan(&d)
			if err != nil {
				if err != io.EOF {
					panic(err)
				}
				break
			}
			// TODO too slow... maybe use byte writer or something?
			inputBytes = append(inputBytes, d...)
		}
	}

	if outputDest == "" {
		panic("outputSource is empty\n")
	}

	// verify input format
	var incFileType string
	if inputFormat == "" {
		if inputSource == "" {
			panic("if stdin is used to provide src, -srcFormat must be defined,")
		}
		incFileType = detType.InferTypeFromFilename(inputSource)
	} else {
		incFileType = inputFormat
	}
	importFn, ok := meshTypeToImportMap[incFileType]
	if !ok {
		panicText := fmt.Sprintf("intput type %v not supported", incFileType)
		panic(panicText)
	}

	// verify output format
	var outFileType string
	if outputFormat == "" {
		outFileType = detType.InferTypeFromFilename(outputDest)
	} else {
		outFileType = outputFormat
	}
	exportFn, ok := meshTypeToExportMap[outFileType]
	if len(inputBytes) == 0 {
		inputBytes = common.GetResource(inputSource)
	}
	meshes := importFn([][]byte{inputBytes})

	// TODO move transform to build config
	// json marshalling makes wasm build difficult
	var xformMatrix common.TransformationMatrix
	xformMatrix.ParseCommaDelimitedString(xformMatrixString)

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
	fragmentOutBuffer := map[string][]byte{}
	outBuffer = exportFn(meshes)

	for key, mesh := range submeshNameToMeshMap {
		fragmentOutBuffer[key] = exportFn([]common.Mesh{mesh})[0]
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
