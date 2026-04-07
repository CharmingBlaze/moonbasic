//go:build cgo || (windows && !cgo)

package mbdraw

import (
	"fmt"
	"image/color"

	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/runtime"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

type prim2DKind uint8

const (
	p2Circle prim2DKind = iota
	p2CircleLines
	p2Ellipse
	p2EllipseLines
	p2Rect
	p2RectLines
	p2Line
	p2Triangle
	p2TriangleLines
	p2Ring
	p2RingLines
	p2Poly
	p2PolyLines
)

// drawPrim2D holds 2D immediate draw state (no native resource).
type drawPrim2D struct {
	release heap.ReleaseOnce
	kind    prim2DKind

	x, y   float32
	x2, y2 float32
	x3, y3 float32
	w, h   float32
	cr, cg, cb, ca int32
	outline        bool
	segments       int32 // ring segment count; also used as poly side count where noted below
	ringInner      float32
	ringOuter      float32
	ringStart      float32
	ringEnd        float32
	polySides      int32
	polyRot        float32
	lineThick      float32
}

func (o *drawPrim2D) Free() { o.release.Do(func() {}) }

func (o *drawPrim2D) TypeTag() uint16 { return heap.TagDrawPrim2D }

func (o *drawPrim2D) TypeName() string {
	switch o.kind {
	case p2Circle:
		return "DRAWCIRCLE2"
	case p2CircleLines:
		return "DRAWCIRCLE2W"
	case p2Ellipse:
		return "DRAWELLIPSE2"
	case p2EllipseLines:
		return "DRAWELLIPSE2W"
	case p2Rect:
		return "DRAWRECT2"
	case p2RectLines:
		return "DRAWRECT2W"
	case p2Line:
		return "DRAWLINE2"
	case p2Triangle:
		return "DRAWTRI2"
	case p2TriangleLines:
		return "DRAWTRI2W"
	case p2Ring:
		return "DRAWRING2"
	case p2RingLines:
		return "DRAWRING2W"
	case p2Poly:
		return "DRAWPOLY2"
	case p2PolyLines:
		return "DRAWPOLY2W"
	default:
		return "DRAWPRIM2D"
	}
}

func (m *Module) prim2DDrawSelf(o *drawPrim2D) error {
	col := color.RGBA{R: uint8(o.cr), G: uint8(o.cg), B: uint8(o.cb), A: uint8(o.ca)}
	switch o.kind {
	case p2Circle:
		if o.outline {
			rl.DrawCircleLines(int32(o.x), int32(o.y), o.w, col)
		} else {
			rl.DrawCircle(int32(o.x), int32(o.y), o.w, col)
		}
	case p2CircleLines:
		rl.DrawCircleLines(int32(o.x), int32(o.y), o.w, col)
	case p2Ellipse:
		if o.outline {
			rl.DrawEllipseLines(int32(o.x), int32(o.y), o.w, o.h, col)
		} else {
			rl.DrawEllipse(int32(o.x), int32(o.y), o.w, o.h, col)
		}
	case p2EllipseLines:
		rl.DrawEllipseLines(int32(o.x), int32(o.y), o.w, o.h, col)
	case p2Rect:
		if o.outline {
			rl.DrawRectangleLines(int32(o.x), int32(o.y), int32(o.w), int32(o.h), col)
		} else {
			rl.DrawRectangle(int32(o.x), int32(o.y), int32(o.w), int32(o.h), col)
		}
	case p2RectLines:
		rl.DrawRectangleLines(int32(o.x), int32(o.y), int32(o.w), int32(o.h), col)
	case p2Line:
		rl.DrawLine(int32(o.x), int32(o.y), int32(o.x2), int32(o.y2), col)
	case p2Triangle:
		if o.outline {
			rl.DrawTriangleLines(rl.Vector2{X: o.x, Y: o.y}, rl.Vector2{X: o.x2, Y: o.y2}, rl.Vector2{X: o.x3, Y: o.y3}, col)
		} else {
			rl.DrawTriangle(rl.Vector2{X: o.x, Y: o.y}, rl.Vector2{X: o.x2, Y: o.y2}, rl.Vector2{X: o.x3, Y: o.y3}, col)
		}
	case p2TriangleLines:
		rl.DrawTriangleLines(rl.Vector2{X: o.x, Y: o.y}, rl.Vector2{X: o.x2, Y: o.y2}, rl.Vector2{X: o.x3, Y: o.y3}, col)
	case p2Ring:
		rl.DrawRing(rl.Vector2{X: o.x, Y: o.y}, o.ringInner, o.ringOuter, o.ringStart, o.ringEnd, o.segments, col)
	case p2RingLines:
		rl.DrawRingLines(rl.Vector2{X: o.x, Y: o.y}, o.ringInner, o.ringOuter, o.ringStart, o.ringEnd, o.segments, col)
	case p2Poly:
		rl.DrawPoly(rl.Vector2{X: o.x, Y: o.y}, o.polySides, o.w, o.polyRot, col)
	case p2PolyLines:
		rl.DrawPolyLinesEx(rl.Vector2{X: o.x, Y: o.y}, o.polySides, o.w, o.polyRot, o.lineThick, col)
	default:
		return fmt.Errorf("draw2d: unknown kind")
	}
	return nil
}

func registerPrim2DWrappers(m *Module, r runtime.Registrar) {
	r.Register("DRAWPRIM2D.POS", "draw", m.prim2DPos)
	r.Register("DRAWPRIM2D.SIZE", "draw", m.prim2DSize)
	r.Register("DRAWPRIM2D.COLOR", "draw", m.prim2DColorCmd)
	r.Register("DRAWPRIM2D.COL", "draw", m.prim2DColorCmd)
	r.Register("DRAWPRIM2D.OUTLINE", "draw", m.prim2DOutline)
	r.Register("DRAWPRIM2D.P2", "draw", m.prim2DP2)
	r.Register("DRAWPRIM2D.P3", "draw", m.prim2DP3)
	r.Register("DRAWPRIM2D.RING", "draw", m.prim2DRingDims)
	r.Register("DRAWPRIM2D.SEGS", "draw", m.prim2DSegs)
	r.Register("DRAWPRIM2D.SIDES", "draw", m.prim2DSides)
	r.Register("DRAWPRIM2D.ROT", "draw", m.prim2DRot)
	r.Register("DRAWPRIM2D.THICK", "draw", m.prim2DThick)
	r.Register("DRAWPRIM2D.DRAW", "draw", m.prim2DDraw)
	r.Register("DRAWPRIM2D.FREE", "draw", m.prim2DFree)

	r.Register("DRAWCIRCLE2", "draw", runtime.AdaptLegacy(m.makeDrawCircle2))
	r.Register("DRAWCIRCLE2W", "draw", runtime.AdaptLegacy(m.makeDrawCircle2W))
	r.Register("DRAWELLIPSE2", "draw", runtime.AdaptLegacy(m.makeDrawEllipse2))
	r.Register("DRAWELLIPSE2W", "draw", runtime.AdaptLegacy(m.makeDrawEllipse2W))
	r.Register("DRAWRECT2", "draw", runtime.AdaptLegacy(m.makeDrawRect2))
	r.Register("DRAWRECT2W", "draw", runtime.AdaptLegacy(m.makeDrawRect2W))
	r.Register("DRAWLINE2", "draw", runtime.AdaptLegacy(m.makeDrawLine2))
	r.Register("DRAWTRI2", "draw", runtime.AdaptLegacy(m.makeDrawTri2))
	r.Register("DRAWTRI2W", "draw", runtime.AdaptLegacy(m.makeDrawTri2W))
	r.Register("DRAWRING2", "draw", runtime.AdaptLegacy(m.makeDrawRing2))
	r.Register("DRAWRING2W", "draw", runtime.AdaptLegacy(m.makeDrawRing2W))
	r.Register("DRAWPOLY2", "draw", runtime.AdaptLegacy(m.makeDrawPoly2))
	r.Register("DRAWPOLY2W", "draw", runtime.AdaptLegacy(m.makeDrawPoly2W))
}

func castPrim2D(h *heap.Store, v value.Value) (*drawPrim2D, error) {
	if v.Kind != value.KindHandle {
		return nil, fmt.Errorf("expected handle")
	}
	return heap.Cast[*drawPrim2D](h, heap.Handle(v.IVal))
}

func (m *Module) prim2DPos(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("DRAWPRIM2D.POS expects (handle, x#, y#)")
	}
	o, err := castPrim2D(m.h, args[0])
	if err != nil {
		return value.Nil, err
	}
	x, ok1 := argFloat(args[1])
	y, ok2 := argFloat(args[2])
	if !ok1 || !ok2 {
		return value.Nil, fmt.Errorf("DRAWPRIM2D.POS: numeric required")
	}
	o.x, o.y = x, y
	return value.Nil, nil
}

