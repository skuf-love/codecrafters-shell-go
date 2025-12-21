package main

import(
	"os"
	"os/exec"
	"strings"
	"testing"
	"fmt"
	"io"
	"bufio"
	"time"
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

func (c ShellTestContext) sendInput(input string) {
	_, err := c.stdin.Write([]byte(input))
	if err != nil {
		c.t.Fatal(err)
	}
}

func eatInitPrompt(reader *bufio.Reader, t *testing.T) {
	_, err := reader.ReadString('$')
	if err != nil {
 		t.Fatal(err)
	}
}
func (c ShellTestContext) readUntilPrompt(t *testing.T) (string, error) {
	var result strings.Builder

	reader := c.stdoutReader
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
			c.log(fmt.Sprintf("Inside goroutine reading byte: %q(%v)", b, b))
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
		case <- time.After(200 * time.Millisecond):
			c.log(fmt.Sprintf("Timeout goroutine, received result: %v", result.String()))
			c.stdin.Write([]byte("echo %\n")) // use % symbol to signal goroutine to stop reading and finis
			<- done
			c.log(fmt.Sprintf("Done received: %v", result.String()))
			return result.String(), nil
		}
	}
}

func (c ShellTestContext) assertCmd(input string, expectedOutput string) {
	c.sendInput(input)
	c.log("  AssertCmd: " + input)

	output, err := c.readUntilPrompt(c.t)

	if err != nil {
	    c.t.Fatalf("Failed to read initial prompt: %v", err)
	}

	trimmed := strings.Trim(output, "\n")

 	if trimmed != expectedOutput {
 		c.t.Errorf("----")
 		c.t.Errorf("Command: %s", input)
 		c.t.Errorf("expected output: %s", expectedOutput)
 		c.t.Errorf("Received output: %s", trimmed)
 		c.t.Errorf("----")
 	}
}

func (c ShellTestContext) tearDown() {
	c.sendInput("exit\n")

	c.logFile.WriteString("< < <\n")
	c.stdin.Close()
	c.logFile.Close()

	stderrBytes, err := io.ReadAll(c.stderr)
	stdoutBytes, err := io.ReadAll(c.stdout)


	if err := c.cmd.Wait(); err != nil {
		c.t.Fatalf("command failed: %s", err)
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
	logFile *os.File
}
func InitTest(t *testing.T) ShellTestContext {
	cmd := exec.Command("./testapp")
	stdin, stdout, stderr := SetupPipes(t, cmd)
	stdoutReader := bufio.NewReader(stdout)
	logFile, _ := os.OpenFile("./test.log", os.O_APPEND|os.O_WRONLY, 0644)
	logFile.WriteString("> > > Example Log Started\n")
//	if err != nil {
//		t.Fatal(err)
//	}

	if err := cmd.Start(); err != nil {
 		t.Fatal(err)
 	}
	
	logFile.WriteString("Command Started\n")
//	eatInitPrompt(stdoutReader, t)

	return ShellTestContext{
		t,
		cmd,
		stdin,
		stdout,
		stderr,
		stdoutReader,
		logFile,
	}
}

func (c ShellTestContext) log(entry string){
	c.logFile.WriteString(entry + "\n")
}

func TestLocateExecutableFiles(t *testing.T) {
	InitTest(t)
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
	
	home_path := os.Getenv("HOME")
	context.assertCmd("pwd\n", home_path + "/study/go/codecrafters-shell-go/app")

	context.tearDown()
}

func TestCd(t *testing.T) { 
	c := InitTest(t)
	home_path := os.Getenv("HOME")

	c.sendInput("cd ~\n")
	c.readUntilPrompt(c.t)
	c.assertCmd("pwd\n", home_path)

	c.sendInput("cd .\n")
	c.readUntilPrompt(c.t)
	c.assertCmd("pwd\n", home_path)

	c.sendInput("cd ..\n")
	c.readUntilPrompt(c.t)
	c.assertCmd("pwd\n", "/Users")

	c.sendInput("cd " + home_path + "/study\n")
	c.readUntilPrompt(c.t)
	c.assertCmd("pwd\n", home_path + "/study")

	c.sendInput("cd go\n")
	c.readUntilPrompt(c.t)
	c.assertCmd("pwd\n", home_path + "/study/go")

	c.assertCmd("cd nope\n", "cd: nope: No such file or directory")

	c.tearDown()
}

func TestEcho(t *testing.T) { 
	context := InitTest(t)

	context.assertCmd("echo 'hello   world'\n", "hello   world")

	context.tearDown()
}

func TestStdout(t *testing.T) {
	c := InitTest(t)

	c.sendInput("echo hello 1> file1.txt\n")
	c.assertCmd("cat file1.txt\n", "hello")
	os.Remove("file1.txt")

	c.sendInput("echo hello > file2.txt\n")
	c.assertCmd("cat file2.txt\n", "hello")
	os.Remove("file2.txt")



	c.sendInput("ls -1 pig > cow/dog.md\n")
	c.assertCmd("cat cow/dog.md\n", "grape\norange\npear")
	os.Remove("cow/dog.md")

	c.assertCmd("cat pig/grape nonexistent 1> cow/fox.md\n", "cat: nonexistent: No such file or directory")

	c.tearDown()
}

func TestStderr(t *testing.T) {
	context := InitTest(t)

	context.assertCmd("cat pig/grape nonexistent 2> cow/fox.md\n", "grape")

	context.assertCmd("cat cow/fox.md\n", "cat: nonexistent: No such file or directory")
	os.Remove("cow/fox.md")

	context.tearDown()
}

func TestStdoutErrRedirectAppend(t *testing.T) {
	c := InitTest(t)

	c.sendInput("cat pig/grape >> stdappend.md\n")
	//readUntilPrompt(c.t)
	c.sendInput("echo ololo 1>> stdappend.md\n")
	//readUntilPrompt(c.t)
	c.assertCmd("cat stdappend.md\n", "grape\nololo")

	os.Remove("stdappend.md")

	c.sendInput("cat pig/grape nonexistent 2>> errappend.md\n")
	c.readUntilPrompt(c.t)
	c.sendInput("cat pig/grape nonexistent 2>> errappend.md\n")
	c.readUntilPrompt(c.t)
	c.assertCmd("cat errappend.md\n", "cat: nonexistent: No such file or directory\ncat: nonexistent: No such file or directory")

	os.Remove("errappend.md")

	c.tearDown()
}
