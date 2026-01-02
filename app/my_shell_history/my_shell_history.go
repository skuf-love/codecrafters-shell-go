package my_shell_history

import(
	"path/filepath"
	"fmt"
	"os"
	"errors"
	"io"
	"bufio"
)

var log []string
var logSinceAppend []string

func Init(){
	if path, varDefined := os.LookupEnv("HISTFILE"); varDefined {
		err := ImportFromFile(path)
		if err != nil {
			fmt.Printf("Failed to read history from file: %v\n", err)
		}
	}
}

func Dump(){
	if path, varDefined := os.LookupEnv("HISTFILE"); varDefined {
		err := ExportToFile(path, false)
		if err != nil {
			fmt.Printf("Failed to export history to file: %v\n", err)
		}
	}
}
func StoreCommand(cmd string){
	log = append(log, cmd)
	logSinceAppend = append(logSinceAppend , cmd)
}

func Log() []string {
	return log
}

func ImportFromFile(pathArg string) error {
		path, err := filepath.Abs(pathArg)
		if err != nil {
			return errors.New(fmt.Sprintf("history: filepath error: %v", err))
		}

		_, err = os.Stat(path)
		if errors.Is(err, os.ErrNotExist) {
			return errors.New(fmt.Sprintf("history: file not exists: %v", err))
		}
		file, err := os.Open(path)

		if err != nil{
			return errors.New(fmt.Sprintf("history: error open file: %v", err))
		}
		defer file.Close()

		reader := bufio.NewReader(file)
		for {
			line , _,err := reader.ReadLine()
			if err != nil {
				if errors.Is(err, io.EOF) {
					return nil
				}else{
					return err
				}
			}
			if len(line) > 0 {
				StoreCommand(string(line))
			}
		}

	return nil
}

func ExportToFile(path string, apnd bool) error {
		file, err := prepareExportFile(path, apnd)
		if err != nil {
			return err
		}
		defer file.Close()

		writer := bufio.NewWriter(file)

		var commandsToWrite []string

		if apnd {
			commandsToWrite = logSinceAppend
		}else{
			commandsToWrite = log
		}
		
		for _, cmd := range commandsToWrite{
			writer.WriteString(fmt.Sprintf("%v\n", cmd))
		}
		err = writer.Flush()
		if err != nil {
			return err
		}
		logSinceAppend = make([]string, 0)

	return nil
}

func prepareExportFile(path string, apnd bool) (*os.File, error) {
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
