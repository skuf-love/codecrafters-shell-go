package shell_args

import (
	"strings"
)

func ParseInput(input string) []string {
	trimmed_input := input[:len(input)-1]

	return strings.Split(trimmed_input, " ")
}
