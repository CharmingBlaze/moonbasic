//go:build cgo || (windows && !cgo)

package mbdraw

import (
	"fmt"
	"image/color"

	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/runtime"
	"moonbasic/runtime/mbmodel3d"
	"moonbasic/runtime/texture"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

// prim3DKind selects immediate 3D draw state for [drawPrim3D].
type prim3DKind uint8

const (
	primCube prim3DKind = iota
	primCubeWires
	primSphere
	primSphereWires
	primCylinder
	primCylinderWires
	primCapsule
	primCapsuleWires
	primPlane
	primBBox
	primRay
	primLine
	primPoint
	primGrid
	primBillboard
	primBillboardRec
)

// drawPrim3D holds mutable state for object-style DRAW3D wrappers (no native resource — Rule 2).
// Ownership: billboard/ray may reference other heap handles (texture, float array); those are not freed here.
type drawPrim3D struct {
	release heap.ReleaseOnce
	kind    prim3DKind

	x, y, z   float32
	x2, y2, z2 float32
	w, h, d   float32
	rTop, rBot float32
	minx, miny, minz, maxx, maxy, maxz float32
	srcX, srcY, srcW, srcH float32

	slices, rings int32
	gridSlices    int32
	spacing       float32

	rayH  heap.Handle
	texH  heap.Handle

	cr, cg, cb, ca int32
	wire           bool
}

func (o *drawPrim3D) Free() {
	o.release.Do(func() {})
}

func (o *drawPrim3D) TypeTag() uint16 { return heap.TagDrawPrim3D }

func (o *drawPrim3D) TypeName() string {
	switch o.kind {
	case primCube:
		return "DRAWCUBE"
	case primCubeWires:
		return "DRAWCUBEWIRES"
	case primSphere:
		return "DRAWSPHERE"
	case primSphereWires:
		return "DRAWSPHEREW"
	case primCylinder:
		return "DRAWCYLINDER"
	case primCylinderWires:
		return "DRAWCYLINDERW"
	case primCapsule:
		return "DRAWCAP"
	case primCapsuleWires:
		return "DRAWCAPW"
	case primPlane:
		return "DRAWPLANE"
	case primBBox:
		return "DRAWBBOX"
	case primRay:
		return "DRAWRAY"
	case primLine:
		return "DRAWLINE3D"
	case primPoint:
		return "DRAWPOINT3D"
	case primGrid:
		return "DRAWGRID3D"
	case primBillboard:
		return "DRAWBILLBOARD"
	case primBillboardRec:
		return "DRAWBILLBOARDREC"
	default:
		return "DRAWPRIM3D"
	}
}

func prim3DColor(o *drawPrim3D) color.RGBA {
	return color.RGBA{R: uint8(o.cr), G: uint8(o.cg), B: uint8(o.cb), A: uint8(o.ca)}
}

func (m *Module) prim3DDrawSelf(o *drawPrim3D) error {
	if m == nil || m.h == nil {
		return fmt.Errorf("draw: heap not bound")
	}
	col := prim3DColor(o)
	switch o.kind {
	case primCube:
		pos := rl.Vector3{X: o.x, Y: o.y, Z: o.z}
		if o.wire {
			drawCubeWiresRL(pos, o.w, o.h, o.d, col)
		} else {
			drawCubeRL(pos, o.w, o.h, o.d, col)
		}
	case primCubeWires:
		drawCubeWiresRL(rl.Vector3{X: o.x, Y: o.y, Z: o.z}, o.w, o.h, o.d, col)
	case primSphere:
		rl.DrawSphere(rl.Vector3{X: o.x, Y: o.y, Z: o.z}, o.w, col)
	case primSphereWires:
		rl.DrawSphereWires(rl.Vector3{X: o.x, Y: o.y, Z: o.z}, o.w, int32(o.slices), int32(o.rings), col)
	case primCylinder:
		rl.DrawCylinder(rl.Vector3{X: o.x, Y: o.y, Z: o.z}, o.rTop, o.rBot, o.h, int32(o.slices), col)
	case primCylinderWires:
		rl.DrawCylinderWires(rl.Vector3{X: o.x, Y: o.y, Z: o.z}, o.rTop, o.rBot, o.h, int32(o.slices), col)
	case primCapsule:
		rl.DrawCapsule(rl.Vector3{X: o.x, Y: o.y, Z: o.z}, rl.Vector3{X: o.x2, Y: o.y2, Z: o.z2}, o.w, int32(o.slices), int32(o.rings), col)
	case primCapsuleWires:
		rl.DrawCapsuleWires(rl.Vector3{X: o.x, Y: o.y, Z: o.z}, rl.Vector3{X: o.x2, Y: o.y2, Z: o.z2}, o.w, int32(o.slices), int32(o.rings), col)
	case primPlane:
		rl.DrawPlane(rl.Vector3{X: o.x, Y: o.y, Z: o.z}, rl.Vector2{X: o.w, Y: o.d}, col)
	case primBBox:
		bbox := rl.BoundingBox{Min: rl.Vector3{X: o.minx, Y: o.miny, Z: o.minz}, Max: rl.Vector3{X: o.maxx, Y: o.maxy, Z: o.maxz}}
		rl.DrawBoundingBox(bbox, col)
	case primRay:
		if o.rayH == 0 {
			return fmt.Errorf("DRAWRAY.Draw: call SetRay(rayArrayHandle) first")
		}
		if m.h.ArrayFlatLen(o.rayH) != 6 {
			return fmt.Errorf("DRAWRAY: ray array must have 6 floats")
		}
		px, _ := m.h.ArrayGetFloat(o.rayH, 0)
		py, _ := m.h.ArrayGetFloat(o.rayH, 1)
		pz, _ := m.h.ArrayGetFloat(o.rayH, 2)
		dx, _ := m.h.ArrayGetFloat(o.rayH, 3)
		dy, _ := m.h.ArrayGetFloat(o.rayH, 4)
		dz, _ := m.h.ArrayGetFloat(o.rayH, 5)
		ray := rl.Ray{
			Position:  rl.Vector3{X: float32(px), Y: float32(py), Z: float32(pz)},
			Direction: rl.Vector3{X: float32(dx), Y: float32(dy), Z: float32(dz)},
		}
		rl.DrawRay(ray, col)
	case primLine:
		rl.DrawLine3D(rl.Vector3{X: o.x, Y: o.y, Z: o.z}, rl.Vector3{X: o.x2, Y: o.y2, Z: o.z2}, col)
	case primPoint:
		rl.DrawPoint3D(rl.Vector3{X: o.x, Y: o.y, Z: o.z}, col)
	case primGrid:
		rl.DrawGrid(int32(o.gridSlices), o.spacing)
	case primBillboard, primBillboardRec:
		cam, in3D := mbmodel3d.ActiveCamera3D()
		if !in3D {
			return fmt.Errorf("DRAWBILLBOARD.Draw: must be inside Camera.Begin/End")
		}
		if o.texH == 0 {
			return fmt.Errorf("DRAWBILLBOARD.Draw: call SetTexture(texHandle) first")
		}
		tex, err := texture.ForBinding(m.h, o.texH)
		if err != nil {
			return err
		}
		if o.kind == primBillboard {
			rl.DrawBillboard(cam, tex, rl.Vector3{X: o.x, Y: o.y, Z: o.z}, o.w, col)
		} else {
			src := rl.Rectangle{X: o.srcX, Y: o.srcY, Width: o.srcW, Height: o.srcH}
			pos := rl.Vector3{X: o.x, Y: o.y, Z: o.z}
			sz := rl.Vector2{X: o.w, Y: o.h}
			rl.DrawBillboardRec(cam, tex, src, pos, sz, col)
		}
	default:
		return fmt.Errorf("draw: unknown prim3D kind")
	}
	return nil
}

func registerPrim3DWrappers(m *Module, r runtime.Registrar) {
	r.Register("DRAWPRIM3D.POS", "draw", m.prim3DPos)
	r.Register("DRAWPRIM3D.SIZE", "draw", m.prim3DSize)
	r.Register("DRAWPRIM3D.COLOR", "draw", m.prim3DColorCmd)
	r.Register("DRAWPRIM3D.COL", "draw", m.prim3DColorCmd)
	r.Register("DRAWPRIM3D.WIRE", "draw", m.prim3DWire)
	r.Register("DRAWPRIM3D.RADIUS", "draw", m.prim3DRadius)
	r.Register("DRAWPRIM3D.ENDPOINT", "draw", m.prim3DEndPoint)
	r.Register("DRAWPRIM3D.CYL", "draw", m.prim3DCylDims)
	r.Register("DRAWPRIM3D.BBOX", "draw", m.prim3DBBox)
	r.Register("DRAWPRIM3D.SLICES", "draw", m.prim3DSlices)
	r.Register("DRAWPRIM3D.RINGS", "draw", m.prim3DRings)
	r.Register("DRAWPRIM3D.GRID", "draw", m.prim3DGrid)
	r.Register("DRAWPRIM3D.SETRAY", "draw", m.prim3DSetRay)
	r.Register("DRAWPRIM3D.SETTEXTURE", "draw", m.prim3DSetTexture)
	r.Register("DRAWPRIM3D.SRCTEX", "draw", m.prim3DSrcTex)
	r.Register("DRAWPRIM3D.DRAW", "draw", m.prim3DDraw)
	r.Register("DRAWPRIM3D.FREE", "draw", m.prim3DFree)

	r.Register("DRAWCUBE", "draw", runtime.AdaptLegacy(m.makeDrawCube))
	r.Register("DRAWCUBEWIRES", "draw", runtime.AdaptLegacy(m.makeDrawCubeWires))
	r.Register("DRAWSPHERE", "draw", runtime.AdaptLegacy(m.makeDrawSphere))
	r.Register("DRAWSPHEREW", "draw", runtime.AdaptLegacy(m.makeDrawSphereW))
	r.Register("DRAWCYLINDER", "draw", runtime.AdaptLegacy(m.makeDrawCylinder))
	r.Register("DRAWCYLINDERW", "draw", runtime.AdaptLegacy(m.makeDrawCylinderW))
	r.Register("DRAWCAP", "draw", runtime.AdaptLegacy(m.makeDrawCap))
	r.Register("DRAWCAPW", "draw", runtime.AdaptLegacy(m.makeDrawCapW))
	r.Register("DRAWPLANE", "draw", runtime.AdaptLegacy(m.makeDrawPlane))
	r.Register("DRAWBBOX", "draw", runtime.AdaptLegacy(m.makeDrawBBoxPrim))
	r.Register("DRAWRAY", "draw", runtime.AdaptLegacy(m.makeDrawRayPrim))
	r.Register("DRAWLINE3D", "draw", runtime.AdaptLegacy(m.makeDrawLinePrim))
	r.Register("DRAWPOINT3D", "draw", runtime.AdaptLegacy(m.makeDrawPointPrim))
	r.Register("DRAWGRID3D", "draw", runtime.AdaptLegacy(m.makeDrawGridPrim))
	r.Register("DRAWBILLBOARD", "draw", runtime.AdaptLegacy(m.makeDrawBillboardPrim))
	r.Register("DRAWBILLBOARDREC", "draw", runtime.AdaptLegacy(m.makeDrawBillboardRecPrim))
}

func castPrim3D(h *heap.Store, v value.Value) (*drawPrim3D, error) {
	if v.Kind != value.KindHandle {
		return nil, fmt.Errorf("expected handle")
	}
	o, err := heap.Cast[*drawPrim3D](h, heap.Handle(v.IVal))
	if err != nil {
		return nil, err
	}
	return o, nil
}

func (m *Module) prim3DPos(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) != 4 {
		return value.Nil, fmt.Errorf("DRAWPRIM3D.POS expects (handle, x#, y#, z#)")
	}
	o, err := castPrim3D(m.h, args[0])
	if err != nil {
		return value.Nil, err
	}
	x, ok1 := argFloat(args[1])
	y, ok2 := argFloat(args[2])
	z, ok3 := argFloat(args[3])
	if !ok1 || !ok2 || !ok3 {
		return value.Nil, fmt.Errorf("DRAWPRIM3D.POS: numeric coordinates required")
	}
	o.x, o.y, o.z = x, y, z
	return value.Nil, nil
}

