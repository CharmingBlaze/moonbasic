// Package pipeline orchestrates the moonBASIC compiler and VM execution stages.
package pipeline

import (
	"fmt"
	"os"

	"moonbasic/compiler/arena"
	"moonbasic/compiler/codegen"
	"moonbasic/compiler/include"
	"moonbasic/compiler/parser"
	"moonbasic/compiler/semantic"
	"moonbasic/vm/moon"
	"moonbasic/vm/opcode"
	"moonbasic/vm"
	"io"
)

// Options carries configuration for the VM execution.
type Options struct {
	Debug bool      // If true, print disassembly before execution
	Trace bool      // If true, print VM state after each opcode
	Out   io.Writer // Output stream for trace and errors (default os.Stderr)

	// ProfileRecorder when non-nil accumulates per-source-line instruction counts during Execute.
	ProfileRecorder *vm.ProfileRecorder

	// HostArgs is argv for ARGC / COMMAND$; nil leaves Registry.HostArgs nil so those builtins use os.Args.
	HostArgs []string
}

// CompileSource parses, analyzes, and generates code from a string.
func CompileSource(name, src string) (*opcode.Program, error) {
	SyncPackageIncludeRoots()
	lines := parser.SplitLines(src)
	ar := arena.NewArena()
	defer ar.Reset()

	// 1. Parsing
	prog, err := parser.ParseSourceWithArena(name, src, ar)
	if err != nil {
		return nil, err
	}
	prog, err = include.ExpandWithArena(name, prog, ar)
	if err != nil {
		return nil, err
	}

	// 2. Semantic Analysis
	an := semantic.DefaultAnalyzer(name, lines)
	if err := an.Run(prog); err != nil {
		return nil, err
	}

	// 3. Code Generation
	g := codegen.New(name, lines)
	bc, err := g.Compile(prog)
	if err != nil {
		return nil, fmt.Errorf("[moonBASIC] CodeGen Error: %v", err)
	}

	return bc, nil
}

// CompileFile reads a file from disk and compiles it.
func CompileFile(path string) (*opcode.Program, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return CompileSource(path, string(data))
}

// CheckFile reads a file from disk and performs only semantic analysis.
func CheckFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return CheckSource(path, string(data))
}

// CheckSource performs parsing and semantic analysis only.
func CheckSource(name, src string) error {
	SyncPackageIncludeRoots()
	ar := arena.NewArena()
	defer ar.Reset()
	prog, err := parser.ParseSourceWithArena(name, src, ar)
	if err != nil {
		return err
	}
	prog, err = include.ExpandWithArena(name, prog, ar)
	if err != nil {
		return err
	}
	an := semantic.DefaultAnalyzer(name, parser.SplitLines(src))
	return an.Run(prog)
}

// EncodeMOON serializes a compiled program to MOON container bytes (.mbc).
func EncodeMOON(prog *opcode.Program) ([]byte, error) {
	return moon.Encode(prog)
}

// DecodeMOON loads a program from MOON bytes after validating magic and version.
func DecodeMOON(data []byte) (*opcode.Program, error) {
	return moon.Decode(data)
}

// DecodeMOONFromFile reads a file from disk and decodes it as MOON bytecode.
func DecodeMOONFromFile(path string) (*opcode.Program, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return DecodeMOON(data)
}
