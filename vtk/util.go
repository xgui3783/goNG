package vtk

import (
	"fmt"
	"gong/common"
	"regexp"
	"strconv"
	"strings"
)

func parseVtk(inputBytes []byte) common.Mesh {
	input := string(inputBytes)
	lines := strings.Split(input, "\n")

	re := regexp.MustCompile("^#\\ vtk\\ DataFile\\ Version\\ [0-9]\\.[0-9]$")
	if !re.MatchString(lines[0]) {
		panic("vtk header incorrect")
	}

	if len(lines[1]) > 256 {
		panic("header size too large")
	}

	outputMesh := common.Mesh{
		Vertices: []common.Vertex{},
		Faces:    []common.Face{},
	}

	// encoding := lines[3]
	readingFlag := false
	readingPoints := false
	numPoints := 0
	readingFaces := false
	numFaces := 0

	reOther := regexp.MustCompile(`POINT_DATA`)
	rePoints := regexp.MustCompile(`^POINTS\ ([0-9]*?)\ `)
	reFaces := regexp.MustCompile(`^POLYGONS\ ([0-9]*?)\ ([0-9]*?)$`)

	reDataset := regexp.MustCompile(`^DATASET`)
	reDatasetPolyData := regexp.MustCompile(`^DATASET\ POLYDATA`)

	for _, line := range lines[3:] {

		// parse begiinning of dataset block
		if reDataset.MatchString(line) {
			if reDatasetPolyData.MatchString(line) {
				readingFlag = true
			} else {
				readingFlag = false
				readingPoints = false
				readingFaces = false
			}

			continue
		}

		if !readingFlag {
			continue
		}

		if reOther.MatchString(line) {
			readingPoints = false
			readingFaces = false
			continue
		}

		if rePoints.MatchString(line) {
			readingPoints = true
			readingFaces = false

			numPointsString := rePoints.FindStringSubmatch(line)[1]
			np, err := strconv.Atoi(numPointsString)
			if err != nil {
				panic(err)
			}
			numPoints = np
			continue
		}

		if reFaces.MatchString(line) {
			readingFaces = true
			readingPoints = false

			numFacesString := reFaces.FindStringSubmatch(line)[1]
			numFacesTotalElementsString := reFaces.FindStringSubmatch(line)[2]
			nf, err := strconv.Atoi(numFacesString)
			if err != nil {
				panic(err)
			}
			numFaces = nf
			numFacesTotal, err := strconv.Atoi(numFacesTotalElementsString)
			if err != nil {
				panic(err)
			}
			if numFaces*4 != numFacesTotal {
				panicText := fmt.Sprintf("numfaces *4 != numFacesTotal. %v * 4 != %v", numFaces, numFacesTotal)
				panic(panicText)
			}
			continue
		}

		if readingPoints {
			splitStr := strings.Split(line, ` `)

			p1, err := strconv.ParseFloat(splitStr[0], 32)
			if err != nil {
				panic(err)
			}
			p2, err := strconv.ParseFloat(splitStr[1], 32)
			if err != nil {
				panic(err)
			}
			p3, err := strconv.ParseFloat(splitStr[2], 32)
			if err != nil {
				panic(err)
			}
			vertex := common.Vertex{float32(p1), float32(p2), float32(p3)}
			outputMesh.Vertices = append(outputMesh.Vertices, vertex)
			continue
		}

		if readingFaces {
			splitStr := strings.Split(line, ` `)

			if splitStr[0] != "3" {
				panictext := fmt.Sprintf("polygon attribute with first entry other than 3 is not yet supported. The vertex entry is %v", splitStr[0])
				panic(panictext)
			}

			p1, err := strconv.ParseUint(splitStr[1], 10, 32)
			if err != nil {
				panic(err)
			}
			p2, err := strconv.ParseUint(splitStr[2], 10, 32)
			if err != nil {
				panic(err)
			}
			p3, err := strconv.ParseUint(splitStr[3], 10, 32)
			if err != nil {
				panic(err)
			}
			face := common.Face{uint32(p1), uint32(p2), uint32(p3)}
			outputMesh.Faces = append(outputMesh.Faces, face)
			continue
		}
	}

	if numPoints != len(outputMesh.Vertices) {
		fmt.Printf("numPoints: %d outputMesh no vertex: %d\n", numPoints, len(outputMesh.Vertices))
	}

	if numFaces != len(outputMesh.Faces) {
		fmt.Printf("numFaces: %d outpushMesh no faces %d\n", numFaces, len(outputMesh.Faces))
	}

	fmt.Printf("num vertex %d num faces%d\n", len(outputMesh.Vertices), len(outputMesh.Faces))

	return outputMesh
}

func writeVtk(mesh common.Mesh) []byte {
	returnBytes := []byte("# vtk DataFile Version 2.0\n")

	headerString := fmt.Sprintf("%v\n", common.HEADER)
	returnBytes = append(returnBytes, []byte(headerString)...)

	block := fmt.Sprintf(`ASCII
DATASET POLYDATA
POINTS %d double
`, len(mesh.Vertices))
	returnBytes = append(returnBytes, []byte(block)...)

	for _, vertex := range mesh.Vertices {
		vertexString := fmt.Sprintf("%.2f %.2f %.2f\n", vertex[0], vertex[1], vertex[2])
		returnBytes = append(returnBytes, []byte(vertexString)...)
	}

	polygonBlock := fmt.Sprintf("POLYGONS %d %d\n", len(mesh.Faces), len(mesh.Faces)*4)
	returnBytes = append(returnBytes, []byte(polygonBlock)...)

	for _, face := range mesh.Faces {
		faceString := fmt.Sprintf("3 %d %d %d\n", face[0], face[1], face[2])
		returnBytes = append(returnBytes, []byte(faceString)...)
	}
	return returnBytes
}

const (
	ASCII  = "ASCII"
	BINARY = "BINARY"
)
