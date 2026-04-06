//go:build cgo || (windows && !cgo)

package window

import (
	"fmt"

	"moonbasic/runtime"
	"moonbasic/vm/value"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func (m *Module) registerWindowStateCommands(reg runtime.Registrar) {
	reg.Register("WINDOW.SETFLAG", "window", m.wSetFlag)
	reg.Register("WINDOW.CLEARFLAG", "window", m.wClearFlag)
	reg.Register("WINDOW.CHECKFLAG", "window", m.wCheckFlag)
	reg.Register("WINDOW.SETSTATE", "window", m.wSetFlag)
	reg.Register("WINDOW.SETMINSIZE", "window", m.wSetMinSize)
	reg.Register("WINDOW.SETMAXSIZE", "window", m.wSetMaxSize)
	reg.Register("WINDOW.GETPOSITIONX", "window", m.wGetPositionX)
	reg.Register("WINDOW.GETPOSITIONY", "window", m.wGetPositionY)
	reg.Register("WINDOW.SETMONITOR", "window", m.wSetMonitor)
	reg.Register("WINDOW.GETMONITORCOUNT", "window", m.wGetMonitorCount)
	reg.Register("WINDOW.GETMONITORNAME", "window", m.wGetMonitorName)
	reg.Register("WINDOW.GETMONITORWIDTH", "window", m.wGetMonitorWidth)
	reg.Register("WINDOW.GETMONITORHEIGHT", "window", m.wGetMonitorHeight)
	reg.Register("WINDOW.GETMONITORREFRESHRATE", "window", m.wGetMonitorRefreshRate)
	reg.Register("WINDOW.GETSCALEDPIX", "window", m.wGetScaleDPIX)
	reg.Register("WINDOW.GETSCALEDPIY", "window", m.wGetScaleDPIY)
	reg.Register("WINDOW.SETICON", "window", m.wSetIcon)
	reg.Register("WINDOW.SETOPACITY", "window", m.wSetOpacity)
}

func (m *Module) requireOpen(name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if !m.opened {
		return fmt.Errorf("%s: window is not open", name)
	}
	return nil
}

func (m *Module) wSetFlag(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("WINDOW.SETFLAG expects 1 argument (flag)")
	}
	f, ok := argInt(args[0])
	if !ok {
		return value.Nil, fmt.Errorf("WINDOW.SETFLAG: flag must be numeric")
	}
	if err := m.requireOpen("WINDOW.SETFLAG"); err != nil {
		return value.Nil, err
	}
	rl.SetWindowState(uint32(f))
	return value.Nil, nil
}

func (m *Module) wClearFlag(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("WINDOW.CLEARFLAG expects 1 argument (flag)")
	}
	f, ok := argInt(args[0])
	if !ok {
		return value.Nil, fmt.Errorf("WINDOW.CLEARFLAG: flag must be numeric")
	}
	if err := m.requireOpen("WINDOW.CLEARFLAG"); err != nil {
		return value.Nil, err
	}
	rl.ClearWindowState(uint32(f))
	return value.Nil, nil
}

func (m *Module) wCheckFlag(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("WINDOW.CHECKFLAG expects 1 argument (flag)")
	}
	f, ok := argInt(args[0])
	if !ok {
		return value.Nil, fmt.Errorf("WINDOW.CHECKFLAG: flag must be numeric")
	}
	if err := m.requireOpen("WINDOW.CHECKFLAG"); err != nil {
		return value.Nil, err
	}
	return value.FromBool(rl.IsWindowState(uint32(f))), nil
}

func (m *Module) wSetMinSize(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("WINDOW.SETMINSIZE expects (w, h)")
	}
	w, ok1 := argInt(args[0])
	h, ok2 := argInt(args[1])
	if !ok1 || !ok2 {
		return value.Nil, fmt.Errorf("WINDOW.SETMINSIZE: w and h must be numeric")
	}
	if err := m.requireOpen("WINDOW.SETMINSIZE"); err != nil {
		return value.Nil, err
	}
	rl.SetWindowMinSize(int(w), int(h))
	return value.Nil, nil
}

func (m *Module) wSetMaxSize(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("WINDOW.SETMAXSIZE expects (w, h)")
	}
	w, ok1 := argInt(args[0])
	h, ok2 := argInt(args[1])
	if !ok1 || !ok2 {
		return value.Nil, fmt.Errorf("WINDOW.SETMAXSIZE: w and h must be numeric")
	}
	if err := m.requireOpen("WINDOW.SETMAXSIZE"); err != nil {
		return value.Nil, err
	}
	rl.SetWindowMaxSize(int(w), int(h))
	return value.Nil, nil
}

func (m *Module) wGetPositionX(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("WINDOW.GETPOSITIONX expects 0 arguments")
	}
	if err := m.requireOpen("WINDOW.GETPOSITIONX"); err != nil {
		return value.Nil, err
	}
	p := rl.GetWindowPosition()
	return value.FromInt(int64(p.X)), nil
}

