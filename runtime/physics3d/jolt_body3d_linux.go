//go:build linux && cgo

package mbphysics3d

import (
	"fmt"
	"strings"
	"unsafe"

	"github.com/bbitechnologies/jolt-go/jolt"

	mbruntime "moonbasic/runtime"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

func parseMotion(s string) jolt.MotionType {
	switch strings.ToUpper(strings.TrimSpace(s)) {
	case "STATIC":
		return jolt.MotionTypeStatic
	case "KINEMATIC":
		return jolt.MotionTypeKinematic
	case "DYNAMIC":
		return jolt.MotionTypeDynamic
	default:
		return jolt.MotionTypeDynamic
	}
}

func BDAddBox(h *heap.Store, args []value.Value) (value.Value, error) {
	if len(args) != 4 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("BODY3D.ADDBOX expects (builder, hw, hh, hd)")
	}
	bu, err := heap.Cast[*BuilderObj](h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	hx, _ := args[1].ToFloat()
	hy, _ := args[2].ToFloat()
	hz, _ := args[3].ToFloat()
	if bu.Shape != nil {
		bu.Shape.Destroy()
	}
	bu.Shape = jolt.CreateBox(jolt.Vec3{X: float32(hx), Y: float32(hy), Z: float32(hz)})
	bu.QKind = 1
	bu.QBox = jolt.Vec3{X: float32(hx), Y: float32(hy), Z: float32(hz)}
	return value.Nil, nil
}

func BDAddSphere(h *heap.Store, args []value.Value) (value.Value, error) {
	if len(args) != 2 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("BODY3D.ADDSPHERE expects (builder, radius)")
	}
	bu, err := heap.Cast[*BuilderObj](h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	r, _ := args[1].ToFloat()
	if bu.Shape != nil {
		bu.Shape.Destroy()
	}
	bu.Shape = jolt.CreateSphere(float32(r))
	bu.QKind = 2
	bu.QSphere = float32(r)
	return value.Nil, nil
}

func BDAddCapsule(h *heap.Store, args []value.Value) (value.Value, error) {
	if len(args) != 3 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("BODY3D.ADDCAPSULE expects (builder, radius, height)")
	}
	bu, err := heap.Cast[*BuilderObj](h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	r, _ := args[1].ToFloat()
	h_val, _ := args[2].ToFloat()
	hh := float32(h_val)/2 - float32(r)
	if hh < 0.05 {
		hh = 0.05
	}
	if bu.Shape != nil {
		bu.Shape.Destroy()
	}
	bu.Shape = jolt.CreateCapsule(hh, float32(r))
	bu.QKind = 3
	bu.QCapH = hh
	bu.QCapR = float32(r)
	return value.Nil, nil
}

func bdAddMesh(m *Module, args []value.Value) (value.Value, error) {
	return value.Nil, fmt.Errorf("BODY3D.ADDMESH: requires Phase D model handle (not implemented)")
}

func BDCommit(h *heap.Store, args []value.Value) (value.Value, error) {
	if h == nil {
		return value.Nil, mbruntime.Errorf("BODY3D.COMMIT: heap not bound")
	}
	if len(args) != 4 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("BODY3D.COMMIT expects (builder, x, y, z)")
	}
	joltMu.Lock()
	bi := joltBi
	joltMu.Unlock()
	if bi == nil {
		return value.Nil, mbruntime.Errorf("BODY3D.COMMIT: physics not started")
	}
	bu, err := heap.Cast[*BuilderObj](h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	if bu.Shape == nil {
		return value.Nil, fmt.Errorf("BODY3D.COMMIT: no shape (call ADDBOX/ADDSPHERE/ADDCAPSULE first)")
	}
	x, _ := args[1].ToFloat()
	y, _ := args[2].ToFloat()
	z, _ := args[3].ToFloat()
	sh := bu.Shape
	bu.Shape = nil
	
	ccd := bu.EnableCCD && bu.Motion == jolt.MotionTypeDynamic
	
	id := bi.CreateBody(sh, jolt.Vec3{X: float32(x), Y: float32(y), Z: float32(z)}, bu.Motion, ccd)
	sh.Destroy()
	
	if bu.Friction > 0 {
		bi.SetFriction(id, bu.Friction)
	}
	if bu.Restitution > 0 {
		bi.SetRestitution(id, bu.Restitution)
	}

	var qshape *jolt.Shape
	switch bu.QKind {
	case 1:
		qshape = jolt.CreateBox(bu.QBox)
	case 2:
		qshape = jolt.CreateSphere(bu.QSphere)
	case 3:
		qshape = jolt.CreateCapsule(bu.QCapH, bu.QCapR)
	}

	h.Free(heap.Handle(args[0].IVal))
	body := &body3dObj{id: id, queryShape: qshape}
	body.setFinalizer()
	bh, err := h.Alloc(body)
	if err != nil {
		if qshape != nil {
			qshape.Destroy()
		}
		if id != nil {
			id.Destroy()
		}
		return value.Nil, err
	}
	joltRegisterBody(id, bh)

	joltBodyMu.Lock()
	bidx := nextBufferIndex
	bufferIndexMap[uintptr(unsafe.Pointer(id))] = bidx
	bufferIndexToBody[bidx] = uintptr(unsafe.Pointer(id))
	body.bufferIndex = bidx
	nextBufferIndex++
	// Grow if needed
	if nextBufferIndex >= matrixBufferAlloc {
		matrixBufferAlloc += 1024
		newBuf := make([]float32, matrixBufferAlloc*16)
		copy(newBuf, matrixBuffer)
		matrixBuffer = newBuf
	}
	joltBodyMu.Unlock()

	registerBufferBodyForCollision(bidx, bh)

	return value.FromHandle(bh), nil
}

func bdSetPos(m *Module, args []value.Value) (value.Value, error) {
	if len(args) != 4 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("BODY3D.SETPOS expects (body, x, y, z)")
	}
	joltMu.Lock()
	bi := joltBi
	joltMu.Unlock()
	if bi == nil {
		return value.Nil, mbruntime.Errorf("BODY3D.SETPOS: physics not started")
	}
	bo, err := heap.Cast[*body3dObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	x, _ := args[1].ToFloat()
	y, _ := args[2].ToFloat()
	z, _ := args[3].ToFloat()
	bi.SetPosition(bo.id, jolt.Vec3{X: float32(x), Y: float32(y), Z: float32(z)})
	return value.Nil, nil
}

func bdActivate(m *Module, args []value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("BODY3D.ACTIVATE expects body handle")
	}
	joltMu.Lock()
	bi := joltBi
	joltMu.Unlock()
	if bi == nil {
		return value.Nil, mbruntime.Errorf("BODY3D.ACTIVATE: physics not started")
	}
	bo, err := heap.Cast[*body3dObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	bi.ActivateBody(bo.id)
	return value.Nil, nil
}

func bdDeactivate(m *Module, args []value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("BODY3D.DEACTIVATE expects body handle")
	}
	joltMu.Lock()
	bi := joltBi
	joltMu.Unlock()
	if bi == nil {
		return value.Nil, mbruntime.Errorf("BODY3D.DEACTIVATE: physics not started")
	}
	bo, err := heap.Cast[*body3dObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	bi.DeactivateBody(bo.id)
	return value.Nil, nil
}

func bdGetPos(m *Module, args []value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, mbruntime.Errorf("BODY3D.GETPOS: heap not bound")
	}
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("BODY3D.GETPOS expects body handle")
	}
	joltMu.Lock()
	bi := joltBi
	joltMu.Unlock()
	if bi == nil {
		return value.Nil, mbruntime.Errorf("BODY3D.GETPOS: physics not started")
	}
	bo, err := heap.Cast[*body3dObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	p := bi.GetPosition(bo.id)
	arr, err := heap.NewArray([]int64{3})
	if err != nil {
		return value.Nil, err
	}
	_ = arr.Set([]int64{0}, float64(p.X))
	_ = arr.Set([]int64{1}, float64(p.Y))
	_ = arr.Set([]int64{2}, float64(p.Z))
	ph, err := m.h.Alloc(arr)
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(ph), nil
}

