//go:build (linux || windows) && cgo

package mbentity

import (
	"fmt"
	"strings"

	"github.com/bbitechnologies/jolt-go/jolt"
	"moonbasic/runtime"
	mbphysics3d "moonbasic/runtime/physics3d"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func registerEntityPhysicsMacroAPI(m *Module, r runtime.Registrar) {
	// One-line setup
	r.Register("ENTITY.PHYSICS", "entity", runtime.AdaptLegacy(m.entPhysicsAuto))
	r.Register("ENTITY.ADDPHYSICS", "entity", runtime.AdaptLegacy(m.entAddPhysics))
	r.Register("PHYSICS.AUTO", "entity", runtime.AdaptLegacy(m.entPhysicsAuto))
	r.Register("ENTITY.PHYSICSMOTION", "entity", runtime.AdaptLegacy(m.entPhysicsMotion))

	// Property-based Builder (Breaks up long argument chains)
	r.Register("PHYSICS.SHAPE", "entity", runtime.AdaptLegacy(m.physSetShape))
	r.Register("PHYSICS.SIZE", "entity", runtime.AdaptLegacy(m.physSetSize))
	r.Register("PHYSICS.FRICTION", "entity", runtime.AdaptLegacy(m.physSetFriction))
	r.Register("PHYSICS.BOUNCE", "entity", runtime.AdaptLegacy(m.physSetBounce))
	r.Register("ENTITY.SETBOUNCINESS", "entity", runtime.AdaptLegacy(m.physSetBounce))
	r.Register("PHYSICS.BUILD", "entity", runtime.AdaptLegacy(m.physBuild))

	r.Register("PHYSICS.IMPULSE", "entity", runtime.AdaptLegacy(m.physImpulse))
	r.Register("PHYSICS.VELOCITY", "entity", runtime.AdaptLegacy(m.physVelocity))
	r.Register("PHYSICS.FORCE", "entity", runtime.AdaptLegacy(m.physForce))
	r.Register("PHYSICS.TORQUE", "entity", runtime.AdaptLegacy(m.physTorque))
	r.Register("PHYSICS.SETROT", "entity", runtime.AdaptLegacy(m.physSetRotation))
	r.Register("PHYSICS.GRAVITY", "entity", runtime.AdaptLegacy(m.physSetGravityFactor))
	r.Register("PHYSICS.WAKE", "entity", runtime.AdaptLegacy(m.physWake))
	r.Register("PHYSICS.CCD", "entity", runtime.AdaptLegacy(m.physSetCCD))
}

// entAddPhysics(entity, motion$, shape$ [, mass#]) — motion "static"|"dynamic", shape "box"|"capsule"|"sphere".
func (m *Module) entAddPhysics(args []value.Value) (value.Value, error) {
	if len(args) < 3 {
		return value.Nil, fmt.Errorf("ENTITY.ADDPHYSICS expects (entity, motion$, shape$ [, mass#])")
	}
	if m.h == nil {
		return value.Nil, fmt.Errorf("ENTITY.ADDPHYSICS: heap not bound")
	}
	motion := ""
	shape := "BOX"
	if args[1].Kind == value.KindString {
		if s, ok := m.h.GetString(int32(args[1].IVal)); ok {
			motion = strings.ToUpper(strings.TrimSpace(s))
		}
	}
	if len(args) >= 3 && args[2].Kind == value.KindString {
		if s, ok := m.h.GetString(int32(args[2].IVal)); ok {
			shape = strings.ToUpper(strings.TrimSpace(s))
		}
	}
	mass := 1.0
	if strings.Contains(motion, "STATIC") || motion == "STATIC" {
		mass = 0
	}
	if len(args) >= 4 {
		if v, ok := args[3].ToFloat(); ok {
			mass = v
		}
	}
	shapeKey := "BOX"
	switch {
	case strings.Contains(shape, "CAPSULE"):
		shapeKey = "CAPSULE"
	case strings.Contains(shape, "SPHERE"):
		shapeKey = "SPHERE"
	case strings.Contains(shape, "MESH"):
		shapeKey = "MESH"
	default:
		shapeKey = "BOX"
	}
	idx := m.h.Intern(shapeKey)
	return m.entPhysicsAuto([]value.Value{args[0], value.FromStringIndex(idx), value.FromFloat(mass)})
}

// entPhysicsAuto(id, type$, mass# [, friction#, bounce#])
func (m *Module) entPhysicsAuto(args []value.Value) (value.Value, error) {
	if len(args) < 1 {
		return value.Nil, fmt.Errorf("ENTITY.PHYSICS expects (id [, type, mass, friction, bounce])")
	}
	id, ok := m.entID(args[0])
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("ENTITY.PHYSICS: invalid entity")
	}
	e := m.store().ents[id]
	if e == nil {
		return value.Nil, fmt.Errorf("ENTITY.PHYSICS: unknown entity")
	}

	typ := "BOX"
	if len(args) >= 2 {
		if s, ok := m.h.GetString(int32(args[1].IVal)); ok {
			typ = strings.ToUpper(s)
		}
	}
	mass := 1.0
	if len(args) >= 3 {
		mass, _ = args[2].ToFloat()
	}
	fric := 1.0
	if len(args) >= 4 {
		fric, _ = args[3].ToFloat()
	}
	bounce := 0.0
	if len(args) >= 5 {
		bounce, _ = args[4].ToFloat()
	}
	enableCCD := false
	if len(args) >= 6 {
		enableCCD = value.Truthy(args[5], nil, m.h)
	}

	// Smart Guess Shape logic if typ is empty or "AUTO"
	if typ == "" || typ == "AUTO" {
		typ = "BOX" // Default Fallback (Stable)
		if e.hasRLModel {
			bb := rl.GetModelBoundingBox(e.rlModel)
			dx := bb.Max.X - bb.Min.X
			dy := bb.Max.Y - bb.Min.Y
			dz := bb.Max.Z - bb.Min.Z
			
			// Tall & Thin -> Capsule (1.2x width)
			maxWidth := dx
			if dz > dx { maxWidth = dz }
			
			if dy > maxWidth * 1.2 {
				typ = "CAPSULE"
			} else if dy < (dx + dz) * 0.15 || dx < (dy + dz) * 0.15 || dz < (dx + dy) * 0.15 {
				// Flat & Wide (any axis < 15% of others) -> Plane-style BOX
				typ = "BOX"
			}
		}
	}

	wp := m.worldPos(e)

	motion := jolt.MotionTypeDynamic
	if mass <= 1e-4 {
		motion = jolt.MotionTypeStatic
	}

	allowedDOFs := 0
	if motion == jolt.MotionTypeDynamic && typ == "CAPSULE" {
		// World-space translation + yaw only — stops capsule tipping like a bowling pin.
		allowedDOFs = jolt.AllowedDOFsPlaneXZ
	}

	// 1. Create Builder
	bh, err := m.h.Alloc(&mbphysics3d.BuilderObj{
		Motion:        motion,
		Friction:      float32(fric),
		Restitution:   float32(bounce),
		EnableCCD:     enableCCD,
		AllowedDOFs:   allowedDOFs,
	})
	if err != nil {
		return value.Nil, err
	}
	bHandle := value.FromHandle(bh)

	// 2. Add Shape
	switch typ {
	case "SPHERE":
		rad := e.radius
		if rad <= 1e-3 && e.hasRLModel {
			bb := rl.GetModelBoundingBox(e.rlModel)
			dx := bb.Max.X - bb.Min.X
			dy := bb.Max.Y - bb.Min.Y
			dz := bb.Max.Z - bb.Min.Z
			rad = (dx + dy + dz) / 6.0 // Average half-dimension
		}
		if rad <= 1e-3 { rad = 0.5 }
		_, err = mbphysics3d.BDAddSphere(m.h, []value.Value{bHandle, value.FromFloat(float64(rad))})
	case "CAPSULE":
		rad := e.radius
		height := e.cylH
		if rad <= 1e-3 && e.hasRLModel {
			bb := rl.GetModelBoundingBox(e.rlModel)
			rad = (bb.Max.X - bb.Min.X) / 2.0
			height = bb.Max.Y - bb.Min.Y
		}
		if rad <= 1e-3 { rad = 0.5 }
		if height <= 1e-3 { height = 2.0 }
		_, err = mbphysics3d.BDAddCapsule(m.h, []value.Value{bHandle, value.FromFloat(float64(rad)), value.FromFloat(float64(height))})
	case "MESH":
		// Not implemented in Jolt driver yet
		return value.Nil, fmt.Errorf("ENTITY.PHYSICS: MESH type not yet supported in Jolt driver")
	default: // BOX
		hw, hh, hd := e.w/2, e.h/2, e.d/2
		if hw <= 1e-3 && e.hasRLModel {
			bb := rl.GetModelBoundingBox(e.rlModel)
			hw = (bb.Max.X - bb.Min.X) / 2.0
			hh = (bb.Max.Y - bb.Min.Y) / 2.0
			hd = (bb.Max.Z - bb.Min.Z) / 2.0
		}
		if hw <= 1e-3 {
			hw, hh, hd = 0.5, 0.5, 0.5
		}
		_, err = mbphysics3d.BDAddBox(m.h, []value.Value{bHandle, value.FromFloat(float64(hw)), value.FromFloat(float64(hh)), value.FromFloat(float64(hd))})
	}
	if err != nil {
		return value.Nil, err
	}

	// 3. Commit (position only; motion is on BuilderObj)
	bodyHandleVal, err := mbphysics3d.BDCommit(m.h, []value.Value{
		bHandle,
		value.FromFloat(float64(wp.X)),
		value.FromFloat(float64(wp.Y)),
		value.FromFloat(float64(wp.Z)),
	})
	if err != nil {
		return value.Nil, err
	}

	// 4. Link
	bidxVal, err := mbphysics3d.BDBufferIndex(m.h, []value.Value{bodyHandleVal})
	if err != nil {
		return value.Nil, err
	}
	_, err = m.entLinkPhysBuffer([]value.Value{args[0], bidxVal})
	return value.Nil, err
}

