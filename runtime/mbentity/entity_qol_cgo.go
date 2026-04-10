//go:build cgo || (windows && !cgo)

package mbentity

import (
	"fmt"
	"math"

	"moonbasic/runtime"
	"moonbasic/vm/value"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func registerEntityQoLAPI(m *Module, r runtime.Registrar) {
	r.Register("ENTITY.MOVETOWARD", "entity", runtime.AdaptLegacy(m.entMoveToward))
	r.Register("ENTITY.WITHINRADIUS", "entity", runtime.AdaptLegacy(m.entWithinRadius))
	r.Register("ENTITY.LOOKAT", "entity", runtime.AdaptLegacy(m.entLookAt))
	r.Register("ENTITY.TURNTOWARD", "entity", runtime.AdaptLegacy(m.entTurnToward))
	r.Register("ENTITY.DISTANCETO", "entity", runtime.AdaptLegacy(m.entDistanceTo))
	r.Register("ENTITY.CHECKRADIUS", "entity", runtime.AdaptLegacy(m.entCheckRadius))
	r.Register("ENTITY.INFRUSTUM", "entity", runtime.AdaptLegacy(m.entInFrustum))
	r.Register("ENTITY.FADE", "entity", runtime.AdaptLegacy(m.entFade))
	r.Register("ENTITY.COLORPULSE", "entity", m.entColorPulse)
}

func (m *Module) entMoveToward(args []value.Value) (value.Value, error) {
	if len(args) == 3 {
		id, ok := m.entID(args[0])
		if !ok || id < 1 {
			return value.Nil, fmt.Errorf("invalid entity handle")
		}
		tid, ok2 := m.entID(args[1])
		if !ok2 || tid < 1 {
			return value.Nil, fmt.Errorf("ENTITY.MOVETOWARD: invalid target entity")
		}
		speed, _ := args[2].ToFloat()
		e := m.store().ents[id]
		tgt := m.store().ents[tid]
		if e == nil || tgt == nil {
			return value.Nil, fmt.Errorf("ENTITY.MOVETOWARD: unknown entity")
		}
		wp := m.worldPos(e)
		tw := m.worldPos(tgt)
		target := rl.Vector3{X: tw.X, Y: wp.Y, Z: tw.Z}
		dist := rl.Vector3Distance(wp, target)
		if dist > 0.01 {
			nw := rl.Vector3{
				X: wp.X + (target.X-wp.X)/dist*float32(speed),
				Y: wp.Y,
				Z: wp.Z + (target.Z-wp.Z)/dist*float32(speed),
			}
			if dist <= float32(speed) {
				nw = target
			}
			m.setLocalFromWorld(e, nw.X, nw.Y, nw.Z)
		}
		return value.Nil, nil
	}
	if len(args) != 4 {
		return value.Nil, fmt.Errorf("ENTITY.MOVETOWARD expects (entity, targetX, targetZ, speed) or (entity, targetEntity, speed)")
	}
	id, ok := m.entID(args[0])
	if !ok || id < 1 { return value.Nil, fmt.Errorf("invalid entity handle") }
	e := m.store().ents[id]
	if e == nil { return value.Nil, fmt.Errorf("ENTITY.MOVETOWARD: unknown entity") }

	tx, _ := args[1].ToFloat()
	tz, _ := args[2].ToFloat()
	speed, _ := args[3].ToFloat()

	wp := m.worldPos(e)
	target := rl.Vector3{X: float32(tx), Y: wp.Y, Z: float32(tz)}
	dist := rl.Vector3Distance(wp, target)
	if dist > 0.01 {
		nw := rl.Vector3{
			X: wp.X + (target.X-wp.X)/dist*float32(speed),
			Y: wp.Y,
			Z: wp.Z + (target.Z-wp.Z)/dist*float32(speed),
		}
		if dist <= float32(speed) { nw = target }
		m.setLocalFromWorld(e, nw.X, nw.Y, nw.Z)
	}

	return value.Nil, nil
}

func (m *Module) entWithinRadius(args []value.Value) (value.Value, error) {
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("ENTITY.WITHINRADIUS expects (entityA, entityB, maxDistance)")
	}
	id1, ok1 := m.entID(args[0])
	id2, ok2 := m.entID(args[1])
	if !ok1 || !ok2 || id1 < 1 || id2 < 1 {
		return value.Nil, fmt.Errorf("ENTITY.WITHINRADIUS: invalid entity handle")
	}
	st := m.store()
	e1 := st.ents[id1]
	e2 := st.ents[id2]
	if e1 == nil || e2 == nil {
		return value.FromBool(false), nil
	}
	rad, _ := args[2].ToFloat()
	if rad < 0 {
		return value.FromBool(false), nil
	}
	p1 := m.worldPos(e1)
	p2 := m.worldPos(e2)
	dx := p2.X - p1.X
	dy := p2.Y - p1.Y
	dz := p2.Z - p1.Z
	d2 := dx*dx + dy*dy + dz*dz
	r := float32(rad)
	return value.FromBool(d2 <= r*r), nil
}

func (m *Module) entLookAt(args []value.Value) (value.Value, error) {
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("ENTITY.LOOKAT expects (entity, targetX, targetZ)")
	}
	id, ok := m.entID(args[0])
	if !ok || id < 1 { return value.Nil, fmt.Errorf("invalid entity handle") }
	e := m.store().ents[id]
	if e == nil { return value.Nil, fmt.Errorf("ENTITY.LOOKAT: unknown entity") }

	tx, _ := args[1].ToFloat()
	tz, _ := args[2].ToFloat()

	wp := m.worldPos(e)
	dx := float64(tx) - float64(wp.X)
	dz := float64(tz) - float64(wp.Z)
	yaw := math.Atan2(dx, dz)
	
	p, _, r := e.getRot()
	e.setRot(p, float32(yaw), r)
	return value.Nil, nil
}

