package stlAscii

import (
	"fmt"
	"io/ioutil"
	"testing"
)

func TestParseAsciiStl(t *testing.T) {
	cube, err := ioutil.ReadFile("../../testData/cube.stl")
	if err != nil {
		panic(err)
	}
	mesh := ParseAsciiStl(cube)
	if len(mesh.Vertices) != 8 {
		panicText := fmt.Sprintf("cube vertices should be 8, but is instead %v", len(mesh.Vertices))
		t.Errorf(panicText)
	}

	if len(mesh.Faces) != 12 {
		panicText := fmt.Sprintf("cube triangles should be 12, but is instead %v", len(mesh.Faces))
		t.Errorf(panicText)
	}
}
