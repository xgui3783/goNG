package common

import (
	"gong/detProtocol"
	"io/ioutil"
	"math"
	"net/http"
	"os"
)

type Vertex [3]float32
type Face [3]uint32

type Mesh struct {
	Vertices []Vertex `json:"vertices"`
	Faces    []Face   `json:"faces"`
}

type MeshMetadata struct {
	Index int
}

func FindMin(nums []float32) float32 {
	if len(nums) == 0 {
		panic("len(nums) cannot be zero")
	}
	min := nums[0]
	for _, v := range nums {
		if v < min {
			min = v
		}
	}

	return min
}

func findMax(nums []float32) float32 {
	if len(nums) == 0 {
		panic("len(nums) cannot be zero")
	}
	max := nums[0]
	for _, v := range nums {
		if v > max {
			max = v
		}
	}

	return max
}

func Sub(v1 [3]float32, v2 [3]float32) (output [3]float32) {
	for idx, v := range v1 {
		output[idx] = v - v2[idx]
	}
	return
}

func Cross(v1 [3]float32, v2 [3]float32) (output [3]float32) {
	output[0] = v1[1]*v2[2] - v1[2]*v2[1]
	output[1] = v1[0]*v2[2] - v1[2]*v2[0]
	output[2] = v1[0]*v2[1] - v1[0]*v2[1]
	return
}

func GetNormal(vertices [3]Vertex) [3]float32 {
	v1 := Sub(vertices[0], vertices[1])
	v2 := Sub(vertices[0], vertices[2])
	return Normalize(Cross(v1, v2))
}

func GetMod(vertex Vertex) float32 {
	return float32(math.Sqrt(float64(vertex[0]*vertex[0] + vertex[1]*vertex[1] + vertex[2]*vertex[2])))
}

func Normalize(vertex Vertex) Vertex {
	mod := GetMod(vertex)
	return Vertex{vertex[0] / mod, vertex[1] / mod, vertex[2] / mod}
}

func CheckFileExists(inputFilepath string) bool {

	_, err := os.Stat(inputFilepath)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	panic(err)
}

func GetResource(rootPath string) []byte {
	switch detProtocol.InferProtocolFromFilename(rootPath) {
	case detProtocol.Local:
		filebytes, err := ioutil.ReadFile(rootPath)
		if err != nil {
			panic(err)
		}
		return filebytes
	case detProtocol.HTTP:
		return GetHTTPResource(rootPath)
	default:
		panic("Get resource error")
	}
}

func GetHTTPResource(rootPath string) []byte {
	resp, err := http.Get(rootPath)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, moreErr := ioutil.ReadAll(resp.Body)
	if moreErr != nil {
		panic(moreErr)
	}
	return body
}

const (
	/**
	* some file formats, such as STL, does not provide vertices as an array
	* in such a case, FLOAT_TOLERANCE will be used to determine the spatial resolution between two points
	 */
	FLOAT_TOLERANCE = 1.0e-6
)
