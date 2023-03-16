package gii

/*
* The types is inferred from

GIFTI Surface Data Format

Version 1.0
14 January 2011

available from https://www.nitrc.org/frs/download.php/2871/GIFTI_Surface_Format.pdf

*/

import (
	"bytes"
	zlib "compress/zlib"
	b64 "encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"gong/common"
	ioutil "io/ioutil"
	"math"
	"strconv"
	"strings"
)

type CoordinateSystemTransformMatrix struct {
	DataSpace        DataSpace        `xml:"DataSpace"`
	MatrixData       MatrixData       `xml:"MatrixData"`
	TransformedSpace TransformedSpace `xml:"TransformedSpace"`
}

type Data struct {
	Value string `xml:",chardata"`
}

type DataArray struct {
	Data                            Data                              `xml:"Data"`                            /* 1 */
	MetaData                        MetaData                          `xml:"MetaData"`                        /* ? */
	CoordinateSystemTransformMatrix []CoordinateSystemTransformMatrix `xml:"CoordinateSystemTransformMatrix"` /* + */

	ArrayIndexingOrder string `xml:"ArrayIndexingOrder,attr"`
	DataType           string `xml:"DataType,attr"`
	Dimensionality     int    `xml:"Dimensionality,attr"`
	Dim0               int    `xml:"Dim0,attr"`
	Dim1               int    `xml:"Dim1,attr"`
	Dim2               int    `xml:"Dim2,attr"`
	// Dim0 Dim1 ... etc

	Encoding           string `xml:"Encoding,attr"`
	Endian             string `xml:"Endian,attr"`
	ExternalFileName   string `xml:"ExternalFileName,attr"`
	ExternalFileOffset string `xml:"ExternalFileOffset,attr"`
	Intent             string `xml:"Intent,attr"`
}

func (da *DataArray) getValue() []byte {
	var data = []byte(da.Data.Value)
	if strings.Contains(da.Encoding, "Base64Binary") {
		var err error
		data, err = b64.StdEncoding.DecodeString(da.Data.Value)
		if err != nil {
			panic("decodeing base64 error")
		}
	}

	if strings.Contains(da.Encoding, "GZip") {
		reader, err := zlib.NewReader(bytes.NewReader(data))
		if err != nil {
			panic("gunzip failed!")
		}
		var _err error
		data, _err = ioutil.ReadAll(reader)
		if _err != nil {
			panic("reading failed!")
		}
	}
	return data
}

type ParseNum struct {
	Uint32  func([]byte) uint32
	Int8    func([]byte) uint32
	Float32 func([]byte) float32
}

func (da *DataArray) parseStringDataAsFloat() ([][3]float32, error) {
	if strings.Contains(da.Encoding, "Base64Binary") {
		return nil, errors.New("Cannot parse data as if string if encoding contains base64binary")
	}
	splitF := common.SplitStringByWhiteSpaceNL(string(da.getValue()))
	if len(splitF)%3 != 0 {
		panicText := fmt.Sprintf("numver of values of NIFTI_INTENT_POINTSET is not a multiple of 3: it is %d", len(splitF))
		panic(panicText)
	}

	output := [][3]float32{}
	for i := 0; i < len(splitF)/3; i++ {
		v := parseStringsToFloat([3]string{splitF[i*3], splitF[i*3+1], splitF[i*3+2]})
		output = append(output, v)
	}
	return output, nil
}

func (da *DataArray) parseStringDataAsInt() ([][3]uint32, error) {
	if strings.Contains(da.Encoding, "Base64Binary") {
		return nil, errors.New("Cannot parse data as if string if encoding contains base64binary")
	}
	splitF := common.SplitStringByWhiteSpaceNL(string(da.getValue()))
	if len(splitF)%3 != 0 {
		panicText := fmt.Sprintf("numver of values of NIFTI_INTENT_POINTSET is not a multiple of 3: it is %d", len(splitF))
		panic(panicText)
	}

	output := [][3]uint32{}
	for i := 0; i < len(splitF)/3; i++ {
		v := parseStringsToInt([3]string{splitF[i*3], splitF[i*3+1], splitF[i*3+2]})
		output = append(output, v)
	}
	return output, nil
}

