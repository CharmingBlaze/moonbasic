//go:build !cgo && !windows

package input

import (
	"moonbasic/vm/value"
)

func (m *Module) inSetGamepadMappings(args []value.Value) (value.Value, error) {
	return value.FromInt(0), nil
}
