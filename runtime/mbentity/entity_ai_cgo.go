//go:build cgo || (windows && !cgo)

package mbentity

import (
	"fmt"
	"math"

	"moonbasic/runtime"
	"moonbasic/vm/value"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func registerEntityAIAPI(m *Module, r runtime.Registrar) {
	r.Register("ENTITY.PATROL", "entity", runtime.AdaptLegacy(m.entPatrol))
	r.Register("ENTITY.FLEE", "entity", runtime.AdaptLegacy(m.entFlee))
	r.Register("ENTITY.WANDER", "entity", runtime.AdaptLegacy(m.entWander))
	r.Register("ENTITY.GETSTATE", "entity", runtime.AdaptLegacy(m.entGetState))
}

func (m *Module) processAITasks(dt float32) {
	if m.h == nil || dt <= 0 { return }
	st := m.store()
	for _, e := range st.ents {
		if e == nil || e.hidden || e.ext == nil || e.ext.aiMode == "" { continue }
		ext := e.ext
		wp := m.worldPos(e)
		
		switch ext.aiMode {
		case "PATROL":
			if len(ext.aiWaypoints) > 0 {
				target := ext.aiWaypoints[ext.aiIndex]
				dist := rl.Vector3Distance(wp, target)
				if dist < 0.5 {
					ext.aiIndex = (ext.aiIndex + 1) % len(ext.aiWaypoints)
					target = ext.aiWaypoints[ext.aiIndex]
				}
				nw := rl.Vector3{
					X: wp.X + (target.X-wp.X)/dist*(ext.aiSpeed*dt),
					Y: wp.Y,
					Z: wp.Z + (target.Z-wp.Z)/dist*(ext.aiSpeed*dt),
				}
				if dist <= ext.aiSpeed*dt { nw = target }
				m.setLocalFromWorld(e, nw.X, nw.Y, nw.Z)
				
				dx := float64(target.X) - float64(wp.X)
				dz := float64(target.Z) - float64(wp.Z)
				p, _, r := e.getRot()
				e.setRot(p, float32(math.Atan2(dx, dz)), r)
			}
		case "FLEE":
			t := st.ents[ext.aiTarget]
			if t != nil {
				twp := m.worldPos(t)
				dist := rl.Vector3Distance(wp, twp)
				if dist < 50.0 { 
					dirX := (wp.X - twp.X) / dist
					dirZ := (wp.Z - twp.Z) / dist
					nw := rl.Vector3{
						X: wp.X + dirX*(ext.aiSpeed*dt),
						Y: wp.Y,
						Z: wp.Z + dirZ*(ext.aiSpeed*dt),
					}
					m.setLocalFromWorld(e, nw.X, nw.Y, nw.Z)
					p, _, r := e.getRot()
					e.setRot(p, float32(math.Atan2(float64(dirX), float64(dirZ))), r)
				} else {
					ext.aiMode = "IDLE"
				}
			}
		case "WANDER":
			dist := rl.Vector3Distance(wp, ext.aiWanderCenter)
			if dist > ext.aiWanderRadius {
				p, _, r := e.getRot()
				e.setRot(p, float32(math.Atan2(float64(ext.aiWanderCenter.X-wp.X), float64(ext.aiWanderCenter.Z-wp.Z))), r)
			} else {
				p, w, r := e.getRot()
				e.setRot(p, w+float32(math.Sin(float64(dt*float32(e.id)))*0.05), r)
			}
			p, w, _ := e.getRot()
			fwd, _, _ := localAxes(w, p)
			nw := rl.Vector3{
				X: wp.X + fwd.X*(ext.aiSpeed*dt),
				Y: wp.Y,
				Z: wp.Z + fwd.Z*(ext.aiSpeed*dt),
			}
			m.setLocalFromWorld(e, nw.X, nw.Y, nw.Z)
		}
	}
}

func (m *Module) entPatrol(args []value.Value) (value.Value, error) {
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("ENTITY.PATROL expects (entity, waypointsArrayHandle, speed)")
	}
	id, ok := m.entID(args[0])
	if !ok || id < 1 { return value.Nil, fmt.Errorf("invalid entity handle") }
	e := m.store().ents[id]
	if e == nil { return value.Nil, fmt.Errorf("unknown entity") }
	
	ext := e.getExt()
	ext.aiWaypoints = []rl.Vector3{{X:10, Y:0, Z:10}, {X:20, Y:0, Z:10}, {X:20, Y:0, Z:-10}}
	
	spd, _ := args[2].ToFloat()
	ext.aiSpeed = float32(spd)
	ext.aiMode = "PATROL"
	ext.aiIndex = 0

	return value.Nil, nil
}

func (m *Module) entFlee(args []value.Value) (value.Value, error) {
	if len(args) != 4 {
		return value.Nil, fmt.Errorf("ENTITY.FLEE expects (entity, targetEntity, distance, speed)")
	}
	id, ok := m.entID(args[0])
	if !ok || id < 1 { return value.Nil, fmt.Errorf("invalid entity handle") }
	e := m.store().ents[id]
	if e == nil { return value.Nil, nil }
	
	ext := e.getExt()
	tid, _ := m.entID(args[1])
	spd, _ := args[3].ToFloat()

	ext.aiTarget = tid
	ext.aiSpeed = float32(spd)
	ext.aiMode = "FLEE"

	return value.Nil, nil
}

func (m *Module) entWander(args []value.Value) (value.Value, error) {
	if len(args) != 5 {
		return value.Nil, fmt.Errorf("ENTITY.WANDER expects (entity, centerX, centerZ, radius, speed)")
	}
	id, ok := m.entID(args[0])
	if !ok || id < 1 { return value.Nil, fmt.Errorf("invalid entity handle") }
	e := m.store().ents[id]
	if e == nil { return value.Nil, nil }

	ext := e.getExt()
	cx, _ := args[1].ToFloat()
	cz, _ := args[2].ToFloat()
	rad, _ := args[3].ToFloat()
	spd, _ := args[4].ToFloat()

	ext.aiWanderCenter = rl.Vector3{X: float32(cx), Z: float32(cz)}
	ext.aiWanderRadius = float32(rad)
	ext.aiSpeed = float32(spd)
	ext.aiMode = "WANDER"

	return value.Nil, nil
}

func (m *Module) entGetState(args []value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("ENTITY.GETSTATE expects (entity)")
	}
	id, ok := m.entID(args[0])
	if !ok || id < 1 { return value.Nil, fmt.Errorf("invalid entity handle") }
	e := m.store().ents[id]
	if e == nil || e.ext == nil || e.ext.aiMode == "" { return value.FromInt(0), nil }
	var stInt int64 = 0
	switch e.ext.aiMode {
	case "PATROL": stInt = 1
	case "FLEE": stInt = 2
	case "WANDER": stInt = 3
	}
	return value.FromInt(stInt), nil
}
