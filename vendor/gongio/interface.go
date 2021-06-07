package gongio

import (
	"common"
	"detType"
	"errors"
	"flag"
	"fmt"
	"gongio/glue"
)

type InputInterface struct {
	SrcFormatHelperText *string
	SrcFormat           *string

	SrcHelperText *string
	Src           *string

	OutFormatHelperText *string
	OutFormat           *string

	OutHelperText *string
	Out           *string
}

func (ii *InputInterface) GetMesh() (meshPtr *[]common.Mesh, e error) {
	defer func() {
		if r := recover(); r != nil {
			errorText := fmt.Sprintf("error: %v", r)
			e = errors.New(errorText)
			return
		}
	}()
	var incFileType string
	if *(ii.SrcFormat) == "" {
		if *(ii.SrcFormat) == "" {
			panic("if stdin is used to provide src, -srcFormat must be defined,")
		}
		incFileType = detType.InferTypeFromFilename(*(ii.Src))
	} else {
		incFileType = *(ii.SrcFormat)
	}
	
	importFn, ok := glue.MeshTypeToImportMap[incFileType]
	if !ok {
		panicText := fmt.Sprintf("intput type %v not supported", incFileType)
		panic(panicText)
	} else {
		iBytes := common.GetResource(*(ii.Src))
		meshes := importFn([][]byte{iBytes})
		meshPtr = &(meshes)
		return
	}
}

func (ii *InputInterface) GetBytes(meshPtr *[]common.Mesh) (rBytes [][]byte, e error) {
	defer func() {
		if r := recover(); r != nil {
			errorText := fmt.Sprintf("error: %v", r)
			e = errors.New(errorText)
			return
		}
	}()

	var outFileType string
	if *(ii.OutFormat) == "" {
		outFileType = detType.InferTypeFromFilename(*(ii.Out))
	} else {
		outFileType = *(ii.OutFormat)
	}
	if exportFn, ok := glue.MeshTypeToExportMap[outFileType]; !ok {
		panicText := fmt.Sprintf("exportFn for %v cannot be found", outFileType)
		panic(panicText)
	} else {
		bByptes := exportFn(*meshPtr)
		rBytes = bByptes
		return
	}
}

func NewInputInterface() *InputInterface {

	dSrcFormatHelperText := "SrcFormatHelperText"
	dSrcHelperText := "SrcHelperText"
	dOutFormatHelperText := "OutFormatHelperText"
	dOutHelperText := "OutHelperText"

	ii := InputInterface{}
	ii.SrcFormatHelperText = &dSrcFormatHelperText
	ii.SrcHelperText = &dSrcHelperText
	ii.OutFormatHelperText = &dOutFormatHelperText
	ii.OutHelperText = &dOutHelperText
	return &ii
}

func (ii *InputInterface) SetupFlag(fs *flag.FlagSet) {

	ii.SrcFormat = (*fs).String("srcFormat", "", *(ii.SrcFormatHelperText))
	ii.Src = (*fs).String("src", "", *(ii.SrcHelperText))
	ii.OutFormat = (*fs).String("outFormat", "", *(ii.OutFormatHelperText))
	ii.Out = (*fs).String("out", "", *(ii.OutHelperText))
}
