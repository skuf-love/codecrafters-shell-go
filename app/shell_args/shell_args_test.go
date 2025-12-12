package shell_args

import(
	"testing"
	"reflect"
)

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
