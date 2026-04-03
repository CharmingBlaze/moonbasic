//go:build cgo

package mbaudio

import (
	"fmt"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

func (m *Module) registerMusic(r runtime.Registrar) {
	r.Register("AUDIO.LOADMUSIC", "audio", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		return value.Nil, fmt.Errorf("AUDIO.LOADMUSIC: not implemented")
	}))
}
