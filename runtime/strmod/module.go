package strmod

import "moonbasic/runtime"

type Module struct{}

func NewModule() *Module {
	return &Module{}
}

func (m *Module) Register(r runtime.Registrar) {
	registerStrings(r)
}

func (m *Module) Shutdown() {}
