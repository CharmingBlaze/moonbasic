//go:build linux && cgo

package mbphysics3d

import (
	"fmt"
	"math"
	"sync"

	"github.com/bbitechnologies/jolt-go/jolt"

	mbruntime "moonbasic/runtime"
	mbcamera "moonbasic/runtime/camera"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

// Pick layer lookup: entity id -> collision layer 0–31. Registered from mbentity.
var (
	pickMu           sync.Mutex
	pickLayerLookup  func(int64) (uint8, bool)
	pickOx, pickOy, pickOz float64
	pickDx, pickDy, pickDz float64
	pickMaxDist      float64 // if > 0, normalize direction and scale to this length
	pickLayerMask    uint32  // 0 = accept all layers; else bit i = accept layer i
	pickRadius       float64

	pickHit       bool
	pickEnt       int64
	pickPX, pickPY, pickPZ float64
	pickNX, pickNY, pickNZ float64
	pickDist      float64
)

// SetPickLayerLookup registers a resolver for ENTITY.COLLISIONLAYER (call from mbentity).
func SetPickLayerLookup(fn func(int64) (uint8, bool)) {
	pickMu.Lock()
	pickLayerLookup = fn
	pickMu.Unlock()
}

func resetPickState() {
	pickMu.Lock()
	pickOx, pickOy, pickOz = 0, 0, 0
	pickDx, pickDy, pickDz = 0, 0, -1
	pickMaxDist = 0
	pickLayerMask = 0
	pickRadius = 0
	clearPickResultLocked()
	pickMu.Unlock()
}

func clearPickResultLocked() {
	pickHit = false
	pickEnt = 0
	pickPX, pickPY, pickPZ = 0, 0, 0
	pickNX, pickNY, pickNZ = 0, 0, 1
	pickDist = 0
}

func registerPickCommands(m *Module, reg mbruntime.Registrar) {
	reg.Register("PICK.ORIGIN", "physics3d", mbruntime.AdaptLegacy(func(a []value.Value) (value.Value, error) { return pickOrigin(m, a) }))
	reg.Register("PICK.DIRECTION", "physics3d", mbruntime.AdaptLegacy(func(a []value.Value) (value.Value, error) { return pickDirection(m, a) }))
	reg.Register("PICK.MAXDIST", "physics3d", mbruntime.AdaptLegacy(func(a []value.Value) (value.Value, error) { return pickMaxDistSet(m, a) }))
	reg.Register("PICK.LAYERMASK", "physics3d", mbruntime.AdaptLegacy(func(a []value.Value) (value.Value, error) { return pickLayerMaskSet(m, a) }))
	reg.Register("PICK.RADIUS", "physics3d", mbruntime.AdaptLegacy(func(a []value.Value) (value.Value, error) { return pickRadiusSet(m, a) }))
	reg.Register("PICK.CAST", "physics3d", mbruntime.AdaptLegacy(func(a []value.Value) (value.Value, error) { return pickCast(m, a) }))
	reg.Register("PICK.FROMCAMERA", "physics3d", mbruntime.AdaptLegacy(func(a []value.Value) (value.Value, error) { return pickFromCamera(m, a) }))
	reg.Register("PICK.SCREENCAST", "physics3d", mbruntime.AdaptLegacy(func(a []value.Value) (value.Value, error) { return pickScreenCast(m, a) }))
	reg.Register("PICK.X", "physics3d", mbruntime.AdaptLegacy(func(a []value.Value) (value.Value, error) { return pickGet(m, a, func() float64 { return pickPX }) }))
	reg.Register("PICK.Y", "physics3d", mbruntime.AdaptLegacy(func(a []value.Value) (value.Value, error) { return pickGet(m, a, func() float64 { return pickPY }) }))
	reg.Register("PICK.Z", "physics3d", mbruntime.AdaptLegacy(func(a []value.Value) (value.Value, error) { return pickGet(m, a, func() float64 { return pickPZ }) }))
	reg.Register("PICK.NX", "physics3d", mbruntime.AdaptLegacy(func(a []value.Value) (value.Value, error) { return pickGet(m, a, func() float64 { return pickNX }) }))
	reg.Register("PICK.NY", "physics3d", mbruntime.AdaptLegacy(func(a []value.Value) (value.Value, error) { return pickGet(m, a, func() float64 { return pickNY }) }))
	reg.Register("PICK.NZ", "physics3d", mbruntime.AdaptLegacy(func(a []value.Value) (value.Value, error) { return pickGet(m, a, func() float64 { return pickNZ }) }))
	reg.Register("PICK.ENTITY", "physics3d", mbruntime.AdaptLegacy(func(a []value.Value) (value.Value, error) { return pickEntityGet(m, a) }))
	reg.Register("PICK.DIST", "physics3d", mbruntime.AdaptLegacy(func(a []value.Value) (value.Value, error) { return pickGet(m, a, func() float64 { return pickDist }) }))
	reg.Register("PICK.HIT", "physics3d", mbruntime.AdaptLegacy(func(a []value.Value) (value.Value, error) { return pickHitGet(m, a) }))
}

func pickOrigin(m *Module, args []value.Value) (value.Value, error) {
	_ = m
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("PICK.ORIGIN expects (x#, y#, z#)")
	}
	x, _ := args[0].ToFloat()
	y, _ := args[1].ToFloat()
	z, _ := args[2].ToFloat()
	pickMu.Lock()
	pickOx, pickOy, pickOz = x, y, z
	pickMu.Unlock()
	return value.Nil, nil
}

