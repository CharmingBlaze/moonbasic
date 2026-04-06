//go:build cgo || (windows && !cgo)

package mbentity

import (
	"fmt"
	"strings"

	"moonbasic/runtime"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// entityRefFreeHookOwner tracks which Module installed heap.EntityFreeHook so Shutdown can clear it.
var entityRefFreeHookOwner *Module

// clearEntityRefFreeHookIfOwner clears the global hook when this module shuts down (avoids dangling purge callbacks).
func clearEntityRefFreeHookIfOwner(m *Module) {
	if entityRefFreeHookOwner == m {
		heap.EntityFreeHook = nil
		entityRefFreeHookOwner = nil
	}
}

// registerBlitzEntityHandles registers CUBE / SPHERE constructors returning ENTITYREF handles
// and wires heap.EntityFreeHook. See docs/reference/BLITZ3D.md ("Dot-syntax entities").
func registerBlitzEntityHandles(m *Module, r runtime.Registrar) {
	entityRefFreeHookOwner = m
	heap.EntityFreeHook = func(id int64) { m.purgeEntityByID(id) }

	r.Register("CUBE", "entity", runtime.AdaptLegacy(m.blitzCube))
	r.Register("SPHERE", "entity", runtime.AdaptLegacy(m.blitzSphere))
}

func (m *Module) wrapEntityRef(id int64) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("CUBE/SPHERE: heap not bound")
	}
	if id < 1 {
		return value.Nil, fmt.Errorf("invalid entity id")
	}
	h, err := m.h.Alloc(&heap.EntityRef{ID: id})
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(h), nil
}

// blitzCube: CUBE() default 1×1×1, or CUBE(w#, h#, d#).
func (m *Module) blitzCube(args []value.Value) (value.Value, error) {
	var w, h, d float32 = 1, 1, 1
	switch len(args) {
	case 0:
	case 3:
		var ok1, ok2, ok3 bool
		w, ok1 = argF32(args[0])
		h, ok2 = argF32(args[1])
		d, ok3 = argF32(args[2])
		if !ok1 || !ok2 || !ok3 {
			return value.Nil, fmt.Errorf("CUBE: dimensions must be numeric")
		}
	default:
		return value.Nil, fmt.Errorf("CUBE expects 0 or 3 arguments (w#, h#, d#)")
	}
	st := m.store()
	id := st.nextID
	st.nextID++
	e := newDefaultEnt(id)
	e.kind = entKindBox
	e.r, e.g, e.b = 180, 180, 200
	e.w, e.h, e.d = w, h, d
	e.static = true
	e.useSphere = false
	e.gravity = 0
	st.ents[id] = e
	return m.wrapEntityRef(id)
}

// blitzSphere: SPHERE(radius# [, segments]) — static sphere for drawing/collision; default 16 segments.
func (m *Module) blitzSphere(args []value.Value) (value.Value, error) {
	if len(args) < 1 || len(args) > 2 {
		return value.Nil, fmt.Errorf("SPHERE expects 1 or 2 arguments (radius# [, segments])")
	}
	rad, ok1 := argF32(args[0])
	if !ok1 || rad <= 0 {
		return value.Nil, fmt.Errorf("SPHERE: radius must be positive and numeric")
	}
	seg := int64(16)
	if len(args) == 2 {
		s, ok := args[1].ToInt()
		if !ok || s < 3 {
			return value.Nil, fmt.Errorf("SPHERE: segments must be int >= 3")
		}
		seg = s
	}
	st := m.store()
	id := st.nextID
	st.nextID++
	e := newDefaultEnt(id)
	e.kind = entKindSphere
	e.radius = rad
	e.segH, e.segV = int32(seg), int32(seg)
	e.static = true
	e.w, e.h, e.d = rad*2, rad*2, rad*2
	st.ents[id] = e
	return m.wrapEntityRef(id)
}

// purgeEntityByID removes an entity from the store and unloads Raylib resources (shared by ENTITY.FREE and heap.EntityRef.Free).
func (m *Module) purgeEntityByID(id int64) {
	if id < 1 {
		return
	}
	st := m.store()
	e := st.ents[id]
	if e == nil {
		return
	}
	if e.hasRLModel {
		if len(e.modelAnims) > 0 {
			rl.UnloadModelAnimations(e.modelAnims)
			e.modelAnims = nil
		}
		rl.UnloadModel(e.rlModel)
	}
	if e.name != "" {
		delete(st.byName, strings.ToUpper(e.name))
	}
	delete(st.ents, id)
}
