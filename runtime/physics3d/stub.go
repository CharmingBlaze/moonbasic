//go:build (!linux && !windows) || !cgo

package mbphysics3d

import (
	"moonbasic/vm/heap"
	"moonbasic/vm/value"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Vec3 struct {
	X, Y, Z float32
}

func GravityVec() Vec3 { return Vec3{0, -9.81, 0} }

type JointObj struct {
	Release heap.ReleaseOnce
}

func (j *JointObj) TypeName() string { return "Joint3D" }
func (j *JointObj) TypeTag() uint16  { return heap.TagPhysicsBody + 10 } // Use a distinct tag for joints
func (j *JointObj) Free()            {}

// stubHint explains why native Jolt is unavailable on this build (see jolt-go cgo_* and lib/{platform}).
const stubHint = "native Jolt is not linked on this build: need (linux or windows) amd64/arm64 with CGO_ENABLED=1 and Jolt static libraries under third_party/jolt-go/jolt/lib/<platform> (see docs/JOLT_WINDOWS_PARITY.md). Builds without CGO use this stub."

type ShapeObj struct {
	Kind int
	F1, F2, F3 float32
}
type body3dObj struct {
	ID          int
	Pos, Rot    Vec3
	LinearVel  Vec3
	AngularVel Vec3
	Shape       *ShapeObj
	Layer       int
	Collision   bool
}

var (
	staticBodies = make(map[heap.Handle]*body3dObj)
	nextBodyID   = 1
)

func phSetOnCollision(m *Module, ha, hb value.Value, cb string) (value.Value, error) {
	_ = m
	_, _, _ = ha, hb, cb
	return value.Nil, nil
}

func phCreateBody(m *Module, motion string) (value.Value, error) {
	_ = motion
	if m == nil || m.h == nil {
		return value.Nil, nil
	}
	id, _ := m.h.Alloc(&BuilderObj{})
	return value.FromHandle(id), nil
}

func (m *Module) phStart(args []value.Value) (value.Value, error)              { return value.Nil, nil }
func (m *Module) phStop(args []value.Value) (value.Value, error)               { return value.Nil, nil }
func (m *Module) phSetGravity(args []value.Value) (value.Value, error)        { return value.Nil, nil }
func (m *Module) phGetGravityX(args []value.Value) (value.Value, error)       { return value.FromFloat(0), nil }
func (m *Module) phGetGravityY(args []value.Value) (value.Value, error)       { return value.FromFloat(-9.81), nil }
func (m *Module) phGetGravityZ(args []value.Value) (value.Value, error)       { return value.FromFloat(0), nil }
func (m *Module) phStep(args []value.Value) (value.Value, error) {
	dt := 1.0 / 60.0
	if len(args) == 1 {
		if v, ok := args[0].ToFloat(); ok {
			dt = v
		}
	}
	
	// Basic Euler Integration for Dynamic Bodies in Stub
	for _, b := range staticBodies {
		if b.Collision { // Use this as a proxy for "dynamic" in stub for now
			// Integrate Velocity -> Position
			b.Pos.X += b.LinearVel.X * float32(dt)
			b.Pos.Y += b.LinearVel.Y * float32(dt)
			b.Pos.Z += b.LinearVel.Z * float32(dt)
			
			// Simple Gravity
			b.LinearVel.Y -= 9.81 * float32(dt)
		}
	}

	// 1. Process Aero (Shared Go Logic)
	m.ProcessAeroDynamics(float32(dt))
	
	return value.Nil, nil
}
func (m *Module) phSetTimeStep(args []value.Value) (value.Value, error)       { return value.Nil, nil }
func (m *Module) phGetMatrixBuffer(args []value.Value) (value.Value, error)   { return value.Nil, nil }
func (m *Module) phSetSubsteps(args []value.Value) (value.Value, error)       { return value.Nil, nil }
func (m *Module) phProcessCollisions(args []value.Value) (value.Value, error) { return value.Nil, nil }
func (m *Module) phRaycast(args []value.Value) (value.Value, error)           { return value.Nil, nil }
func (m *Module) phSpherecast(args []value.Value) (value.Value, error)        { return value.Nil, nil }
func (m *Module) phBoxcast(args []value.Value) (value.Value, error)           { return value.Nil, nil }
func (m *Module) phEnable(args []value.Value) (value.Value, error)            { return value.Nil, nil }
func (m *Module) phDisable(args []value.Value) (value.Value, error)           { return value.Nil, nil }
func (m *Module) phSyncWasmToPhysRegs(args []value.Value) (value.Value, error) { return value.Nil, nil }
func (m *Module) phDebugDraw(args []value.Value) (value.Value, error)         { return value.Nil, nil }

func (m *Module) shCreateBox(args []value.Value) (value.Value, error) {
	hx, _ := args[0].ToFloat(); hy, _ := args[1].ToFloat(); hz, _ := args[2].ToFloat()
	id, _ := m.h.Alloc(&ShapeObj{1, float32(hx), float32(hy), float32(hz)})
	return value.FromHandle(id), nil
}
func (m *Module) shCreateSphere(args []value.Value) (value.Value, error) {
	r, _ := args[0].ToFloat()
	id, _ := m.h.Alloc(&ShapeObj{2, float32(r), 0, 0})
	return value.FromHandle(id), nil
}
func (m *Module) shCreateCapsule(args []value.Value) (value.Value, error) {
	r, _ := args[0].ToFloat(); hi, _ := args[1].ToFloat()
	id, _ := m.h.Alloc(&ShapeObj{3, float32(r), float32(hi), 0})
	return value.FromHandle(id), nil
}
func (m *Module) shCreateCylinder(args []value.Value) (value.Value, error) {
	r, _ := args[0].ToFloat(); hi, _ := args[1].ToFloat()
	id, _ := m.h.Alloc(&ShapeObj{4, float32(r), float32(hi), 0})
	return value.FromHandle(id), nil
}
func (m *Module) shGetType(args []value.Value) (value.Value, error) {
	sh, err := heap.Cast[*ShapeObj](m.h, heap.Handle(args[0].IVal))
	if err != nil { return value.Nil, err }
	return value.FromInt(int64(sh.Kind)), nil
}
func (m *Module) shGetDim1(args []value.Value) (value.Value, error) {
	sh, err := heap.Cast[*ShapeObj](m.h, heap.Handle(args[0].IVal))
	if err != nil { return value.Nil, err }
	return value.FromFloat(float64(sh.F1)), nil
}
func (m *Module) shGetDim2(args []value.Value) (value.Value, error) {
	sh, err := heap.Cast[*ShapeObj](m.h, heap.Handle(args[0].IVal))
	if err != nil { return value.Nil, err }
	return value.FromFloat(float64(sh.F2)), nil
}
func (m *Module) shGetDim3(args []value.Value) (value.Value, error) {
	sh, err := heap.Cast[*ShapeObj](m.h, heap.Handle(args[0].IVal))
	if err != nil { return value.Nil, err }
	return value.FromFloat(float64(sh.F3)), nil
}
func (m *Module) shFree(args []value.Value) (value.Value, error) {
	m.h.Free(heap.Handle(args[0].IVal))
	return value.Nil, nil
}

func (m *Module) knCreate(args []value.Value) (value.Value, error) {
	sh, _ := heap.Cast[*ShapeObj](m.h, heap.Handle(args[0].IVal))
	body := &body3dObj{ID: nextBodyID, Shape: sh, Collision: true}
	nextBodyID++
	id, _ := m.h.Alloc(body)
	return value.FromHandle(id), nil
}
func (m *Module) stCreate(args []value.Value) (value.Value, error) {
	sh, _ := heap.Cast[*ShapeObj](m.h, heap.Handle(args[0].IVal))
	body := &body3dObj{ID: nextBodyID, Shape: sh, Collision: true}
	nextBodyID++
	id, _ := m.h.Alloc(body)
	staticBodies[id] = body
	return value.FromHandle(id), nil
}
func (m *Module) trCreate(args []value.Value) (value.Value, error) {
	sh, _ := heap.Cast[*ShapeObj](m.h, heap.Handle(args[0].IVal))
	body := &body3dObj{ID: nextBodyID, Shape: sh, Collision: true}
	nextBodyID++
	id, _ := m.h.Alloc(body)
	return value.FromHandle(id), nil
}

func (m *Module) bdSetPos(args []value.Value) (value.Value, error) {
	bo, _ := m.hGetBody(args[0])
	if bo != nil && len(args) >= 4 {
		x, _ := args[1].ToFloat(); y, _ := args[2].ToFloat(); z, _ := args[3].ToFloat()
		bo.Pos = Vec3{X: float32(x), Y: float32(y), Z: float32(z)}
	}
	return value.Nil, nil
}
func (m *Module) bdGetPos(args []value.Value) (value.Value, error) {
	bo, _ := m.hGetBody(args[0])
	if bo == nil { return valueVec3FromFloats(m.h, 0, 0, 0) }
	return valueVec3FromFloats(m.h, float64(bo.Pos.X), float64(bo.Pos.Y), float64(bo.Pos.Z))
}
func (m *Module) bdActivate(args []value.Value) (value.Value, error)     { return value.Nil, nil }
func (m *Module) bdDeactivate(args []value.Value) (value.Value, error)   { return value.Nil, nil }
func (m *Module) bdSetRotation(args []value.Value) (value.Value, error) { return value.Nil, nil }
func (m *Module) bdGetRotZero(args []value.Value) (value.Value, error)  { return value.FromFloat(0), nil }
func (m *Module) bdSetMass(args []value.Value) (value.Value, error)      { return value.Nil, nil }
func (m *Module) bdSetFriction(args []value.Value) (value.Value, error)  { return value.Nil, nil }
func (m *Module) bdSetRestitution(args []value.Value) (value.Value, error) { return value.Nil, nil }
func (m *Module) bdApplyForce(args []value.Value) (value.Value, error)   { return value.Nil, nil }
func (m *Module) bdApplyImpulse(args []value.Value) (value.Value, error) { return value.Nil, nil }
func (m *Module) bdAxis(args []value.Value, axis int) (value.Value, error) {
	bo, _ := m.hGetBody(args[0])
	if bo == nil { return value.FromFloat(0), nil }
	switch axis {
	case 0: return value.FromFloat(float64(bo.Pos.X)), nil
	case 1: return value.FromFloat(float64(bo.Pos.Y)), nil
	case 2: return value.FromFloat(float64(bo.Pos.Z)), nil
	}
	return value.FromFloat(0), nil
}
func (m *Module) bdFree(args []value.Value) (value.Value, error) {
	bh := heap.Handle(args[0].IVal)
	delete(staticBodies, bh)
	m.h.Free(bh)
	return value.Nil, nil
}
func (m *Module) bdCollided3D(args []value.Value) (value.Value, error)      { return value.FromBool(false), nil }
func (m *Module) bdCollisionOther3D(args []value.Value) (value.Value, error) { return value.FromInt(0), nil }
func (m *Module) bdCollisionPoint3D(args []value.Value) (value.Value, error) { return value.FromFloat(0), nil }
func (m *Module) bdCollisionNormal3D(args []value.Value) (value.Value, error) { return value.FromFloat(0), nil }
func (m *Module) bdNoOp(args []value.Value) (value.Value, error)           { return value.Nil, nil }

func (m *Module) brSetPos(args []value.Value) (value.Value, error) {
	bo, _ := heap.Cast[*body3dObj](m.h, heap.Handle(args[0].IVal))
	if bo != nil && len(args) >= 4 {
		x, _ := args[1].ToFloat(); y, _ := args[2].ToFloat(); z, _ := args[3].ToFloat()
		bo.Pos = Vec3{X: float32(x), Y: float32(y), Z: float32(z)}
	}
	return value.Nil, nil
}
func (m *Module) brSetLayer(args []value.Value) (value.Value, error)   { return value.Nil, nil }
func (m *Module) brEnableColl(args []value.Value) (value.Value, error) { return value.Nil, nil }
func (m *Module) brFree(args []value.Value) (value.Value, error) {
	bh := heap.Handle(args[0].IVal)
	delete(staticBodies, bh)
	m.h.Free(bh)
	return value.Nil, nil
}

func (m *Module) phJointFixed(args []value.Value) (value.Value, error)  { return value.Nil, nil }
func (m *Module) phJointHinge(args []value.Value) (value.Value, error)  { return value.Nil, nil }
func (m *Module) phJointSlider(args []value.Value) (value.Value, error) { return value.Nil, nil }
func (m *Module) phJointCone(args []value.Value) (value.Value, error)   { return value.Nil, nil }
func (m *Module) phJointDelete(args []value.Value) (value.Value, error) { return value.Nil, nil }

func phLevelStatic(m *Module, args []value.Value) (value.Value, error) {
	return value.Nil, nil
}

func (m *Module) hGetBody(v value.Value) (*body3dObj, error) {
	return heap.Cast[*body3dObj](m.h, heap.Handle(v.IVal))
}

func (s *ShapeObj) TypeName() string { return "Shape" }
func (s *ShapeObj) TypeTag() uint16  { return heap.TagShape }
func (s *ShapeObj) Free()            {}

func (b *body3dObj) TypeName() string { return "Body3D" }
func (b *body3dObj) TypeTag() uint16 {
	return heap.TagPhysicsBody // Fallback sharing
}
func (b *body3dObj) Free() {}

func shutdownPhysics3D(m *Module) { _ = m }

// Exported for charcontroller/stub.go
func GetStaticBodyRegistry() map[heap.Handle]*body3dObj { return staticBodies }

type body3dObjExport struct {
	Pos   Vec3
	Shape *ShapeObj
}

func ApplyImpulseToIndex(idx int, x, y, z float32)     {}
func GetLinearVelocityToIndex(idx int) (x, y, z float32) { return 0, 0, 0 }
func GetBodyQuaternionForBufferIndex(idx int) (x, y, z, w float32, ok bool) { return 0, 0, 0, 1, false }
func SetVelocityToIndex(idx int, x, y, z float32)      {}
func SetPositionToIndex(idx int, x, y, z float32)      {}
func SetFrictionToIndex(idx int, x float32)            {}
func SetRestitutionToIndex(idx int, x float32)         {}
func WakeIndex(idx int)                                {}
func ApplyForceToIndex(idx int, x, y, z float32)       {}
func RotateToIndex(idx int, p, y, r float32)           {}
func SetGravityFactorToIndex(idx int, x float32)      {}

func (m *Module) phWorldSetup(args []value.Value) (value.Value, error) {
	return value.Nil, nil
}

func (m *Module) bdAddMesh(args []value.Value) (value.Value, error) { return value.Nil, nil }
func (m *Module) bdAddBox(args []value.Value) (value.Value, error) { return value.Nil, nil }
func (m *Module) bdAddSphere(args []value.Value) (value.Value, error) { return value.Nil, nil }
func (m *Module) bdAddCapsule(args []value.Value) (value.Value, error) { return value.Nil, nil }
func (m *Module) bdCommit(args []value.Value) (value.Value, error)  { return value.Nil, nil }
func (m *Module) bdBufferIndex(args []value.Value) (value.Value, error) {
	return value.FromInt(-1), nil
}

// Joint Stubs
func (m *Module) phCreateHingeJoint(args []value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, nil
	}
	id, _ := m.h.Alloc(&JointObj{})
	return value.FromHandle(id), nil
}
func (m *Module) phCreatePointJoint(args []value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, nil
	}
	id, _ := m.h.Alloc(&JointObj{})
	return value.FromHandle(id), nil
}

