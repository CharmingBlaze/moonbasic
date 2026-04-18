//go:build cgo || (windows && !cgo)

package water

import (
	"fmt"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/runtime/mbmatrix"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

// module implementation below

// PointInWaterVolume reports whether (x,y,z) lies inside any water column (between bed and wavy surface).
func (m *Module) PointInWaterVolume(x, y, z float32) bool {
	if m.h == nil {
		return false
	}
	for _, id := range m.waters {
		o, err := castW(m, id)
		if err != nil {
			continue
		}
		halfW := o.Width * 0.5
		halfD := o.Depth * 0.5
		if x < o.PX-halfW || x > o.PX+halfW || z < o.PZ-halfD || z > o.PZ+halfD {
			continue
		}
		bob := float32(math.Sin(float64(o.WaveT))) * o.WaveAmp * 0.15
		surfY := o.PY + bob
		if y <= surfY && y >= o.BedY {
			return true
		}
	}
	return false
}

func castW(m *Module, h heap.Handle) (*WaterObject, error) {
	return heap.Cast[*WaterObject](m.h, h)
}

func (m *Module) wMake(args []value.Value) (value.Value, error) {
	if m.h == nil || len(args) != 2 {
		return value.Nil, fmt.Errorf("WATER.MAKE expects width, depth")
	}
	wf, ok1 := args[0].ToFloat()
	df, ok2 := args[1].ToFloat()
	if !ok1 || !ok2 {
		return value.Nil, fmt.Errorf("WATER.MAKE: arguments must be numeric")
	}
	w := float32(wf)
	d := float32(df)
	if w <= 0 || d <= 0 {
		return value.Nil, fmt.Errorf("WATER.MAKE: dimensions must be > 0")
	}
	mesh := rl.GenMeshPlane(w, d, 32, 32)
	mat := rl.LoadMaterialDefault()
	o := &WaterObject{
		Mesh:     mesh,
		Mat:      mat,
		Width:    w,
		Depth:    d,
		WaveAmp:  0.35,
		WaveFreq: 1.2,
		Shallow:  rl.Color{R: 0, G: 160, B: 200, A: 180},
		Deep:     rl.Color{R: 0, G: 50, B: 100, A: 230},
		BedY:     -50,
		ScaleX:   1, ScaleY: 1, ScaleZ: 1,
	}
	id, err := m.h.Alloc(o)
	if err != nil {
		return value.Nil, err
	}
	m.waters = append(m.waters, id)
	return value.FromHandle(int32(id)), nil
}

func (m *Module) wCreate(args []value.Value) (value.Value, error) {
	if m.h == nil || len(args) != 5 {
		return value.Nil, fmt.Errorf("WATER.CREATE expects x, z, width, depth, level")
	}
	x, _ := args[0].ToFloat()
	z, _ := args[1].ToFloat()
	wf, _ := args[2].ToFloat()
	df, _ := args[3].ToFloat()
	level, _ := args[4].ToFloat()
	w := float32(wf)
	d := float32(df)
	if w <= 0 || d <= 0 {
		return value.Nil, fmt.Errorf("WATER.CREATE: width and depth must be > 0")
	}
	mesh := rl.GenMeshPlane(w, d, 32, 32)
	mat := rl.LoadMaterialDefault()
	o := &WaterObject{
		Mesh:     mesh,
		Mat:      mat,
		Width:    w,
		Depth:    d,
		WaveAmp:  0.35,
		WaveFreq: 1.2,
		Shallow:  rl.Color{R: 0, G: 160, B: 200, A: 180},
		Deep:     rl.Color{R: 0, G: 50, B: 100, A: 230},
		BedY:     -50,
		PX:       float32(x),
		PY:       float32(level),
		PZ:       float32(z),
		ScaleX:   1, ScaleY: 1, ScaleZ: 1,
	}
	o.BedY = o.PY - 12
	id, err := m.h.Alloc(o)
	if err != nil {
		return value.Nil, err
	}
	m.waters = append(m.waters, id)
	return value.FromHandle(int32(id)), nil
}

func (m *Module) wSetWave(args []value.Value) (value.Value, error) {
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("WATER.SETWAVE expects water, speed, height")
	}
	h := heap.Handle(args[0].IVal)
	o, err := castW(m, h)
	if err != nil {
		return value.Nil, err
	}
	sp, _ := args[1].ToFloat()
	amp, _ := args[2].ToFloat()
	o.WaveFreq = float32(sp)
	o.WaveAmp = float32(amp)
	return args[0], nil
}

