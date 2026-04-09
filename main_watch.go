//go:build fullruntime

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"

	"moonbasic/compiler/pipeline"
)

func runWatch(path string, opts pipeline.Options) int {
	if strings.EqualFold(filepath.Ext(path), ".mbc") {
		fmt.Fprintln(os.Stderr, "error: --watch expects a source .mb file")
		return 2
	}
	run := func() int {
		prog, err := pipeline.CompileFile(path)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return 2
		}
		if err := pipeline.RunProgram(prog, opts); err != nil {
			fmt.Fprintln(os.Stderr, err)
			return 3
		}
		return 0
	}
	if code := run(); code != 0 {
		// Keep watching even if first run fails (syntax errors while editing).
		fmt.Fprintln(os.Stderr, "(watch continues; fix errors and save)")
	}

	w, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Fprintf(os.Stderr, "fsnotify: %v\n", err)
		return 1
	}
	defer w.Close()

	dir := filepath.Dir(path)
	if err := w.Add(dir); err != nil {
		fmt.Fprintf(os.Stderr, "watch %s: %v\n", dir, err)
		return 1
	}
	absPath, err := filepath.Abs(path)
	if err != nil {
		absPath = path
	}

	var debounce *time.Timer
	resetDebounce := func() {
		if debounce != nil {
			debounce.Stop()
		}
		debounce = time.AfterFunc(200*time.Millisecond, func() {
			fmt.Fprintf(os.Stderr, "\n--- moonbasic watch: recompiling %s ---\n", filepath.Base(absPath))
			run()
		})
	}

	for {
		select {
		case ev, ok := <-w.Events:
			if !ok {
				return 0
			}
			evPath, err := filepath.Abs(ev.Name)
			if err != nil {
				evPath = ev.Name
			}
			if !strings.EqualFold(evPath, absPath) {
				continue
			}
			if ev.Op&(fsnotify.Write|fsnotify.Create|fsnotify.Rename) != 0 {
				resetDebounce()
			}
		case err, ok := <-w.Errors:
			if !ok {
				return 0
			}
			fmt.Fprintf(os.Stderr, "watch: %v\n", err)
		}
	}
}
