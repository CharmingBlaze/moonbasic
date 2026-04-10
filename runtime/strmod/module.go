package strmod

import "moonbasic/runtime"

type Module struct{}

func NewModule() *Module {
	return &Module{}
}

func (m *Module) Register(r runtime.Registrar) {
	registerStringBuiltins(r)
}

func (m *Module) Shutdown() {}

func (m *Module) Reset() {}

