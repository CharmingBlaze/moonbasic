package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"moonbasic/compiler/parser"
	"moonbasic/compiler/pipeline"
)

func disasmMBC(path string) int {
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return 1
	}
	prog, err := pipeline.DecodeMOON(data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "load bytecode: %v\n", err)
		return 2
	}
	var lines []string
	base := strings.TrimSuffix(path, filepath.Ext(path))
	if b, err := os.ReadFile(base + ".mb"); err == nil {
		lines = parser.SplitLines(string(b))
	}
	fmt.Fprintf(os.Stdout, "// disassembly: %s\n", path)
	pipeline.PrintProgramDisassembly(prog, os.Stdout, lines)
	return 0
}
