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
				(*binExecutables)[dirEntry.Name()] = Executable{dirEntry.Name(), false, binPath}
			}
		}
	}
	
}
var cmdMap map[string]Executable

func (ex Executable) Run(args []string){
	cmd := exec.Command(ex.name, args...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	
	fmt.Printf("%s", string(output))
}
func main() {

	cmdMap = make(map[string]Executable)
	LoadBinPaths(&cmdMap)

	cmdMap["exit"] = Executable{"exit", true, "builtin",}
	cmdMap["echo"] = Executable{"echo",  true, "builtin",}
	cmdMap["type"] = Executable{"type", true, "builtin",}
	cmdMap["pwd"] = Executable{"pwd", true, "builtin",}
	cmdMap["cd"] = Executable{"cd", true, "builtin",}

	for {
		fmt.Print("$ ")
		input, read_err := bufio.NewReader(os.Stdin).ReadString('\n')
		if read_err != nil {
			fmt.Fprintln(os.Stderr, "Error reading input:", read_err)
			os.Exit(1)
		}

		parsedInput := shell_args.ParseInput(input)

		cmd_name := parsedInput.CommandName
		args := parsedInput.Arguments

		//if len(args) > 1 && strings.HasSuffix(args[len(args)-2], ">") {
		//	fmt.Printf("Args: %v\n", args)
		//	stdout_path := args[len(args)-1]
		//	fmt.Printf("Args: %v\n", stdout_path)
		//	stdout, _ := os.Create(stdout_path)

		//	stdout.WriteString("test file write")
		//	stdout.Close()
		//	args = args[0:len(args)-2]
		//	fmt.Printf("Args withot stdout redirect : %v\n", args)
		//}

		cmd, cmd_map_ok := cmdMap[cmd_name]
		
		if cmd_map_ok != true {
			fmt.Println(cmd_name + ": command not found")
			continue
		}

		if cmd.name == "exit"{
			exitExecutable()
		}

		if cmd.name == "echo" {
			echoExecutable(args)
			continue
		}

		if cmd.name == "type" {
			typeExecutable(args)
			continue
		}
		if cmd.name == "pwd" {
			pwdExecutable()
			continue
		}
		if cmd.name == "cd" {
			cdExecutable(args[0])
			continue
		}

		cmd.Run(args)

	}
}

