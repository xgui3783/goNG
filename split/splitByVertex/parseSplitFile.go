package splitByVertex

import (
	"bytes"
	"common"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

func parseSplitMeshByVertexLine(line string, meshMap *MeshMap, vertexMap *VertexMap) {
	trimmedLine := common.TrimStartingEndingWhiteSpaces(string(line))
	separatedLine := strings.SplitN(trimmedLine, " ", 3)
	if len(separatedLine) < 2 {
		return
	}
	vertexIndex, err := strconv.ParseInt(separatedLine[0], 10, 32)
	if err != nil {
		return
	}

	label := separatedLine[1]
	vertexIndices, ok := (*meshMap)[label]
	if ok {
		(*meshMap)[label] = append(vertexIndices, uint32(vertexIndex))
	} else {
		(*meshMap)[label] = []uint32{uint32(vertexIndex)}
	}

	_, vOk := (*vertexMap)[uint32(vertexIndex)]
	if vOk {
		panicText := fmt.Sprintf("duplicated vertex label for %v", vertexIndex)
		panic(panicText)
	} else {
		(*vertexMap)[uint32(vertexIndex)] = label
	}
	return
}

func processSplitMeshByVertexfile(pathToFile string) (rtMeshMaps MeshMap, rtVertexMap VertexMap) {
	rtMeshMaps = MeshMap{}
	rtVertexMap = VertexMap{}
	if pathToFile == "" {
		return
	}
	if common.CheckFileExists(pathToFile) == false {
		panic("splitMeshByVertex does not exist")
	}
	readbytes, err := ioutil.ReadFile(pathToFile)
	if err != nil {
		panic(err)
	}

	splitBytes := bytes.Split(readbytes, []byte("\n"))

	for _, line := range splitBytes {
		lineString := string(line)
		common.TrimHashComments(&lineString)
		parseSplitMeshByVertexLine(lineString, &rtMeshMaps, &rtVertexMap)
	}
	return
}
