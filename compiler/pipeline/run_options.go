package pipeline

import (
	"io"

	"moonbasic/lineprof"
)

// Options carries configuration for the VM execution.
type Options struct {
	Debug bool      // If true, print disassembly before execution
	Trace bool      // If true, print VM state after each opcode
	Out   io.Writer // Output stream for trace and errors (default os.Stderr)

	// ProfileRecorder when non-nil accumulates per-source-line instruction counts during Execute.
	ProfileRecorder lineprof.LineProfiler

	// HostArgs is argv for ARGC / COMMAND$; nil leaves Registry.HostArgs nil so those builtins use os.Args.
	HostArgs []string
}
