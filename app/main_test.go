package main

import(
	"os"
	"os/exec"
	"strings"
	"testing"
	"fmt"
	"io"
	"bufio"
)


func TestMain(m *testing.M) {
	err := exec.Command("go", "build", "-o", "testapp").Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failded to build app: %v\n", err)
		os.Exit(1)
	}

	exitCode := m.Run()
	os.Remove("testapp")
	
	os.Exit(exitCode)
}
func SetupPipes(t *testing.T, cmd *exec.Cmd) (stdin io.WriteCloser, stdout, stderr io.ReadCloser){
	fmt.Println("Start pipe setups")
 	stdin, err := cmd.StdinPipe()
 	if err != nil {
 		t.Fatal(err)
 	}
	
	stdout, err = cmd.StdoutPipe()
	if err != nil {
		t.Fatal(err)
	}
	
	stderr, err = cmd.StderrPipe()
	if err != nil {
		t.Fatal(err)
	}
	return stdin, stdout, stderr
}

func sendInput(input string, t *testing.T, stdin io.WriteCloser) {
	_, err := stdin.Write([]byte(input))
	if err != nil {
		t.Fatal(err)
	}
}

func eatInitPrompt(reader *bufio.Reader, t *testing.T) {
	_, err := reader.ReadString('$')
	if err != nil {
 		t.Fatal(err)
	}
}
func readUntilPrompt(reader *bufio.Reader, t *testing.T) (string, error) {
	
	line, err := reader.ReadString('$')
	if err == io.EOF {
		return "", err
	}

	line, _ = strings.CutSuffix(line, "\n$")
	return line, nil
}

func (c ShellTestContext) assertCmd(input string, expectedOutput string) {
	sendInput(input, c.t, c.stdin)

	output, err := readUntilPrompt(c.stdoutReader, c.t)

	if err != nil {
	    c.t.Fatalf("Failed to read initial prompt: %v", err)
	}

 	if strings.Trim(output, " ") != expectedOutput {
 		c.t.Errorf("----")
 		c.t.Errorf("Command: %s", input)
 		c.t.Errorf("expected output: %s", expectedOutput)
 		c.t.Errorf("Received output: %s", output)
 		c.t.Errorf("----")
 	}
}

func (c ShellTestContext) tearDown() {
	sendInput("exit\n", c.t, c.stdin)

	c.stdin.Close()

	stderrBytes, err := io.ReadAll(c.stderr)
	stdoutBytes, err := io.ReadAll(c.stdout)

	if err := c.cmd.Wait(); err != nil {
		c.t.Fatalf("command failed: %s", err)
	}

	if err != nil {
		c.t.Fatal(err)
	}


	if err != nil {
		c.t.Fatal(err)
	}
	
	fmt.Printf("%v",string(stdoutBytes))


	if len(stderrBytes) > 0 {
		c.t.Logf("Stderr: %s", string(stderrBytes))
	}
}

type ShellTestContext struct {
	t *testing.T
	cmd *exec.Cmd
	stdin io.WriteCloser
	stdout io.ReadCloser
	stderr io.ReadCloser
	stdoutReader *bufio.Reader
}
func InitTest(t *testing.T) ShellTestContext {
	cmd := exec.Command("./testapp")
	stdin, stdout, stderr := SetupPipes(t, cmd)
	stdoutReader := bufio.NewReader(stdout)

 	if err := cmd.Start(); err != nil {
 		t.Fatal(err)
 	}
	
	eatInitPrompt(stdoutReader, t)

	return ShellTestContext{
		t,
		cmd,
		stdin,
		stdout,
		stderr,
		stdoutReader,
	}
}

func TestLocateExecutableFiles(t *testing.T) {
	
	context := InitTest(t)

	context.assertCmd("type echo\n", "echo is a shell builtin")
	context.assertCmd("type type\n", "type is a shell builtin")
	context.assertCmd("type exit\n", "exit is a shell builtin")

	context.assertCmd("type grep\n", "grep is /usr/bin/grep")
	context.assertCmd("type invalid_command\n", "invalid_command: not found")

	context.tearDown()
}

func TestPwd(t *testing.T) { 
	context := InitTest(t)
	
	context.assertCmd("pwd\n", "/Users/kostyamalinovskiy/study/go/codecrafters-shell-go/app")

	context.tearDown()
}
