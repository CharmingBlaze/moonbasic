//go:build !linux || !cgo

package mbphysics3d

import (
	"fmt"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

// stubHint explains why native Jolt is unavailable on this build (see jolt-go cgo_* files: Linux/Darwin only today).
const stubHint = "native Jolt is not linked on this build (need Linux/macOS amd64/arm64 + CGO; github.com/bbitechnologies/jolt-go has no Windows CGO libs yet). Use Linux CI, WSL2 for dev parity, or contribute Windows static libs to jolt-go."

func registerPhysics3DCommands(m *Module, reg runtime.Registrar) {
	reg.Register("PHYSICS3D.GETSCRATCHFLOAT", "physics3d", runtime.AdaptLegacy(func(a []value.Value) (value.Value, error) { return phGetScratchFloat(m, a) }))
	stub := func(name string) runtime.BuiltinFn {
		return func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
			_ = rt
			_ = args
			return value.Nil, fmt.Errorf("%s: %s", name, stubHint)
		}
	}
	keys := []string{
		"PHYSICS3D.START", "PHYSICS3D.STOP", "PHYSICS3D.SETGRAVITY", "PHYSICS3D.STEP",
		"PHYSICS3D.SETSUBSTEPS", "PHYSICS3D.DEBUGDRAW", "PHYSICS3D.SETTIMESTEP", "PHYSICS3D.GETMATRIXBUFFER",
		"PHYSICS3D.ONCOLLISION", "PHYSICS3D.PROCESSCOLLISIONS", "PHYSICS3D.RAYCAST", "PHYSICS3D.SYNCWASMTOPHYSREGS",
		"PHYSICS.START", "PHYSICS.STOP", "PHYSICS.SETGRAVITY", "PHYSICS.STEP", "PHYSICS.SETSUBSTEPS",
		"PHYSICS.RAYCAST", "PHYSICS.SPHERECAST", "PHYSICS.BOXCAST", "PHYSICS.ENABLE", "PHYSICS.DISABLE",
		"BODY3D.MAKE", "BODY3D.ADDBOX", "BODY3D.ADDSPHERE", "BODY3D.ADDCAPSULE", "BODY3D.ADDMESH",
		"BODY3D.COMMIT", "BODY3D.SETPOS", "BODY3D.SETPOSITION", "BODY3D.GETPOS", "BODY3D.ACTIVATE", "BODY3D.DEACTIVATE",
		"BODY3D.SETROT", "BODY3D.GETROT",
		"BODY3D.SETMASS", "BODY3D.SETFRICTION", "BODY3D.SETRESTITUTION",
		"BODY3D.APPLYFORCE", "BODY3D.APPLYIMPULSE", "BODY3D.SETLINEARVEL", "BODY3D.SETANGULARVEL",
		"BODY3D.X", "BODY3D.Y", "BODY3D.Z", "BODY3D.FREE", "BODY3D.BUFFERINDEX",
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

// ErrPhysicsUnavailable is returned by stub physics builtins when native Jolt is not linked.
func ErrPhysicsUnavailable() error {
	return fmt.Errorf("%s", stubHint)
}

// Exported stubs for mbentity
type BuilderObj struct {
	Friction    float32
	Restitution float32
}

func (b *BuilderObj) TypeName() string { return "Body3DBuilder" }
func (b *BuilderObj) TypeTag() uint16  { return 0 }
func (b *BuilderObj) Free()            {}

func BDAddBox(h any, args []value.Value) (value.Value, error) {
	_, _ = h, args
	return value.Nil, fmt.Errorf("BDAddBox: %s", stubHint)
}
func BDAddSphere(h any, args []value.Value) (value.Value, error) {
	_, _ = h, args
	return value.Nil, fmt.Errorf("BDAddSphere: %s", stubHint)
}
func BDAddCapsule(h any, args []value.Value) (value.Value, error) {
	_, _ = h, args
	return value.Nil, fmt.Errorf("BDAddCapsule: %s", stubHint)
}
func BDCommit(h any, args []value.Value) (value.Value, error) {
	_, _ = h, args
	return value.Nil, fmt.Errorf("BDCommit: %s", stubHint)
}
func BDBufferIndex(h any, args []value.Value) (value.Value, error) {
	_, _ = h, args
	return value.Nil, fmt.Errorf("BDBufferIndex: %s", stubHint)
}

func ApplyImpulseToIndex(idx int, x, y, z float32)     {}
func SetVelocityToIndex(idx int, x, y, z float32)      {}
func SetFrictionToIndex(idx int, x float32)            {}
func SetRestitutionToIndex(idx int, x float32)         {}
func WakeIndex(idx int)                                {}
func ApplyForceToIndex(idx int, x, y, z float32)       {}
func RotateToIndex(idx int, p, y, r float32)           {}
func SetGravityFactorToIndex(idx int, x float32)      {}
