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
	"context"
	"github.com/codecrafters-io/shell-starter-go/app/my_shell_history"
)

type Executable struct {
	name string
	builtIn bool
	path string
	executable func([]string, []byte) []byte
}

func nullExecutable([]string, []byte) []byte { return make([]byte, 0)}

func addBuiltIn(cmdMap map[string]Executable, name string){
	cmdMap[name] = Executable{
		name: name, 
		builtIn: true,
		path: "builtin",
	}
}


func LoadBinPaths(binExecutables *map[string]Executable) {
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
					executable: nullExecutable,
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
	StdoutPipe() (io.ReadCloser, error)
	Start() error
	Wait() error
}

func (ex Executable) BuildCmd(cmdArgs shell_args.ParsedArgs, ctx context.Context) (CmdInterface, error){
	var cmd CmdInterface
	var err error

	if ex.builtIn {
		cmd, err = CommandContext(ctx, ex.name, cmdArgs.Arguments...)
		if err != nil {
			return nil, err
		}
	} else {
		cmd = &ExecCmdWraper{exec.CommandContext(ctx, ex.name, cmdArgs.Arguments...)}
	}

	return cmd, nil
}

func PrepareRedirectFile(path string, apnd bool) (*os.File, error) {
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
	} else if fileInfo.Mode().IsRegular() && apnd {
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
	
	addBuiltIn(cmdMap, "exit")
	addBuiltIn(cmdMap, "echo")
	addBuiltIn(cmdMap, "type")
	addBuiltIn(cmdMap, "pwd")
	addBuiltIn(cmdMap, "cd")
	addBuiltIn(cmdMap, "history")

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
		if len(input) == 0 {
			continue
		}

		my_shell_history.StoreCommand(input)
		parsedInput := shell_args.ParseInput(input)

		err = runPipeline(parsedInput)
		if err != nil {
			fmt.Println(err)
			continue
		}
	}
}

func runPipeline(commandList []shell_args.ParsedArgs) error {
	//var err error
	stdin := os.Stdin
	stdout := os.Stdout
	stderr := os.Stderr
	bgCtx := context.Background()
	ctx, cancel := context.WithCancel(bgCtx)
	commandsCount := len(commandList)
	osCommands := make([]CmdInterface, commandsCount)

	for i, cmdArgs := range(commandList) {
		cmdName := cmdArgs.CommandName
		cmd, ok := cmdMap[cmdName]
		if ok != true {
			return errors.New(fmt.Sprintf("%v: command not found", cmdName))
		}
		stderr = stdout

		osCmd, err := cmd.BuildCmd(cmdArgs, ctx)
		if err != nil {
			return err
		}
		if commandsCount == 1 {
			if cmdArgs.IsStdoutRedirected() {
				stdout, err = PrepareRedirectFile(cmdArgs.StdoutPath, cmdArgs.AppendStdout)
				if err != nil {
					return err
				}
				defer stdout.Close()
			}

			if cmdArgs.IsStderrRedirected() {
				stderr, err = PrepareRedirectFile(cmdArgs.StderrPath, cmdArgs.AppendStderr)
				if err != nil {
					return err
				}
				defer stderr.Close()
			}
			osCmd.SetStdout(stdout)
			osCmd.SetStderr(stderr)
			osCmd.SetStdin(stdin)
		}else{
			if i > 0 {
				stdoutPipe, err := osCommands[i-1].StdoutPipe()
				if err != nil {
					fmt.Printf("Cmd %v stdout pipe error: %v\n", commandList[i-1].CommandName, err.Error())
					return err
				}
				osCmd.SetStdin(stdoutPipe)
				if i == commandsCount - 1 {
					osCmd.SetStdout(stdout)
				}
			}

		}

		osCommands[i] = osCmd
	}

	done := make(chan struct{})
	for i, osCmd := range osCommands {
		osCmd.Start()

		go func(){
			osCmd.Wait()
			if i == len(osCommands) - 1 {
				cancel()
				close(done)
			}
		}()

	}

	<- done
	return nil
}

