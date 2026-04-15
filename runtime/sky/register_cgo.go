//go:build cgo || (windows && !cgo)

package sky

import (
	"fmt"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/runtime"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

func registerSky(m *Module, r runtime.Registrar) {
	r.Register("SKY.MAKE", "sky", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return skyMake(m, rt, "SKY.MAKE", args...) })
	r.Register("SKY.CREATE", "sky", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return skyMake(m, rt, "SKY.CREATE", args...) })
	r.Register("SKY.FREE", "sky", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return skyFree(m, rt, args...) })
	r.Register("SKY.UPDATE", "sky", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return skyUpdate(m, rt, args...) })
	r.Register("SKY.DRAW", "sky", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return skyDraw(m, rt, args...) })
	r.Register("SKY.SETTIME", "sky", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return skySetTime(m, rt, args...) })
	r.Register("SKY.SETDAYLENGTH", "sky", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return skySetDayLength(m, rt, args...) })
	r.Register("SKY.GETTIMEHOURS", "sky", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return skyGetTimeHours(m, rt, args...) })
	r.Register("SKY.ISNIGHT", "sky", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return skyIsNight(m, rt, args...) })
}

func castSky(m *Module, h heap.Handle) (*SkyObject, error) {
	return heap.Cast[*SkyObject](m.h, h)
}

func skyColor(o *SkyObject) rl.Color {
	t := float64(o.Time)
	day := 0.5 + 0.5*math.Sin(t*2*math.Pi)
	r := uint8(20 + day*100)
	g := uint8(30 + day*140)
	b := uint8(60 + day*160)
	return rl.Color{R: r, G: g, B: b, A: 255}
}

func skyMake(m *Module, rt *runtime.Runtime, op string, args ...value.Value) (value.Value, error) {
	if m.h == nil || len(args) != 0 {
		return value.Nil, fmt.Errorf("%s expects no arguments", op)
	}
	o := &SkyObject{Time: 0.45}
	id, err := m.h.Alloc(o)
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(int32(id)), nil
}

func skyFree(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if m.h == nil || len(args) != 1 {
		return value.Nil, fmt.Errorf("SKY.FREE expects handle")
	}
	h, err := rt.ArgHandle(args, 0)
	if err != nil {
		return value.Nil, err
	}
	m.h.Free(heap.Handle(h))
	return value.Nil, nil
}

func skyUpdate(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("SKY.UPDATE expects sky, dt")
	}
	h, err := rt.ArgHandle(args, 0)
	if err != nil {
		return value.Nil, err
	}
	dt, err := rt.ArgFloat(args, 1)
	if err != nil {
		return value.Nil, err
	}
	o, err := castSky(m, heap.Handle(h))
	if err != nil {
		return value.Nil, err
	}
	if o.DayLength > 0 {
		o.Time += float32(dt) / o.DayLength
		for o.Time >= 1 {
			o.Time--
		}
	}
	return value.Nil, nil
}

func skyDraw(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("SKY.DRAW expects sky")
	}
	h, err := rt.ArgHandle(args, 0)
	if err != nil {
		return value.Nil, err
	}
	o, err := castSky(m, heap.Handle(h))
	if err != nil {
		return value.Nil, err
	}
	col := skyColor(o)
	rl.DisableBackfaceCulling()
	rl.DisableDepthTest()
	rl.DrawSphere(rl.Vector3{}, 220, col)
	rl.EnableDepthTest()
	rl.EnableBackfaceCulling()
	return value.Nil, nil
}

func skySetTime(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("SKY.SETTIME expects sky, t")
	}
	h, err := rt.ArgHandle(args, 0)
	if err != nil {
		return value.Nil, err
	}
	t, err := rt.ArgFloat(args, 1)
	if err != nil {
		return value.Nil, err
	}
	o, err := castSky(m, heap.Handle(h))
	if err != nil {
		return value.Nil, err
	}
	o.Time = float32(t)
	return value.Nil, nil
}

func skySetDayLength(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("SKY.SETDAYLENGTH expects sky, seconds")
	}
	h, err := rt.ArgHandle(args, 0)
	if err != nil {
		return value.Nil, err
	}
	sec, err := rt.ArgFloat(args, 1)
	if err != nil {
		return value.Nil, err
	}
	o, err := castSky(m, heap.Handle(h))
	if err != nil {
		return value.Nil, err
	}
	o.DayLength = float32(sec)
	return value.Nil, nil
}

func skyGetTimeHours(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("SKY.GETTIMEHOURS expects sky")
	}
	h, err := rt.ArgHandle(args, 0)
	if err != nil {
		return value.Nil, err
	}
	o, err := castSky(m, heap.Handle(h))
	if err != nil {
		return value.Nil, err
	}
	return value.FromFloat(float64(o.Time * 24)), nil
}

func skyIsNight(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("SKY.ISNIGHT expects sky")
	}
	h, err := rt.ArgHandle(args, 0)
	if err != nil {
		return value.Nil, err
	}
	o, err := castSky(m, heap.Handle(h))
	if err != nil {
		return value.Nil, err
	}
	return value.FromBool(o.Time < 0.25 || o.Time > 0.75), nil
}