func (da *DataArray) getParseNum() ParseNum {
	if da.Endian == LittleEndian {
		return ParseNum{
			Uint32: func(b []byte) uint32 {
				if len(b) != 4 {
					panic("Uint32 must be exactly 4 bytes")
				}

				return binary.LittleEndian.Uint32(b)
			},
			Int8: func(b []byte) uint32 {
				if len(b) != 1 {
					panic("Int8 input must be a single byte")
				}
				return uint32(b[0])
			},
			Float32: func(b []byte) float32 {
				if len(b) != 4 {
					panic("Uint32 must be exactly 4 bytes")
				}
				return math.Float32frombits(binary.LittleEndian.Uint32(b))
			},
		}
	}
	if da.Endian == BigEndian {
		return ParseNum{
			Uint32: func(b []byte) uint32 {
				if len(b) != 4 {
					panic("Uint32 must be exactly 4 bytes")
				}
				return binary.BigEndian.Uint32(b)
			},
			Int8: func(b []byte) uint32 {
				if len(b) != 1 {
					panic("Int8 input must be a single byte")
				}
				return uint32(b[0])
			},
			Float32: func(b []byte) float32 {
				if len(b) != 4 {
					panic("Uint32 must be exactly 4 bytes")
				}
				return math.Float32frombits(binary.BigEndian.Uint32(b))
			},
		}
	}
	panic("Endian is not set")
}

func (da *DataArray) getFloatTriplets() [][3]float32 {
	if da.DataType != NIFTI_TYPE_FLOAT32 {
		panicText := fmt.Sprintf("datatype is not float32, but is instead %v", da.DataType)
		panic(panicText)
	}
	if da.Dim1 != 3 {
		panicText := fmt.Sprintf("dim1 needs to be 3, but is instead %v", da.Dim1)
		panic(panicText)
	}
	bin := da.getValue()
	if da.Dim0*da.Dim1*4 != len(bin) {
		panicText := fmt.Sprintf("dim0 * dim1 * 4 is %v, but len of expanded is instead %v", da.Dim0*da.Dim1*4, da.Dim1)
		panic(panicText)
	}
	parseNum := da.getParseNum()
	returnArray := make([][3]float32, 0)
	for idx := 0; idx < da.Dim0; idx++ {
		startIdx := idx * 4 * da.Dim1
		newTriplet := [3]float32{
			parseNum.Float32(bin[startIdx : startIdx+4]),
			parseNum.Float32(bin[startIdx+4 : startIdx+8]),
			parseNum.Float32(bin[startIdx+8 : startIdx+12]),
		}
		returnArray = append(returnArray, newTriplet)
	}
	return returnArray
}

func (da *DataArray) getIntTriplets() [][3]uint32 {
	if da.DataType != NIFTI_TYPE_INT32 && da.DataType != NIFTI_TYPE_UINT8 {
		panicText := fmt.Sprintf("datatype is not float32, but is instead %v", da.DataType)
		panic(panicText)
	}
	if da.Dim1 != 3 {
		panicText := fmt.Sprintf("dim1 needs to be 3, but is instead %v", da.Dim1)
		panic(panicText)
	}
	bin := da.getValue()
	size_per_num := 0
	var fn func([]byte) uint32
	parseNum := da.getParseNum()
	if da.DataType == NIFTI_TYPE_INT32 {
		size_per_num = 4
		fn = parseNum.Uint32
	}
	if da.DataType == NIFTI_TYPE_UINT8 {
		size_per_num = 1
		fn = parseNum.Int8
	}
	if da.Dim0*da.Dim1*size_per_num != len(bin) {
		panicText := fmt.Sprintf("dim0 * dim1 * %v is %v, but len of expanded is instead %v", size_per_num, da.Dim0*da.Dim1*4, da.Dim1)
		panic(panicText)
	}
	returnArray := make([][3]uint32, 0)
	for idx := 0; idx < da.Dim0; idx++ {
		startIdx := idx * size_per_num * da.Dim1
		newTriplet := [3]uint32{
			fn(bin[startIdx : startIdx+size_per_num]),
			fn(bin[startIdx+size_per_num : startIdx+(size_per_num*2)]),
			fn(bin[startIdx+(size_per_num*2) : startIdx+(size_per_num*3)]),
		}
		returnArray = append(returnArray, newTriplet)
	}
	return returnArray
}

type DataSpace struct {
	Value string `xml:",chardata"`
}

type GIFTI struct {
	MetaData   MetaData    `xml:"MetaData"`   /* ? */
	LabelTable LabelTable  `xml:"LabelTable"` /* ? */
	DataArray  []DataArray `xml:"DataArray"`  /* + */

	NumberOfDataArrays           int    `xml:"NumberOfDataArrays,attr"`
	Version                      string `xml:"Version,attr"`
	XMLNSXSI                     string `xml:"xmlns:xsi,attr"`
	XSINONAMESPACESCHEMALOCATION string `xml:"xsi:noNamespaceSchemaLocation,attr"`
}

type Label struct {
	Key   int     `xml:"Key,attr"`
	Red   float32 `xml:"Red,attr"`
	Blue  float32 `xml:"Blue,attr"`
	Green float32 `xml:"Green,attr"`
	Alpha float32 `xml:"Alpha,attr"`
}

