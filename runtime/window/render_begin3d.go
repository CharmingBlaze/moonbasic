package window

import (
	"fmt"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

// rBegin3D and rEnd3D delegate to CAMERA.BEGIN / CAMERA.END (shared by CGO, purego, and stub builds).

func (m *Module) rBegin3D(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("RENDER.BEGIN3D expects 1 argument (camera handle)")
	}
	reg := runtime.ActiveRegistry()
	if reg == nil {
		return value.Nil, fmt.Errorf("RENDER.BEGIN3D: registry not active")
	}
	return reg.Call("CAMERA.BEGIN", args)
}

func (m *Module) rEnd3D(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("RENDER.END3D expects 0 arguments")
	}
	reg := runtime.ActiveRegistry()
	if reg == nil {
		return value.Nil, fmt.Errorf("RENDER.END3D: registry not active")
	}
	return reg.Call("CAMERA.END", nil)
}
