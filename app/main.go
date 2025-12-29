package main

import (
	"fmt"
	"os"
	"strings"
	"os/exec"
	"github.com/codecrafters-io/shell-starter-go/app/shell_args"
	"path/filepath"
	"errors"
	"github.com/chzyer/readline"
	"maps"
	"slices"
	"github.com/codecrafters-io/shell-starter-go/app/custom_prefix_completer"
	"io"
)


type Executable struct {
	name string
	builtIn bool
	path string
	executable func([]string) []byte
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
				(*binExecutables)[dirEntry.Name()] = Executable{
					name:	dirEntry.Name(),
					builtIn: false,
					path: binPath,
					executable: func([]string) []byte { return make([]byte, 0)},
				}
			}
		}
	}
	
}
var cmdMap map[string]Executable

type CmdInterface interface{
	Run() error
	SetStdin(io.Reader)
	SetStdout(io.Writer)
	SetStderr(io.Writer)
}

func (ex Executable) Run(cmdArgs shell_args.ParsedArgs, stdin io.Reader, stdout io.Writer, stderr io.Writer){
	//var err error
	var cmd CmdInterface
	if ex.builtIn {
		cmd = Command(ex.name, cmdArgs.Arguments...)
	} else {
		cmd = &ExecCmdWraper{exec.Command(ex.name, cmdArgs.Arguments...)}
	}

	cmd.SetStdin(stdin)
	cmd.SetStdout(stdout)
	cmd.SetStderr(stderr)

	cmd.Run()
	//if err != nil {
	//	fmt.Println(err)
	//}

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
	completer := readline.NewPrefixCompleter(execCompleters...)
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
		input, err := rl.Readline()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading input:", err)
			break
		}

		parsedInput := shell_args.ParseInput(input)

		err = runPipeline(parsedInput)
		if err != nil {
			fmt.Printf("%q\n", err)
			continue
		}

	}
}

func runPipeline(commandList []shell_args.ParsedArgs) error {
	var err error
	stdin := os.Stdin
	stdout := os.Stdout
	stderr := os.Stderr
	for _, cmdArgs := range(commandList) {
		cmdName := cmdArgs.CommandName
		cmd, ok := cmdMap[cmdName]
		if ok != true {
			return errors.New(fmt.Sprintf("%v: command not found", cmdName))
		}
		stderr = stdout
		if cmdArgs.IsStdoutRedirected() {
			stdout, err = PrepareRedirectFile(cmdArgs.StdoutPath, cmdArgs.AppendStdout)
			if err != nil {
				return err
			}
		}

		if cmdArgs.IsStderrRedirected() {
			stderr, err = PrepareRedirectFile(cmdArgs.StderrPath, cmdArgs.AppendStderr)
			if err != nil {
				return err
			}
		}
	
		cmd.Run(cmdArgs, stdin, stdout, stderr)

		if cmdArgs.IsStdoutRedirected() {
			err = stdout.Close()
			if err != nil {
				return err
			}
		}

		if cmdArgs.IsStderrRedirected() {
			err = stderr.Close()
			if err != nil {
				return err
			}
		}
		stdout = os.Stdout
		stderr = os.Stderr
	}
	return nil
}

