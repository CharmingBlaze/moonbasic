//go:build compiler_only || nonet

package pipeline

import (
	"moonbasic/runtime"
	"moonbasic/vm"
)

func registerNetModules(reg *runtime.Registry) {
	// NOP
}

func wireNetCallbacks(reg *runtime.Registry, machine *vm.VM) {
	// NOP
}
