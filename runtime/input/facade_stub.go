//go:build !cgo && !windows

package input

import "moonbasic/vm/value"

func (m *Module) mouseDX(args []value.Value) (value.Value, error)       { return value.FromFloat(0), nil }
func (m *Module) mouseDY(args []value.Value) (value.Value, error)       { return value.FromFloat(0), nil }
func (m *Module) mouseWheel(args []value.Value) (value.Value, error)    { return value.FromFloat(0), nil }
func (m *Module) mouseDown(args []value.Value) (value.Value, error)     { return value.FromBool(false), nil }
func (m *Module) mousePressed(args []value.Value) (value.Value, error)  { return value.FromBool(false), nil }
func (m *Module) mouseReleased(args []value.Value) (value.Value, error) { return value.FromBool(false), nil }
func (m *Module) makeMouse(args []value.Value) (value.Value, error)     { return value.Nil, nil }
func (m *Module) mouseHitGlobal(args []value.Value) (value.Value, error) {
	return value.FromBool(false), nil
}
func (m *Module) mouseXGlobal(args []value.Value) (value.Value, error) { return value.FromInt(0), nil }
func (m *Module) mouseYGlobal(args []value.Value) (value.Value, error) { return value.FromInt(0), nil }

func (m *Module) keyDown(args []value.Value) (value.Value, error) { return value.FromBool(false), nil }
func (m *Module) keyHit(args []value.Value) (value.Value, error)  { return value.FromBool(false), nil }
func (m *Module) keyUp(args []value.Value) (value.Value, error)   { return value.FromBool(false), nil }
func (m *Module) makeKey(args []value.Value) (value.Value, error) { return value.Nil, nil }

func (m *Module) gpAxis(args []value.Value) (value.Value, error)   { return value.FromFloat(0), nil }
func (m *Module) gpButton(args []value.Value) (value.Value, error) { return value.FromBool(false), nil }
func (m *Module) makeGamepad(args []value.Value) (value.Value, error) {
	return value.Nil, nil
}
