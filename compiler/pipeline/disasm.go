package pipeline

import (
	"fmt"
	"io"
	"sort"

	"moonbasic/vm/opcode"
)

// PrintProgramDisassembly writes a human-readable listing of every chunk in prog.
// If sourceLines is non-nil and indexed by 1-based line, a source snippet is appended.
func PrintProgramDisassembly(prog *opcode.Program, w io.Writer, sourceLines []string) {
	if prog == nil || prog.Main == nil {
		fmt.Fprintln(w, "(nil program)")
		return
	}
	printChunkDisasm(w, prog.Main, sourceLines)
	names := sortedChunkNames(prog.Functions)
	for _, name := range names {
		printChunkDisasm(w, prog.Functions[name], sourceLines)
	}
}

func sortedChunkNames(m map[string]*opcode.Chunk) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

func printChunkDisasm(w io.Writer, ch *opcode.Chunk, sourceLines []string) {
	if ch == nil {
		return
	}
	fmt.Fprintf(w, "=== %s ===\n", ch.Name)
	for i, instr := range ch.Instructions {
		line := int32(-1)
		if i < len(ch.SourceLines) {
			line = ch.SourceLines[i]
		}
		fmt.Fprintf(w, "%04d  L%-4d  %s", i, line, instr.String())
		if sourceLines != nil && line >= 1 && int(line) <= len(sourceLines) {
			fmt.Fprintf(w, "  | %s", sourceLines[line-1])
		}
		fmt.Fprintln(w)
	}
	fmt.Fprintln(w)
}
