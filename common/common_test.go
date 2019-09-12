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
	defer func () {
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
	assertPanic(t, func(){
		xformM.ParseCommaDelimitedString("1,2,3")
	})
	assertPanic(t, func(){
		xformM.ParseCommaDelimitedString("1,0,0,0,0,1,0,0,0,0,1,test")
	})

	testScaling(t)
	testTranslation(t)
}

func TestFlipTriangle(t *testing.T) {
	mesh := Mesh{
		Vertices: []Vertex{
			Vertex{0.0, 1.0, 2.0},
			Vertex{-1.0, -2.0, -3.0},
		},
		Faces: []Face{
			Face{},
		},
	}

	mesh.FlipTriangleOrder()

	if mesh.Vertices[0][0] != 1.0 || mesh.Vertices[0][1] != 0.0 || mesh.Vertices[0][2] != 2.0 {
		t.Errorf("Expected triangle order flipped, did not")
	}
	if mesh.Vertices[1][0] != -2.0 || mesh.Vertices[1][1] != -1.0 || mesh.Vertices[1][2] != -3.0 {
		t.Errorf("Expected triangle order flipped, did not, neg")
	}
}

func TestDet(t *testing.T){
	
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