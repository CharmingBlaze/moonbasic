//go:build cgo || (windows && !cgo)

package mbentity

import (
	"fmt"

	"moonbasic/runtime"
	"moonbasic/runtime/mbtween"
	"moonbasic/vm/value"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func registerEntityTweenAPI(m *Module, r runtime.Registrar) {
	r.Register("ENTITY.ANIMATETOWARD", "entity", runtime.AdaptLegacy(m.entAnimateToward))
	r.Register("ENTITY.TWEEN", "entity", m.entTween)
}

// processEntityTweens advances ENTITY.ANIMATETOWARD as well as new complex ENTITY.TWEEN calls.
func (m *Module) processEntityTweens(dt float32) {
	if dt <= 0 { return }
	st := m.store()
	for _, e := range st.ents {
		if e == nil || e.ext == nil { continue }
		ext := e.ext

		// 1. Legacy ENTITY.ANIMATETOWARD logic
		if ext.tweenActive {
			ext.tweenElapsed += dt
			u := float32(1)
			if ext.tweenDuration > 1e-6 {
				u = ext.tweenElapsed / ext.tweenDuration
			}
			if u >= 1 {
				m.setLocalFromWorld(e, ext.tweenTX, ext.tweenTY, ext.tweenTZ)
				ext.tweenActive = false
			} else {
				x := ext.tweenSX + (ext.tweenTX-ext.tweenSX)*u
				y := ext.tweenSY + (ext.tweenTY-ext.tweenSY)*u
				z := ext.tweenSZ + (ext.tweenTZ-ext.tweenSZ)*u
				m.setLocalFromWorld(e, x, y, z)
			}
			if ext.tweenFading {
				e.alpha = ext.tweenAlphaStart + (ext.tweenAlphaEnd-ext.tweenAlphaStart)*u
				if u >= 1 { ext.tweenFading = false }
			}
		}

		// 2. Complex persistent tweens (multi-property)
		if len(ext.complexTweens) > 0 {
			active := ext.complexTweens[:0]
			for i := range ext.complexTweens {
				tw := &ext.complexTweens[i]
				tw.elapsed += dt
				p := float64(1.0)
				if tw.duration > 1e-6 {
					p = float64(tw.elapsed / tw.duration)
				}
				if p > 1.0 { p = 1.0 }
				
				val := tw.start + (tw.target-tw.start)*float32(mbtween.Ease(p, tw.ease))
				m.setEntityProperty(e, tw.prop, val)

				if p < 1.0 {
					active = append(active, *tw)
				}
			}
			ext.complexTweens = active
		}

		if ext.tweenPulsing {
			ext.pulseT += dt * ext.pulseSpeed
		}
	}
}

func (m *Module) setEntityProperty(e *ent, prop string, val float32) {
	switch prop {
	case "x": 
		pos := e.getPos()
		pos.X = val
		e.setPos(pos)
	case "y": 
		pos := e.getPos()
		pos.Y = val
		e.setPos(pos)
	case "z": 
		pos := e.getPos()
		pos.Z = val
		e.setPos(pos)
	case "scale": e.scale = rl.Vector3{X: val, Y: val, Z: val}
	case "scale_x": e.scale.X = val
	case "scale_y": e.scale.Y = val
	case "scale_z": e.scale.Z = val
	case "rotation_x", "pitch":
		_, w, r := e.getRot()
		e.setRot(val, w, r)
	case "rotation_y", "yaw":
		p, _, r := e.getRot()
		e.setRot(p, val, r)
	case "rotation_z", "roll":
		p, w, _ := e.getRot()
		e.setRot(p, w, val)
	case "alpha": e.alpha = val
	case "r": e.r = uint8(val)
	case "g": e.g = uint8(val)
	case "b": e.b = uint8(val)
	}
}

func (m *Module) entTween(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 5 {
		return value.Nil, fmt.Errorf("ENTITY.TWEEN expects (entity#, property$, targetValue, duration, easeType$)")
	}
	id, ok := m.entID(args[0])
	if !ok || id < 1 { return value.Nil, fmt.Errorf("invalid entity handle") }
	e := m.store().ents[id]
	if e == nil { return value.Nil, fmt.Errorf("unknown entity") }

	prop, _ := rt.ArgString(args, 1)
	target, _ := args[2].ToFloat()
	dur, _ := args[3].ToFloat()
	ease, _ := rt.ArgString(args, 4)

	start := m.getEntityProperty(e, prop)
	
	ext := e.getExt()
	// Remove existing tween for same property
	active := ext.complexTweens[:0]
	for _, tw := range ext.complexTweens {
		if tw.prop != prop {
			active = append(active, tw)
		}
	}
	ext.complexTweens = append(active, complexTween{
		prop:     prop,
		start:    start,
		target:   float32(target),
		duration: float32(dur),
		ease:     ease,
	})
	return value.Nil, nil
}

func (m *Module) getEntityProperty(e *ent, prop string) float32 {
	switch prop {
	case "x": return e.getPos().X
	case "y": return e.getPos().Y
	case "z": return e.getPos().Z
	case "scale": return e.scale.X
	case "scale_x": return e.scale.X
	case "scale_y": return e.scale.Y
	case "scale_z": return e.scale.Z
	case "rotation_x", "pitch":
		p, _, _ := e.getRot()
		return p
	case "rotation_y", "yaw":
		_, w, _ := e.getRot()
		return w
	case "rotation_z", "roll":
		_, _, r := e.getRot()
		return r
	case "alpha": return e.alpha
	case "r": return float32(e.r)
	case "g": return float32(e.g)
	case "b": return float32(e.b)
	}
	return 0
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
	ext := e.getExt()
	ext.tweenSX, ext.tweenSY, ext.tweenSZ = wp.X, wp.Y, wp.Z
	ext.tweenTX, ext.tweenTY, ext.tweenTZ = float32(x), float32(y), float32(z)
	ext.tweenElapsed = 0
	ext.tweenDuration = float32(dur)
	ext.tweenActive = true
	return value.Nil, nil
}
