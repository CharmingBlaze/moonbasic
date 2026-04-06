//go:build !cgo && !windows

package scatter

import (
	"fmt"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

func registerScatter(m *Module, r runtime.Registrar) {
	h := func(n string) func(*runtime.Runtime, ...value.Value) (value.Value, error) {
		return func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
			return value.Nil, fmt.Errorf("%s requires CGO", n)
		}
	}
	r.Register("SCATTER.CREATE", "scatter", h("SCATTER.CREATE"))
	r.Register("SCATTER.FREE", "scatter", h("SCATTER.FREE"))
	r.Register("SCATTER.APPLY", "scatter", h("SCATTER.APPLY"))
	r.Register("SCATTER.DRAWALL", "scatter", h("SCATTER.DRAWALL"))
	r.Register("PROP.PLACE", "prop", h("PROP.PLACE"))
	r.Register("PROP.FREE", "prop", h("PROP.FREE"))
	r.Register("PROP.DRAWALL", "prop", h("PROP.DRAWALL"))
}
