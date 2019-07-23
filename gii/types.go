package gii

/*
* The types is inferred from

GIFTI Surface Data Format

Version 1.0
14 January 2011

available from https://www.nitrc.org/frs/download.php/2871/GIFTI_Surface_Format.pdf

*/

import (
	"fmt"
	"gong/common"
	"strconv"
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
		// if &(m[row]) == nil {
		// 	m[row] = [4]float64{}
		// }
		fmt.Printf("%d, %d, %f\n", row, col, parsedFloat)
		// m[row][col] = parsedFloat
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