func (m *Module) prim2DSize(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	o, err := castPrim2D(m.h, args[0])
	if err != nil {
		return value.Nil, err
	}
	switch len(args) {
	case 2:
		rad, ok := argFloat(args[1])
		if !ok {
			return value.Nil, fmt.Errorf("DRAWPRIM2D.SIZE: radius must be numeric")
		}
		switch o.kind {
		case p2Circle, p2CircleLines, p2Poly, p2PolyLines:
			o.w = rad
		default:
			return value.Nil, fmt.Errorf("DRAWPRIM2D.SIZE: use (handle, w#, h#) for this shape")
		}
	case 3:
		a, ok1 := argFloat(args[1])
		b, ok2 := argFloat(args[2])
		if !ok1 || !ok2 {
			return value.Nil, fmt.Errorf("DRAWPRIM2D.SIZE: numeric required")
		}
		switch o.kind {
		case p2Ring, p2RingLines:
			o.ringInner, o.ringOuter = a, b
		default:
			o.w, o.h = a, b
		}
	default:
		return value.Nil, fmt.Errorf("DRAWPRIM2D.SIZE expects (handle, radius#) or (handle, w#, h#)")
	}
	return value.Nil, nil
}

func (m *Module) prim2DColorCmd(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) != 5 {
		return value.Nil, fmt.Errorf("DRAWPRIM2D.COLOR expects (handle, r,g,b,a)")
	}
	o, err := castPrim2D(m.h, args[0])
	if err != nil {
		return value.Nil, err
	}
	r, ok1 := argInt(args[1])
	g, ok2 := argInt(args[2])
	b, ok3 := argInt(args[3])
	a, ok4 := argInt(args[4])
	if !ok1 || !ok2 || !ok3 || !ok4 {
		return value.Nil, fmt.Errorf("DRAWPRIM2D.COLOR: numeric")
	}
	o.cr, o.cg, o.cb, o.ca = int32(r), int32(g), int32(b), int32(a)
	return value.Nil, nil
}

