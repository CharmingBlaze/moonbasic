//go:build !cgo && !windows

package input

import (
	"fmt"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

func (m *Module) registerBlitzAliases(r runtime.Registrar) {
	stub := func(name string) runtime.BuiltinFn {
		return func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
			_ = rt
			_ = args
			return value.Nil, fmt.Errorf("%s requires CGO (raylib)", name)
		}
	}
	for _, n := range []string{
		"INPUT.KEYHIT", "INPUT.MOUSEXSPEED", "INPUT.MOUSEYSPEED",
		"INPUT.JOYX", "INPUT.JOYY", "INPUT.JOYBUTTON", "INPUT.JOYDOWN",
	} {
		r.Register(n, "input", stub(n))
	}
}