func (m *Module) prim3DSize(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) != 2 && len(args) != 4 {
		return value.Nil, fmt.Errorf("DRAWPRIM3D.SIZE expects (handle, w#) or (handle, w#, h#, d#)")
	}
	o, err := castPrim3D(m.h, args[0])
	if err != nil {
		return value.Nil, err
	}
	switch len(args) {
	case 2:
		w, ok := argFloat(args[1])
		if !ok {
			return value.Nil, fmt.Errorf("DRAWPRIM3D.SIZE: width must be numeric")
		}
		switch o.kind {
		case primSphere, primSphereWires:
			o.w = w
		default:
			o.w, o.h, o.d = w, w, w
		}
	case 4:
		w, ok1 := argFloat(args[1])
		h, ok2 := argFloat(args[2])
		d, ok3 := argFloat(args[3])
		if !ok1 || !ok2 || !ok3 {
			return value.Nil, fmt.Errorf("DRAWPRIM3D.SIZE: numeric dimensions required")
		}
		o.w, o.h, o.d = w, h, d
	}
	return value.Nil, nil
}

func (m *Module) prim3DColorCmd(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) != 5 {
		return value.Nil, fmt.Errorf("DRAWPRIM3D.COLOR expects (handle, r, g, b, a)")
	}
	o, err := castPrim3D(m.h, args[0])
	if err != nil {
		return value.Nil, err
	}
	r, ok1 := argInt(args[1])
	g, ok2 := argInt(args[2])
	b, ok3 := argInt(args[3])
	a, ok4 := argInt(args[4])
	if !ok1 || !ok2 || !ok3 || !ok4 {
		return value.Nil, fmt.Errorf("DRAWPRIM3D.COLOR: components must be numeric")
	}
	o.cr, o.cg, o.cb, o.ca = int32(r), int32(g), int32(b), int32(a)
	return value.Nil, nil
}

