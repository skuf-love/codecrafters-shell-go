package main

import (
	"fmt"
	"os"
	"strings"
	"github.com/codecrafters-io/shell-starter-go/app/shell_args"
)

func exitExecutable(shell_args.ParsedArgs) []byte{
	os.Exit(0)
	return make([]byte, 0)
}
func echoExecutable(cmdArgs shell_args.ParsedArgs) []byte{
	output := fmt.Sprintln(strings.Join(cmdArgs.Arguments, " "))
	return []byte(output)
}

func typeExecutable(cmdArgs shell_args.ParsedArgs) []byte {
	output := ""
	cmd, ok := cmdMap[cmdArgs.Arguments[0]]
	if ok {
		if cmd.builtIn {
			output = fmt.Sprintln(cmd.name + " is a shell builtin")
		} else {
			output = fmt.Sprintln(cmd.name + " is " + cmd.path)
		}
	}else{
		output = fmt.Sprintln(cmdArgs.Arguments[0] + ": not found")
	}
	return []byte(output)
}
func cdExecutable(cmdArgs shell_args.ParsedArgs) []byte{
	output := ""
	path := cmdArgs.Arguments[0]
	if path == "~" {
		path = os.Getenv("HOME")
	}
	stat, err := os.Stat(path)
	if err != nil {
		output = fmt.Sprintf("cd: %v: No such file or directory\n" ,path)
		return []byte(output)
	}
	if !stat.IsDir() {
		output = fmt.Sprintf("%v\n", err)
		return []byte(output)
	}
	err = os.Chdir(path)
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	return []byte(output)
}
func pwdExecutable(shell_args.ParsedArgs)  []byte{
	output := ""
	wd, err := os.Getwd()
	if err != nil {
		output = fmt.Sprintf("%v", err)
	}
	output = fmt.Sprintln(wd) 
	return []byte(output)
}

