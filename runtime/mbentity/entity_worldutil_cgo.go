//go:build cgo || (windows && !cgo)

package mbentity

import (
	"fmt"
	"math"

	mbcamera "moonbasic/runtime/camera"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// entityMatElement matches MAT4.GETELEMENT / column-major Raylib layout (see mbmatrix.helpers_cgo matElement).
func entityMatElement(mat rl.Matrix, row, col int32) float32 {
	switch col {
	case 0:
		switch row {
		case 0:
			return mat.M0
		case 1:
			return mat.M1
		case 2:
			return mat.M2
		case 3:
			return mat.M3
		}
	case 1:
		switch row {
		case 0:
			return mat.M4
		case 1:
			return mat.M5
		case 2:
			return mat.M6
		case 3:
			return mat.M7
		}
	case 2:
		switch row {
		case 0:
			return mat.M8
		case 1:
			return mat.M9
		case 2:
			return mat.M10
		case 3:
			return mat.M11
		}
	case 3:
		switch row {
		case 0:
			return mat.M12
		case 1:
			return mat.M13
		case 2:
			return mat.M14
		case 3:
			return mat.M15
		}
	}
	return 0
}

func matMulDir(m rl.Matrix, v rl.Vector3) rl.Vector3 {
	return rl.Vector3{
		X: m.M0*v.X + m.M4*v.Y + m.M8*v.Z,
		Y: m.M1*v.X + m.M5*v.Y + m.M9*v.Z,
		Z: m.M2*v.X + m.M6*v.Y + m.M10*v.Z,
	}
}

func (m *Module) entMatrixElement(args []value.Value) (value.Value, error) {
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("ENTITY.MATRIXELEMENT expects (entity#, row, col) with row,col in 0..3 (column-major world matrix)")
	}
	id, ok := m.entID(args[0])
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("invalid entity")
	}
	e := m.store().ents[id]
	if e == nil {
		return value.Nil, fmt.Errorf("unknown entity")
	}
	ri, ok1 := args[1].ToInt()
	ci, ok2 := args[2].ToInt()
	if !ok1 {
		if f, ok := args[1].ToFloat(); ok {
			ri = int64(f)
			ok1 = true
		}
	}
	if !ok2 {
		if f, ok := args[2].ToFloat(); ok {
			ci = int64(f)
			ok2 = true
		}
	}
	if !ok1 || !ok2 || ri < 0 || ri > 3 || ci < 0 || ci > 3 {
		return value.Nil, fmt.Errorf("row and col must be 0..3")
	}
	wm := m.worldMatrix(e)
	v := entityMatElement(wm, int32(ri), int32(ci))
	return value.FromFloat(float64(v)), nil
}

func (m *Module) entDeltaAxis(args []value.Value, axis int) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("ENTITY.DELTA* expects (entityA#, entityB#)")
	}
	ia, ok1 := m.entID(args[0])
	ib, ok2 := m.entID(args[1])
	if !ok1 || !ok2 || ia < 1 || ib < 1 {
		return value.Nil, fmt.Errorf("invalid entity ids")
	}
	a := m.store().ents[ia]
	b := m.store().ents[ib]
	if a == nil || b == nil {
		return value.Nil, fmt.Errorf("unknown entity")
	}
	pa := m.worldPos(a)
	pb := m.worldPos(b)
	var d float32
	switch axis {
	case 0:
		d = pb.X - pa.X
	case 1:
		d = pb.Y - pa.Y
	default:
		d = pb.Z - pa.Z
	}
	return value.FromFloat(float64(d)), nil
}

