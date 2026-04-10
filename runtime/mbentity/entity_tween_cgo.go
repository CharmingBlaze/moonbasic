//go:build cgo || (windows && !cgo)

package mbentity

import (
	"fmt"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

func registerEntityTweenAPI(m *Module, r runtime.Registrar) {
	r.Register("ENTITY.ANIMATETOWARD", "entity", runtime.AdaptLegacy(m.entAnimateToward))
	r.Register("EntityAnimateToward", "entity", runtime.AdaptLegacy(m.entAnimateToward))
}

// processEntityTweens advances ENTITY.ANIMATETOWARD lerps (call from ENTITY.UPDATE).
func (m *Module) processEntityTweens(dt float32) {
	if dt <= 0 {
		return
	}
	st := m.store()
	for _, e := range st.ents {
		if e == nil || !e.tweenActive {
			continue
		}
		e.tweenElapsed += dt
		u := float32(1)
		if e.tweenDuration > 1e-6 {
			u = e.tweenElapsed / e.tweenDuration
		}
		if u >= 1 {
			m.setLocalFromWorld(e, e.tweenTX, e.tweenTY, e.tweenTZ)
			e.tweenActive = false
			continue
		}
		x := e.tweenSX + (e.tweenTX-e.tweenSX)*u
		y := e.tweenSY + (e.tweenTY-e.tweenSY)*u
		z := e.tweenSZ + (e.tweenTZ-e.tweenSZ)*u
		m.setLocalFromWorld(e, x, y, z)
	}
}

func (m *Module) entAnimateToward(args []value.Value) (value.Value, error) {
	if len(args) != 5 {
		return value.Nil, fmt.Errorf("ENTITY.ANIMATETOWARD expects (entity#, x#, y#, z#, duration#)")
	}
	id, ok := m.entID(args[0])
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("ENTITY.ANIMATETOWARD: invalid entity")
	}
	e := m.store().ents[id]
	if e == nil {
		return value.Nil, fmt.Errorf("ENTITY.ANIMATETOWARD: unknown entity")
	}
	x, ok1 := args[1].ToFloat()
	y, ok2 := args[2].ToFloat()
	z, ok3 := args[3].ToFloat()
	dur, ok4 := args[4].ToFloat()
	if !ok1 || !ok2 || !ok3 || !ok4 || dur <= 0 {
		return value.Nil, fmt.Errorf("ENTITY.ANIMATETOWARD: bad numeric args (duration must be > 0)")
	}
	wp := m.worldPos(e)
	e.tweenSX, e.tweenSY, e.tweenSZ = wp.X, wp.Y, wp.Z
	e.tweenTX, e.tweenTY, e.tweenTZ = float32(x), float32(y), float32(z)
	e.tweenElapsed = 0
	e.tweenDuration = float32(dur)
	e.tweenActive = true
	return value.Nil, nil
}
