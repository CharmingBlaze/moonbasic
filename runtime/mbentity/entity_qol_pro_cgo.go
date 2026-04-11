//go:build cgo || (windows && !cgo)

package mbentity

import (
	"fmt"

	"moonbasic/runtime"
	mbmatrix "moonbasic/runtime/mbmatrix"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func registerEntityQoLProAPI(m *Module, r runtime.Registrar) {
	r.Register("ENTITY.OUTLINE", "entity", runtime.AdaptLegacy(m.entOutline))
	r.Register("ENTITY.SNAPTO", "entity", runtime.AdaptLegacy(m.entSnapTo))
}

func (m *Module) entOutline(args []value.Value) (value.Value, error) {
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("ENTITY.OUTLINE expects (entity, thickness, color)")
	}
	id, ok := m.entID(args[0])
	if !ok || id < 1 { return value.Nil, fmt.Errorf("invalid entity handle") }
	e := m.store().ents[id]
	if e == nil { return value.Nil, nil }

	thick, _ := args[1].ToFloat()
	colorH := heap.Handle(args[2].IVal)
	
	ext := e.getExt()
	ext.outlineThickness = float32(thick)
	if colorH != 0 {
		col, err := mbmatrix.GetColor(m.h, colorH)
		if err == nil {
			ext.outlineColor = col
		} else {
			ext.outlineColor = rl.Black
		}
	} else {
		ext.outlineColor = rl.Black
	}
	
	return value.Nil, nil
}

func (m *Module) entSnapTo(args []value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("ENTITY.SNAPTO expects (entity, targetEntity)")
	}
	id1, ok1 := m.entID(args[0])
	id2, ok2 := m.entID(args[1])
	if !ok1 || !ok2 || id1 < 1 || id2 < 1 {
		return value.Nil, fmt.Errorf("ENTITY.SNAPTO: invalid entity handles")
	}

	st := m.store()
	e1 := st.ents[id1]
	e2 := st.ents[id2]
	if e1 == nil || e2 == nil { return value.Nil, nil }

	wp2 := m.worldPos(e2)
	m.setLocalFromWorld(e1, wp2.X, wp2.Y, wp2.Z)
	
	p2, w2, r2 := e2.getRot()
	e1.setRot(p2, w2, r2)

	return value.Nil, nil
}