func (m *Module) entTFormPoint(args []value.Value) (value.Value, error) {
	if len(args) != 5 {
		return value.Nil, fmt.Errorf("ENTITY.TFORMPOINT expects (x#, y#, z#, srcEntity#, dstEntity#)")
	}
	x, ok1 := args[0].ToFloat()
	y, ok2 := args[1].ToFloat()
	z, ok3 := args[2].ToFloat()
	if !ok1 || !ok2 || !ok3 {
		return value.Nil, fmt.Errorf("position must be numeric")
	}
	sid, ok4 := m.entID(args[3])
	did, ok5 := m.entID(args[4])
	if !ok4 || !ok5 || sid < 1 || did < 1 {
		return value.Nil, fmt.Errorf("invalid entity ids")
	}
	se := m.store().ents[sid]
	de := m.store().ents[did]
	if se == nil || de == nil {
		return value.Nil, fmt.Errorf("unknown entity")
	}
	Ms := m.worldMatrix(se)
	Md := m.worldMatrix(de)
	invMd := rl.MatrixInvert(Md)
	T := rl.MatrixMultiply(invMd, Ms)
	p := rl.Vector3Transform(rl.Vector3{X: float32(x), Y: float32(y), Z: float32(z)}, T)
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

func (m *Module) entTFormVector(args []value.Value) (value.Value, error) {
	if len(args) != 5 {
		return value.Nil, fmt.Errorf("ENTITY.TFORMVECTOR expects (x#, y#, z#, srcEntity#, dstEntity#)")
	}
	x, ok1 := args[0].ToFloat()
	y, ok2 := args[1].ToFloat()
	z, ok3 := args[2].ToFloat()
	if !ok1 || !ok2 || !ok3 {
		return value.Nil, fmt.Errorf("vector must be numeric")
	}
	sid, ok4 := m.entID(args[3])
	did, ok5 := m.entID(args[4])
	if !ok4 || !ok5 || sid < 1 || did < 1 {
		return value.Nil, fmt.Errorf("invalid entity ids")
	}
	se := m.store().ents[sid]
	de := m.store().ents[did]
	if se == nil || de == nil {
		return value.Nil, fmt.Errorf("unknown entity")
	}
	Ms := m.worldMatrix(se)
	Md := m.worldMatrix(de)
	invMd := rl.MatrixInvert(Md)
	v := rl.Vector3{X: float32(x), Y: float32(y), Z: float32(z)}
	vw := matMulDir(Ms, v)
	vd := matMulDir(invMd, vw)
	arr, err := heap.NewArray([]int64{3})
	if err != nil {
		return value.Nil, err
	}
	_ = arr.Set([]int64{0}, float64(vd.X))
	_ = arr.Set([]int64{1}, float64(vd.Y))
	_ = arr.Set([]int64{2}, float64(vd.Z))
	ph, err := m.h.Alloc(arr)
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(ph), nil
}

func (m *Module) entInView(args []value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, fmt.Errorf("ENTITY.INVIEW: heap not bound")
	}
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("ENTITY.INVIEW expects (entity#, camera)")
	}
	eid, ok := m.entID(args[0])
	if !ok || eid < 1 {
		return value.Nil, fmt.Errorf("invalid entity")
	}
	e := m.store().ents[eid]
	if e == nil {
		return value.Nil, fmt.Errorf("unknown entity")
	}
	if args[1].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("camera must be a handle")
	}
	ch := heap.Handle(args[1].IVal)
	cam, err := mbcamera.RayCamera3D(m.h, ch)
	if err != nil {
		return value.Nil, err
	}
	rw := float32(rl.GetScreenWidth())
	rh := float32(rl.GetScreenHeight())
	aspect := float32(16.0 / 9.0)
	if rh > 1e-3 {
		aspect = rw / rh
	}
	f := mbcamera.ExtractFrustum(cam, aspect)
	vis := entityInFrustum(m, e, f)
	return value.FromBool(vis), nil
}

func entityInFrustum(m *Module, e *ent, f mbcamera.Frustum) bool {
	if e.cullMode == 1 {
		return true
	}
	if e.cullMode == 2 {
		return false
	}
	wp := m.worldPos(e)
	switch e.kind {
	case entKindSphere:
		ms := e.scale.X
		if e.scale.Y > ms {
			ms = e.scale.Y
		}
		if e.scale.Z > ms {
			ms = e.scale.Z
		}
		return f.SphereVisible(wp.X, wp.Y, wp.Z, e.radius*ms)
	case entKindBox, entKindMesh, entKindModel:
		mn, mx := m.aabbWorldMinMax(e)
		return f.AABBVisible(mn.X, mn.Y, mn.Z, mx.X, mx.Y, mx.Z)
	case entKindCylinder:
		h := e.cylH * e.scale.Y
		rs := e.scale.X
		if e.scale.Z > rs {
			rs = e.scale.Z
		}
		rt := e.radius * rs
		rad := float32(math.Sqrt(float64(rt*rt + h*h*0.25)))
		return f.SphereVisible(wp.X, wp.Y, wp.Z, rad)
	case entKindPlane:
		sx := e.w * e.scale.X * 0.5
		sz := e.d * e.scale.Z * 0.5
		if sx < 1e-3 {
			sx = 0.5
		}
		if sz < 1e-3 {
			sz = 0.5
		}
		return f.AABBVisible(wp.X-sx, wp.Y-0.01, wp.Z-sz, wp.X+sx, wp.Y+0.01, wp.Z+sz)
	default:
		if e.ext != nil && e.ext.isSprite {
			sx := e.w * e.scale.X * 0.5
			sy := e.h * e.scale.Y * 0.5
			if sx < 1e-3 {
				sx = 0.5
			}
			if sy < 1e-3 {
				sy = 0.5
			}
			r := sx
			if sy > r {
				r = sy
			}
			return f.SphereVisible(wp.X, wp.Y, wp.Z, r)
		}
		// entKindEmpty and unknown: point test
		return f.PointVisible(wp.X, wp.Y, wp.Z)
	}
}
