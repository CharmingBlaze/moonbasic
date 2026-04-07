//go:build !compiler_only && !nonet

package pipeline

import (
	mbnet "moonbasic/runtime/net"
	"moonbasic/runtime"
	"moonbasic/vm"
)

func registerNetModules(reg *runtime.Registry) {
	reg.RegisterModule(mbnet.NewModule())
}

func wireNetCallbacks(reg *runtime.Registry, machine *vm.VM) {
	var netMod *mbnet.Module
	for _, m := range reg.Modules {
		if mod, ok := m.(*mbnet.Module); ok {
			netMod = mod
			break
		}
	}
	if netMod != nil {
		netMod.SetUserInvoker(machine.CallUserFunction)
	}
	mbnet.SeedMultiplayerGlobals(machine.Globals)
}
