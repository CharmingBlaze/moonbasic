//go:build !cgo && !windows

package input

import "moonbasic/vm/value"

const hint = "INPUT.* requires CGO: set CGO_ENABLED=1 and install a C compiler, then rebuild"

func (m *Module) Shutdown() {}

func (m *Module) inKeyDown(args []value.Value) (value.Value, error) {
	return value.FromBool(false), nil
}

func (m *Module) inKeyPressed(args []value.Value) (value.Value, error) {
	return value.FromBool(false), nil
}

func (m *Module) inKeyReleased(args []value.Value) (value.Value, error) {
	return value.FromBool(false), nil
}

func (m *Module) inGetKeyName(args []value.Value) (value.Value, error) {
	return value.Nil, nil
}

func (m *Module) inGetInactivity(args []value.Value) (value.Value, error) {
	return value.FromFloat(0), nil
}
