//go:build !cgo

package mbutil

import (
	"fmt"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

func (m *Module) registerDroppedFiles(r runtime.Registrar) {
	r.Register("UTIL.ISFILEDROPPED", "util", runtime.AdaptLegacy(func([]value.Value) (value.Value, error) {
		return value.FromBool(false), nil
	}))
	r.Register("UTIL.GETDROPPEDFILES", "util", runtime.AdaptLegacy(func([]value.Value) (value.Value, error) {
		return value.Nil, fmt.Errorf("UTIL.GETDROPPEDFILES requires CGO (Raylib)")
	}))
	r.Register("UTIL.CLEARDROPPEDFILES", "util", runtime.AdaptLegacy(func([]value.Value) (value.Value, error) {
		return value.Nil, nil
	}))
}
