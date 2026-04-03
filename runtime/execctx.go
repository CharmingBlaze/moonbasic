package runtime

import (
	"io"
	"os"

	"moonbasic/vm/value"
)

// execReg holds the registry active during a nested Call chain (single-threaded VM).
var execReg *Registry

// ActiveRegistry returns the registry currently executing a native builtin, or nil.
func ActiveRegistry() *Registry {
	return execReg
}

func enterCall(r *Registry) func() {
	prev := execReg
	execReg = r
	return func() { execReg = prev }
}

// DiagWriter returns the host diagnostic stream (pipeline Options.Out), or os.Stderr.
func DiagWriter() io.Writer {
	if execReg != nil && execReg.DiagOut != nil {
		return execReg.DiagOut
	}
	return os.Stderr
}

// ArgString resolves a string Value using the active program's string pool.
func ArgString(v value.Value) string {
	if execReg == nil || execReg.Prog == nil {
		return value.StringAt(v, nil)
	}
	return value.StringAt(v, execReg.Prog.StringTable)
}

// RetString interns s into the active program pool and returns a string Value.
func RetString(s string) value.Value {
	if execReg == nil || execReg.Prog == nil {
		panic("moonbasic/runtime: RetString called without active VM program")
	}
	return value.FromStringIndex(execReg.Prog.InternString(s))
}
