package main

import (
	"fmt"
	"bufio"
	"os"
	"strings"
	"os/exec"
	"github.com/codecrafters-io/shell-starter-go/app/shell_args"
)


type Executable struct {
	name string
	builtIn bool
	path string
	executable func(shell_args.ParsedArgs) []byte
}


func LoadBinPaths(binExecutables *map[string]Executable)  {
	pathVar := os.Getenv("PATH")

	paths := strings.Split(pathVar, ":")
		
	for _, path := range paths {

		dirEntries, err := os.ReadDir(path)
		if err != nil {
			continue
		}

		for _, dirEntry := range dirEntries {
			binPath := path + "/" + dirEntry.Name()


			fileInfo, err := os.Stat(binPath)
			if err != nil {
				continue
			}
			if fileInfo.IsDir() {
				continue
			}
			if _, ok := (*binExecutables)[dirEntry.Name()]; ok {
				continue
			}
			mode := fileInfo.Mode()
			if mode.IsRegular() && mode.Perm()&0111 != 0 {
				(*binExecutables)[dirEntry.Name()] = Executable{dirEntry.Name(), false, binPath, func(shell_args.ParsedArgs) []byte { return make([]byte, 0)},}
			}
		}
	}
	
}
var cmdMap map[string]Executable

func (ex Executable) Run(cmdArgs shell_args.ParsedArgs){
	output := make([]byte, 0)
	var err error
	if ex.builtIn {
		builtinOutput := ex.executable(cmdArgs)
		if builtinOutput != nil {
			output = builtinOutput
		}
		fmt.Printf("%s", string(output))
	} else {
		cmd := exec.Command(ex.name, cmdArgs.Arguments...)

		output, err = cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
		fmt.Printf("%s", string(output))
		
	}
}
func main() {

	cmdMap = make(map[string]Executable)
	LoadBinPaths(&cmdMap)

	cmdMap["exit"] = Executable{"exit", true, "builtin", exitExecutable,}
	cmdMap["echo"] = Executable{"echo",  true, "builtin", echoExecutable,}
	cmdMap["type"] = Executable{"type", true, "builtin", typeExecutable,}
	cmdMap["pwd"] = Executable{"pwd", true, "builtin", pwdExecutable,}
	cmdMap["cd"] = Executable{"cd", true, "builtin", cdExecutable,}

	for {
		fmt.Print("$ ")
		input, read_err := bufio.NewReader(os.Stdin).ReadString('\n')
		if read_err != nil {
			fmt.Fprintln(os.Stderr, "Error reading input:", read_err)
			os.Exit(1)
		}

		parsedInput := shell_args.ParseInput(input)

		cmd_name := parsedInput.CommandName

		cmd, cmd_map_ok := cmdMap[cmd_name]
		
		if cmd_map_ok != true {
			fmt.Println(cmd_name + ": command not found")
			continue
		}


		cmd.Run(parsedInput)

	}
}

