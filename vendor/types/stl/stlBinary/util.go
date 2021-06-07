package stlBinary

import (
	"bytes"
	"common"
	"encoding/binary"
	"types/stl/stlCommon"
)

type binaryStlQuintStruct struct {
	Normal    [3]float32
	Vertex1   [3]float32
	Vertex2   [3]float32
	Vertex3   [3]float32
	Attribute uint16
}

func parseBytesAsFloat32(input []byte) (strPtr *binaryStlQuintStruct) {

	r := bytes.NewReader(input)

	if err := binary.Read(r, binary.LittleEndian, strPtr); err != nil {
		panic("parseBytesAsFloat32 error\n")
	}

	return
}

func assemble(vertices [3]common.Vertex) []byte {
	buf := new(bytes.Buffer)
	normal := common.GetNormal(vertices)

	putFloatTriplet := func(coord common.Vertex) {
		stlCommon.AppendFloat32ToBuffer(coord[0], buf)
		stlCommon.AppendFloat32ToBuffer(coord[1], buf)
		stlCommon.AppendFloat32ToBuffer(coord[2], buf)
	}

	putFloatTriplet(normal)

	for _, v := range vertices {
		putFloatTriplet(v)
	}

	padding := make([]byte, 2, 2)

	binary.Write(buf, binary.LittleEndian, padding)

	return buf.Bytes()
}

func getStlHeader() []byte {
	maxLength := 80
	stringLiteral := "Converted with goNG. Author Xiao Gui<panda@pandamakes.com.au>. MIT licensed"

	var header []byte

	if maxLength < len(stringLiteral) {
		header = []byte(stringLiteral)[:maxLength]
	} else {
		diff := maxLength - len(stringLiteral)
		header = append([]byte(stringLiteral), make([]byte, diff, diff)...)
	}
	return header
}

func WriteBinaryStlFromMesh(mesh common.Mesh) []byte {

	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, getStlHeader())
	if err != nil {
		panic(err)
	}

	numFaces := len(mesh.Faces)

	err = binary.Write(buf, binary.LittleEndian, uint32(numFaces))
	if err != nil {
		panic(err)
	}

	for _, f := range mesh.Faces {
		vertices := [3]common.Vertex{}
		for idx, vIndex := range f {
			vertices[idx] = mesh.Vertices[vIndex]
		}
		assembled := assemble(vertices)
		err = binary.Write(buf, binary.LittleEndian, assembled)
		if err != nil {
			panic(err)
		}
	}

	return buf.Bytes()
}

func parseBinaryStl(input []byte) *(common.Mesh) {
	mesh := common.Mesh{}
	meshData := input[84:]
	numTriangles := binary.LittleEndian.Uint32(input[80:84])
	stlAppendFace := stlCommon.PrepareStlAppendFace(&mesh)
	for i := uint32(0); i < numTriangles; i++ {
		startIdx := 14 * i
		endIdx := startIdx + 14
		strPtr := parseBytesAsFloat32(meshData[startIdx:endIdx])
		stlAppendFace([3]common.Vertex{
			(*strPtr).Vertex1,
			(*strPtr).Vertex2,
			(*strPtr).Vertex3,
		})
	}
	return &mesh
}
