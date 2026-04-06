//go:build !cgo && !windows

package input

import (
	"fmt"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

const gestureHint = "GESTURE.* natives require CGO: set CGO_ENABLED=1 and install a C compiler, then rebuild"

func registerGesture(r runtime.Registrar) {
	stub := func(name string) runtime.BuiltinFn {
		return func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
			return value.Nil, fmt.Errorf("%s: %s", name, gestureHint)
		}
	}
	r.Register("GESTURE.ENABLE", "input", stub("GESTURE.ENABLE"))
	r.Register("GESTURE.ISDETECTED", "input", stub("GESTURE.ISDETECTED"))
	r.Register("GESTURE.GETDETECTED", "input", stub("GESTURE.GETDETECTED"))
	r.Register("GESTURE.GETHOLDDURATION", "input", stub("GESTURE.GETHOLDDURATION"))
	r.Register("GESTURE.GETDRAGVECTORX", "input", stub("GESTURE.GETDRAGVECTORX"))
	r.Register("GESTURE.GETDRAGVECTORY", "input", stub("GESTURE.GETDRAGVECTORY"))
	r.Register("GESTURE.GETDRAGANGLE", "input", stub("GESTURE.GETDRAGANGLE"))
	r.Register("GESTURE.GETPINCHVECTORX", "input", stub("GESTURE.GETPINCHVECTORX"))
	r.Register("GESTURE.GETPINCHVECTORY", "input", stub("GESTURE.GETPINCHVECTORY"))
	r.Register("GESTURE.GETPINCHANGLE", "input", stub("GESTURE.GETPINCHANGLE"))
}