func (m *Module) entPhysicsMotion(args []value.Value) (value.Value, error) {
	// ... mass/motion control logic
	return value.Nil, nil
}

// Map to track builders associated with entities for the property-setter pattern
var pendingBuilders = make(map[int64]heap.Handle)

func (m *Module) physSetShape(args []value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("PHYSICS.SHAPE expects (id, type)")
	}
	id, ok := m.entID(args[0])
	if !ok || id < 1 { return value.Nil, fmt.Errorf("invalid entity") }
	_ = args[1].Kind 
	
	// Create builder if not exists
	bh, ok := pendingBuilders[id]
	if !ok {
		var err error
		bh, err = m.h.Alloc(&mbphysics3d.BuilderObj{})
		if err != nil { return value.Nil, err }
		pendingBuilders[id] = bh
	}
	
	return value.Nil, nil
}

func (m *Module) physSetSize(args []value.Value) (value.Value, error) {
	if len(args) != 4 {
		return value.Nil, fmt.Errorf("PHYSICS.SIZE expects (id, x, y, z)")
	}
	id, ok := m.entID(args[0])
	if !ok || id < 1 { return value.Nil, nil }
	bh, ok := pendingBuilders[id]
	if !ok { return value.Nil, fmt.Errorf("no pending physics builder for entity") }
	
	x, _ := args[1].ToFloat()
	y, _ := args[2].ToFloat()
	z, _ := args[3].ToFloat()
	
	// Delegate to Jolt builder
	return mbphysics3d.BDAddBox(m.h, []value.Value{value.FromHandle(bh), value.FromFloat(x), value.FromFloat(y), value.FromFloat(z)})
}

