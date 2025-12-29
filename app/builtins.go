package main

import (
	"fmt"
	"os"
	"strings"
	"io"
)

type Cmd struct{
	name string
	Stdin io.Reader
	Stdout io.Writer
	Stderr io.Writer
	executable func([]string) []byte
	args []string
	writePipe *io.PipeWriter
	readPipe *io.PipeReader
}

func Command(name string, args ...string) *Cmd {
	var executable func([]string) []byte
	switch name {
	case "exit":
		executable = exitExecutable
	case "echo":
		executable = echoExecutable
	case "type":
		executable = typeExecutable
	case "cd":
		executable = cdExecutable
	case "pwd":
		executable = pwdExecutable
	}

	rp, wp := io.Pipe()
	return &Cmd {
		name: name,
		executable: executable,
		args: args,
		writePipe: wp,
		readPipe: rp,
	}
}

func (cmd *Cmd) Run() error{
	result := cmd.executable(cmd.args)
	cmd.Stdout.Write(result)
	return nil
	
}

func (c *Cmd) SetStdin(stdin io.Reader) {
	c.Stdin = stdin
}

func (c *Cmd) SetStdout(stdout io.Writer) {
	c.Stdout = stdout
}

func (c *Cmd) SetStderr(stderr io.Writer) {
	c.Stderr = stderr
}

func (c *Cmd) StdoutPipe() (io.ReadCloser, error){
	return c.readPipe, nil
}

func (c *Cmd) Start() error{
	return c.Run()
}
func (c *Cmd) Wait() error{
	return nil
}

func exitExecutable([]string) []byte{
	os.Exit(0)
	return make([]byte, 0)
}
func echoExecutable(args []string) []byte{
	output := fmt.Sprintln(strings.Join(args, " "))
	return []byte(output)
}

func typeExecutable(args []string) []byte {
	output := ""
	cmd, ok := cmdMap[args[0]]
	if ok {
		if cmd.builtIn {
			output = fmt.Sprintln(cmd.name + " is a shell builtin")
		} else {
			output = fmt.Sprintln(cmd.name + " is " + cmd.path)
		}
	}else{
		output = fmt.Sprintln(args[0] + ": not found")
	}
	return []byte(output)
}
func cdExecutable(args []string) []byte{
	output := ""
	path := args[0]
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
func pwdExecutable(args []string)  []byte{
	output := ""
	wd, err := os.Getwd()
	if err != nil {
		output = fmt.Sprintf("%v", err)
	}
	output = fmt.Sprintln(wd) 
	return []byte(output)
}

