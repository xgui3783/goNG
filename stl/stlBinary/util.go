package stlBinary

import (
	"bytes"
	"encoding/binary"
	"gong/common"
	"gong/stl/stlCommon"
)

type binaryStlQuintStruct struct {
	Normal    [3]float32
	Vertex1   [3]float32
	Vertex2   [3]float32
	Vertex3   [3]float32
	Attribute uint16
}

func parseBytesAsFloat32(input []byte) (strPtr *binaryStlQuintStruct) {
	strPtr = &(binaryStlQuintStruct{})

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
	meshData := input[84:]
	numTriangles := binary.LittleEndian.Uint32(input[80:84])
	
	mapVertIdx := make(map[float32]map[float32]map[float32]uint32, 0)
	vertexIndex := uint32(0)

	makeDeepMap := func(x float32, y float32, z float32) {
		_, okx := mapVertIdx[x]
		if okx == false {
			mapVertIdx[x] = make(map[float32]map[float32]uint32)
		}
		_, oky := mapVertIdx[x][y]
		if oky == false {
			mapVertIdx[x][y] = make(map[float32]uint32)
		}
	}
	vertices := make([]common.Vertex, 0, vertexIndex)
	getVertexIndex := func(vertex common.Vertex) uint32 {
		// TODO relying on float32 equality, which is always a shaky thing to rely on
		idx, ok := mapVertIdx[vertex[0]][vertex[1]][vertex[2]]
		if ok == false {
			makeDeepMap(vertex[0], vertex[1], vertex[2])

			defer func() {
				vertexIndex++
			}()
			mapVertIdx[vertex[0]][vertex[1]][vertex[2]] = vertexIndex
			vertices = append(vertices, vertex)
			return vertexIndex
		} else {
			return idx
		}
	}
	faces := []common.Face{}

	for i := uint32(0); i < numTriangles; i++ {
		startIdx := 50 * i
		endIdx := startIdx + 50
		strPtr := parseBytesAsFloat32(meshData[startIdx:endIdx])
		parsedVertices := [3]common.Vertex{
			(*strPtr).Vertex1,
			(*strPtr).Vertex2,
			(*strPtr).Vertex3,
		}
		faceVertexIndex := [3]uint32{
			getVertexIndex(parsedVertices[0]),
			getVertexIndex(parsedVertices[1]),
			getVertexIndex(parsedVertices[2]),
		}

		faces = append(faces, common.Face{faceVertexIndex[0], faceVertexIndex[1], faceVertexIndex[2]})
	}
	return &common.Mesh{
		Faces:    faces,
		Vertices: vertices,
	}
}
