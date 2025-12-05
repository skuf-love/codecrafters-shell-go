package main

import (
	"fmt"
	"bufio"
	"os"
	"strings"
)


type Cmd struct {
	name string
	builtIn bool
}

var cmdMap = map[string]Cmd{
	"exit":  Cmd{"exit", true,},
	"echo":  Cmd{"echo",  true,},
	"type":  Cmd{"type", true,},
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
		}
	} else {
		
		fmt.Println(cmdArgs[0] + ": not found")
	}

}
func main() {
	for {
		fmt.Print("$ ")
		input, read_err := bufio.NewReader(os.Stdin).ReadString('\n')
		if read_err != nil {
			fmt.Println(os.Stderr, "Error reading input:", read_err)
			os.Exit(1)
		}

		trimmed_input := input[:len(input)-1]
		split_input := strings.Split(trimmed_input, " ")

		cmd_name := split_input[0]
		args := split_input[1:len(split_input)]

		cmd, cmd_map_ok := cmdMap[cmd_name]
		
		if cmd_map_ok != true {
			fmt.Println(cmd_name + ": command not found")
			os.Exit(0)
		}

		if cmd.name == "exit"{
			exitCmd()
		}

		if cmd.name == "echo" {
			echoCmd(args)
			continue
		}

		if cmd.name == "type" {
			typeCmd(args)
			continue
		}

	}
}
