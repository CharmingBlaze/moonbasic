//go:build !linux || !cgo

package mbphysics3d

import (
	"moonbasic/runtime"
	"moonbasic/vm/value"
)

const stubHint = "PHYSICS3D/BODY3D require Linux x86_64/arm64 with CGO and Jolt (github.com/bbitechnologies/jolt-go). Windows builds use stubs until bindings exist."

func registerPhysics3DCommands(m *Module, reg runtime.Registrar) {
	reg.Register("PHYSICS3D.GETSCRATCHFLOAT", "physics3d", runtime.AdaptLegacy(func(a []value.Value) (value.Value, error) { return phGetScratchFloat(m, a) }))
	stub := func(name string) runtime.BuiltinFn {
		return func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
			_ = rt
			// No-op on Windows/non-CGO to prevent crashes
			return value.Nil, nil
		}
	}
	keys := []string{
		"PHYSICS3D.START", "PHYSICS3D.STOP", "PHYSICS3D.SETGRAVITY", "PHYSICS3D.STEP",
		"PHYSICS3D.SETSUBSTEPS", "PHYSICS3D.DEBUGDRAW", "PHYSICS3D.SETTIMESTEP",
		"PHYSICS3D.ONCOLLISION", "PHYSICS3D.PROCESSCOLLISIONS", "PHYSICS3D.RAYCAST", "PHYSICS3D.SYNCWASMTOPHYSREGS",
		"PHYSICS.START", "PHYSICS.STOP", "PHYSICS.SETGRAVITY", "PHYSICS.STEP", "PHYSICS.SETSUBSTEPS",
		"PHYSICS.RAYCAST", "PHYSICS.SPHERECAST", "PHYSICS.BOXCAST", "PHYSICS.ENABLE", "PHYSICS.DISABLE",
		"BODY3D.MAKE", "BODY3D.ADDBOX", "BODY3D.ADDSPHERE", "BODY3D.ADDCAPSULE", "BODY3D.ADDMESH",
		"BODY3D.COMMIT", "BODY3D.SETPOS", "BODY3D.SETPOSITION", "BODY3D.GETPOS", "BODY3D.ACTIVATE", "BODY3D.DEACTIVATE",
		"BODY3D.SETROT", "BODY3D.GETROT",
		"BODY3D.SETMASS", "BODY3D.SETFRICTION", "BODY3D.SETRESTITUTION",
		"BODY3D.APPLYFORCE", "BODY3D.APPLYIMPULSE", "BODY3D.SETLINEARVEL", "BODY3D.SETANGULARVEL",
		"BODY3D.X", "BODY3D.Y", "BODY3D.Z", "BODY3D.FREE",
		"BODY3D.COLLIDED", "BODY3D.COLLISIONOTHER", "BODY3D.COLLISIONPOINT", "BODY3D.COLLISIONNORMAL",
		"JOINT3D.FIXED", "JOINT3D.HINGE", "JOINT3D.SLIDER", "JOINT3D.CONE", "JOINT3D.DELETE",
		"PICK.ORIGIN", "PICK.DIRECTION", "PICK.MAXDIST", "PICK.LAYERMASK", "PICK.RADIUS",
		"PICK.CAST", "PICK.FROMCAMERA", "PICK.SCREENCAST",
		"PICK.X", "PICK.Y", "PICK.Z", "PICK.NX", "PICK.NY", "PICK.NZ", "PICK.ENTITY", "PICK.DIST", "PICK.HIT",
	}
	for _, k := range keys {
		reg.Register(k, "physics3d", stub(k))
	}
}

func shutdownPhysics3D(m *Module) { _ = m }
