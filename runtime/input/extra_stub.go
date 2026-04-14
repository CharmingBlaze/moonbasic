//go:build !cgo && !windows

package input

import (
	"fmt"
	"moonbasic/runtime"
	"moonbasic/vm/value"
)

// extra input methods (stubs) - most are now in mouse_extra_stub.go to avoid duplication
// but we keep the ones registered as AdaptLegacy if needed.

func (m *Module) inMapKey(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	return value.Nil, nil
}
func (m *Module) inMapGamepadButton(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	return value.Nil, nil
}
func (m *Module) inMapGamepadAxis(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	return value.Nil, nil
}
func (m *Module) inActionPressed(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	return value.False, nil
}
func (m *Module) inActionDown(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	return value.False, nil
}
func (m *Module) inActionReleased(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	return value.False, nil
}
func (m *Module) inActionAxis(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	return value.FromFloat(0), nil
}
func (m *Module) inSaveMappings(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	return value.Nil, nil
}
func (m *Module) inLoadMappings(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	return value.Nil, nil
}

func (m *Module) inSetGamepadMappings(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	return value.FromInt(0), nil
}
