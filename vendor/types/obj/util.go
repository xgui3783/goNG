package obj

import (
	"common"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func ParseObjBytes(input []byte) []common.Mesh {
	strippedComments := stripComments(string(input))
	lines := common.SplitByNewLine(strippedComments)
	mesh := common.Mesh{}
	for _, line := range lines {
		parseObjLine(line, &mesh)
	}
	return []common.Mesh{mesh}
}

func WriteMeshToBytes(mesh common.Mesh) []byte {
	returnybtes := []byte{}

	objName := fmt.Sprintf("o fragment\n")
	returnybtes = append(returnybtes, []byte(objName)...)

	// vertices
	for _, vertex := range mesh.Vertices {
		s := fmt.Sprintf("v %.f %.f %.f\n", vertex[0], vertex[1], vertex[2])
		returnybtes = append(returnybtes, []byte(s)...)
	}

	s := []byte("s off\n")
	returnybtes = append(returnybtes, s...)

	//faces
	for _, face := range mesh.Faces {
		s := fmt.Sprintf("f %d %d %d\n", face[0]+1, face[1]+1, face[2]+1)
		returnybtes = append(returnybtes, []byte(s)...)
	}

	return returnybtes
}

func validateInputTobeParsed(input string) bool {
	if len(input) == 0 {
		return false
	}
	return true
}

func parseObjLine(line string, mesh *common.Mesh) {
	trimmed := common.TrimStartingEndingWhiteSpaces(line)
	if !validateInputTobeParsed(trimmed) {
		return
	}
	firstTwochar := string([]byte(trimmed)[:2])
	switch firstTwochar {
	default:
		panicText := fmt.Sprintf("parsing obj line error. unknown first two characters of line: %v", firstTwochar)
		panic(panicText)
	case "v ":
		stringtoBeParsed := common.TrimStartingEndingWhiteSpaces(string([]byte(trimmed)[2:]))
		floats := common.ParseStringAsFloatsWDelimiter(stringtoBeParsed, " ")
		(*mesh).Vertices = append((*mesh).Vertices, common.Vertex{floats[0], floats[1], floats[2]})
	case "f ":
		stringtoBeParsed := common.TrimStartingEndingWhiteSpaces(string([]byte(trimmed)[2:]))
		objFace := parseFaceLine(stringtoBeParsed)
		(*mesh).Faces = append((*mesh).Faces, common.Face{objFace.Vertex[0], objFace.Vertex[1], objFace.Vertex[2]})
	case "vt":
		fmt.Printf("vt: vertex texture not yet implemented")
	case "vn":
		fmt.Printf("vn: vertex normal not yet implemented")
	case "l ":
		fmt.Printf("line has yet to be implemented")
	case "vp":
		fmt.Printf("parameter space has yet been implemented")
	case "mt":
		fmt.Printf("mtl external file has yet been immplemented")
	case "us":
		fmt.Printf("usemtl has yet been implemented")
	case "o ":
		fallthrough
	case "g ":
		fallthrough
	case "s ":
		s := fmt.Sprintf("NYI: %v", firstTwochar)
		os.Stderr.WriteString(s)
	}

	return
}

func parseFaceLine(line string) ObjFace {
	r := ObjFace{}
	terms := strings.Split(line, ` `)
	for _, term := range terms {
		if !validateInputTobeParsed(term) {
			continue
		}
		digits := strings.Split(term, `/`)
		for idx, digit := range digits {
			switch idx {
			default:
				fmt.Printf("")
			case 0:
				parsedUint, err := strconv.ParseUint(digit, 10, 32)
				if err != nil {
					panic(err)
				}
				r.Vertex = append(r.Vertex, uint32(parsedUint)-1)
			case 1:
				if digit == "" {
					continue
				}
				parsedUint, err := strconv.ParseUint(digit, 10, 32)
				if err != nil {
					panic(err)
				}
				r.VertexTexture = append(r.Vertex, uint32(parsedUint)-1)
			case 2:
				if digit == "" {
					continue
				}
				parsedUint, err := strconv.ParseUint(digit, 10, 32)
				if err != nil {
					panic(err)
				}
				r.VertexNormal = append(r.Vertex, uint32(parsedUint)-1)
			}
		}
	}
	return r
}

func stripComments(input string) string {
	re := regexp.MustCompile(`(?m)\#.*?$`)
	return re.ReplaceAllString(input, "")
}

type ObjFace struct {
	Vertex        []uint32
	VertexNormal  []uint32
	VertexTexture []uint32
}
