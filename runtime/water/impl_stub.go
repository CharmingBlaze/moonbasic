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

func (m *Module) wGetPos(args []value.Value) (value.Value, error) {
	return value.Nil, fmt.Errorf("WATER.GETPOS requires CGO")
}

func (m *Module) wSetRot(args []value.Value) (value.Value, error) {
	return value.Nil, fmt.Errorf("WATER.SETROT requires CGO")
}

func (m *Module) wGetRot(args []value.Value) (value.Value, error) {
	return value.Nil, fmt.Errorf("WATER.GETROT requires CGO")
}

func (m *Module) wSetScale(args []value.Value) (value.Value, error) {
	return value.Nil, fmt.Errorf("WATER.SETSCALE requires CGO")
}

func (m *Module) wGetScale(args []value.Value) (value.Value, error) {
	return value.Nil, fmt.Errorf("WATER.GETSCALE requires CGO")
}

func (m *Module) wGetColor(args []value.Value) (value.Value, error) {
	return value.Nil, fmt.Errorf("WATER.GETCOLOR requires CGO")
}

func (m *Module) wGetWaveHeight(args []value.Value) (value.Value, error) {
	return value.FromFloat(0), nil
}

func (m *Module) wGetWaveSpeed(args []value.Value) (value.Value, error) {
	return value.FromFloat(0), nil
}

func (m *Module) wGetShallowColor(args []value.Value) (value.Value, error) {
	return value.Nil, fmt.Errorf("WATER.GETSHALLOWCOLOR requires CGO")
}

func (m *Module) wGetDeepColor(args []value.Value) (value.Value, error) {
	return value.Nil, fmt.Errorf("WATER.GETDEEPCOLOR requires CGO")
}