func pickDirection(m *Module, args []value.Value) (value.Value, error) {
	_ = m
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("PICK.DIRECTION expects (dx#, dy#, dz#)")
	}
	x, _ := args[0].ToFloat()
	y, _ := args[1].ToFloat()
	z, _ := args[2].ToFloat()
	pickMu.Lock()
	pickDx, pickDy, pickDz = x, y, z
	pickMu.Unlock()
	return value.Nil, nil
}

func pickMaxDistSet(m *Module, args []value.Value) (value.Value, error) {
	_ = m
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("PICK.MAXDIST expects (maxDist#)")
	}
	d, _ := args[0].ToFloat()
	pickMu.Lock()
	pickMaxDist = d
	pickMu.Unlock()
	return value.Nil, nil
}

func pickLayerMaskSet(m *Module, args []value.Value) (value.Value, error) {
	_ = m
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("PICK.LAYERMASK expects (mask#)")
	}
	v, _ := args[0].ToFloat()
	pickMu.Lock()
	pickLayerMask = uint32(int64(v)) //nolint:gosec // BASIC mask fits uint32
	pickMu.Unlock()
	return value.Nil, nil
}

func pickRadiusSet(m *Module, args []value.Value) (value.Value, error) {
	_ = m
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("PICK.RADIUS expects (radius#)")
	}
	r, _ := args[0].ToFloat()
	if r > 1e-9 {
		return value.Nil, fmt.Errorf("PICK.RADIUS: non-zero sphere pick not implemented")
	}
	pickMu.Lock()
	pickRadius = r
	pickMu.Unlock()
	return value.Nil, nil
}

func pickFromCamera(m *Module, args []value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, mbruntime.Errorf("PICK.FROMCAMERA: heap not bound")
	}
	if len(args) != 3 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("PICK.FROMCAMERA expects (camera, screenX#, screenY#)")
	}
	sx, _ := args[1].ToFloat()
	sy, _ := args[2].ToFloat()
	pos, dir, err := mbcamera.WorldRayFromScreen(m.h, heap.Handle(args[0].IVal), float32(sx), float32(sy))
	if err != nil {
		return value.Nil, err
	}
	pickMu.Lock()
	pickOx, pickOy, pickOz = float64(pos.X), float64(pos.Y), float64(pos.Z)
	pickDx, pickDy, pickDz = float64(dir.X), float64(dir.Y), float64(dir.Z)
	// Screen rays use normalized direction; require MAXDIST for length unless user set already
	if pickMaxDist <= 0 {
		pickMaxDist = 10_000
	}
	pickMu.Unlock()
	return value.Nil, nil
}

