//go:build linux && cgo

package mbphysics3d

import (
	"fmt"
	"strings"
	"sync"

	"github.com/bbitechnologies/jolt-go/jolt"

	"moonbasic/runtime"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
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

type collRule struct {
	ha, hb heap.Handle
	cb     string
}

type collEvent struct {
	ha, hb heap.Handle
	cb     string
}

type body3dObj struct {
	id      *jolt.BodyID
	release heap.ReleaseOnce
}

func (b *body3dObj) TypeName() string { return "Body3D" }

func (b *body3dObj) TypeTag() uint16 { return heap.TagPhysicsBody }

func (b *body3dObj) Free() {
	b.release.Do(func() {
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

func registerPhysics3DCommands(m *Module, reg runtime.Registrar) {
	reg.Register("PHYSICS3D.START", "physics3d", runtime.AdaptLegacy(func(a []value.Value) (value.Value, error) { return phStart(m, a) }))
	reg.Register("PHYSICS3D.STOP", "physics3d", runtime.AdaptLegacy(func(a []value.Value) (value.Value, error) { return phStop(m, a) }))
	reg.Register("PHYSICS3D.SETGRAVITY", "physics3d", runtime.AdaptLegacy(func(a []value.Value) (value.Value, error) { return phSetGravity(m, a) }))
	reg.Register("PHYSICS3D.STEP", "physics3d", runtime.AdaptLegacy(func(a []value.Value) (value.Value, error) { return phStep(m, a) }))
	reg.Register("PHYSICS3D.SETSUBSTEPS", "physics3d", runtime.AdaptLegacy(func(a []value.Value) (value.Value, error) { return phSetSubsteps(m, a) }))
	reg.Register("PHYSICS3D.ONCOLLISION", "physics3d", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 3 || args[0].Kind != value.KindHandle || args[1].Kind != value.KindHandle || args[2].Kind != value.KindString {
			return value.Nil, fmt.Errorf("PHYSICS3D.ONCOLLISION expects (handle, handle, string)")
		}
		cb, err := rt.ArgString(args, 2)
		if err != nil {
			return value.Nil, err
		}
		joltMu.Lock()
		collRules = append(collRules, collRule{ha: heap.Handle(args[0].IVal), hb: heap.Handle(args[1].IVal), cb: cb})
		joltMu.Unlock()
		return value.Nil, nil
	})
	reg.Register("PHYSICS3D.PROCESSCOLLISIONS", "physics3d", runtime.AdaptLegacy(func(a []value.Value) (value.Value, error) { return phProcessCollisions(m, a) }))
	reg.Register("PHYSICS3D.RAYCAST", "physics3d", runtime.AdaptLegacy(func(a []value.Value) (value.Value, error) { return phRaycast(m, a) }))
	reg.Register("BODY3D.MAKE", "physics3d", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if m.h == nil {
			return value.Nil, runtime.Errorf("BODY3D.MAKE: heap not bound")
		}
		motion := "dynamic"
		if len(args) == 0 {
			// default motion type
		} else if len(args) == 1 && args[0].Kind == value.KindString {
			var err error
			motion, err = rt.ArgString(args, 0)
			if err != nil {
				return value.Nil, err
			}
		} else {
			return value.Nil, fmt.Errorf("BODY3D.MAKE expects 0 arguments (default DYNAMIC) or 1 motion string (STATIC, KINEMATIC, DYNAMIC)")
		}
		b := &builderObj{motion: parseMotion(motion)}
		bid, err := m.h.Alloc(b)
		if err != nil {
			return value.Nil, err
		}
		return value.FromHandle(bid), nil
	})
	reg.Register("BODY3D.ADDBOX", "physics3d", runtime.AdaptLegacy(func(a []value.Value) (value.Value, error) { return bdAddBox(m, a) }))
	reg.Register("BODY3D.ADDSPHERE", "physics3d", runtime.AdaptLegacy(func(a []value.Value) (value.Value, error) { return bdAddSphere(m, a) }))
	reg.Register("BODY3D.ADDCAPSULE", "physics3d", runtime.AdaptLegacy(func(a []value.Value) (value.Value, error) { return bdAddCapsule(m, a) }))
	reg.Register("BODY3D.ADDMESH", "physics3d", runtime.AdaptLegacy(func(a []value.Value) (value.Value, error) { return bdAddMesh(m, a) }))
	reg.Register("BODY3D.COMMIT", "physics3d", runtime.AdaptLegacy(func(a []value.Value) (value.Value, error) { return bdCommit(m, a) }))
	reg.Register("BODY3D.SETPOS", "physics3d", runtime.AdaptLegacy(func(a []value.Value) (value.Value, error) { return bdSetPos(m, a) }))
	reg.Register("BODY3D.SETPOSITION", "physics3d", runtime.AdaptLegacy(func(a []value.Value) (value.Value, error) { return bdSetPos(m, a) }))
	reg.Register("BODY3D.GETPOS", "physics3d", runtime.AdaptLegacy(func(a []value.Value) (value.Value, error) { return bdGetPos(m, a) }))
	reg.Register("BODY3D.ACTIVATE", "physics3d", runtime.AdaptLegacy(func(a []value.Value) (value.Value, error) { return bdActivate(m, a) }))
	reg.Register("BODY3D.DEACTIVATE", "physics3d", runtime.AdaptLegacy(func(a []value.Value) (value.Value, error) { return bdDeactivate(m, a) }))
	reg.Register("BODY3D.SETROT", "physics3d", runtime.AdaptLegacy(func(a []value.Value) (value.Value, error) { return bdNoOp(m, a) }))
	reg.Register("BODY3D.GETROT", "physics3d", runtime.AdaptLegacy(func(a []value.Value) (value.Value, error) { return bdGetRotZero(m, a) }))
	reg.Register("BODY3D.SETMASS", "physics3d", runtime.AdaptLegacy(func(a []value.Value) (value.Value, error) { return bdNoOp(m, a) }))
	reg.Register("BODY3D.SETFRICTION", "physics3d", runtime.AdaptLegacy(func(a []value.Value) (value.Value, error) { return bdNoOp(m, a) }))
	reg.Register("BODY3D.SETRESTITUTION", "physics3d", runtime.AdaptLegacy(func(a []value.Value) (value.Value, error) { return bdNoOp(m, a) }))
	reg.Register("BODY3D.APPLYFORCE", "physics3d", runtime.AdaptLegacy(func(a []value.Value) (value.Value, error) { return bdNoOp(m, a) }))
	reg.Register("BODY3D.APPLYIMPULSE", "physics3d", runtime.AdaptLegacy(func(a []value.Value) (value.Value, error) { return bdNoOp(m, a) }))
	reg.Register("BODY3D.SETLINEARVEL", "physics3d", runtime.AdaptLegacy(func(a []value.Value) (value.Value, error) { return bdNoOp(m, a) }))
	reg.Register("BODY3D.SETANGULARVEL", "physics3d", runtime.AdaptLegacy(func(a []value.Value) (value.Value, error) { return bdNoOp(m, a) }))
	reg.Register("BODY3D.X", "physics3d", runtime.AdaptLegacy(func(a []value.Value) (value.Value, error) { return bdAxis(m, a, 0) }))
	reg.Register("BODY3D.Y", "physics3d", runtime.AdaptLegacy(func(a []value.Value) (value.Value, error) { return bdAxis(m, a, 1) }))
	reg.Register("BODY3D.Z", "physics3d", runtime.AdaptLegacy(func(a []value.Value) (value.Value, error) { return bdAxis(m, a, 2) }))
	reg.Register("BODY3D.FREE", "physics3d", runtime.AdaptLegacy(func(a []value.Value) (value.Value, error) { return bdFree(m, a) }))
}

func shutdownPhysics3D(m *Module) {
	_, _ = phStop(m, nil)
}

func phStart(m *Module, args []value.Value) (value.Value, error) {
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("PHYSICS3D.START expects 0 arguments")
	}
	joltMu.Lock()
	defer joltMu.Unlock()
	if joltSys != nil {
		return value.Nil, runtime.Errorf("PHYSICS3D.START: already started")
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

func phSetSubsteps(m *Module, args []value.Value) (value.Value, error) {
	_ = m
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("PHYSICS3D.SETSUBSTEPS expects 1 int argument")
	}
	// Reserved for future Jolt sub-step tuning; fixed 1/60 step today.
	return value.Nil, nil
}

func phStep(m *Module, args []value.Value) (value.Value, error) {
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("PHYSICS3D.STEP expects 0 arguments")
	}
	joltMu.Lock()
	ps := joltSys
	joltMu.Unlock()
	if ps == nil {
		return value.Nil, runtime.Errorf("PHYSICS3D.STEP: physics not started")
	}
	ps.Update(1.0 / 60.0)
	return value.Nil, nil
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
		return value.Nil, runtime.Errorf("PHYSICS3D.RAYCAST: heap not bound")
	}
	if len(args) != 7 {
		return value.Nil, fmt.Errorf("PHYSICS3D.RAYCAST expects 7 floats (ox,oy,oz, dx,dy,dz, maxdist)")
	}
	joltMu.Lock()
	ps := joltSys
	joltMu.Unlock()
	if ps == nil {
		return value.Nil, runtime.Errorf("PHYSICS3D.RAYCAST: physics not started")
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

func bdAddBox(m *Module, args []value.Value) (value.Value, error) {
	if len(args) != 4 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("BODY3D.ADDBOX expects (builder, hw, hh, hd)")
	}
	bu, err := heap.Cast[*builderObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	hx, _ := args[1].ToFloat()
	hy, _ := args[2].ToFloat()
	hz, _ := args[3].ToFloat()
	if bu.shape != nil {
		bu.shape.Destroy()
	}
	bu.shape = jolt.CreateBox(jolt.Vec3{X: float32(hx), Y: float32(hy), Z: float32(hz)})
	return value.Nil, nil
}

func bdAddSphere(m *Module, args []value.Value) (value.Value, error) {
	if len(args) != 2 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("BODY3D.ADDSPHERE expects (builder, radius)")
	}
	bu, err := heap.Cast[*builderObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	r, _ := args[1].ToFloat()
	if bu.shape != nil {
		bu.shape.Destroy()
	}
	bu.shape = jolt.CreateSphere(float32(r))
	return value.Nil, nil
}

func bdAddCapsule(m *Module, args []value.Value) (value.Value, error) {
	if len(args) != 3 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("BODY3D.ADDCAPSULE expects (builder, radius, height)")
	}
	bu, err := heap.Cast[*builderObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	r, _ := args[1].ToFloat()
	h, _ := args[2].ToFloat()
	hh := float32(h)/2 - float32(r)
	if hh < 0.05 {
		hh = 0.05
	}
	if bu.shape != nil {
		bu.shape.Destroy()
	}
	bu.shape = jolt.CreateCapsule(hh, float32(r))
	return value.Nil, nil
}

func bdAddMesh(m *Module, args []value.Value) (value.Value, error) {
	return value.Nil, fmt.Errorf("BODY3D.ADDMESH: requires Phase D model handle (not implemented)")
}

func bdCommit(m *Module, args []value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("BODY3D.COMMIT: heap not bound")
	}
	if len(args) != 4 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("BODY3D.COMMIT expects (builder, x, y, z)")
	}
	joltMu.Lock()
	bi := joltBi
	joltMu.Unlock()
	if bi == nil {
		return value.Nil, runtime.Errorf("BODY3D.COMMIT: physics not started")
	}
	bu, err := heap.Cast[*builderObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	if bu.shape == nil {
		return value.Nil, fmt.Errorf("BODY3D.COMMIT: no shape (call ADDBOX/ADDSPHERE/ADDCAPSULE first)")
	}
	x, _ := args[1].ToFloat()
	y, _ := args[2].ToFloat()
	z, _ := args[3].ToFloat()
	sh := bu.shape
	bu.shape = nil
	id := bi.CreateBody(sh, jolt.Vec3{X: float32(x), Y: float32(y), Z: float32(z)}, bu.motion, false)
	sh.Destroy()
	m.h.Free(heap.Handle(args[0].IVal))
	body := &body3dObj{id: id}
	bh, err := m.h.Alloc(body)
	if err != nil {
		if id != nil {
			id.Destroy()
		}
		return value.Nil, err
	}
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
		return value.Nil, runtime.Errorf("BODY3D.SETPOS: physics not started")
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
		return value.Nil, runtime.Errorf("BODY3D.ACTIVATE: physics not started")
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
		return value.Nil, runtime.Errorf("BODY3D.DEACTIVATE: physics not started")
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
		return value.Nil, runtime.Errorf("BODY3D.GETPOS: heap not bound")
	}
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("BODY3D.GETPOS expects body handle")
	}
	joltMu.Lock()
	bi := joltBi
	joltMu.Unlock()
	if bi == nil {
		return value.Nil, runtime.Errorf("BODY3D.GETPOS: physics not started")
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

func bdGetRotZero(m *Module, args []value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("BODY3D.GETROT: heap not bound")
	}
	arr, err := heap.NewArray([]int64{3})
	if err != nil {
		return value.Nil, err
	}
	rh, err := m.h.Alloc(arr)
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(rh), nil
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

func bdFree(m *Module, args []value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("BODY3D.FREE expects handle")
	}
	m.h.Free(heap.Handle(args[0].IVal))
	return value.Nil, nil
}
