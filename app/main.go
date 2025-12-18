package main

import (
	"fmt"
	"bufio"
	"os"
	"strings"
	"os/exec"
	"github.com/codecrafters-io/shell-starter-go/app/shell_args"
	"bytes"
	"path/filepath"
	"errors"
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
	var stdoutBuf, stderrBuf bytes.Buffer
	var err error
	if ex.builtIn {
		stdoutBuf.Write(ex.executable(cmdArgs))
	} else {
		cmd := exec.Command(ex.name, cmdArgs.Arguments...)

		cmd.Stdout = &stdoutBuf
		cmd.Stderr = &stderrBuf

		err = cmd.Run()

	}
	stdout := stdoutBuf.Bytes()
	stderr := stderrBuf.Bytes()

	if cmdArgs.IsStdoutRedirected() {
		file, err := PrepareRedirectFile(cmdArgs.StdoutPath, cmdArgs.AppendStdout)
		if err != nil {
			fmt.Printf("%v\n", err)
			return
		}
		defer file.Close()
		_, err = file.Write(stdout)

		if err != nil {
			fmt.Printf("%v\n", err)
			return
		}
	}else{
		fmt.Printf("%s", string(stdout))
	}
	if cmdArgs.IsStderrRedirected() {
		errFile, err := PrepareRedirectFile(cmdArgs.StderrPath, cmdArgs.AppendStderr)
		if err != nil {
			fmt.Printf("%v\n", err)
			return
		}
		defer errFile.Close()
		_, err = errFile.Write(stderr)

		if err != nil {
			fmt.Printf("%v\n", err)
			return
		}
	}else{
		if err != nil {
			fmt.Printf("%v", string(stderr))
		}
	}
}

func PrepareRedirectFile(path string, append bool) (*os.File, error) {
	//fmt.Printf("preparing %v\n", path)
	path, err := filepath.Abs(path)
	//defer fmt.Printf("prepared %v\n", path)
	if err != nil {
		return nil, err
	}


	fileInfo, err := os.Stat(path)
	if errors.Is(err, os.ErrNotExist) {
		return os.Create(path)
	}else if err != nil {
		return nil, err
	}

	if fileInfo.Mode().IsDir() {
		return nil, fmt.Errorf("Path %q is a directory", path)
	} else if fileInfo.Mode().IsRegular() && append {
		//fmt.Printf("Appending %v\n", path)
		return os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0644)
	}else{
		return os.Create(path)
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

