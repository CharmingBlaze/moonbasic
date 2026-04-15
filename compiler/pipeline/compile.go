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

	// StrictDeprecated treats deprecated manifest aliases (e.g. *.MAKE, *.SETPOSITION) as compile errors.
	StrictDeprecated bool
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
	an.StrictDeprecated = opts.StrictDeprecated
	if err := an.Run(prog); err != nil {
		return nil, err
	}
	for _, w := range an.DeprecationWarnings() {
		fmt.Fprintf(os.Stderr, "[moonBASIC] Warning: %s\n", w)
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

// CheckOptions configures standalone semantic analysis (check / LSP).
type CheckOptions struct {
	StrictDeprecated bool
}

// CheckFile reads a file from disk and performs only semantic analysis.
func CheckFile(path string) error {
	return CheckFileWithOptions(path, CheckOptions{})
}

// CheckFileWithOptions is like CheckFile with extra semantic options.
func CheckFileWithOptions(path string, opts CheckOptions) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return CheckSourceWithOptions(path, string(data), opts)
}

// CheckSourceWithNotices runs the same analysis as CheckSource and returns deprecation notices.
// If semantic analysis fails, err is non-nil; notices may still contain entries from code analyzed before the failure.
func CheckSourceWithNotices(name, src string, opts CheckOptions) ([]semantic.DeprecationNotice, error) {
	SyncPackageIncludeRoots()
	ar := arena.NewArena()
	defer ar.Reset()
	prog, err := parser.ParseSourceWithArena(name, src, ar)
	if err != nil {
		return nil, err
	}
	prog, err = include.ExpandWithArena(name, prog, ar)
	if err != nil {
		return nil, err
	}
	an := semantic.DefaultAnalyzer(name, parser.SplitLines(src))
	an.StrictDeprecated = opts.StrictDeprecated
	if err := an.Run(prog); err != nil {
		return an.DeprecationNotices(), err
	}
	return an.DeprecationNotices(), nil
}

// CheckSource performs parsing and semantic analysis only.
func CheckSource(name, src string) error {
	return CheckSourceWithOptions(name, src, CheckOptions{})
}

// CheckSourceWithOptions is CheckSource with semantic options.
func CheckSourceWithOptions(name, src string, opts CheckOptions) error {
	notices, err := CheckSourceWithNotices(name, src, opts)
	if opts.StrictDeprecated {
		return err
	}
	for _, n := range notices {
		fmt.Fprintf(os.Stderr, "[moonBASIC] Warning: %s\n", n.String())
	}
	return err
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
