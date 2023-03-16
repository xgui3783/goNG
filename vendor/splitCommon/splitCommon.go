package splitCommon

import (
	"common"
	"errors"
	"fmt"
	"math"
	"strings"
)

type VertexMetadata struct {
	VertexPtr  common.Vertex
	Polarity   float32
	MeshMapPtr *map[*common.Mesh]uint32
}

func (vm *VertexMetadata) GetNewIndex(meshPtr *common.Mesh) (index uint32, found bool) {
	index, found = (*vm.MeshMapPtr)[meshPtr]
	return
}

// NB this type specifically refer to the scenario where an edge crosses plane
// the OldIdx and Vtxmetadata refers to the **opposite** vertex to the edge!
type CrossBnd struct {
	VtxMetadataPtr    *VertexMetadata
	PosVtxMetadataPtr *VertexMetadata
	NegVtxMetadataPtr *VertexMetadata
	CrossIntercept    [3]float32
	NextVtx           *VertexMetadata
	PrevVtx           *VertexMetadata
}

func SplitMeshByPlane(meshPtr *common.Mesh, planeVec3 [3][3]float32) (error, *map[string]*common.Mesh) {
	planeNormalPtr := common.GetPlaneNormalPtr(planeVec3)
	planePtPtr := &(planeVec3[0])
	posMesh := common.Mesh{
		Vertices: []common.Vertex{},
		Faces:    []common.Face{},
	}
	negMesh := common.Mesh{
		Vertices: []common.Vertex{},
		Faces:    []common.Face{},
	}

	vertexIdxMap := map[uint32]*VertexMetadata{}

	// vertex polarity need only to be calculated once
	// when they are needed again, memoized value can be used
	for idx, vtx := range meshPtr.Vertices {
		vtxIdx := uint32(idx)
		if _, found := vertexIdxMap[vtxIdx]; !found {
			polarity := common.GetVec3PlanePolarity(vtx, planeNormalPtr, planePtPtr)

			meshMap := map[*common.Mesh]uint32{}
			// if a vertex lies exactly on the plane
			// it will be in both hemisphere
			if polarity >= 0 {
				newVtxIdx := posMesh.AddVertex(vtx)
				meshMap[&posMesh] = newVtxIdx
			}
			if polarity <= 0 {
				newVtxIdx := negMesh.AddVertex(vtx)
				meshMap[&negMesh] = newVtxIdx
			}
			vtxMetadata := VertexMetadata{
				Polarity:   polarity,
				MeshMapPtr: &meshMap,
				VertexPtr:  vtx,
			}
			vertexIdxMap[vtxIdx] = &vtxMetadata
		}
	}

	for _, vIndices := range meshPtr.Faces {
		polarityMeta := [3]*VertexMetadata{}
		for idx, vIdx := range vIndices {
			if v, ok := vertexIdxMap[vIdx]; !ok {
				errorText := fmt.Sprintf("vertex with index %v not found in idx map. Exiting \n", vIdx)
				return errors.New(errorText), nil
			} else {
				polarityMeta[idx] = v
			}
		}

		// polarity total is useful, in the case where face does not cross boundry, which polarity to append the new face
		polarityTotal := float32(0)

		// crossBnd is used later to easily figure out how to dissect a face
		crossBnd := []CrossBnd{}
		crossBndMainFlag := false
		for idx, vtxMetaPtr := range polarityMeta {
			otherMeta1Ptr := polarityMeta[(idx+1)%3]
			otherMeta2Ptr := polarityMeta[(idx+2)%3]

			crossBndFlag := otherMeta1Ptr.Polarity*otherMeta2Ptr.Polarity < 0
			if crossBndFlag {
				crossBndMainFlag = true

				line := [2][3]float32{
					otherMeta1Ptr.VertexPtr,
					otherMeta2Ptr.VertexPtr,
				}
				if err, xInter := common.FindIntersect(line, planeNormalPtr, planePtPtr, true); err != nil {
					return err, nil
				} else {
					var posPtPtr *VertexMetadata
					var negPtPtr *VertexMetadata
					if otherMeta1Ptr.Polarity > 0 {
						posPtPtr = otherMeta1Ptr
						negPtPtr = otherMeta2Ptr
					} else {
						posPtPtr = otherMeta2Ptr
						negPtPtr = otherMeta1Ptr
					}
					crossBnd = append(crossBnd, CrossBnd{
						VtxMetadataPtr:    vtxMetaPtr,
						CrossIntercept:    xInter,
						PosVtxMetadataPtr: posPtPtr,
						NegVtxMetadataPtr: negPtPtr,
						NextVtx:           otherMeta1Ptr,
						PrevVtx:           otherMeta2Ptr,
					})
				}
			}
			polarityTotal += vtxMetaPtr.Polarity
		}

		newPosFaces := [][3]uint32{}
		newNegFaces := [][3]uint32{}
		appendFace := func(newFacePtr *[][3]uint32, newVtxIdx [3]uint32) {
			*newFacePtr = append(*newFacePtr, newVtxIdx)
		}

		getVertexIndex := func(vtxMetadataPtr *VertexMetadata, meshPtr *common.Mesh) uint32 {
			returnVal := (*vtxMetadataPtr.MeshMapPtr)[meshPtr]
			return returnVal
		}

		// if does not cross plane boundary, no splitting necessary
		if !crossBndMainFlag {

			// in rare occassions, if polarity total === 0, safest to append face to both polarities
			if polarityTotal >= 0 {
				newVtxIdx := [3]uint32{}
				for idx, ptr := range polarityMeta {
					newVtxIdx[idx] = getVertexIndex(ptr, &posMesh)
				}
				appendFace(&newPosFaces, newVtxIdx)
			}
			if polarityTotal <= 0 {
				newVtxIdx := [3]uint32{}
				for idx, ptr := range polarityMeta {
					newVtxIdx[idx] = getVertexIndex(ptr, &negMesh)
				}
				appendFace(&newNegFaces, newVtxIdx)
			}
		} else {

			// do not expect more than 2 edges cross plane
			if len(crossBnd) > 2 {
				errorText := fmt.Sprintf("crossBnd len larger than 2 (%v)\nAborting...\n", len(crossBnd))
				return errors.New(errorText), nil
			}
			// expect at least 1 plane crossing plane
			if len(crossBnd) == 0 {
				errorText := fmt.Sprintf("crossBnd len == 0 \nAborting...\n")
				return errors.New(errorText), nil
			}

			// crossed boundary, split mesh
			firstCrossFlag := true
			var lastPosFace [3]uint32
			var lastNegFace [3]uint32
			var nextIsNewFlag bool
			for _, cross := range crossBnd {
				newPosMeshVtxIdx := posMesh.AddVertex(cross.CrossIntercept)
				newNegMeshVtxIdx := negMesh.AddVertex(cross.CrossIntercept)

				if cross.VtxMetadataPtr.Polarity >= 0 {
					var nextVtxIdx uint32
					var prevVtxIdx uint32
					// self, corsspt, positive point
					selfIdx := (*(cross.VtxMetadataPtr).MeshMapPtr)[&posMesh]
					posIdx := (*(cross.PosVtxMetadataPtr).MeshMapPtr)[&posMesh]

					nextIsNewFlag = (cross.PrevVtx == cross.PosVtxMetadataPtr)
					if nextIsNewFlag {
						nextVtxIdx = newPosMeshVtxIdx
						prevVtxIdx = posIdx
					} else {
						nextVtxIdx = posIdx
						prevVtxIdx = newPosMeshVtxIdx
					}

					newVtxIdx := [3]uint32{
						selfIdx,
						nextVtxIdx,
						prevVtxIdx,
					}
					appendFace(&newPosFaces, newVtxIdx)
				}
				if cross.VtxMetadataPtr.Polarity <= 0 {

					var nextVtxIdx uint32
					var prevVtxIdx uint32
					// self, corsspt, positive point
					selfIdx := (*(cross.VtxMetadataPtr).MeshMapPtr)[&negMesh]
					negIdx := (*(cross.NegVtxMetadataPtr).MeshMapPtr)[&negMesh]

					nextIsNewFlag = (cross.PrevVtx == cross.NegVtxMetadataPtr)
					if nextIsNewFlag {
						nextVtxIdx = newNegMeshVtxIdx
						prevVtxIdx = negIdx
					} else {
						nextVtxIdx = negIdx
						prevVtxIdx = newNegMeshVtxIdx
					}

					newVtxIdx := [3]uint32{
						selfIdx,
						nextVtxIdx,
						prevVtxIdx,
					}
					appendFace(&newNegFaces, newVtxIdx)
				}

				// If vertex polarity is not zero
				// expect 1 face in each pos/neg to be appended
				if cross.VtxMetadataPtr.Polarity != 0 {

					// lastPosFace is set when the first new vertex is found
					// if nil, set the first new vertex
					if firstCrossFlag {
						firstCrossFlag = false

						if cross.VtxMetadataPtr.Polarity > 0 {
							existingPosVtxIdx, pExistOk := cross.VtxMetadataPtr.GetNewIndex(&posMesh)
							if !pExistOk {
								errorText := fmt.Sprintf("posvtkidx does not exist!")
								return errors.New(errorText), nil
							}
							lastPosFace = [3]uint32{
								0,
								newPosMeshVtxIdx,
								existingPosVtxIdx,
							}

							negIdx, negIdxFound := cross.NegVtxMetadataPtr.GetNewIndex(&negMesh)
							if !negIdxFound {
								errorText := fmt.Sprintf("corresponding neg idx does not exist!")
								return errors.New(errorText), nil
							}
							lastNegFace = [3]uint32{
								0,
								negIdx,
								newNegMeshVtxIdx,
							}
						} else {
							existingNegVtxIdx, nExistOk := cross.VtxMetadataPtr.GetNewIndex(&negMesh)
							if !nExistOk {
								errorText := fmt.Sprintf("negidx does not exist!")
								return errors.New(errorText), nil
							}
							lastNegFace = [3]uint32{
								0,
								newNegMeshVtxIdx,
								existingNegVtxIdx,
							}

							posIdx, posIdxFound := cross.PosVtxMetadataPtr.GetNewIndex(&posMesh)
							if !posIdxFound {
								errorText := fmt.Sprintf("corresponding pos idx does not exist!")
								return errors.New(errorText), nil
							}
							lastPosFace = [3]uint32{
								0,
								posIdx,
								newPosMeshVtxIdx,
							}
						}

						if !nextIsNewFlag {
							lastPosFace = [3]uint32{
								lastPosFace[0],
								lastPosFace[2],
								lastPosFace[1],
							}

							lastNegFace = [3]uint32{
								lastNegFace[0],
								lastNegFace[2],
								lastNegFace[1],
							}
						}
					} else {
						lastPosFace[0] = newPosMeshVtxIdx
						lastNegFace[0] = newNegMeshVtxIdx
						// fmt.Printf("debug %v\n", lastPosFace)
						appendFace(&newPosFaces, lastPosFace)
						appendFace(&newNegFaces, lastNegFace)
					}
				}
			}
		}

		// append new faces
		for _, newFace := range newPosFaces {
			if err := posMesh.AddFace(newFace); err != nil {
				return err, nil
			}
		}
		for _, newFace := range newNegFaces {
			if err := negMesh.AddFace(newFace); err != nil {
				return err, nil
			}
		}
	}

	returnMesh := map[string]*common.Mesh{}
	if len(posMesh.Vertices) > 0 {
		returnMesh["p"] = &posMesh
	} else {
		returnMesh["p"] = nil
	}
	if len(negMesh.Vertices) > 0 {
		returnMesh["n"] = &negMesh
	} else {
		returnMesh["n"] = nil
	}

	return nil, &returnMesh
}

