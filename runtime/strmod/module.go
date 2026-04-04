package strmod

import "moonbasic/runtime"

type Module struct{}

func NewModule() *Module {
	return &Module{}
}

func (m *Module) Register(_ runtime.Registrar) {
	// String builtins (LSET$, FORMAT$, etc.) live in runtime core; avoid duplicate registry keys.
}

func (m *Module) Shutdown() {}
