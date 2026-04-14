//go:build !cgo && !windows

package input

import (
	"moonbasic/vm/value"
)

func (m *Module) inFlushMouse(args []value.Value) (value.Value, error) { return value.Nil, nil }
func (m *Module) inFlushKeys(args []value.Value) (value.Value, error)  { return value.Nil, nil }

func (m *Module) inWaitMouse(args []value.Value) (value.Value, error) { return value.FromInt(0), nil }
func (m *Module) inWaitKey(args []value.Value) (value.Value, error)   { return value.FromInt(0), nil }
func (m *Module) inGetKey(args []value.Value) (value.Value, error)    { return value.FromInt(0), nil }

func (m *Module) inMoveMouse(args []value.Value) (value.Value, error) { return value.Nil, nil }
func (m *Module) inHidePointer(args []value.Value) (value.Value, error) { return value.Nil, nil }
func (m *Module) inShowPointer(args []value.Value) (value.Value, error) { return value.Nil, nil }

func (m *Module) blitzMouseDown(rt interface{}, args ...value.Value) (value.Value, error) {
	return value.False, nil
}
func (m *Module) blitzMouseHit(rt interface{}, args ...value.Value) (value.Value, error) {
	return value.False, nil
}
func (m *Module) blitzKeyDown(rt interface{}, args ...value.Value) (value.Value, error) {
	return value.False, nil
}
func (m *Module) blitzKeyHit(rt interface{}, args ...value.Value) (value.Value, error) {
	return value.False, nil
}
func (m *Module) blitzGetJoy(rt interface{}, args ...value.Value) (value.Value, error) {
	return value.False, nil
}
func (m *Module) blitzJoyDown(rt interface{}, args ...value.Value) (value.Value, error) {
	return value.False, nil
}
func (m *Module) blitzJoyHit(rt interface{}, args ...value.Value) (value.Value, error) {
	return value.False, nil
}
func (m *Module) blitzJoyX(rt interface{}, args ...value.Value) (value.Value, error) {
	return value.FromFloat(0), nil
}
func (m *Module) blitzJoyY(rt interface{}, args ...value.Value) (value.Value, error) {
	return value.FromFloat(0), nil
}
func (m *Module) blitzJoyZ(rt interface{}, args ...value.Value) (value.Value, error) {
	return value.FromFloat(0), nil
}
func (m *Module) blitzJoyU(rt interface{}, args ...value.Value) (value.Value, error) {
	return value.FromFloat(0), nil
}
func (m *Module) blitzJoyV(rt interface{}, args ...value.Value) (value.Value, error) {
	return value.FromFloat(0), nil
}
func (m *Module) blitzJoyPitch(rt interface{}, args ...value.Value) (value.Value, error) {
	return value.FromFloat(0), nil
}
func (m *Module) blitzJoyYaw(rt interface{}, args ...value.Value) (value.Value, error) {
	return value.FromFloat(0), nil
}
func (m *Module) blitzJoyRoll(rt interface{}, args ...value.Value) (value.Value, error) {
	return value.FromFloat(0), nil
}
func (m *Module) blitzJoyHat(rt interface{}, args ...value.Value) (value.Value, error) {
	return value.FromInt(0), nil
}
func (m *Module) blitzJoyXDir(rt interface{}, args ...value.Value) (value.Value, error) {
	return value.FromInt(0), nil
}
func (m *Module) blitzJoyYDir(rt interface{}, args ...value.Value) (value.Value, error) {
	return value.FromInt(0), nil
}
