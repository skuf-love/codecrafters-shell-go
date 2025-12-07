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

func readUntil (reader *bufio.Reader, expected string, timeout time.Duration) (string, error) {
	//return "", nil
	result := ""
	deadline := time.Now().Add(timeout)
	
	for time.Now().Before(deadline) {
		// Try to read with a short timeout
		line, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			return result, err
		}
		fmt.Printf("received string: %s", line)
		
		result = line
		
		if strings.Contains(line, expected) {
			return result, nil
		}

		if strings.Contains(line, "$") {
			continue
		}
		
		// Small delay to avoid busy waiting
		time.Sleep(10 * time.Millisecond)
	}
	
	return result, fmt.Errorf("timeout waiting for: %q, got: %q", expected, result)
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
	//var result strings.Builder
	
	line, err := reader.ReadString('$')
	if err == io.EOF {
		//return result.String(), err
		return "", err
	}

	line, _ = strings.CutSuffix(line, "\n$")
	return line, nil
	
	//	result.WriteString(line)
	//	return result.String(), nil
	//for {
	//	line, err := reader.ReadString('\n')
	//	if err == io.EOF {
	//		return result.String(), err
	//	}

	//	trimmed := strings.TrimSpace(line)
	//	if trimmed == "$" {
	//		return result.String(), nil
	//	}
	//	result.WriteString(line)
	//}
}

func assertCmd(input string, expectedOutput string, stdin io.WriteCloser, stdoutReader *bufio.Reader, t *testing.T) {
	sendInput(input, t, stdin)

	output, err := readUntilPrompt(stdoutReader, t)

	if err != nil {
	    t.Fatalf("Failed to read initial prompt: %v", err)
	}

 	if !strings.Contains(output, expectedOutput) {
 		t.Errorf("----")
 		t.Errorf("Command: %s", input)
 		t.Errorf("expected output: %s", expectedOutput)
 		t.Errorf("Received output: %s", output)
 		t.Errorf("----")
 	}
}

func TestLocateExecutableFiles(t *testing.T) {
	cmd := exec.Command("./testapp")
	stdin, stdout, stderr := SetupPipes(t, cmd)
	 
	stdoutReader := bufio.NewReader(stdout)
	// stderrReader := bufio.NewReader(stderr)

 	if err := cmd.Start(); err != nil {
 		t.Fatal(err)
 	}
//	// Wait for initial prompt
	eatInitPrompt(stdoutReader, t)


	//output, err := readUntilPrompt(stdoutReader, t)
	//if err != nil {
	//    t.Fatalf("Failed to read initial prompt: %v", err)
	//}
	//t.Logf("Initial prompt: %s", output)
//
// 	input := "type echo\n"
//	
//
//	output, err = readUntil(stdoutReader, "echo", 20*time.Millisecond)
//        if err != nil {
//		t.Fatalf("Command failed: %v", err)
//	}
//
//	if !strings.Contains(output, "echo is a shell builtin") {
//		t.Errorf("Expected output to contain 'echo is a shell builtin', got: %s", output)
//	}
//	t.Logf("Received: %s", strings.TrimSpace(output))
	
	//stdin.Close()

	//
 	//if !strings.Contains(output, "echo is a shell builtin") {
 	//	t.Errorf("unexpected output: %s", output)
 	//}

	assertCmd("type echo\n", "echo is a shell builtin", stdin, stdoutReader, t)
	assertCmd("type type\n", "type is a shell builtin", stdin, stdoutReader, t)
	assertCmd("type exit\n", "exit is a shell builtin", stdin, stdoutReader, t)

	sendInput("exit\n", t, stdin)

	stdin.Close()

	stderrBytes, err := io.ReadAll(stderr)
	stdoutBytes, err := io.ReadAll(stdout)

	if err := cmd.Wait(); err != nil {
		t.Fatalf("command failed: %s", err)
	}

	if err != nil {
		t.Fatal(err)
	}


	if err != nil {
		t.Fatal(err)
	}
	
	fmt.Printf("%v",string(stdoutBytes))


	if len(stderrBytes) > 0 {
		t.Logf("Stderr: %s", string(stderrBytes))
	}


}

