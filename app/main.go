package main

import (
	"fmt"
	"os"
	"strings"
	"os/exec"
	"github.com/codecrafters-io/shell-starter-go/app/shell_args"
	"bytes"
	"path/filepath"
	"errors"
	"github.com/chzyer/readline"
	"maps"
	"slices"
	"github.com/codecrafters-io/shell-starter-go/app/custom_prefix_completer"
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
		DumpStream(cmdArgs.StdoutPath, cmdArgs.AppendStdout, stdout)
	}else{
		fmt.Printf("%s", string(stdout))
	}

	if cmdArgs.IsStderrRedirected() {
		DumpStream(cmdArgs.StderrPath, cmdArgs.AppendStderr, stderr)
	}else{
		if err != nil {
			fmt.Printf("%v", string(stderr))
		}
	}
}

func DumpStream(destPath string, doAppend bool, buffer []byte) {
	file, err := PrepareRedirectFile(destPath, doAppend)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	defer file.Close()
	_, err = file.Write(buffer)

	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
}

func PrepareRedirectFile(path string, append bool) (*os.File, error) {
	path, err := filepath.Abs(path)
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
		return os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0644)
	}else{
		return os.Create(path)
	}

}

func notFound(string) []string {
	fmt.Print("\x07")
	return make([]string, 0)
}
//var completer = readline.NewPrefixCompleter(
//	readline.PcItem("echo"),
//	readline.PcItem("exit"),
//	readline.PcItemDynamic(notFound),
//)

func PcItemsFromCmds(cmdMap map[string]Executable) []readline.PrefixCompleterInterface {
	cmdNames := slices.Sorted(maps.Keys(cmdMap))
	cmdCount := len(cmdNames)
	completers := make([]readline.PrefixCompleterInterface, cmdCount)

	for i, cmdName := range cmdNames {
		completers[i] = readline.PcItem(cmdName)
	}
	return completers
}

func main() {
	//var termios syscall.Termios
	//fd := int(os.Stdout.Fd())
	//r1, _, serr := syscall.Syscall6(
	//	syscall.SYS_IOCTL,
	//	uintptr(fd),
	//	syscall.TIOCGETA,  // Попытка получить настройки терминала
	//	uintptr(unsafe.Pointer(&termios)),
        //0, 0, 0,
    	//)
	//fmt.Printf("syscall.Syscall6 r1: %v", r1) 
	//if serr == 0 {
	//	fmt.Println("This is a TTY!")
	//	fmt.Printf("Termios: %+v\n", termios)
	//    } else {
	//	fmt.Printf("Not a TTY, error: %v\n", serr)
	//    }

	cmdMap = make(map[string]Executable)
	LoadBinPaths(&cmdMap)


	cmdMap["exit"] = Executable{"exit", true, "builtin", exitExecutable,}
	cmdMap["echo"] = Executable{"echo",  true, "builtin", echoExecutable,}
	cmdMap["type"] = Executable{"type", true, "builtin", typeExecutable,}
	cmdMap["pwd"] = Executable{"pwd", true, "builtin", pwdExecutable,}
	cmdMap["cd"] = Executable{"cd", true, "builtin", cdExecutable,}

	execCompleters := PcItemsFromCmds(cmdMap)
	//execCompleters = append(execCompleters, readline.PcItemDynamic(notFound))
	completer := readline.NewPrefixCompleter(execCompleters...)
	//customCompleter := &CustomPrefixCompleter{completer, 0, make([][]rune,0), int(0), make([]rune, 0), "$ "}
	customCompleter := custom_prefix_completer.New(completer, "$ ")

	rl, err := readline.NewEx(&readline.Config{
		Prompt: "$ ",
		AutoComplete: &customCompleter,
	})
	if err != nil {
		panic(err)
	}
	defer rl.Close()

	for {
		//fmt.Print("$ ")
		//input, read_err := bufio.NewReader(os.Stdin).ReadString('\n')
		input, err := rl.Readline()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading input:", err)
			//os.Exit(1)
			break
		}

		parsedInput := shell_args.ParseInput(input)[0]

		cmd_name := parsedInput.CommandName

		cmd, cmd_map_ok := cmdMap[cmd_name]
		
		if cmd_map_ok != true {
			fmt.Println(cmd_name + ": command not found")
			continue
		}


		cmd.Run(parsedInput)

	}
}