func (m *Module) prim2DOutline(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("DRAWPRIM2D.OUTLINE expects (handle, on)")
	}
	o, err := castPrim2D(m.h, args[0])
	if err != nil {
		return value.Nil, err
	}
	var pool []string
	if rt.Prog != nil {
		pool = rt.Prog.StringTable
	}
	o.outline = value.Truthy(args[1], pool, m.h)
	return value.Nil, nil
}

func (m *Module) prim2DP2(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("DRAWPRIM2D.P2 expects (handle, x#, y#) second point for line/triangle")
	}
	o, err := castPrim2D(m.h, args[0])
	if err != nil {
		return value.Nil, err
	}
	x, ok1 := argFloat(args[1])
	y, ok2 := argFloat(args[2])
	if !ok1 || !ok2 {
		return value.Nil, fmt.Errorf("DRAWPRIM2D.P2: numeric")
	}
	o.x2, o.y2 = x, y
	return value.Nil, nil
}

func (m *Module) prim2DP3(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("DRAWPRIM2D.P3 expects (handle, x#, y#) third vertex")
	}
	o, err := castPrim2D(m.h, args[0])
	if err != nil {
		return value.Nil, err
	}
	x, ok1 := argFloat(args[1])
	y, ok2 := argFloat(args[2])
	if !ok1 || !ok2 {
		return value.Nil, fmt.Errorf("DRAWPRIM2D.P3: numeric")
	}
	o.x3, o.y3 = x, y
	return value.Nil, nil
}

func (m *Module) prim2DRingDims(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) != 6 {
		return value.Nil, fmt.Errorf("DRAWPRIM2D.RING expects (handle, inner#, outer#, start#, end#)")
	}
	o, err := castPrim2D(m.h, args[0])
	if err != nil {
		return value.Nil, err
	}
	if o.kind != p2Ring && o.kind != p2RingLines {
		return value.Nil, fmt.Errorf("DRAWPRIM2D.RING: not a ring wrapper")
	}
	inner, ok1 := argFloat(args[1])
	outer, ok2 := argFloat(args[2])
	start, ok3 := argFloat(args[3])
	end, ok4 := argFloat(args[4])
	if !ok1 || !ok2 || !ok3 || !ok4 {
		return value.Nil, fmt.Errorf("DRAWPRIM2D.RING: numeric")
	}
	o.ringInner, o.ringOuter, o.ringStart, o.ringEnd = inner, outer, start, end
	return value.Nil, nil
}

func (m *Module) prim2DSegs(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("DRAWPRIM2D.SEGS expects (handle, segments#)")
	}
	o, err := castPrim2D(m.h, args[0])
	if err != nil {
		return value.Nil, err
	}
	if o.kind != p2Ring && o.kind != p2RingLines {
		return value.Nil, fmt.Errorf("DRAWPRIM2D.SEGS: ring only")
	}
	n, ok := argInt(args[1])
	if !ok || n < 3 {
		return value.Nil, fmt.Errorf("DRAWPRIM2D.SEGS: int >= 3")
	}
	o.segments = n
	return value.Nil, nil
}

