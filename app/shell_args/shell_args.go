package shell_args

import (
	"strings"
)

func ParseInput(input string) []string {
	trimmed_input, _ := strings.CutSuffix(input, "\n")
	
	args := make([]string, 0)
	current_arg := make([]rune, 0)
	inside_quotes := false
	for _, char := range trimmed_input {
		if char == '\'' {
			inside_quotes = !inside_quotes
			continue
		}
		if char == ' ' && !inside_quotes {
			if len(current_arg) > 0 {
				args = append(args, string(current_arg))
				current_arg = make([]rune, 0)
			}
		} else {
			current_arg = append(current_arg, char)
		}

	}

	args = append(args, string(current_arg))
	return args
}
