//go:build (linux || windows) && cgo

package mbphysics3d

import (
	"fmt"

	"github.com/bbitechnologies/jolt-go/jolt"
	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/internal/joltwasm"
	"moonbasic/runtime/mbmatrix"
	mbruntime "moonbasic/runtime"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

func (m *Module) phStart(args []value.Value) (value.Value, error) {
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("PHYSICS3D.START expects 0 arguments")
	}
	joltMu.Lock()
	defer joltMu.Unlock()
	if joltSys != nil {
		return value.Nil, mbruntime.Errorf("PHYSICS3D.START: already started")
	}
	if !joltCoreInit {
		if err := jolt.Init(); err != nil {
			return value.Nil, err
		}
		joltCoreInit = true
	}
	joltSys = jolt.NewPhysicsSystem()
	joltBi = joltSys.GetBodyInterface()
	collRules = nil
	collPending = nil
	joltBodyMu.Lock()
	joltBodyToHandle = make(map[*jolt.BodyID]heap.Handle)
	joltBodyPacked = make(map[uint32]heap.Handle)
	joltBodyDynamic = make(map[*jolt.BodyID]bool)
	joltBodyMu.Unlock()
	nextBufferIndex = 0
	matrixBufferAlloc = 1024
	matrixBuffer = make([]float32, matrixBufferAlloc*16)
	prevMatrixBuffer = make([]float32, matrixBufferAlloc*16)
	bufferIndexMap = make(map[*jolt.BodyID]int)
	bufferIndexToBody = make(map[int]*jolt.BodyID)
	m.accumulator = 0
	m.fixedStep = 1.0 / 60.0
	m.simTime = 0
	resetCollisionBridgeState()
	resetPickState()
	return value.Nil, nil
}

func (m *Module) phStop(args []value.Value) (value.Value, error) {
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("PHYSICS3D.STOP expects 0 arguments")
	}
	joltMu.Lock()
	defer joltMu.Unlock()
	if joltSys != nil {
		joltSys.Destroy()
		joltSys = nil
		joltBi = nil
	}
	if joltCoreInit {
		jolt.Shutdown()
		joltCoreInit = false
	}
	collRules = nil
	collPending = nil
	joltBodyMu.Lock()
	joltBodyToHandle = nil
	joltBodyPacked = nil
	joltBodyDynamic = nil
	joltBodyMu.Unlock()
	matrixBuffer = nil
	prevMatrixBuffer = nil
	bufferIndexMap = nil
	bufferIndexToBody = nil
	resetCollisionBridgeState()
	resetPickState()
	return value.Nil, nil
}

func (m *Module) phSetGravity(args []value.Value) (value.Value, error) {
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("PHYSICS3D.SETGRAVITY expects 3 float arguments")
	}
	x, _ := args[0].ToFloat()
	y, _ := args[1].ToFloat()
	z, _ := args[2].ToFloat()
	joltMu.Lock()
	gravX, gravY, gravZ = float32(x), float32(y), float32(z)
	joltMu.Unlock()
	return value.Nil, nil
}

func (m *Module) phGetGravityX(args []value.Value) (value.Value, error) {
	_ = m
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("PHYSICS3D.GETGRAVITYX expects 0 arguments")
	}
	g := GravityVec()
	return value.FromFloat(float64(g.X)), nil
}

func (m *Module) phGetGravityY(args []value.Value) (value.Value, error) {
	_ = m
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("PHYSICS3D.GETGRAVITYY expects 0 arguments")
	}
	g := GravityVec()
	return value.FromFloat(float64(g.Y)), nil
}

func (m *Module) phGetGravityZ(args []value.Value) (value.Value, error) {
	_ = m
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("PHYSICS3D.GETGRAVITYZ expects 0 arguments")
	}
	g := GravityVec()
	return value.FromFloat(float64(g.Z)), nil
}

func (m *Module) phSetTimeStep(args []value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("PHYSICS3D.SETTIMESTEP expects 1 Hz argument (e.g. 60)")
	}
	v, ok := args[0].ToFloat()
	if !ok || v < 1 {
		return value.Nil, fmt.Errorf("PHYSICS3D.SETTIMESTEP: invalid rate")
	}
	m.fixedStep = 1.0 / v
	return value.Nil, nil
}

