package custom_prefix_completer

import(
	"fmt"
	"strings"
	"bytes"
	"github.com/chzyer/readline"
	"sort"
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
	//fmt.Printf("\nLine: %q, len: %v\n", line, aLen)
	//fmt.Printf("\nCandidates: %q\n", candidates)
	candidates, aLen = cpc.TryIncrementalComplete(candidates)
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

func (cpc *Completer) TryIncrementalComplete(candidates [][]rune) ([][]rune, int){

	candidatesCount := len(candidates)
	//fmt.Printf("\n not processed stringified candidates: %v\n", candidates)
	if candidatesCount > 0 {
		if len(candidates[0]) == 1 && candidates[0][0] == ' ' {
			candidates = candidates[1:len(candidates)]
			candidatesCount = len(candidates)
		}
	}

	//fmt.Printf("\n not stringified candidates: %v\n", candidates)
	if candidatesCount < 2 {
		return candidates, candidatesCount
	}
	stringCandidates := make([]string, candidatesCount)

	for i, cand :=  range candidates {
		stringCandidates[i] = strings.Trim(string(cand), " ")
	}

	sort.Strings(stringCandidates)
	//fmt.Printf("\n stringified candidates: %v\n", stringCandidates)
	
	//fmt.Printf("%q\n", stringCandidates)
	for i := 1; i < candidatesCount; i++{
		if !strings.HasPrefix(stringCandidates[i], stringCandidates[i-1]){
			return candidates, candidatesCount
		}
	}
	//fmt.Printf("%q\n", candidates)
	// figure out if the pattern is present
	//sort candidates by length ASC
	//return array with with only first candidate

	return [][]rune{[]rune(stringCandidates[0])}, 1
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
