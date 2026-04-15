//go:build (linux || windows) && cgo

/*
Package mbphysics3d provides the Jolt Physics integration for MoonBASIC.

ARCHITECTURE: SHARED-MEMORY SYNCHRONIZATION
This module implements a zero-copy shared memory architecture for physics-to-entity synchronization.
Physics simulation results (position and rotation) are written directly into a contiguous float32
buffer (the "matrix buffer") which is shared with the entity module (mbentity).

CROSS-PLATFORM CGO
Since v1.3.1, this module is platform-agnostic for CGO builds. On Windows, it links against
static libraries (libJolt.a, libjolt_wrapper.a) built from the local third_party/jolt-go tree.
*/
package mbphysics3d

import (
	goruntime "runtime"
	"sync"

	"github.com/bbitechnologies/jolt-go/jolt"

	"moonbasic/vm/heap"
)

var (
	joltMu       sync.Mutex
	joltSys      *jolt.PhysicsSystem
	joltBi       *jolt.BodyInterface
	gravX        float32 = 0
	gravY        float32 = -9.81
	gravZ        float32 = 0
	collRules    []collRule
	collPending  []collEvent
	joltCoreInit bool
)

var (
	matrixBuffer      []float32
	prevMatrixBuffer  []float32 // snapshot at start of STEP (for optional render interpolation + collision reuse)
	nextBufferIndex   int
	bufferIndexMap    map[*jolt.BodyID]int // body -> index
	bufferIndexToBody map[int]*jolt.BodyID // index -> body
	matrixBufferAlloc int
	// Written at end of PHYSICS3D.STEP; used by PhysicsMatrixInterpAlpha.
	matrixInterpAccum float64
	matrixInterpFixed float64
)

type collRule struct {
	ha, hb heap.Handle
	cb     string
}

type collEvent struct {
	ha, hb heap.Handle
	cb     string
}

var (
	joltBodyMu       sync.Mutex
	joltBodyToHandle map[*jolt.BodyID]heap.Handle
	joltBodyPacked   map[uint32]heap.Handle // Jolt BodyID index+sequence → BODY3D heap handle (KCC contact fan-in)
	joltBodyDynamic  map[*jolt.BodyID]bool
)

func joltRegisterBody(id *jolt.BodyID, h heap.Handle) {
	if id == nil {
		return
	}
	joltBodyMu.Lock()
	defer joltBodyMu.Unlock()
	if joltBodyToHandle == nil {
		joltBodyToHandle = make(map[*jolt.BodyID]heap.Handle)
	}
	joltBodyToHandle[id] = h
	packed := id.IndexAndSequenceNumber()
	if joltBodyPacked == nil {
		joltBodyPacked = make(map[uint32]heap.Handle)
	}
	joltBodyPacked[packed] = h
}

func joltMarkBodyDynamic(id *jolt.BodyID, dynamic bool) {
	if id == nil {
		return
	}
	joltBodyMu.Lock()
	defer joltBodyMu.Unlock()
	if joltBodyDynamic == nil {
		joltBodyDynamic = make(map[*jolt.BodyID]bool)
	}
	joltBodyDynamic[id] = dynamic
}

func joltUnregisterBody(id *jolt.BodyID) {
	if id == nil {
		return
	}
	joltBodyMu.Lock()
	defer joltBodyMu.Unlock()
	packed := id.IndexAndSequenceNumber()
	delete(joltBodyToHandle, id)
	delete(joltBodyDynamic, id)
	if joltBodyPacked != nil {
		delete(joltBodyPacked, packed)
	}
}

func joltLookupHandle(id *jolt.BodyID) (heap.Handle, bool) {
	if id == nil {
		return 0, false
	}
	joltBodyMu.Lock()
	defer joltBodyMu.Unlock()
	h, ok := joltBodyToHandle[id]
	return h, ok
}

// LookupBodyHeapByPacked resolves CharacterContactEvent.BodyB (index+sequence) to a registered BODY3D handle.
func LookupBodyHeapByPacked(packed uint32) (heap.Handle, bool) {
	joltBodyMu.Lock()
	defer joltBodyMu.Unlock()
	if joltBodyPacked == nil {
		return 0, false
	}
	h, ok := joltBodyPacked[packed]
	return h, ok
}

// JointObj is a script-visible joint handle. The vendored jolt-go wrapper does not expose hinge/point
// constraints yet; JOINT.CREATE* allocates a placeholder (parity with stub builds).
type JointObj struct {
	Release heap.ReleaseOnce
}