func (m *Module) physSetFriction(args []value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("PHYSICS.FRICTION expects (id, friction)")
	}
	id, ok := m.entID(args[0])
	if !ok || id < 1 { return value.Nil, nil }
	
	val, _ := args[1].ToFloat()
	
	// If body already exists, set live
	e := m.store().ents[id]
	if e != nil && e.physBufIndex >= 0 {
		mbphysics3d.SetFrictionToIndex(e.physBufIndex, float32(val))
	}

	// Also store in pending builder
	bh, ok := pendingBuilders[id]
	if ok {
		bu, _ := heap.Cast[*mbphysics3d.BuilderObj](m.h, bh)
		if bu != nil {
			bu.Friction = float32(val)
		}
	}
	return value.Nil, nil
}

func (m *Module) physSetBounce(args []value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("PHYSICS.BOUNCE expects (id, restitution)")
	}
	id, ok := m.entID(args[0])
	if !ok || id < 1 { return value.Nil, nil }
	
	val, _ := args[1].ToFloat()
	
	e := m.store().ents[id]
	if e != nil && e.physBufIndex >= 0 {
		mbphysics3d.SetRestitutionToIndex(e.physBufIndex, float32(val))
	}
	// Also store in pending builder
	bh, ok := pendingBuilders[id]
	if ok {
		bu, _ := heap.Cast[*mbphysics3d.BuilderObj](m.h, bh)
		if bu != nil {
			bu.Restitution = float32(val)
		}
	}
	return value.Nil, nil
}

func (m *Module) physBuild(args []value.Value) (value.Value, error) {
	if len(args) < 1 { return value.Nil, fmt.Errorf("PHYSICS.BUILD expects (id [, mass])") }
	id, ok := m.entID(args[0])
	if !ok || id < 1 { return value.Nil, nil }
	bh, ok := pendingBuilders[id]
	if !ok { return value.Nil, fmt.Errorf("no pending physics builder for entity") }
	
	mass := 1.0
	if len(args) >= 2 { mass, _ = args[1].ToFloat() }
	_ = mass 

	e := m.store().ents[id]
	wp := m.worldPos(e)
	
	// Commit and Link
	bodyVal, err := mbphysics3d.BDCommit(m.h, []value.Value{
		value.FromHandle(bh),
		value.FromFloat(float64(wp.X)),
		value.FromFloat(float64(wp.Y)),
		value.FromFloat(float64(wp.Z)),
	})
	if err != nil { return value.Nil, err }
	
	bidxVal, _ := mbphysics3d.BDBufferIndex(m.h, []value.Value{bodyVal})
	m.entLinkPhysBuffer([]value.Value{args[0], bidxVal})
	
	delete(pendingBuilders, id)
	return value.Nil, nil
}

