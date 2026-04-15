// Package mblight implements LIGHT.* handles for PBR + directional shadow mapping.
package mblight

import (
	"sync"

	"moonbasic/runtime"
	"moonbasic/vm/heap"
)

// Module registers LIGHT.* builtins.
type Module struct {
	h *heap.Store
}

var (
	shadowMu           sync.Mutex
	shadowCasterHandle heap.Handle // light handle with shadow enabled; 0 = none
)

// NewModule creates the light module.
func NewModule() *Module { return &Module{} }

// BindHeap implements runtime.HeapAware.
func (m *Module) BindHeap(h *heap.Store) { m.h = h }

// Register implements runtime.Module.
func (m *Module) Register(r runtime.Registrar) {
	r.Register("LIGHT.MAKE", "light", m.lightMake)
	r.Register("LIGHT.CREATE", "light", m.lightMake)
	r.Register("LIGHT.CREATEPOINT", "light", m.lightCreatePoint)
	r.Register("LIGHT.CREATEDIRECTIONAL", "light", m.lightCreateDirectional)
	r.Register("LIGHT.CREATESPOT", "light", m.lightCreateSpot)
	r.Register("LIGHT.FREE", "light", m.lightFree)
	r.Register("LIGHT.SETDIR", "light", m.lightSetDir)
	r.Register("LIGHT.SETSHADOW", "light", m.lightSetShadow)
	r.Register("LIGHT.SETCOLOR", "light", m.lightSetColor)
	r.Register("LIGHT.SETINTENSITY", "light", m.lightSetIntensity)
	r.Register("LIGHT.GETPOS", "light", m.lightGetPos)
	r.Register("LIGHT.GETDIR", "light", m.lightGetDir)
	r.Register("LIGHT.GETCOLOR", "light", m.lightGetColor)
	r.Register("LIGHT.SETPOSITION", "light", m.lightSetPosition)
	r.Register("LIGHT.SETPOS", "light", m.lightSetPosition)
	r.Register("LIGHT.SETTARGET", "light", m.lightSetTarget)
	r.Register("LIGHT.SETSHADOWBIAS", "light", m.lightSetShadowBias)
	r.Register("LIGHT.SETINNERCONE", "light", m.lightSetInnerCone)
	r.Register("LIGHT.SETOUTERCONE", "light", m.lightSetOuterCone)
	r.Register("LIGHT.SETRANGE", "light", m.lightSetRange)
	r.Register("LIGHT.ENABLE", "light", m.lightEnable)
	r.Register("LIGHT.SETSTATE", "light", m.lightEnable)
	r.Register("LIGHT.ISENABLED", "light", m.lightIsEnabled)

	r.Register("CreateLight", "light", m.blitzCreateLight)
	r.Register("LightRange", "light", m.blitzLightRange)
	r.Register("LightColor", "light", m.blitzLightColor)
	m.registerPointLightBlitz(r)
}

// Shutdown implements runtime.Module.
func (m *Module) Shutdown() {}

func (m *Module) Reset() {}
