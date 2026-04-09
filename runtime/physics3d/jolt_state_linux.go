//go:build linux && cgo

package mbphysics3d

import (
	goruntime "runtime"
	"sync"
	"unsafe"

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
	nextBufferIndex   int
	bufferIndexMap    map[uintptr]int // *BodyID -> index
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

// builderObj owns a Jolt shape; destroy shape before COMMIT discards builder (host order).
type builderObj struct {
	motion  jolt.MotionType
	shape   *jolt.Shape
	release heap.ReleaseOnce
	// Query template (rebuild after COMMIT for overlap tests).
	qKind   uint8 // 1 box, 2 sphere, 3 capsule
	qBox    jolt.Vec3
	qSphere float32
	qCapH   float32
	qCapR   float32
}

func (b *builderObj) TypeName() string { return "Body3DBuilder" }

func (b *builderObj) TypeTag() uint16 { return heap.TagPhysicsBuilder }

func (b *builderObj) Free() {
	b.release.Do(func() {
		if b.shape != nil {
			b.shape.Destroy()
			b.shape = nil
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
