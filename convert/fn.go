package convert

import (
	"common"
	"errors"
	"fmt"
	"gongio"
	"path/filepath"
	"regexp"
)

func convert(ii *gongio.InputInterface, xformMatrixString string, flipTriangle bool, forceTriangleFlag bool) (returnError error) {
	defer func() {
		if r := recover(); r != nil {
			errorText := fmt.Sprintf("error: %v", r)
			returnError = errors.New(errorText)
			return
		}
	}()
	meshes, err := ii.GetMesh()
	if err != nil {
		panic(err)
	}
	var xformMatrix common.TransformationMatrix
	xformMatrix.ParseCommaDelimitedString(xformMatrixString)

	for mIdx, _ := range *meshes {
		for vIdx, _ := range (*meshes)[mIdx].Vertices {
			(*meshes)[mIdx].Vertices[vIdx].Transform(xformMatrix)
		}

		if forceTriangleFlag {
			if flipTriangle {
				(*meshes)[mIdx].FlipTriangleOrder()
			}
		} else {
			if xformMatrix.Det() < 0 {
				(*meshes)[mIdx].FlipTriangleOrder()
			}
		}
	}

	outBuffer, err2 := ii.GetBytes(meshes)

	if err2 != nil {
		panic(err2)
	}

	if len(outBuffer) == 1 {
		gongio.WriteBytesToFile(*(ii.Out), outBuffer[0])
	} else {
		for idx, bytes := range outBuffer {
			ext := filepath.Ext(*(ii.Out))
			regexString := fmt.Sprintf("(%v)?$", ext)
			re := regexp.MustCompile(regexString)
			filename := re.ReplaceAllStringFunc(*(ii.Out), func(ext string) string {
				return fmt.Sprintf("_%d%v", idx, ext)
			})
			gongio.WriteBytesToFile(filename, bytes)
		}
	}

	return nil
}