func (m *Module) prim3DWire(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("DRAWPRIM3D.WIRE expects (handle, on)")
	}
	o, err := castPrim3D(m.h, args[0])
	if err != nil {
		return value.Nil, err
	}
	var pool []string
	if rt.Prog != nil {
		pool = rt.Prog.StringTable
	}
	o.wire = value.Truthy(args[1], pool, m.h)
	return value.Nil, nil
}

func (m *Module) prim3DRadius(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("DRAWPRIM3D.RADIUS expects (handle, radius#)")
	}
	o, err := castPrim3D(m.h, args[0])
	if err != nil {
		return value.Nil, err
	}
	rad, ok := argFloat(args[1])
	if !ok {
		return value.Nil, fmt.Errorf("DRAWPRIM3D.RADIUS: must be numeric")
	}
	o.w = rad
	return value.Nil, nil
}

func (m *Module) prim3DEndPoint(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) != 4 {
		return value.Nil, fmt.Errorf("DRAWPRIM3D.ENDPOINT expects (handle, x#, y#, z#) capsule/line end")
	}
	o, err := castPrim3D(m.h, args[0])
	if err != nil {
		return value.Nil, err
	}
	x, ok1 := argFloat(args[1])
	y, ok2 := argFloat(args[2])
	z, ok3 := argFloat(args[3])
	if !ok1 || !ok2 || !ok3 {
		return value.Nil, fmt.Errorf("DRAWPRIM3D.ENDPOINT: numeric required")
	}
	o.x2, o.y2, o.z2 = x, y, z
	return value.Nil, nil
}

