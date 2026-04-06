//go:build cgo || (windows && !cgo)

package mbcamera

import (
	"fmt"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/runtime"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

var cullStats struct {
	total, culled, visible            int
	frustumCulled, distanceCulled     int
	horizonCulled, occlusionCulled    int
}

// Phase B: occluder handles (reserved; no hot-path work in Phase A).
var occlusionEnabled bool
var occluderHandles []heap.Handle

func (m *Module) registerCull(r runtime.Registrar) {
	r.Register("CULL.SPHEREVISIBLE", "camera", m.cullSphereVisible)
	r.Register("CULL.AABBVISIBLE", "camera", m.cullAABBVisible)
	r.Register("CULL.POINTVISIBLE", "camera", m.cullPointVisible)
	r.Register("CULL.INRANGE", "camera", m.cullInRange)
	r.Register("CULL.DISTANCE", "camera", m.cullDistance)
	r.Register("CULL.DISTANCESQ", "camera", m.cullDistanceSq)
	r.Register("CULL.BEHINDHORIZON", "camera", m.cullBehindHorizon)
	r.Register("CULL.BATCHSPHERE", "camera", m.cullBatchSphere)
	r.Register("CULL.OCCLUSIONENABLE", "camera", m.cullOcclusionEnable)
	r.Register("CULL.OCCLUDERADD", "camera", m.cullOccluderAdd)
	r.Register("CULL.OCCLUDERCLEAR", "camera", m.cullOccluderClear)
	r.Register("CULL.ISOCCLUDED", "camera", m.cullIsOccluded)
	r.Register("CULL.SETMAXDISTANCE", "camera", m.cullSetMaxDistance)
	r.Register("CULL.GETMAXDISTANCE", "camera", m.cullGetMaxDistance)
	r.Register("CULL.STATSRESET", "camera", m.cullStatsReset)
	r.Register("CULL.STATSTOTAL", "camera", m.cullStatsTotal)
	r.Register("CULL.STATSCULLED", "camera", m.cullStatsCulled)
	r.Register("CULL.STATSVISIBLE", "camera", m.cullStatsVisible)
	r.Register("CULL.STATSFRUSTUMCULLED", "camera", m.cullStatsFrustumCulled)
	r.Register("CULL.STATSDISTANCECULLED", "camera", m.cullStatsDistanceCulled)
	r.Register("CULL.STATSHORIZONCULLED", "camera", m.cullStatsHorizonCulled)
	r.Register("CULL.STATSOCCLUSIONCULLED", "camera", m.cullStatsOcclusionCulled)
	r.Register("CULL.SETBACKFACECULLING", "camera", m.cullSetBackfaceCulling)
}

func argF32(v value.Value) float32 {
	if f, ok := v.ToFloat(); ok {
		return float32(f)
	}
	if i, ok := v.ToInt(); ok {
		return float32(i)
	}
	return 0
}

func (m *Module) cullSphereVisible(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 4 {
		return value.Nil, fmt.Errorf("CULL.SPHEREVISIBLE expects (cx, cy, cz, r)")
	}
	if !frustumValid {
		return rt.RetBool(true), nil
	}
	cx := argF32(args[0])
	cy := argF32(args[1])
	cz := argF32(args[2])
	rad := argF32(args[3])

	cullStats.total++
	if !WithinDistance(cx, cy, cz, activeCamPos[0], activeCamPos[1], activeCamPos[2], globalMaxDist) {
		cullStats.distanceCulled++
		cullStats.culled++
		return rt.RetBool(false), nil
	}
	if !activeFrustum.SphereVisible(cx, cy, cz, rad) {
		cullStats.frustumCulled++
		cullStats.culled++
		return rt.RetBool(false), nil
	}
	cullStats.visible++
	return rt.RetBool(true), nil
}

func (m *Module) cullAABBVisible(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 6 {
		return value.Nil, fmt.Errorf("CULL.AABBVISIBLE expects (minX,minY,minZ,maxX,maxY,maxZ)")
	}
	if !frustumValid {
		return rt.RetBool(true), nil
	}
	minX := argF32(args[0])
	minY := argF32(args[1])
	minZ := argF32(args[2])
	maxX := argF32(args[3])
	maxY := argF32(args[4])
	maxZ := argF32(args[5])

	cx := (minX + maxX) * 0.5
	cy := (minY + maxY) * 0.5
	cz := (minZ + maxZ) * 0.5

	cullStats.total++
	if !WithinDistance(cx, cy, cz, activeCamPos[0], activeCamPos[1], activeCamPos[2], globalMaxDist) {
		cullStats.distanceCulled++
		cullStats.culled++
		return rt.RetBool(false), nil
	}
	if !activeFrustum.AABBVisible(minX, minY, minZ, maxX, maxY, maxZ) {
		cullStats.frustumCulled++
		cullStats.culled++
		return rt.RetBool(false), nil
	}
	cullStats.visible++
	return rt.RetBool(true), nil
}

func (m *Module) cullPointVisible(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("CULL.POINTVISIBLE expects (x, y, z)")
	}
	if !frustumValid {
		return rt.RetBool(true), nil
	}
	x := argF32(args[0])
	y := argF32(args[1])
	z := argF32(args[2])

	cullStats.total++
	if !WithinDistance(x, y, z, activeCamPos[0], activeCamPos[1], activeCamPos[2], globalMaxDist) {
		cullStats.distanceCulled++
		cullStats.culled++
		return rt.RetBool(false), nil
	}
	if !activeFrustum.PointVisible(x, y, z) {
		cullStats.frustumCulled++
		cullStats.culled++
		return rt.RetBool(false), nil
	}
	cullStats.visible++
	return rt.RetBool(true), nil
}

