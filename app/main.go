package main

import (
	"fmt"
	"bufio"
	"os"
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
		fmt.Println(trimmed_command + ": command not found")
	}
}
