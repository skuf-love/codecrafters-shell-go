package shell_args

import(
	"testing"
	"reflect"
)

func AssertParse(input string, expected []string, t *testing.T) {
	result := ParseInput(input)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("ParseInput(%q) = %v, expected %v", input, result, expected)
	}
}

func TestQuotesAgrgsParse(t *testing.T) {

	input := "echo"
	expected := []string{"echo"}
	result := ParseInput(input)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("ParseInput(%q) = %v, expected %v", input, result, expected)
	}


	// echo 'hello    world'	hello    world
	input = "echo 'hello    world'"
	expected = []string{"echo", "hello    world"}
	result = ParseInput(input)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("ParseInput(%q) = %v, expected %v", input, result, expected)
	}


	// echo hello    world		hello world
	input = "echo hello    world"
	expected = []string{"echo", "hello", "world"}
	result = ParseInput(input)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("ParseInput(%q) = %v, expected %v", input, result, expected)
	}

	// echo 'hello''world'		helloworld

	input = "echo 'hello''world'"
	expected = []string{"echo", "helloworld"}
	result = ParseInput(input)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("ParseInput(%q) = %v, expected %v", input, result, expected)
	}
	// echo hello''world		helloworld
	input = "echo 'hello''world"
	expected = []string{"echo", "helloworld"}
	result = ParseInput(input)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("ParseInput(%q) = %v, expected %v", input, result, expected)
	}

}

func TestDoubleQuotesAgrgsParse(t *testing.T) {

	input := "echo"
	expected := []string{"echo"}
	result := ParseInput(input)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("ParseInput(%q) = %v, expected %v", input, result, expected)
	}


	// echo 'hello    world'	hello    world
	input = "echo \"hello    world\""
	expected = []string{"echo", "hello    world"}
	result = ParseInput(input)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("ParseInput(%q) = %v, expected %v", input, result, expected)
	}



	input = "echo \"hello\"\"world\""
	expected = []string{"echo", "helloworld"}
	result = ParseInput(input)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("ParseInput(%q) = %v, expected %v", input, result, expected)
	}

	input = "echo \"hello\" \"world\""
	expected = []string{"echo", "hello", "world"}
	result = ParseInput(input)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("ParseInput(%q) = %v, expected %v", input, result, expected)
	}

	input = "echo \"shell's test\""
	expected = []string{"echo", "shell's test"}
	result = ParseInput(input)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("ParseInput(%q) = %v, expected %v", input, result, expected)
	}

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

