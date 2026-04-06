// Command moonbasic is the moonBASIC compiler and runtime driver.
package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"moonbasic/compiler/pipeline"
	"moonbasic/internal/bench"
	"moonbasic/lsp"
	"moonbasic/vm"
)

// Release tags should match this string (e.g. v1.2.1) for distributor scripts that parse --version.
const version = "1.2.1" // Milestone 5 stable

// printRuntimeBanner writes the same runtime line as --version (for --info runs).
func printRuntimeBanner(w io.Writer) {
	if w == nil {
		w = os.Stderr
	}
	fmt.Fprintf(w, "moonBASIC %s | Runtime: Go 1.22 | raylib 5.5 | Jolt 5.1 | Box2D 3.0 | ENet 1.3\n", version)
}

func main() {
	os.Exit(run())
}

func run() int {
	if len(os.Args) >= 2 {
		switch strings.ToLower(strings.TrimSpace(os.Args[1])) {
		case "install", "list", "publish":
			return runPackageCLI(os.Args[1:])
		}
	}

	var (
		compileOnly = flag.Bool("compile", false, "compile source to .mbc (MOON bytecode)")
		runFile     = flag.String("run", "", "run a precompiled .mbc file (path argument)")
		checkOnly   = flag.Bool("check", false, "parse and type-check only")
		debug       = flag.Bool("info", false, "print runtime banner + bytecode disassembly before execution")
		trace       = flag.Bool("trace", false, "dump VM state after every instruction (Golden Trace)")
		showVer     = flag.Bool("version", false, "print version and exit")
		lspMode     = flag.Bool("lsp", false, "run Language Server Protocol (stdio) for editors")
		disasm      = flag.Bool("disasm", false, "print human-readable bytecode for a .mbc file")
		profile     = flag.Bool("profile", false, "count VM instructions per source line; print top 10 (source .mb only)")
		watch       = flag.Bool("watch", false, "recompile and rerun when the source .mb file changes")
		benchmark   = flag.String("benchmark", "", "run .mb as a benchmark; capture PRINT lines containing MOONBENCH")
		listBuiltins = flag.Bool("list-builtins", false, "print all registered built-in command keys and exit")
	)

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage:\n")
		fmt.Fprintf(os.Stderr, "  moonbasic [flags] <source.mb>     compile and run from source\n")
		fmt.Fprintf(os.Stderr, "  moonbasic [flags] <file.mbc>      run MOON bytecode (same as --run)\n")
		fmt.Fprintf(os.Stderr, "  moonbasic --run <file.mbc>        run MOON bytecode only\n")
		fmt.Fprintf(os.Stderr, "  moonbasic --disasm <file.mbc>     disassemble bytecode (optional paired .mb for source text)\n")
		fmt.Fprintf(os.Stderr, "  moonbasic --profile <source.mb>  run with per-line instruction profiling\n")
		fmt.Fprintf(os.Stderr, "  moonbasic --watch <source.mb>     watch and rerun on save\n")
		fmt.Fprintf(os.Stderr, "  moonbasic --lsp                   language server on stdio\n")
		fmt.Fprintf(os.Stderr, "  moonbasic install|list|publish    package manager (see docs/PACKAGES.md)\n")
		fmt.Fprintf(os.Stderr, "  moonbasic --benchmark <file.mb>   run benchmark harness\n")
		fmt.Fprintf(os.Stderr, "Flags:\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	if *benchmark != "" {
		opts := pipeline.Options{Out: os.Stderr}
		if err := bench.Run(*benchmark, opts, os.Stderr); err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			return 3
		}
		return 0
	}

	if *lspMode {
		if err := lsp.Serve(); err != nil {
			fmt.Fprintf(os.Stderr, "lsp: %v\n", err)
			return 1
		}
		return 0
	}

	if *showVer {
		fmt.Printf("moonBASIC %s\n", version)
		fmt.Fprintln(os.Stdout, "Runtime: Go 1.22 | raylib 5.5 | Jolt 5.1 | Box2D 3.0 | ENet 1.3")
		return 0
	}

	if *listBuiltins {
		keys := pipeline.ListBuiltins()
		for _, k := range keys {
			fmt.Println(k)
		}
		return 0
	}

	args := flag.Args()
	opts := pipeline.Options{
		Debug: *debug,
		Trace: *trace,
		Out:   os.Stderr,
	}

	if *runFile != "" {
		return runMBC(*runFile, opts)
	}

	if len(args) == 0 {
		flag.Usage()
		return 1
	}
	path := args[0]

	if *disasm {
		if !strings.EqualFold(filepath.Ext(path), ".mbc") {
			fmt.Fprintln(os.Stderr, "error: --disasm requires a .mbc file")
			return 2
		}
		return disasmMBC(path)
	}

	if *watch {
		return runWatch(path, opts)
	}

	if *profile {
		if strings.EqualFold(filepath.Ext(path), ".mbc") {
			fmt.Fprintln(os.Stderr, "error: --profile requires a source .mb file")
			return 2
		}
		rec := vm.NewProfileRecorder()
		opts.ProfileRecorder = rec
		prog, err := pipeline.CompileFile(path)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return 2
		}
		if *debug {
			printRuntimeBanner(opts.Out)
		}
		if err := pipeline.RunProgram(prog, opts); err != nil {
			fmt.Fprintln(os.Stderr, err)
			return 3
		}
		top := vm.TopProfileLines(rec, 10)
		fmt.Fprintln(os.Stderr, "Instruction hotspots (top 10 source lines):")
		for _, e := range top {
			fmt.Fprintf(os.Stderr, "  line %5d  %10d\n", e.Line, e.Count)
		}
		return 0
	}

	// Positional .mbc: run bytecode (no --run flag needed).
	if !*checkOnly && !*compileOnly && strings.EqualFold(filepath.Ext(path), ".mbc") {
		return runMBC(path, opts)
	}

	if *checkOnly {
		if err := pipeline.CheckFile(path); err != nil {
			fmt.Fprintln(os.Stderr, err)
			return 2
		}
		fmt.Println("Check: OK")
		return 0
	}

	if *compileOnly {
		return compileToMBC(path)
	}

	prog, err := pipeline.CompileFile(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 2
	}

	if *debug {
		printRuntimeBanner(opts.Out)
	}
	if err := pipeline.RunProgram(prog, opts); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 3
	}

	return 0
}

func compileToMBC(path string) int {
	if strings.EqualFold(filepath.Ext(path), ".mbc") {
		fmt.Fprintln(os.Stderr, "error: --compile expects a source file, not a .mbc bytecode file")
		return 2
	}
	prog, err := pipeline.CompileFile(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 2
	}
	out := mbcOutPath(path)
	data, err := pipeline.EncodeMOON(prog)
	if err != nil {
		fmt.Fprintf(os.Stderr, "encode: %v\n", err)
		return 2
	}
	if err := os.WriteFile(out, data, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "write %s: %v\n", out, err)
		return 1
	}
	fmt.Fprintf(os.Stderr, "wrote %s\n", out)
	return 0
}

func mbcOutPath(src string) string {
	ext := filepath.Ext(src)
	base := strings.TrimSuffix(src, ext)
	return base + ".mbc"
}

func runMBC(path string, opts pipeline.Options) int {
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
	if opts.Debug {
		printRuntimeBanner(opts.Out)
	}
	if err := pipeline.RunProgram(prog, opts); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 3
	}
	return 0
}
