package detType

import "testing"

func TestInferTypeFromFilename(t *testing.T) {

	ng_mesh := "https://neuroglancer.humanbrainproject.org/precomputed/JuBrain/v2.2c/colin27_seg"
	vtk := "./something.vtk"

	appleVal := InferTypeFromFilename(ng_mesh)
	if appleVal != NG_MESH {
		t.Errorf("not right")
	}
	notAppleVal := InferTypeFromFilename(vtk)
	if notAppleVal != VTK {
		t.Errorf("vtk infer failed")
	}
}
