//go:build !cgo && !windows

package weathermod

import (
	"fmt"

	"moonbasic/vm/value"
)

func (m *Module) wMake(args []value.Value) (value.Value, error) {
	return value.Nil, fmt.Errorf("WEATHER.MAKE requires CGO")
}

func (m *Module) wFree(args []value.Value) (value.Value, error) {
	return value.Nil, nil
}

func (m *Module) wUpdate(args []value.Value) (value.Value, error) {
	return value.Nil, nil
}

func (m *Module) wDraw(args []value.Value) (value.Value, error) {
	return value.Nil, nil
}

func (m *Module) wSetType(args []value.Value) (value.Value, error) {
	return value.Nil, nil
}

func (m *Module) wGetCoverage(args []value.Value) (value.Value, error) {
	return value.FromFloat(0), nil
}

func (m *Module) wGetType(args []value.Value) (value.Value, error) {
	return value.FromInt(0), nil
}

func (m *Module) fogEnable(args []value.Value) (value.Value, error) {
	return value.Nil, nil
}

func (m *Module) fogSetNear(args []value.Value) (value.Value, error) {
	return value.Nil, nil
}

func (m *Module) fogSetFar(args []value.Value) (value.Value, error) {
	return value.Nil, nil
}

func (m *Module) fogSetRange(args []value.Value) (value.Value, error) {
	return value.Nil, nil
}

func (m *Module) fogSetColor(args []value.Value) (value.Value, error) {
	return value.Nil, nil
}

func (m *Module) windSet(args []value.Value) (value.Value, error) {
	return value.Nil, nil
}

func (m *Module) windGetStrength(args []value.Value) (value.Value, error) {
	return value.FromFloat(0), nil
}
