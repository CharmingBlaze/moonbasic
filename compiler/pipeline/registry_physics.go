//go:build fullruntime && !nophysics

package pipeline

import (
	mbcharcontroller "moonbasic/runtime/charcontroller"
	mbphysics2d "moonbasic/runtime/physics2d"
	mbphysics3d "moonbasic/runtime/physics3d"
	"moonbasic/runtime"
	"moonbasic/vm"
)

func registerPhysicsModules(reg *runtime.Registry) {
	// Physics / character: register before manifest so natives win; char before physics3d
	// so Shutdown frees CharacterVirtual instances before the Jolt world is torn down.
	reg.RegisterModule(mbcharcontroller.NewModule())
	reg.RegisterModule(mbphysics2d.NewModule())
	reg.RegisterModule(mbphysics3d.NewModule())
}

func wirePhysicsCallbacks(reg *runtime.Registry, machine *vm.VM) {
	var p3 *mbphysics3d.Module
	for _, m := range reg.Modules {
		if mod, ok := m.(*mbphysics3d.Module); ok {
			p3 = mod
			break
		}
	}
	if p3 != nil {
		p3.SetUserInvoker(machine.CallUserFunction)
		p3.SetVM(machine)
	}
}
