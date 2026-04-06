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

func registerDraw3DCmds(m *Module, r runtime.Registrar) {
	r.Register("DRAW3D.GRID", "draw", runtime.AdaptLegacy(m.drawGrid))
	r.Register("DRAW3D.LINE", "draw", runtime.AdaptLegacy(m.drawLine3D))
	r.Register("DRAW3D.POINT", "draw", runtime.AdaptLegacy(m.drawPoint3D))
	r.Register("DRAW3D.SPHERE", "draw", runtime.AdaptLegacy(m.drawSphere))
	r.Register("DRAW3D.SPHEREWIRES", "draw", runtime.AdaptLegacy(m.drawSphereWires))
	r.Register("DRAW3D.CUBE", "draw", runtime.AdaptLegacy(m.drawCube))
	r.Register("DRAW3D.CUBEWIRES", "draw", runtime.AdaptLegacy(m.drawCubeWires))
	r.Register("DRAW3D.CYLINDER", "draw", runtime.AdaptLegacy(m.drawCylinder))
	r.Register("DRAW3D.CYLINDERWIRES", "draw", runtime.AdaptLegacy(m.drawCylinderWires))
	r.Register("DRAW3D.CAPSULE", "draw", runtime.AdaptLegacy(m.drawCapsule))
	r.Register("DRAW3D.CAPSULEWIRES", "draw", runtime.AdaptLegacy(m.drawCapsuleWires))
	r.Register("DRAW3D.PLANE", "draw", runtime.AdaptLegacy(m.drawPlane))
	r.Register("DRAW3D.BBOX", "draw", runtime.AdaptLegacy(m.drawBBox))
	r.Register("DRAW3D.RAY", "draw", runtime.AdaptLegacy(m.drawRay))
	r.Register("DRAW3D.BILLBOARD", "draw", runtime.AdaptLegacy(m.drawBillboard))
	r.Register("DRAW3D.BILLBOARDREC", "draw", runtime.AdaptLegacy(m.drawBillboardRec))
}

func (m *Module) drawRay(args []value.Value) (value.Value, error) {
	if len(args) != 5 {
		return value.Nil, fmt.Errorf("DRAW3D.RAY expects 5 arguments (rayHandle, r,g,b,a)")
	}
	if args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("DRAW3D.RAY: first argument must be a ray handle")
	}
	arrH := heap.Handle(args[0].IVal)
	if m.h.ArrayFlatLen(arrH) != 6 {
		return value.Nil, fmt.Errorf("DRAW3D.RAY: ray handle must be a 6-element float array")
	}
	px, ok0 := m.h.ArrayGetFloat(arrH, 0)
	py, ok1 := m.h.ArrayGetFloat(arrH, 1)
	pz, ok2 := m.h.ArrayGetFloat(arrH, 2)
	dx, ok3 := m.h.ArrayGetFloat(arrH, 3)
	dy, ok4 := m.h.ArrayGetFloat(arrH, 4)
	dz, ok5 := m.h.ArrayGetFloat(arrH, 5)
	if !ok0 || !ok1 || !ok2 || !ok3 || !ok4 || !ok5 {
		return value.Nil, fmt.Errorf("DRAW3D.RAY: ray array data invalid")
	}
	ray := rl.Ray{
		Position:  rl.Vector3{X: float32(px), Y: float32(py), Z: float32(pz)},
		Direction: rl.Vector3{X: float32(dx), Y: float32(dy), Z: float32(dz)},
	}
	r, ok1 := argInt(args[1])
	g, ok2 := argInt(args[2])
	b, ok3 := argInt(args[3])
	a, ok4 := argInt(args[4])
	if !ok1 || !ok2 || !ok3 || !ok4 {
		return value.Nil, fmt.Errorf("DRAW3D.RAY: color components must be numeric")
	}
	col := color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)}
	rl.DrawRay(ray, col)
	return value.Nil, nil
}

