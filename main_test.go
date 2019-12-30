package main

import (
	"io/ioutil"
	"testing"
)

func TestMain(t *testing.T) {
	// check that test files can be read
	cube, err := ioutil.ReadFile("testData/cube.stl")
	if err != nil {
		t.Errorf("reading testData/cube.stl fail. test data not readable")
	}
	if string(cube[:5]) != "solid" {
		t.Errorf("cube.stl first 5 characters does not equal solid, instead it is: %v\nEntire file is:\n%v", string(cube[:5]), string(cube))
	}
}
