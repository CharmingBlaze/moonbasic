//go:build fullruntime

// moonBASIC Engine (CLI)
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"moonbasic/compiler/pipeline"
	"moonbasic/internal/driver"
	mbphysics3d "moonbasic/runtime/physics3d"
)

const version = "1.2.8"

func init() {
	// Pin the main goroutine before any work. OpenGL/GLFW contexts (Raylib) must stay on the OS
	// thread that created them; integrated GPUs (e.g. Intel Iris Xe) are especially strict.
	// pipeline.RunProgram also calls LockOSThread before the VM runs; this locks earlier so any
	// future code in main() before RunProgram stays on the same thread.
	runtime.LockOSThread()
}

func main() {
	var (
		debug     = flag.Bool("info", false, "print runtime banner + bytecode disassembly")
		trace     = flag.Bool("trace", false, "VM state trace")
		showVer   = flag.Bool("version", false, "print version")
	)

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "moonBASIC Engine %s\n", version)
		fmt.Fprintf(os.Stderr, "Usage:\n")
		fmt.Fprintf(os.Stderr, "  moonrun [flags] <file.mbc>        run MOON bytecode\n")
		fmt.Fprintf(os.Stderr, "  moonrun [flags] <source.mb>       compile and run source\n")
		fmt.Fprintf(os.Stderr, "Flags:\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	if *showVer {
		fmt.Printf("moonBASIC Engine %s\n", version)
		fmt.Fprintln(os.Stdout, "Runtime: raylib 5.5 | Jolt 5.1 | Box2D 3.0 | ENet 1.3")
		return
	}

	args := flag.Args()
	if len(args) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	path := args[0]
	opts := pipeline.Options{
		Debug: *debug,
		Trace: *trace,
		Out:   os.Stderr,
	}

	if *debug {
		sel := driver.GetDefaultDriver()
		fmt.Fprintf(os.Stderr, "driver: %s\n", sel.String())
	}

	mbphysics3d.LogJoltPhysicsBackendHint()

	if strings.EqualFold(filepath.Ext(path), ".mbc") {
		prog, err := pipeline.DecodeMOONFromFile(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "load error: %v\n", err)
			os.Exit(2)
		}
		if err := pipeline.RunProgram(prog, opts); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(3)
		}
	} else {
		// Compile and run from source (convenience for development)
		prog, err := pipeline.CompileFile(path)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(2)
		}
		if err := pipeline.RunProgram(prog, opts); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(3)
		}
	}
}
