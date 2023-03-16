package stlBinary

import (
	"bytes"
	"common"
	"io/ioutil"
	"testing"
	"types/stl/stlAscii"
)

func TestWriteBinaryStlFromMesh(t *testing.T) {
	singleFile := common.GetResource("../../testData/cube.stl")
	cubeMesh := stlAscii.Import([][]byte{singleFile})
	outgoingBytes := make([][]byte, 0)
	for _, mesh := range cubeMesh {
		outgoingBytes = append(outgoingBytes, WriteBinaryStlFromMesh(mesh))
	}
	testCubeBin, err := ioutil.ReadFile("../../testData/cube_binary.stl")
	if err != nil {
		panic(err)
	}
	if !bytes.Equal(testCubeBin, outgoingBytes[0]) {
		t.Errorf("created stl bin is not the same as the one in the file")
	}
}