func SplitMeshByPointCardinalPlanes(mesh *common.Mesh, pt [3]float32) (returnErr error, returnMeshPtr *map[string]*common.Mesh) {
	ptx := common.Add(pt, [3]float32{pt[0], 0, 0})
	pty := common.Add(pt, [3]float32{0, pt[1], 0})
	ptz := common.Add(pt, [3]float32{0, 0, pt[2]})
	ptxy := common.Add(pt, [3]float32{pt[0], pt[1], 0})
	ptxz := common.Add(pt, [3]float32{pt[0], 0, pt[2]})
	ptyz := common.Add(pt, [3]float32{0, pt[1], pt[2]})
	plane0 := [3][3]float32{
		pt,
		ptyz,
		pty,
	}

	plane1 := [3][3]float32{
		pt,
		ptxz,
		ptz,
	}

	plane2 := [3][3]float32{
		pt,
		ptxy,
		ptx,
	}

	targetMeshesPtr := &map[string]*common.Mesh{"": mesh}
	for _, pln := range [3][3][3]float32{plane0, plane1, plane2} {
		newMeshPtrs := &map[string]*common.Mesh{}
		for meshKey, meshPtr := range *targetMeshesPtr {
			if meshPtr == nil {
				continue
			}
			err, splitMeshesPtr := SplitMeshByPlane(meshPtr, pln)
			if err != nil {
				returnErr = err
				return
			}
			for splitMeshKey, splitMesh := range *splitMeshesPtr {
				if splitMesh == nil {
					continue
				}
				newKey := fmt.Sprintf("%v%v", meshKey, splitMeshKey)
				(*newMeshPtrs)[newKey] = splitMesh
			}
		}
		targetMeshesPtr = newMeshPtrs
	}
	proxyMap := map[string]*common.Mesh{}
	for meshKey, mesh := range *targetMeshesPtr {
		if mesh == nil {
			continue
		}
		pnKeys := strings.Split(meshKey, "")
		counter := 0
		for pnKeyIdx, pnKey := range pnKeys {
			if pnKey == "n" {
				counter += int(math.Pow(2, float64(pnKeyIdx)))
			}
		}
		meshKey := fmt.Sprintf("%v", counter)
		proxyMap[meshKey] = mesh
	}
	returnMeshPtr = &proxyMap
	return
}

func CutMeshByPointCardinalPlanes(mesh *common.Mesh, pt [3]float32) (error, *common.Mesh) {
	err, meshes := SplitMeshByPointCardinalPlanes(mesh, pt)
	if err != nil {
		return err, nil
	}
	returnMeshPtr := &common.Mesh{}
	for _, meshPtr := range *meshes {
		returnMeshPtr.AppendMesh(meshPtr)
	}
	return nil, returnMeshPtr
}

func SplitCommon() {
	fmt.Println("hello world")
}
