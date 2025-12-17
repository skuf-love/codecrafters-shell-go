package shell_args

import(
	"testing"
	"reflect"
)

func AssertParse(input string, expected []string, t *testing.T) {
	result := ParseInput(input)

	command := expected[0]
	arguments := make([]string, 0)
	if len(expected) > 1 {
		arguments = expected[1:]
	}
	if !reflect.DeepEqual(result.Arguments, arguments) || result.CommandName != command {
		t.Errorf("ParseInput(%q)", input)
		t.Errorf("Command = %v, expected %v", result.CommandName, command)
		t.Errorf("Arguments = %v, expected %v", result.Arguments, arguments)
		t.Errorf("---")
	}
}

func TestQuotesAgrgsParse(t *testing.T) {

	AssertParse("echo", []string{"echo"}, t)

	AssertParse("echo 'hello    world'", []string{"echo", "hello    world"}, t)

	AssertParse("echo hello    world", []string{"echo", "hello", "world"}, t)


	AssertParse("echo 'hello''world'", []string{"echo", "helloworld"}, t)
	AssertParse("echo 'hello''world", []string{"echo", "helloworld"}, t)

}

func TestDoubleQuotesAgrgsParse(t *testing.T) {

	AssertParse("echo \"hello    world\"", []string{"echo", "hello    world"}, t)

	AssertParse("echo \"hello\"\"world\"", []string{"echo", "helloworld"}, t)

	AssertParse("echo \"hello\" \"world\"", []string{"echo", "hello", "world"}, t)

	AssertParse("echo \"shell's test\"", []string{"echo", "shell's test"}, t)
}

func TestBackslashParse(t *testing.T)  {
	AssertParse("echo \\'\\\"hello world\\\"\\'", []string{"echo", "'\"hello", "world\"'"}, t)

	AssertParse("echo world\\ \\ \\ \\ \\ \\ script", []string{"echo", "world      script"}, t)
}

func TestDoubleQuoteBackslashParse(t *testing.T)  {
	// echo "hello'script'\\n'world"
	// echo hello'script'\n'world
	AssertParse("echo \"hello'script'\\n'world\"", []string{"echo", "hello'script'\\n'world"}, t)

	//echo "hello\"insidequotes"script\"
	//hello"insidequotesscript"
	AssertParse("echo \"hello\\\"insidequotes\"script\\\"", []string{"echo", "hello\"insidequotesscript\""}, t)
}

func TestStdout(t *testing.T) {

	result := ParseInput("echo hello 1> file.txt")



	if result.StdoutPath != "file.txt" {
		t.Errorf("Expected stdoutPath: %v; Result: %v", "file.txt", result.StdoutPath)
	}
	if  !result.isStdoutRedirected() {
		t.Errorf("Expected isStdoutRedirected to be true but got %v", result.isStdoutRedirected())
	}

	result = ParseInput("echo hello > file2.txt")

	if result.StdoutPath != "file2.txt" {
		t.Errorf("Expected stdoutPath: %v; Result: %v", "file2.txt", result.StdoutPath)
	}
	if  !result.isStdoutRedirected() {
		t.Errorf("Expected isStdoutRedirected to be true but got %v", result.isStdoutRedirected())
	}

}