func (m *Module) prim2DSides(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("DRAWPRIM2D.SIDES expects (handle, sides#)")
	}
	o, err := castPrim2D(m.h, args[0])
	if err != nil {
		return value.Nil, err
	}
	if o.kind != p2Poly && o.kind != p2PolyLines {
		return value.Nil, fmt.Errorf("DRAWPRIM2D.SIDES: polygon only")
	}
	n, ok := argInt(args[1])
	if !ok || n < 3 {
		return value.Nil, fmt.Errorf("DRAWPRIM2D.SIDES: int >= 3")
	}
	o.polySides = n
	return value.Nil, nil
}

func (m *Module) prim2DRot(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("DRAWPRIM2D.ROT expects (handle, rotation#)")
	}
	o, err := castPrim2D(m.h, args[0])
	if err != nil {
		return value.Nil, err
	}
	if o.kind != p2Poly && o.kind != p2PolyLines {
		return value.Nil, fmt.Errorf("DRAWPRIM2D.ROT: polygon only")
	}
	rad, ok := argFloat(args[1])
	if !ok {
		return value.Nil, fmt.Errorf("DRAWPRIM2D.ROT: numeric")
	}
	o.polyRot = rad
	return value.Nil, nil
}

func (m *Module) prim2DThick(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("DRAWPRIM2D.THICK expects (handle, thick#)")
	}
	o, err := castPrim2D(m.h, args[0])
	if err != nil {
		return value.Nil, err
	}
	if o.kind != p2PolyLines {
		return value.Nil, fmt.Errorf("DRAWPRIM2D.THICK: DRAWPOLY2W only")
	}
	t, ok := argFloat(args[1])
	if !ok || t <= 0 {
		return value.Nil, fmt.Errorf("DRAWPRIM2D.THICK: positive number")
	}
	o.lineThick = t
	return value.Nil, nil
}

func (m *Module) prim2DDraw(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("DRAWPRIM2D.DRAW expects (handle)")
	}
	o, err := castPrim2D(m.h, args[0])
	if err != nil {
		return value.Nil, err
	}
	return value.Nil, m.prim2DDrawSelf(o)
}

func (m *Module) prim2DFree(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("DRAWPRIM2D.FREE expects handle")
	}
	if _, err := castPrim2D(m.h, args[0]); err != nil {
		return value.Nil, err
	}
	return value.Nil, rt.Heap.Free(heap.Handle(args[0].IVal))
}

func (m *Module) allocPrim2(o *drawPrim2D) (value.Value, error) {
	if m.h == nil {
		return value.Nil, fmt.Errorf("draw: heap not bound")
	}
	h, err := m.h.Alloc(o)
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(h), nil
}

func (m *Module) makeDrawCircle2(args []value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("DRAWCIRCLE2 expects 1 argument (radius#)")
	}
	rad, ok := argFloat(args[0])
	if !ok || rad <= 0 {
		return value.Nil, fmt.Errorf("DRAWCIRCLE2: positive radius required")
	}
	o := &drawPrim2D{kind: p2Circle, w: rad, cr: 255, cg: 255, cb: 255, ca: 255}
	return m.allocPrim2(o)
}

func (m *Module) makeDrawCircle2W(args []value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("DRAWCIRCLE2W expects 1 argument (radius#)")
	}
	rad, ok := argFloat(args[0])
	if !ok || rad <= 0 {
		return value.Nil, fmt.Errorf("DRAWCIRCLE2W: positive radius required")
	}
	o := &drawPrim2D{kind: p2CircleLines, w: rad, cr: 255, cg: 255, cb: 255, ca: 255}
	return m.allocPrim2(o)
}

func (m *Module) makeDrawEllipse2(args []value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("DRAWELLIPSE2 expects 2 arguments (rx#, ry#)")
	}
	rx, ok1 := argFloat(args[0])
	ry, ok2 := argFloat(args[1])
	if !ok1 || !ok2 {
		return value.Nil, fmt.Errorf("DRAWELLIPSE2: numeric")
	}
	o := &drawPrim2D{kind: p2Ellipse, w: rx, h: ry, cr: 255, cg: 255, cb: 255, ca: 255}
	return m.allocPrim2(o)
}