func (m *Module) bdGetAxis(args []value.Value, axis int) (value.Value, error) { return value.FromFloat(0), nil }

// Advanced Body Stubs (Restored)
func (m *Module) bdSetDamping(args []value.Value) (value.Value, error)       { return value.Nil, nil }
func (m *Module) bdLockAxis(args []value.Value) (value.Value, error)        { return value.Nil, nil }
func (m *Module) bdSetGravityFactor(args []value.Value) (value.Value, error) { return value.Nil, nil }
func (m *Module) btdSetCCD(args []value.Value) (value.Value, error)          { return value.Nil, nil }

func (m *Module) bdGetLinearVel(args []value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, nil
	}
	b, _ := heap.Cast[*body3dObj](m.h, heap.Handle(args[0].IVal))
	if b == nil {
		return valueVec3FromFloats(m.h, 0, 0, 0)
	}
	return valueVec3FromFloats(m.h, float64(b.LinearVel.X), float64(b.LinearVel.Y), float64(b.LinearVel.Z))
}

func (m *Module) bdSetLinearVel(args []value.Value) (value.Value, error) {
	if len(args) == 4 && args[0].Kind == value.KindHandle {
		b, _ := heap.Cast[*body3dObj](m.h, heap.Handle(args[0].IVal))
		if b != nil {
			vx, _ := args[1].ToFloat()
			vy, _ := args[2].ToFloat()
			vz, _ := args[3].ToFloat()
			b.LinearVel = Vec3{X: float32(vx), Y: float32(vy), Z: float32(vz)}
		}
	}
	return value.Nil, nil
}

