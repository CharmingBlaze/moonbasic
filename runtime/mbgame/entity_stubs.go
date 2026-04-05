package mbgame

import "moonbasic/runtime"

func (m *Module) registerEntityStubs(r runtime.Registrar) {
	_ = m
	_ = r
	// Reserved for lightweight ENTITY.* helpers that compose heap types; avoid duplicating ECS modules.
}
