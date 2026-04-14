//go:build !cgo && !windows

package input

import (
	"moonbasic/runtime"
	"moonbasic/vm/value"
)

func (m *Module) inKeyHit(args []value.Value) (value.Value, error) { return value.False, nil }

func (m *Module) inMouseXSpeed(args []value.Value) (value.Value, error) {
	return value.FromFloat(0), nil
}

func (m *Module) inMouseYSpeed(args []value.Value) (value.Value, error) {
	return value.FromFloat(0), nil
}

func (m *Module) inJoyX(args []value.Value) (value.Value, error) { return value.FromFloat(0), nil }
func (m *Module) inJoyY(args []value.Value) (value.Value, error) { return value.FromFloat(0), nil }

func (m *Module) inJoyButton(args []value.Value) (value.Value, error) {
	return value.False, nil
}

func (m *Module) inMouseWheel(args []value.Value) (value.Value, error) {
	return value.FromFloat(0), nil
}

func (m *Module) inKeyUp(args []value.Value) (value.Value, error) { return value.False, nil }

func (m *Module) inAxis(args []value.Value) (value.Value, error) { return value.FromFloat(0), nil }

func (m *Module) inputMouseDelta(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	return allocInputTuple2(nil, 0, 0)
}

func (m *Module) inputMoveDir(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	return allocInputTuple2(nil, 0, 0)
}
