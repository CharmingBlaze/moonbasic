//go:build (linux || windows) && cgo

package mbphysics3d

import (
	"fmt"

	"github.com/bbitechnologies/jolt-go/jolt"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

// createBodyLowLevel is the unified allocator for modern physics bodies.
func createBodyLowLevel(h *heap.Store, shapeH heap.Handle, motion jolt.MotionType, sensor bool, tag uint16) (value.Value, error) {
	shObj, err := heap.Cast[*ShapeObj](h, shapeH)
	if err != nil {
		return value.Nil, fmt.Errorf("invalid shape handle")
	}
	joltMu.Lock()
	bi := joltBi
	joltMu.Unlock()
	if bi == nil {
		return value.Nil, fmt.Errorf("PHYSICS3D not started")
	}

	// Create body in Jolt
	id := bi.CreateBody(shObj.Shape, jolt.Vec3{}, motion, sensor, 0.5, 0.2, 0)
	if id == nil {
		return value.Nil, fmt.Errorf("failed to create Jolt body")
	}
	if sensor {
		// Set sensor flag
	}

	// Double check if we need a query shape
	var qshape *jolt.Shape
	switch shObj.Kind {
	case 1:
		qshape = jolt.CreateBox(jolt.Vec3{X: shObj.Dim1, Y: shObj.Dim2, Z: shObj.Dim3})
	case 2:
		qshape = jolt.CreateSphere(shObj.Dim1)
	case 3:
		qshape = jolt.CreateCapsule(shObj.Dim2/2-shObj.Dim1, shObj.Dim1)
	}

	body := &body3dObj{id: id, queryShape: qshape}
	body.setFinalizer()

	// Determine the exact tag for the handle
	bh, err := h.Alloc(body)
	if err != nil {
		if qshape != nil {
			qshape.Destroy()
		}
		id.Destroy()
		return value.Nil, err
	}
	joltRegisterBody(id, bh)

	// Register in matrix buffer for visual sync
	joltBodyMu.Lock()
	bidx := nextBufferIndex
	bufferIndexMap[id] = bidx
	bufferIndexToBody[bidx] = id
	body.bufferIndex = bidx
	nextBufferIndex++
	// Grow if needed
	if nextBufferIndex >= matrixBufferAlloc {
		matrixBufferAlloc += 1024
		newBuf := make([]float32, matrixBufferAlloc*16)
		copy(newBuf, matrixBuffer)
		matrixBuffer = newBuf
		newPrev := make([]float32, len(newBuf))
		if len(prevMatrixBuffer) > 0 {
			copy(newPrev, prevMatrixBuffer)
		}
		prevMatrixBuffer = newPrev
	}
	joltBodyMu.Unlock()

	registerBufferBodyForCollision(bidx, bh)

	return value.FromHandle(bh), nil
}

func knCreate(h *heap.Store, args []value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("KINEMATIC.CREATE expects shape handle")
	}
	return createBodyLowLevel(h, heap.Handle(args[0].IVal), jolt.MotionTypeKinematic, false, heap.TagKinematicBody)
}

func stCreate(h *heap.Store, args []value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("STATIC.CREATE expects shape handle")
	}
	return createBodyLowLevel(h, heap.Handle(args[0].IVal), jolt.MotionTypeStatic, false, heap.TagStaticBody)
}

func trCreate(h *heap.Store, args []value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("TRIGGER.CREATE expects shape handle")
	}
	return createBodyLowLevel(h, heap.Handle(args[0].IVal), jolt.MotionTypeStatic, true, heap.TagTriggerBody)
}

// Shared BODYREF methods

func brSetPos(h *heap.Store, args []value.Value) (value.Value, error) {
	if len(args) != 4 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("BODY.SETPOSITION expects (body, x, y, z)")
	}
	bo, err := castBody(h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	x, _ := args[1].ToFloat()
	y, _ := args[2].ToFloat()
	z, _ := args[3].ToFloat()
	joltMu.Lock()
	bi := joltBi
	joltMu.Unlock()
	if bi == nil {
		return value.Nil, nil
	}
	bi.SetPosition(bo.id, jolt.Vec3{X: float32(x), Y: float32(y), Z: float32(z)})
	return value.Nil, nil
}

func castBody(h *heap.Store, bh heap.Handle) (*body3dObj, error) {
	// Polymorphic cast
	if obj, err := heap.Cast[*body3dObj](h, bh); err == nil {
		return obj, nil
	}
	// Try specific types if needed, but since they all use body3dObj struct:
	obj, ok := h.Get(bh)
	if !ok {
		return nil, fmt.Errorf("invalid body handle")
	}
	if b, ok := obj.(*body3dObj); ok {
		return b, nil
	}
	return nil, fmt.Errorf("not a physics body handle")
}

func brFree(h *heap.Store, args []value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("BODY.FREE expects handle")
	}
	h.Free(heap.Handle(args[0].IVal))
	return value.Nil, nil
}
