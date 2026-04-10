package window

import (
	"fmt"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

func (m *Module) registerLoadingModeCommands(reg runtime.Registrar) {
	reg.Register("WINDOW.SETLOADINGMODE", "window", m.wSetLoadingMode)
	reg.Register("WINDOW.LOADINGMODE", "window", m.wGetLoadingMode)
}

func (m *Module) wSetLoadingMode(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if rt == nil {
		return value.Nil, fmt.Errorf("WINDOW.SETLOADINGMODE: runtime not available")
	}
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("WINDOW.SETLOADINGMODE expects 1 argument (enabled?)")
	}
	b, err := rt.ArgBool(args, 0)
	if err != nil {
		return value.Nil, err
	}
	rt.SetLoadingMode(b)
	return value.Nil, nil
}

func (m *Module) wGetLoadingMode(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("WINDOW.LOADINGMODE expects 0 arguments")
	}
	if rt == nil {
		return value.FromBool(false), nil
	}
	return value.FromBool(rt.LoadingMode()), nil
}
