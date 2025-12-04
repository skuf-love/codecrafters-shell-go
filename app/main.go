package main

import (
	"fmt"
	"bufio"
	"os"
	"strings"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Print

func main() {
	for {
		fmt.Print("$ ")
		command, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			fmt.Println(os.Stderr, "Error reading input:", err)
			os.Exit(1)
		}
		trimmed_command := command[:len(command)-1]
		if trimmed_command == "exit"{
			os.Exit(0)
		}
		echo_cmd := "echo"
		
		cmd := strings.Split(trimmed_command, " ")
		if cmd[0] == echo_cmd {
			args := cmd[1:len(cmd)]
			fmt.Println(strings.Join(args, " "))
			continue
		}
		fmt.Println(trimmed_command + ": command not found")
	}
}