func (m *Module) drawBillboard(args []value.Value) (value.Value, error) {
	if len(args) != 9 {
		return value.Nil, fmt.Errorf("DRAW3D.BILLBOARD expects 9 arguments (tex, x,y,z, size, r,g,b,a)")
	}
	cam, in3D := mbmodel3d.ActiveCamera3D()
	if !in3D {
		return value.Nil, fmt.Errorf("DRAW3D.BILLBOARD must be called within a CAMERA.BEGIN/END block")
	}
	tex, err := texture.ForBinding(m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	x, ok1 := argFloat(args[1])
	y, ok2 := argFloat(args[2])
	z, ok3 := argFloat(args[3])
	size, ok4 := argFloat(args[4])
	if !ok1 || !ok2 || !ok3 || !ok4 {
		return value.Nil, fmt.Errorf("DRAW3D.BILLBOARD: geometry arguments must be numeric")
	}
	r, ok5 := argInt(args[5])
	g, ok6 := argInt(args[6])
	b, ok7 := argInt(args[7])
	a, ok8 := argInt(args[8])
	if !ok5 || !ok6 || !ok7 || !ok8 {
		return value.Nil, fmt.Errorf("DRAW3D.BILLBOARD: color components must be numeric")
	}
	col := color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)}
	rl.DrawBillboard(cam, tex, rl.Vector3{X: x, Y: y, Z: z}, size, col)
	return value.Nil, nil
}

func (m *Module) drawBillboardRec(args []value.Value) (value.Value, error) {
	if len(args) != 14 {
		return value.Nil, fmt.Errorf("DRAW3D.BILLBOARDREC expects 14 arguments (tex, srcx,srcy,srcw,srch, x,y,z, w,h, r,g,b,a)")
	}
	cam, in3D := mbmodel3d.ActiveCamera3D()
	if !in3D {
		return value.Nil, fmt.Errorf("DRAW3D.BILLBOARDREC must be called within a CAMERA.BEGIN/END block")
	}
	tex, err := texture.ForBinding(m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	srcx, ok1 := argFloat(args[1])
	srcy, ok2 := argFloat(args[2])
	srcw, ok3 := argFloat(args[3])
	srch, ok4 := argFloat(args[4])
	x, ok5 := argFloat(args[5])
	y, ok6 := argFloat(args[6])
	z, ok7 := argFloat(args[7])
	w, ok8 := argFloat(args[8])
	h, ok9 := argFloat(args[9])
	if !ok1 || !ok2 || !ok3 || !ok4 || !ok5 || !ok6 || !ok7 || !ok8 || !ok9 {
		return value.Nil, fmt.Errorf("DRAW3D.BILLBOARDREC: geometry arguments must be numeric")
	}
	r, ok10 := argInt(args[10])
	g, ok11 := argInt(args[11])
	b, ok12 := argInt(args[12])
	a, ok13 := argInt(args[13])
	if !ok10 || !ok11 || !ok12 || !ok13 {
		return value.Nil, fmt.Errorf("DRAW3D.BILLBOARDREC: color components must be numeric")
	}
	col := color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)}
	sourceRec := rl.Rectangle{X: srcx, Y: srcy, Width: srcw, Height: srch}
	pos := rl.Vector3{X: x, Y: y, Z: z}
	size := rl.Vector2{X: w, Y: h}
	rl.DrawBillboardRec(cam, tex, sourceRec, pos, size, col)
	return value.Nil, nil
}