func (m *Module) phSetSubsteps(args []value.Value) (value.Value, error) {
	_ = m
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("PHYSICS3D.SETSUBSTEPS expects 1 int argument")
	}
	// Reserved for future Jolt sub-step tuning; fixed 1/60 step today.
	return value.Nil, nil
}

func (m *Module) phStep(args []value.Value) (value.Value, error) {
	dt := 1.0 / 60.0
	if len(args) == 1 {
		if v, ok := args[0].ToFloat(); ok {
			dt = v
		}
	} else if len(args) > 1 {
		return value.Nil, fmt.Errorf("PHYSICS3D.STEP/UPDATE expects 0 or 1 argument (dt#)")
	}

	joltMu.Lock()
	ps := joltSys
	joltMu.Unlock()
	if ps == nil {
		return value.Nil, mbruntime.Errorf("PHYSICS3D.STEP/UPDATE: physics not started")
	}

	// Snapshot last published poses for render interpolation (translation blend).
	if len(matrixBuffer) > 0 && len(prevMatrixBuffer) == len(matrixBuffer) {
		copy(prevMatrixBuffer, matrixBuffer)
	} else if len(matrixBuffer) > 0 && (len(prevMatrixBuffer) != len(matrixBuffer)) {
		prevMatrixBuffer = make([]float32, len(matrixBuffer))
		copy(prevMatrixBuffer, matrixBuffer)
	}

	m.accumulator += dt
	steps := 0
	// Standard semi-fixed timestep accumulator with 5-step cap.
	for m.accumulator >= m.fixedStep && steps < 5 {
		m.applyConfiguredGravity(float32(m.fixedStep))
		ps.Update(float32(m.fixedStep))
		m.accumulator -= m.fixedStep
		m.simTime += m.fixedStep
		steps++
	}

	m.syncSharedBuffers()
	matrixInterpAccum = m.accumulator
	matrixInterpFixed = m.fixedStep
	if afterPhysicsMatrixSync != nil {
		afterPhysicsMatrixSync()
	}
	m.SyncWasmPhysicsAfterStep()
	collectContactsAfterStep(m)
	m.queueCollisionCallbacksFromRules()
	if physicsKCCFanIn != nil {
		physicsKCCFanIn(m)
	}

	// Process Aero (Shared Go Logic)
	m.ProcessAeroDynamics(float32(dt))

	return value.Nil, nil
}

// queueCollisionCallbacksFromRules converts contact-frame pairs into PHYSICS3D.ONCOLLISION callbacks.
// This runs after collectContactsAfterStep so PHYSICS3D.PROCESSCOLLISIONS sees a stable queue.
func (m *Module) queueCollisionCallbacksFromRules() {
	joltMu.Lock()
	rules := append([]collRule(nil), collRules...)
	joltMu.Unlock()
	if len(rules) == 0 {
		return
	}
	pending := make([]collEvent, 0, len(rules))
	for _, r := range rules {
		ea, oka := entityIDForCollisionRuleHandle(m, r.ha)
		eb, okb := entityIDForCollisionRuleHandle(m, r.hb)
		if !oka || !okb {
			continue
		}
		if _, hit := PairCollidedThisFrame(ea, eb); !hit {
			continue
		}
		pending = append(pending, collEvent{ha: r.ha, hb: r.hb, cb: r.cb})
	}
	if len(pending) == 0 {
		return
	}
	joltMu.Lock()
	collPending = append(collPending, pending...)
	joltMu.Unlock()
}

// applyConfiguredGravity applies PHYSICS3D.SETGRAVITY to dynamic BODY3D instances.
// The current jolt-go wrapper does not expose PhysicsSystem::SetGravity, so we
// integrate gravity by velocity for dynamic bodies each fixed simulation step.
func (m *Module) applyConfiguredGravity(dt float32) {
	if dt <= 0 {
		return
	}
	joltMu.Lock()
	bi := joltBi
	gx, gy, gz := gravX, gravY, gravZ
	joltMu.Unlock()
	if bi == nil {
		return
	}
	joltBodyMu.Lock()
	var gravIDs []*jolt.BodyID
	for id, dynamic := range joltBodyDynamic {
		if !dynamic {
			continue
		}
		if m.h != nil {
			if h, ok := joltLookupHandle(id); ok {
				if obj, ok := m.h.Get(h); ok {
					if bo, ok := obj.(*body3dObj); !ok || bo == nil || bo.motion != jolt.MotionTypeDynamic {
						continue
					}
				}
			}
		}
		gravIDs = append(gravIDs, id)
	}
	joltBodyMu.Unlock()

	if len(gravIDs) == 0 {
		return
	}
	dvx, dvy, dvz := gx*dt, gy*dt, gz*dt
	bi.BatchApplyGravityDelta(gravIDs, dvx, dvy, dvz)
}