func (m *Module) cullInRange(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	switch len(args) {
	case 3:
		if !frustumValid {
			return rt.RetBool(true), nil
		}
		cx := argF32(args[0])
		cy := argF32(args[1])
		cz := argF32(args[2])
		ok := WithinDistance(cx, cy, cz, activeCamPos[0], activeCamPos[1], activeCamPos[2], globalMaxDist)
		return rt.RetBool(ok), nil
	case 4:
		if !frustumValid {
			return rt.RetBool(true), nil
		}
		cx := argF32(args[0])
		cy := argF32(args[1])
		cz := argF32(args[2])
		maxd := argF32(args[3])
		ok := WithinDistance(cx, cy, cz, activeCamPos[0], activeCamPos[1], activeCamPos[2], maxd)
		return rt.RetBool(ok), nil
	default:
		return value.Nil, fmt.Errorf("CULL.INRANGE expects (cx,cy,cz) or (cx,cy,cz,maxdist)")
	}
}

func (m *Module) cullDistance(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("CULL.DISTANCE expects (cx, cy, cz)")
	}
	if !frustumValid {
		return rt.RetFloat(0), nil
	}
	cx := argF32(args[0])
	cy := argF32(args[1])
	cz := argF32(args[2])
	dx := cx - activeCamPos[0]
	dy := cy - activeCamPos[1]
	dz := cz - activeCamPos[2]
	d := math.Sqrt(float64(dx*dx + dy*dy + dz*dz))
	return rt.RetFloat(d), nil
}

func (m *Module) cullDistanceSq(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("CULL.DISTANCESQ expects (cx, cy, cz)")
	}
	if !frustumValid {
		return rt.RetFloat(0), nil
	}
	cx := argF32(args[0])
	cy := argF32(args[1])
	cz := argF32(args[2])
	dx := cx - activeCamPos[0]
	dy := cy - activeCamPos[1]
	dz := cz - activeCamPos[2]
	return rt.RetFloat(float64(dx*dx + dy*dy + dz*dz)), nil
}