func (m *Module) drawPlane(args []value.Value) (value.Value, error) {
	if len(args) != 9 {
		return value.Nil, fmt.Errorf("DRAW3D.PLANE expects 9 arguments (x,y,z, w,d, r,g,b,a)")
	}
	x, ok1 := argFloat(args[0])
	y, ok2 := argFloat(args[1])
	z, ok3 := argFloat(args[2])
	w, ok4 := argFloat(args[3])
	d, ok5 := argFloat(args[4])
	if !ok1 || !ok2 || !ok3 || !ok4 || !ok5 {
		return value.Nil, fmt.Errorf("DRAW3D.PLANE: geometry arguments must be numeric")
	}
	r, ok6 := argInt(args[5])
	g, ok7 := argInt(args[6])
	b, ok8 := argInt(args[7])
	a, ok9 := argInt(args[8])
	if !ok6 || !ok7 || !ok8 || !ok9 {
		return value.Nil, fmt.Errorf("DRAW3D.PLANE: color components must be numeric")
	}
	col := color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)}
	rl.DrawPlane(rl.Vector3{X: x, Y: y, Z: z}, rl.Vector2{X: w, Y: d}, col)
	return value.Nil, nil
}

func (m *Module) drawBBox(args []value.Value) (value.Value, error) {
	if len(args) != 10 {
		return value.Nil, fmt.Errorf("DRAW3D.BBOX expects 10 arguments (minx,miny,minz, maxx,maxy,maxz, r,g,b,a)")
	}
	minx, ok1 := argFloat(args[0])
	miny, ok2 := argFloat(args[1])
	minz, ok3 := argFloat(args[2])
	maxx, ok4 := argFloat(args[3])
	maxy, ok5 := argFloat(args[4])
	maxz, ok6 := argFloat(args[5])
	if !ok1 || !ok2 || !ok3 || !ok4 || !ok5 || !ok6 {
		return value.Nil, fmt.Errorf("DRAW3D.BBOX: coordinates must be numeric")
	}
	r, ok7 := argInt(args[6])
	g, ok8 := argInt(args[7])
	b, ok9 := argInt(args[8])
	a, ok10 := argInt(args[9])
	if !ok7 || !ok8 || !ok9 || !ok10 {
		return value.Nil, fmt.Errorf("DRAW3D.BBOX: color components must be numeric")
	}
	col := color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)}
	bbox := rl.BoundingBox{Min: rl.Vector3{X: minx, Y: miny, Z: minz}, Max: rl.Vector3{X: maxx, Y: maxy, Z: maxz}}
	rl.DrawBoundingBox(bbox, col)
	return value.Nil, nil
}

func (m *Module) drawCapsule(args []value.Value) (value.Value, error) {
	if len(args) != 13 {
		return value.Nil, fmt.Errorf("DRAW3D.CAPSULE expects 13 arguments (sx,sy,sz, ex,ey,ez, radius, slices, rings, r,g,b,a)")
	}
	sx, ok1 := argFloat(args[0])
	sy, ok2 := argFloat(args[1])
	sz, ok3 := argFloat(args[2])
	ex, ok4 := argFloat(args[3])
	ey, ok5 := argFloat(args[4])
	ez, ok6 := argFloat(args[5])
	radius, ok7 := argFloat(args[6])
	slices, ok8 := argInt(args[7])
	rings, ok9 := argInt(args[8])
	if !ok1 || !ok2 || !ok3 || !ok4 || !ok5 || !ok6 || !ok7 || !ok8 || !ok9 {
		return value.Nil, fmt.Errorf("DRAW3D.CAPSULE: geometry arguments must be numeric")
	}
	r, ok10 := argInt(args[9])
	g, ok11 := argInt(args[10])
	b, ok12 := argInt(args[11])
	a, ok13 := argInt(args[12])
	if !ok10 || !ok11 || !ok12 || !ok13 {
		return value.Nil, fmt.Errorf("DRAW3D.CAPSULE: color components must be numeric")
	}
	col := color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)}
	rl.DrawCapsule(rl.Vector3{X: sx, Y: sy, Z: sz}, rl.Vector3{X: ex, Y: ey, Z: ez}, radius, int32(slices), int32(rings), col)
	return value.Nil, nil
}