func (j *JointObj) TypeName() string { return "Joint3D" }
func (j *JointObj) TypeTag() uint16  { return heap.TagPhysicsBody + 10 }
func (j *JointObj) Free() {
	j.Release.Do(func() {})
}

type body3dObj struct {
	id          *jolt.BodyID
	queryShape  *jolt.Shape // duplicate shape for CollideShape queries (same geometry as body)
	bufferIndex int         // index into sharedMatrixBuffer
	motion      jolt.MotionType
	release     heap.ReleaseOnce
	// Primitive template from COMMIT / SHAPE.CREATE* (qKind 0 = mesh or unknown)
	qKind      uint8
	qBox       jolt.Vec3
	qSphere    float32
	qCapH      float32
	qCapR      float32
	sx, sy, sz float32 // collision scale factors (default 1)
}

func (b *body3dObj) TypeName() string { return "Body3D" }

func (b *body3dObj) setFinalizer() {
	goruntime.SetFinalizer(b, func(o *body3dObj) {
		o.Free()
	})
}

func (b *body3dObj) TypeTag() uint16 { return heap.TagPhysicsBody }

func (b *body3dObj) Free() {
	b.release.Do(func() {
		if b.queryShape != nil {
			b.queryShape.Destroy()
			b.queryShape = nil
		}
		unregisterBufferBodyForCollision(b.bufferIndex)
		joltUnregisterBody(b.id)
		if b.id != nil {
			b.id.Destroy()
			b.id = nil
		}
	})
}

// BuilderObj owns a Jolt shape; destroy shape before COMMIT discards builder (host order).
type BuilderObj struct {
	Motion  jolt.MotionType
	Shape   *jolt.Shape
	Release heap.ReleaseOnce
	// Query template (rebuild after COMMIT for overlap tests).
	QKind          uint8 // 1 box, 2 sphere, 3 capsule
	QBox           jolt.Vec3
	QSphere        float32
	QCapH          float32
	QCapR          float32
	Friction       float32
	Restitution    float32
	LinearDamping  float32
	AngularDamping float32
	EnableCCD      bool
	// AllowedDOFs: 0 = Jolt default (all DOFs). Non-zero = EAllowedDOFs bitmask (e.g. platformer = 0x17).
	AllowedDOFs int
	// ForceSensor: trigger volumes (kinematic sensor); uses Jolt sensor object layer when true.
	ForceSensor bool
	// ObjectLayer: jolt.ObjectLayerAuto (-1) to derive from motion/sensor; else explicit layer (0..4, see physics_layers.h: ONE_WAY=4).
	ObjectLayer int
}

func (b *BuilderObj) TypeName() string { return "Body3DBuilder" }

func (b *BuilderObj) TypeTag() uint16 { return heap.TagPhysicsBuilder }

func (b *BuilderObj) Free() {
	b.Release.Do(func() {
		if b.Shape != nil {
			b.Shape.Destroy()
			b.Shape = nil
		}
	})
}

// ActiveJoltPhysics exposes the live Jolt world for charcontroller (linux only).
func ActiveJoltPhysics() *jolt.PhysicsSystem {
	joltMu.Lock()
	defer joltMu.Unlock()
	return joltSys
}

// GravityVec returns the last PHYSICS3D.SETGRAVITY values (for character updates).
func GravityVec() jolt.Vec3 {
	joltMu.Lock()
	defer joltMu.Unlock()
	return jolt.Vec3{X: gravX, Y: gravY, Z: gravZ}
}

// Helpers for Entity-Centric interaction (used by mbentity)

func ApplyImpulseToIndex(idx int, x, y, z float32) {
	joltMu.Lock()
	bi := joltBi
	joltBodyMu.Lock()
	id, ok := bufferIndexToBody[idx]
	joltBodyMu.Unlock()
	joltMu.Unlock()

	if !ok || bi == nil {
		return
	}
	bi.AddImpulse(id, jolt.Vec3{X: x, Y: y, Z: z})
	bi.ActivateBody(id)
}

func GetLinearVelocityToIndex(idx int) (x, y, z float32) {
	joltMu.Lock()
	bi := joltBi
	joltBodyMu.Lock()
	id, ok := bufferIndexToBody[idx]
	joltBodyMu.Unlock()
	joltMu.Unlock()

	if !ok || bi == nil {
		return 0, 0, 0
	}
	v := bi.GetLinearVelocity(id)
	return v.X, v.Y, v.Z
}

