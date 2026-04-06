//go:build !cgo && !windows

package input

import (
	"fmt"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

func (m *Module) registerInputAdvanced(reg runtime.Registrar) {
	hint := "INPUT mouse/touch/gamepad requires CGO (Raylib)"
	reg.Register("INPUT.MOUSEX", "input", stubErr(hint, "INPUT.MOUSEX"))
	reg.Register("INPUT.MOUSEY", "input", stubErr(hint, "INPUT.MOUSEY"))
	reg.Register("INPUT.MOUSEDOWN", "input", stubErr(hint, "INPUT.MOUSEDOWN"))
	reg.Register("INPUT.SETMOUSESCALE", "input", stubErr(hint, "INPUT.SETMOUSESCALE"))
	reg.Register("INPUT.SETMOUSEOFFSET", "input", stubErr(hint, "INPUT.SETMOUSEOFFSET"))
	reg.Register("INPUT.GETMOUSEWORLDPOS", "input", stubErr(hint, "INPUT.GETMOUSEWORLDPOS"))
	reg.Register("INPUT.TOUCHCOUNT", "input", stubErr(hint, "INPUT.TOUCHCOUNT"))
	reg.Register("INPUT.TOUCHX", "input", stubErr(hint, "INPUT.TOUCHX"))
	reg.Register("INPUT.TOUCHY", "input", stubErr(hint, "INPUT.TOUCHY"))
	reg.Register("INPUT.TOUCHPRESSED", "input", stubErr(hint, "INPUT.TOUCHPRESSED"))
	reg.Register("INPUT.GETTOUCHPOINTID", "input", stubErr(hint, "INPUT.GETTOUCHPOINTID"))
	reg.Register("INPUT.GAMEPADBUTTONCOUNT", "input", stubErr(hint, "INPUT.GAMEPADBUTTONCOUNT"))
	reg.Register("INPUT.GAMEPADAXISCOUNT", "input", stubErr(hint, "INPUT.GAMEPADAXISCOUNT"))
	reg.Register("INPUT.SETGAMEPADMAPPINGS", "input", stubErr(hint, "INPUT.SETGAMEPADMAPPINGS"))
}

func stubErr(hint, name string) runtime.BuiltinFn {
	return func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		_ = rt
		return value.Nil, fmt.Errorf("%s: %s", name, hint)
	}
}
