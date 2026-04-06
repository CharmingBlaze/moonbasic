package mbcollision

import (
	"fmt"
	"math"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

// registerRay2DBuiltins registers 2D ray vs primitive tests (pure math; works without CGO).
func (m *Module) registerRay2DBuiltins(reg runtime.Registrar) {
	reg.Register("RAY2D.HITCIRCLE_HIT", "collision", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		hit, _, _, _, err := ray2dHitCircleFromArgs(args)
		if err != nil {
			return value.Nil, err
		}
		return value.FromBool(hit), nil
	}))
	reg.Register("RAY2D.HITCIRCLE_DISTANCE", "collision", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		hit, dist, _, _, err := ray2dHitCircleFromArgs(args)
		if err != nil {
			return value.Nil, err
		}
		if !hit {
			return value.FromFloat(0.0), nil
		}
		return value.FromFloat(dist), nil
	}))
	reg.Register("RAY2D.HITCIRCLE_POINTX", "collision", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		hit, _, px, _, err := ray2dHitCircleFromArgs(args)
		if err != nil {
			return value.Nil, err
		}
		if !hit {
			return value.FromFloat(0.0), nil
		}
		return value.FromFloat(px), nil
	}))
	reg.Register("RAY2D.HITCIRCLE_POINTY", "collision", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		hit, _, _, py, err := ray2dHitCircleFromArgs(args)
		if err != nil {
			return value.Nil, err
		}
		if !hit {
			return value.FromFloat(0.0), nil
		}
		return value.FromFloat(py), nil
	}))

	reg.Register("RAY2D.HITRECT_HIT", "collision", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		hit, _, _, _, err := ray2dHitRectFromArgs(args)
		if err != nil {
			return value.Nil, err
		}
		return value.FromBool(hit), nil
	}))
	reg.Register("RAY2D.HITRECT_DISTANCE", "collision", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		hit, dist, _, _, err := ray2dHitRectFromArgs(args)
		if err != nil {
			return value.Nil, err
		}
		if !hit {
			return value.FromFloat(0.0), nil
		}
		return value.FromFloat(dist), nil
	}))
	reg.Register("RAY2D.HITRECT_POINTX", "collision", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		hit, _, px, _, err := ray2dHitRectFromArgs(args)
		if err != nil {
			return value.Nil, err
		}
		if !hit {
			return value.FromFloat(0.0), nil
		}
		return value.FromFloat(px), nil
	}))
	reg.Register("RAY2D.HITRECT_POINTY", "collision", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		hit, _, _, py, err := ray2dHitRectFromArgs(args)
		if err != nil {
			return value.Nil, err
		}
		if !hit {
			return value.FromFloat(0.0), nil
		}
		return value.FromFloat(py), nil
	}))

	reg.Register("RAY2D.HITSEGMENT_HIT", "collision", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		hit, _, _, _, err := ray2dHitSegmentFromArgs(args)
		if err != nil {
			return value.Nil, err
		}
		return value.FromBool(hit), nil
	}))
	reg.Register("RAY2D.HITSEGMENT_DISTANCE", "collision", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		hit, dist, _, _, err := ray2dHitSegmentFromArgs(args)
		if err != nil {
			return value.Nil, err
		}
		if !hit {
			return value.FromFloat(0.0), nil
		}
		return value.FromFloat(dist), nil
	}))
	reg.Register("RAY2D.HITSEGMENT_POINTX", "collision", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		hit, _, px, _, err := ray2dHitSegmentFromArgs(args)
		if err != nil {
			return value.Nil, err
		}
		if !hit {
			return value.FromFloat(0.0), nil
		}
		return value.FromFloat(px), nil
	}))
	reg.Register("RAY2D.HITSEGMENT_POINTY", "collision", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		hit, _, _, py, err := ray2dHitSegmentFromArgs(args)
		if err != nil {
			return value.Nil, err
		}
		if !hit {
			return value.FromFloat(0.0), nil
		}
		return value.FromFloat(py), nil
	}))
}

func argF64(v value.Value) (float64, bool) {
	if f, ok := v.ToFloat(); ok {
		return f, true
	}
	if i, ok := v.ToInt(); ok {
		return float64(i), true
	}
	return 0, false
}

// ray2dHitCircleFromArgs: (ox, oy, dx, dy, cx, cy, r) — direction (dx,dy) is normalized internally.
func ray2dHitCircleFromArgs(args []value.Value) (hit bool, t float64, px, py float64, err error) {
	if len(args) != 7 {
		return false, 0, 0, 0, fmt.Errorf("RAY2D.HITCIRCLE_* expects 7 arguments (ox, oy, dx, dy, cx, cy, r)")
	}
	fs := make([]float64, 7)
	for i := range fs {
		var ok bool
		fs[i], ok = argF64(args[i])
		if !ok {
			return false, 0, 0, 0, fmt.Errorf("RAY2D.HITCIRCLE_*: argument %d must be numeric", i+1)
		}
	}
	ox, oy, dx, dy := fs[0], fs[1], fs[2], fs[3]
	cx, cy, r := fs[4], fs[5], fs[6]
	if r < 0 {
		r = -r
	}
	lenD := math.Hypot(dx, dy)
	if lenD < 1e-12 {
		return false, 0, 0, 0, fmt.Errorf("RAY2D.HITCIRCLE_*: ray direction length must be > 0")
	}
	dx /= lenD
	dy /= lenD
	vx := ox - cx
	vy := oy - cy
	b := vx*dx + vy*dy
	c := vx*vx + vy*vy - r*r
	disc := b*b - c
	if disc < 0 {
		return false, 0, 0, 0, nil
	}
	sqrtD := math.Sqrt(disc)
	t0 := -b - sqrtD
	t1 := -b + sqrtD
	var tHit float64
	if t0 >= 0 {
		tHit = t0
	} else if t1 >= 0 {
		tHit = t1
	} else {
		return false, 0, 0, 0, nil
	}
	px = ox + tHit*dx
	py = oy + tHit*dy
	return true, tHit, px, py, nil
}

