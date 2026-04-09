//go:build fullruntime

// Mario64 entities example — embedded MOON bytecode (see game.mbc).
// Rebuild bytecode from repo root:
//
//	go run . examples/mario64/main_entities.mb
//	copy /Y examples\mario64\main_entities.mbc cmd\mario64_entities\game.mbc
//
// Then: go build -tags fullruntime -o path\Mario64Entities.exe ./cmd/mario64_entities
package main

import (
	_ "embed"
	"flag"
	"fmt"
	"os"

	"moonbasic/compiler/pipeline"
)

//go:embed game.mbc
var gameMBC []byte

const version = "1.0.0"

func main() {
	var (
		debug   = flag.Bool("info", false, "print bytecode disassembly")
		trace   = flag.Bool("trace", false, "VM state trace")
		showVer = flag.Bool("version", false, "print version")
	)
	flag.Parse()

	if *showVer {
		fmt.Printf("Mario64 entities hop %s (moonBASIC embedded)\n", version)
		return
	}

	prog, err := pipeline.DecodeMOON(gameMBC)
	if err != nil {
		fmt.Fprintf(os.Stderr, "decode embedded game: %v\n", err)
		os.Exit(2)
	}

	opts := pipeline.Options{
		Debug: *debug,
		Trace: *trace,
		Out:   os.Stderr,
	}
	if err := pipeline.RunProgram(prog, opts); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(3)
	}
}
