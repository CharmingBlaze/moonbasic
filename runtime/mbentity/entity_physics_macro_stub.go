//go:build !linux || !cgo

package mbentity

import (
	"moonbasic/runtime"
	mbphysics3d "moonbasic/runtime/physics3d"
	"moonbasic/vm/value"
)

// Jolt body builder helpers (BDAdd*, BuilderObj, etc.) exist only on linux+cgo.
// Full API: entity_physics_macro_cgo.go (linux && cgo).

func registerEntityPhysicsMacroAPI(m *Module, r runtime.Registrar) {
	noop := runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		_ = args
		return value.Nil, mbphysics3d.ErrPhysicsUnavailable()
	})
	keys := []string{
		"ENTITY.PHYSICS", "PHYSICS.AUTO", "ENTITY.PHYSICSMOTION",
		"PHYSICS.SHAPE", "PHYSICS.SIZE", "PHYSICS.FRICTION", "PHYSICS.BOUNCE", "PHYSICS.BUILD",
		"PHYSICS.IMPULSE", "PHYSICS.VELOCITY", "PHYSICS.FORCE", "PHYSICS.TORQUE", "PHYSICS.SETROT",
		"PHYSICS.GRAVITY", "PHYSICS.WAKE", "PHYSICS.CCD",
	}
	for _, k := range keys {
		r.Register(k, "entity", noop)
	}
	_ = m
}
