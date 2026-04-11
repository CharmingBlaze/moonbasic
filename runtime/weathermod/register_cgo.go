//go:build cgo || (windows && !cgo)

package weathermod

import (
	"fmt"
	"math"

	"moonbasic/runtime"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

func registerWeather(m *Module, r runtime.Registrar) {
	r.Register("WEATHER.MAKE", "weather", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return wMake(m, rt, args...) })
	r.Register("WEATHER.FREE", "weather", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return wFree(m, rt, args...) })
	r.Register("WEATHER.UPDATE", "weather", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return wUpdate(m, rt, args...) })
	r.Register("WEATHER.DRAW", "weather", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return wDraw(m, rt, args...) })
	r.Register("WEATHER.SETTYPE", "weather", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return wSetType(m, rt, args...) })
	r.Register("WEATHER.GETCOVERAGE", "weather", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return wGetCoverage(m, rt, args...) })
	r.Register("WEATHER.GETTYPE", "weather", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return wGetType(m, rt, args...) })

	r.Register("FOG.ENABLE", "fog", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return fogEnable(m, rt, args...) })
	r.Register("FOG.SETNEAR", "fog", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return fogSetNear(m, rt, args...) })
	r.Register("FOG.SETFAR", "fog", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return fogSetFar(m, rt, args...) })
	r.Register("FOG.SETCOLOR", "fog", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return fogSetColor(m, rt, args...) })
	r.Register("FOG.SETRANGE", "fog", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return fogSetRange(m, rt, args...) })

	r.Register("WIND.SET", "wind", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return windSet(m, rt, args...) })
	r.Register("WIND.GETSTRENGTH", "wind", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return windGetStrength(m, rt, args...) })
}

func wMake(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if m.h == nil || len(args) != 0 {
		return value.Nil, fmt.Errorf("WEATHER.MAKE expects no args")
	}
	id, err := m.h.Alloc(&WeatherObject{Kind: "clear", Coverage: 0.2})
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(int32(id)), nil
}

func wFree(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("WEATHER.FREE expects handle")
	}
	h, err := rt.ArgHandle(args, 0)
	if err != nil {
		return value.Nil, err
	}
	m.h.Free(heap.Handle(h))
	return value.Nil, nil
}

func wUpdate(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("WEATHER.UPDATE expects weather, dt")
	}
	h, err := rt.ArgHandle(args, 0)
	if err != nil {
		return value.Nil, err
	}
	dt, err := rt.ArgFloat(args, 1)
	if err != nil {
		return value.Nil, err
	}
	o, err := heap.Cast[*WeatherObject](m.h, heap.Handle(h))
	if err != nil {
		return value.Nil, err
	}
	o.Intensity = float32(math.Min(1, float64(o.Intensity)+dt*0.01))
	return value.Nil, nil
}

func wDraw(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("WEATHER.DRAW expects weather")
	}
	_, _ = rt.ArgHandle(args, 0)
	return value.Nil, nil
}

func wSetType(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("WEATHER.SETTYPE expects weather, type")
	}
	h, err := rt.ArgHandle(args, 0)
	if err != nil {
		return value.Nil, err
	}
	s, err := rt.ArgString(args, 1)
	if err != nil {
		return value.Nil, err
	}
	o, err := heap.Cast[*WeatherObject](m.h, heap.Handle(h))
	if err != nil {
		return value.Nil, err
	}
	o.Kind = s
	if s == "rain" || s == "storm" {
		o.Coverage = 0.7
	}
	return value.Nil, nil
}

func wGetCoverage(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("WEATHER.GETCOVERAGE expects weather")
	}
	h, err := rt.ArgHandle(args, 0)
	if err != nil {
		return value.Nil, err
	}
	o, err := heap.Cast[*WeatherObject](m.h, heap.Handle(h))
	if err != nil {
		return value.Nil, err
	}
	return value.FromFloat(float64(o.Coverage)), nil
}

func wGetType(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("WEATHER.GETTYPE expects weather")
	}
	h, err := rt.ArgHandle(args, 0)
	if err != nil {
		return value.Nil, err
	}
	o, err := heap.Cast[*WeatherObject](m.h, heap.Handle(h))
	if err != nil {
		return value.Nil, err
	}
	return rt.RetString(o.Kind), nil
}

func fogEnable(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("FOG.ENABLE expects bool")
	}
	b, err := rt.ArgBool(args, 0)
	if err != nil {
		return value.Nil, err
	}
	m.FogOn = b
	return value.Nil, nil
}

func fogSetNear(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("FOG.SETNEAR expects metres")
	}
	v, err := rt.ArgFloat(args, 0)
	if err != nil {
		return value.Nil, err
	}
	m.FogNear = float32(v)
	return value.Nil, nil
}

func fogSetFar(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("FOG.SETFAR expects metres")
	}
	v, err := rt.ArgFloat(args, 0)
	if err != nil {
		return value.Nil, err
	}
	m.FogFar = float32(v)
	return value.Nil, nil
}

func fogSetRange(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("FOG.SETRANGE expects (near, far)")
	}
	near, err := rt.ArgFloat(args, 0)
	if err != nil {
		return value.Nil, err
	}
	far, err := rt.ArgFloat(args, 1)
	if err != nil {
		return value.Nil, err
	}
	m.FogNear = float32(near)
	m.FogFar = float32(far)
	return value.Nil, nil
}

func fogSetColor(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 4 {
		return value.Nil, fmt.Errorf("FOG.SETCOLOR expects r, g, b, a")
	}
	r_, _ := rt.ArgInt(args, 0)
	g_, _ := rt.ArgInt(args, 1)
	b_, _ := rt.ArgInt(args, 2)
	a_, _ := rt.ArgInt(args, 3)
	m.FogR, m.FogG, m.FogB, m.FogA = int(r_), int(g_), int(b_), int(a_)
	_ = m.FogOn
	return value.Nil, nil
}

func windSet(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("WIND.SET expects strength, dx, dz")
	}
	s, err := rt.ArgFloat(args, 0)
	if err != nil {
		return value.Nil, err
	}
	dx, err := rt.ArgFloat(args, 1)
	if err != nil {
		return value.Nil, err
	}
	dz, err := rt.ArgFloat(args, 2)
	if err != nil {
		return value.Nil, err
	}
	m.WindStr = float32(s)
	m.WindDirX = float32(dx)
	m.WindDirZ = float32(dz)
	return value.Nil, nil
}

func windGetStrength(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("WIND.GETSTRENGTH expects no args")
	}
	return value.FromFloat(float64(m.WindStr)), nil
}