func (m *Module) drawCapsuleWires(args []value.Value) (value.Value, error) {
	if len(args) != 13 {
		return value.Nil, fmt.Errorf("DRAW3D.CAPSULEWIRES expects 13 arguments (sx,sy,sz, ex,ey,ez, radius, slices, rings, r,g,b,a)")
	}
	sx, ok1 := argFloat(args[0])
	sy, ok2 := argFloat(args[1])
	sz, ok3 := argFloat(args[2])
	ex, ok4 := argFloat(args[3])
	ey, ok5 := argFloat(args[4])
	ez, ok6 := argFloat(args[5])
	radius, ok7 := argFloat(args[6])
	slices, ok8 := argInt(args[7])
	rings, ok9 := argInt(args[8])
	if !ok1 || !ok2 || !ok3 || !ok4 || !ok5 || !ok6 || !ok7 || !ok8 || !ok9 {
		return value.Nil, fmt.Errorf("DRAW3D.CAPSULEWIRES: geometry arguments must be numeric")
	}
	r, ok10 := argInt(args[9])
	g, ok11 := argInt(args[10])
	b, ok12 := argInt(args[11])
	a, ok13 := argInt(args[12])
	if !ok10 || !ok11 || !ok12 || !ok13 {
		return value.Nil, fmt.Errorf("DRAW3D.CAPSULEWIRES: color components must be numeric")
	}
	col := color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)}
	rl.DrawCapsuleWires(rl.Vector3{X: sx, Y: sy, Z: sz}, rl.Vector3{X: ex, Y: ey, Z: ez}, radius, int32(slices), int32(rings), col)
	return value.Nil, nil
}

func (m *Module) drawCylinder(args []value.Value) (value.Value, error) {
	if len(args) != 11 {
		return value.Nil, fmt.Errorf("DRAW3D.CYLINDER expects 11 arguments (x,y,z, rTop,rBot, h, slices, r,g,b,a)")
	}
	x, ok1 := argFloat(args[0])
	y, ok2 := argFloat(args[1])
	z, ok3 := argFloat(args[2])
	rTop, ok4 := argFloat(args[3])
	rBot, ok5 := argFloat(args[4])
	h, ok6 := argFloat(args[5])
	slices, ok7 := argInt(args[6])
	if !ok1 || !ok2 || !ok3 || !ok4 || !ok5 || !ok6 || !ok7 {
		return value.Nil, fmt.Errorf("DRAW3D.CYLINDER: geometry arguments must be numeric")
	}
	r, ok8 := argInt(args[7])
	g, ok9 := argInt(args[8])
	b, ok10 := argInt(args[9])
	a, ok11 := argInt(args[10])
	if !ok8 || !ok9 || !ok10 || !ok11 {
		return value.Nil, fmt.Errorf("DRAW3D.CYLINDER: color components must be numeric")
	}
	col := color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)}
	rl.DrawCylinder(rl.Vector3{X: x, Y: y, Z: z}, rTop, rBot, h, int32(slices), col)
	return value.Nil, nil
}

func (m *Module) drawCylinderWires(args []value.Value) (value.Value, error) {
	if len(args) != 11 {
		return value.Nil, fmt.Errorf("DRAW3D.CYLINDERWIRES expects 11 arguments (x,y,z, rTop,rBot, h, slices, r,g,b,a)")
	}
	x, ok1 := argFloat(args[0])
	y, ok2 := argFloat(args[1])
	z, ok3 := argFloat(args[2])
	rTop, ok4 := argFloat(args[3])
	rBot, ok5 := argFloat(args[4])
	h, ok6 := argFloat(args[5])
	slices, ok7 := argInt(args[6])
	if !ok1 || !ok2 || !ok3 || !ok4 || !ok5 || !ok6 || !ok7 {
		return value.Nil, fmt.Errorf("DRAW3D.CYLINDERWIRES: geometry arguments must be numeric")
	}
	r, ok8 := argInt(args[7])
	g, ok9 := argInt(args[8])
	b, ok10 := argInt(args[9])
	a, ok11 := argInt(args[10])
	if !ok8 || !ok9 || !ok10 || !ok11 {
		return value.Nil, fmt.Errorf("DRAW3D.CYLINDERWIRES: color components must be numeric")
	}
	col := color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)}
	rl.DrawCylinderWires(rl.Vector3{X: x, Y: y, Z: z}, rTop, rBot, h, int32(slices), col)
	return value.Nil, nil
}