func (m *Module) cullBehindHorizon(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 4 {
		return value.Nil, fmt.Errorf("CULL.BEHINDHORIZON expects (camera, maxY, cx, cz)")
	}
	h, err := rt.ArgHandle(args, 0)
	if err != nil {
		return value.Nil, err
	}
	o, err := heap.Cast[*camObj](m.h, heap.Handle(h))
	if err != nil {
		return value.Nil, err
	}
	maxY := argF32(args[1])
	cx := argF32(args[2])
	cz := argF32(args[3])

	fwd := rl.Vector3Subtract(o.cam.Target, o.cam.Position)
	fwd = rl.Vector3Normalize(fwd)
	pitch := float32(math.Asin(float64(fwd.Y))) * (180.0 / math.Pi)

	behind := BehindHorizon(
		o.cam.Position.X, o.cam.Position.Y, o.cam.Position.Z,
		maxY, cx, cz,
		pitch, o.cam.Fovy,
	)

	cullStats.total++
	if behind {
		cullStats.horizonCulled++
		cullStats.culled++
	} else {
		cullStats.visible++
	}
	return rt.RetBool(behind), nil
}

func (m *Module) cullBatchSphere(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("CULL.BATCHSPHERE expects (positions, radii, results)")
	}
	ph, err := rt.ArgHandle(args, 0)
	if err != nil {
		return value.Nil, err
	}
	rh, err := rt.ArgHandle(args, 1)
	if err != nil {
		return value.Nil, err
	}
	resh, err := rt.ArgHandle(args, 2)
	if err != nil {
		return value.Nil, err
	}
	posArr, err := heap.Cast[*heap.Array](m.h, heap.Handle(ph))
	if err != nil {
		return value.Nil, err
	}
	radArr, err := heap.Cast[*heap.Array](m.h, heap.Handle(rh))
	if err != nil {
		return value.Nil, err
	}
	outArr, err := heap.Cast[*heap.Array](m.h, heap.Handle(resh))
	if err != nil {
		return value.Nil, err
	}
	if radArr.Kind != heap.ArrayKindFloat {
		return value.Nil, fmt.Errorf("CULL.BATCHSPHERE: radii must be numeric array")
	}
	if posArr.Kind != heap.ArrayKindFloat {
		return value.Nil, fmt.Errorf("CULL.BATCHSPHERE: positions must be numeric array")
	}
	if outArr.Kind != heap.ArrayKindBool {
		return value.Nil, fmt.Errorf("CULL.BATCHSPHERE: results must be bool array")
	}

	n := radArr.TotalElements()
	if n <= 0 || posArr.TotalElements() < n*3 || outArr.TotalElements() < n {
		return value.Nil, nil
	}

	camX, camY, camZ := activeCamPos[0], activeCamPos[1], activeCamPos[2]
	maxDistSq := globalMaxDist * globalMaxDist

	for i := 0; i < n; i++ {
		cx, _ := posArr.GetFloat([]int64{int64(i * 3)})
		cy, _ := posArr.GetFloat([]int64{int64(i*3 + 1)})
		cz, _ := posArr.GetFloat([]int64{int64(i*3 + 2)})
		rf, _ := radArr.GetFloat([]int64{int64(i)})
		cxf := float32(cx)
		cyf := float32(cy)
		czf := float32(cz)
		rf32 := float32(rf)

		var visible bool
		if !frustumValid {
			visible = true
		} else {
			dx := cxf - camX
			dy := cyf - camY
			dz := czf - camZ
			distSq := dx*dx + dy*dy + dz*dz
			visible = distSq <= maxDistSq && activeFrustum.SphereVisible(cxf, cyf, czf, rf32)
		}

		_ = outArr.SetFloat([]int64{int64(i)}, boolToFloat(visible))

		cullStats.total++
		if visible {
			cullStats.visible++
		} else {
			cullStats.culled++
			if frustumValid {
				dx := cxf - camX
				dy := cyf - camY
				dz := czf - camZ
				if dx*dx+dy*dy+dz*dz > maxDistSq {
					cullStats.distanceCulled++
				} else {
					cullStats.frustumCulled++
				}
			}
		}
	}
	return value.Nil, nil
}

