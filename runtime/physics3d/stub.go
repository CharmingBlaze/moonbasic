//go:build !linux || !cgo

package mbphysics3d

import (
	"fmt"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

const stubHint = "PHYSICS3D/BODY3D require Linux x86_64/arm64 with CGO and Jolt (github.com/bbitechnologies/jolt-go). Windows builds use stubs until bindings exist."

func registerPhysics3DCommands(m *Module, reg runtime.Registrar) {
	_ = m
	stub := func(name string) runtime.BuiltinFn {
		return func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
			_ = rt
			return value.Nil, fmt.Errorf("%s: %s", name, stubHint)
		}
	}
	keys := []string{
		"PHYSICS3D.START", "PHYSICS3D.STOP", "PHYSICS3D.SETGRAVITY", "PHYSICS3D.STEP",
		"PHYSICS3D.SETSUBSTEPS",
		"PHYSICS3D.ONCOLLISION", "PHYSICS3D.PROCESSCOLLISIONS", "PHYSICS3D.RAYCAST",
		"BODY3D.MAKE", "BODY3D.ADDBOX", "BODY3D.ADDSPHERE", "BODY3D.ADDCAPSULE", "BODY3D.ADDMESH",
		"BODY3D.COMMIT", "BODY3D.SETPOS", "BODY3D.GETPOS", "BODY3D.SETROT", "BODY3D.GETROT",
		"BODY3D.SETMASS", "BODY3D.SETFRICTION", "BODY3D.SETRESTITUTION",
		"BODY3D.APPLYFORCE", "BODY3D.APPLYIMPULSE", "BODY3D.SETLINEARVEL", "BODY3D.SETANGULARVEL",
		"BODY3D.X", "BODY3D.Y", "BODY3D.Z", "BODY3D.FREE",
	}
	for _, k := range keys {
		reg.Register(k, "physics3d", stub(k))
	}
}

func shutdownPhysics3D(m *Module) { _ = m }
