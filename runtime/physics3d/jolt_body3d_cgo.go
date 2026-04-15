//go:build (linux || windows) && cgo

package mbphysics3d

import (
	"fmt"
	"strings"
	"unsafe"

	"github.com/bbitechnologies/jolt-go/jolt"
	rl "github.com/gen2brain/raylib-go/raylib"

	mbruntime "moonbasic/runtime"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

func parseMotion(s string) jolt.MotionType {
	switch strings.ToUpper(strings.TrimSpace(s)) {
	case "STATIC":
		return jolt.MotionTypeStatic
	case "KINEMATIC", "TRIGGER":
		return jolt.MotionTypeKinematic
	case "DYNAMIC":
		return jolt.MotionTypeDynamic
	default:
		return jolt.MotionTypeDynamic
	}
}

func phCreateBody(m *Module, motion string) (value.Value, error) {
	if m.h == nil {
		return value.Nil, mbruntime.Errorf("BODY3D.CREATE: heap not bound")
	}
	mType := parseMotion(motion)

	// Professional Architect Defaults (Tri-Tier Architecture)
	friction := float32(0.5) // Default for props
	restitution := float32(0.0)
	linearDamping := float32(0.0)
	angularDamping := float32(0.0)

	switch mType {
	case jolt.MotionTypeStatic:
		friction = 0.5 // High friction for stage floors
	case jolt.MotionTypeDynamic:
		linearDamping = 0.05 // Prevents sliding on ice
		angularDamping = 0.05
	}

	bu := &BuilderObj{
		Motion:         mType,
		Friction:       friction,
		Restitution:    restitution,
		LinearDamping:  linearDamping,
		AngularDamping: angularDamping,
		AllowedDOFs:    0,
		ObjectLayer:    -1,
	}
	if strings.EqualFold(strings.TrimSpace(motion), "TRIGGER") {
		bu.ForceSensor = true
	}
	id, err := m.h.Alloc(bu)
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}

func (m *Module) bdAddBox(args []value.Value) (value.Value, error) {
	if len(args) != 4 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("BODY3D.ADDBOX expects (builder, hw, hh, hd)")
	}
	bu, err := heap.Cast[*BuilderObj](m.h, heap.Handle(args[0].IVal))
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

func (m *Module) bdAddSphere(args []value.Value) (value.Value, error) {
	if len(args) != 2 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("BODY3D.ADDSPHERE expects (builder, radius)")
	}
	bu, err := heap.Cast[*BuilderObj](m.h, heap.Handle(args[0].IVal))
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

func (m *Module) bdAddCapsule(args []value.Value) (value.Value, error) {
	if len(args) != 3 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("BODY3D.ADDCAPSULE expects (builder, radius, height)")
	}
	bu, err := heap.Cast[*BuilderObj](m.h, heap.Handle(args[0].IVal))
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

func (m *Module) bdAddMesh(args []value.Value) (value.Value, error) {
	if len(args) != 2 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("BODY3D.ADDMESH expects (builder, entityID)")
	}
	bu, err := heap.Cast[*BuilderObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	eid, _ := args[1].ToInt()

	if m.meshLookup == nil {
		return value.Nil, fmt.Errorf("BODY3D.ADDMESH: mesh lookup not wired (engine bridge missing)")
	}
	meshes := m.meshLookup(eid)
	if len(meshes) == 0 {
		return value.Nil, fmt.Errorf("BODY3D.ADDMESH: entity %d has no meshes", eid)
	}

	var allVerts []jolt.Vec3
	var allIndices []int32

	for _, mesh := range meshes {
		off := int32(len(allVerts))
		vCount := int(mesh.VertexCount)
		iCount := int(mesh.TriangleCount) * 3

		// vertices are X,Y,Z floats
		vPtr := (*[1 << 30]float32)(unsafe.Pointer(mesh.Vertices))[: vCount*3 : vCount*3]
		for i := 0; i < vCount; i++ {
			allVerts = append(allVerts, jolt.Vec3{X: vPtr[i*3], Y: vPtr[i*3+1], Z: vPtr[i*3+2]})
		}

		// indices are uint16
		iPtr := (*[1 << 30]uint16)(unsafe.Pointer(mesh.Indices))[:iCount:iCount]
		for i := 0; i < iCount; i++ {
			allIndices = append(allIndices, int32(iPtr[i])+off)
		}
	}

	if bu.Shape != nil {
		bu.Shape.Destroy()
	}
	bu.Shape = jolt.CreateMesh(allVerts, allIndices)
	bu.QKind = 0
	return value.Nil, nil
}

func (m *Module) bdCommit(args []value.Value) (value.Value, error) {
	if m.h == nil {
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
	bu, err := heap.Cast[*BuilderObj](m.h, heap.Handle(args[0].IVal))
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

	motion := bu.Motion
	isSensor := bu.ForceSensor

	var id *jolt.BodyID
	if bu.ObjectLayer >= 0 {
		id = bi.CreateBody(sh, jolt.Vec3{X: float32(x), Y: float32(y), Z: float32(z)}, motion, isSensor,
			bu.Friction, bu.Restitution, bu.AllowedDOFs, bu.ObjectLayer)
	} else {
		id = bi.CreateBody(sh, jolt.Vec3{X: float32(x), Y: float32(y), Z: float32(z)}, motion, isSensor,
			bu.Friction, bu.Restitution, bu.AllowedDOFs)
	}
	_ = bu.EnableCCD && motion == jolt.MotionTypeDynamic // reserved for future CCD exposure on body creation

	// Apply Tri-Tier Damping
	if motion == jolt.MotionTypeDynamic && joltSys != nil {
		joltSys.SetBodyDamping(id, bu.LinearDamping, bu.AngularDamping)
	}

	sh.Destroy()

	var qshape *jolt.Shape
	switch bu.QKind {
	case 1:
		qshape = jolt.CreateBox(bu.QBox)
	case 2:
		qshape = jolt.CreateSphere(bu.QSphere)
	case 3:
		qshape = jolt.CreateCapsule(bu.QCapH, bu.QCapR)
	}

	m.h.Free(heap.Handle(args[0].IVal))
	body := &body3dObj{id: id, queryShape: qshape, motion: motion}
	body.setFinalizer()
	bh, err := m.h.Alloc(body)
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
	joltMarkBodyDynamic(id, motion == jolt.MotionTypeDynamic)

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

func (m *Module) bdSetPos(args []value.Value) (value.Value, error) {
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

func (m *Module) bdActivate(args []value.Value) (value.Value, error) {
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

func (m *Module) bdDeactivate(args []value.Value) (value.Value, error) {
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

func (m *Module) bdGetPos(args []value.Value) (value.Value, error) {
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

// bdGetRotZero returns a 3-element array [pitch, yaw, roll] in ENTITY rotation order (QuaternionToEuler Y,Z,X).
func (m *Module) bdGetRotZero(args []value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, mbruntime.Errorf("BODY3D.GETROT: heap not bound")
	}
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("BODY3D.GETROT expects body handle")
	}
	bo, err := heap.Cast[*body3dObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	joltMu.Lock()
	bi := joltBi
	joltMu.Unlock()
	if bi == nil {
		return value.Nil, mbruntime.Errorf("BODY3D.GETROT: physics not started")
	}
	q := bi.GetRotation(bo.id)
	v := rl.QuaternionToEuler(rl.Quaternion{X: q.X, Y: q.Y, Z: q.Z, W: q.W})
	pitch, yaw, roll := v.Y, v.Z, v.X
	arr, err := heap.NewArray([]int64{3})
	if err != nil {
		return value.Nil, err
	}
	_ = arr.Set([]int64{0}, float64(pitch))
	_ = arr.Set([]int64{1}, float64(yaw))
	_ = arr.Set([]int64{2}, float64(roll))
	ph, err := m.h.Alloc(arr)
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(ph), nil
}

func (m *Module) bdSetRotation(args []value.Value) (value.Value, error) {
	if len(args) != 4 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("BODY3D.SETROT expects (body, p, y, r)")
	}
	joltMu.Lock()
	bi := joltBi
	joltMu.Unlock()
	if bi == nil {
		return value.Nil, mbruntime.Errorf("physics not started")
	}
	if _, err := heap.Cast[*body3dObj](m.h, heap.Handle(args[0].IVal)); err != nil {
		return value.Nil, err
	}
	_, _ = args[1].ToFloat()
	_, _ = args[2].ToFloat()
	_, _ = args[3].ToFloat()
	// SetRotation is not exposed on jolt-go BodyInterface (v0.8.x).
	return value.Nil, nil
}

func (m *Module) bdSetMass(args []value.Value) (value.Value, error) {
	if len(args) != 2 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("BODY3D.SETMASS expects (body, mass#)")
	}
	// Mass calculation handled by Jolt at Commit/Setup; dynamic update via MassProperties not exposed.
	return value.Nil, nil
}

func (m *Module) bdSetFriction(args []value.Value) (value.Value, error) {
	if len(args) != 2 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("BODY3D.SETFRICTION expects (body, val#)")
	}
	joltMu.Lock()
	bi := joltBi
	joltMu.Unlock()
	if bi == nil {
		return value.Nil, nil
	}
	bo, err := heap.Cast[*body3dObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	v, _ := args[1].ToFloat()
	bi.SetFriction(bo.id, float32(v))
	return value.Nil, nil
}

func (m *Module) bdSetRestitution(args []value.Value) (value.Value, error) {
	if len(args) != 2 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("BODY3D.SETRESTITUTION expects (body, val#)")
	}
	joltMu.Lock()
	bi := joltBi
	joltMu.Unlock()
	if bi == nil {
		return value.Nil, nil
	}
	bo, err := heap.Cast[*body3dObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	v, _ := args[1].ToFloat()
	bi.SetRestitution(bo.id, float32(v))
	return value.Nil, nil
}

func (m *Module) bdApplyForce(args []value.Value) (value.Value, error) {
	if len(args) != 4 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("BODY3D.APPLYFORCE expects (body, x, y, z)")
	}
	joltMu.Lock()
	bi := joltBi
	joltMu.Unlock()
	if bi == nil {
		return value.Nil, nil
	}
	bo, err := heap.Cast[*body3dObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	x, _ := args[1].ToFloat()
	y, _ := args[2].ToFloat()
	z, _ := args[3].ToFloat()

	// Implementation note: vendored Jolt wrapper only has AddImpulse.
	// Force is impulse / dt. So impulse = force * dt.
	dt := float32(m.fixedStep)
	if dt <= 0 {
		dt = 1.0 / 60.0
	}

	bi.AddImpulse(bo.id, jolt.Vec3{X: float32(x) * dt, Y: float32(y) * dt, Z: float32(z) * dt})
	bi.ActivateBody(bo.id)
	return value.Nil, nil
}

func (m *Module) bdApplyImpulse(args []value.Value) (value.Value, error) {
	if len(args) != 4 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("BODY3D.APPLYIMPULSE expects (body, x, y, z)")
	}
	joltMu.Lock()
	bi := joltBi
	joltMu.Unlock()
	if bi == nil {
		return value.Nil, nil
	}
	bo, err := heap.Cast[*body3dObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	x, _ := args[1].ToFloat()
	y, _ := args[2].ToFloat()
	z, _ := args[3].ToFloat()
	bi.AddImpulse(bo.id, jolt.Vec3{X: float32(x), Y: float32(y), Z: float32(z)})
	bi.ActivateBody(bo.id)
	return value.Nil, nil
}

func (m *Module) bdSetLinearVel(args []value.Value) (value.Value, error) {
	if len(args) != 4 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("BODY3D.SETLINEARVEL expects (body, x, y, z)")
	}
	joltMu.Lock()
	bi := joltBi
	joltMu.Unlock()
	if bi == nil {
		return value.Nil, nil
	}
	bo, err := heap.Cast[*body3dObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	x, _ := args[1].ToFloat()
	y, _ := args[2].ToFloat()
	z, _ := args[3].ToFloat()
	bi.SetLinearVelocity(bo.id, jolt.Vec3{X: float32(x), Y: float32(y), Z: float32(z)})
	bi.ActivateBody(bo.id)
	return value.Nil, nil
}

func (m *Module) bdNoOp(args []value.Value) (value.Value, error) {
	return value.Nil, nil
}

func (m *Module) bdAxis(args []value.Value, axis int) (value.Value, error) {
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

func (m *Module) bdBufferIndex(args []value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("BODY3D.BUFFERINDEX expects handle")
	}
	bo, err := heap.Cast[*body3dObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	return value.FromInt(int64(bo.bufferIndex)), nil
}

func (m *Module) bdFree(args []value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("BODY3D.FREE expects handle")
	}
	m.h.Free(heap.Handle(args[0].IVal))
	return value.Nil, nil
}

func (m *Module) brSetPos(args []value.Value) (value.Value, error) {
	return m.bdSetPos(args)
}

func (m *Module) brSetLayer(args []value.Value) (value.Value, error) {
	return value.Nil, nil
}

func (m *Module) brEnableColl(args []value.Value) (value.Value, error) {
	return value.Nil, nil
}

func (m *Module) brFree(args []value.Value) (value.Value, error) {
	return m.bdFree(args)
}

func (m *Module) knCreate(args []value.Value) (value.Value, error) {
	return phCreateBody(m, "KINEMATIC")
}

func (m *Module) stCreate(args []value.Value) (value.Value, error) {
	return phCreateBody(m, "STATIC")
}

func (m *Module) trCreate(args []value.Value) (value.Value, error) {
	return phCreateBody(m, "TRIGGER")
}

func (m *Module) bdCollided3D(args []value.Value) (value.Value, error) {
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
	for _, hit := range hits {
		if hit.BodyID == nil {
			continue
		}
		if hit.BodyID == bo.id {
			continue
		}
		return value.FromInt(1), nil
	}
	return value.FromInt(0), nil
}

func (m *Module) bdCollisionOther3D(args []value.Value) (value.Value, error) {
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
	for _, hit := range hits {
		if hit.BodyID == nil {
			continue
		}
		if hit.BodyID == bo.id {
			continue
		}
		if h, ok := joltLookupHandle(hit.BodyID); ok {
			return value.FromHandle(h), nil
		}
	}
	return value.FromHandle(0), nil
}

func (m *Module) bdCollisionPoint3D(args []value.Value) (value.Value, error) {
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
	for _, hit := range hits {
		if hit.BodyID == nil || hit.BodyID == bo.id {
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

func (m *Module) bdCollisionNormal3D(args []value.Value) (value.Value, error) {
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
