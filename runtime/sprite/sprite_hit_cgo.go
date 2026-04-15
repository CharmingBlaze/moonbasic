//go:build cgo || (windows && !cgo)

package mbsprite

import (
	"fmt"
	"math"

	"moonbasic/runtime"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

type hitVec2 struct {
	x, y float64
}

// spriteHitGeometry matches DrawTexturePro / raylib rtextures.c: destination rect, scaled origin,
// and rotation in degrees (same as drawSpriteAtScreen).
func spriteHitGeometry(s *spriteObj) (destX, destY, dw, dh, ox, oy, rotDeg float64) {
	destX = float64(s.x)
	destY = float64(s.y)
	sx := float64(s.scaleX)
	sy := float64(s.scaleY)
	if sx == 0 {
		sx = 1
	}
	if sy == 0 {
		sy = 1
	}
	dw = float64(s.frameW) * sx
	dh = float64(s.frameH) * sy
	ox = float64(s.originX) * sx
	oy = float64(s.originY) * sy
	rotDeg = float64(s.rotRad) * 180.0 / math.Pi
	return
}

// spriteQuadCorners mirrors raylib DrawTexturePro corner math (rtextures.c).
func spriteQuadCorners(destX, destY, dw, dh, ox, oy, rotDeg float64) (tl, bl, br, tr hitVec2) {
	if math.Abs(rotDeg) < 1e-12 {
		x := destX - ox
		y := destY - oy
		tl = hitVec2{x, y}
		bl = hitVec2{x, y + dh}
		br = hitVec2{x + dw, y + dh}
		tr = hitVec2{x + dw, y}
		return
	}
	cos := math.Cos(rotDeg * math.Pi / 180.0)
	sin := math.Sin(rotDeg * math.Pi / 180.0)
	x := destX
	y := destY
	dx := -ox
	dy := -oy
	tl.x = x + dx*cos - dy*sin
	tl.y = y + dx*sin + dy*cos
	tr.x = x + (dx+dw)*cos - dy*sin
	tr.y = y + (dx+dw)*sin + dy*cos
	bl.x = x + dx*cos - (dy+dh)*sin
	bl.y = y + dx*sin + (dy+dh)*cos
	br.x = x + (dx+dw)*cos - (dy+dh)*sin
	br.y = y + (dx+dw)*sin + (dy+dh)*cos
	return
}

func spritePointHit(destX, destY, dw, dh, ox, oy, rotDeg float64, px, py float32) bool {
	if dw <= 0 || dh <= 0 {
		return false
	}
	pxf := float64(px)
	pyf := float64(py)
	if math.Abs(rotDeg) < 1e-12 {
		lx := pxf - (destX - ox)
		ly := pyf - (destY - oy)
		return lx >= 0 && lx < dw && ly >= 0 && ly < dh
	}
	cos := math.Cos(rotDeg * math.Pi / 180.0)
	sin := math.Sin(rotDeg * math.Pi / 180.0)
	Wx := pxf - destX
	Wy := pyf - destY
	lx := ox + Wx*cos + Wy*sin
	ly := oy - Wx*sin + Wy*cos
	return lx >= 0 && lx < dw && ly >= 0 && ly < dh
}

func projectMinMax(pts []hitVec2, ax, ay float64) (min, max float64) {
	min = pts[0].x*ax + pts[0].y*ay
	max = min
	for i := 1; i < len(pts); i++ {
		p := pts[i].x*ax + pts[i].y*ay
		if p < min {
			min = p
		}
		if p > max {
			max = p
		}
	}
	return min, max
}

func axisSeparates(a, b []hitVec2, ax, ay float64) bool {
	al := math.Hypot(ax, ay)
	if al < 1e-15 {
		return false
	}
	nx := ax / al
	ny := ay / al
	amin, amax := projectMinMax(a, nx, ny)
	bmin, bmax := projectMinMax(b, nx, ny)
	return amax < bmin || bmax < amin
}

// satConvexOverlap SAT test for two convex polygons (same vertex order as raylib quad).
func satConvexOverlap(a, b []hitVec2) bool {
	axes := make([]hitVec2, 0, 8)
	addAxes := func(pts []hitVec2) {
		n := len(pts)
		for i := 0; i < n; i++ {
			j := (i + 1) % n
			ex := pts[j].x - pts[i].x
			ey := pts[j].y - pts[i].y
			axes = append(axes, hitVec2{x: -ey, y: ex})
		}
	}
	addAxes(a)
	addAxes(b)
	for _, ax := range axes {
		if axisSeparates(a, b, ax.x, ax.y) {
			return false
		}
	}
	return true
}

func spriteHitOverlap(a, b *spriteObj) bool {
	adx, ady, adw, adh, aox, aoy, ar := spriteHitGeometry(a)
	bdx, bdy, bdw, bdh, box, boy, br := spriteHitGeometry(b)
	if adw <= 0 || adh <= 0 || bdw <= 0 || bdh <= 0 {
		return false
	}
	atl, abl, abr, atr := spriteQuadCorners(adx, ady, adw, adh, aox, aoy, ar)
	btl, bbl, bbr, btr := spriteQuadCorners(bdx, bdy, bdw, bdh, box, boy, br)
	pa := []hitVec2{atl, abl, abr, atr}
	pb := []hitVec2{btl, bbl, bbr, btr}
	return satConvexOverlap(pa, pb)
}

func (m *Module) spHit(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("SPRITE.HIT expects 2 arguments (handleA, handleB)")
	}
	ha, ok := argHandle(args[0])
	if !ok {
		return value.Nil, fmt.Errorf("SPRITE.HIT: invalid handle A")
	}
	hb, ok := argHandle(args[1])
	if !ok {
		return value.Nil, fmt.Errorf("SPRITE.HIT: invalid handle B")
	}
	a, err := heap.Cast[*spriteObj](m.h, ha)
	if err != nil {
		return value.Nil, err
	}
	b, err := heap.Cast[*spriteObj](m.h, hb)
	if err != nil {
		return value.Nil, err
	}
	hit := spriteHitOverlap(a, b)
	return value.FromBool(hit), nil
}

func (m *Module) spPointHit(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("SPRITE.POINTHIT expects 3 arguments (handle, x, y)")
	}
	h, ok := argHandle(args[0])
	if !ok {
		return value.Nil, fmt.Errorf("SPRITE.POINTHIT: invalid sprite handle")
	}
	s, err := heap.Cast[*spriteObj](m.h, h)
	if err != nil {
		return value.Nil, err
	}
	px, ok1 := argF(args[1])
	py, ok2 := argF(args[2])
	if !ok1 || !ok2 {
		return value.Nil, fmt.Errorf("SPRITE.POINTHIT: x and y must be numeric")
	}
	dx, dy, dw, dh, ox, oy, rotDeg := spriteHitGeometry(s)
	inside := spritePointHit(dx, dy, dw, dh, ox, oy, rotDeg, px, py)
	return value.FromBool(inside), nil
}
