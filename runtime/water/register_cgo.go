//go:build cgo

package water

import (
	"fmt"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/runtime"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

func registerWater(m *Module, r runtime.Registrar) {
	r.Register("WATER.MAKE", "water", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return wMake(m, rt, args...) })
	r.Register("WATER.FREE", "water", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return wFree(m, rt, args...) })
	r.Register("WATER.SETPOS", "water", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return wSetPos(m, rt, args...) })
	r.Register("WATER.DRAW", "water", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return wDraw(m, rt, args...) })
	r.Register("WATER.UPDATE", "water", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return wUpdate(m, rt, args...) })
	r.Register("WATER.SETWAVEHEIGHT", "water", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return wSetWaveHeight(m, rt, args...) })
	r.Register("WATER.GETWAVEY", "water", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return wGetWaveY(m, rt, args...) })
	r.Register("WATER.GETDEPTH", "water", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return wGetDepth(m, rt, args...) })
	r.Register("WATER.ISUNDER", "water", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return wIsUnder(m, rt, args...) })
	r.Register("WATER.SETSHALLOWCOLOR", "water", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return wSetShallow(m, rt, args...) })
	r.Register("WATER.SETDEEPCOLOR", "water", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return wSetDeep(m, rt, args...) })
}

func castW(m *Module, h heap.Handle) (*WaterObject, error) {
	return heap.Cast[*WaterObject](m.h, h)
}

func wMake(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if m.h == nil || len(args) != 2 {
		return value.Nil, fmt.Errorf("WATER.MAKE expects width#, depth#")
	}
	wf, err := rt.ArgFloat(args, 0)
	if err != nil {
		return value.Nil, err
	}
	df, err := rt.ArgFloat(args, 1)
	if err != nil {
		return value.Nil, err
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
	}
	id, err := m.h.Alloc(o)
	if err != nil {
		return value.Nil, err
	}
	m.waters = append(m.waters, id)
	return value.FromHandle(int32(id)), nil
}

func wFree(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if m.h == nil || len(args) != 1 {
		return value.Nil, fmt.Errorf("WATER.FREE expects handle")
	}
	h, err := rt.ArgHandle(args, 0)
	if err != nil {
		return value.Nil, err
	}
	hh := heap.Handle(h)
	for i, id := range m.waters {
		if id == hh {
			m.waters = append(m.waters[:i], m.waters[i+1:]...)
			break
		}
	}
	m.h.Free(hh)
	return value.Nil, nil
}

func wSetPos(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 4 {
		return value.Nil, fmt.Errorf("WATER.SETPOS expects water, x#, y#, z#")
	}
	h, err := rt.ArgHandle(args, 0)
	if err != nil {
		return value.Nil, err
	}
	o, err := castW(m, heap.Handle(h))
	if err != nil {
		return value.Nil, err
	}
	x, err := rt.ArgFloat(args, 1)
	if err != nil {
		return value.Nil, err
	}
	y, err := rt.ArgFloat(args, 2)
	if err != nil {
		return value.Nil, err
	}
	z, err := rt.ArgFloat(args, 3)
	if err != nil {
		return value.Nil, err
	}
	o.PX = float32(x)
	o.PY = float32(y)
	o.PZ = float32(z)
	o.BedY = o.PY - 12
	return value.Nil, nil
}

func wDraw(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("WATER.DRAW expects handle")
	}
	h, err := rt.ArgHandle(args, 0)
	if err != nil {
		return value.Nil, err
	}
	o, err := castW(m, heap.Handle(h))
	if err != nil {
		return value.Nil, err
	}
	bob := float32(math.Sin(float64(o.WaveT))) * o.WaveAmp * 0.15
	mat := rl.MatrixTranslate(o.PX, o.PY+bob, o.PZ)
	rl.DrawMesh(o.Mesh, o.Mat, mat)
	return value.Nil, nil
}

func wUpdate(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("WATER.UPDATE expects dt#")
	}
	dt, err := rt.ArgFloat(args, 0)
	if err != nil {
		return value.Nil, err
	}
	for _, id := range m.waters {
		o, err := castW(m, id)
		if err != nil {
			continue
		}
		o.WaveT += float32(dt) * o.WaveFreq
	}
	return value.Nil, nil
}

