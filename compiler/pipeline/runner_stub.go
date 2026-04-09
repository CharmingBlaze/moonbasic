//go:build !fullruntime

package pipeline

import (
	"errors"
	"moonbasic/vm/opcode"
)

// RunProgram is a stub when the game runtime is not linked.
func RunProgram(prog *opcode.Program, opts Options) error {
	return errors.New("moonBASIC runtime engine is not included in this build (rebuild with -tags fullruntime)")
}

func ListBuiltins() []string {
	// For compiler-only builds, we should probably still return the manifest keys if we want ListBuiltins to work.
	// But it usually requires the registry. We can just return an error or a message.
	return []string{"[Runtime disabled: build with -tags fullruntime]"}
}
