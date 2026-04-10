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
	return &Module{FogNear: 100, FogFar: 500}
}

func (m *Module) BindHeap(h *heap.Store) { m.h = h }
func (m *Module) Register(r runtime.Registrar) { registerWeather(m, r) }
func (m *Module) Shutdown() {}

func (m *Module) Reset() {}

