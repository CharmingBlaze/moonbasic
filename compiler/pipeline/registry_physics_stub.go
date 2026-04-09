//go:build fullruntime && nophysics

package pipeline

import (
	"moonbasic/runtime"
	"moonbasic/vm"
)

func registerPhysicsModules(reg *runtime.Registry) {
	// NOP
}

func wirePhysicsCallbacks(reg *runtime.Registry, machine *vm.VM) {
	// NOP
}
