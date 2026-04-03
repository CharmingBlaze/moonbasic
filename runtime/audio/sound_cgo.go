//go:build cgo

package mbaudio

import (
	"fmt"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

func (m *Module) registerSound(r runtime.Registrar) {
	r.Register("AUDIO.LOADSOUND", "audio", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		return value.Nil, fmt.Errorf("AUDIO.LOADSOUND: not implemented")
	}))
}