func (m *Module) wFree(args []value.Value) (value.Value, error) {
	if m.h == nil || len(args) != 1 {
		return value.Nil, fmt.Errorf("WATER.FREE expects handle")
	}
	hh := heap.Handle(args[0].IVal)
	for i, id := range m.waters {
		if id == hh {
			m.waters = append(m.waters[:i], m.waters[i+1:]...)
			break
		}
	}
	m.h.Free(hh)
	return value.Nil, nil
}

func (m *Module) wSetPos(args []value.Value) (value.Value, error) {
	if len(args) != 4 {
		return value.Nil, fmt.Errorf("WATER.SETPOS expects water, x, y, z")
	}
	h := heap.Handle(args[0].IVal)
	o, err := castW(m, h)
	if err != nil {
		return value.Nil, err
	}
	x, _ := args[1].ToFloat()
	y, _ := args[2].ToFloat()
	z, _ := args[3].ToFloat()
	o.PX = float32(x)
	o.PY = float32(y)
	o.PZ = float32(z)
	o.BedY = o.PY - 12
	return args[0], nil
}

func (m *Module) wDraw(args []value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("WATER.DRAW expects handle")
	}
	h := heap.Handle(args[0].IVal)
	o, err := castW(m, h)
	if err != nil {
		return value.Nil, err
	}
	bob := float32(math.Sin(float64(o.WaveT))) * o.WaveAmp * 0.15
	trans := rl.MatrixTranslate(o.PX, o.PY+bob, o.PZ)
	rot := rl.MatrixRotateXYZ(rl.Vector3{X: o.RotX * rl.Deg2rad, Y: o.RotY * rl.Deg2rad, Z: o.RotZ * rl.Deg2rad})
	scale := rl.MatrixScale(o.ScaleX, o.ScaleY, o.ScaleZ)
	mat := rl.MatrixMultiply(rl.MatrixMultiply(scale, rot), trans)
	rl.DrawMesh(o.Mesh, o.Mat, mat)
	return args[0], nil
}

func (m *Module) wUpdate(args []value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("WATER.UPDATE expects dt")
	}
	dt, _ := args[0].ToFloat()
	for _, id := range m.waters {
		o, err := castW(m, id)
		if err != nil {
			continue
		}
		o.WaveT += float32(dt) * o.WaveFreq
	}
	return value.Nil, nil
}

func (m *Module) wSetWaveHeight(args []value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("WATER.SETWAVEHEIGHT expects water, height")
	}
	h := heap.Handle(args[0].IVal)
	o, err := castW(m, h)
	if err != nil {
		return value.Nil, err
	}
	v, _ := args[1].ToFloat()
	o.WaveAmp = float32(v)
	return args[0], nil
}

func (m *Module) wGetWaveY(args []value.Value) (value.Value, error) {
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("WATER.GETWAVEY expects water, x, z")
	}
	h := heap.Handle(args[0].IVal)
	o, err := castW(m, h)
	if err != nil {
		return value.Nil, err
	}
	bob := float32(math.Sin(float64(o.WaveT))) * o.WaveAmp * 0.15
	return value.FromFloat(float64(o.PY + bob)), nil
}

func (m *Module) wGetDepth(args []value.Value) (value.Value, error) {
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("WATER.GETDEPTH expects water, x, z")
	}
	h := heap.Handle(args[0].IVal)
	o, err := castW(m, h)
	if err != nil {
		return value.Nil, err
	}
	bob := float32(math.Sin(float64(o.WaveT))) * o.WaveAmp * 0.15
	surfY := float64(o.PY + bob)
	return value.FromFloat(surfY - float64(o.BedY)), nil
}

func (m *Module) wIsUnder(args []value.Value) (value.Value, error) {
	if len(args) != 4 {
		return value.Nil, fmt.Errorf("WATER.ISUNDER expects water, x, y, z")
	}
	h := heap.Handle(args[0].IVal)
	o, err := castW(m, h)
	if err != nil {
		return value.Nil, err
	}
	y, _ := args[2].ToFloat()
	bob := float32(math.Sin(float64(o.WaveT))) * o.WaveAmp * 0.15
	surfY := float64(o.PY + bob)
	return value.FromBool(y < surfY), nil
}

