package split

import (
	"common"
	"errors"
	"flag"
	"fmt"
	"gong/split/splitByVertex"
	"gong/split/plane"
	"os"
	"gongcommon"
)

const SubCmd = "split"
const HelperText = "temp helper text"

var splitMethodsMap = map[string]*SplitMethod{}
var methodList = []string{}
var subCmdMap = map[string]*gongcommon.SubCmd{}

func init() {
	// register splitByVertex method
	// splitMethodsMap[splitByVertex.MethodName] = &SplitMethod{
	// 	name:      splitByVertex.MethodName,
	// 	parseFlag: splitByVertex.ParseFlag,
	// 	split:     splitByVertex.Split,
	// }

	// register split by vertex subcommand
	methodList = append(methodList, splitByVertex.MethodName)
	splitByVertexSubCmd := gongcommon.SubCmd{
		Name: splitByVertex.MethodName,
		HelperText: splitByVertex.HelperText,
		Parse: splitByVertex.Parse,
	}
	subCmdMap[splitByVertex.MethodName] = &splitByVertexSubCmd

	methodList = append(methodList, plane.MethodName)
	planeSubCmd := gongcommon.SubCmd{
		Name: plane.MethodName,
		HelperText: plane.HelperText,
		Parse: plane.Parse,
	}
	subCmdMap[plane.MethodName] = &planeSubCmd

}

func Parse(fs *flag.FlagSet) error {
	arg := os.Args[2:]
	if len(arg) < 1 {
		return errors.New("Subcommand is required!")
	}
	if subCmd, found := subCmdMap[arg[0]]; !found {
		errorText := fmt.Sprintf(
			"Valid subcommands %v \nSubcommand %v not found",
			methodList,
			arg[0],
		)
		return errors.New(errorText)
	} else {
		return subCmd.Parse(fs)
	}
}

type SplitMethod struct {
	name      string
	parseFlag func(*flag.FlagSet)
	split     func(*common.Mesh) (*map[string]common.Mesh, error)
}