func (m *Module) phSyncWasmToPhysRegs(args []value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("PHYSICS3D.SYNCWASMTOPHYSREGS expects (count, firstReg)")
	}
	ci, ok1 := args[0].ToInt()
	ri, ok2 := args[1].ToInt()
	if !ok1 || !ok2 {
		return value.Nil, fmt.Errorf("PHYSICS3D.SYNCWASMTOPHYSREGS: count and firstReg must be numeric")
	}
	if ci < 0 || ci > 256 || ri < 0 || ri > 255 {
		return value.Nil, fmt.Errorf("PHYSICS3D.SYNCWASMTOPHYSREGS: count must be 0..256, firstReg 0..255")
	}
	m.vmMu.Lock()
	v := m.vmRef
	m.vmMu.Unlock()
	if v == nil {
		return value.Nil, mbruntime.Errorf("PHYSICS3D.SYNCWASMTOPHYSREGS: VM not bound (engine wiring)")
	}
	m.wasmMu.Lock()
	view := m.wasmPhysicsView
	m.wasmMu.Unlock()
	if view.Mem == nil {
		return value.Nil, mbruntime.Errorf("PHYSICS3D.SYNCWASMTOPHYSREGS: WASM physics view not bound (call BindWasmPhysicsView from host)")
	}
	joltwasm.UpdateVMPhysics(v, view)
	if err := v.ExecSyncPhysics(uint8(ri), int(ci)); err != nil {
		return value.Nil, err
	}
	return value.Nil, nil
}

// matrix16FromPosQuatRL writes a column-major 4×4 world matrix (Raylib MODEL.DRAW layout) from
// Jolt position + quaternion. Matches three.js makeRotationFromQuaternion column packing.
func matrix16FromPosQuatRL(dest []float32, p jolt.Vec3, q jolt.Quat) {
	x, y, z, w := q.X, q.Y, q.Z, q.W
	x2, y2, z2 := x+x, y+y, z+z
	xx, xy, xz := x*x2, x*y2, x*z2
	yy, yz, zz := y*y2, y*z2, z*z2
	wx, wy, wz := w*x2, w*y2, w*z2

	dest[0] = 1 - (yy + zz)
	dest[1] = xy + wz
	dest[2] = xz - wy
	dest[3] = 0
	dest[4] = xy - wz
	dest[5] = 1 - (xx + zz)
	dest[6] = yz + wx
	dest[7] = 0
	dest[8] = xz + wy
	dest[9] = yz - wx
	dest[10] = 1 - (xx + yy)
	dest[11] = 0
	dest[12] = p.X
	dest[13] = p.Y
	dest[14] = p.Z
	dest[15] = 1
}

func (m *Module) syncSharedBuffers() {
	joltMu.Lock()
	bi := joltBi
	joltMu.Unlock()
	if bi == nil {
		return
	}

	joltBodyMu.Lock()
	n := nextBufferIndex
	if n <= 0 {
		joltBodyMu.Unlock()
		return
	}
	ids := make([]*jolt.BodyID, 0, n)
	idxs := make([]int, 0, n)
	for idx := 0; idx < n; idx++ {
		if id, ok := bufferIndexToBody[idx]; ok && id != nil {
			ids = append(ids, id)
			idxs = append(idxs, idx)
		}
	}
	joltBodyMu.Unlock()

	if len(ids) == 0 {
		return
	}
	batch := make([]float32, len(ids)*8)
	bi.BatchGetBodyTransforms(ids, batch)
	for j := range ids {
		off := j * 8
		pos := jolt.Vec3{X: batch[off], Y: batch[off+1], Z: batch[off+2]}
		rot := jolt.Quat{X: batch[off+3], Y: batch[off+4], Z: batch[off+5], W: batch[off+6]}
		idx := idxs[j]
		dest := matrixBuffer[idx*16 : (idx+1)*16]
		matrix16FromPosQuatRL(dest, pos, rot)
	}
}

