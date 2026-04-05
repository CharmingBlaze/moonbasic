package mbgame

import "moonbasic/runtime"

func (m *Module) registerPure(r runtime.Registrar) {
	m.registerCollisionBuiltins(r)
	m.registerMathBuiltins(r)
	m.registerColorFormatBuiltins(r)
	m.registerEaseNoiseRandBuiltins(r)
}
