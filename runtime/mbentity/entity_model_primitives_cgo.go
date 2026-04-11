//go:build cgo || (windows && !cgo)

package mbentity

import (
	"fmt"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

// MODEL.CREATECAPSULE / MODEL.CREATEBOX are entity constructors (EntityRef) for Modern Blitz-style samples.
func registerModelEntityPrimitives(m *Module, r runtime.Registrar) {
	r.Register("MODEL.CREATECAPSULE", "entity", runtime.AdaptLegacy(m.modelCreateCapsule))
	r.Register("MODEL.CREATEBOX", "entity", runtime.AdaptLegacy(m.modelCreateBox))
}

func (m *Module) modelCreateCapsule(args []value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("MODEL.CREATECAPSULE expects 2 arguments (radius#, height#)")
	}
	rad, ok1 := argF32(args[0])
	h, ok2 := argF32(args[1])
	if !ok1 || !ok2 || rad <= 0 || h <= 0 {
		return value.Nil, fmt.Errorf("MODEL.CREATECAPSULE: positive radius and height required")
	}
	st := m.store()
	id := st.nextID
	st.nextID++
	st.ensureSlices(int(id))
	e := newDefaultEnt(id, &st.spatial)
	e.kind = entKindCapsule
	e.radius = rad
	e.cylH = h
	e.segV = 16
	e.w, e.h, e.d = rad*2, h, rad*2
	e.static = false
	e.useSphere = true
	e.physBottomOffset = rad + h*0.5
	e.gravity = -28
	st.ents[id] = e
	st.dynamicEnts = append(st.dynamicEnts, e)
	return m.wrapEntityRef(id)
}

func (m *Module) modelCreateBox(args []value.Value) (value.Value, error) {
	var w, h, d float32
	var staticEnt bool = true
	switch len(args) {
	case 3:
		var ok1, ok2, ok3 bool
		w, ok1 = argF32(args[0])
		h, ok2 = argF32(args[1])
		d, ok3 = argF32(args[2])
		if !ok1 || !ok2 || !ok3 {
			return value.Nil, fmt.Errorf("MODEL.CREATEBOX: dimensions must be numeric")
		}
	case 4:
		var ok1, ok2, ok3 bool
		w, ok1 = argF32(args[0])
		h, ok2 = argF32(args[1])
		d, ok3 = argF32(args[2])
		if !ok1 || !ok2 || !ok3 {
			return value.Nil, fmt.Errorf("MODEL.CREATEBOX: dimensions must be numeric")
		}
		staticEnt = value.Truthy(args[3], nil, m.h)
	default:
		return value.Nil, fmt.Errorf("MODEL.CREATEBOX expects 3 (w#, h#, d#) or 4 (w#, h#, d#, static?)")
	}
	st := m.store()
	id := st.nextID
	st.nextID++
	st.ensureSlices(int(id))
	e := newDefaultEnt(id, &st.spatial)
	e.kind = entKindBox
	e.r, e.g, e.b = 180, 180, 200
	e.w, e.h, e.d = w, h, d
	e.static = staticEnt
	e.useSphere = false
	e.physBottomOffset = h * 0.5
	e.gravity = 0
	st.ents[id] = e
	if staticEnt {
		st.staticEnts = append(st.staticEnts, e)
	} else {
		st.dynamicEnts = append(st.dynamicEnts, e)
	}
	return m.wrapEntityRef(id)
}
