package shell_test_context

import (
	"testing"
	"os/exec"
	"io"
	"bufio"
	"os"
	"strings"
	"fmt"
	"time"
)

type ShellTestContext struct {
	T *testing.T
	Cmd *exec.Cmd
	Stdin io.WriteCloser
	Stdout io.ReadCloser
	Stderr io.ReadCloser
	StdoutReader *bufio.Reader
	LogFile *os.File
}

func (c ShellTestContext) SendInput(input string) {
	_, err := c.Stdin.Write([]byte(input))
	if err != nil {
		c.T.Fatal(err)
	}
}

func (c ShellTestContext) ReadUntilPrompt() (string, error) {
	var result strings.Builder

	t := c.T
	reader := c.StdoutReader
	buf := make(chan byte, 1)
	done := make(chan struct{})
	c.log("About to start goroutine")
	go func(){
		c.log("Goroutine started")
		var b byte
		var err error
		defer close(done)
		for {
			b, _ = reader.ReadByte()
			//c.log(fmt.Sprintf("Inside goroutine reading byte: %q(%v)", b, b))
			if err != nil {
				c.log(fmt.Sprintf("Inside goroutine error: %v", err))
				t.Fatal(err)
				close(buf)
				break
			}
			if b == byte('%') {
				c.log("Inside goroutine closing the channel")
				reader.ReadByte() // read \n after %
				close(buf)
				break
			}else{
				buf <- b
			}
		}
	}()
		
	for {
		select {
		case anotherByte := <- buf:
			result.WriteByte(anotherByte)
		case <- time.After(220 * time.Millisecond):
			c.log(fmt.Sprintf("Timeout goroutine, received result: %v", result.String()))
			c.Stdin.Write([]byte("echo %\n")) // use % symbol to signal goroutine to stop reading and finis
			<- done
			//c.log(fmt.Sprintf("Done received: %v", result.String()))
			return result.String(), nil
		}
	}
}

func (c ShellTestContext) AssertCmd(input string, expectedOutput string) {
	c.SendInput(input)
	c.log("  AssertCmd: " + input)

	output, err := c.ReadUntilPrompt()

	if err != nil {
	    c.T.Fatalf("Failed to read initial prompt: %v", err)
	}

	trimmed := strings.Trim(output, "\n")

 	if trimmed != expectedOutput {
 		c.T.Errorf("----")
 		c.T.Errorf("Command: %s", input)
 		c.T.Errorf("expected output: %s", expectedOutput)
 		c.T.Errorf("Received output: %s", trimmed)
 		c.T.Errorf("----")
 	}
}

func (c ShellTestContext) TearDown() {
	c.SendInput("exit\n")

	c.LogFile.WriteString("< < <\n")
	c.Stdin.Close()
	c.LogFile.Close()

	stderrBytes, err := io.ReadAll(c.Stderr)
	stdoutBytes, err := io.ReadAll(c.Stdout)


	if err := c.Cmd.Wait(); err != nil {
		c.T.Fatalf("command failed: %s", err)
	}

	if err != nil {
		c.T.Fatal(err)
	}
	


	fmt.Printf("%v",string(stdoutBytes))
	if len(stderrBytes) > 0 {
		c.T.Logf("Stderr: %s", string(stderrBytes))
	}
}

func (c ShellTestContext) log(entry string){
	c.LogFile.WriteString(entry + "\n")
}
