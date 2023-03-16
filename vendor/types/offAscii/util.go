package offAscii

import (
	"errors"
	"fmt"
	"common"
	"regexp"
	"strconv"
)

const (
	OFF   = iota
	COFF  = iota
	NOFF  = iota
	CNOFF = iota

	UNDETERMINED = iota
)

// scan first 128 lines for header info. If does not find, throw
const HEADER_LIMIT = 128

func WriteMeshToBytes(mesh common.Mesh) (output []byte) {
	summaryString := fmt.Sprintf("%d %d 1\n", len(mesh.Vertices), len(mesh.Faces))
	output = append(output, []byte(summaryString)...)
	for _, ver := range mesh.Vertices {
		vertexString := fmt.Sprintf("%g %g %g\n", ver[0], ver[1], ver[2])
		output = append(output, []byte(vertexString)...)
	}

	for _, face := range mesh.Faces {
		faceString := fmt.Sprintf("%d %d %d", face[0], face[1], face[2])
		output = append(output, []byte(faceString)...)
	}

	return
}

func ParseOffAscii(buffer []byte) (meshes []common.Mesh) {

	lines := common.SplitByNewLine(string(buffer))
	numVertex, numFaces, offsetToData, _, err := parseHeader(&lines)
	if err != nil {
		panic(err)
	}
	meshes = []common.Mesh{
		common.Mesh{
			Vertices: []common.Vertex{},
			Faces:    []common.Face{},
		},
	}
	vp := &meshes[0].Vertices
	fp := &meshes[0].Faces

	for int64(len(*vp)) < numVertex {
		if ver, err := parseVertex(lines[offsetToData]); err == nil {
			*vp = append(*vp, ver)
		}
		offsetToData++
	}

	for int64(len(*fp)) < numFaces {
		if face, err := parseFace(lines[offsetToData]); err == nil {
			*fp = append(*fp, face)
		}
		offsetToData++
	}
	return
}

func parseHeader(lines *[]string) (numVertex int64, numFaces int64, offsetToData int, offType int, returnError error) {
	regexpString := fmt.Sprintf("^.*?(C)?(N)?OFF")
	re := regexp.MustCompile(regexpString)
	for index, line := range *lines {

		// comment block... ignore line
		if string(line[0]) == "#" {
			continue
		}

		if _, err := strconv.ParseInt(string(line[0]), 10, 8); err == nil {
			nums := common.SplitStringByWhiteSpaceNL(line)
			if len(nums) < 2 {
				errorString := fmt.Sprintf("Expected vertex face tally to be >=2, but got %d. Raw input: %v\n", len(nums), line)
				returnError = errors.New(errorString)
				return
			} else {
				if numVertex, returnError = strconv.ParseInt(nums[0], 10, 64); returnError != nil {
					return
				}

				if numFaces, returnError = strconv.ParseInt(nums[1], 10, 64); returnError != nil {
					return
				}
				offsetToData = index + 1
				return
			}
		}

		headerRe := re.FindString(line)
		if len(headerRe) > 0 {
			switch headerRe {
			case "OFF":
				offType = OFF
			case "NOFF":
				offType = NOFF
			case "COFF":
				offType = COFF
			case "CNOFF":
				offType = CNOFF
			default:
				offType = UNDETERMINED
			}
		}

		if index >= HEADER_LIMIT {
			errString := fmt.Sprintf("Header not found after %d lines... end.", HEADER_LIMIT)
			returnError = errors.New(errString)
			return
		}
	}

	errString := fmt.Sprintf("header wasn't found. Is this an ascii off file?")
	returnError = errors.New(errString)
	return
}

func parseFace(line string) (face common.Face, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("parseFace error")
		}
	}()

	ints := common.ParseStringAsIntsWDelimiter(line, " ")
	if len(ints) < 4 {
		panic("parseFace not getting 3 results in return")
	}
	if ints[0] != 3 {
		panic("first index of face does not equal 3")
	}
	face = common.Face{ints[1], ints[2], ints[3]}
	return
}

func parseVertex(line string) (ver common.Vertex, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("parseVertex error")
		}
	}()
	floats := common.ParseStringAsFloatsWDelimiter(line, " ")
	if len(floats) < 3 {
		panic("parseVertex not getting 3 results in return")
	}
	ver = common.Vertex{floats[0], floats[1], floats[2]}
	return
}
