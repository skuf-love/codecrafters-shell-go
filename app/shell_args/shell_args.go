package shell_args

import (
	"strings"
)

func ParseInput(input string) []string {
	trimmed_input, _ := strings.CutSuffix(input, "\n")

	return strings.Split(trimmed_input, " ")
}
