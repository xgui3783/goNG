package stlBinary

import (
	"bytes"
	"gong/common"
	"gong/stlAscii"
	"io/ioutil"
	"testing"
)

func TestWriteBinaryStlFromMesh(t *testing.T) {
	cubeMesh := stlAscii.Import("../testData/cube.stl")
	outgoingBytes := make([][]byte, 0)
	for idx, mesh := range cubeMesh {
		outgoingBytes = append(outgoingBytes, WriteBinaryStlFromMesh(mesh, common.MeshMetadata{Index: idx}))
	}
	testCubeBin, err := ioutil.ReadFile("../testData/cube_binary.stl")
	if err != nil {
		panic(err)
	}
	if !bytes.Equal(testCubeBin, outgoingBytes[0]) {
		t.Errorf("created stl bin is not the same as the one in the file")
	}
}
