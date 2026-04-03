// Command moondoc generates static HTML reference from commands.json.
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"moonbasic/internal/docgen"
)

func main() {
	out := flag.String("out", "site", "output directory")
	jsonPath := flag.String("commands", "", "path to commands.json (default: compiler/builtinmanifest/commands.json relative to cwd)")
	flag.Parse()
	jp := *jsonPath
	if jp == "" {
		jp = docgen.DefaultCommandsPath()
	}
	if !filepath.IsAbs(jp) {
		if abs, err := filepath.Abs(jp); err == nil {
			jp = abs
		}
	}
	if err := docgen.Generate(jp, *out); err != nil {
		fmt.Fprintf(os.Stderr, "moondoc: %v\n", err)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stderr, "wrote site to %s\n", *out)
}