// bdGetRotZero returns a 3-element array [pitch, yaw, roll]; placeholder until GetRotation is exposed on BodyInterface.
func bdGetRotZero(m *Module, args []value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, mbruntime.Errorf("BODY3D.GETROT: heap not bound")
	}
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("BODY3D.GETROT expects body handle")
	}
	arr, err := heap.NewArray([]int64{3})
	if err != nil {
		return value.Nil, err
	}
	_ = arr.Set([]int64{0}, 0)
	_ = arr.Set([]int64{1}, 0)
	_ = arr.Set([]int64{2}, 0)
	ph, err := m.h.Alloc(arr)
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(ph), nil
}

func bdSetRotation(m *Module, args []value.Value) (value.Value, error) {
	if len(args) != 4 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("BODY3D.SETROT expects (body, p, y, r)")
	}
	joltMu.Lock()
	bi := joltBi
	joltMu.Unlock()
	if bi == nil { return value.Nil, mbruntime.Errorf("physics not started") }
	bo, err := heap.Cast[*body3dObj](m.h, heap.Handle(args[0].IVal))
	if err != nil { return value.Nil, err }
	p, _ := args[1].ToFloat()
	y, _ := args[2].ToFloat()
	r, _ := args[3].ToFloat()
	q := rl.QuaternionFromEuler(float32(p), float32(y), float32(r))
	bi.SetRotation(bo.id, jolt.Quat{X: q.X, Y: q.Y, Z: q.Z, W: q.W}, jolt.ActivationActivate)
	return value.Nil, nil
}