func (m *Module) phGetMatrixBuffer(args []value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, mbruntime.Errorf("PHYSICS3D.GETMATRIXBUFFER: heap not bound")
	}
	// We return a handle to a Special shared numeric array
	arr, _ := heap.NewSharedArrayF32(matrixBuffer)
	id, err := m.h.Alloc(arr)
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}

func (m *Module) phProcessCollisions(args []value.Value) (value.Value, error) {
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("PHYSICS3D.PROCESSCOLLISIONS expects 0 arguments")
	}
	joltMu.Lock()
	sys := joltSys
	joltMu.Unlock()
	if sys == nil {
		return value.Nil, nil
	}
	
	// Drain rigid body contacts (sensors & triggers)
	events := sys.DrainContactQueue(256)
	for _, ev := range events {
		ha, okA := joltLookupHandle(ev.BodyA)
		hb, okB := joltLookupHandle(ev.BodyB)
		if !okA || !okB {
			continue
		}

		joltMu.Lock()
		for _, rule := range collRules {
			// Match bidirectional rules (h1=ha, h2=hb OR h1=hb, h2=ha)
			// A rule with ha=0 or hb=0 could be a wildcard, but MoonBASIC currently requires explicit pairs.
			if (rule.ha == ha && rule.hb == hb) || (rule.ha == hb && rule.hb == ha) {
				collPending = append(collPending, collEvent{
					ha: ha,
					hb: hb,
					cb: rule.cb,
				})
			}
		}
		joltMu.Unlock()
	}

	// Character collisions already populate collPending elsewhere (usually in Step)
	joltMu.Lock()
	q := collPending
	collPending = nil
	joltMu.Unlock()
	if m.invoke != nil {
		for _, ev := range q {
			_, _ = m.invoke(ev.cb, []value.Value{value.FromHandle(ev.ha), value.FromHandle(ev.hb)})
		}
	}
	return value.Nil, nil
}

func (m *Module) phRaycast(args []value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, mbruntime.Errorf("PHYSICS3D.RAYCAST: heap not bound")
	}
	if len(args) != 7 {
		return value.Nil, fmt.Errorf("PHYSICS3D.RAYCAST expects 7 floats (ox,oy,oz, dx,dy,dz, maxdist)")
	}
	joltMu.Lock()
	ps := joltSys
	joltMu.Unlock()
	if ps == nil {
		return value.Nil, mbruntime.Errorf("PHYSICS3D.RAYCAST: physics not started")
	}
	ox, _ := args[0].ToFloat()
	oy, _ := args[1].ToFloat()
	oz, _ := args[2].ToFloat()
	dx, _ := args[3].ToFloat()
	dy, _ := args[4].ToFloat()
	dz, _ := args[5].ToFloat()
	maxd, _ := args[6].ToFloat()
	origin := jolt.Vec3{X: float32(ox), Y: float32(oy), Z: float32(oz)}
	dir := jolt.Vec3{X: float32(dx), Y: float32(dy), Z: float32(dz)}
	L := dir.Length()
	if L > 1e-6 && float64(L) > maxd {
		s := float32(maxd / float64(L))
		dir = dir.Mul(s)
	}
	hit, ok := ps.CastRay(origin, dir)
	arr, err := heap.NewArray([]int64{6})
	if err != nil {
		return value.Nil, err
	}
	if !ok {
		_ = arr.Set([]int64{0}, 0)
		for i := int64(1); i < 6; i++ {
			_ = arr.Set([]int64{i}, 0)
		}
		id, err := m.h.Alloc(arr)
		if err != nil {
			return value.Nil, err
		}
		return value.FromHandle(id), nil
	}
	_ = arr.Set([]int64{0}, 1)
	_ = arr.Set([]int64{1}, float64(hit.Normal.X))
	_ = arr.Set([]int64{2}, float64(hit.Normal.Y))
	_ = arr.Set([]int64{3}, float64(hit.Normal.Z))
	_ = arr.Set([]int64{4}, float64(hit.Fraction))
	_ = arr.Set([]int64{5}, 0)
	id, err := m.h.Alloc(arr)
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}

