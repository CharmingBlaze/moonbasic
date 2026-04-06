//go:build !cgo && !windows

package sky

import (
	"fmt"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

func registerSky(m *Module, r runtime.Registrar) {
	h := func(n string) func(*runtime.Runtime, ...value.Value) (value.Value, error) {
		return func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
			return value.Nil, fmt.Errorf("%s requires CGO", n)
		}
	}
	r.Register("SKY.MAKE", "sky", h("SKY.MAKE"))
	r.Register("SKY.FREE", "sky", h("SKY.FREE"))
	r.Register("SKY.UPDATE", "sky", h("SKY.UPDATE"))
	r.Register("SKY.DRAW", "sky", h("SKY.DRAW"))
	r.Register("SKY.SETTIME", "sky", h("SKY.SETTIME"))
	r.Register("SKY.SETDAYLENGTH", "sky", h("SKY.SETDAYLENGTH"))
	r.Register("SKY.GETTIMEHOURS", "sky", h("SKY.GETTIMEHOURS"))
	r.Register("SKY.ISNIGHT", "sky", h("SKY.ISNIGHT"))
}
