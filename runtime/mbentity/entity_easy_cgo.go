//go:build cgo || (windows && !cgo)

package mbentity

import (
	"fmt"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

func registerEntityEasyAPI(m *Module, r runtime.Registrar) {
	// Simple transform shorthands
	r.Register("ENTITY.POS", "entity", runtime.AdaptLegacy(m.entSetPosition))
	r.Register("ENTITY.ROT", "entity", runtime.AdaptLegacy(m.entRotateEntityAbs))
	r.Register("ENTITY.SCA", "entity", runtime.AdaptLegacy(m.entScaleEntity))

	// Component-level getters/setters (Easy Mode)
	r.Register("ENTITY.X", "entity", runtime.AdaptLegacy(m.entX))
	r.Register("ENTITY.Y", "entity", runtime.AdaptLegacy(m.entY))
	r.Register("ENTITY.Z", "entity", runtime.AdaptLegacy(m.entZ))
	r.Register("ENTITY.P", "entity", runtime.AdaptLegacy(m.entP))
	r.Register("ENTITY.W", "entity", runtime.AdaptLegacy(m.entYaw))
	r.Register("ENTITY.R", "entity", runtime.AdaptLegacy(m.entR))
	
	// Aesthetic shorthands
	r.Register("ENTITY.RGB", "entity", runtime.AdaptLegacy(m.entRGB))
	r.Register("ENTITY.ALPHA", "entity", runtime.AdaptLegacy(m.entAlpha))
}

func (m *Module) entX(args []value.Value) (value.Value, error) {
	if len(args) < 1 { return value.Nil, fmt.Errorf("ENTITY.X(id [, val#])") }
	id, ok := m.entID(args[0])
	if !ok || id < 1 { return value.FromFloat(0), nil }
	e := m.store().ents[id]
	if e == nil { return value.FromFloat(0), nil }

	if len(args) >= 2 {
		val, _ := args[1].ToFloat()
		pos := e.getPos()
		pos.X = float32(val)
		e.setPos(pos)
		return value.Nil, nil
	}
	return value.FromFloat(float64(e.getPos().X)), nil
}

func (m *Module) entY(args []value.Value) (value.Value, error) {
	if len(args) < 1 { return value.Nil, fmt.Errorf("ENTITY.Y(id [, val#])") }
	id, ok := m.entID(args[0])
	if !ok || id < 1 { return value.FromFloat(0), nil }
	e := m.store().ents[id]
	if e == nil { return value.FromFloat(0), nil }

	if len(args) >= 2 {
		val, _ := args[1].ToFloat()
		pos := e.getPos()
		pos.Y = float32(val)
		e.setPos(pos)
		return value.Nil, nil
	}
	return value.FromFloat(float64(e.getPos().Y)), nil
}

func (m *Module) entZ(args []value.Value) (value.Value, error) {
	if len(args) < 1 { return value.Nil, fmt.Errorf("ENTITY.Z(id [, val#])") }
	id, ok := m.entID(args[0])
	if !ok || id < 1 { return value.FromFloat(0), nil }
	e := m.store().ents[id]
	if e == nil { return value.FromFloat(0), nil }

	if len(args) >= 2 {
		val, _ := args[1].ToFloat()
		pos := e.getPos()
		pos.Z = float32(val)
		e.setPos(pos)
		return value.Nil, nil
	}
	return value.FromFloat(float64(e.getPos().Z)), nil
}

func (m *Module) entP(args []value.Value) (value.Value, error) {
	if len(args) < 1 { return value.Nil, fmt.Errorf("ENTITY.P(id [, val#])") }
	id, ok := m.entID(args[0])
	if !ok || id < 1 { return value.FromFloat(0), nil }
	e := m.store().ents[id]
	if e == nil { return value.FromFloat(0), nil }

	if len(args) >= 2 {
		val, _ := args[1].ToFloat()
		_, w, r := e.getRot()
		e.setRot(float32(val), w, r)
		return value.Nil, nil
	}
	p, _, _ := e.getRot()
	return value.FromFloat(float64(p)), nil
}

func (m *Module) entYaw(args []value.Value) (value.Value, error) {
	if len(args) < 1 { return value.Nil, fmt.Errorf("ENTITY.W(id [, val#])") }
	id, ok := m.entID(args[0])
	if !ok || id < 1 { return value.FromFloat(0), nil }
	e := m.store().ents[id]
	if e == nil { return value.FromFloat(0), nil }

	if len(args) >= 2 {
		val, _ := args[1].ToFloat()
		p, _, r := e.getRot()
		e.setRot(p, float32(val), r)
		return value.Nil, nil
	}
	_, w, _ := e.getRot()
	return value.FromFloat(float64(w)), nil
}

func (m *Module) entR(args []value.Value) (value.Value, error) {
	if len(args) < 1 { return value.Nil, fmt.Errorf("ENTITY.R(id [, val#])") }
	id, ok := m.entID(args[0])
	if !ok || id < 1 { return value.FromFloat(0), nil }
	e := m.store().ents[id]
	if e == nil { return value.FromFloat(0), nil }

	if len(args) >= 2 {
		val, _ := args[1].ToFloat()
		p, w, _ := e.getRot()
		e.setRot(p, w, float32(val))
		return value.Nil, nil
	}
	_, _, r := e.getRot()
	return value.FromFloat(float64(r)), nil
}

func (m *Module) entRGB(args []value.Value) (value.Value, error) {
	if len(args) != 4 { return value.Nil, fmt.Errorf("ENTITY.RGB(id, r, g, b)") }
	// Use m.entColor(args) which handles EntityColor(id, r, g, b)
	return m.entColor(args)
}
