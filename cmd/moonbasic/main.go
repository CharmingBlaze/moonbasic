//go:build !fullruntime

// moonBASIC Compiler (CLI)
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"moonbasic/compiler/pipeline"
	"moonbasic/lsp"
)

const version = "1.2.10"

func main() {
	var (
		checkOnly   = flag.Bool("check", false, "parse and type-check only")
		showVer     = flag.Bool("version", false, "print version and exit")
		lspMode     = flag.Bool("lsp", false, "run Language Server Protocol (stdio) for editors")
		disasm      = flag.Bool("disasm", false, "print human-readable bytecode for a .mbc file")
		symbolsOut  = flag.String("symbols-out", "", "also write symbol table JSON (path, globals with Persistent, funcs, types) to this file when compiling")
	)

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "moonBASIC Compiler %s\n", version)
		fmt.Fprintf(os.Stderr, "Usage:\n")
		fmt.Fprintf(os.Stderr, "  moonbasic [flags] <source.mb>     compile to .mbc\n")
		fmt.Fprintf(os.Stderr, "  moonbasic --check <source.mb>     parse and type-check only\n")
		fmt.Fprintf(os.Stderr, "  moonbasic --lsp                   language server on stdio\n")
		fmt.Fprintf(os.Stderr, "  moonbasic --disasm <file.mbc>     disassemble bytecode\n")
		fmt.Fprintf(os.Stderr, "Flags:\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	if *showVer {
		fmt.Printf("moonBASIC Compiler %s\n", version)
		return
	}

	if *lspMode {
		if err := lsp.Serve(); err != nil {
			fmt.Fprintf(os.Stderr, "lsp: %v\n", err)
			os.Exit(1)
		}
		return
	}

	args := flag.Args()
	if len(args) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	path := args[0]

	if *disasm {
		if !strings.EqualFold(filepath.Ext(path), ".mbc") {
			fmt.Fprintln(os.Stderr, "error: --disasm requires a .mbc file")
			os.Exit(2)
		}
		// TODO: Implement or call pipeline disasm
		prog, err := pipeline.DecodeMOONFromFile(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "decode: %v\n", err)
			os.Exit(2)
		}
		pipeline.PrintProgramDisassembly(prog, os.Stdout, nil)
		return
	}

	if *checkOnly {
		if err := pipeline.CheckFile(path); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(2)
		}
		fmt.Println("Check: OK")
		return
	}

	// Default: Compile to MBC
	if err := compileToMBC(path, *symbolsOut); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
}

func compileToMBC(path, symbolsPath string) error {
	if strings.EqualFold(filepath.Ext(path), ".mbc") {
		return fmt.Errorf("error: compiler expects a source file (.mb), not a .mbc file")
	}
	srcBytes, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	src := string(srcBytes)
	prog, err := pipeline.CompileSource(path, src)
	if err != nil {
		return err
	}
	if symbolsPath != "" {
		raw, err := pipeline.ExportSymbolTableJSON(path, src)
		if err != nil {
			return fmt.Errorf("symbols: %w", err)
		}
		if err := os.WriteFile(symbolsPath, raw, 0644); err != nil {
			return fmt.Errorf("write symbols: %w", err)
		}
		fmt.Fprintf(os.Stderr, "wrote %s\n", symbolsPath)
	}
	out := mbcOutPath(path)
	data, err := pipeline.EncodeMOON(prog)
	if err != nil {
		return fmt.Errorf("encode: %v", err)
	}
	if err := os.WriteFile(out, data, 0644); err != nil {
		return fmt.Errorf("write %s: %v", out, err)
	}
	fmt.Fprintf(os.Stderr, "wrote %s\n", out)
	return nil
}

func mbcOutPath(src string) string {
	ext := filepath.Ext(src)
	base := strings.TrimSuffix(src, ext)
	return base + ".mbc"
}
