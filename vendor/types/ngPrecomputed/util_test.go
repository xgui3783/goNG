package ngPrecomputed

import (
	"fmt"
	"io/ioutil"
	"testing"
)

func testStringArrays(t *testing.T, arr1 []string, arr2 []string) {

	map1 := map[string]bool{}
	map2 := map[string]bool{}

	// set maps
	for _, key := range arr1 {
		map1[key] = true
	}

	for _, key := range arr2 {
		map2[key] = true
	}

	// test maps
	for _, key := range arr1 {
		if map2[key] != true {
			t.Errorf("arr1 map2 with key %v not found", key)
		}
	}

	for _, key := range arr2 {
		if map1[key] != true {
			t.Errorf("arr2 map1 with key %v not found", key)
		}
	}
}

func testIntArrays(t *testing.T, arr1 []int, arr2 []int) {

	map1 := map[int]bool{}
	map2 := map[int]bool{}

	// set maps
	for _, key := range arr1 {
		map1[key] = true
	}

	for _, key := range arr2 {
		map2[key] = true
	}

	// test maps
	for _, key := range arr1 {
		if map2[key] != true {
			t.Errorf("arr1 map2 with key %v not found", key)
		}
	}

	for _, key := range arr2 {
		if map1[key] != true {
			t.Errorf("arr2 map1 with key %v not found", key)
		}
	}
}

const MESH_DIR = "../testData/ngMeshDir"

func TestScanLocalDir(t *testing.T) {
	labelIndicies := ScanLocalDir(MESH_DIR)
	expected := []int{1, 2, 3, 4, 5, 6, 7}
	testIntArrays(t, expected, labelIndicies)
}

func TestGetLocalFragments(t *testing.T) {
	localFragments := GetLocalFragments(MESH_DIR, 1)
	expected := []string{"fragment_1_0", "fragment_1_1"}
	testStringArrays(t, expected, localFragments)
}

func TestParseFragmentBuffer(t *testing.T) {
	buffer, err := ioutil.ReadFile(fmt.Sprintf("%v/fragment_1_0", MESH_DIR))
	if err != nil {
		panic(err)
	}
	ParseFragmentBuffer(buffer)
}
