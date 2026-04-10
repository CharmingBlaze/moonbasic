//go:build !(cgo || (windows && !cgo))

package player

import (
	"fmt"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

func registerPlayerTerrainCommands(m *Module, reg runtime.Registrar) {
	_ = m
	reg.Register("PLAYER.SNAPTOGROUND", "player", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		_ = rt
		_ = args
		return value.Nil, fmt.Errorf("PLAYER.SNAPTOGROUND requires CGO and Raylib terrain")
	})
	reg.Register("PLAYER.ISSWIMMING", "player", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		_ = rt
		_ = args
		return value.FromBool(false), nil
	})
}
