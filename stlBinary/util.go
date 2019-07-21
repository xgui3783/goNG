package stlBinary

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"gong/common"
	"math"
)

func putFloatTriplet(coord common.Vertex) []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, math.Float32bits(coord[0]))
	binary.Write(buf, binary.LittleEndian, math.Float32bits(coord[1]))
	binary.Write(buf, binary.LittleEndian, math.Float32bits(coord[2]))
	return buf.Bytes()
}

func assemble(vertices [3]common.Vertex) []byte {
	buf := new(bytes.Buffer)
	normal := common.GetNormal(vertices)

	binary.Write(buf, binary.LittleEndian, putFloatTriplet(normal))

	for _, v := range vertices {
		binary.Write(buf, binary.LittleEndian, putFloatTriplet(v))
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

func WriteBinaryStlFromMesh(mesh common.Mesh, metadata common.MeshMetadata) []byte {

	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, getStlHeader())
	if err != nil {
		panic(err)
	}

	// fmt.Printf("header written, buffer size %v\n", strconv.Itoa(len(buf.Bytes())))

	numVertices := len(mesh.Vertices)
	numFaces := len(mesh.Faces)

	err = binary.Write(buf, binary.LittleEndian, uint32(numFaces))
	if err != nil {
		panic(err)
	}
	fmt.Printf("num vertex written\n")

	fmt.Printf("vertices length %v\n", numVertices)
	for _, f := range mesh.Faces {
		vertices := [3]common.Vertex{}
		for idx, vIndex := range f {
			fmt.Printf("index: %v, vindex: %v\n", idx, vIndex)
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