func (m *Module) phCreateHingeJoint(args []value.Value) (value.Value, error) {
	if len(args) != 8 || args[0].Kind != value.KindHandle || args[1].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("JOINT.CREATEHINGE expects (b1, b2, px, py, pz, ax, ay, az)")
	}
	if _, err := heap.Cast[*body3dObj](m.h, heap.Handle(args[0].IVal)); err != nil {
		return value.Nil, err
	}
	if _, err := heap.Cast[*body3dObj](m.h, heap.Handle(args[1].IVal)); err != nil {
		return value.Nil, err
	}
	joltMu.Lock()
	ok := joltSys != nil
	joltMu.Unlock()
	if !ok {
		return value.Nil, fmt.Errorf("physics not started")
	}
	id, err := m.h.Alloc(&JointObj{})
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}

func (m *Module) phCreatePointJoint(args []value.Value) (value.Value, error) {
	if len(args) != 5 || args[0].Kind != value.KindHandle || args[1].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("JOINT.CREATEPOINT expects (b1, b2, px, py, pz)")
	}
	if _, err := heap.Cast[*body3dObj](m.h, heap.Handle(args[0].IVal)); err != nil {
		return value.Nil, err
	}
	if _, err := heap.Cast[*body3dObj](m.h, heap.Handle(args[1].IVal)); err != nil {
		return value.Nil, err
	}
	joltMu.Lock()
	ok := joltSys != nil
	joltMu.Unlock()
	if !ok {
		return value.Nil, fmt.Errorf("physics not started")
	}
	id, err := m.h.Alloc(&JointObj{})
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}

func (m *Module) phJointDelete(args []value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("JOINT.FREE expects (joint)")
	}
	j, err := heap.Cast[*JointObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	j.Free()
	m.h.Free(heap.Handle(args[0].IVal))
	return value.Nil, nil
}

func (m *Module) phJointFixed(args []value.Value) (value.Value, error) {
	_ = args
	return value.Nil, fmt.Errorf("JOINT.FIXED is not implemented on native backend yet")
}
func (m *Module) phJointHinge(args []value.Value) (value.Value, error) {
	if len(args) != 8 || args[0].Kind != value.KindHandle || args[1].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("JOINT.HINGE expects (b1, b2, px, py, pz, ax, ay, az)")
	}
	bo1, err := heap.Cast[*body3dObj](m.h, heap.Handle(args[0].IVal))
	if err != nil { return value.Nil, err }
	bo2, err := heap.Cast[*body3dObj](m.h, heap.Handle(args[1].IVal))
	if err != nil { return value.Nil, err }
	px, _ := args[2].ToFloat()
	py, _ := args[3].ToFloat()
	pz, _ := args[4].ToFloat()
	ax, _ := args[5].ToFloat()
	ay, _ := args[6].ToFloat()
	az, _ := args[7].ToFloat()

	joltMu.Lock()
	sys := joltSys
	joltMu.Unlock()
	if sys == nil { return value.Nil, nil }

	joint := sys.CreateHingeJoint(bo1.id, bo2.id, jolt.Vec3{X: float32(px), Y: float32(py), Z: float32(pz)}, jolt.Vec3{X: float32(ax), Y: float32(ay), Z: float32(az)})
	if joint == nil {
		return value.Nil, fmt.Errorf("failed to create hinge joint")
	}
	// We might want to return a joint handle here, but the manifest says returns null.
	return value.Nil, nil
}
func (m *Module) phJointSlider(args []value.Value) (value.Value, error) {
	_ = args
	return value.Nil, fmt.Errorf("JOINT.SLIDER is not implemented on native backend yet")
}
func (m *Module) phJointCone(args []value.Value) (value.Value, error) {
	_ = args
	return value.Nil, fmt.Errorf("JOINT.CONE is not implemented on native backend yet")
}

func (m *Module) bdSetGravityFactor(args []value.Value) (value.Value, error) {
	if len(args) != 2 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("BODY3D.SETGRAVITYFACTOR expects (body, factor#)")
	}
	bo, err := heap.Cast[*body3dObj](m.h, heap.Handle(args[0].IVal))
	if err != nil || bo.id == nil {
		return value.Nil, err
	}
	factor, _ := args[1].ToFloat()
	joltMu.Lock()
	bi := joltBi
	joltMu.Unlock()
	if bi != nil {
		bi.SetGravityFactor(bo.id, float32(factor))
	}
	bo.gravFactor = float32(factor)
	return args[0], nil
}