func (m *Module) prim3DCylDims(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) != 4 {
		return value.Nil, fmt.Errorf("DRAWPRIM3D.CYL expects (handle, rTop#, rBot#, h#)")
	}
	o, err := castPrim3D(m.h, args[0])
	if err != nil {
		return value.Nil, err
	}
	rtop, ok1 := argFloat(args[1])
	rbot, ok2 := argFloat(args[2])
	h, ok3 := argFloat(args[3])
	if !ok1 || !ok2 || !ok3 {
		return value.Nil, fmt.Errorf("DRAWPRIM3D.CYL: numeric required")
	}
	o.rTop, o.rBot, o.h = rtop, rbot, h
	return value.Nil, nil
}

func (m *Module) prim3DBBox(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) != 7 {
		return value.Nil, fmt.Errorf("DRAWPRIM3D.BBOX expects (handle, minx,miny,minz, maxx,maxy,maxz)")
	}
	o, err := castPrim3D(m.h, args[0])
	if err != nil {
		return value.Nil, err
	}
	o.minx, _ = argFloat(args[1])
	o.miny, _ = argFloat(args[2])
	o.minz, _ = argFloat(args[3])
	o.maxx, _ = argFloat(args[4])
	o.maxy, _ = argFloat(args[5])
	o.maxz, _ = argFloat(args[6])
	return value.Nil, nil
}

