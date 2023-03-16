package common

import "testing"

func TestSplitStringByWhiteSpaceNL(t *testing.T) {

	testReturn := func(inputStringPtr *[]string) {
		expectedSingleLine := []string{
			"1.000",
			"world",
		}
		if len(expectedSingleLine) != len(*inputStringPtr) {
			t.Errorf("len different, expecting %v, got %v", len(expectedSingleLine), len(*inputStringPtr))
			return
		}
		for idx, el := range expectedSingleLine {
			if el != (*inputStringPtr)[idx] {
				t.Errorf("cannot parse single line correctly, expected %v, got %v", el, (*inputStringPtr)[idx])
			}
		}

		for idx, el := range *inputStringPtr {
			if el != expectedSingleLine[idx] {
				t.Errorf("cannot parse single line correctly, expected %v, got %v", expectedSingleLine[idx], el)
			}
		}
	}
	singleLine := "1.000 world"
	splitSingleLine := SplitStringByWhiteSpaceNL(singleLine)
	testReturn(&splitSingleLine)

	multiLine := `1.000
	world`
	splitMultiLine := SplitStringByWhiteSpaceNL(multiLine)
	testReturn(&splitMultiLine)

	singleLineNeedsTrimming := "  1.000   world    "
	splitSingleLineNeedsTrimming := SplitStringByWhiteSpaceNL(singleLineNeedsTrimming)
	testReturn(&splitSingleLineNeedsTrimming)

	multipleNeedsTrimming := `
	
	1.000
  
	world


	`

	splitMultiLineNeedsTrimming := SplitStringByWhiteSpaceNL(multipleNeedsTrimming)
	testReturn(&splitMultiLineNeedsTrimming)
}

func assertPanic(t *testing.T, f func()) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected to panic, did not")
		}
	}()

	f()
}

func testScaling(t *testing.T) {
	var xformM TransformationMatrix
	xformM.ParseCommaDelimitedString("2,0,0,0,0,2,0,0,0,0,2,0")

	vertex := Vertex{1.0, 2.0, 3.0}
	vertex.Transform(xformM)

	if vertex[0] != 2.0 || vertex[1] != 4.0 || vertex[2] != 6.0 {
		t.Errorf("unexpected scalling output %v", vertex)
	}
}

func testTranslation(t *testing.T) {

	var xformM TransformationMatrix
	xformM.ParseCommaDelimitedString("1,0,0,2,0,1,0,2,0,0,1,2")

	vertex := Vertex{1.0, 2.0, 3.0}
	vertex.Transform(xformM)

	if vertex[0] != 3.0 || vertex[1] != 4.0 || vertex[2] != 5.0 {
		t.Errorf("unexpected scalling output %v", vertex)
	}
}

func TestParseCommaDelimitedString(t *testing.T) {
	var xformM TransformationMatrix
	assertPanic(t, func() {
		xformM.ParseCommaDelimitedString("1,2,3")
	})
	assertPanic(t, func() {
		xformM.ParseCommaDelimitedString("1,0,0,0,0,1,0,0,0,0,1,test")
	})

	testScaling(t)
	testTranslation(t)
}

func TestFlipTriangle(t *testing.T) {
	mesh := Mesh{
		Vertices: []Vertex{},
		Faces: []Face{
			Face{0, 1, 2},
			Face{0, 2, 3},
		},
	}

	mesh.FlipTriangleOrder()

	if mesh.Faces[0][0] != 2 || mesh.Faces[0][1] != 1 || mesh.Faces[0][2] != 0 {
		t.Errorf("Expected face order flipped, did not")
	}
	if mesh.Faces[1][0] != 3 || mesh.Faces[1][1] != 2 || mesh.Faces[1][2] != 0 {
		t.Errorf("Expected face order flipped, did not, 2")
	}
}

func TestDet(t *testing.T) {

	var xformM TransformationMatrix
	xformM.ParseCommaDelimitedString("1,0,0,0,0,1,0,0,0,0,1,0")

	if xformM.Det() != 1.0 {
		t.Errorf("expected det to equal 1, but does not")
	}

	xformM.ParseCommaDelimitedString("0,1,0,0,1,0,0,0,0,0,1,0")

	if xformM.Det() != -1.0 {
		t.Errorf("expeected det to equal -1, but does not")
	}
}

func TestSub(t *testing.T) {
	v1 := [3]float32{1,2,3}
	v2 := [3]float32{0.1,0.2,0.3}
	vr := Sub(v1,v2)
	vexpect := [3]float32{0.9, 1.8, 2.7}
	for idx := range(vr) {
		if vr[idx] != vexpect[idx] {
			t.Errorf("Sub fails, at idx %v. Expect %v, got %v\n", idx, vexpect[idx], vr[idx])
		}
	}
}

func TestCross(t *testing.T) {

	// basic usage
	{
		v1 := [3]float32{2,0,0}
		v2 := [3]float32{0,2,0}
		vr := Cross(v1,v2)
		vexpect := [3]float32{0, 0, 4}
		for idx := range(vr) {
			if vr[idx] != vexpect[idx] {
				t.Errorf("Cross #1 fails, at idx %v. Expect %v, got %v\n", idx, vexpect[idx], vr[idx])
			}
		}
	}

	// more advanced usage?
	{
		v1 := [3]float32{2,2,0}
		v2 := [3]float32{2,2,2}
		vr := Cross(v1,v2)
		vexpect := [3]float32{4, -4, 0}
		for idx := range(vr) {
			if vr[idx] != vexpect[idx] {
				t.Errorf("Cross #1 fails, at idx %v. Expect %v, got %v\n", idx, vexpect[idx], vr[idx])
			}
		}
	}
}

func TestGetNormal(t *testing.T) {

	threshold := float32(1e-3)

	v0 := [3]float32{0,0,0}
	v1 := [3]float32{2,2,0}
	v2 := [3]float32{2,2,2}
	vertex := [3]Vertex{v0, v1, v2}
	vr := GetNormal(vertex)
	vexpect := [3]float32{0.7071, -0.7071, 0}
	for idx := range(vr) {
		if diff := vr[idx] - vexpect[idx]; (diff * diff) > threshold {
			t.Errorf("GetNormal fails, at idx %v. Expect %v, got %v\n", idx, vexpect[idx], vr[idx])
		}
	}
}