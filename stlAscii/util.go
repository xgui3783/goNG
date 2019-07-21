package stlAscii

import (
	"bytes"
	"fmt"
	"gong/common"
	"regexp"
	"strconv"
)

func parseStlAsciiLine(line []byte) common.Vertex {

	snRegexpString := `[0-9]\.[0-9]*[eE]-?[0-9][0-9]*`
	vertReString := fmt.Sprintf(`(%v)\W*?(%v)\W*?(%v)`, snRegexpString, snRegexpString, snRegexpString)
	re3 := regexp.MustCompile(vertReString)
	match := re3.FindSubmatch(line)
	x, err := strconv.ParseFloat(string(match[1]), 32)
	if err != nil {
		panic(err)
	}
	y, err := strconv.ParseFloat(string(match[2]), 32)
	if err != nil {
		panic(err)
	}
	z, err := strconv.ParseFloat(string(match[3]), 32)
	if err != nil {
		panic(err)
	}
	return common.Vertex{float32(x), float32(y), float32(z)}
}
func ParseAsciiStl(buffer []byte) common.Mesh {

	re := regexp.MustCompile(`(?s)facet(.*?)endfacet`)
	triangles := re.FindAll(buffer, -1)

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
	for _, triangle := range triangles {
		re2 := regexp.MustCompile(`(?m)vertex\W(.*?)$`)
		vertices := re2.FindAll(triangle, -1)

		parsedVertices := [3]common.Vertex{
			parseStlAsciiLine(vertices[0]),
			parseStlAsciiLine(vertices[1]),
			parseStlAsciiLine(vertices[2]),
		}

		faceVertexIndex := [3]uint32{
			getVertexIndex(parsedVertices[0]),
			getVertexIndex(parsedVertices[1]),
			getVertexIndex(parsedVertices[2]),
		}

		faces = append(faces, common.Face{faceVertexIndex[0], faceVertexIndex[1], faceVertexIndex[2]})
	}
	return common.Mesh{
		Faces:    faces,
		Vertices: vertices,
	}
}

func WriteAsciiStlFromMesh(mesh common.Mesh, meshMetadata common.MeshMetadata) []byte {
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
