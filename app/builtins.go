package main

import (
	"fmt"
	"os"
	"strings"
	"io"
	"context"
	"bufio"
	"github.com/codecrafters-io/shell-starter-go/app/my_shell_history"
)

type Cmd struct{
	name string
	Stdin io.Reader
	Stdout io.Writer
	Stderr io.Writer
	executable func([]string, []byte) []byte
	args []string
	StdoutWritePipe *io.PipeWriter
	StdoutReadPipe *io.PipeReader
	ctx context.Context
	done chan struct{}
}

func CommandContext(ctx context.Context, name string, args ...string) (*Cmd, error) {
	cmd, err := Command(name, args...)
	if err != nil {
		return nil, err
	}
	cmd.ctx = ctx
	return cmd, nil
}
func Command(name string, args ...string) (*Cmd, error){
	var executable func([]string, []byte) []byte
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
	case "history":
		executable = historyExecutable
	default:
		return nil, fmt.Errorf("%b buildIn not defined", name)
	}

	rp, wp := io.Pipe()
	return &Cmd {
		name: name,
		executable: executable,
		args: args,
		StdoutWritePipe: wp,
		StdoutReadPipe: rp,
	}, nil
}

func (cmd *Cmd) Run() error{
	result := cmd.executable(cmd.args, make([]byte, 0))
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
	if c.Stdout != nil {
		return nil, fmt.Errorf("Stdout is already set")
	}
	c.Stdout = c.StdoutWritePipe
	return c.StdoutReadPipe, nil
}

func (c *Cmd) Start() error{
	c.done = make(chan struct{})
	go func(){
		//read stdin 
		reader := bufio.NewReader(c.Stdin)
		stdin := make([]byte, reader.Buffered())
		_, err := reader.Read(stdin)
		if err != nil {
			close(c.done)
			return
		}
//		for  {
//			inbyte, err := reader.ReadByte()
//			if err != nil {
//				break
//			}
//			stdin = append(stdin, inbyte)
//		}
		result := c.executable(c.args, stdin)
		//fmt.Printf("Write result %v\n", result)
		c.Stdout.Write(result)
		c.StdoutWritePipe.Close()
		close(c.done)
	}()
	return nil
}
func (c *Cmd) Wait() error{

	//fmt.Println("before c.done")
	<- c.done
	//fmt.Println("after c.done")
	return nil
}

func exitExecutable([]string, []byte) []byte{
	os.Exit(0)
	return make([]byte, 0)
}
func echoExecutable(args []string, stdin []byte) []byte{
	output := fmt.Sprintln(strings.Join(args, " "))
	return []byte(output)
}

func typeExecutable(args []string, stdin []byte) []byte {
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
func cdExecutable(args []string, stdin []byte) []byte{
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
func pwdExecutable(args []string, stdin []byte)  []byte{
	output := ""
	wd, err := os.Getwd()
	if err != nil {
		output = fmt.Sprintf("%v", err)
	}
	output = fmt.Sprintln(wd) 
	return []byte(output)
}

func historyExecutable([]string, []byte) []byte{
	history := ""
	for i, cmdLine := range my_shell_history.Log() {
		history += fmt.Sprintf("    %v  %v\n", i+1, cmdLine)
	}
	return []byte(history)
}

