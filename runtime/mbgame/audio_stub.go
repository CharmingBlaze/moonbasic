//go:build !cgo

package mbgame

import (
	"fmt"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

func (m *Module) registerAudioHelpers(r runtime.Registrar) {
	_ = m
	stub := func(name string) runtime.BuiltinFn {
		return func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
			_ = rt
			_ = args
			return value.Nil, fmt.Errorf("%s requires CGO_ENABLED=1", name)
		}
	}
	r.Register("GAME.SETMASTERVOLUME", "game", stub("GAME.SETMASTERVOLUME"))
	r.Register("GAME.GETMASTERVOLUME", "game", stub("GAME.GETMASTERVOLUME"))
}
