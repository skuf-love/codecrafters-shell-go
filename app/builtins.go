package main

import (
	"fmt"
	"os"
	"strings"
)

func exitExecutable() {
	os.Exit(0)
}
func echoExecutable(cmdArgs []string) {
	fmt.Println(strings.Join(cmdArgs, " "))
}

func typeExecutable(cmdArgs []string) {
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
func cdExecutable(path string) {
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
func pwdExecutable() {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Printf("%v", err)
	}
	fmt.Println(wd) 
}

