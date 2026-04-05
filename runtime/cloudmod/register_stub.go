//go:build !cgo

package cloudmod

import (
	"fmt"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

func registerCloud(m *Module, r runtime.Registrar) {
	h := func(n string) func(*runtime.Runtime, ...value.Value) (value.Value, error) {
		return func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
			return value.Nil, fmt.Errorf("%s requires CGO", n)
		}
	}
	r.Register("CLOUD.MAKE", "cloud", h("CLOUD.MAKE"))
	r.Register("CLOUD.FREE", "cloud", h("CLOUD.FREE"))
	r.Register("CLOUD.UPDATE", "cloud", h("CLOUD.UPDATE"))
	r.Register("CLOUD.DRAW", "cloud", h("CLOUD.DRAW"))
	r.Register("CLOUD.SETCOVERAGE", "cloud", h("CLOUD.SETCOVERAGE"))
}
