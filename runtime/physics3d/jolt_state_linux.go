//go:build linux && cgo

package mbphysics3d

import (
	goruntime "runtime"
	"sync"
	"unsafe"

	"github.com/bbitechnologies/jolt-go/jolt"
	rl "github.com/gen2brain/raylib-go/raylib"

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
	nextBufferIndex   int
	bufferIndexMap    map[uintptr]int // *BodyID -> index
	bufferIndexToBody map[int]uintptr // index -> *BodyID
	matrixBufferAlloc int
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
	joltBodyToHandle map[uintptr]heap.Handle // *BodyID pointer -> VM handle
)

func joltRegisterBody(id *jolt.BodyID, h heap.Handle) {
	if id == nil {
		return
	}
	joltBodyMu.Lock()
	defer joltBodyMu.Unlock()
	if joltBodyToHandle == nil {
		joltBodyToHandle = make(map[uintptr]heap.Handle)
	}
	joltBodyToHandle[uintptr(unsafe.Pointer(id))] = h
}

func joltUnregisterBody(id *jolt.BodyID) {
	if id == nil {
		return
	}
	joltBodyMu.Lock()
	defer joltBodyMu.Unlock()
	delete(joltBodyToHandle, uintptr(unsafe.Pointer(id)))
}

func joltLookupHandle(id *jolt.BodyID) (heap.Handle, bool) {
	if id == nil {
		return 0, false
	}
	joltBodyMu.Lock()
	defer joltBodyMu.Unlock()
	h, ok := joltBodyToHandle[uintptr(unsafe.Pointer(id))]
	return h, ok
}

type body3dObj struct {
	id          *jolt.BodyID
	queryShape  *jolt.Shape // duplicate shape for CollideShape queries (same geometry as body)
	bufferIndex int         // index into sharedMatrixBuffer
	release     heap.ReleaseOnce
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
	QKind   uint8 // 1 box, 2 sphere, 3 capsule
	QBox    jolt.Vec3
	QSphere float32
	QCapH   float32
	QCapR   float32
	Friction    float32
	Restitution float32
	EnableCCD   bool
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
	ptr, ok := bufferIndexToBody[idx]
	joltBodyMu.Unlock()
	joltMu.Unlock()

	if !ok || bi == nil { return }
	id := (*jolt.BodyID)(unsafe.Pointer(ptr))
	bi.AddImpulse(id, jolt.Vec3{X: x, Y: y, Z: z})
}

func SetVelocityToIndex(idx int, x, y, z float32) {
	joltMu.Lock()
	bi := joltBi
	joltBodyMu.Lock()
	ptr, ok := bufferIndexToBody[idx]
	joltBodyMu.Unlock()
	joltMu.Unlock()

	if !ok || bi == nil { return }
	id := (*jolt.BodyID)(unsafe.Pointer(ptr))
	bi.SetLinearVelocity(id, jolt.Vec3{X: x, Y: y, Z: z})
}

func WakeIndex(idx int) {
	joltMu.Lock()
	bi := joltBi
	joltBodyMu.Lock()
	ptr, ok := bufferIndexToBody[idx]
	joltBodyMu.Unlock()
	joltMu.Unlock()

	if !ok || bi == nil { return }
	id := (*jolt.BodyID)(unsafe.Pointer(ptr))
	bi.ActivateBody(id)
}
func ApplyForceToIndex(idx int, x, y, z float32) {
	joltMu.Lock()
	bi := joltBi
	joltBodyMu.Lock()
	ptr, ok := bufferIndexToBody[idx]
	joltBodyMu.Unlock()
	joltMu.Unlock()

	if !ok || bi == nil { return }
	id := (*jolt.BodyID)(unsafe.Pointer(ptr))
	bi.AddForce(id, jolt.Vec3{X: x, Y: y, Z: z})
}

func SetFrictionToIndex(idx int, val float32) {
	joltMu.Lock()
	bi := joltBi
	joltBodyMu.Lock()
	ptr, ok := bufferIndexToBody[idx]
	joltBodyMu.Unlock()
	joltMu.Unlock()

	if !ok || bi == nil { return }
	id := (*jolt.BodyID)(unsafe.Pointer(ptr))
	bi.SetFriction(id, val)
}

func SetRestitutionToIndex(idx int, val float32) {
	joltMu.Lock()
	bi := joltBi
	joltBodyMu.Lock()
	ptr, ok := bufferIndexToBody[idx]
	joltBodyMu.Unlock()
	joltMu.Unlock()

	if !ok || bi == nil { return }
	id := (*jolt.BodyID)(unsafe.Pointer(ptr))
	bi.SetRestitution(id, val)
}

func SetGravityFactorToIndex(idx int, val float32) {
	joltMu.Lock()
	bi := joltBi
	joltBodyMu.Lock()
	ptr, ok := bufferIndexToBody[idx]
	joltBodyMu.Unlock()
	joltMu.Unlock()

	if !ok || bi == nil { return }
	id := (*jolt.BodyID)(unsafe.Pointer(ptr))
	bi.SetGravityFactor(id, val)
}

func RotateToIndex(idx int, p, y, r float32) {
	joltMu.Lock()
	bi := joltBi
	joltBodyMu.Lock()
	ptr, ok := bufferIndexToBody[idx]
	joltBodyMu.Unlock()
	joltMu.Unlock()

	if !ok || bi == nil { return }
	id := (*jolt.BodyID)(unsafe.Pointer(ptr))
	
	q := rl.QuaternionFromEuler(p, y, r)
	bi.SetRotation(id, jolt.Quat{X: q.X, Y: q.Y, Z: q.Z, W: q.W}, jolt.ActivationActivate)
}