func (m *Module) bdSetDamping(args []value.Value) (value.Value, error) {
	if len(args) != 3 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("BODY3D.SETDAMPING expects (body, linear#, angular#)")
	}
	bo, err := heap.Cast[*body3dObj](m.h, heap.Handle(args[0].IVal))
	if err != nil || bo.id == nil {
		return value.Nil, err
	}
	lin, _ := args[1].ToFloat()
	ang, _ := args[2].ToFloat()

	joltMu.Lock()
	sys := joltSys
	joltMu.Unlock()
	if sys == nil {
		return value.Nil, nil
	}
	sys.SetBodyDamping(bo.id, float32(lin), float32(ang))
	bo.linDamp = float32(lin)
	bo.angDamp = float32(ang)
	return args[0], nil
}

func (m *Module) bdLockAxis(args []value.Value) (value.Value, error) {
	if len(args) != 2 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("BODY3D.LOCKAXIS expects (body, axis_flags)")
	}
	bo, err := heap.Cast[*body3dObj](m.h, heap.Handle(args[0].IVal))
	if err != nil || bo.id == nil {
		return value.Nil, err
	}
	flags, ok := args[1].ToInt()
	if !ok {
		return value.Nil, fmt.Errorf("invalid axis flags")
	}
	
	joltMu.Lock()
	sys := joltSys
	joltMu.Unlock()
	if sys != nil {
		sys.SetAllowedDOFs(bo.id, int(flags))
	}
	return args[0], nil
}

func (m *Module) btdSetCCD(args []value.Value) (value.Value, error) {
	if len(args) != 2 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("BODY3D.SETCCD expects (body, toggle)")
	}
	bo, err := heap.Cast[*body3dObj](m.h, heap.Handle(args[0].IVal))
	if err != nil || bo.id == nil {
		return value.Nil, err
	}
	quality := jolt.MotionQualityDiscrete
	if value.Truthy(args[1], nil, nil) {
		quality = jolt.MotionQualityLinearCast
	}
	joltMu.Lock()
	bi := joltBi
	joltMu.Unlock()
	if bi != nil {
		bi.SetMotionQuality(bo.id, quality)
	}
	bo.ccd = value.Truthy(args[1], nil, nil)
	return args[0], nil
}

func (m *Module) bdGetGravityFactor(args []value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("BODY3D.GETGRAVITYFACTOR expects handle")
	}
	bo, err := heap.Cast[*body3dObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.FromFloat(1), nil
	}
	return value.FromFloat(float64(bo.gravFactor)), nil
}

func (m *Module) bdGetDamping(args []value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("BODY3D.GETDAMPING expects handle")
	}
	bo, err := heap.Cast[*body3dObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.FromFloat(0), nil
	}
	// Return as array [lin, ang]
	return mbmatrix.AllocVec2Value(m.h, bo.linDamp, bo.angDamp)
}

func (m *Module) bdGetCCD(args []value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("BODY3D.GETCCD expects handle")
	}
	bo, err := heap.Cast[*body3dObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.FromBool(false), nil
	}
	return value.FromBool(bo.ccd), nil
}

func (m *Module) phDebugDraw(args []value.Value) (value.Value, error) {
	_ = args
	return value.Nil, fmt.Errorf("PHYSICS3D.DEBUGDRAW not implemented on native backend")
}
func (m *Module) phSpherecast(args []value.Value) (value.Value, error) {
	_ = args
	return value.Nil, fmt.Errorf("PHYSICS3D.SPHERECAST not implemented on native backend")
}
func (m *Module) phBoxcast(args []value.Value) (value.Value, error) {
	_ = args
	return value.Nil, fmt.Errorf("PHYSICS3D.BOXCAST not implemented on native backend")
}
func (m *Module) phEnable(args []value.Value) (value.Value, error) {
	_ = args
	return value.Nil, fmt.Errorf("PHYSICS3D.ENABLE not implemented on native backend")
}
func (m *Module) phDisable(args []value.Value) (value.Value, error) {
	_ = args
	return value.Nil, fmt.Errorf("PHYSICS3D.DISABLE not implemented on native backend")
}