func (m *Module) makeDrawEllipse2W(args []value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("DRAWELLIPSE2W expects 2 arguments (rx#, ry#)")
	}
	rx, ok1 := argFloat(args[0])
	ry, ok2 := argFloat(args[1])
	if !ok1 || !ok2 {
		return value.Nil, fmt.Errorf("DRAWELLIPSE2W: numeric")
	}
	o := &drawPrim2D{kind: p2EllipseLines, w: rx, h: ry, cr: 255, cg: 255, cb: 255, ca: 255}
	return m.allocPrim2(o)
}

func (m *Module) makeDrawRect2(args []value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("DRAWRECT2 expects 2 arguments (w, h)")
	}
	w, ok1 := argFloat(args[0])
	h, ok2 := argFloat(args[1])
	if !ok1 || !ok2 {
		return value.Nil, fmt.Errorf("DRAWRECT2: numeric")
	}
	o := &drawPrim2D{kind: p2Rect, w: w, h: h, cr: 255, cg: 255, cb: 255, ca: 255}
	return m.allocPrim2(o)
}

func (m *Module) makeDrawRect2W(args []value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("DRAWRECT2W expects 2 arguments (w, h)")
	}
	w, ok1 := argFloat(args[0])
	h, ok2 := argFloat(args[1])
	if !ok1 || !ok2 {
		return value.Nil, fmt.Errorf("DRAWRECT2W: numeric")
	}
	o := &drawPrim2D{kind: p2RectLines, w: w, h: h, cr: 255, cg: 255, cb: 255, ca: 255}
	return m.allocPrim2(o)
}

func (m *Module) makeDrawLine2(args []value.Value) (value.Value, error) {
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("DRAWLINE2 expects 0 arguments (use Pos and P2)")
	}
	o := &drawPrim2D{kind: p2Line, x2: 1, y2: 1, cr: 255, cg: 255, cb: 255, ca: 255}
	return m.allocPrim2(o)
}

func (m *Module) makeDrawTri2(args []value.Value) (value.Value, error) {
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("DRAWTRI2 expects 0 arguments (use Pos, P2, P3)")
	}
	o := &drawPrim2D{kind: p2Triangle, x2: 1, y2: 0, x3: 0, y3: 1, cr: 255, cg: 255, cb: 255, ca: 255}
	return m.allocPrim2(o)
}

func (m *Module) makeDrawTri2W(args []value.Value) (value.Value, error) {
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("DRAWTRI2W expects 0 arguments")
	}
	o := &drawPrim2D{kind: p2TriangleLines, x2: 1, y2: 0, x3: 0, y3: 1, cr: 255, cg: 255, cb: 255, ca: 255}
	return m.allocPrim2(o)
}

func (m *Module) makeDrawRing2(args []value.Value) (value.Value, error) {
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("DRAWRING2 expects 0 arguments (use Pos, Size or RING, Segs, Color, Draw)")
	}
	o := &drawPrim2D{
		kind: p2Ring, ringInner: 20, ringOuter: 40, ringStart: 0, ringEnd: 360,
		segments: 32, cr: 255, cg: 255, cb: 255, ca: 255,
	}
	return m.allocPrim2(o)
}

func (m *Module) makeDrawRing2W(args []value.Value) (value.Value, error) {
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("DRAWRING2W expects 0 arguments")
	}
	o := &drawPrim2D{
		kind: p2RingLines, ringInner: 20, ringOuter: 40, ringStart: 0, ringEnd: 360,
		segments: 32, cr: 255, cg: 255, cb: 255, ca: 255,
	}
	return m.allocPrim2(o)
}

func (m *Module) makeDrawPoly2(args []value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("DRAWPOLY2 expects 1 argument (sides#)")
	}
	sides, ok := argInt(args[0])
	if !ok || sides < 3 {
		return value.Nil, fmt.Errorf("DRAWPOLY2: sides >= 3")
	}
	o := &drawPrim2D{
		kind: p2Poly, w: 50, polySides: sides, polyRot: 0,
		cr: 255, cg: 255, cb: 255, ca: 255,
	}
	return m.allocPrim2(o)
}

func (m *Module) makeDrawPoly2W(args []value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("DRAWPOLY2W expects 1 argument (sides#)")
	}
	sides, ok := argInt(args[0])
	if !ok || sides < 3 {
		return value.Nil, fmt.Errorf("DRAWPOLY2W: sides >= 3")
	}
	o := &drawPrim2D{
		kind: p2PolyLines, w: 50, polySides: sides, polyRot: 0, lineThick: 1,
		cr: 255, cg: 255, cb: 255, ca: 255,
	}
	return m.allocPrim2(o)
}
