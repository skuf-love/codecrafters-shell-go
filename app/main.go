package main

import (
	"fmt"
	"bufio"
	"os"
	"strings"
	"os/exec"
)


type Cmd struct {
	name string
	builtIn bool
	path string
}

func exitCmd() {
	os.Exit(0)
}
func echoCmd(cmdArgs []string) {
	fmt.Println(strings.Join(cmdArgs, " "))
}

func typeCmd(cmdArgs []string) {
	cmd, ok := cmdMap[cmdArgs[0]]
	if ok {
		if cmd.builtIn {
			fmt.Println(cmd.name + " is a shell builtin")
		} else {
			fmt.Println(cmd.name + " is " + cmd.path)
		}
		return
	}
	fmt.Println(cmdArgs[0] + ": not found")
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
	cmd := exec.Command(ex.path, args...)

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

	for {
		fmt.Print("$ ")
		input, read_err := bufio.NewReader(os.Stdin).ReadString('\n')
		if read_err != nil {
			fmt.Fprintln(os.Stderr, "Error reading input:", read_err)
			os.Exit(1)
		}

		trimmed_input := input[:len(input)-1]
		split_input := strings.Split(trimmed_input, " ")

		cmd_name := split_input[0]
		args := split_input[1:len(split_input)]

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

		cmd.Run(args)

	}
}

