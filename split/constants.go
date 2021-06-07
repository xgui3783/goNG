package split

import "fmt"

func getSplitMethodHelperText() string {
	return fmt.Sprintf(
		`Method by which mesh should be split. Supports %v`,
		methodList,
	)	
}
