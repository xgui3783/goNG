package gii

import (
	"encoding/xml"
	"fmt"
	"gong/common"
	"strconv"
)

func parseString(s string) float32 {
	x, errx := strconv.ParseFloat(s, 32)
	if errx != nil {
		panic(errx)
	}
	return float32(x)
}

func parseStringsToFloat(arrStrings [3]string) [3]float32 {
	x := parseString(arrStrings[0])
	y := parseString(arrStrings[1])
	z := parseString(arrStrings[2])
	return [3]float32{x, y, z}
}

func parseStringInt(s string) uint32 {
	x, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		panic(err)
	}
	return uint32(x)
}

func parseStringsToInt(arrStrings [3]string) [3]uint32 {
	return [3]uint32{
		parseStringInt(arrStrings[0]),
		parseStringInt(arrStrings[1]),
		parseStringInt(arrStrings[2]),
	}
}

func ParseGii(gii []byte) []common.Mesh {
	giiStruct := GIFTI{}
	if err := xml.Unmarshal(gii, &giiStruct); err != nil {
		fmt.Println(err)
		panic(err)
	}
	var vertices DataArray
	var faces DataArray
	for _, dataArray := range giiStruct.DataArray {
		if !vertices.IsAssigned() && dataArray.Intent == NIFTI_INTENT_POINTSET {
			vertices = dataArray
		}
		if !faces.IsAssigned() && dataArray.Intent == NIFTI_INTENT_TRIANGLE {
			faces = dataArray
		}
	}

	if !vertices.IsAssigned() {
		panic("no dataarray with NIFTI_INTENT_POINTSET found")
	}
	if !faces.IsAssigned() {
		panic("no dataarray with NIFTI_INTENT_TRIANGLE found")
	}

	/** process vertices */
	outputVertices := []common.Vertex{}
	for _, triplet := range vertices.getFloatTriplets() {
		outputVertices = append(outputVertices, common.Vertex(triplet))
	}

	/** process faces */
	outputFaces := []common.Face{}
	for _, triplet := range faces.getIntTriplets() {
		outputFaces = append(outputFaces, common.Face(triplet))
	}

	return []common.Mesh{
		common.Mesh{
			Vertices: outputVertices,
			Faces:    outputFaces,
		},
	}
}

func WriteGii(mesh common.Mesh) []byte {
	vertices := mesh.Vertices
	faces := mesh.Faces

	verticesBytes := []byte{}
	facesBytes := []byte{}

	for _, vertex := range vertices {
		if len(verticesBytes) > 0 {
			verticesBytes = append(verticesBytes, []byte(" ")...)
		}

		floatToBeAppended := fmt.Sprintf("%f %f %f", vertex[0], vertex[1], vertex[2])
		verticesBytes = append(verticesBytes, []byte(floatToBeAppended)...)
	}

	for _, face := range faces {
		if len(facesBytes) > 0 {
			facesBytes = append(facesBytes, []byte(" ")...)
		}

		intToBeAppended := fmt.Sprintf("%d %d %d", face[0], face[1], face[2])
		facesBytes = append(facesBytes, []byte(intToBeAppended)...)
	}

	gii := GIFTI{
		MetaData: MetaData{
			MD: []MD{
				MD{
					Name: CData{
						Value: `conversion software`,
					},
					Value: CData{
						Value: `goNG <https://github.com/xgui3783/goNG>`,
					},
				},
			},
		},
		DataArray: []DataArray{
			// vertices
			DataArray{
				ArrayIndexingOrder: RowMajorOrder,
				DataType:           NIFTI_TYPE_FLOAT32,
				Intent:             NIFTI_INTENT_POINTSET,
				Dimensionality:     2,
				Dim0:               len(vertices),
				Dim1:               3,
				Encoding:           ASCII,
				Data: Data{
					Value: string(verticesBytes),
				},
			},
			// faces
			DataArray{
				ArrayIndexingOrder: RowMajorOrder,
				DataType:           NIFTI_TYPE_INT32,
				Intent:             NIFTI_INTENT_TRIANGLE,
				Dimensionality:     2,
				Dim0:               len(faces),
				Dim1:               3,
				Encoding:           ASCII,
				Data: Data{
					Value: string(facesBytes),
				},
			},
		},
	}

	// giiBytes, err := xml.Marshal(gii)
	giiBytes, err := xml.MarshalIndent(gii, "", "  ")
	if err != nil {
		panic(err)
	}
	return append([]byte(GII_HEADER), giiBytes...)
}

const GII_HEADER = `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE GIFTI SYSTEM "http://gifti.projects.nitrc.org/gifti.dtd">
`
