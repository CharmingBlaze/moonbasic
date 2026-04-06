//go:build cgo || (windows && !cgo)

package window

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

func (m *Module) registerWindowPlacementCommands(reg runtime.Registrar) {
	reg.Register("WINDOW.SETPOSITION", "window", runtime.AdaptLegacy(m.wSetPosition))
	reg.Register("WINDOW.SETSIZE", "window", runtime.AdaptLegacy(m.wSetSize))
	reg.Register("WINDOW.MINIMIZE", "window", runtime.AdaptLegacy(m.wMinimize))
	reg.Register("WINDOW.MAXIMIZE", "window", runtime.AdaptLegacy(m.wMaximize))
	reg.Register("WINDOW.RESTORE", "window", runtime.AdaptLegacy(m.wRestore))
	reg.Register("WINDOW.SETTARGETFPS", "window", m.wSetFPS)
}

func (m *Module) wSetPosition(args []value.Value) (value.Value, error) {
	if err := m.requireOpen("WINDOW.SETPOSITION"); err != nil {
		return value.Nil, err
	}
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("WINDOW.SETPOSITION expects (x, y)")
	}
	x, okx := argInt(args[0])
	y, oky := argInt(args[1])
	if !okx || !oky {
		return value.Nil, fmt.Errorf("WINDOW.SETPOSITION: x and y must be numeric")
	}
	rl.SetWindowPosition(int(x), int(y))
	return value.Nil, nil
}

func (m *Module) wSetSize(args []value.Value) (value.Value, error) {
	if err := m.requireOpen("WINDOW.SETSIZE"); err != nil {
		return value.Nil, err
	}
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("WINDOW.SETSIZE expects (width, height)")
	}
	w, okw := argInt(args[0])
	h, okh := argInt(args[1])
	if !okw || !okh || w < 1 || h < 1 {
		return value.Nil, fmt.Errorf("WINDOW.SETSIZE: width and height must be positive numerics")
	}
	rl.SetWindowSize(int(w), int(h))
	return value.Nil, nil
}

func (m *Module) wMinimize(args []value.Value) (value.Value, error) {
	if err := m.requireOpen("WINDOW.MINIMIZE"); err != nil {
		return value.Nil, err
	}
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("WINDOW.MINIMIZE expects 0 arguments")
	}
	rl.MinimizeWindow()
	return value.Nil, nil
}

func (m *Module) wMaximize(args []value.Value) (value.Value, error) {
	if err := m.requireOpen("WINDOW.MAXIMIZE"); err != nil {
		return value.Nil, err
	}
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("WINDOW.MAXIMIZE expects 0 arguments")
	}
	rl.MaximizeWindow()
	return value.Nil, nil
}

func (m *Module) wRestore(args []value.Value) (value.Value, error) {
	if err := m.requireOpen("WINDOW.RESTORE"); err != nil {
		return value.Nil, err
	}
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("WINDOW.RESTORE expects 0 arguments")
	}
	rl.RestoreWindow()
	return value.Nil, nil
}
