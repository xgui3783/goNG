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
