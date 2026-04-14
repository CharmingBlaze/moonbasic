//go:build (!linux && !windows) || !cgo

package mbentity

import (
	"moonbasic/runtime"
	"moonbasic/vm/value"
)

// Jolt body builder helpers (BDAdd*, BuilderObj, etc.) exist only on (linux||windows)+cgo.
// Full API: entity_physics_macro_cgo.go ((linux || windows) && cgo).

func registerEntityPhysicsMacroAPI(m *Module, r runtime.Registrar) {
	// Succeed without Jolt so the same .mb runs on Windows; physics has no effect until Linux+CGO Jolt.
	noop := runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		_ = args
		return value.Nil, nil
	})
	keys := []string{
		"ENTITY.PHYSICS", "ENTITY.ADDPHYSICS", "PHYSICS.AUTO", "ENTITY.PHYSICSMOTION",
		"PHYSICS.SHAPE", "PHYSICS.SIZE", "PHYSICS.FRICTION", "PHYSICS.BOUNCE", "ENTITY.SETBOUNCINESS", "PHYSICS.BUILD",
		"PHYSICS.IMPULSE", "PHYSICS.VELOCITY", "PHYSICS.FORCE", "PHYSICS.TORQUE", "PHYSICS.SETROT",
		"PHYSICS.GRAVITY", "PHYSICS.WAKE", "PHYSICS.CCD",
	}
	for _, k := range keys {
		r.Register(k, "entity", noop)
	}
	_ = m
}