func boolToFloat(b bool) float64 {
	if b {
		return 1
	}
	return 0
}

func (m *Module) cullOcclusionEnable(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("CULL.OCCLUSIONENABLE expects (enable)")
	}
	b, err := rt.ArgBool(args, 0)
	if err != nil {
		return value.Nil, err
	}
	occlusionEnabled = b
	return rt.RetBool(true), nil
}

func (m *Module) cullOccluderAdd(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("CULL.OCCLUDERADD expects (model)")
	}
	h, err := rt.ArgHandle(args, 0)
	if err != nil {
		return value.Nil, err
	}
	occluderHandles = append(occluderHandles, heap.Handle(h))
	_ = occlusionEnabled
	return value.Nil, nil
}

func (m *Module) cullOccluderClear(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("CULL.OCCLUDERCLEAR expects no arguments")
	}
	occluderHandles = occluderHandles[:0]
	return value.Nil, nil
}

func (m *Module) cullIsOccluded(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 4 {
		return value.Nil, fmt.Errorf("CULL.ISOCCLUDED expects (cx, cy, cz, r)")
	}
	return rt.RetBool(false), nil
}

func (m *Module) cullSetMaxDistance(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("CULL.SETMAXDISTANCE expects (maxdist)")
	}
	d := argF32(args[0])
	if d < 0 {
		d = 0
	}
	SetGlobalMaxDistance(d)
	return value.Nil, nil
}

func (m *Module) cullGetMaxDistance(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("CULL.GETMAXDISTANCE expects no arguments")
	}
	return rt.RetFloat(float64(GlobalMaxDistance())), nil
}

func (m *Module) cullStatsReset(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("CULL.STATSRESET expects no arguments")
	}
	cullStats.total = 0
	cullStats.culled = 0
	cullStats.visible = 0
	cullStats.frustumCulled = 0
	cullStats.distanceCulled = 0
	cullStats.horizonCulled = 0
	cullStats.occlusionCulled = 0
	return value.Nil, nil
}

func (m *Module) cullStatsTotal(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("CULL.STATSTOTAL expects no arguments")
	}
	return rt.RetInt(int64(cullStats.total)), nil
}

func (m *Module) cullStatsCulled(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("CULL.STATSCULLED expects no arguments")
	}
	return rt.RetInt(int64(cullStats.culled)), nil
}

func (m *Module) cullStatsVisible(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("CULL.STATSVISIBLE expects no arguments")
	}
	return rt.RetInt(int64(cullStats.visible)), nil
}

func (m *Module) cullStatsFrustumCulled(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("CULL.STATSFRUSTUMCULLED expects no arguments")
	}
	return rt.RetInt(int64(cullStats.frustumCulled)), nil
}

func (m *Module) cullStatsDistanceCulled(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("CULL.STATSDISTANCECULLED expects no arguments")
	}
	return rt.RetInt(int64(cullStats.distanceCulled)), nil
}

func (m *Module) cullStatsHorizonCulled(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("CULL.STATSHORIZONCULLED expects no arguments")
	}
	return rt.RetInt(int64(cullStats.horizonCulled)), nil
}

func (m *Module) cullStatsOcclusionCulled(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("CULL.STATSOCCLUSIONCULLED expects no arguments")
	}
	return rt.RetInt(int64(cullStats.occlusionCulled)), nil
}

func (m *Module) cullSetBackfaceCulling(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("CULL.SETBACKFACECULLING expects (enable)")
	}
	b, err := rt.ArgBool(args, 0)
	if err != nil {
		return value.Nil, err
	}
	if b {
		rl.EnableBackfaceCulling()
	} else {
		rl.DisableBackfaceCulling()
	}
	return value.Nil, nil
}