func (m *Module) bdGetAngularVel(args []value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, nil
	}
	b, _ := heap.Cast[*body3dObj](m.h, heap.Handle(args[0].IVal))
	if b == nil {
		return valueVec3FromFloats(m.h, 0, 0, 0)
	}
	return valueVec3FromFloats(m.h, float64(b.AngularVel.X), float64(b.AngularVel.Y), float64(b.AngularVel.Z))
}

func (m *Module) bdSetAngularVel(args []value.Value) (value.Value, error) {
	if len(args) == 4 && args[0].Kind == value.KindHandle {
		b, _ := heap.Cast[*body3dObj](m.h, heap.Handle(args[0].IVal))
		if b != nil {
			ax, _ := args[1].ToFloat()
			ay, _ := args[2].ToFloat()
			az, _ := args[3].ToFloat()
			b.AngularVel = Vec3{X: float32(ax), Y: float32(ay), Z: float32(az)}
		}
	}
	return value.Nil, nil
}

func (m *Module) bdApplyTorque(args []value.Value) (value.Value, error) { return value.Nil, nil }
func (m *Module) bdGetMass(args []value.Value) (value.Value, error)      { return value.FromFloat(1.0), nil }

// Aero Stubs (Logic implemented in Go, shared between platforms; handles in module.go)
func (m *Module) arSetLift(args []value.Value) (value.Value, error)   { return value.Nil, nil }
func (m *Module) arSetThrust(args []value.Value) (value.Value, error) { return value.Nil, nil }
func (m *Module) arSetDrag(args []value.Value) (value.Value, error)   { return value.Nil, nil }