func (m *Module) prim3DSlices(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("DRAWPRIM3D.SLICES expects (handle, slices)")
	}
	o, err := castPrim3D(m.h, args[0])
	if err != nil {
		return value.Nil, err
	}
	s, ok := argInt(args[1])
	if !ok {
		return value.Nil, fmt.Errorf("DRAWPRIM3D.SLICES: must be int")
	}
	o.slices = int32(s)
	return value.Nil, nil
}

func (m *Module) prim3DRings(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("DRAWPRIM3D.RINGS expects (handle, rings)")
	}
	o, err := castPrim3D(m.h, args[0])
	if err != nil {
		return value.Nil, err
	}
	s, ok := argInt(args[1])
	if !ok {
		return value.Nil, fmt.Errorf("DRAWPRIM3D.RINGS: must be int")
	}
	o.rings = int32(s)
	return value.Nil, nil
}

func (m *Module) prim3DGrid(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("DRAWPRIM3D.GRID expects (handle, slices, spacing#)")
	}
	o, err := castPrim3D(m.h, args[0])
	if err != nil {
		return value.Nil, err
	}
	sl, ok1 := argInt(args[1])
	sp, ok2 := argFloat(args[2])
	if !ok1 || !ok2 {
		return value.Nil, fmt.Errorf("DRAWPRIM3D.GRID: numeric required")
	}
	o.gridSlices = int32(sl)
	o.spacing = sp
	return value.Nil, nil
}

func (m *Module) prim3DSetRay(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("DRAWPRIM3D.SETRAY expects (handle, rayArrayHandle)")
	}
	o, err := castPrim3D(m.h, args[0])
	if err != nil {
		return value.Nil, err
	}
	if args[1].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("DRAWPRIM3D.SETRAY: second arg must be handle")
	}
	o.rayH = heap.Handle(args[1].IVal)
	return value.Nil, nil
}

func (m *Module) prim3DSetTexture(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("DRAWPRIM3D.SETTEXTURE expects (handle, texHandle)")
	}
	o, err := castPrim3D(m.h, args[0])
	if err != nil {
		return value.Nil, err
	}
	if args[1].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("DRAWPRIM3D.SETTEXTURE: texture handle required")
	}
	o.texH = heap.Handle(args[1].IVal)
	return value.Nil, nil
}

func (m *Module) prim3DSrcTex(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) != 5 {
		return value.Nil, fmt.Errorf("DRAWPRIM3D.SRCTEX expects (handle, srcX#, srcY#, srcW#, srcH#)")
	}
	o, err := castPrim3D(m.h, args[0])
	if err != nil {
		return value.Nil, err
	}
	o.srcX, _ = argFloat(args[1])
	o.srcY, _ = argFloat(args[2])
	o.srcW, _ = argFloat(args[3])
	o.srcH, _ = argFloat(args[4])
	return value.Nil, nil
}

func (m *Module) prim3DDraw(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("DRAWPRIM3D.DRAW expects (handle)")
	}
	o, err := castPrim3D(m.h, args[0])
	if err != nil {
		return value.Nil, err
	}
	return value.Nil, m.prim3DDrawSelf(o)
}

func (m *Module) prim3DFree(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("DRAWPRIM3D.FREE expects (handle)")
	}
	if args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("DRAWPRIM3D.FREE: handle required")
	}
	if _, err := castPrim3D(m.h, args[0]); err != nil {
		return value.Nil, err
	}
	return value.Nil, rt.Heap.Free(heap.Handle(args[0].IVal))
}

