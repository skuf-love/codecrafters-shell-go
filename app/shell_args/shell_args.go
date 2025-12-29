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
	allCommands []ParsedArgs
	parsedArgs ParsedArgs
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
		c.mode = "single_quote"
		return
	}
	if char == '"' {
		c.mode = "double_quote"
		return
	}
	if char == ' ' && len(c.currentArg) == 0 {
		return 
	}
	if char == ' ' {
		c.args = append(c.args, string(c.currentArg))
		c.currentArg = make([]rune, 0)
	} else if char == '|' {
		c.flushCommand()
	}else{
		c.currentArg = append(c.currentArg, char)
	}
}
func (c *parseContext) singleQuoteRead(char rune) {
	if char == '\'' {
		c.mode = "normal"
		return
	}
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

func (c *parseContext) flushCommand(){
	if len(c.currentArg) > 0 {
		c.args = append(c.args, string(c.currentArg))
	}

	commandName := c.args[0]
	commandArguments := make([]string, 0)
	if len(c.args) > 1 {
		commandArguments = c.args[1:]
	}
	commandArguments, stdoutPath, stderrPath, appendStdout, appendStderr := parseRedirects(commandArguments)
	c.parsedArgs.CommandName = commandName
	c.parsedArgs.Arguments = commandArguments
	c.parsedArgs.StdoutPath = stdoutPath
	c.parsedArgs.StderrPath = stderrPath
	c.parsedArgs.AppendStdout = appendStdout
	c.parsedArgs.AppendStderr = appendStderr

	c.allCommands = append(c.allCommands, c.parsedArgs)

	c.parsedArgs = newParsedArgs()
	c.currentArg = make([]rune, 0)
	c.args = make([]string, 0)
	
}

func newParsedArgs() ParsedArgs{
	return ParsedArgs{
		"",
		make([]string, 0),
		"",
		"",
		false,
		false,
	}
}
func ParseInput(input string) []ParsedArgs {
	trimmed_input, _ := strings.CutSuffix(input, "\n")
	//fmt.Println(trimmed_input)
	
	parsedArgs := newParsedArgs()
	context := parseContext{
		make([]rune, 0),
		make([]string, 0),
		"normal",
		false,
		make([]ParsedArgs, 0),
		parsedArgs,
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
	
	context.flushCommand()

	//context.args = append(context.args, string(context.currentArg))

	//commandName := context.args[0]
	//commandArguments := make([]string, 0)
	//if len(context.args) > 1 {
	//	commandArguments = context.args[1:]
	//}
	//commandArguments, stdoutPath, stderrPath, appendStdout, appendStderr := parseRedirects(commandArguments)

	//return ParsedArgs{
	//	commandName,
	//	commandArguments,
	//	stdoutPath,
	//	stderrPath,
	//	appendStdout,
	//	appendStderr,
	//} 
	return context.allCommands
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
				appendStderr = strings.HasSuffix(mayBeRedirect, ">>")
			}else{
				stdoutPath = commandArguments[pathIndex]
				appendStdout = strings.HasSuffix(mayBeRedirect, ">>")
			}
			commandArguments = commandArguments[0:(symIndex)]
		} else {
			break
		}
	}
	return commandArguments, stdoutPath, stderrPath, appendStdout, appendStderr
}

