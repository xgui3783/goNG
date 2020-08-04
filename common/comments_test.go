package common

import (
	"testing"
)

func TestTrimHashComments(t *testing.T) {
	testString := "hello world #blessed"
	comments := TrimHashComments(&testString)
	if comments != "blessed" {
		t.Errorf("comments is not stripped")
	}
	if testString != "hello world " {
		t.Errorf("input not trimmed")
	}
}
