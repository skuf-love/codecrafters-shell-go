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

func StoreCommand(cmd string){
	log = append(log, cmd)
}

func Log() []string {
	return log
}

func ImportFromFile(pathArg string) error {
		path, err := filepath.Abs(pathArg)
		if err != nil {
			return errors.New(fmt.Sprintf("history: filepath error: %v\n", err))
		}

		_, err = os.Stat(path)
		if errors.Is(err, os.ErrNotExist) {
			return errors.New(fmt.Sprintf("history: file not exists: %v\n", err))
		}
		file, err := os.Open(path)

		if err != nil{
			return errors.New(fmt.Sprintf("history: error open file: %v\n", err))
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
		
		for _, cmd := range log{
			writer.WriteString(fmt.Sprintf("%v\n", cmd))
		}
		err = writer.Flush()
		if err != nil {
			return err
		}


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
