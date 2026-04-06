//go:build !cgo && !windows

package input

import (
	"fmt"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

const cursorHint = "CURSOR.* natives require CGO: set CGO_ENABLED=1 and install a C compiler, then rebuild"

func registerCursor(r runtime.Registrar) {
	stub := func(name string) runtime.BuiltinFn {
		return func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
			return value.Nil, fmt.Errorf("%s: %s", name, cursorHint)
		}
	}
	r.Register("CURSOR.SHOW", "input", stub("CURSOR.SHOW"))
	r.Register("CURSOR.HIDE", "input", stub("CURSOR.HIDE"))
	r.Register("CURSOR.ISHIDDEN", "input", stub("CURSOR.ISHIDDEN"))
	r.Register("CURSOR.ISONSCREEN", "input", stub("CURSOR.ISONSCREEN"))
	r.Register("CURSOR.ENABLE", "input", stub("CURSOR.ENABLE"))
	r.Register("CURSOR.DISABLE", "input", stub("CURSOR.DISABLE"))
	r.Register("CURSOR.SET", "input", stub("CURSOR.SET"))
}
