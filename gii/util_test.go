package gii

import (
	"testing"
)

func TestUmarshalText(t *testing.T) {
	value := `1.000000 0.000000 0.000000 0.000000 0.000000 1.000000 0.000000 0.000000 0.000000 0.000000 1.000000 0.000000 0.000000 0.000000 0.000000 1.000000`
	xform := TransformationMatrix{}
	err := xform.UnmarshalText([]byte(value))
	if err != nil {
		panic(err)
	}
	if !xform.IsIdentity() {
		t.Errorf("should be identity")
	}
}
