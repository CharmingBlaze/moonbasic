package bitwise

import "moonbasic/runtime"

type Module struct{}

func NewModule() *Module {
	return &Module{}
}

func (m *Module) Register(r runtime.Registrar) {
	registerBitwise(r)
}

func (m *Module) Shutdown() {}
