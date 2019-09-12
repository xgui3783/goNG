package ngPrecomputed

import (
	"gong/common"
	"gong/detProtocol"
	"gong/gltf"
	"io/ioutil"
)

func Import(rootPath string, labelMap map[int]gltf.GltfMaterial) []common.Mesh {

	// defaultMaterial := gltf.GltfMaterial{
	// 	Name: "default-material",
	// 	PbrMetallicRoughness: gltf.PbrMetallicRoughness{
	// 		RoughnessFactor: 0.7,
	// 		MetallicFactor:  0.9,
	// 		BaseColorFactor: [4]float32{1.0, 1.0, 1.0, 1.0},
	// 	},
	// }

	protocol := detProtocol.InferProtocolFromFilename(rootPath)

	fragments := make([]string, 0)

	// should not panic, as the protocol is either HTTP or local
	if protocol != detProtocol.HTTP && protocol != detProtocol.Local {
		panic("Unknown protocol")
	}

	if protocol == detProtocol.HTTP {
		if labelMap == nil {
			panic("If ng mesh has http protocol, label map needs to be provided")
		}

		for key := range labelMap {
			for _, fragment := range GetHttpFragments(rootPath, key) {
				fragments = append(fragments, fragment)
			}
		}
	}
	if protocol == detProtocol.Local {

		// TODO directly process info file
		if IsInfoFile(rootPath) {
			panic("cannot process info file yet. use root level (remove /info) and try again")
		}

		// if the file exists, and not an info file, it can only be a single fragment file
		if common.CheckFileExists(rootPath) {
			buf, err := ioutil.ReadFile(rootPath)
			if err != nil {
				panic(err)
			}
			return []common.Mesh{ParseFragmentBuffer(buf)}
		}

		// get labelIndices
		labelIndicies := make([]int, 0)
		if labelMap == nil {
			labelIndicies = ScanLocalDir(rootPath)
		} else {
			for key := range labelMap {
				labelIndicies = append(labelIndicies, key)
			}
		}

		// fetch fragments
		fragments := make([]string, 0)
		for _, labelIndex := range labelIndicies {
			fragments = append(fragments, GetLocalFragments(rootPath, labelIndex)...)
		}
	}

	panic("not yet implemented")
}

func Export(meshes []common.Mesh) [][]byte {
	returnBytes := [][]byte{}
	for _, mesh := range meshes {
		returnBytes = append(returnBytes, WriteFragmentFromMesh(mesh))
	}
	return returnBytes
}