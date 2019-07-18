package stl

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
	binary.Write(buf, binary.LittleEndian, math.Float32bits(coord[0]))
	binary.Write(buf, binary.LittleEndian, math.Float32bits(coord[0]))
	return buf.Bytes()
}

func assemble(vertices [3]common.Vertex) []byte {
	buf := new(bytes.Buffer)
	normal := common.GetNormal(vertices)

	binary.Write(buf, binary.LittleEndian, putFloatTriplet(normal))

	for _, v := range vertices {
		binary.Write(buf, binary.LittleEndian, putFloatTriplet(v))
	}

	padding := make([]byte, 0, 2)

	binary.Write(buf, binary.LittleEndian, padding)

	return buf.Bytes()
}

func getStlHeader() []byte {
	header := make([]byte, 0, 80)
	return header
}

func WriteBinaryStlFromMesh(mesh common.Mesh) []byte {

	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, getStlHeader())
	if err != nil {
		panic(err)
	}

	numVertices := len(mesh.Vertices)

	err = binary.Write(buf, binary.LittleEndian, uint32(numVertices))
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

func WriteAsciiStlFromMesh(mesh common.Mesh) []byte {
	buf := new(bytes.Buffer)
	buf.Write([]byte("solid Untitled\n"))

	for _, f := range mesh.Faces {

		v1 := mesh.Vertices[f[0]]
		v2 := mesh.Vertices[f[1]]
		v3 := mesh.Vertices[f[2]]

		normal := common.GetNormal([3]common.Vertex{v1, v2, v3})

		buf.Write([]byte(fmt.Sprintf("facet normal %e %e %e\n", normal[0], normal[1], normal[2])))
		buf.Write([]byte("outer loop\n"))
		buf.Write([]byte(fmt.Sprintf("vertex %e %e %e\n", v1[0], v1[1], v1[2])))
		buf.Write([]byte(fmt.Sprintf("vertex %e %e %e\n", v2[0], v2[1], v2[2])))
		buf.Write([]byte(fmt.Sprintf("vertex %e %e %e\n", v3[0], v3[1], v3[2])))
		buf.Write([]byte(fmt.Sprintf("endloop\n")))
		buf.Write([]byte(fmt.Sprintf("endfacet\n")))
	}

	buf.Write([]byte("endsolid Untitled"))
	return buf.Bytes()
}
