package bench

import (
	"fmt"
	"os"
	"regexp"
	goruntime "runtime"
	"strings"
	"time"

	"moonbasic/compiler/pipeline"
)

var moonBenchLine = regexp.MustCompile(`(?i)MOONBENCH\s+(.+)`)

// Run compiles and executes the script, captures PRINT output for MOONBENCH lines, and prints metrics to w.
func Run(scriptPath string, opts pipeline.Options, w *os.File) error {
	if w == nil {
		w = os.Stderr
	}
	var m0, m1 goruntime.MemStats
	goruntime.ReadMemStats(&m0)

	prog, err := pipeline.CompileFile(scriptPath)
	if err != nil {
		return err
	}

	start := time.Now()
	out, runErr := withStdoutCapture(func() error {
		return pipeline.RunProgram(prog, opts)
	})
	elapsed := time.Since(start)
	goruntime.ReadMemStats(&m1)

	alloc := int64(m1.TotalAlloc - m0.TotalAlloc)
	lines := strings.Split(string(out), "\n")
	var benchLine string
	for _, ln := range lines {
		if m := moonBenchLine.FindStringSubmatch(strings.TrimSpace(ln)); m != nil {
			benchLine = strings.TrimSpace(m[1])
			break
		}
	}

	fmt.Fprintf(w, "BENCH file=%s wall_ms=%.3f heap_alloc_delta=%d sys=%d num_gc=%d\n",
		scriptPath, elapsed.Seconds()*1000, alloc, int64(m1.Sys-m0.Sys), int(m1.NumGC-m0.NumGC))
	if benchLine != "" {
		fmt.Fprintf(w, "BENCH %s\n", benchLine)
	} else {
		fmt.Fprintf(w, "BENCH (no MOONBENCH line in program output; add PRINT with substring MOONBENCH)\n")
	}
	return runErr
}