// ray2dHitRectFromArgs: axis-aligned box (minx, miny, maxx, maxy) vs ray (ox,oy,dx,dy).
func ray2dHitRectFromArgs(args []value.Value) (hit bool, t float64, px, py float64, err error) {
	if len(args) != 8 {
		return false, 0, 0, 0, fmt.Errorf("RAY2D.HITRECT_* expects 8 arguments (ox, oy, dx, dy, minx, miny, maxx, maxy)")
	}
	fs := make([]float64, 8)
	for i := range fs {
		var ok bool
		fs[i], ok = argF64(args[i])
		if !ok {
			return false, 0, 0, 0, fmt.Errorf("RAY2D.HITRECT_*: argument %d must be numeric", i+1)
		}
	}
	ox, oy, dx, dy := fs[0], fs[1], fs[2], fs[3]
	minx, miny, maxx, maxy := fs[4], fs[5], fs[6], fs[7]
	if minx > maxx {
		minx, maxx = maxx, minx
	}
	if miny > maxy {
		miny, maxy = maxy, miny
	}
	lenD := math.Hypot(dx, dy)
	if lenD < 1e-12 {
		return false, 0, 0, 0, fmt.Errorf("RAY2D.HITRECT_*: ray direction length must be > 0")
	}
	dx /= lenD
	dy /= lenD
	tMin := -1e100
	tMax := 1e100
	if math.Abs(dx) < 1e-12 {
		if ox < minx || ox > maxx {
			return false, 0, 0, 0, nil
		}
	} else {
		t1 := (minx - ox) / dx
		t2 := (maxx - ox) / dx
		if t1 > t2 {
			t1, t2 = t2, t1
		}
		tMin = math.Max(tMin, t1)
		tMax = math.Min(tMax, t2)
		if tMin > tMax {
			return false, 0, 0, 0, nil
		}
	}
	if math.Abs(dy) < 1e-12 {
		if oy < miny || oy > maxy {
			return false, 0, 0, 0, nil
		}
	} else {
		t1 := (miny - oy) / dy
		t2 := (maxy - oy) / dy
		if t1 > t2 {
			t1, t2 = t2, t1
		}
		tMin = math.Max(tMin, t1)
		tMax = math.Min(tMax, t2)
		if tMin > tMax {
			return false, 0, 0, 0, nil
		}
	}
	tHit := tMin
	if tHit < 0 {
		if tMax >= 0 {
			tHit = 0
		} else {
			return false, 0, 0, 0, nil
		}
	}
	px = ox + tHit*dx
	py = oy + tHit*dy
	return true, tHit, px, py, nil
}

// ray2dHitSegmentFromArgs: ray (ox,oy,dx,dy) vs segment (x1,y1)-(x2,y2).
func ray2dHitSegmentFromArgs(args []value.Value) (hit bool, t float64, px, py float64, err error) {
	if len(args) != 8 {
		return false, 0, 0, 0, fmt.Errorf("RAY2D.HITSEGMENT_* expects 8 arguments (ox, oy, dx, dy, x1, y1, x2, y2)")
	}
	fs := make([]float64, 8)
	for i := range fs {
		var ok bool
		fs[i], ok = argF64(args[i])
		if !ok {
			return false, 0, 0, 0, fmt.Errorf("RAY2D.HITSEGMENT_*: argument %d must be numeric", i+1)
		}
	}
	ox, oy, rdx, rdy := fs[0], fs[1], fs[2], fs[3]
	x1, y1, x2, y2 := fs[4], fs[5], fs[6], fs[7]
	lenR := math.Hypot(rdx, rdy)
	if lenR < 1e-12 {
		return false, 0, 0, 0, fmt.Errorf("RAY2D.HITSEGMENT_*: ray direction length must be > 0")
	}
	rdx /= lenR
	rdy /= lenR
	sdx := x2 - x1
	sdy := y2 - y1
	den := rdx*sdy - rdy*sdx
	if math.Abs(den) < 1e-12 {
		return false, 0, 0, 0, nil
	}
	tRay := ((x1-ox)*sdy - (y1-oy)*sdx) / den
	u := ((x1-ox)*rdy - (y1-oy)*rdx) / den
	if tRay < 0 || u < 0 || u > 1 {
		return false, 0, 0, 0, nil
	}
	px = ox + tRay*rdx
	py = oy + tRay*rdy
	return true, tRay, px, py, nil
}
