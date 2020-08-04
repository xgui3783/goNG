package common

import (
	"strings"
)

func TrimHashComments(input *string) (comments string) {
	output := strings.SplitN(*input, "#", 2)
	*input = output[0]
	if len(output) < 2 {
		return ""
	}
	return output[1]
}
