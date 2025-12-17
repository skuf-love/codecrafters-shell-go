package shell_args

import (
	"strings"
	"slices"
	//"fmt"
)

type parseContext struct{
	currentArg []rune
	args []string
	mode string
	escape bool
}

type ParsedArgs struct{
	CommandName string
	Arguments []string
	StdoutPath string
}

func (pa ParsedArgs) isStdoutRedirected() bool {
	return len(pa.StdoutPath) > 0
}


func (c *parseContext) normalRead(char rune) {
	if c.escape {
		c.currentArg = append(c.currentArg, char)
		c.escape = false
		return
	}
	if char == '\\' {
		c.escape = true
		return
	}
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

var doubleQuoteEscaped = []rune{'\\', '"'}

func (c *parseContext) doubleQuoteRead(char rune) {
	if c.escape {
		if !slices.Contains(doubleQuoteEscaped, char){
			c.currentArg = append(c.currentArg, '\\')
		}
		c.currentArg = append(c.currentArg, char)
		c.escape = false
		return
	}
	if char == '\\' {
		c.escape = true
		return
	}
	if char == '"' {
		//fmt.Printf("%v - Switching from D to N\n", string(char))
		c.mode = "normal"
		return
	}
	//fmt.Printf("DOUBLE Append '%v'\n", string(char))
	c.currentArg = append(c.currentArg, char)
}
func ParseInput(input string) ParsedArgs {
	trimmed_input, _ := strings.CutSuffix(input, "\n")
	//fmt.Println(trimmed_input)
	
	context := parseContext{
		make([]rune, 0),
		make([]string, 0),
		"normal",
		false,
	}
	for _, char := range trimmed_input {
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

	commandName := context.args[0]
	commandArguments := make([]string, 0)
	stdoutPath := ""
	if len(context.args) > 1 {
		commandArguments = context.args[1:]
		if len(commandArguments) > 1 {
			symIndex := len(commandArguments) - 2
			stdoutPathIndex := len(commandArguments) - 1
			if commandArguments[symIndex] == ">" || commandArguments[symIndex] == "1>" {
				stdoutPath = commandArguments[stdoutPathIndex]
				commandArguments = commandArguments[0:(stdoutPathIndex+1)]
			}
		}
	}

	return ParsedArgs{
		commandName,
		commandArguments,
		stdoutPath,
	} 
}