func phSetOnCollision(m *Module, ha, hb value.Value, cb string) (value.Value, error) {
	joltMu.Lock()
	defer joltMu.Unlock()
	collRules = append(collRules, collRule{
		ha: heap.Handle(ha.IVal),
		hb: heap.Handle(hb.IVal),
		cb: cb,
	})
	return value.Nil, nil
}

func (m *Module) phWorldSetup(args []value.Value) (value.Value, error) {
	grav := -9.81
	if len(args) > 0 {
		if v, ok := args[0].ToFloat(); ok {
			grav = v
		}
	}
	joltMu.Lock()
	started := joltSys != nil
	joltMu.Unlock()
	if !started {
		if _, err := m.phStart(nil); err != nil {
			return value.Nil, err
		}
	}
	return m.phSetGravity([]value.Value{value.FromFloat(0), value.FromFloat(grav), value.FromFloat(0)})
}

// Internal Bridges for shared Go solvers (aero_host.go, vehicle_host.go)

func (m *Module) getBodyTransform(b *body3dObj) (rl.Vector3, rl.Quaternion, bool) {
	joltMu.Lock()
	defer joltMu.Unlock()
	if joltBi == nil {
		return rl.Vector3{}, rl.QuaternionIdentity(), false
	}
	p := joltBi.GetPosition(b.id)
	q := joltBi.GetRotation(b.id)
	return rl.Vector3{X: p.X, Y: p.Y, Z: p.Z}, rl.Quaternion{X: q.X, Y: q.Y, Z: q.Z, W: q.W}, true
}

func (m *Module) getBodyVelocity(b *body3dObj) rl.Vector3 {
	joltMu.Lock()
	defer joltMu.Unlock()
	if joltBi == nil {
		return rl.Vector3{}
	}
	p := joltBi.GetLinearVelocity(b.id)
	return rl.Vector3{X: p.X, Y: p.Y, Z: p.Z}
}

func (m *Module) applyBodyForce(b *body3dObj, f rl.Vector3) {
	dt := float32(m.fixedStep)
	if dt <= 0 {
		dt = 1.0 / 60.0
	}
	joltMu.Lock()
	defer joltMu.Unlock()
	if joltBi == nil {
		return
	}
	joltBi.AddImpulse(b.id, jolt.Vec3{X: f.X * dt, Y: f.Y * dt, Z: f.Z * dt})
}

// physicsKCCFanIn is registered by the player/KCC bridge to merge CharacterVirtual contact drains into ONCOLLISION.
var physicsKCCFanIn func(*Module)

// SetPhysicsKCCFanIn registers a hook invoked near the end of PHYSICS3D.STEP (after rigid contact rules are queued).
func SetPhysicsKCCFanIn(fn func(*Module)) {
	physicsKCCFanIn = fn
}

// FanInCharacterContactEvents appends matching PHYSICS3D.ONCOLLISION callbacks when KCC contacts involve playerEid and a registered rigid body.
func (m *Module) FanInCharacterContactEvents(playerEid int64, events []jolt.CharacterContactEvent) {
	if m == nil || m.h == nil || len(events) == 0 {
		return
	}
	joltMu.Lock()
	rules := append([]collRule(nil), collRules...)
	joltMu.Unlock()
	if len(rules) == 0 {
		return
	}
	seen := make(map[string]struct{}, len(events)*2)
	for _, ev := range events {
		bh, ok := LookupBodyHeapByPacked(ev.BodyB)
		if !ok {
			continue
		}
		otherEid, ok := EntityIDForBodyHandle(bh)
		if !ok {
			continue
		}
		for _, r := range rules {
			ea, oka := entityIDForCollisionRuleHandle(m, r.ha)
			eb, okb := entityIDForCollisionRuleHandle(m, r.hb)
			if !oka || !okb {
				continue
			}
			if !((ea == playerEid && eb == otherEid) || (eb == playerEid && ea == otherEid)) {
				continue
			}
			key := pairKey(ea, eb) + ":" + r.cb
			if _, dup := seen[key]; dup {
				continue
			}
			seen[key] = struct{}{}
			joltMu.Lock()
			collPending = append(collPending, collEvent{ha: r.ha, hb: r.hb, cb: r.cb})
			joltMu.Unlock()
		}
	}
}