func bdSetFriction(m *Module, args []value.Value) (value.Value, error) {
	if len(args) != 2 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("BODY3D.SETFRICTION expects (body, val#)")
	}
	joltMu.Lock()
	bi := joltBi
	joltMu.Unlock()
	if bi == nil { return value.Nil, nil }
	bo, err := heap.Cast[*body3dObj](m.h, heap.Handle(args[0].IVal))
	if err != nil { return value.Nil, err }
	val, _ := args[1].ToFloat()
	bi.SetFriction(bo.id, float32(val))
	return value.Nil, nil
}

func bdSetRestitution(m *Module, args []value.Value) (value.Value, error) {
	if len(args) != 2 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("BODY3D.SETRESTITUTION expects (body, val#)")
	}
	joltMu.Lock()
	bi := joltBi
	joltMu.Unlock()
	if bi == nil { return value.Nil, nil }
	bo, err := heap.Cast[*body3dObj](m.h, heap.Handle(args[0].IVal))
	if err != nil { return value.Nil, err }
	val, _ := args[1].ToFloat()
	bi.SetRestitution(bo.id, float32(val))
	return value.Nil, nil
}

func bdApplyForce(m *Module, args []value.Value) (value.Value, error) {
	if len(args) != 4 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("BODY3D.APPLYFORCE expects (body, x, y, z)")
	}
	joltMu.Lock()
	bi := joltBi
	joltMu.Unlock()
	if bi == nil { return value.Nil, nil }
	bo, err := heap.Cast[*body3dObj](m.h, heap.Handle(args[0].IVal))
	if err != nil { return value.Nil, err }
	x, _ := args[1].ToFloat()
	y, _ := args[2].ToFloat()
	z, _ := args[3].ToFloat()
	bi.AddForce(bo.id, jolt.Vec3{X: float32(x), Y: float32(y), Z: float32(z)})
	return value.Nil, nil
}

func bdApplyImpulse(m *Module, args []value.Value) (value.Value, error) {
	if len(args) != 4 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("BODY3D.APPLYIMPULSE expects (body, x, y, z)")
	}
	joltMu.Lock()
	bi := joltBi
	joltMu.Unlock()
	if bi == nil { return value.Nil, nil }
	bo, err := heap.Cast[*body3dObj](m.h, heap.Handle(args[0].IVal))
	if err != nil { return value.Nil, err }
	x, _ := args[1].ToFloat()
	y, _ := args[2].ToFloat()
	z, _ := args[3].ToFloat()
	bi.AddImpulse(bo.id, jolt.Vec3{X: float32(x), Y: float32(y), Z: float32(z)})
	return value.Nil, nil
}

func bdSetLinearVel(m *Module, args []value.Value) (value.Value, error) {
	if len(args) != 4 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("BODY3D.SETLINEARVEL expects (body, x, y, z)")
	}
	joltMu.Lock()
	bi := joltBi
	joltMu.Unlock()
	if bi == nil { return value.Nil, nil }
	bo, err := heap.Cast[*body3dObj](m.h, heap.Handle(args[0].IVal))
	if err != nil { return value.Nil, err }
	x, _ := args[1].ToFloat()
	y, _ := args[2].ToFloat()
	z, _ := args[3].ToFloat()
	bi.SetLinearVelocity(bo.id, jolt.Vec3{X: float32(x), Y: float32(y), Z: float32(z)})
	return value.Nil, nil
}

func bdNoOp(m *Module, args []value.Value) (value.Value, error) {
	return value.Nil, nil
}

func bdAxis(m *Module, args []value.Value, axis int) (value.Value, error) {
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("BODY3D axis getter expects handle")
	}
	joltMu.Lock()
	bi := joltBi
	joltMu.Unlock()
	if bi == nil {
		return value.FromFloat(0), nil
	}
	bo, err := heap.Cast[*body3dObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.FromFloat(0), nil
	}
	p := bi.GetPosition(bo.id)
	switch axis {
	case 0:
		return value.FromFloat(float64(p.X)), nil
	case 1:
		return value.FromFloat(float64(p.Y)), nil
	default:
		return value.FromFloat(float64(p.Z)), nil
	}
}

