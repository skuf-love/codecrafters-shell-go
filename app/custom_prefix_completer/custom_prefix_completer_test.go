package custom_prefix_completer

import(
	"testing"
	"github.com/chzyer/readline"
	"fmt"
)

func testSomething(t *testing.T){ 
	cmdNames :=[]string{"xyz", "xyz_foo", "xyz_foo_bar", "xyz_foo_bar_baz"}
	completers := make([]readline.PrefixCompleterInterface, len(cmdNames))
	
	for i, cmdName := range cmdNames {
		completers[i] = readline.PcItem(cmdName)
	}
	baseCompleter := readline.NewPrefixCompleter(completers...)
	completer := New(baseCompleter, "$ ")
	fmt.Printf("completer %v", completer) 
}
