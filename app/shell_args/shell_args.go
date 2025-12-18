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
	StderrPath string
	AppendStdout bool
	AppendStderr bool
}

func (pa ParsedArgs) IsStdoutRedirected() bool {
	return len(pa.StdoutPath) > 0
}

func (pa ParsedArgs) IsStderrRedirected() bool {
	return len(pa.StderrPath) > 0
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
	if len(context.args) > 1 {
		commandArguments = context.args[1:]
	}
	commandArguments, stdoutPath, stderrPath, appendStdout, appendStderr := parseRedirects(commandArguments)

	return ParsedArgs{
		commandName,
		commandArguments,
		stdoutPath,
		stderrPath,
		appendStdout,
		appendStderr,
	} 
}

func parseRedirects(commandArguments []string) ([]string, string, string,
bool, bool) {
	stdoutPath, stderrPath, mayBeRedirect := "", "", ""
	var symIndex, pathIndex int
	appendStdout, appendStderr := false, false

	for {
		if len(commandArguments) < 2 {
			break
		}
		symIndex = len(commandArguments) - 2
		pathIndex = len(commandArguments) - 1
		mayBeRedirect = commandArguments[symIndex]
		if strings.HasSuffix(mayBeRedirect, ">") {
			if strings.HasPrefix(mayBeRedirect, "2") {
				stderrPath = commandArguments[pathIndex]
				if strings.HasSuffix(mayBeRedirect, ">>") {
					appendStderr = true
				}
			}else{
				stdoutPath = commandArguments[pathIndex]
				if strings.HasSuffix(mayBeRedirect, ">>") {
					appendStdout = true
				}
			}
			commandArguments = commandArguments[0:(symIndex)]
		} else {
			break
		}
	}
	return commandArguments, stdoutPath, stderrPath, appendStdout, appendStderr
}