func (m *Module) wSetShallow(args []value.Value) (value.Value, error) {
	if len(args) != 5 {
		return value.Nil, fmt.Errorf("WATER.SETSHALLOWCOLOR expects water, r, g, b, a")
	}
	h := heap.Handle(args[0].IVal)
	o, err := castW(m, h)
	if err != nil {
		return value.Nil, err
	}
	ri, ok1 := args[1].ToInt()
	gi, ok2 := args[2].ToInt()
	bi, ok3 := args[3].ToInt()
	ai, ok4 := args[4].ToInt()
	if !ok1 || !ok2 || !ok3 || !ok4 {
		return value.Nil, fmt.Errorf("WATER.SETSHALLOWCOLOR: colors must be integers")
	}
	o.Shallow = rl.Color{R: uint8(ri), G: uint8(gi), B: uint8(bi), A: uint8(ai)}
	return args[0], nil
}

func (m *Module) wSetDeep(args []value.Value) (value.Value, error) {
	if len(args) != 5 {
		return value.Nil, fmt.Errorf("WATER.SETDEEPCOLOR expects water, r, g, b, a")
	}
	h := heap.Handle(args[0].IVal)
	o, err := castW(m, h)
	if err != nil {
		return value.Nil, err
	}
	ri, ok1 := args[1].ToInt()
	gi, ok2 := args[2].ToInt()
	bi, ok3 := args[3].ToInt()
	ai, ok4 := args[4].ToInt()
	if !ok1 || !ok2 || !ok3 || !ok4 {
		return value.Nil, fmt.Errorf("WATER.SETDEEPCOLOR: colors must be integers")
	}
	o.Deep = rl.Color{R: uint8(ri), G: uint8(gi), B: uint8(bi), A: uint8(ai)}
	return args[0], nil
}

// wSetColor sets shallow tint from packed RGB (bits 16..23, 8..15, 0..7) and scales both shallow/deep alpha from clarity (0..1 or 0..255).
func (m *Module) wSetColor(args []value.Value) (value.Value, error) {
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("WATER.SETCOLOR expects water, diffuse, clarity")
	}
	h := heap.Handle(args[0].IVal)
	o, err := castW(m, h)
	if err != nil {
		return value.Nil, err
	}
	di, ok1 := args[1].ToInt()
	cl, ok2 := args[2].ToFloat()
	if !ok1 || !ok2 {
		return value.Nil, fmt.Errorf("WATER.SETCOLOR: diffuse and clarity must be numeric")
	}
	r := uint8((di >> 16) & 0xff)
	g := uint8((di >> 8) & 0xff)
	b := uint8(di & 0xff)
	var a uint8
	if cl > 1 {
		if cl > 255 {
			a = 255
		} else {
			a = uint8(cl)
		}
	} else if cl < 0 {
		a = 0
	} else {
		a = uint8(cl * 255)
	}
	o.Shallow = rl.Color{R: r, G: g, B: b, A: a}
	// Deep: darker tint for depth gradient
	o.Deep = rl.Color{
		R: uint8(minInt(255, int(r)*2/3)),
		G: uint8(minInt(255, int(g)*2/3)),
		B: uint8(minInt(255, int(b)*2/3)),
		A: uint8(minInt(255, int(a)*13/10)),
	}
	return args[0], nil
}

func (m *Module) wGetPos(args []value.Value) (value.Value, error) {
	if m.h == nil || len(args) != 1 {
		return value.Nil, fmt.Errorf("WATER.GETPOS expects handle")
	}
	h := heap.Handle(args[0].IVal)
	o, err := castW(m, h)
	if err != nil {
		return value.Nil, err
	}
	return mbmatrix.AllocVec3Value(m.h, o.PX, o.PY, o.PZ)
}

