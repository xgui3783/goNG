package ngPrecomputed

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"gong/common"
	"io/ioutil"
	"math"
	"path/filepath"
	"regexp"
	"strconv"
)

func ScanLocalDir(dir string) []int {
	fileInfos, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}
	labelIndicies := make([]int, 0)
	for _, fileInfo := range fileInfos {
		if fileInfo.IsDir() {
			continue
		}
		r := regexp.MustCompile(`^([0-9]*?)(:0)?$`)
		matches := r.FindStringSubmatch(fileInfo.Name())

		if len(matches) == 0 || matches[1] == "" {
			continue
		}
		idx, err := strconv.Atoi(matches[1])
		if err != nil {
			panic(err)
		}
		labelIndicies = append(labelIndicies, idx)
	}
	return labelIndicies
}

func parseMeshInfoString(meshInfoString []byte) []string {
	meshInfo := NgMeshInfo{}
	if err := json.Unmarshal(meshInfoString, &meshInfo); err != nil {
		panic(err)
	}
	return meshInfo.Fragments
}

func GetHttpFragments(rootPath string, index int) []string {
	strKey := strconv.Itoa(index)
	url := rootPath + strKey + ":0"
	body := common.GetHTTPResource(url)
	return parseMeshInfoString(body)
}

func GetLocalFragments(rootPath string, index int) []string {
	filename := strconv.Itoa(index)
	fullFilename := ""
	if common.CheckFileExists(filepath.Join(rootPath, filename)) {
		fullFilename = filepath.Join(rootPath, filename)
	} else if common.CheckFileExists(filepath.Join(rootPath, filename+":0")) {
		fullFilename = filepath.Join(rootPath, filename+":0")
	} else {
		panic("neither " + filename + " nor " + filename + ":0 exists")
	}

	jsonBody, err := ioutil.ReadFile(fullFilename)
	if err != nil {
		panic(err)
	}

	return parseMeshInfoString(jsonBody)
}

func WriteFragmentFromMesh(mesh common.Mesh) []byte {
	output := make([]byte, 0)
	numVertex := len(mesh.Vertices)
	binary.LittleEndian.PutUint32(output[0:4], uint32(numVertex))

	for vIndex, v := range mesh.Vertices {
		for cIndex, c := range v {
			offset := vIndex*12 + cIndex*4 + 4
			bits := math.Float32bits(c)
			binary.LittleEndian.PutUint32(output[offset:offset+4], bits)
		}
	}

	initFaceOffset := 4 + numVertex*4*3
	for vIndex, v := range mesh.Faces {
		for cIndex, c := range v {
			offset := initFaceOffset + vIndex*12 + cIndex*4
			binary.LittleEndian.PutUint32(output[offset:offset+4], c)
		}
	}

	return output
}

func ParseFragmentBuffer(buffer []byte) common.Mesh {
	bufferLength := len(buffer)

	fmt.Printf("buffer length %v\n", strconv.Itoa(bufferLength))
	numVertex := binary.LittleEndian.Uint32(buffer[0:4])
	fmt.Printf("num vertex %v\n", strconv.Itoa(int(numVertex)))
	numTriangles := (bufferLength - 4 - (int(numVertex) * 12)) / 12
	fmt.Printf("num of triangles %v\n", strconv.Itoa(int(numTriangles)))
	vertexBufferOffset := 4
	vertices := []common.Vertex{}
	vertexIndex := 0
	for uint32(vertexIndex) < numVertex {

		xBits := binary.LittleEndian.Uint32(buffer[vertexBufferOffset+0+vertexIndex*12 : vertexBufferOffset+4+vertexIndex*12])
		x := math.Float32frombits(xBits)
		yBits := binary.LittleEndian.Uint32(buffer[vertexBufferOffset+4+vertexIndex*12 : vertexBufferOffset+8+vertexIndex*12])
		y := math.Float32frombits(yBits)
		zBits := binary.LittleEndian.Uint32(buffer[vertexBufferOffset+8+vertexIndex*12 : vertexBufferOffset+12+vertexIndex*12])
		z := math.Float32frombits(zBits)

		vertices = append(vertices, common.Vertex{x, y, z})

		vertexIndex++
	}

	faceBufferOffset := 4 + int(numVertex)*12
	faces := []common.Face{}
	faceIndex := 0
	for faceIndex < numTriangles {
		v1 := binary.LittleEndian.Uint32(buffer[faceBufferOffset+0+faceIndex*12 : faceBufferOffset+4+faceIndex*12])
		v2 := binary.LittleEndian.Uint32(buffer[faceBufferOffset+4+faceIndex*12 : faceBufferOffset+8+faceIndex*12])
		v3 := binary.LittleEndian.Uint32(buffer[faceBufferOffset+8+faceIndex*12 : faceBufferOffset+12+faceIndex*12])

		faces = append(faces, common.Face{v1, v2, v3})
		faceIndex++
	}

	return common.Mesh{Vertices: vertices, Faces: faces}
}

func IsInfoFile(inputFilepath string) bool {
	return regexp.MustCompile(`\/info$`).Match([]byte(inputFilepath))
}

type NgMeshInfo struct {
	Fragments []string `json:"fragments"`
}
