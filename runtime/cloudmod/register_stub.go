//go:build !cgo && !windows

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
	r.Register("CLOUD.CREATE", "cloud", h("CLOUD.CREATE"))
	r.Register("CLOUD.MAKE", "cloud", h("CLOUD.MAKE"))
	r.Register("CLOUD.FREE", "cloud", h("CLOUD.FREE"))
	r.Register("CLOUD.UPDATE", "cloud", h("CLOUD.UPDATE"))
	r.Register("CLOUD.DRAW", "cloud", h("CLOUD.DRAW"))
	r.Register("CLOUD.SETCOVERAGE", "cloud", h("CLOUD.SETCOVERAGE"))
	r.Register("CLOUD.GETCOVERAGE", "cloud", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		return value.FromFloat(0.3), nil
	})
}