func (m *Module) allocPrim(o *drawPrim3D) (value.Value, error) {
	if m.h == nil {
		return value.Nil, fmt.Errorf("draw: heap not bound")
	}
	h, err := m.h.Alloc(o)
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(h), nil
}

func (m *Module) makeDrawCube(args []value.Value) (value.Value, error) {
	o := &drawPrim3D{kind: primCube, w: 1, h: 1, d: 1, cr: 255, cg: 255, cb: 255, ca: 255}
	switch len(args) {
	case 0:
	case 3:
		w, ok1 := argFloat(args[0])
		h, ok2 := argFloat(args[1])
		d, ok3 := argFloat(args[2])
		if !ok1 || !ok2 || !ok3 {
			return value.Nil, fmt.Errorf("DRAWCUBE: dimensions must be numeric")
		}
		o.w, o.h, o.d = w, h, d
	default:
		return value.Nil, fmt.Errorf("DRAWCUBE expects 0 or 3 arguments (w#, h#, d#)")
	}
	return m.allocPrim(o)
}

func (m *Module) makeDrawCubeWires(args []value.Value) (value.Value, error) {
	o := &drawPrim3D{kind: primCubeWires, w: 1, h: 1, d: 1, cr: 255, cg: 255, cb: 255, ca: 255, wire: true}
	switch len(args) {
	case 0:
	case 3:
		w, ok1 := argFloat(args[0])
		h, ok2 := argFloat(args[1])
		d, ok3 := argFloat(args[2])
		if !ok1 || !ok2 || !ok3 {
			return value.Nil, fmt.Errorf("DRAWCUBEWIRES: dimensions must be numeric")
		}
		o.w, o.h, o.d = w, h, d
	default:
		return value.Nil, fmt.Errorf("DRAWCUBEWIRES expects 0 or 3 arguments")
	}
	return m.allocPrim(o)
}

func (m *Module) makeDrawSphere(args []value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("DRAWSPHERE expects 1 argument (radius#)")
	}
	rad, ok := argFloat(args[0])
	if !ok || rad <= 0 {
		return value.Nil, fmt.Errorf("DRAWSPHERE: radius must be positive")
	}
	o := &drawPrim3D{kind: primSphere, w: rad, cr: 255, cg: 255, cb: 255, ca: 255}
	return m.allocPrim(o)
}

func (m *Module) makeDrawSphereW(args []value.Value) (value.Value, error) {
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("DRAWSPHEREW expects 3 arguments (radius#, rings, slices)")
	}
	rad, ok1 := argFloat(args[0])
	rings, ok2 := argInt(args[1])
	slices, ok3 := argInt(args[2])
	if !ok1 || !ok2 || !ok3 || rad <= 0 {
		return value.Nil, fmt.Errorf("DRAWSPHEREW: invalid arguments")
	}
	o := &drawPrim3D{kind: primSphereWires, w: rad, rings: int32(rings), slices: int32(slices), cr: 255, cg: 255, cb: 255, ca: 255, wire: true}
	return m.allocPrim(o)
}

func (m *Module) makeDrawCylinder(args []value.Value) (value.Value, error) {
	o := &drawPrim3D{kind: primCylinder, rTop: 1, rBot: 1, h: 1, slices: 16, cr: 255, cg: 255, cb: 255, ca: 255}
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("DRAWCYLINDER expects 0 arguments (use Pos/Cyl/Slices/Color)")
	}
	return m.allocPrim(o)
}

func (m *Module) makeDrawCylinderW(args []value.Value) (value.Value, error) {
	o := &drawPrim3D{kind: primCylinderWires, rTop: 1, rBot: 1, h: 1, slices: 16, cr: 255, cg: 255, cb: 255, ca: 255, wire: true}
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("DRAWCYLINDERW expects 0 arguments")
	}
	return m.allocPrim(o)
}

func (m *Module) makeDrawCap(args []value.Value) (value.Value, error) {
	o := &drawPrim3D{
		kind: primCapsule,
		x: 0, y: 0, z: 0, x2: 0, y2: 1, z2: 0,
		w: 0.5, slices: 16, rings: 16,
		cr: 255, cg: 255, cb: 255, ca: 255,
	}
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("DRAWCAP expects 0 arguments")
	}
	return m.allocPrim(o)
}

