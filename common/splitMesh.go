package common

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

type MeshMap map[string][]uint32
type VertexMap map[uint32]string

func parseSplitMeshByVertexLine(line string, meshMap *MeshMap, vertexMap *VertexMap) {
	trimmedLine := TrimStartingEndingWhiteSpaces(string(line))
	separatedLine := strings.SplitN(trimmedLine, " ", 3)
	if len(separatedLine) < 2 {
		return
	}
	vertexIndex, err := strconv.ParseInt(separatedLine[0], 10, 32)
	if err != nil {
		return
	}

	label := separatedLine[1]
	vertexIndices, ok := (*meshMap)[label]
	if ok {
		(*meshMap)[label] = append(vertexIndices, uint32(vertexIndex))
	} else {
		(*meshMap)[label] = []uint32{uint32(vertexIndex)}
	}

	_, vOk := (*vertexMap)[uint32(vertexIndex)]
	if vOk {
		panicText := fmt.Sprintf("duplicated vertex label for %v", vertexIndex)
		panic(panicText)
	} else {
		(*vertexMap)[uint32(vertexIndex)] = label
	}
	return
}

func ProcessSplitMeshByVertexfile(pathToFile string) (rtMeshMaps MeshMap, rtVertexMap VertexMap) {
	rtMeshMaps = MeshMap{}
	rtVertexMap = VertexMap{}
	if pathToFile == "" {
		return
	}
	if CheckFileExists(pathToFile) == false {
		panic("splitMeshByVertex does not exist")
	}
	readbytes, err := ioutil.ReadFile(pathToFile)
	if err != nil {
		panic(err)
	}

	splitBytes := bytes.Split(readbytes, []byte("\n"))

	for _, line := range splitBytes {
		lineString := string(line)
		TrimHashComments(&lineString)
		parseSplitMeshByVertexLine(lineString, &rtMeshMaps, &rtVertexMap)
	}
	return
}

type MeshObj struct {
	mesh              Mesh
	vertexPtrToIdxMap map[*Vertex]uint32
}

func appendMesh(meshObjMap *map[string]*MeshObj, label string, vertex0 *Vertex, vertex1 *Vertex, vertex2 *Vertex) {

	// check if the current label mesh has already been constructed
	if mesh, ok := (*meshObjMap)[label]; ok {

		// populate vertices not in the list
		triangle := [](*Vertex){vertex0, vertex1, vertex2}
		faceToBeAppended := Face{}
		for faceVertexIdx, vertex := range triangle {
			if vertexIdx, ok := mesh.vertexPtrToIdxMap[vertex]; ok {
				faceToBeAppended[faceVertexIdx] = vertexIdx
			} else {

				// vertex index not found

				// new vertex index
				newIdx := uint32(len(mesh.vertexPtrToIdxMap))

				// append new vertexPtr to idx mapping
				mesh.vertexPtrToIdxMap[vertex] = newIdx

				// append the new vertex to the mesh
				mesh.mesh.Vertices = append(mesh.mesh.Vertices, *vertex)

				faceToBeAppended[faceVertexIdx] = newIdx
			}
		}

		// all missing vertices should be appended already. append face
		mesh.mesh.Faces = append(mesh.mesh.Faces, faceToBeAppended)

	} else {
		// if the mesh has not yet been constructed
		(*meshObjMap)[label] = &MeshObj{
			mesh: Mesh{
				Vertices: []Vertex{
					*vertex0,
					*vertex1,
					*vertex2,
				},
				Faces: []Face{
					Face{
						uint32(0),
						uint32(1),
						uint32(2),
					},
				},
			},
			vertexPtrToIdxMap: map[*Vertex]uint32{
				vertex0: 0,
				vertex1: 1,
				vertex2: 2,
			},
		}
	}
}

func SplitMesh(inputMesh *Mesh, meshMap *MeshMap, vertexMap *VertexMap, splitMeshConfig *SplitMeshConfig) (splitMeshesPtr *map[string]Mesh) {
	splitMeshesPtr = &(map[string]Mesh{})
	meshObjMap := map[string]*MeshObj{}

	for _, face := range (*inputMesh).Faces {

		vertexLabel0 := (*vertexMap)[face[0]]
		vertexLabel1 := (*vertexMap)[face[1]]
		vertexLabel2 := (*vertexMap)[face[2]]

		vertex0 := &((*inputMesh).Vertices[face[0]])
		vertex1 := &((*inputMesh).Vertices[face[1]])
		vertex2 := &((*inputMesh).Vertices[face[2]])

		// if all three vertex indices agree on label, append label
		if vertexLabel0 == vertexLabel1 && vertexLabel1 == vertexLabel2 {

			// check if the current label mesh has already been constructed

			appendMesh(&meshObjMap, vertexLabel0, vertex0, vertex1, vertex2)
		} else {
			// not all three vertex indices agree to label...

			switch (*splitMeshConfig).UntangleAmbiguityMethod {
			case MAJORITY_OR_FIRST_INDEX:
				if vertexLabel0 == vertexLabel1 || vertexLabel0 == vertexLabel2 {
					appendMesh(&meshObjMap, vertexLabel0, vertex0, vertex1, vertex2)
				} else if vertexLabel1 == vertexLabel2 {
					appendMesh(&meshObjMap, vertexLabel1, vertex0, vertex1, vertex2)
				} else {
					// if all three vertex disagrees, use first index
					appendMesh(&meshObjMap, vertexLabel0, vertex0, vertex1, vertex2)
				}
			case EMPTY_LABEL:
				fallthrough
			default:
				appendMesh(&meshObjMap, EMPTY_LABEL, vertex0, vertex1, vertex2)
			}
		}
	}

	for key, meshObj := range meshObjMap {
		(*splitMeshesPtr)[key] = (*meshObj).mesh
	}

	return
}

type SplitMeshConfig struct {
	SplitMeshByVerticesPath string
	UntangleAmbiguityMethod string
}

const (
	EMPTY_LABEL             = "EMPTY_LABEL"
	MAJORITY_OR_FIRST_INDEX = "MAJORITY_OR_FIRST_INDEX"
)
