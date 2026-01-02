package main

import(
	"os"
	"os/exec"
	"testing"
	"fmt"
	"io"
	"bufio"
	"github.com/codecrafters-io/shell-starter-go/app/shell_test_context"
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


func eatInitPrompt(reader *bufio.Reader, t *testing.T) {
	_, err := reader.ReadString('$')
	if err != nil {
 		t.Fatal(err)
	}
}



func InitTest(t *testing.T) shell_test_context.ShellTestContext {
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

	return shell_test_context.ShellTestContext{
		t,
		cmd,
		stdin,
		stdout,
		stderr,
		stdoutReader,
		logFile,
	}
}

func TestLocateExecutableFiles(t *testing.T) {
	context := InitTest(t)

	context.AssertCmd("type echo\n", "echo is a shell builtin")
	context.AssertCmd("type type\n", "type is a shell builtin")
	context.AssertCmd("type exit\n", "exit is a shell builtin")

	context.AssertCmd("type grep\n", "grep is /usr/bin/grep")
	context.AssertCmd("type invalid_command\n", "invalid_command: not found")

	context.TearDown()
}

func TestPwd(t *testing.T) { 
	context := InitTest(t)
	
	home_path := os.Getenv("HOME")
	context.AssertCmd("pwd\n", home_path + "/study/go/codecrafters-shell-go/app")

	context.TearDown()
}

func TestCd(t *testing.T) { 
	c := InitTest(t)
	home_path := os.Getenv("HOME")

	c.SendInput("cd ~\n")
	c.ReadUntilPrompt()
	c.AssertCmd("pwd\n", home_path)

	c.SendInput("cd .\n")
	c.ReadUntilPrompt()
	c.AssertCmd("pwd\n", home_path)

	c.SendInput("cd ..\n")
	c.ReadUntilPrompt()
	c.AssertCmd("pwd\n", "/Users")

	c.SendInput("cd " + home_path + "/study\n")
	c.ReadUntilPrompt()
	c.AssertCmd("pwd\n", home_path + "/study")

	c.SendInput("cd go\n")
	c.ReadUntilPrompt()
	c.AssertCmd("pwd\n", home_path + "/study/go")

	c.AssertCmd("cd nope\n", "cd: nope: No such file or directory")

	c.TearDown()
}

func TestEcho(t *testing.T) { 
	context := InitTest(t)

	context.AssertCmd("echo 'hello   world'\n", "hello   world")

	context.TearDown()
}

func TestStdout(t *testing.T) {
	c := InitTest(t)

	c.SendInput("echo hello 1> file1.txt\n")
	c.AssertCmd("cat file1.txt\n", "hello")
	os.Remove("file1.txt")

	c.SendInput("echo hello > file2.txt\n")
	c.AssertCmd("cat file2.txt\n", "hello")
	os.Remove("file2.txt")



	c.SendInput("ls -1 pig > cow/dog.md\n")
	c.AssertCmd("cat cow/dog.md\n", "grape\norange\npear")
	os.Remove("cow/dog.md")

	c.AssertCmd("cat pig/grape nonexistent 1> cow/fox.md\n", "cat: nonexistent: No such file or directory")

	c.TearDown()
}

func TestStderr(t *testing.T) {
	context := InitTest(t)

	context.AssertCmd("cat pig/grape nonexistent 2> cow/fox.md\n", "grape")

	context.AssertCmd("cat cow/fox.md\n", "cat: nonexistent: No such file or directory")
	os.Remove("cow/fox.md")

	context.TearDown()
}

func TestStdoutErrRedirectAppend(t *testing.T) {
	c := InitTest(t)

	c.SendInput("cat pig/grape >> stdappend.md\n")
	c.SendInput("echo ololo 1>> stdappend.md\n")
	c.AssertCmd("cat stdappend.md\n", "grape\nololo")

	os.Remove("stdappend.md")

	c.SendInput("cat pig/grape nonexistent 2>> errappend.md\n")
	c.ReadUntilPrompt()
	c.SendInput("cat pig/grape nonexistent 2>> errappend.md\n")
	c.ReadUntilPrompt()
	c.AssertCmd("cat errappend.md\n", "cat: nonexistent: No such file or directory\ncat: nonexistent: No such file or directory")

	os.Remove("errappend.md")

	c.TearDown()
}
//
//func TestDualCommandPipeline(t *testing.T) {
//	c := InitTest(t)
//	
//	//c.AssertCmd("cat ~/tmp/file | wc", "5      10      77")
//
//	//c.AssertCmd("il -f ~/tmp/file-1 | head -n 3", "raspberry strawberry\npear mango\npineapple apple")
//
//	c.TearDown()
//}
