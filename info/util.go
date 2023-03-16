package info

import (
	"common"
	"encoding/json"
	"fmt"
	"gongio"
)

type InfoStruct struct {
	Min     *common.Vertex `json:"min"`
	Max     *common.Vertex `json:"max"`
	Center  [3]float32     `json:"center"`
	NumVert uint32         `json:"num_vertices"`
	NumFace uint32         `json:"num_faces"`
}

func writeInfo(ii *gongio.InputInterface, numOfLod int) (returnError error) {

	meshes, err := ii.GetMesh()
	if err != nil {
		panic(err)
	}

	if len(*meshes) != 1 {
		panicText := fmt.Sprintf("Expected only a single mesh, got %d", len(*meshes))
		panic(panicText)
	}

	meshOfInterestPtr := &(*meshes)[0]
	err, minMax := common.GetVertexMinMax(meshOfInterestPtr)
	if err != nil {
		panic(err)
	}
	outputInfo := InfoStruct{
		Min:     minMax[0],
		Max:     minMax[1],
		Center:  common.MulScalar(common.Add(*minMax[0], *minMax[1]), 0.5),
		NumVert: uint32(len(meshOfInterestPtr.Vertices)),
		NumFace: uint32(len(meshOfInterestPtr.Faces)),
	}
	jsonOutputInfo, err := json.Marshal(outputInfo)
	if *ii.Out == "" {
		fmt.Printf(string(jsonOutputInfo))
		return nil
	}
	gongio.WriteBytesToFile(*ii.Out, jsonOutputInfo)
	return nil
	// max := minMax[1]
	// buf := new(bytes.Buffer)

	// //chunk_shape (3x float32le finest octtree)
	// for _, coord := range max {
	// 	chunk_shape_frag := coord / float32(math.Pow(2, float64(numOfLod)))
	// 	binary.Write(buf, binary.LittleEndian, float32(chunk_shape_frag))
	// }

	// //grid origin
	// binary.Write(buf, binary.LittleEndian, float32(0))
	// binary.Write(buf, binary.LittleEndian, float32(0))
	// binary.Write(buf, binary.LittleEndian, float32(0))

	// //num_lods
	// binary.Write(buf, binary.LittleEndian, uint32(numOfLod))

	// //lod_scales, `num_lods` iterations
	// for i := 0; i <= numOfLod; i++ {
	// 	binary.Write(buf, binary.LittleEndian, float32(1))
	// }
	// //vertex_offsets, `num_lods*3` iterations
	// for i := 0; i <= numOfLod*3; i++ {
	// 	binary.Write(buf, binary.LittleEndian, float32(0))
	// }
	// //num_fragments_per_lod, `num_lods`
	// for i := numOfLod; i >= 0; i-- {
	// 	// TODO in more efficient/advanced mesh octtree-ification,
	// 	// one may be able to remove a few octrangle (if no vertices exist)
	// 	// but, for now, let's assume the most in efficient octtree-ification
	// 	binary.Write(buf, binary.LittleEndian, uint32(math.Pow(8, float64(i))))
	// }
	// //lod specific
	// for i := numOfLod; i >= 0; i-- {
	// 	//fragment_positions
	// 	// TODO see above... advance usecases may allow us remove some oct-corners
	// 	num_frag := int(math.Pow(8, float64(i)))
	// 	x_idx := 0
	// 	y_idx := 0
	// 	z_idx := 0
	// 	for frag_idx := 0; frag_idx < num_frag; frag_idx++ {
	// 		binary.Write(buf, binary.LittleEndian, uint32(x_idx))
	// 		binary.Write(buf, binary.LittleEndian, uint32(y_idx))
	// 		binary.Write(buf, binary.LittleEndian, uint32(z_idx))

	// 		if x_idx <= y_idx {
	// 			x_idx++
	// 		} else if y_idx <= z_idx {
	// 			y_idx++
	// 		} else {
	// 			z_idx++
	// 		}
	// 	}

	// 	for frag_idx := 0; frag_idx < num_frag; frag_idx++ {
	// 		binary.Write(buf, binary.LittleEndian, uint32(x_idx))
	// 	}

	// }

	// return nil
}
