package stlCommon

import (
	"fmt"
	"testing"
)

func TestSerializeVertexTripletToString(t *testing.T) {

	v1 := [3]float32{1.0, 2.0, 3.0}
	v2 := [3]float32{1.00, 2.00, 3.00}
	v3 := [3]float32{1.01, 2.00, 3.00}

	if serializeVertexTripletToString(v1) != serializeVertexTripletToString(v2) {
		errorText := fmt.Sprintf("expecting the serialize fn to produce same hash, but did not")
		t.Errorf(errorText)
	}

	if serializeVertexTripletToString(v1) == serializeVertexTripletToString(v3) {
		errorText := fmt.Sprintf("expecting the serialize fn to produce different hash, but produced same hash")
		t.Errorf(errorText)
	}
}