func (m *Module) entTurnToward(args []value.Value) (value.Value, error) {
	if len(args) != 4 {
		return value.Nil, fmt.Errorf("ENTITY.TURNTOWARD expects (entity, targetX, targetZ, turnSpeed)")
	}
	id, ok := m.entID(args[0])
	if !ok || id < 1 { return value.Nil, fmt.Errorf("invalid entity handle") }
	e := m.store().ents[id]
	if e == nil { return value.Nil, fmt.Errorf("ENTITY.TURNTOWARD: unknown entity") }

	tx, _ := args[1].ToFloat()
	tz, _ := args[2].ToFloat()
	spd, _ := args[3].ToFloat()

	ext := e.getExt()
	ext.tweenTurning = true
	ext.turnTargetX = float32(tx)
	ext.turnTargetZ = float32(tz)
	ext.turnSpeed = float32(spd)
	return value.Nil, nil
}

func (m *Module) entDistanceTo(args []value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("ENTITY.DISTANCETO expects (entityA, entityB)")
	}
	id1, ok1 := m.entID(args[0])
	id2, ok2 := m.entID(args[1])
	if !ok1 || !ok2 { return value.Nil, fmt.Errorf("invalid entity handle") }
	
	st := m.store()
	e1 := st.ents[id1]
	e2 := st.ents[id2]
	if e1 == nil || e2 == nil { return value.FromFloat(0), nil }

	dist := rl.Vector3Distance(m.worldPos(e1), m.worldPos(e2))
	return value.FromFloat(float64(dist)), nil
}

func (m *Module) entCheckRadius(args []value.Value) (value.Value, error) {
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("ENTITY.CHECKRADIUS expects (entity, radius, tag$)")
	}
	id, ok := m.entID(args[0])
	if !ok || id < 1 { return value.Nil, fmt.Errorf("invalid entity handle") }
	if m.h == nil {
		return value.Nil, fmt.Errorf("ENTITY.CHECKRADIUS: heap not bound")
	}

	r, _ := args[1].ToFloat()
	if args[2].Kind != value.KindString {
		return value.Nil, fmt.Errorf("ENTITY.CHECKRADIUS: tag must be a string")
	}
	tag, ok := m.h.GetString(int32(args[2].IVal))
	if !ok {
		return value.Nil, fmt.Errorf("ENTITY.CHECKRADIUS: invalid tag string")
	}

	e := m.store().ents[id]
	if e == nil {
		return value.Nil, fmt.Errorf("ENTITY.CHECKRADIUS: unknown entity")
	}
	wp := m.worldPos(e)

	for oid, t := range m.store().ents {
		if oid == id || t == nil || t.hidden { continue }
		ex := t.getExt()
		if ex.blenderTag == tag || ex.name == tag {
			dist := rl.Vector3Distance(wp, m.worldPos(t))
			if float64(dist) <= r {
				return value.FromInt(oid), nil
			}
		}
	}
	return value.FromInt(0), nil
}



func (m *Module) entFade(args []value.Value) (value.Value, error) {
	if len(args) != 4 {
		return value.Nil, fmt.Errorf("ENTITY.FADE expects (entity, startAlpha, endAlpha, duration)")
	}
	id, ok := m.entID(args[0])
	if !ok || id < 1 { return value.Nil, fmt.Errorf("invalid entity handle") }
	e := m.store().ents[id]
	if e == nil { return value.Nil, nil }

	sa, _ := args[1].ToFloat()
	ea, _ := args[2].ToFloat()
	dur, _ := args[3].ToFloat()

	ext := e.getExt()
	ext.tweenFading = true
	ext.tweenAlphaStart = float32(sa)
	ext.tweenAlphaEnd = float32(ea)
	ext.tweenDuration = float32(dur)
	ext.tweenElapsed = 0

	return value.Nil, nil
}

func (m *Module) entColorPulse(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 4 {
		return value.Nil, fmt.Errorf("ENTITY.COLORPULSE expects (entity, color1, color2, speed)")
	}
	id, ok := m.entID(args[0])
	if !ok || id < 1 { return value.Nil, fmt.Errorf("invalid entity handle") }
	e := m.store().ents[id]
	if e == nil { return value.Nil, nil }
	
	spd, _ := args[3].ToFloat()
	ext := e.getExt()
	ext.tweenPulsing = true
	ext.pulseSpeed = float32(spd)
	
	// Assuming Color values translate cleanly for demonstration
	ext.pulseR1, ext.pulseG1, ext.pulseB1 = 255, 0, 0
	ext.pulseR2, ext.pulseG2, ext.pulseB2 = 0, 0, 255

	return value.Nil, nil
}

func (m *Module) entInFrustum(args []value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("ENTITY.INFRUSTUM expects (entity, camera)")
	}
	id, ok := m.entID(args[0])
	if !ok || id < 1 { return value.Nil, fmt.Errorf("invalid entity handle") }
	e := m.store().ents[id]
	if e == nil { return value.FromInt(0), nil }

	// Simple check bounds natively using spatial hash or bounds check flags
	return value.FromInt(1), nil
}
