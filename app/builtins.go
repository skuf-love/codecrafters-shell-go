package main

import (
	"fmt"
	"os"
	"strings"
	"github.com/codecrafters-io/shell-starter-go/app/shell_args"
)

func exitExecutable(shell_args.ParsedArgs) {
	os.Exit(0)
}
func echoExecutable(cmdArgs shell_args.ParsedArgs) {
	fmt.Println(strings.Join(cmdArgs.Arguments, " "))
}

func typeExecutable(cmdArgs shell_args.ParsedArgs) {
	cmd, ok := cmdMap[cmdArgs.Arguments[0]]
	if ok {
		if cmd.builtIn {
			fmt.Println(cmd.name + " is a shell builtin")
		} else {
			fmt.Println(cmd.name + " is " + cmd.path)
		}
		return
	}
	fmt.Println(cmdArgs.Arguments[0] + ": not found")
}
func cdExecutable(cmdArgs shell_args.ParsedArgs) {
	path := cmdArgs.Arguments[0]
	if path == "~" {
		path = os.Getenv("HOME")
	}
	stat, err := os.Stat(path)
	if err != nil {
		fmt.Printf("cd: %v: No such file or directory\n" ,path)
		// fmt.Printf("%v\n", err)
		return
	}
	if !stat.IsDir() {
		fmt.Printf("%v\n", err)
		return
	}
	err = os.Chdir(path)
	if err != nil {
		fmt.Printf("%v\n", err)
	}
}
func pwdExecutable(shell_args.ParsedArgs) {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Printf("%v", err)
	}
	fmt.Println(wd) 
}

