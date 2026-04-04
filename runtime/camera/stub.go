//go:build !cgo

package mbcamera

import (
	"fmt"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

const hint = "CAMERA.* requires CGO: set CGO_ENABLED=1 and install a C compiler, then rebuild"

// Register implements runtime.Module.
func (m *Module) Register(r runtime.Registrar) {
	stub := func(name string) runtime.BuiltinFn {
		return func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
			return value.Nil, fmt.Errorf("%s: %s", name, hint)
		}
	}
	r.Register("CAMERA.MAKE", "camera", stub("CAMERA.MAKE"))
	r.Register("CAMERA.SETPOS", "camera", stub("CAMERA.SETPOS"))
	r.Register("CAMERA.SETPOSITION", "camera", stub("CAMERA.SETPOSITION"))
	r.Register("CAMERA.SETTARGET", "camera", stub("CAMERA.SETTARGET"))
	r.Register("CAMERA.SETFOV", "camera", stub("CAMERA.SETFOV"))
	r.Register("CAMERA.BEGIN", "camera", stub("CAMERA.BEGIN"))
	r.Register("CAMERA.END", "camera", stub("CAMERA.END"))
	r.Register("CAMERA.MOVE", "camera", stub("CAMERA.MOVE"))
	r.Register("CAMERA.GETRAY", "camera", stub("CAMERA.GETRAY"))
	r.Register("CAMERA.GETVIEWRAY", "camera", stub("CAMERA.GETVIEWRAY"))
	r.Register("CAMERA.GETMATRIX", "camera", stub("CAMERA.GETMATRIX"))
	r.Register("CAMERA.GETPOS", "camera", stub("CAMERA.GETPOS"))
	r.Register("CAMERA.GETTARGET", "camera", stub("CAMERA.GETTARGET"))
	r.Register("CAMERA.SETUP", "camera", stub("CAMERA.SETUP"))
	r.Register("CAMERA.FREE", "camera", stub("CAMERA.FREE"))
	r.Register("MATRIX.FREE", "camera", stub("MATRIX.FREE"))
	r.Register("CAMERA2D.MAKE", "camera", stub("CAMERA2D.MAKE"))
	r.Register("CAMERA2D.SETTARGET", "camera", stub("CAMERA2D.SETTARGET"))
	r.Register("CAMERA2D.SETOFFSET", "camera", stub("CAMERA2D.SETOFFSET"))
	r.Register("CAMERA2D.SETZOOM", "camera", stub("CAMERA2D.SETZOOM"))
	r.Register("CAMERA2D.SETROTATION", "camera", stub("CAMERA2D.SETROTATION"))
	r.Register("CAMERA2D.BEGIN", "camera", stub("CAMERA2D.BEGIN"))
	r.Register("CAMERA2D.END", "camera", stub("CAMERA2D.END"))
	r.Register("CAMERA2D.GETMATRIX", "camera", stub("CAMERA2D.GETMATRIX"))
	r.Register("CAMERA2D.WORLDTOSCREEN", "camera", stub("CAMERA2D.WORLDTOSCREEN"))
	r.Register("CAMERA2D.SCREENTOWORLD", "camera", stub("CAMERA2D.SCREENTOWORLD"))
}

// Shutdown implements runtime.Module.
func (m *Module) Shutdown() {}
