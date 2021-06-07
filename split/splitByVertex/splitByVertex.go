package splitByVertex

import (
	"common"
	"flag"
	"path"
	"path/filepath"
	"errors"
	"fmt"
	"gongio"
	"os"
)

const (
	MethodName = "splitByVertex"
	HelperText = "TODO add helper text"
)

var config *SplitByVertexConfig

func Parse(fs *flag.FlagSet) error {

	ii := gongio.NewInputInterface()
	ii.SetupFlag(fs)
	
	filepathPtr := fs.String("splitByVertexPath", "", splitMeshByVertexHelperText)
	resolveAmbPtr := fs.String("splitByVertexResolveAmb", EMPTY_LABEL, splitMeshAmbiguousStrategyHelperTxt)
	if filepathPtr == nil {
		panicText := fmt.Sprintf("splitByVertexPath must be provided for %v", MethodName)
		return errors.New(panicText)
	}
	config = &SplitByVertexConfig{
		FilePath:         *filepathPtr,
		ResolveAmbiguity: EMPTY_LABEL,
	}
	if resolveAmbPtr != nil {
		config.ResolveAmbiguity = *resolveAmbPtr
	}

	fs.Parse(os.Args[3:])

	meshes, err := ii.GetMesh()
	if err != nil {
		return err
	}
	if len(*meshes) != 1 {
		panicText := fmt.Sprintf("split by vertex only works with a single mesh input.")
		panic(panicText)
	}
	meshOfInterest := (*meshes)[0]
	meshMap, vertexMap := processSplitMeshByVertexfile(config.FilePath)
	mapNameToMeshPtr := splitMesh(&meshOfInterest, &meshMap, &vertexMap, config)
	gongio.Mkdir(*ii.Out)
	for meshName, mesh := range *mapNameToMeshPtr {
		rbytes, err := ii.GetBytes(&[]common.Mesh{mesh})
		if err != nil {
			return err
		}

		fragmentFilename := path.Join(*ii.Out, meshName)
		ext := filepath.Ext(*ii.Out)
		filename := fmt.Sprintf("%v%v", fragmentFilename, ext)

		gongio.WriteBytesToFile(filename, rbytes[0])
	}
	return nil
}

const (
	EMPTY_LABEL             = "EMPTY_LABEL"
	MAJORITY_OR_FIRST_INDEX = "MAJORITY_OR_FIRST_INDEX"
)

func appendMesh(meshObjMap *map[string]*MeshObj, label string, vertex0 *common.Vertex, vertex1 *common.Vertex, vertex2 *common.Vertex) {

	// check if the current label mesh has already been constructed
	if mesh, ok := (*meshObjMap)[label]; ok {

		// populate vertices not in the list
		triangle := [](*common.Vertex){vertex0, vertex1, vertex2}
		faceToBeAppended := common.Face{}
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
			mesh: common.Mesh{
				Vertices: []common.Vertex{
					*vertex0,
					*vertex1,
					*vertex2,
				},
				Faces: []common.Face{
					common.Face{
						uint32(0),
						uint32(1),
						uint32(2),
					},
				},
			},
			vertexPtrToIdxMap: map[*common.Vertex]uint32{
				vertex0: 0,
				vertex1: 1,
				vertex2: 2,
			},
		}
	}
}

func splitMesh(inputMesh *common.Mesh, meshMap *MeshMap, vertexMap *VertexMap, splitMeshConfig *SplitByVertexConfig) (splitMeshesPtr *map[string]common.Mesh) {
	splitMeshesPtr = &(map[string]common.Mesh{})
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

			switch (*splitMeshConfig).ResolveAmbiguity {
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
