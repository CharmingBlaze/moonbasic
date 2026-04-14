//go:build !cgo && !windows

package input

import (
	"fmt"
	"moonbasic/vm/value"
)

const mouseExtraHint = "INPUT mouse/touch/gamepad requires CGO (Raylib)"

func (m *Module) inMousePressed(args []value.Value) (value.Value, error) {
	return value.False, nil
}

func (m *Module) inMouseReleased(args []value.Value) (value.Value, error) {
	return value.False, nil
}

func (m *Module) inMouseWheelMove(args []value.Value) (value.Value, error) {
	return value.FromFloat(0), nil
}

func (m *Module) inMouseDeltaX(args []value.Value) (value.Value, error) {
	return value.FromFloat(0), nil
}

func (m *Module) inMouseDeltaY(args []value.Value) (value.Value, error) {
	return value.FromFloat(0), nil
}

func (m *Module) inSetMousePos(args []value.Value) (value.Value, error) {
	return value.Nil, nil
}

func (m *Module) inCharPressed(args []value.Value) (value.Value, error) {
	return value.FromInt(0), nil
}

func (m *Module) inIsGamepadAvailable(args []value.Value) (value.Value, error) {
	return value.False, nil
}

func (m *Module) inGetGamepadAxisValue(args []value.Value) (value.Value, error) {
	return value.FromFloat(0), nil
}

func (m *Module) inMouseX(args []value.Value) (value.Value, error) { return value.FromInt(0), nil }
func (m *Module) inMouseY(args []value.Value) (value.Value, error) { return value.FromInt(0), nil }
func (m *Module) inMouseDown(args []value.Value) (value.Value, error) { return value.False, nil }
func (m *Module) inMouseHit(args []value.Value) (value.Value, error)  { return value.False, nil }
func (m *Module) inSetMouseScale(args []value.Value) (value.Value, error) { return value.Nil, nil }
func (m *Module) inSetMouseOffset(args []value.Value) (value.Value, error) { return value.Nil, nil }
func (m *Module) inGetMouseWorldPos(args []value.Value) (value.Value, error) {
	return value.Nil, fmt.Errorf("INPUT.GETMOUSEWORLDPOS: %s", mouseExtraHint)
}
func (m *Module) inTouchCount(args []value.Value) (value.Value, error)     { return value.FromInt(0), nil }
func (m *Module) inTouchX(args []value.Value) (value.Value, error)         { return value.FromInt(0), nil }
func (m *Module) inTouchY(args []value.Value) (value.Value, error)         { return value.FromInt(0), nil }
func (m *Module) inTouchPressed(args []value.Value) (value.Value, error)   { return value.False, nil }
func (m *Module) inGetTouchPointID(args []value.Value) (value.Value, error) { return value.FromInt(0), nil }
func (m *Module) inGamepadButtonCount(args []value.Value) (value.Value, error) { return value.FromInt(0), nil }
func (m *Module) inGamepadAxisCount(args []value.Value) (value.Value, error)   { return value.FromInt(0), nil }
func (m *Module) inSetGamepadMappings(rt interface{}, args ...value.Value) (value.Value, error) {
	return value.FromInt(0), nil
}
