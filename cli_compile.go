package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"moonbasic/compiler/pipeline"
)

func compileToMBC(path string) error {
	if strings.EqualFold(filepath.Ext(path), ".mbc") {
		return fmt.Errorf("error: compiler expects a source file (.mb), not a .mbc file")
	}
	prog, err := pipeline.CompileFile(path)
	if err != nil {
		return err
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
