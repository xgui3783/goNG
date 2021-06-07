package detType

import "testing"

func TestInferTypeFromFilename(t *testing.T) {

	ng_mesh := "https://neuroglancer.humanbrainproject.org/precomputed/JuBrain/v2.2c/colin27_seg"
	vtk := "./something.vtk"
	stlBinary := "../testData/cube_binary.stl"
	stlAscii := "../testData/cube.stl"

	appleVal := InferTypeFromFilename(ng_mesh)
	if appleVal != NG_MESH {
		t.Errorf("not right")
	}
	notAppleVal := InferTypeFromFilename(vtk)
	if notAppleVal != VTK {
		t.Errorf("vtk infer failed")
	}

	if InferTypeFromFilename(stlBinary) != STL_BINARY {
		t.Errorf("stl binary file does not return STL_BINARY filetype")
	}

	if InferTypeFromFilename(stlAscii) != STL_ASCII {
		t.Errorf("stl ascii file does not return STL_ASCII filetype")
	}
}
