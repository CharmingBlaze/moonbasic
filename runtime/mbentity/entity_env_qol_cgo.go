//go:build cgo || (windows && !cgo)

package mbentity

import (
	"fmt"

	"moonbasic/runtime"
	"moonbasic/vm/value"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func registerEntityEnvQoLAPI(m *Module, r runtime.Registrar) {
	r.Register("WORLD.DAYNIGHTCYCLE", "entity", runtime.AdaptLegacy(m.worldDayNightCycle))
	r.Register("TRIGGER.CREATEZONE", "entity", m.triggerCreateZone)
	r.Register("ENTITY.ATTACH", "entity", runtime.AdaptLegacy(m.entAttach))
	r.Register("ENTITY.EMITPARTICLES", "entity", runtime.AdaptLegacy(m.entEmitParticles))
	r.Register("CAMERA.FOLLOW", "entity", runtime.AdaptLegacy(m.camFollow))
	r.Register("WORLD.SCREENSHAKE", "entity", m.worldScreenShake)
	r.Register("WORLD.SHAKE", "entity", m.worldScreenShake)
	r.Register("WORLD.HITSTOP", "entity", m.worldHitStop)
}

func (m *Module) worldDayNightCycle(args []value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("WORLD.DAYNIGHTCYCLE expects (durationSeconds#)")
	}
	// Implementation binds to mblight/sun if available.
	// We'd store this in global runtime or mblight state to advance smoothly.
	// For compilation limits: stubbed
	return value.Nil, nil
}

func (m *Module) triggerCreateZone(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 7 {
		return value.Nil, fmt.Errorf("TRIGGER.CREATEZONE expects (x, y, z, w, h, d, tag$)")
	}
	x, _ := args[0].ToFloat()
	y, _ := args[1].ToFloat()
	z, _ := args[2].ToFloat()
	w, _ := args[3].ToFloat()
	h_sz, _ := args[4].ToFloat()
	d, _ := args[5].ToFloat()
	tag, _ := rt.ArgString(args, 6)

	// Creates a hidden SENSOR entity box natively
	st := m.store()
	nid := st.nextID
	st.nextID++
	st.ensureSlices(int(nid))
	ent := newDefaultEnt(nid, &st.spatial)
	ent.kind = entKindBox
	ent.setPos(rl.Vector3{X: float32(x), Y: float32(y), Z: float32(z)})
	ent.w = float32(w)
	ent.h = float32(h_sz)
	ent.d = float32(d)
	ent.hidden = true
	ent.static = true
	ent.getExt().blenderTag = tag
	ent.getExt().collType = 99 // Sensor bounds

	m.store().ents[nid] = ent
	return value.FromInt(nid), nil
}

func (m *Module) entAttach(args []value.Value) (value.Value, error) {
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("ENTITY.ATTACH expects (child#, parent#, offsetArrayHandle)")
	}
	child, ok1 := m.entID(args[0])
	parent, ok2 := m.entID(args[1])
	if !ok1 || !ok2 || child < 1 || parent < 1 {
		return value.Nil, fmt.Errorf("ENTITY.ATTACH: invalid entities")
	}
	
	// Uses internal Parent binding
	st := m.store()
	st.ents[child].getExt().parentID = parent
	return value.Nil, nil
}

func (m *Module) entEmitParticles(args []value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("ENTITY.EMITPARTICLES expects (entity#, template#)")
	}
	id, ok := m.entID(args[0])
	if !ok || id < 1 { return value.Nil, fmt.Errorf("invalid entity handle") }
	e := m.store().ents[id]
	if e == nil { return value.Nil, nil }
	
	// Mock: binds the active tracking natively passing execution logic out.
	// Ex: mbparticles.Attach(id, templateHandle)
	return value.Nil, nil
}

func (m *Module) camFollow(args []value.Value) (value.Value, error) {
	if len(args) != 4 {
		return value.Nil, fmt.Errorf("CAMERA.FOLLOW expects (camera#, entity#, distance#, smoothness#)")
	}
	return m.camFollowEntity(nil, args[0], args[1], args[2], value.FromFloat(10), args[3])
}

// worldHitStop freezes gameplay time for a wall-clock duration (impact frames); see Registry.HitStopEndAt + mbtime.DeltaSeconds.
func (m *Module) worldHitStop(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("WORLD.HITSTOP expects (durationSeconds#)")
	}
	if rt == nil {
		return value.Nil, fmt.Errorf("WORLD.HITSTOP: runtime not available")
	}
	sec, ok := args[0].ToFloat()
	if !ok || sec < 0 {
		return value.Nil, fmt.Errorf("WORLD.HITSTOP: duration must be a non-negative number")
	}
	rt.HitStopEndAt = float64(rl.GetTime()) + sec
	return value.Nil, nil
}

func (m *Module) worldScreenShake(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("WORLD.SCREENSHAKE expects (intensity, duration)")
	}
	if m.cam == nil {
		return value.Nil, fmt.Errorf("WORLD.SCREENSHAKE: camera module not bound")
	}
	camH := m.cam.ActiveCameraHandle()
	if camH == 0 {
		return value.Nil, nil // No active camera to shake
	}
	
	// Call CAMERA.SHAKE(camH, intensity, duration)
	return rt.Call("CAMERA.SHAKE", []value.Value{
		value.FromHandle(camH),
		args[0],
		args[1],
	})
}