func BDBufferIndex(h *heap.Store, args []value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("BODY3D.BUFFERINDEX expects handle")
	}
	bo, err := heap.Cast[*body3dObj](h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	return value.FromInt(int64(bo.bufferIndex)), nil
}

func bdFree(m *Module, args []value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("BODY3D.FREE expects handle")
	}
	m.h.Free(heap.Handle(args[0].IVal))
	return value.Nil, nil
}

func bdCollided3D(m *Module, args []value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("BODY3D.COLLIDED expects body handle")
	}
	bo, err := heap.Cast[*body3dObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	if bo.queryShape == nil {
		return value.FromInt(0), nil
	}
	joltMu.Lock()
	ps := joltSys
	bi := joltBi
	joltMu.Unlock()
	if ps == nil || bi == nil {
		return value.Nil, mbruntime.Errorf("BODY3D.COLLIDED: physics not started")
	}
	pos := bi.GetPosition(bo.id)
	hits := ps.CollideShapeGetHits(bo.queryShape, pos, 8, 1e-3)
	self := uintptr(unsafe.Pointer(bo.id))
	for _, hit := range hits {
		if hit.BodyID == nil {
			continue
		}
		if uintptr(unsafe.Pointer(hit.BodyID)) == self {
			continue
		}
		return value.FromInt(1), nil
	}
	return value.FromInt(0), nil
}

func bdCollisionOther3D(m *Module, args []value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("BODY3D.COLLISIONOTHER expects body handle")
	}
	bo, err := heap.Cast[*body3dObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	if bo.queryShape == nil {
		return value.FromHandle(0), nil
	}
	joltMu.Lock()
	ps := joltSys
	bi := joltBi
	joltMu.Unlock()
	if ps == nil || bi == nil {
		return value.Nil, mbruntime.Errorf("BODY3D.COLLISIONOTHER: physics not started")
	}
	pos := bi.GetPosition(bo.id)
	hits := ps.CollideShapeGetHits(bo.queryShape, pos, 8, 1e-3)
	self := uintptr(unsafe.Pointer(bo.id))
	for _, hit := range hits {
		if hit.BodyID == nil {
			continue
		}
		if uintptr(unsafe.Pointer(hit.BodyID)) == self {
			continue
		}
		if h, ok := joltLookupHandle(hit.BodyID); ok {
			return value.FromHandle(h), nil
		}
	}
	return value.FromHandle(0), nil
}

func bdCollisionPoint3D(m *Module, args []value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("BODY3D.COLLISIONPOINT expects body handle")
	}
	bo, err := heap.Cast[*body3dObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	if bo.queryShape == nil {
		return bdCollisionZeroArray3(m)
	}
	joltMu.Lock()
	ps := joltSys
	bi := joltBi
	joltMu.Unlock()
	if ps == nil || bi == nil {
		return value.Nil, mbruntime.Errorf("BODY3D.COLLISIONPOINT: physics not started")
	}
	pos := bi.GetPosition(bo.id)
	hits := ps.CollideShapeGetHits(bo.queryShape, pos, 1, 1e-3)
	self := uintptr(unsafe.Pointer(bo.id))
	for _, hit := range hits {
		if hit.BodyID == nil || uintptr(unsafe.Pointer(hit.BodyID)) == self {
			continue
		}
		arr, err := heap.NewArray([]int64{3})
		if err != nil {
			return value.Nil, err
		}
		_ = arr.Set([]int64{0}, float64(hit.ContactPoint.X))
		_ = arr.Set([]int64{1}, float64(hit.ContactPoint.Y))
		_ = arr.Set([]int64{2}, float64(hit.ContactPoint.Z))
		id, err := m.h.Alloc(arr)
		if err != nil {
			return value.Nil, err
		}
		return value.FromHandle(id), nil
	}
	return bdCollisionZeroArray3(m)
}

func bdCollisionNormal3D(m *Module, args []value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("BODY3D.COLLISIONNORMAL expects body handle")
	}
	// CollideShapeGetHits does not return contact normals; use PHYSICS3D.RAYCAST for surface normals.
	return bdCollisionZeroArray3(m)
}

func bdCollisionZeroArray3(m *Module) (value.Value, error) {
	arr, err := heap.NewArray([]int64{3})
	if err != nil {
		return value.Nil, err
	}
	_ = arr.Set([]int64{0}, 0)
	_ = arr.Set([]int64{1}, 0)
	_ = arr.Set([]int64{2}, 1)
	id, err := m.h.Alloc(arr)
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}