func wSetWaveHeight(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("WATER.SETWAVEHEIGHT expects water, height#")
	}
	h, err := rt.ArgHandle(args, 0)
	if err != nil {
		return value.Nil, err
	}
	o, err := castW(m, heap.Handle(h))
	if err != nil {
		return value.Nil, err
	}
	v, err := rt.ArgFloat(args, 1)
	if err != nil {
		return value.Nil, err
	}
	o.WaveAmp = float32(v)
	return value.Nil, nil
}

func wGetWaveY(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("WATER.GETWAVEY expects water, x#, z#")
	}
	h, err := rt.ArgHandle(args, 0)
	if err != nil {
		return value.Nil, err
	}
	o, err := castW(m, heap.Handle(h))
	if err != nil {
		return value.Nil, err
	}
	_, err = rt.ArgFloat(args, 1)
	if err != nil {
		return value.Nil, err
	}
	_, err = rt.ArgFloat(args, 2)
	if err != nil {
		return value.Nil, err
	}
	bob := float32(math.Sin(float64(o.WaveT))) * o.WaveAmp * 0.15
	return value.FromFloat(float64(o.PY + bob)), nil
}

func wGetDepth(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("WATER.GETDEPTH expects water, x#, z#")
	}
	h, err := rt.ArgHandle(args, 0)
	if err != nil {
		return value.Nil, err
	}
	o, err := castW(m, heap.Handle(h))
	if err != nil {
		return value.Nil, err
	}
	if _, err = rt.ArgFloat(args, 1); err != nil {
		return value.Nil, err
	}
	if _, err = rt.ArgFloat(args, 2); err != nil {
		return value.Nil, err
	}
	bob := float32(math.Sin(float64(o.WaveT))) * o.WaveAmp * 0.15
	surfY := float64(o.PY + bob)
	return value.FromFloat(surfY - float64(o.BedY)), nil
}

func wIsUnder(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 4 {
		return value.Nil, fmt.Errorf("WATER.ISUNDER expects water, x#, y#, z#")
	}
	h, err := rt.ArgHandle(args, 0)
	if err != nil {
		return value.Nil, err
	}
	o, err := castW(m, heap.Handle(h))
	if err != nil {
		return value.Nil, err
	}
	_, err = rt.ArgFloat(args, 1)
	if err != nil {
		return value.Nil, err
	}
	y, err := rt.ArgFloat(args, 2)
	if err != nil {
		return value.Nil, err
	}
	_, err = rt.ArgFloat(args, 3)
	if err != nil {
		return value.Nil, err
	}
	bob := float32(math.Sin(float64(o.WaveT))) * o.WaveAmp * 0.15
	surfY := float64(o.PY + bob)
	return value.FromBool(y < surfY), nil
}

func wSetShallow(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 5 {
		return value.Nil, fmt.Errorf("WATER.SETSHALLOWCOLOR expects water, r, g, b, a")
	}
	h, err := rt.ArgHandle(args, 0)
	if err != nil {
		return value.Nil, err
	}
	o, err := castW(m, heap.Handle(h))
	if err != nil {
		return value.Nil, err
	}
	ri, _ := rt.ArgInt(args, 1)
	gi, _ := rt.ArgInt(args, 2)
	bi, _ := rt.ArgInt(args, 3)
	ai, _ := rt.ArgInt(args, 4)
	o.Shallow = rl.Color{R: uint8(ri), G: uint8(gi), B: uint8(bi), A: uint8(ai)}
	return value.Nil, nil
}

func wSetDeep(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 5 {
		return value.Nil, fmt.Errorf("WATER.SETDEEPCOLOR expects water, r, g, b, a")
	}
	h, err := rt.ArgHandle(args, 0)
	if err != nil {
		return value.Nil, err
	}
	o, err := castW(m, heap.Handle(h))
	if err != nil {
		return value.Nil, err
	}
	ri, _ := rt.ArgInt(args, 1)
	gi, _ := rt.ArgInt(args, 2)
	bi, _ := rt.ArgInt(args, 3)
	ai, _ := rt.ArgInt(args, 4)
	o.Deep = rl.Color{R: uint8(ri), G: uint8(gi), B: uint8(bi), A: uint8(ai)}
	return value.Nil, nil
}
