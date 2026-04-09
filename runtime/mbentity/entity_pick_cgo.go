//go:build cgo || (windows && !cgo)

package mbentity

import (
	"fmt"
	"sync"

	mbcamera "moonbasic/runtime/camera"
	"moonbasic/runtime"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var pickMu sync.Mutex
var lastPick pickRec

type pickRec struct {
	valid      bool
	hitX       float32
	hitY       float32
	hitZ       float32
	nX, nY, nZ float32
	entID      int64
	dist       float32
	tri, surf  int32
}

func registerPickBlitz(m *Module, r runtime.Registrar) {
	r.Register("LinePick", "entity", runtime.AdaptLegacy(m.linePick))
	r.Register("CameraPick", "entity", m.cameraPick)
	r.Register("PickedX", "entity", runtime.AdaptLegacy(m.pickedX))
	r.Register("PickedY", "entity", runtime.AdaptLegacy(m.pickedY))
	r.Register("PickedZ", "entity", runtime.AdaptLegacy(m.pickedZ))
	r.Register("PickedNX", "entity", runtime.AdaptLegacy(m.pickedNX))
	r.Register("PickedNY", "entity", runtime.AdaptLegacy(m.pickedNY))
	r.Register("PickedNZ", "entity", runtime.AdaptLegacy(m.pickedNZ))
	r.Register("PickedEntity", "entity", runtime.AdaptLegacy(m.pickedEntity))
	r.Register("PickedDistance", "entity", runtime.AdaptLegacy(m.pickedDistance))
	r.Register("PickedSurface", "entity", runtime.AdaptLegacy(m.pickedSurface))
	r.Register("PickedTriangle", "entity", runtime.AdaptLegacy(m.pickedTriangle))
}

// LinePick(ox, oy, oz, dx, dy, dz [, radius]) — ray vs static entity AABBs; radius reserved (swept test not implemented).
func (m *Module) linePick(args []value.Value) (value.Value, error) {
	if len(args) != 6 && len(args) != 7 {
		return value.Nil, fmt.Errorf("LinePick expects (x#, y#, z#, dx#, dy#, dz# [, radius#])")
	}
	ox, ok1 := argF32(args[0])
	oy, ok2 := argF32(args[1])
	oz, ok3 := argF32(args[2])
	dx, ok4 := argF32(args[3])
	dy, ok5 := argF32(args[4])
	dz, ok6 := argF32(args[5])
	if !ok1 || !ok2 || !ok3 || !ok4 || !ok5 || !ok6 {
		return value.Nil, fmt.Errorf("LinePick: numeric arguments required")
	}
	_ = len(args)
	origin := rl.Vector3{X: ox, Y: oy, Z: oz}
	dir := rl.Vector3{X: dx, Y: dy, Z: dz}
	if rl.Vector3Length(dir) < 1e-8 {
		m.clearPick()
		return value.FromInt(0), nil
	}
	dir = rl.Vector3Normalize(dir)
	end := rl.Vector3Add(origin, rl.Vector3Scale(dir, 1e6))
	bestID := int64(0)
	bestT := float32(1e30)
	var bestHit, bestN rl.Vector3
	st := m.store()
	for _, e := range st.ents {
		if e == nil || e.hidden || !e.static {
			continue
		}
		switch e.kind {
		case entKindBox, entKindPlane, entKindCylinder, entKindCone, entKindMesh, entKindModel:
			mn, mx := m.aabbWorldMinMax(e)
			t := rayAABB(origin, end, mn, mx)
			if t >= 0 && t < bestT {
				bestT = t
				bestID = e.id
				bestHit = rl.Vector3Add(origin, rl.Vector3Scale(dir, t))
				bestN = aabbHitNormal(mn, mx, bestHit)
			}
		}
	}
	if bestID == 0 {
		m.clearPick()
		return value.FromInt(0), nil
	}
	m.storePick(bestID, bestHit, bestN, bestT)
	return value.FromInt(bestID), nil
}

func aabbHitNormal(mn, mx, hit rl.Vector3) rl.Vector3 {
	c := rl.Vector3{X: (mn.X + mx.X) * 0.5, Y: (mn.Y + mx.Y) * 0.5, Z: (mn.Z + mx.Z) * 0.5}
	dd := rl.Vector3Subtract(hit, c)
	ax := absF32(dd.X)
	ay := absF32(dd.Y)
	az := absF32(dd.Z)
	switch {
	case ax >= ay && ax >= az:
		if dd.X >= 0 {
			return rl.Vector3{X: 1, Y: 0, Z: 0}
		}
		return rl.Vector3{X: -1, Y: 0, Z: 0}
	case ay >= ax && ay >= az:
		if dd.Y >= 0 {
			return rl.Vector3{X: 0, Y: 1, Z: 0}
		}
		return rl.Vector3{X: 0, Y: -1, Z: 0}
	default:
		if dd.Z >= 0 {
			return rl.Vector3{X: 0, Y: 0, Z: 1}
		}
		return rl.Vector3{X: 0, Y: 0, Z: -1}
	}
}

func absF32(f float32) float32 {
	if f < 0 {
		return -f
	}
	return f
}

func (m *Module) clearPick() {
	pickMu.Lock()
	lastPick = pickRec{}
	pickMu.Unlock()
}

func (m *Module) storePick(id int64, hit rl.Vector3, n rl.Vector3, t float32) {
	pickMu.Lock()
	defer pickMu.Unlock()
	if id < 1 {
		lastPick = pickRec{}
		return
	}
	lastPick = pickRec{
		valid: true,
		hitX:  hit.X,
		hitY:  hit.Y,
		hitZ:  hit.Z,
		nX:    n.X,
		nY:    n.Y,
		nZ:    n.Z,
		entID: id,
		dist:  t,
		tri:   -1,
		surf:  -1,
	}
}

// CameraPick(camera, screenX, screenY) — screen-space pick using active render size.
func (m *Module) cameraPick(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("CameraPick: heap not bound")
	}
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("CameraPick expects (camera, screenX#, screenY#)")
	}
	if args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("CameraPick: camera handle required")
	}
	ch := heap.Handle(args[0].IVal)
	sx, ok1 := argF32(args[1])
	sy, ok2 := argF32(args[2])
	if !ok1 || !ok2 {
		return value.Nil, fmt.Errorf("CameraPick: screen coords must be numeric")
	}
	cam, err := mbcamera.RayCamera3D(m.h, ch)
	if err != nil {
		return value.Nil, err
	}
	rw := float32(rl.GetRenderWidth())
	rh := float32(rl.GetRenderHeight())
	if rw <= 0 || rh <= 0 {
		return value.Nil, fmt.Errorf("CameraPick: render size not ready")
	}
	ray := rl.GetScreenToWorldRayEx(rl.Vector2{X: sx, Y: sy}, cam, int32(rw), int32(rh))
	origin := ray.Position
	dir := rl.Vector3Normalize(ray.Direction)
	end := rl.Vector3Add(origin, rl.Vector3Scale(dir, 1e6))
	bestID := int64(0)
	bestT := float32(1e30)
	var bestHit, bestN rl.Vector3
	st := m.store()
	for _, e := range st.ents {
		if e == nil || e.hidden || !e.static {
			continue
		}
		switch e.kind {
		case entKindBox, entKindPlane, entKindCylinder, entKindCone, entKindMesh, entKindModel:
			mn, mx := m.aabbWorldMinMax(e)
			t := rayAABB(origin, end, mn, mx)
			if t >= 0 && t < bestT {
				bestT = t
				bestID = e.id
				bestHit = rl.Vector3Add(origin, rl.Vector3Scale(dir, t))
				bestN = aabbHitNormal(mn, mx, bestHit)
			}
		}
	}
	_ = rt
	if bestID < 1 {
		m.clearPick()
		return value.FromInt(0), nil
	}
	m.storePick(bestID, bestHit, bestN, bestT)
	return value.FromInt(bestID), nil
}

