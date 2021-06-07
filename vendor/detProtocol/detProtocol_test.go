package detProtocol

import "testing"

func TestInferProtocol(t *testing.T) {
	httpPath := "https://neuroglancer.humanbrainproject.org/precomputed/JuBrain/v2.2c/colin27_seg"
	localPath := "./something.vtk"

	httpInferProtocol := InferProtocolFromFilename(httpPath)
	if httpInferProtocol != HTTP {
		t.Errorf("HTTP test failed")
	}

	localInferProtocol := InferProtocolFromFilename(localPath)
	if localInferProtocol != Local {
		t.Errorf("local test failed")
	}
}