func (m *Module) drawCube(args []value.Value) (value.Value, error) {
	if len(args) != 10 {
		return value.Nil, fmt.Errorf("DRAW3D.CUBE expects 10 arguments (x,y,z, w,h,d, r,g,b,a)")
	}
	x, ok1 := argFloat(args[0])
	y, ok2 := argFloat(args[1])
	z, ok3 := argFloat(args[2])
	w, ok4 := argFloat(args[3])
	h, ok5 := argFloat(args[4])
	d, ok6 := argFloat(args[5])
	if !ok1 || !ok2 || !ok3 || !ok4 || !ok5 || !ok6 {
		return value.Nil, fmt.Errorf("DRAW3D.CUBE: geometry arguments must be numeric")
	}
	r, ok7 := argInt(args[6])
	g, ok8 := argInt(args[7])
	b, ok9 := argInt(args[8])
	a, ok10 := argInt(args[9])
	if !ok7 || !ok8 || !ok9 || !ok10 {
		return value.Nil, fmt.Errorf("DRAW3D.CUBE: color components must be numeric")
	}
	col := color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)}
	rl.DrawCube(rl.Vector3{X: x, Y: y, Z: z}, w, h, d, col)
	return value.Nil, nil
}

func (m *Module) drawCubeWires(args []value.Value) (value.Value, error) {
	if len(args) != 10 {
		return value.Nil, fmt.Errorf("DRAW3D.CUBEWIRES expects 10 arguments (x,y,z, w,h,d, r,g,b,a)")
	}
	x, ok1 := argFloat(args[0])
	y, ok2 := argFloat(args[1])
	z, ok3 := argFloat(args[2])
	w, ok4 := argFloat(args[3])
	h, ok5 := argFloat(args[4])
	d, ok6 := argFloat(args[5])
	if !ok1 || !ok2 || !ok3 || !ok4 || !ok5 || !ok6 {
		return value.Nil, fmt.Errorf("DRAW3D.CUBEWIRES: geometry arguments must be numeric")
	}
	r, ok7 := argInt(args[6])
	g, ok8 := argInt(args[7])
	b, ok9 := argInt(args[8])
	a, ok10 := argInt(args[9])
	if !ok7 || !ok8 || !ok9 || !ok10 {
		return value.Nil, fmt.Errorf("DRAW3D.CUBEWIRES: color components must be numeric")
	}
	col := color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)}
	rl.DrawCubeWires(rl.Vector3{X: x, Y: y, Z: z}, w, h, d, col)
	return value.Nil, nil
}

func (m *Module) drawSphere(args []value.Value) (value.Value, error) {
	if len(args) != 8 {
		return value.Nil, fmt.Errorf("DRAW3D.SPHERE expects 8 arguments (x,y,z, radius, r,g,b,a)")
	}
	x, ok1 := argFloat(args[0])
	y, ok2 := argFloat(args[1])
	z, ok3 := argFloat(args[2])
	radius, ok4 := argFloat(args[3])
	if !ok1 || !ok2 || !ok3 || !ok4 {
		return value.Nil, fmt.Errorf("DRAW3D.SPHERE: geometry arguments must be numeric")
	}
	r, ok5 := argInt(args[4])
	g, ok6 := argInt(args[5])
	b, ok7 := argInt(args[6])
	a, ok8 := argInt(args[7])
	if !ok5 || !ok6 || !ok7 || !ok8 {
		return value.Nil, fmt.Errorf("DRAW3D.SPHERE: color components must be numeric")
	}
	col := color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)}
	rl.DrawSphere(rl.Vector3{X: x, Y: y, Z: z}, radius, col)
	return value.Nil, nil
}