func (m *Module) wSetRot(args []value.Value) (value.Value, error) {
	if len(args) < 2 {
		return value.Nil, fmt.Errorf("WATER.SETROT expects water, rot_y [or x,y,z]")
	}
	o, err := castW(m, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	if len(args) == 4 {
		rx, _ := args[1].ToFloat()
		ry, _ := args[2].ToFloat()
		rz, _ := args[3].ToFloat()
		o.RotX, o.RotY, o.RotZ = float32(rx), float32(ry), float32(rz)
	} else {
		ry, _ := args[1].ToFloat()
		o.RotY = float32(ry)
	}
	return args[0], nil
}

func (m *Module) wGetRot(args []value.Value) (value.Value, error) {
	o, err := castW(m, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	return mbmatrix.AllocVec3Value(m.h, o.RotX, o.RotY, o.RotZ)
}

func (m *Module) wSetScale(args []value.Value) (value.Value, error) {
	if len(args) < 2 {
		return value.Nil, fmt.Errorf("WATER.SETSCALE expects water, scale [or sx,sy,sz]")
	}
	o, err := castW(m, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	if len(args) == 4 {
		sx, _ := args[1].ToFloat()
		sy, _ := args[2].ToFloat()
		sz, _ := args[3].ToFloat()
		o.ScaleX, o.ScaleY, o.ScaleZ = float32(sx), float32(sy), float32(sz)
	} else {
		s, _ := args[1].ToFloat()
		o.ScaleX, o.ScaleY, o.ScaleZ = float32(s), float32(s), float32(s)
	}
	return args[0], nil
}

func (m *Module) wGetScale(args []value.Value) (value.Value, error) {
	o, err := castW(m, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	return mbmatrix.AllocVec3Value(m.h, o.ScaleX, o.ScaleY, o.ScaleZ)
}

func (m *Module) wGetColor(args []value.Value) (value.Value, error) {
	if m.h == nil || len(args) != 1 {
		return value.Nil, fmt.Errorf("WATER.GETCOLOR expects handle")
	}
	h := heap.Handle(args[0].IVal)
	o, err := castW(m, h)
	if err != nil {
		return value.Nil, err
	}
	arr, err := heap.NewArrayOfKind([]int64{4}, heap.ArrayKindFloat, 0)
	if err != nil {
		return value.Nil, err
	}
	arr.Floats[0] = float64(o.Shallow.R)
	arr.Floats[1] = float64(o.Shallow.G)
	arr.Floats[2] = float64(o.Shallow.B)
	arr.Floats[3] = float64(o.Shallow.A)
	id, err := m.h.Alloc(arr)
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}

func (m *Module) wGetWaveHeight(args []value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("WATER.GETWAVEHEIGHT expects handle")
	}
	h := heap.Handle(args[0].IVal)
	o, err := castW(m, h)
	if err != nil {
		return value.Nil, err
	}
	return value.FromFloat(float64(o.WaveAmp)), nil
}

func (m *Module) wGetWaveSpeed(args []value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("WATER.GETWAVESPEED expects handle")
	}
	h := heap.Handle(args[0].IVal)
	o, err := castW(m, h)
	if err != nil {
		return value.Nil, err
	}
	return value.FromFloat(float64(o.WaveFreq)), nil
}

func (m *Module) wGetShallowColor(args []value.Value) (value.Value, error) {
	if m.h == nil || len(args) != 1 {
		return value.Nil, fmt.Errorf("WATER.GETSHALLOWCOLOR expects handle")
	}
	o, err := castW(m, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	arr, err := heap.NewArrayOfKind([]int64{4}, heap.ArrayKindFloat, 0)
	if err != nil {
		return value.Nil, err
	}
	arr.Floats[0] = float64(o.Shallow.R)
	arr.Floats[1] = float64(o.Shallow.G)
	arr.Floats[2] = float64(o.Shallow.B)
	arr.Floats[3] = float64(o.Shallow.A)
	id, err := m.h.Alloc(arr)
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}

func (m *Module) wGetDeepColor(args []value.Value) (value.Value, error) {
	if m.h == nil || len(args) != 1 {
		return value.Nil, fmt.Errorf("WATER.GETDEEPCOLOR expects handle")
	}
	o, err := castW(m, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	arr, err := heap.NewArrayOfKind([]int64{4}, heap.ArrayKindFloat, 0)
	if err != nil {
		return value.Nil, err
	}
	arr.Floats[0] = float64(o.Deep.R)
	arr.Floats[1] = float64(o.Deep.G)
	arr.Floats[2] = float64(o.Deep.B)
	arr.Floats[3] = float64(o.Deep.A)
	id, err := m.h.Alloc(arr)
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
