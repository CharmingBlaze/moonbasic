//go:build !cgo && !windows

package water

import (
	"fmt"

	"moonbasic/vm/value"
)

func (m *Module) wMake(args []value.Value) (value.Value, error) {
	return value.Nil, fmt.Errorf("WATER.MAKE requires CGO")
}

func (m *Module) wCreate(args []value.Value) (value.Value, error) {
	return value.Nil, fmt.Errorf("WATER.CREATE requires CGO")
}

func (m *Module) wFree(args []value.Value) (value.Value, error) {
	return value.Nil, nil
}

func (m *Module) wSetPos(args []value.Value) (value.Value, error) {
	return value.Nil, nil
}

func (m *Module) wDraw(args []value.Value) (value.Value, error) {
	return value.Nil, nil
}

func (m *Module) wUpdate(args []value.Value) (value.Value, error) {
	return value.Nil, nil
}

func (m *Module) wSetWaveHeight(args []value.Value) (value.Value, error) {
	return value.Nil, nil
}

func (m *Module) wSetWave(args []value.Value) (value.Value, error) {
	return value.Nil, nil
}

func (m *Module) wGetWaveY(args []value.Value) (value.Value, error) {
	return value.FromFloat(0), nil
}

func (m *Module) wGetDepth(args []value.Value) (value.Value, error) {
	return value.FromFloat(0), nil
}

func (m *Module) wIsUnder(args []value.Value) (value.Value, error) {
	return value.FromBool(false), nil
}

func (m *Module) wSetShallow(args []value.Value) (value.Value, error) {
	return value.Nil, nil
}

func (m *Module) wSetDeep(args []value.Value) (value.Value, error) {
	return value.Nil, nil
}

func (m *Module) wSetColor(args []value.Value) (value.Value, error) {
	return value.Nil, nil
}
