// Package mblight implements LIGHT.* handles (configuration for future lighting hooks).
package mblight

import (
	"fmt"
	"math"
	"strings"
	"sync"

	"moonbasic/runtime"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
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
	r.Register("LIGHT.SETDIR", "light", m.lightSetDir)
	r.Register("LIGHT.SETSHADOW", "light", m.lightSetShadow)
}

// Shutdown implements runtime.Module.
func (m *Module) Shutdown() {}

type lightObj struct {
	kind      string
	r, g, b   float32
	intensity float32
	// Normalized direction the light travels (e.g. toward the scene); used by 3D shadow + PBR.
	dirX, dirY, dirZ float32
	shadow           bool
}

func (o *lightObj) TypeName() string { return "Light" }

func (o *lightObj) TypeTag() uint16 { return heap.TagLight }

func (o *lightObj) Free() {}

func normalizeDir3(x, y, z float32) (float32, float32, float32) {
	l := float32(math.Sqrt(float64(x*x + y*y + z*z)))
	if l < 1e-8 {
		return 0, -1, 0
	}
	return x / l, y / l, z / l
}

// ShadowCasterHandle returns the heap handle of the light marked for shadow maps, or 0.
func ShadowCasterHandle() heap.Handle {
	shadowMu.Lock()
	defer shadowMu.Unlock()
	return shadowCasterHandle
}

// LightDirection returns the normalized light travel direction for a light handle.
func LightDirection(hs *heap.Store, h heap.Handle) (x, y, z float32, ok bool) {
	if hs == nil || h == 0 {
		return 0, -1, 0, false
	}
	o, err := heap.Cast[*lightObj](hs, h)
	if err != nil {
		return 0, -1, 0, false
	}
	return o.dirX, o.dirY, o.dirZ, true
}

// LightDiffuse scales RGB by intensity for shading (PBR sun color).
func LightDiffuse(hs *heap.Store, h heap.Handle) (r, g, b float32) {
	if hs == nil || h == 0 {
		return 1, 1, 1
	}
	o, err := heap.Cast[*lightObj](hs, h)
	if err != nil {
		return 1, 1, 1
	}
	return o.r * o.intensity, o.g * o.intensity, o.b * o.intensity
}

func (m *Module) lightSetDir(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("LIGHT.SETDIR: heap not bound")
	}
	if len(args) != 4 {
		return value.Nil, fmt.Errorf("LIGHT.SETDIR expects (light, x, y, z)")
	}
	if args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("LIGHT.SETDIR: light must be a handle")
	}
	o, err := heap.Cast[*lightObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	x, ok1 := argF32(args[1])
	y, ok2 := argF32(args[2])
	z, ok3 := argF32(args[3])
	if !ok1 || !ok2 || !ok3 {
		return value.Nil, fmt.Errorf("LIGHT.SETDIR: x, y, z must be numeric")
	}
	o.dirX, o.dirY, o.dirZ = normalizeDir3(x, y, z)
	return value.Nil, nil
}

func argF32(v value.Value) (float32, bool) {
	if f, ok := v.ToFloat(); ok {
		return float32(f), true
	}
	if i, ok := v.ToInt(); ok {
		return float32(i), true
	}
	return 0, false
}

func (m *Module) lightSetShadow(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("LIGHT.SETSHADOW: heap not bound")
	}
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("LIGHT.SETSHADOW expects (light, enabled?)")
	}
	if args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("LIGHT.SETSHADOW: light must be a handle")
	}
	h := heap.Handle(args[0].IVal)
	o, err := heap.Cast[*lightObj](m.h, h)
	if err != nil {
		return value.Nil, err
	}
	en, ok := argBoolLight(args[1])
	if !ok {
		return value.Nil, fmt.Errorf("LIGHT.SETSHADOW: enabled must be bool or numeric")
	}
	shadowMu.Lock()
	defer shadowMu.Unlock()
	if en {
		o.shadow = true
		shadowCasterHandle = h
	} else {
		o.shadow = false
		if shadowCasterHandle == h {
			shadowCasterHandle = 0
		}
	}
	return value.Nil, nil
}

func argBoolLight(v value.Value) (bool, bool) {
	if v.Kind == value.KindBool {
		return v.IVal != 0, true
	}
	if i, ok := v.ToInt(); ok {
		return i != 0, true
	}
	if f, ok := v.ToFloat(); ok {
		return f != 0, true
	}
	return false, false
}

func (m *Module) lightMake(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("LIGHT.MAKE: heap not bound")
	}
	kind := "directional"
	if len(args) > 1 {
		return value.Nil, fmt.Errorf("LIGHT.MAKE expects 0 arguments (default directional, white, intensity 1) or 1 kind string")
	}
	if len(args) == 1 {
		if args[0].Kind != value.KindString {
			return value.Nil, fmt.Errorf("LIGHT.MAKE: optional argument must be a string kind (e.g. \"directional\")")
		}
		s, err := rt.ArgString(args, 0)
		if err != nil {
			return value.Nil, err
		}
		kind = strings.ToLower(strings.TrimSpace(s))
	}
	o := &lightObj{
		kind:      kind,
		r:         1,
		g:         1,
		b:         1,
		intensity: 1,
		dirX:      -0.4082483,
		dirY:      -0.8164966,
		dirZ:      -0.4082483, // normalize(-1,-2,-1)
	}
	id, err := m.h.Alloc(o)
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}
