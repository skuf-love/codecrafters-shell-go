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
	output := make([]byte, 0)
	fmt.Println(strings.Join(cmdArgs.Arguments, " "))
	return output 
}

func typeExecutable(cmdArgs shell_args.ParsedArgs) []byte {
	output := make([]byte, 0)
	cmd, ok := cmdMap[cmdArgs.Arguments[0]]
	if ok {
		if cmd.builtIn {
			fmt.Println(cmd.name + " is a shell builtin")
		} else {
			fmt.Println(cmd.name + " is " + cmd.path)
		}
		return output
	}
	fmt.Println(cmdArgs.Arguments[0] + ": not found")
	return output 
}
func cdExecutable(cmdArgs shell_args.ParsedArgs) []byte{
	output := make([]byte, 0)
	path := cmdArgs.Arguments[0]
	if path == "~" {
		path = os.Getenv("HOME")
	}
	stat, err := os.Stat(path)
	if err != nil {
		fmt.Printf("cd: %v: No such file or directory\n" ,path)
		// fmt.Printf("%v\n", err)
		return output
	}
	if !stat.IsDir() {
		fmt.Printf("%v\n", err)
		return output
	}
	err = os.Chdir(path)
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	return output
}
func pwdExecutable(shell_args.ParsedArgs)  []byte{
	output := make([]byte, 0)
	wd, err := os.Getwd()
	if err != nil {
		fmt.Printf("%v", err)
	}
	fmt.Println(wd) 
	return output
}

