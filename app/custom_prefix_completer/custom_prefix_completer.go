package custom_prefix_completer

import(
	"fmt"
	"strings"
	"bytes"
	"github.com/chzyer/readline"
)

type Completer struct{
	prefixCompleter *readline.PrefixCompleter
	tabCount int32
	prevCandidates [][]rune
	prevLength int
	prevLine []rune
	prompt string
}

func New(baseCompleter *readline.PrefixCompleter, prompt string) Completer {
	return Completer{baseCompleter, 0, make([][]rune,0), int(0), make([]rune, 0), prompt}
}

func (cpc *Completer) Do(line []rune, pos int) (newLine [][]rune, length int) {
	lineStr := string(line[:pos])
	if string(cpc.prevLine) != lineStr {
		cpc.tabCount = 0
		cpc.prevCandidates = make([][]rune, 0)
	}
	cpc.tabCount++
	cpc.prevLine = line
	candidates, aLen := cpc.prefixCompleter.Do(line, pos)
	if len(candidates) > 1 {

		if cpc.tabCount == 1 {
			fmt.Print("\x07")
			cpc.prevCandidates = candidates
			return make([][]rune, 0), 0
		}else{
			cpc.tabCount = 0
			stringCandidates := make([]string, 0)
			var expanded string
			for _, cand := range cpc.prevCandidates {
				expanded = lineStr + string(cand)
				stringCandidates = append(stringCandidates, expanded)
			}
			fmt.Printf("\n%v\n", strings.Join(stringCandidates, " "))
			fmt.Printf("%v%v", cpc.prompt, lineStr)
			return [][]rune{}, len(lineStr)
		}
	} else if len(candidates) == 0{
		fmt.Print("\x07")
		return make([][]rune, 0), 0
	}else{
		return candidates, aLen
	}
}

func (cpc *Completer) Print(prefix string, level int, buf *bytes.Buffer) {
	cpc.prefixCompleter.Print(prefix, level, buf)
}
func (cpc *Completer) GetName() []rune {
	return cpc.prefixCompleter.GetName()
}
func (cpc *Completer) GetChildren() []readline.PrefixCompleterInterface {
	return cpc.prefixCompleter.GetChildren()
}
func (cpc *Completer) SetChildren(children []readline.PrefixCompleterInterface) {
	cpc.prefixCompleter.SetChildren(children)
}