func (m *Module) wGetPositionY(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("WINDOW.GETPOSITIONY expects 0 arguments")
	}
	if err := m.requireOpen("WINDOW.GETPOSITIONY"); err != nil {
		return value.Nil, err
	}
	p := rl.GetWindowPosition()
	return value.FromInt(int64(p.Y)), nil
}

func (m *Module) wSetMonitor(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("WINDOW.SETMONITOR expects 1 argument (idx)")
	}
	idx, ok := argInt(args[0])
	if !ok {
		return value.Nil, fmt.Errorf("WINDOW.SETMONITOR: idx must be numeric")
	}
	if err := m.requireOpen("WINDOW.SETMONITOR"); err != nil {
		return value.Nil, err
	}
	rl.SetWindowMonitor(int(idx))
	return value.Nil, nil
}

func (m *Module) wGetMonitorCount(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("WINDOW.GETMONITORCOUNT expects 0 arguments")
	}
	return value.FromInt(int64(rl.GetMonitorCount())), nil
}

func (m *Module) wGetMonitorName(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("WINDOW.GETMONITORNAME expects 1 argument (idx)")
	}
	idx, ok := argInt(args[0])
	if !ok {
		return value.Nil, fmt.Errorf("WINDOW.GETMONITORNAME: idx must be numeric")
	}
	name := rl.GetMonitorName(int(idx))
	return rt.RetString(name), nil
}

func (m *Module) wGetMonitorWidth(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("WINDOW.GETMONITORWIDTH expects 1 argument (idx)")
	}
	idx, ok := argInt(args[0])
	if !ok {
		return value.Nil, fmt.Errorf("WINDOW.GETMONITORWIDTH: idx must be numeric")
	}
	return value.FromInt(int64(rl.GetMonitorWidth(int(idx)))), nil
}

func (m *Module) wGetMonitorHeight(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("WINDOW.GETMONITORHEIGHT expects 1 argument (idx)")
	}
	idx, ok := argInt(args[0])
	if !ok {
		return value.Nil, fmt.Errorf("WINDOW.GETMONITORHEIGHT: idx must be numeric")
	}
	return value.FromInt(int64(rl.GetMonitorHeight(int(idx)))), nil
}

func (m *Module) wGetMonitorRefreshRate(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("WINDOW.GETMONITORREFRESHRATE expects 1 argument (idx)")
	}
	idx, ok := argInt(args[0])
	if !ok {
		return value.Nil, fmt.Errorf("WINDOW.GETMONITORREFRESHRATE: idx must be numeric")
	}
	return value.FromInt(int64(rl.GetMonitorRefreshRate(int(idx)))), nil
}

func (m *Module) wGetScaleDPIX(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("WINDOW.GETSCALEDPIX expects 0 arguments")
	}
	if err := m.requireOpen("WINDOW.GETSCALEDPIX"); err != nil {
		return value.Nil, err
	}
	s := rl.GetWindowScaleDPI()
	return value.FromFloat(float64(s.X)), nil
}

func (m *Module) wGetScaleDPIY(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("WINDOW.GETSCALEDPIY expects 0 arguments")
	}
	if err := m.requireOpen("WINDOW.GETSCALEDPIY"); err != nil {
		return value.Nil, err
	}
	s := rl.GetWindowScaleDPI()
	return value.FromFloat(float64(s.Y)), nil
}

func (m *Module) wSetIcon(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Kind != value.KindString {
		return value.Nil, fmt.Errorf("WINDOW.SETICON expects 1 argument (image path$)")
	}
	if err := m.requireOpen("WINDOW.SETICON"); err != nil {
		return value.Nil, err
	}
	path, err := rt.ArgString(args, 0)
	if err != nil {
		return value.Nil, err
	}
	img := rl.LoadImage(path)
	if img == nil {
		return value.Nil, runtime.Errorf("WINDOW.SETICON: failed to load %q", path)
	}
	defer rl.UnloadImage(img)
	rl.SetWindowIcon(*img)
	return value.Nil, nil
}

func (m *Module) wSetOpacity(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("WINDOW.SETOPACITY expects 1 argument (alpha#)")
	}
	var a float64
	var ok bool
	if args[0].Kind == value.KindFloat {
		a, ok = args[0].ToFloat()
	} else if i, o2 := args[0].ToInt(); o2 {
		a, ok = float64(i), true
	}
	if !ok {
		return value.Nil, fmt.Errorf("WINDOW.SETOPACITY: alpha must be numeric")
	}
	if err := m.requireOpen("WINDOW.SETOPACITY"); err != nil {
		return value.Nil, err
	}
	rl.SetWindowOpacity(float32(a))
	return value.Nil, nil
}
