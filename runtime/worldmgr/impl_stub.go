//go:build !cgo && !windows

package worldmgr

import (
	"fmt"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

func (m *Module) worldSetReflection(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	return value.Nil, fmt.Errorf("WORLD.SETREFLECTION requires CGO")
}

func (m *Module) worldFogMode(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	return value.Nil, fmt.Errorf("WORLD.FOGMODE requires CGO")
}

func (m *Module) worldFogColor(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	return value.Nil, fmt.Errorf("WORLD.FOGCOLOR requires CGO")
}

func (m *Module) worldFogDensity(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	return value.Nil, fmt.Errorf("WORLD.FOGDENSITY requires CGO")
}

func worldSetCenter(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	return value.Nil, fmt.Errorf("WORLD.SETCENTER requires CGO")
}

func worldSetCenterEntity(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	return value.Nil, fmt.Errorf("WORLD.SETCENTERENTITY requires CGO")
}

func worldUpdate(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	return value.Nil, nil
}

func worldStreamEnable(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	return value.Nil, nil
}

func worldPreload(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	return value.Nil, nil
}

func worldStatus(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	return rt.RetString("stubs-only")
}

func worldIsReady(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	return value.FromBool(false), nil
}

func worldSetVegetation(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	return value.Nil, fmt.Errorf("WORLD.SETVEGETATION requires CGO")
}
