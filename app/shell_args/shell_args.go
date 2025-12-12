package shell_args

import (
	"strings"
	//"fmt"
)

type parseContext struct{
	currentArg []rune
	args []string
	mode string
	escape bool
}

func (c *parseContext) normalRead(char rune) {
	if char == '\'' {
		//fmt.Printf("%v - Switching from N to S\n", string(char))
		c.mode = "single_quote"
		return
	}
	if char == '"' {
		//fmt.Printf("%v - Switching from N to D\n", string(char))
		c.mode = "double_quote"
		return
	}
	if char == ' ' && len(c.currentArg) == 0 {
		return 
	}
	if char == ' ' {
		c.args = append(c.args, string(c.currentArg))
		c.currentArg = make([]rune, 0)
	}else{
		//fmt.Printf("NORMAL Append '%v'\n", string(char))
		c.currentArg = append(c.currentArg, char)
	}
}
func (c *parseContext) singleQuoteRead(char rune) {
	if char == '\'' {
		//fmt.Printf("%v - Switching from S to N\n", string(char))
		c.mode = "normal"
		return
	}
	//fmt.Printf("SINGLE Append '%v'\n", string(char))
	c.currentArg = append(c.currentArg, char)
}

func (c *parseContext) doubleQuoteRead(char rune) {
	if char == '"' {
		//fmt.Printf("%v - Switching from D to N\n", string(char))
		c.mode = "normal"
		return
	}
	//fmt.Printf("DOUBLE Append '%v'\n", string(char))
	c.currentArg = append(c.currentArg, char)
}
func ParseInput(input string) []string {
	trimmed_input, _ := strings.CutSuffix(input, "\n")
	//fmt.Println(trimmed_input)
	
	context := parseContext{
		make([]rune, 0),
		make([]string, 0),
		"normal",
		false,
	}
	for _, char := range trimmed_input {
		if context.escape {
			context.currentArg = append(context.currentArg, char)
			context.escape = false
			continue
		}
		if char == '\\' {
			context.escape = true
			continue
		}
		switch context.mode {
		case "normal":
			context.normalRead(char)
		case "single_quote":
			context.singleQuoteRead(char)
		case "double_quote":
			context.doubleQuoteRead(char)
		}
	}

	context.args = append(context.args, string(context.currentArg))
	return context.args
}


