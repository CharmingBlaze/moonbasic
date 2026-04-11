//go:build cgo || (windows && !cgo)

package mbentity

import (
	"fmt"

	"moonbasic/runtime"
	"moonbasic/vm/value"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func registerEntityPhysicsQoLAPI(m *Module, r runtime.Registrar) {
	r.Register("PHYSICS.EXPLOSION", "entity", runtime.AdaptLegacy(m.physExplosion))
	r.Register("WORLD.EXPLOSION", "entity", runtime.AdaptLegacy(m.physExplosion))
	r.Register("ENTITY.SETWEIGHT", "entity", runtime.AdaptLegacy(m.physSetWeight))
	r.Register("ENTITY.APPLYTORQUE", "entity", runtime.AdaptLegacy(m.physApplyTorque))
	r.Register("ENTITY.ONHIT", "entity", m.physOnHit)
	r.Register("ENTITY.GHOSTMODE", "entity", runtime.AdaptLegacy(m.physGhostMode))
	r.Register("ENTITY.EXPLODE", "entity", runtime.AdaptLegacy(m.physExplode))
}

func (m *Module) physExplosion(args []value.Value) (value.Value, error) {
	if len(args) != 5 {
		return value.Nil, fmt.Errorf("PHYSICS.EXPLOSION expects (x, y, z, force, radius)")
	}
	// Coordinates
	x, _ := args[0].ToFloat()
	y, _ := args[1].ToFloat()
	z, _ := args[2].ToFloat()
	force, _ := args[3].ToFloat()
	rad, _ := args[4].ToFloat()

	epicenter := rl.Vector3{X: float32(x), Y: float32(y), Z: float32(z)}

	st := m.store()
	for _, e := range st.ents {
		if e == nil || e.hidden || e.static { continue }
		wp := m.worldPos(e)
		dist := rl.Vector3Distance(epicenter, wp)
		
		if float64(dist) < rad {
			// Radial impulse mapping dynamically
			inf := 1.0 - (float64(dist) / rad) // Attenuation curve
			dir := rl.Vector3{X: wp.X - epicenter.X, Y: wp.Y - epicenter.Y + 0.5, Z: wp.Z - epicenter.Z} // Upward bias
			dir = rl.Vector3Normalize(dir)
			
			// If physBufIndex >= 0 (Jolt), we can pass directly to a backend Jolt method.
			// For demonstration, we simply map it to `e.vel` applying native acceleration.
			e.vel.X += dir.X * float32(force*inf)
			e.vel.Y += dir.Y * float32(force*inf)
			e.vel.Z += dir.Z * float32(force*inf)
		}
	}
	return value.Nil, nil
}

func (m *Module) physSetWeight(args []value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("ENTITY.SETWEIGHT expects (entity, mass)")
	}
	id, ok := m.entID(args[0])
	if !ok || id < 1 { return value.Nil, fmt.Errorf("invalid entity handle") }
	e := m.store().ents[id]
	if e == nil { return value.Nil, nil }
	
	mass, _ := args[1].ToFloat()
	e.mass = float32(mass)
	// Would sync to Jolt via `joltwasm.SetBodyMass(id, e.mass)` natively.
	return value.Nil, nil
}

func (m *Module) physApplyTorque(args []value.Value) (value.Value, error) {
	if len(args) != 4 {
		return value.Nil, fmt.Errorf("ENTITY.APPLYTORQUE expects (entity, tx, ty, tz)")
	}
	id, ok := m.entID(args[0])
	if !ok || id < 1 { return value.Nil, fmt.Errorf("invalid entity handle") }
	e := m.store().ents[id]
	if e == nil { return value.Nil, nil }
	
	// tx, _ := args[1].ToFloat()
	// ty, _ := args[2].ToFloat()
	// tz, _ := args[3].ToFloat()
	
	// Simulating native rotation impulse mappings for simple objects:
	// e.roll += float32(tx)
	// e.yaw += float32(ty)
	// e.pitch += float32(tz)
	return value.Nil, nil
}

func (m *Module) physOnHit(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("ENTITY.ONHIT expects (entity, macroFuncName)")
	}
	id, ok := m.entID(args[0])
	if !ok || id < 1 { return value.Nil, fmt.Errorf("invalid entity handle") }
	e := m.store().ents[id]
	if e == nil { return value.Nil, nil }
	
	funcStr, _ := rt.ArgString(args, 1)
	e.getExt().onHitAction = funcStr
	return value.Nil, nil
}

func (m *Module) physGhostMode(args []value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("ENTITY.GHOSTMODE expects (entity, duration)")
	}
	id, ok := m.entID(args[0])
	if !ok || id < 1 { return value.Nil, fmt.Errorf("invalid entity handle") }
	e := m.store().ents[id]
	if e == nil { return value.Nil, nil }
	
	// dur, _ := args[1].ToFloat()
	// Toggle Jolt collision group to sensors or none
	// e.ghostMode = true
	// e.ghostTimer = float32(dur)
	return value.Nil, nil
}

func (m *Module) physExplode(args []value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("ENTITY.EXPLODE expects (entity, pieces)")
	}
	id, ok := m.entID(args[0])
	if !ok || id < 1 { return value.Nil, fmt.Errorf("invalid entity handle") }
	e := m.store().ents[id]
	if e == nil { return value.Nil, nil }
	
	// pieces, _ := args[1].ToInt()
	// Unparent sub-nodes and apply random impulses
	wp := m.worldPos(e)
	m.physExplosion([]value.Value{
		value.FromFloat(float64(wp.X)),
		value.FromFloat(float64(wp.Y)),
		value.FromFloat(float64(wp.Z)),
		value.FromFloat(10.0), // force
		value.FromFloat(5.0),  // radius
	})
	
	e.hidden = true // Base entity poofs
	return value.Nil, nil
}
