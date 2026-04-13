//go:build fullruntime

package pipeline

import (
	"fmt"
	goruntime "runtime"
	"os"

	"moonbasic/runtime"
	"moonbasic/vm"
	"moonbasic/vm/heap"
	"moonbasic/vm/opcode"
)

// RunProgram initializes the runtime and executes the program in the VM.
func RunProgram(prog *opcode.Program, opts Options) error {
	if opts.Out == nil {
		opts.Out = os.Stderr
	}

	if opts.Debug {
		fmt.Fprintln(opts.Out, prog.Main.Disassemble())
	}

	goruntime.LockOSThread()

	// 1. Initialize Runtime with compile-time default driver
	h := heap.New()
	d := DefaultDriver()
	reg := runtime.NewRegistry(h, d)
	setupRegistry(reg, h, opts)

	// 2. Setup VM
	machine := vm.New(reg, h)
	// Wire up modules that need to call back into the VM (using the machine just created)
	wireRegistryCallbacks(reg, machine)

	machine.Trace = opts.Trace
	machine.TraceOut = opts.Out
	machine.Profiler = opts.ProfileRecorder

	defer reg.Shutdown() // Raylib + heap cleanup on success or VM error

	// 3. Execution
	if err := machine.Execute(prog); err != nil {
		return err
	}

	return nil
}