func (m *Module) physImpulse(args []value.Value) (value.Value, error) {
	if len(args) != 4 {
		return value.Nil, fmt.Errorf("PHYSICS.IMPULSE expects (id, x, y, z)")
	}
	id, ok := m.entID(args[0])
	if !ok || id < 1 { return value.Nil, nil }
	e := m.store().ents[id]
	if e == nil || e.physBufIndex < 0 { return value.Nil, nil }

	x, _ := args[1].ToFloat()
	y, _ := args[2].ToFloat()
	z, _ := args[3].ToFloat()
	mbphysics3d.ApplyImpulseToIndex(e.physBufIndex, float32(x), float32(y), float32(z))
	return value.Nil, nil
}

func (m *Module) physVelocity(args []value.Value) (value.Value, error) {
	if len(args) != 4 {
		return value.Nil, fmt.Errorf("PHYSICS.VELOCITY expects (id, x, y, z)")
	}
	id, ok := m.entID(args[0])
	if !ok || id < 1 { return value.Nil, nil }
	e := m.store().ents[id]
	if e == nil || e.physBufIndex < 0 { return value.Nil, nil }

	x, _ := args[1].ToFloat()
	y, _ := args[2].ToFloat()
	z, _ := args[3].ToFloat()
	mbphysics3d.SetVelocityToIndex(e.physBufIndex, float32(x), float32(y), float32(z))
	return value.Nil, nil
}

func (m *Module) physWake(args []value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("PHYSICS.WAKE expects (id)")
	}
	id, ok := m.entID(args[0])
	if !ok || id < 1 { return value.Nil, nil }
	e := m.store().ents[id]
	if e == nil || e.physBufIndex < 0 { return value.Nil, nil }

	mbphysics3d.WakeIndex(e.physBufIndex)
	return value.Nil, nil
}
func (m *Module) physForce(args []value.Value) (value.Value, error) {
	if len(args) != 4 {
		return value.Nil, fmt.Errorf("PHYSICS.FORCE expects (id, x, y, z)")
	}
	id, ok := m.entID(args[0])
	if !ok || id < 1 { return value.Nil, nil }
	e := m.store().ents[id]
	if e == nil || e.physBufIndex < 0 { return value.Nil, nil }

	x, _ := args[1].ToFloat()
	y, _ := args[2].ToFloat()
	z, _ := args[3].ToFloat()
	mbphysics3d.ApplyForceToIndex(e.physBufIndex, float32(x), float32(y), float32(z))
	return value.Nil, nil
}

func (m *Module) physTorque(args []value.Value) (value.Value, error) {
	// Torque not yet implemented in low-level jolt_state
	return value.Nil, nil
}

func (m *Module) physSetRotation(args []value.Value) (value.Value, error) {
	if len(args) != 4 {
		return value.Nil, fmt.Errorf("PHYSICS.SETROT expects (id, p, y, r)")
	}
	id, ok := m.entID(args[0])
	if !ok || id < 1 { return value.Nil, nil }
	e := m.store().ents[id]
	if e == nil || e.physBufIndex < 0 { return value.Nil, nil }

	p, _ := args[1].ToFloat()
	y, _ := args[2].ToFloat()
	r, _ := args[3].ToFloat()
	mbphysics3d.RotateToIndex(e.physBufIndex, float32(p), float32(y), float32(r))
	return value.Nil, nil
}

func (m *Module) physSetGravityFactor(args []value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("PHYSICS.GRAVITY expects (id, factor)")
	}
	id, ok := m.entID(args[0])
	if !ok || id < 1 { return value.Nil, nil }
	e := m.store().ents[id]
	if e == nil || e.physBufIndex < 0 { return value.Nil, nil }

	val, _ := args[1].ToFloat()
	mbphysics3d.SetGravityFactorToIndex(e.physBufIndex, float32(val))
	return value.Nil, nil
}

func (m *Module) physSetCCD(args []value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("PHYSICS.CCD expects (id, bool)")
	}
	id, ok := m.entID(args[0])
	if !ok || id < 1 { return value.Nil, nil }
	
	val := value.Truthy(args[1], nil, m.h)
	
	bh, ok := pendingBuilders[id]
	if ok {
		bu, _ := heap.Cast[*mbphysics3d.BuilderObj](m.h, bh)
		if bu != nil {
			bu.EnableCCD = val
		}
	}
	return value.Nil, nil
}