func (m *Module) drawSphereWires(args []value.Value) (value.Value, error) {
	if len(args) != 10 {
		return value.Nil, fmt.Errorf("DRAW3D.SPHEREWIRES expects 10 arguments (x,y,z, radius, rings, slices, r,g,b,a)")
	}
	x, ok1 := argFloat(args[0])
	y, ok2 := argFloat(args[1])
	z, ok3 := argFloat(args[2])
	radius, ok4 := argFloat(args[3])
	rings, ok5 := argInt(args[4])
	slices, ok6 := argInt(args[5])
	if !ok1 || !ok2 || !ok3 || !ok4 || !ok5 || !ok6 {
		return value.Nil, fmt.Errorf("DRAW3D.SPHEREWIRES: geometry arguments must be numeric")
	}
	r, ok7 := argInt(args[6])
	g, ok8 := argInt(args[7])
	b, ok9 := argInt(args[8])
	a, ok10 := argInt(args[9])
	if !ok7 || !ok8 || !ok9 || !ok10 {
		return value.Nil, fmt.Errorf("DRAW3D.SPHEREWIRES: color components must be numeric")
	}
	col := color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)}
	rl.DrawSphereWires(rl.Vector3{X: x, Y: y, Z: z}, radius, int32(rings), int32(slices), col)
	return value.Nil, nil
}

func (m *Module) drawLine3D(args []value.Value) (value.Value, error) {
	if len(args) != 10 {
		return value.Nil, fmt.Errorf("DRAW3D.LINE expects 10 arguments (x1,y1,z1, x2,y2,z2, r,g,b,a)")
	}
	x1, ok1 := argFloat(args[0])
	y1, ok2 := argFloat(args[1])
	z1, ok3 := argFloat(args[2])
	x2, ok4 := argFloat(args[3])
	y2, ok5 := argFloat(args[4])
	z2, ok6 := argFloat(args[5])
	if !ok1 || !ok2 || !ok3 || !ok4 || !ok5 || !ok6 {
		return value.Nil, fmt.Errorf("DRAW3D.LINE: coordinates must be numeric")
	}
	r, ok7 := argInt(args[6])
	g, ok8 := argInt(args[7])
	b, ok9 := argInt(args[8])
	a, ok10 := argInt(args[9])
	if !ok7 || !ok8 || !ok9 || !ok10 {
		return value.Nil, fmt.Errorf("DRAW3D.LINE: color components must be numeric")
	}
	col := color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)}
	rl.DrawLine3D(rl.Vector3{X: x1, Y: y1, Z: z1}, rl.Vector3{X: x2, Y: y2, Z: z2}, col)
	return value.Nil, nil
}

func (m *Module) drawPoint3D(args []value.Value) (value.Value, error) {
	if len(args) != 7 {
		return value.Nil, fmt.Errorf("DRAW3D.POINT expects 7 arguments (x,y,z, r,g,b,a)")
	}
	x, ok1 := argFloat(args[0])
	y, ok2 := argFloat(args[1])
	z, ok3 := argFloat(args[2])
	if !ok1 || !ok2 || !ok3 {
		return value.Nil, fmt.Errorf("DRAW3D.POINT: coordinates must be numeric")
	}
	r, ok4 := argInt(args[3])
	g, ok5 := argInt(args[4])
	b, ok6 := argInt(args[5])
	a, ok7 := argInt(args[6])
	if !ok4 || !ok5 || !ok6 || !ok7 {
		return value.Nil, fmt.Errorf("DRAW3D.POINT: color components must be numeric")
	}
	col := color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)}
	rl.DrawPoint3D(rl.Vector3{X: x, Y: y, Z: z}, col)
	return value.Nil, nil
}

func (m *Module) drawGrid(args []value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("GRID expects 2 arguments (slices, spacing)")
	}
	slices, ok1 := argInt(args[0])
	spacing, ok2 := argFloat(args[1])
	if !ok1 || !ok2 {
		return value.Nil, fmt.Errorf("GRID: arguments must be numeric")
	}
	rl.DrawGrid(slices, spacing)
	return value.Nil, nil
}
