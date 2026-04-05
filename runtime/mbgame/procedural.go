package mbgame

import "moonbasic/runtime"

func (m *Module) registerProcedural(r runtime.Registrar) {
	_ = m
	_ = r
	// Extra procedural helpers (beyond PERLIN/SIMPLEX/FBM in register_ease_noise_rand) go here.
}