type BuilderObj struct {
	Motion   int
	Friction float32
	Shape    *ShapeObj
}

func (b *BuilderObj) TypeName() string { return "Body3DBuilder" }
func (b *BuilderObj) TypeTag() uint16  { return heap.TagPhysicsBuilder }
func (b *BuilderObj) Free()            {}

// Internal Bridges for shared Go solvers (aero_host.go, vehicle_host.go)

func (m *Module) getBodyTransform(b *body3dObj) (rl.Vector3, rl.Quaternion, bool) {
	// Simple stub behavior: return stored position/rotation
	return rl.Vector3{X: b.Pos.X, Y: b.Pos.Y, Z: b.Pos.Z}, rl.QuaternionIdentity(), true
}

func (m *Module) getBodyVelocity(b *body3dObj) rl.Vector3 {
	return rl.Vector3{X: b.LinearVel.X, Y: b.LinearVel.Y, Z: b.LinearVel.Z}
}

func (m *Module) applyBodyForce(b *body3dObj, f rl.Vector3) {
	// Simple stub behavior: integrate force to linear velocity (mass=1.0)
	// dt assumed for simulation here (this is a simplified stub bridge)
	b.LinearVel.X += f.X * 0.016 // Approx 1/60 step
	b.LinearVel.Y += f.Y * 0.016
	b.LinearVel.Z += f.Z * 0.016
}
