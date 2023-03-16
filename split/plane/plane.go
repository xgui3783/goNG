package plane

import (
	"common"
	"errors"
	"flag"
	"fmt"
	"gongcommon"
	"gongio"
	"math"
	"os"
	"path"
	"path/filepath"
	"splitCommon"
	"strconv"
	"strings"
)

const (
	MethodName = "byPlane"
	HelperText = "TODO add helper text"
)

var subCmdMap = map[string]*gongcommon.SubCmd{}

func init() {
	// append subcommand here
}

func Parse(fs *flag.FlagSet) (returnErr error) {
	defer func() {
		if r := recover(); r != nil {
			errorText := fmt.Sprintf("An error occured: %v\n", r)
			returnErr = errors.New(errorText)
			return
		}
	}()

	ii := gongio.NewInputInterface()
	ii.SetupFlag(fs)

	pointsOnPlaneHelperTxt := fmt.Sprint("3 vectors describing the plane, comma separated.\ne.g. 0,0,0,1,1,1,1,1,0 will correspond to v1(0,0,0), v2(1,1,1), v3(1,1,0)\nIgnored if octtree flag is set")
	pointsOnPlanePtr := fs.String("pts", "", pointsOnPlaneHelperTxt)

	octTreeLevelHelperText := fmt.Sprintf("Subdivide mesh via oct tree division")
	octTreeLevel := fs.Int("octtree", 0, octTreeLevelHelperText)

	octTreeMinHelperText := fmt.Sprintf("Comma delimited triplet of the min of octtree division. Each value of the triplet. If undefined, will use min vertex of mesh. Each element must be < corresponding element in octree-max")
	octTreeMin := fs.String("octtree-min", "", octTreeMinHelperText)

	octTreeMaxHelperText := fmt.Sprintf("Comma delimited triplet of the max of octtree division. Each value of the triplet. If undefined, will use max vertex of mesh. Each element must be > corresponding element in octree-min")
	octTreeMax := fs.String("octtree-max", "", octTreeMaxHelperText)

	octTreeFinalPartitionHelperText := fmt.Sprintf("Partition the final meshes into octtree, but do not further divide them.")
	octTreeFinalPartition := fs.Bool("octtree-part", false, octTreeFinalPartitionHelperText)

	octTreeDebugHelperText := fmt.Sprintf("Debug intermediate meshes")
	octTreeDebug := fs.Bool("octtree-debug", false, octTreeDebugHelperText)

	// octTreeFinalImprintHelperText := fmt.Sprintf("Bool. After final octtree decimation, octtree cut the final meshes into 8, but do not separate the mesh. This is to satisfy the neuroglancer precomputed multiresolution meshes.")
	// octTreeFinalImprint := fs.Bool("octtree-final-imprint", false, octTreeFinalImprintHelperText)

	fs.Parse(os.Args[3:])

	var octTreeMinVertexPtr *[3]float32
	var octTreeMaxVertexPtr *[3]float32
	if *octTreeMin != "" {
		tmpMin := common.ParseStringAsFloatsWDelimiter(*octTreeMin, ",")
		if len(tmpMin) != 3 {
			returnErr = errors.New("octtree-min must be three floats comma separated\n")
			return
		}
		octTreeMinVertexPtr = &[3]float32{tmpMin[0], tmpMin[1], tmpMin[2]}
	}
	if *octTreeMax != "" {
		tmpMax := common.ParseStringAsFloatsWDelimiter(*octTreeMax, ",")
		if len(tmpMax) != 3 {
			returnErr = errors.New("octtree-max must be three floats comma separated\n")
			return
		}
		octTreeMaxVertexPtr = &[3]float32{tmpMax[0], tmpMax[1], tmpMax[2]}
	}

	splitMeshesMap := map[string]*common.Mesh{}

	meshes, err := ii.GetMesh()
	if err != nil {
		returnErr = err
		return
	}
	if *octTreeLevel != 0 {
		for meshIdx, mesh := range *meshes {
			currSplitLvl := 0
			srcMeshMapPtr := &map[string]*common.Mesh{"": &mesh}
			if octTreeMinVertexPtr == nil || octTreeMaxVertexPtr == nil {
				err, minMax := common.GetVertexMinMax(&mesh)
				if err != nil {
					returnErr = err
					return
				}
				if octTreeMinVertexPtr == nil {
					octTreeMinVertexPtr = &[3]float32{minMax[0][0], minMax[0][1], minMax[0][2]}
				}
				if octTreeMaxVertexPtr == nil {
					fmt.Printf("-octtree-max unset, using mesh max [%v, %v, %v] as oct tree max\n", minMax[1][0], minMax[1][1], minMax[1][2])
					octTreeMaxVertexPtr = &[3]float32{minMax[1][0], minMax[1][1], minMax[1][2]}
				}
			}

			deltaX := float32(octTreeMaxVertexPtr[0] - octTreeMinVertexPtr[0])
			deltaY := float32(octTreeMaxVertexPtr[1] - octTreeMinVertexPtr[1])
			deltaZ := float32(octTreeMaxVertexPtr[2] - octTreeMinVertexPtr[2])

			getPtFromMeshKey := func(meshKey string) (*[3]float32, error) {

				splitKeys := strings.Split(meshKey, "")
				xFraction := float32(0.5)
				yFraction := float32(0.5)
				zFraction := float32(0.5)
				for splitIdx, splitKey := range splitKeys {
					octtreeIdx, err := strconv.Atoi(splitKey)
					if err != nil {
						return nil, err
					}
					if octtreeIdx%2 == 1 {
						xFraction += 1 / float32(math.Pow(2, float64(splitIdx+2)))
					} else {
						xFraction -= 1 / float32(math.Pow(2, float64(splitIdx+2)))
					}
					if octtreeIdx%4 > 1 {
						yFraction += 1 / float32(math.Pow(2, float64(splitIdx+2)))
					} else {
						yFraction -= 1 / float32(math.Pow(2, float64(splitIdx+2)))
					}
					if octtreeIdx > 3 {
						zFraction += 1 / float32(math.Pow(2, float64(splitIdx+2)))
					} else {
						zFraction -= 1 / float32(math.Pow(2, float64(splitIdx+2)))
					}
				}
				return &([3]float32{
					octTreeMinVertexPtr[0] + deltaX*xFraction,
					octTreeMinVertexPtr[1] + deltaY*yFraction,
					octTreeMinVertexPtr[2] + deltaZ*zFraction,
				}), nil
			}
			for {
				resultSplitMeshPtr := &map[string]*common.Mesh{}
				for key, srcMesh := range *srcMeshMapPtr {
					pt, err := getPtFromMeshKey(key)
					if err != nil {
						returnErr = err
						return
					}
					err, returnMeshPtr := splitCommon.SplitMeshByPointCardinalPlanes(srcMesh, *pt)
					if err != nil {
						returnErr = err
						return
					}

					for idx, returnMesh := range *returnMeshPtr {
						meshKey := fmt.Sprintf("%v%v", key, idx)
						(*resultSplitMeshPtr)[meshKey] = returnMesh
					}
				}
				if *octTreeDebug {
					for key, mesh := range *resultSplitMeshPtr {
						(*srcMeshMapPtr)[key] = mesh
					}
				} else {
					srcMeshMapPtr = &(*resultSplitMeshPtr)
				}
				currSplitLvl += 1
				if currSplitLvl >= *octTreeLevel {
					if *octTreeFinalPartition {
						fmt.Printf("octtree-part set, cutting final mesh into octtree\n")
						for meshKey, mesh := range *srcMeshMapPtr {
							pt, err := getPtFromMeshKey(meshKey)
							if err != nil {
								returnErr = err
								return
							}
							err, meshPtr := splitCommon.CutMeshByPointCardinalPlanes(mesh, *pt)
							if err != nil {
								returnErr = err
								return
							}
							(*srcMeshMapPtr)[meshKey] = meshPtr
						}
					}
					break
				}
			}

			for subMeshKey, mesh := range *srcMeshMapPtr {
				finalKey := fmt.Sprintf("%v_%v", meshIdx, subMeshKey)
				splitMeshesMap[finalKey] = mesh
			}
		}
	} else if *pointsOnPlanePtr != "" {

		parsedPoints := common.ParseStringAsFloatsWDelimiter(*pointsOnPlanePtr, ",")
		if len(parsedPoints) != 9 {
			errorText := fmt.Sprintf(
				"pts argument requires exactly 9 elements, delimited by commas (,). However, %v was found with provided argument %v\n",
				len(parsedPoints),
				*pointsOnPlanePtr)
			returnErr = errors.New(errorText)
			return
		}

		p1 := common.Vec3(parsedPoints[0:3])
		p2 := common.Vec3(parsedPoints[3:6])
		p3 := common.Vec3(parsedPoints[6:9])

		if err != nil {
			panic(err)
		}
		for meshIdx, mesh := range *meshes {
			err, splitMeshPtr := splitCommon.SplitMeshByPlane(&mesh, [3][3]float32{p1, p2, p3})
			if err != nil {
				returnErr = err
				return
			}
			for subMeshKey, sMPtr := range *splitMeshPtr {
				if sMPtr == nil {
					continue
				}
				meshKey := fmt.Sprintf("%d_%v", meshIdx, subMeshKey)
				splitMeshesMap[meshKey] = sMPtr
			}
		}
	} else {
		returnErr = errors.New("either --octtree or --pts must be set for byPlane")
		return
	}

	gongio.Mkdir(*ii.Out)

	for meshKey, mesh := range splitMeshesMap {
		// do not write mesh if empty mesh
		if len(mesh.Vertices) == 0 {
			continue
		}
		rbytes, err := ii.GetBytes(&[]common.Mesh{*mesh})
		if err != nil {
			return err
		}
		fragmentFilename := path.Join(*ii.Out, meshKey)
		ext := filepath.Ext(*ii.Out)
		filename := fmt.Sprintf("%v%v", fragmentFilename, ext)

		gongio.WriteBytesToFile(filename, rbytes[0])
	}
	return
}
