//go:build cgo || (windows && !cgo)

package weathermod

import (
	"fmt"
	"math"

	"moonbasic/runtime/mbmodel3d"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

func syncFogGPU(m *Module) {
	mbmodel3d.SyncSceneFogWeather(m.FogOn, m.FogNear, m.FogFar, m.FogR, m.FogG, m.FogB)
}

// module implementation below

func (m *Module) wMake(args []value.Value) (value.Value, error) {
	if m.h == nil || len(args) != 0 {
		return value.Nil, fmt.Errorf("WEATHER.MAKE expects no args")
	}
	id, err := m.h.Alloc(&WeatherObject{Kind: "clear", Coverage: 0.2})
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(int32(id)), nil
}

func (m *Module) wFree(args []value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("WEATHER.FREE expects handle")
	}
	h := heap.Handle(args[0].IVal)
	m.h.Free(h)
	return value.Nil, nil
}

func (m *Module) wUpdate(args []value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("WEATHER.UPDATE expects weather, dt")
	}
	h := heap.Handle(args[0].IVal)
	dt, _ := args[1].ToFloat()
	o, err := heap.Cast[*WeatherObject](m.h, h)
	if err != nil {
		return value.Nil, err
	}
	o.Intensity = float32(math.Min(1, float64(o.Intensity)+dt*0.01))
	return args[0], nil
}

func (m *Module) wDraw(args []value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("WEATHER.DRAW expects weather")
	}
	return args[0], nil
}

func (m *Module) wSetType(args []value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("WEATHER.SETTYPE expects weather, type")
	}
	h := heap.Handle(args[0].IVal)
	s := args[1].String()
	o, err := heap.Cast[*WeatherObject](m.h, h)
	if err != nil {
		return value.Nil, err
	}
	o.Kind = s
	if s == "rain" || s == "storm" {
		o.Coverage = 0.7
	}
	return args[0], nil
}

func (m *Module) wGetCoverage(args []value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("WEATHER.GETCOVERAGE expects weather")
	}
	h := heap.Handle(args[0].IVal)
	o, err := heap.Cast[*WeatherObject](m.h, h)
	if err != nil {
		return value.Nil, err
	}
	return value.FromFloat(float64(o.Coverage)), nil
}

func (m *Module) wGetType(args []value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("WEATHER.GETTYPE expects weather")
	}
	h := heap.Handle(args[0].IVal)
	o, err := heap.Cast[*WeatherObject](m.h, h)
	if err != nil {
		return value.Nil, err
	}
	idx := m.requireHeap().Intern(o.Kind)
	return value.FromStringIndex(idx), nil
}

func (m *Module) fogEnable(args []value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("FOG.ENABLE expects bool")
	}
	m.FogOn = value.Truthy(args[0], nil, nil)
	syncFogGPU(m)
	return value.Nil, nil
}

func (m *Module) fogSetNear(args []value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("FOG.SETNEAR expects metres")
	}
	v, _ := args[0].ToFloat()
	m.FogNear = float32(v)
	syncFogGPU(m)
	return value.Nil, nil
}

func (m *Module) fogSetFar(args []value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("FOG.SETFAR expects metres")
	}
	v, _ := args[0].ToFloat()
	m.FogFar = float32(v)
	syncFogGPU(m)
	return value.Nil, nil
}

func (m *Module) fogSetRange(args []value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("FOG.SETRANGE expects (near, far)")
	}
	near, _ := args[0].ToFloat()
	far, _ := args[1].ToFloat()
	m.FogNear = float32(near)
	m.FogFar = float32(far)
	syncFogGPU(m)
	return value.Nil, nil
}

func (m *Module) fogSetColor(args []value.Value) (value.Value, error) {
	if len(args) != 4 {
		return value.Nil, fmt.Errorf("FOG.SETCOLOR expects r, g, b, a")
	}
	r, _ := args[0].ToInt()
	g, _ := args[1].ToInt()
	b, _ := args[2].ToInt()
	a, _ := args[3].ToInt()
	m.FogR, m.FogG, m.FogB, m.FogA = int(r), int(g), int(b), int(a)
	syncFogGPU(m)
	return value.Nil, nil
}

func (m *Module) windSet(args []value.Value) (value.Value, error) {
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("WIND.SET expects strength, dx, dz")
	}
	s, _ := args[0].ToFloat()
	dx, _ := args[1].ToFloat()
	dz, _ := args[2].ToFloat()
	m.WindStr = float32(s)
	m.WindDirX = float32(dx)
	m.WindDirZ = float32(dz)
	return value.Nil, nil
}

func (m *Module) windGetStrength(args []value.Value) (value.Value, error) {
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("WIND.GETSTRENGTH expects no args")
	}
	return value.FromFloat(float64(m.WindStr)), nil
}
