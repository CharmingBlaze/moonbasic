// Package pipeline orchestrates the entire moonBASIC compilation process.
// It acts as the primary entry point for turning raw source code into executable VM bytecode.
//
// The compilation pipeline consists of the following automated stages:
//  1. Parsing: Source is broken into tokens (Lexer) and structured into an AST (Parser).
//  2. Expansion: Included files (INCLUDE "file.mb") are recursively expanded into the AST.
//  3. Symbol Building: Types and functions are harvested into a symbol table.
//  4. Semantic Analysis: The AST is validated for correct types, scoped variables, and logic.
//  5. Code Generation: The validated AST is emitted as executable *opcode.Program bytecode.
//
// Developers integrating moonBASIC as a scripting language into a Go host should use
// CompileSource or CompileFile from this package to generate executable bytecode.
package pipeline

import (
	"fmt"
	"os"

	"moonbasic/compiler/arena"
	"moonbasic/compiler/codegen"
	"moonbasic/compiler/include"
	"moonbasic/compiler/parser"
	"moonbasic/compiler/semantic"
	"moonbasic/compiler/symtable"
	"moonbasic/vm/moon"
	"moonbasic/vm/opcode"
)

// CompileOptions configures the compilation process.
type CompileOptions struct {
	// ImplicitDeclaration enables modern syntax without VAR.
	// First assignment declares the variable with inferred type.
	ImplicitDeclaration bool

	// TypeInference enables automatic type detection from expressions.
	// When disabled, variables default to INT unless suffix is present.
	TypeInference bool

	// Debug enables verbose output during compilation.
	Debug bool
}

// CompileSource parses, analyzes, and generates code from a string.
func CompileSource(name, src string) (*opcode.Program, error) {
	return CompileSourceWithOptions(name, src, CompileOptions{
		ImplicitDeclaration: true, // Enable modern syntax by default
		TypeInference:       true,
	})
}

// CompileSourceWithOptions compiles with explicit options for modern syntax support.
func CompileSourceWithOptions(name, src string, opts CompileOptions) (*opcode.Program, error) {
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

	// 2. Two-Pass Symbol Table Builder (for implicit declaration)
	var symbols *symtable.Table
	if opts.ImplicitDeclaration {
		builder := symtable.NewBuilder()
		symbols = builder.Build(prog)
		if opts.Debug {
			fmt.Fprintf(os.Stderr, "[moonBASIC] Implicit declaration: collected %d globals\n", len(symbols.Funcs()))
		}
	}

	// 3. Semantic Analysis
	an := semantic.DefaultAnalyzer(name, lines)
	if err := an.Run(prog); err != nil {
		return nil, err
	}

	// 4. Code Generation (passing symbol table if using implicit declaration)
	g := codegen.NewWithSymbols(name, lines, symbols)
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