type LabelTable struct {
	Label []Label `xml:"Label"`
}

type MatrixData struct {
	Value string `xml:",chardata"`
	// Matrix common.TransformationMatrix `xml:",chardata"`
}

type MetaData struct {
	MD []MD `xml:"MD"` /* * */
}

type MD struct {
	Name  CData `xml:"Name"`
	Value CData `xml:"Value"`
}

// Name

type TransformedSpace struct {
	Value string `xml:",chardata"`
}

// Value

type CData struct {
	Value string `xml:",cdata"`
}

type TransformationMatrix common.TransformationMatrix

func (m *TransformationMatrix) UnmarshalText(text []byte) error {
	splitString := common.SplitStringByWhiteSpaceNL(string(text))
	if len(splitString) != 16 {
		panicText := fmt.Sprintf("Unmarshalling matrix string split length does not equal 16, but is instead: %d\n", len(splitString))
		panic(panicText)
	}

	for idx, s := range splitString {

		row := int(idx / 4)
		col := int(idx % 4)

		parsedFloat, err := strconv.ParseFloat(s, 64)
		if err != nil {
			fmt.Printf("return error")
			return err
		}
		m[row][col] = parsedFloat
	}
	return nil
}

func (m TransformationMatrix) MarshalText() ([]byte, error) {
	returnString := ""
	for _, r := range m {
		for _, c := range r {
			if returnString != "" {
				returnString = returnString + " "
			}

			returnString = fmt.Sprintf("%v%f", returnString, c)
		}
	}
	return []byte(returnString), nil
}

const (
	SpatialUnit = "mm"

	// ArrayIndexingOrder
	RowMajorOrder    = `RowMajorOrder`
	ColumnMajorOrder = `ColumnMajorOrder`

	// DataType
	NIFTI_TYPE_UINT8   = `NIFTI_TYPE_UINT8`
	NIFTI_TYPE_INT32   = `NIFTI_TYPE_INT32`
	NIFTI_TYPE_FLOAT32 = `NIFTI_TYPE_FLOAT32`

	// Encoding
	ASCII              = `ASCII`
	Base64Binary       = `Base64Binary`
	GZipBase64Binary   = `GZipBase64Binary`
	ExternalFileBinary = `ExternalFileBinary`

	// Endian
	BigEndian    = `BigEndian`
	LittleEndian = `LittleEndian`

	// Intent
	NIFTI_INTENT_GENMATRIX   = "NIFTI_INTENT_GENMATRIX"
	NIFTI_INTENT_LABEL       = "NIFTI_INTENT_LABEL"
	NIFTI_INTENT_NODE_INDEX  = "NIFTI_INTENT_NODE_INDEX"
	NIFTI_INTENT_POINTSET    = "NIFTI_INTENT_POINTSET"
	NIFTI_INTENT_RGB_VECTOR  = "NIFTI_INTENT_RGB_VECTOR"
	NIFTI_INTENT_RGBA_VECTOR = "NIFTI_INTENT_RGBA_VECTOR"
	NIFTI_INTENT_SHAPE       = "NIFTI_INTENT_SHAPE"
	NIFTI_INTENT_TIME_SERIES = "NIFTI_INTENT_TIME_SERIES"
	NIFTI_INTENT_TRIANGLE    = "NIFTI_INTENT_TRIANGLE"
	NIFTI_INTENT_NONE        = "NIFTI_INTENT_NONE"
	NIFTI_INTENT_VECTOR      = "NIFTI_INTENT_VECTOR"

	// Dataspace || TransformedSpace
	NIFTI_XFORM_UNKNOWN      = `NIFTI_XFORM_UNKNOWN`
	NIFTI_XFORM_SCANNER_ANAT = `NIFTI_XFORM_SCANNER_ANAT`
	NIFTI_XFORM_ALIGNED_ANAT = `NIFTI_XFORM_ALIGNED_ANAT`
	NIFTI_XFORM_TALAIRACH    = `NIFTI_XFORM_TALAIRACH`
	NIFTI_XFORM_MNI_152      = `NIFTI_XFORM_MNI_152`
)

/**
* zero value for slice is nil
 */
func (dataArray *DataArray) IsAssigned() bool {
	return dataArray.Intent != ""
}

func (m *TransformationMatrix) IsIdentity() bool {
	for rIdx, r := range *m {
		for cIdx, c := range r {
			if rIdx == cIdx {
				if c != 1.0 {
					return false
				}
			} else {
				if c != 0.0 {
					return false
				}
			}
		}
	}
	return true
}
