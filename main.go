package main

import (
	"errors"
	"flag"
	"fmt"
	"gong/convert"
	"gong/info"
	"gong/split"
	"gongcommon"
	"gongio/glue"
	"os"
)

var subCmds = []string{}

func init() {
	validSrcFormats := glue.GetSupportedIncTypes()
	validOutputFormats := glue.GetSupportedOutTypes()
	convert.SetFormats(validSrcFormats, validOutputFormats)

	// add split subcommand
	splitSubCmd := gongcommon.SubCmd{
		Name:       split.SubCmd,
		Parse:      split.Parse,
		HelperText: split.HelperText,
	}
	splitSubCmd.Init()
	SubCmdMap[split.SubCmd] = &splitSubCmd
	subCmds = append(subCmds, split.SubCmd)

	convert.MeshTypeToImportMap = &(glue.MeshTypeToImportMap)
	convert.MeshTypeToExportMap = &(glue.MeshTypeToExportMap)

	// add convert subcommand
	convertSubCmd := gongcommon.SubCmd{
		Name:       convert.SubCmd,
		Parse:      convert.Parse,
		HelperText: convert.HelperText,
	}
	convertSubCmd.Init()
	SubCmdMap[convert.SubCmd] = &convertSubCmd
	subCmds = append(subCmds, convert.SubCmd)

	// add info subcommand
	infoSubCmd := gongcommon.SubCmd{
		Name:       info.SubCmd,
		Parse:      info.Parse,
		HelperText: info.HelperText,
	}
	infoSubCmd.Init()
	SubCmdMap[info.SubCmd] = &infoSubCmd
	subCmds = append(subCmds, info.SubCmd)
}

func parseSubcmd(arg []string) error {
	if len(arg) < 1 {
		return errors.New("Subcommand is required!")
	}
	if subCmd, found := SubCmdMap[arg[0]]; !found {
		errorText := fmt.Sprintf(
			"Valid subcommands %v \nSubcommand %v not found",
			subCmds,
			arg[0],
		)
		flag.PrintDefaults()
		return errors.New(errorText)
	} else {
		return subCmd.RunParse()
	}
}

func main() {
	if err := parseSubcmd(os.Args[1:]); err != nil {
		fmt.Println(err)
		flag.PrintDefaults()
		os.Exit(1)
	}
}

var SubCmdMap = map[string]*gongcommon.SubCmd{}
