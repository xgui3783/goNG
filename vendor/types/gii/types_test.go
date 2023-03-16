package gii

import (
	"encoding/xml"
	"io/ioutil"
	"testing"
)

func TestUnmarshal(t *testing.T) {
	testBytes, err := ioutil.ReadFile("../testData/test.gii")
	if err != nil {
		panic(err)
	}
	testGii := GIFTI{}
	xml.Unmarshal(testBytes, &testGii)

}
