//go:build !cgo

package mbutil

import (
	"fmt"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

func (m *Module) registerDroppedFiles(r runtime.Registrar) {
	isDropped := runtime.AdaptLegacy(func([]value.Value) (value.Value, error) {
		return value.FromBool(false), nil
	})
	noRaylib := runtime.AdaptLegacy(func([]value.Value) (value.Value, error) {
		return value.Nil, fmt.Errorf("UTIL.GETDROPPEDFILES requires CGO (Raylib)")
	})
	clearDropped := runtime.AdaptLegacy(func([]value.Value) (value.Value, error) {
		return value.Nil, nil
	})
	r.Register("UTIL.ISFILEDROPPED", "util", isDropped)
	r.Register("ISFILEDROPPED", "util", isDropped)
	r.Register("UTIL.GETDROPPEDFILES", "util", noRaylib)
	r.Register("GETDROPPEDFILES", "util", noRaylib)
	r.Register("UTIL.CLEARDROPPEDFILES", "util", clearDropped)
}
