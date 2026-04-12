package weathermod

import (
	"moonbasic/runtime"
	"moonbasic/vm/heap"
)

type Module struct {
	h *heap.Store

	FogOn       bool
	FogNear     float32
	FogFar      float32
	FogR        int
	FogG        int
	FogB        int
	FogA        int
	WindStr     float32
	WindDirX    float32
	WindDirZ    float32
}

func NewModule() *Module {
	// Defaults tuned for typical character-scale scenes; override with FOG.SETRANGE / SETFOG.
	return &Module{FogNear: 8, FogFar: 72}
}

func (m *Module) BindHeap(h *heap.Store) { m.h = h }
func (m *Module) Register(r runtime.Registrar) { registerWeather(m, r) }
func (m *Module) Shutdown() {}

func (m *Module) Reset() {}

