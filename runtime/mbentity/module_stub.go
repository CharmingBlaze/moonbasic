//go:build !cgo && !windows

package mbentity

import (
	"fmt"

	"moonbasic/runtime"
	"moonbasic/runtime/texture"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

type Module struct {
	h            *heap.Store
	reg          runtime.Registrar
	tex          *texture.Module
	autoBuoyancy bool
}

func NewModule() *Module { return &Module{} }

func (m *Module) BindTextureModule(t *texture.Module) { m.tex = t }
func (m *Module) BindCamera(c runtime.Module)         {}

func (m *Module) BindHeap(h *heap.Store) { m.h = h }

func registerWaterAutoPhysics(m *Module, r runtime.Registrar) {
	r.Register("WATER.AUTOPHYSICS", "entity", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 1 {
			return value.Nil, fmt.Errorf("WATER.AUTOPHYSICS expects (toggle)")
		}
		on, _ := rt.ArgBool(args, 0)
		m.autoBuoyancy = on
		return value.Nil, nil
	})
}