func (m *Module) pickedX(args []value.Value) (value.Value, error) {
	pickMu.Lock()
	p := lastPick
	pickMu.Unlock()
	if !p.valid {
		return value.FromFloat(0), nil
	}
	return value.FromFloat(float64(p.hitX)), nil
}
func (m *Module) pickedY(args []value.Value) (value.Value, error) {
	pickMu.Lock()
	p := lastPick
	pickMu.Unlock()
	if !p.valid {
		return value.FromFloat(0), nil
	}
	return value.FromFloat(float64(p.hitY)), nil
}
func (m *Module) pickedZ(args []value.Value) (value.Value, error) {
	pickMu.Lock()
	p := lastPick
	pickMu.Unlock()
	if !p.valid {
		return value.FromFloat(0), nil
	}
	return value.FromFloat(float64(p.hitZ)), nil
}
func (m *Module) pickedNX(args []value.Value) (value.Value, error) {
	pickMu.Lock()
	p := lastPick
	pickMu.Unlock()
	if !p.valid {
		return value.FromFloat(0), nil
	}
	return value.FromFloat(float64(p.nX)), nil
}
func (m *Module) pickedNY(args []value.Value) (value.Value, error) {
	pickMu.Lock()
	p := lastPick
	pickMu.Unlock()
	if !p.valid {
		return value.FromFloat(0), nil
	}
	return value.FromFloat(float64(p.nY)), nil
}
func (m *Module) pickedNZ(args []value.Value) (value.Value, error) {
	pickMu.Lock()
	p := lastPick
	pickMu.Unlock()
	if !p.valid {
		return value.FromFloat(0), nil
	}
	return value.FromFloat(float64(p.nZ)), nil
}
func (m *Module) pickedEntity(args []value.Value) (value.Value, error) {
	pickMu.Lock()
	p := lastPick
	pickMu.Unlock()
	if !p.valid {
		return value.FromInt(0), nil
	}
	return value.FromInt(p.entID), nil
}
func (m *Module) pickedDistance(args []value.Value) (value.Value, error) {
	pickMu.Lock()
	p := lastPick
	pickMu.Unlock()
	if !p.valid {
		return value.FromFloat(0), nil
	}
	return value.FromFloat(float64(p.dist)), nil
}
func (m *Module) pickedSurface(args []value.Value) (value.Value, error) {
	pickMu.Lock()
	p := lastPick
	pickMu.Unlock()
	if !p.valid {
		return value.FromInt(-1), nil
	}
	return value.FromInt(int64(p.surf)), nil
}
func (m *Module) pickedTriangle(args []value.Value) (value.Value, error) {
	pickMu.Lock()
	p := lastPick
	pickMu.Unlock()
	if !p.valid {
		return value.FromInt(-1), nil
	}
	return value.FromInt(int64(p.tri)), nil
}
