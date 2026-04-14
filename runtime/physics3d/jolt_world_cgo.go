//go:build (linux || windows) && cgo

package mbphysics3d

import (
	"fmt"

	"github.com/bbitechnologies/jolt-go/jolt"
	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/internal/joltwasm"
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
	joltBodyMu.Unlock()
	nextBufferIndex = 0
	matrixBufferAlloc = 1024
	matrixBuffer = make([]float32, matrixBufferAlloc*16)
	bufferIndexMap = make(map[*jolt.BodyID]int)
	bufferIndexToBody = make(map[int]*jolt.BodyID)
	m.accumulator = 0
	m.fixedStep = 1.0 / 60.0
	resetCollisionBridgeState()
	resetPickState()
	return value.Nil, nil
}

func (m *Module) phStop(args []value.Value) (value.Value, error) {
	if args != nil && len(args) != 0 {
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
	joltBodyMu.Unlock()
	matrixBuffer = nil
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
		return value.Nil, fmt.Errorf("PHYSICS3D.STEP expects 0 or 1 argument (dt#)")
	}

	joltMu.Lock()
	ps := joltSys
	joltMu.Unlock()
	if ps == nil {
		return value.Nil, mbruntime.Errorf("PHYSICS3D.STEP: physics not started")
	}

	m.accumulator += dt
	steps := 0
	// Standard semi-fixed timestep accumulator with 5-step cap.
	for m.accumulator >= m.fixedStep && steps < 5 {
		ps.Update(float32(m.fixedStep))
		m.accumulator -= m.fixedStep
		steps++
	}

	m.syncSharedBuffers()
	if afterPhysicsMatrixSync != nil {
		afterPhysicsMatrixSync()
	}
	m.SyncWasmPhysicsAfterStep()
	collectContactsAfterStep(m)

	// Process Aero (Shared Go Logic)
	m.ProcessAeroDynamics(float32(dt))

	return value.Nil, nil
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
	defer joltBodyMu.Unlock()

	// Sync every registered body into the shared matrix buffer.
	for id := range joltBodyToHandle {
		idx, ok := bufferIndexMap[id]
		if !ok {
			continue
		}
		pos := bi.GetPosition(id)
		rot := bi.GetRotation(id)
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
	if m.invoke == nil {
		return value.Nil, nil
	}
	joltMu.Lock()
	q := collPending
	collPending = nil
	joltMu.Unlock()
	for _, ev := range q {
		_, _ = m.invoke(ev.cb, []value.Value{value.FromHandle(ev.ha), value.FromHandle(ev.hb)})
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

func (m *Module) phJointFixed(args []value.Value) (value.Value, error) { return value.Nil, nil }
func (m *Module) phJointHinge(args []value.Value) (value.Value, error) { return value.Nil, nil }
func (m *Module) phJointSlider(args []value.Value) (value.Value, error) { return value.Nil, nil }
func (m *Module) phJointCone(args []value.Value) (value.Value, error) { return value.Nil, nil }

func (m *Module) bdSetGravityFactor(args []value.Value) (value.Value, error) {
	if len(args) != 2 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("BODY3D.SETGRAVITYFACTOR expects (body, factor#)")
	}
	if _, err := heap.Cast[*body3dObj](m.h, heap.Handle(args[0].IVal)); err != nil {
		return value.Nil, err
	}
	_, _ = args[1].ToFloat()
	// Not exposed by vendored jolt-go BodyInterface (see docs/PHYSICS.md, vendored Jolt table).
	return value.Nil, nil
}

func (m *Module) bdSetDamping(args []value.Value) (value.Value, error) {
	if len(args) != 3 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("BODY3D.SETDAMPING expects (body, linear#, angular#)")
	}
	if _, err := heap.Cast[*body3dObj](m.h, heap.Handle(args[0].IVal)); err != nil {
		return value.Nil, err
	}
	_, _ = args[1].ToFloat()
	_, _ = args[2].ToFloat()
	return value.Nil, nil
}

func (m *Module) bdLockAxis(args []value.Value) (value.Value, error) {
	if len(args) != 2 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("BODY3D.LOCKAXIS expects (body, axis_flags)")
	}
	if _, err := heap.Cast[*body3dObj](m.h, heap.Handle(args[0].IVal)); err != nil {
		return value.Nil, err
	}
	if _, ok := args[1].ToInt(); !ok {
		return value.Nil, fmt.Errorf("invalid axis flags")
	}
	return value.Nil, nil
}

func (m *Module) btdSetCCD(args []value.Value) (value.Value, error) {
	if len(args) != 2 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("BODY3D.SETCCD expects (body, toggle)")
	}
	if _, err := heap.Cast[*body3dObj](m.h, heap.Handle(args[0].IVal)); err != nil {
		return value.Nil, err
	}
	_ = value.Truthy(args[1], nil, nil)
	// Motion quality / CCD not exposed on vendored jolt-go.
	return value.Nil, nil
}

func (m *Module) phDebugDraw(args []value.Value) (value.Value, error) { return value.Nil, nil }
func (m *Module) phSpherecast(args []value.Value) (value.Value, error) { return value.Nil, nil }
func (m *Module) phBoxcast(args []value.Value) (value.Value, error) { return value.Nil, nil }
func (m *Module) phEnable(args []value.Value) (value.Value, error) { return value.Nil, nil }
func (m *Module) phDisable(args []value.Value) (value.Value, error) { return value.Nil, nil }

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

func (m *Module) bdGetLinearVel(args []value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("BODY3D.GETLINEARVEL expects (body)")
	}
	b, err := heap.Cast[*body3dObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	joltMu.Lock()
	var v jolt.Vec3
	if joltBi != nil {
		v = joltBi.GetLinearVelocity(b.id)
	}
	joltMu.Unlock()
	return valueVec3FromFloats(m.h, float64(v.X), float64(v.Y), float64(v.Z))
}


func (m *Module) bdGetAngularVel(args []value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("BODY3D.GETANGULARVEL expects (body)")
	}
	if _, err := heap.Cast[*body3dObj](m.h, heap.Handle(args[0].IVal)); err != nil {
		return value.Nil, err
	}
	// Angular velocity getters are not in vendored jolt-go; return a zero vector for API stability.
	return valueVec3FromFloats(m.h, 0, 0, 0)
}

func (m *Module) bdSetAngularVel(args []value.Value) (value.Value, error) {
	if len(args) != 4 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("BODY3D.SETANGULARVEL expects (body, x#, y#, z#)")
	}
	if _, err := heap.Cast[*body3dObj](m.h, heap.Handle(args[0].IVal)); err != nil {
		return value.Nil, err
	}
	_, _ = args[1].ToFloat()
	_, _ = args[2].ToFloat()
	_, _ = args[3].ToFloat()
	return value.Nil, nil
}


func (m *Module) bdApplyTorque(args []value.Value) (value.Value, error) {
	if len(args) != 4 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("BODY3D.APPLYTORQUE expects (body, x#, y#, z#)")
	}
	if _, err := heap.Cast[*body3dObj](m.h, heap.Handle(args[0].IVal)); err != nil {
		return value.Nil, err
	}
	_, _ = args[1].ToFloat()
	_, _ = args[2].ToFloat()
	_, _ = args[3].ToFloat()
	return value.Nil, nil
}

func (m *Module) bdGetMass(args []value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("BODY3D.GETMASS expects (body)")
	}
	if _, err := heap.Cast[*body3dObj](m.h, heap.Handle(args[0].IVal)); err != nil {
		return value.Nil, err
	}
	return value.FromFloat(1.0), nil
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
