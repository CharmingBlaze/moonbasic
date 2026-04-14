//go:build !cgo && !windows

package input

import (
	"fmt"
	"moonbasic/vm/value"
)

const cursorHint = "CURSOR.* natives require CGO: set CGO_ENABLED=1 and install a C compiler, then rebuild"

func (m *Module) inputLockMouse(args []value.Value) (value.Value, error) {
	return value.Nil, fmt.Errorf("INPUT.LOCKMOUSE: %s", cursorHint)
}

func (m *Module) cursorShow(args []value.Value) (value.Value, error) {
	return value.Nil, nil
}

func (m *Module) cursorHide(args []value.Value) (value.Value, error) {
	return value.Nil, nil
}

func (m *Module) cursorIsHidden(args []value.Value) (value.Value, error) {
	return value.False, nil
}

func (m *Module) cursorIsOnScreen(args []value.Value) (value.Value, error) {
	return value.False, nil
}

func (m *Module) cursorEnable(args []value.Value) (value.Value, error) {
	return value.Nil, nil
}

func (m *Module) cursorDisable(args []value.Value) (value.Value, error) {
	return value.Nil, nil
}

func (m *Module) cursorIsEnabled(args []value.Value) (value.Value, error) {
	return value.False, nil
}

func (m *Module) cursorSet(args []value.Value) (value.Value, error) {
	return value.Nil, nil
}