// GetBodyQuaternionForBufferIndex returns the Jolt body's world rotation for a matrix buffer index
// (see syncSharedBuffers / ENTITY.LINKPHYSBUFFER). ok is false when physics is down or the index is unknown.
func GetBodyQuaternionForBufferIndex(idx int) (x, y, z, w float32, ok bool) {
	joltMu.Lock()
	bi := joltBi
	joltBodyMu.Lock()
	id, have := bufferIndexToBody[idx]
	joltBodyMu.Unlock()
	joltMu.Unlock()

	if !have || bi == nil {
		return 0, 0, 0, 1, false
	}
	q := bi.GetRotation(id)
	return q.X, q.Y, q.Z, q.W, true
}

func SetVelocityToIndex(idx int, x, y, z float32) {
	joltMu.Lock()
	bi := joltBi
	joltBodyMu.Lock()
	id, ok := bufferIndexToBody[idx]
	joltBodyMu.Unlock()
	joltMu.Unlock()

	if !ok || bi == nil {
		return
	}
	bi.SetLinearVelocity(id, jolt.Vec3{X: x, Y: y, Z: z})
	bi.ActivateBody(id)
}

// SetPositionToIndex teleports a body to world position and activates it (traffic-cop path for ENTITY.SETPOSITION).
func SetPositionToIndex(idx int, x, y, z float32) {
	joltMu.Lock()
	bi := joltBi
	joltBodyMu.Lock()
	id, ok := bufferIndexToBody[idx]
	joltBodyMu.Unlock()
	joltMu.Unlock()

	if !ok || bi == nil {
		return
	}
	bi.SetPosition(id, jolt.Vec3{X: x, Y: y, Z: z})
	bi.ActivateBody(id)
}

func WakeIndex(idx int) {
	joltMu.Lock()
	bi := joltBi
	joltBodyMu.Lock()
	id, ok := bufferIndexToBody[idx]
	joltBodyMu.Unlock()
	joltMu.Unlock()

	if !ok || bi == nil {
		return
	}
	bi.ActivateBody(id)
}
func ApplyForceToIndex(idx int, x, y, z float32) {
	joltMu.Lock()
	bi := joltBi
	joltBodyMu.Lock()
	_, ok := bufferIndexToBody[idx]
	joltBodyMu.Unlock()
	joltMu.Unlock()

	if !ok || bi == nil {
		return
	}
	// No AddForce on BodyInterface in vendored jolt-go.
	_ = x
	_ = y
	_ = z
}

func SetFrictionToIndex(idx int, val float32) {
	joltMu.Lock()
	bi := joltBi
	joltBodyMu.Lock()
	id, ok := bufferIndexToBody[idx]
	joltBodyMu.Unlock()
	joltMu.Unlock()

	if !ok || bi == nil {
		return
	}
	bi.SetFriction(id, val)
}

func SetRestitutionToIndex(idx int, val float32) {
	joltMu.Lock()
	bi := joltBi
	joltBodyMu.Lock()
	id, ok := bufferIndexToBody[idx]
	joltBodyMu.Unlock()
	joltMu.Unlock()

	if !ok || bi == nil {
		return
	}
	bi.SetRestitution(id, val)
}

func SetGravityFactorToIndex(idx int, val float32) {
	joltMu.Lock()
	bi := joltBi
	joltBodyMu.Lock()
	_, ok := bufferIndexToBody[idx]
	joltBodyMu.Unlock()
	joltMu.Unlock()

	if !ok || bi == nil {
		return
	}
	// No SetGravityFactor on BodyInterface in vendored jolt-go.
	_ = val
}

func RotateToIndex(idx int, p, y, r float32) {
	joltMu.Lock()
	bi := joltBi
	joltBodyMu.Lock()
	_, ok := bufferIndexToBody[idx]
	joltBodyMu.Unlock()
	joltMu.Unlock()

	if !ok || bi == nil {
		return
	}
	// No SetRotation on BodyInterface in vendored jolt-go.
	_ = p
	_ = y
	_ = r
}

// GetStaticBodyRegistry exists on the heap static-body stub (see stub.go). Linux+Jolt uses real
// bodies without that Pos/Shape layout; mbentity host kinematic helpers only need this symbol to
// compile — iterating nil is a no-op (entity AABBs + Jolt paths cover collision).
func GetStaticBodyRegistry() map[heap.Handle]*body3dObj {
	return nil
}
