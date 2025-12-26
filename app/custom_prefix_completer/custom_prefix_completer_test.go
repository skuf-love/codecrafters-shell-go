package custom_prefix_completer

import(
	"testing"
	"github.com/chzyer/readline"
	"fmt"
	"reflect"
)

func TestIncrementalComplete(t *testing.T){ 
	cmdNames :=[]string{"xyz ", "xyz_foo ", "xyz_foo_bar ", "xyz_foo_bar_baz "}
	completers := make([]readline.PrefixCompleterInterface, len(cmdNames))
	
	for i, cmdName := range cmdNames {
		completers[i] = readline.PcItem(cmdName)
	}
	baseCompleter := readline.NewPrefixCompleter(completers...)
	completer := New(baseCompleter, "$ ")

	cmdNameRunes := make([][]rune, len(cmdNames))
	
	for i, cmdName := range cmdNames {
		cmdNameRunes[i] = []rune(cmdName)
	}

	candidates, _ := completer.TryIncrementalComplete(cmdNameRunes) 

	expectedCandidate := []rune{'x', 'y', 'z'}
	expectedCandidates := [][]rune{expectedCandidate}
	if !reflect.DeepEqual(candidates, expectedCandidates) {
		t.Fatal(fmt.Sprintf("expected: %q, received: %q", expectedCandidates, candidates))
	}
}
