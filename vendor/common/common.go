package common

import (
	"detProtocol"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func Vec3(slice []float32) [3]float32 {
	if len(slice) != 3 {
		panicText := fmt.Sprintf("Vec3 input does not have length 3!")
		panic(panicText)
	}
	return [3]float32{slice[0], slice[1], slice[2]}
}

type Vertex [3]float32
type Face [3]uint32
type Normal [3]float32

func (v *Vertex) Transform(m TransformationMatrix) {
	var prev Vertex
	for idx := 0; idx <= 2; idx++ {
		prev[idx] = v[idx]
	}
	for idx := 0; idx <= 2; idx++ {
		v[idx] = prev[0]*float32(m[idx][0]) + prev[1]*float32(m[idx][1]) + prev[2]*float32(m[idx][2]) + 1*float32(m[idx][3])
	}
}

type Mesh struct {
	Vertices []Vertex `json:"vertices"`
	Faces    []Face   `json:"faces"`
	// VertexNormals []Normal `json:"vertexNormals"`
}

func (meshPtr *Mesh) AddVertex(newVertex [3]float32) uint32 {
	meshPtr.Vertices = append(meshPtr.Vertices, newVertex)
	return uint32(len(meshPtr.Vertices) - 1)
}

func (meshPtr *Mesh) AddFace(vtxIndices [3]uint32) error {
	size := uint32(len(meshPtr.Vertices))
	if size < vtxIndices[0] && size < vtxIndices[1] && size < vtxIndices[2] {
		errorText := fmt.Sprintf("Index out of bound! %v", vtxIndices)
		return errors.New(errorText)
	}
	meshPtr.Faces = append(meshPtr.Faces, vtxIndices)
	return nil
}

func (meshPtr *Mesh) AppendMesh(anotherMeshPtr *Mesh) error {
	prevVerticies := len(meshPtr.Vertices)
	meshPtr.Vertices = append(meshPtr.Vertices, anotherMeshPtr.Vertices...)

	tmpFacesCollection := []Face{}
	for _, face := range anotherMeshPtr.Faces {
		tmpFacesCollection = append(tmpFacesCollection, Face{
			face[0] + uint32(prevVerticies),
			face[1] + uint32(prevVerticies),
			face[2] + uint32(prevVerticies),
		})
	}
	meshPtr.Faces = append(meshPtr.Faces, tmpFacesCollection...)
	return nil
}

func GetVertexMinMax(meshPtr *Mesh) (error, [2]*Vertex) {
	if len(meshPtr.Vertices) == 0 {
		return errors.New("Mesh has 0 vertex"), [2]*Vertex{}
	}

	returnVal := [2]*Vertex{
		&Vertex{1e32, 1e32, 1e32},
		&Vertex{-1e32, -1e32, -1e32},
	}
	for _, vertex := range meshPtr.Vertices {
		if vertex[0] < returnVal[0][0] {
			returnVal[0][0] = vertex[0]
		}
		if vertex[1] < returnVal[0][1] {
			returnVal[0][1] = vertex[1]
		}
		if vertex[2] < returnVal[0][2] {
			returnVal[0][2] = vertex[2]
		}

		if vertex[0] > returnVal[1][0] {
			returnVal[1][0] = vertex[0]
		}
		if vertex[1] > returnVal[1][1] {
			returnVal[1][1] = vertex[1]
		}
		if vertex[2] > returnVal[1][2] {
			returnVal[1][2] = vertex[2]
		}
	}

	return nil, returnVal
}

type TransformationMatrix [4][4]float64

func (m *TransformationMatrix) ParseCommaDelimitedString(input string) {
	vals := strings.Split(input, ",")
	if len(vals) != 12 && len(vals) != 16 {
		panicText := fmt.Sprintf("ParseCommaDelimitedString needs to have 12 or 16 elements, but instead has %v elements\n", len(vals))
		panic(panicText)
	}
	for idx, val := range vals {
		row := int(idx / 4)
		col := int(idx % 4)
		var err error
		m[row][col], err = strconv.ParseFloat(val, 64)
		if err != nil {
			panicText := fmt.Sprintf("Parse float error! %v", err)
			panic(panicText)
		}
	}
	if len(vals) == 12 {
		m[3] = [4]float64{0.0, 0.0, 0.0, 1.0}
	}
}

// adopted from http://glmatrix.net/docs/mat4.js.html#line341
func (m *TransformationMatrix) Det() float64 {
	b00 := m[0][0]*m[1][1] - m[0][1]*m[1][0]
	b01 := m[0][0]*m[1][2] - m[0][2]*m[1][0]
	b02 := m[0][0]*m[1][3] - m[0][3]*m[1][0]
	b03 := m[0][1]*m[1][2] - m[0][2]*m[1][1]
	b04 := m[0][1]*m[1][3] - m[0][3]*m[1][1]
	b05 := m[0][2]*m[1][3] - m[0][3]*m[1][2]
	b06 := m[2][0]*m[3][1] - m[2][1]*m[3][0]
	b07 := m[2][0]*m[3][2] - m[2][2]*m[3][0]
	b08 := m[2][0]*m[3][3] - m[2][3]*m[3][0]
	b09 := m[2][1]*m[3][2] - m[2][2]*m[3][1]
	b10 := m[2][1]*m[3][3] - m[2][3]*m[3][1]
	b11 := m[2][2]*m[3][3] - m[2][3]*m[3][2]

	return b00*b11 - b01*b10 + b02*b09 + b03*b08 - b04*b07 + b05*b06
}

func (mesh *Mesh) FlipTriangleOrder() {
	for fIdx, face := range mesh.Faces {
		mesh.Faces[fIdx][0] = face[2]
		mesh.Faces[fIdx][2] = face[0]
	}
}

type MeshMetadata struct {
	Index int
}

func GetPlaneNormalPtr(plane [3][3]float32) *[3]float32 {
	planeV1 := Sub(plane[0], plane[1])
	planeV2 := Sub(plane[0], plane[2])
	normal := Cross(planeV1, planeV2)
	return &normal
}

func GetVec3PlanePolarity(v [3]float32, planeNormalPtr *[3]float32, planePtPtr *[3]float32) float32 {
	return Dot(Sub(v, *planePtPtr), *planeNormalPtr)
}

func FindIntersect(line [2][3]float32, planeNormalPtr *[3]float32, planePtPtr *[3]float32, internalFlag bool) (error, [3]float32) {

	if GetMod(*planeNormalPtr) == 0 {
		return errors.New("plane normal has mod of zero."), [3]float32{}
	}

	/*
		solving intersect between plane (normal + point) and line segment (2x point)

		line segment

		v1 = x1, y1, z1
		v2 = x2, y2, z2

		parametric form:
		(when t = 0, at v1, t = 1 at v2)
		x = x1 + (x2 - x1) t
		y = y1 + (y2 - y1) t
		z = z1 + (z2 - z1) t


		plane:
		v3 = x3, y3, z3
		v4 = x4, y4, z4
		v5 = x5, y5, z5

		normal to plane:
		vn = xn, yn, zn

		(x - x3, y - y3, z - z3) . norm = 0

		xn . x - x3 . xn +
		yn . y - y3 . yn +
		zn . z - z3 . zn = 0

		xn . [x1 + (x2 - x1) t ] - x3 . xn +
		yn . [y1 + (y2 - y1) t ] - y3 . yn +
		zn . [z1 + (z2 - z1) t ] - z3 . zn = 0


		xn . x1 + xn . (x2 - x1) t  - x3 . xn +
		yn . y1 + yn . (y2 - y1) t  - y3 . yn +
		zn . z1 + zn . (z2 - z1) t  - z3 . zn = 0


		[xn (x2 - x1) +
		yn (y2 - y1) +
		zn (z2 - z1) ] t =
		( x3 . xn +
			y3 . yn +
			z3 . zn ) -
		( xn . x1 +
			yn . y1 +
			zn . z1 )

	*/

	coef := Dot(Sub(line[1], line[0]), *planeNormalPtr)
	rhs := Dot(*planePtPtr, *planeNormalPtr) - Dot(line[0], *planeNormalPtr)
	solveT := rhs / coef
	if internalFlag && (solveT < 0 || solveT > 1) {
		errorText := fmt.Sprintf("intersection outside line segment with t == %v\n", solveT)
		return errors.New(errorText), [3]float32{}
	}
	return nil, Add(line[0], MulScalar(Sub(line[1], line[0]), solveT))
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

func Add(v1 [3]float32, v2 [3]float32) (output [3]float32) {
	for idx, v := range v1 {
		output[idx] = v + v2[idx]
	}
	return
}

func Sub(v1 [3]float32, v2 [3]float32) (output [3]float32) {
	for idx, v := range v1 {
		output[idx] = v - v2[idx]
	}
	return
}

func Dot(v1 [3]float32, v2 [3]float32) (output float32) {
	output = 0
	output += v1[0] * v2[0]
	output += v1[1] * v2[1]
	output += v1[2] * v2[2]
	return output
}

func Cross(v1 [3]float32, v2 [3]float32) (output [3]float32) {
	output[0] = v1[1]*v2[2] - v1[2]*v2[1]
	output[1] = v1[2]*v2[0] - v1[0]*v2[2]
	output[2] = v1[0]*v2[1] - v1[1]*v2[0]
	return
}

func GetNormal(vertices [3]Vertex) [3]float32 {
	v1 := Sub(vertices[0], vertices[1])
	v2 := Sub(vertices[0], vertices[2])
	return Normalize(Cross(v1, v2))
}

func MulScalar(v [3]float32, s float32) [3]float32 {
	return [3]float32{
		v[0] * s,
		v[1] * s,
		v[2] * s,
	}
}

func GetMod(vertex Vertex) float32 {
	return float32(math.Sqrt(float64(vertex[0]*vertex[0] + vertex[1]*vertex[1] + vertex[2]*vertex[2])))
}

func Normalize(vertex Vertex) Vertex {
	mod := GetMod(vertex)
	if mod == 0 {
		panic("Cannot normalize zero vector!")
	}
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

func SplitByNewLine(input string) []string {
	re := regexp.MustCompile(`\n`)
	return re.Split(input, -1)
}

func TrimStartingEndingWhiteSpaces(inputBytes string) string {
	reTrimBeginning := regexp.MustCompile(`(^[\s]+|[\s]+$)`)
	return reTrimBeginning.ReplaceAllString(inputBytes, "")
}

func SplitStringByWhiteSpaceNL(inputBytes string) []string {
	trimmedString := TrimStartingEndingWhiteSpaces(inputBytes)
	re := regexp.MustCompile(`(?m)[\s]+`)
	return re.Split(trimmedString, -1)
}

func ParseStringAsFloatsWDelimiter(input string, delimiter string) []float32 {
	splitString := strings.Split(input, delimiter)
	returnFloats := []float32{}
	for _, s := range splitString {
		f, err := strconv.ParseFloat(s, 32)
		if err != nil {
			panic(err)
		}
		returnFloats = append(returnFloats, float32(f))
	}

	return returnFloats
}

func ParseStringAsIntsWDelimiter(input string, delimiter string) []uint32 {
	splitString := strings.Split(input, delimiter)
	returnInts := []uint32{}
	for _, s := range splitString {
		pInt, err := strconv.ParseInt(s, 10, 32)
		if err != nil {
			panic(err)
		}
		returnInts = append(returnInts, uint32(pInt))
	}

	return returnInts
}

const (
	/**
	* some file formats, such as STL, does not provide vertices as an array
	* in such a case, FLOAT_TOLERANCE will be used to determine the spatial resolution between two points
	 */
	FLOAT_TOLERANCE = 1.0e-6
	HEADER          = `generated by goNG <https://github.com/xgui3783/goNG>`
)

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
	panic("does not support HTTP protocol")

	// resp, err := http.Get(rootPath)
	// if err != nil {
	// 	panic(err)
	// }
	// defer resp.Body.Close()
	// body, moreErr := ioutil.ReadAll(resp.Body)
	// if moreErr != nil {
	// 	panic(moreErr)
	// }
	// return body
}
