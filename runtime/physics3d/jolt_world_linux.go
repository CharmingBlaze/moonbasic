//go:build linux && cgo

package mbphysics3d

import (
	"fmt"
	"unsafe"

	"github.com/bbitechnologies/jolt-go/jolt"

	"moonbasic/internal/joltwasm"
	mbruntime "moonbasic/runtime"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

func phStart(m *Module, args []value.Value) (value.Value, error) {
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
	joltBodyToHandle = make(map[uintptr]heap.Handle)
	joltBodyMu.Unlock()
	nextBufferIndex = 0
	matrixBufferAlloc = 1024
	matrixBuffer = make([]float32, matrixBufferAlloc*16)
	bufferIndexMap = make(map[uintptr]int)
	bufferIndexToBody = make(map[int]uintptr)
	m.accumulator = 0
	m.fixedStep = 1.0 / 60.0
	resetCollisionBridgeState()
	resetPickState()
	return value.Nil, nil
}

func phStop(m *Module, args []value.Value) (value.Value, error) {
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

func phSetGravity(m *Module, args []value.Value) (value.Value, error) {
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

func phGetGravityX(m *Module, args []value.Value) (value.Value, error) {
	_ = m
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("PHYSICS3D.GETGRAVITYX expects 0 arguments")
	}
	g := GravityVec()
	return value.FromFloat(float64(g.X)), nil
}

func phGetGravityY(m *Module, args []value.Value) (value.Value, error) {
	_ = m
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("PHYSICS3D.GETGRAVITYY expects 0 arguments")
	}
	g := GravityVec()
	return value.FromFloat(float64(g.Y)), nil
}

func phGetGravityZ(m *Module, args []value.Value) (value.Value, error) {
	_ = m
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("PHYSICS3D.GETGRAVITYZ expects 0 arguments")
	}
	g := GravityVec()
	return value.FromFloat(float64(g.Z)), nil
}

func phSetTimeStep(m *Module, args []value.Value) (value.Value, error) {
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

func phSetSubsteps(m *Module, args []value.Value) (value.Value, error) {
	_ = m
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("PHYSICS3D.SETSUBSTEPS expects 1 int argument")
	}
	// Reserved for future Jolt sub-step tuning; fixed 1/60 step today.
	return value.Nil, nil
}

func phStep(m *Module, args []value.Value) (value.Value, error) {
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
	return value.Nil, nil
}

func phSyncWasmToPhysRegs(m *Module, args []value.Value) (value.Value, error) {
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
	for ptr := range joltBodyToHandle {
		idx, ok := bufferIndexMap[ptr]
		if !ok {
			continue
		}
		id := (*jolt.BodyID)(unsafe.Pointer(ptr))
		pos := bi.GetPosition(id)
		rot := bi.GetRotation(id)
		dest := matrixBuffer[idx*16 : (idx+1)*16]
		matrix16FromPosQuatRL(dest, pos, rot)
	}
}

func phGetMatrixBuffer(m *Module, args []value.Value) (value.Value, error) {
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

func phProcessCollisions(m *Module, args []value.Value) (value.Value, error) {
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

func phRaycast(m *Module, args []value.Value) (value.Value, error) {
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