func pickScreenCast(m *Module, args []value.Value) (value.Value, error) {
	if _, err := pickFromCamera(m, args); err != nil {
		return value.Nil, err
	}
	return pickCast(m, nil)
}

func pickCast(m *Module, args []value.Value) (value.Value, error) {
	if args != nil && len(args) != 0 {
		return value.Nil, fmt.Errorf("PICK.CAST expects 0 arguments")
	}
	if pickRadius > 1e-9 {
		return value.Nil, fmt.Errorf("PICK.CAST: set PICK.RADIUS to 0 (sphere pick not implemented)")
	}
	joltMu.Lock()
	ps := joltSys
	joltMu.Unlock()
	if ps == nil {
		return value.Nil, mbruntime.Errorf("PICK.CAST: physics not started")
	}

	pickMu.Lock()
	ox, oy, oz := pickOx, pickOy, pickOz
	dx, dy, dz := pickDx, pickDy, pickDz
	maxD := pickMaxDist
	mask := pickLayerMask
	pickMu.Unlock()

	lensq := dx*dx + dy*dy + dz*dz
	if lensq < 1e-30 {
		pickMu.Lock()
		clearPickResultLocked()
		pickMu.Unlock()
		return value.FromInt(0), nil
	}

	var jdir jolt.Vec3
	if maxD > 0 {
		inv := 1.0 / math.Sqrt(lensq)
		jdir = jolt.Vec3{
			X: float32(dx * inv * maxD),
			Y: float32(dy * inv * maxD),
			Z: float32(dz * inv * maxD),
		}
	} else {
		jdir = jolt.Vec3{X: float32(dx), Y: float32(dy), Z: float32(dz)}
	}

	origin := jolt.Vec3{X: float32(ox), Y: float32(oy), Z: float32(oz)}
	hits := ps.CastRayGetHits(origin, jdir, 64)

	pickMu.Lock()
	defer pickMu.Unlock()
	clearPickResultLocked()

	rayLen := math.Sqrt(float64(jdir.X)*float64(jdir.X) + float64(jdir.Y)*float64(jdir.Y) + float64(jdir.Z)*float64(jdir.Z))

	for _, hit := range hits {
		if hit.BodyID == nil {
			continue
		}
		bh, ok := joltLookupHandle(hit.BodyID)
		if !ok {
			continue
		}
		eid, ok := EntityIDForBodyHandle(bh)
		if !ok || eid < 1 {
			continue
		}
		if !pickLayerAllowedLocked(eid, mask) {
			continue
		}
		pickHit = true
		pickEnt = eid
		pickPX = float64(hit.HitPoint.X)
		pickPY = float64(hit.HitPoint.Y)
		pickPZ = float64(hit.HitPoint.Z)
		pickNX = float64(hit.Normal.X)
		pickNY = float64(hit.Normal.Y)
		pickNZ = float64(hit.Normal.Z)
		pickDist = float64(hit.Fraction) * rayLen
		return value.FromInt(eid), nil
	}

	return value.FromInt(0), nil
}

func pickLayerAllowedLocked(entID int64, mask uint32) bool {
	if mask == 0 {
		return true
	}
	if pickLayerLookup == nil {
		return true
	}
	ly, ok := pickLayerLookup(entID)
	if !ok {
		return (mask & 1) != 0 // layer 0 default
	}
	if ly >= 32 {
		return false
	}
	return (mask & (1 << ly)) != 0
}

func pickGet(m *Module, args []value.Value, f func() float64) (value.Value, error) {
	_ = m
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("PICK getter expects 0 arguments")
	}
	pickMu.Lock()
	v := f()
	pickMu.Unlock()
	return value.FromFloat(v), nil
}

func pickEntityGet(m *Module, args []value.Value) (value.Value, error) {
	_ = m
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("PICK.ENTITY expects 0 arguments")
	}
	pickMu.Lock()
	e := pickEnt
	pickMu.Unlock()
	return value.FromInt(e), nil
}

func pickHitGet(m *Module, args []value.Value) (value.Value, error) {
	_ = m
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("PICK.HIT expects 0 arguments")
	}
	pickMu.Lock()
	h := pickHit
	pickMu.Unlock()
	return value.FromBool(h), nil
}