func (m *Module) makeDrawCapW(args []value.Value) (value.Value, error) {
	o := &drawPrim3D{
		kind: primCapsuleWires,
		x: 0, y: 0, z: 0, x2: 0, y2: 1, z2: 0,
		w: 0.5, slices: 16, rings: 16,
		cr: 255, cg: 255, cb: 255, ca: 255, wire: true,
	}
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("DRAWCAPW expects 0 arguments")
	}
	return m.allocPrim(o)
}

func (m *Module) makeDrawPlane(args []value.Value) (value.Value, error) {
	o := &drawPrim3D{kind: primPlane, w: 10, d: 10, cr: 255, cg: 255, cb: 255, ca: 255}
	switch len(args) {
	case 0:
	case 2:
		w, ok1 := argFloat(args[0])
		d, ok2 := argFloat(args[1])
		if !ok1 || !ok2 {
			return value.Nil, fmt.Errorf("DRAWPLANE: width/depth must be numeric")
		}
		o.w, o.d = w, d
	default:
		return value.Nil, fmt.Errorf("DRAWPLANE expects 0 or 2 arguments (width#, depth#)")
	}
	return m.allocPrim(o)
}

func (m *Module) makeDrawBBoxPrim(args []value.Value) (value.Value, error) {
	o := &drawPrim3D{
		kind: primBBox,
		minx: -1, miny: -1, minz: -1, maxx: 1, maxy: 1, maxz: 1,
		cr: 255, cg: 255, cb: 255, ca: 255,
	}
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("DRAWBBOX expects 0 arguments (use BBox method)")
	}
	return m.allocPrim(o)
}

func (m *Module) makeDrawRayPrim(args []value.Value) (value.Value, error) {
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("DRAWRAY expects 0 arguments (use SetRay)")
	}
	o := &drawPrim3D{kind: primRay, cr: 255, cg: 255, cb: 255, ca: 255}
	return m.allocPrim(o)
}

func (m *Module) makeDrawLinePrim(args []value.Value) (value.Value, error) {
	o := &drawPrim3D{kind: primLine, x2: 0, y2: 0, z2: 1, cr: 255, cg: 255, cb: 255, ca: 255}
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("DRAWLINE3D expects 0 arguments")
	}
	return m.allocPrim(o)
}

func (m *Module) makeDrawPointPrim(args []value.Value) (value.Value, error) {
	o := &drawPrim3D{kind: primPoint, cr: 255, cg: 255, cb: 255, ca: 255}
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("DRAWPOINT3D expects 0 arguments")
	}
	return m.allocPrim(o)
}

func (m *Module) makeDrawGridPrim(args []value.Value) (value.Value, error) {
	o := &drawPrim3D{kind: primGrid, gridSlices: 10, spacing: 1, cr: 255, cg: 255, cb: 255, ca: 255}
	switch len(args) {
	case 0:
	case 2:
		sl, ok1 := argInt(args[0])
		sp, ok2 := argFloat(args[1])
		if !ok1 || !ok2 {
			return value.Nil, fmt.Errorf("DRAWGRID3D: numeric required")
		}
		o.gridSlices = int32(sl)
		o.spacing = sp
	default:
		return value.Nil, fmt.Errorf("DRAWGRID3D expects 0 or 2 arguments (slices, spacing#)")
	}
	return m.allocPrim(o)
}

func (m *Module) makeDrawBillboardPrim(args []value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("DRAWBILLBOARD expects 1 argument (textureHandle)")
	}
	if args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("DRAWBILLBOARD: texture handle required")
	}
	o := &drawPrim3D{kind: primBillboard, texH: heap.Handle(args[0].IVal), w: 1, cr: 255, cg: 255, cb: 255, ca: 255}
	return m.allocPrim(o)
}

func (m *Module) makeDrawBillboardRecPrim(args []value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("DRAWBILLBOARDREC expects 1 argument (textureHandle)")
	}
	if args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("DRAWBILLBOARDREC: texture handle required")
	}
	o := &drawPrim3D{kind: primBillboardRec, texH: heap.Handle(args[0].IVal), w: 1, h: 1, cr: 255, cg: 255, cb: 255, ca: 255}
	return m.allocPrim(o)
}
