//go:build !cgo

package water

import (
	"fmt"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

func registerWater(m *Module, r runtime.Registrar) {
	h := func(name string) func(*runtime.Runtime, ...value.Value) (value.Value, error) {
		return func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
			return value.Nil, fmt.Errorf("%s requires CGO", name)
		}
	}
	r.Register("WATER.MAKE", "water", h("WATER.MAKE"))
	r.Register("WATER.FREE", "water", h("WATER.FREE"))
	r.Register("WATER.SETPOS", "water", h("WATER.SETPOS"))
	r.Register("WATER.DRAW", "water", h("WATER.DRAW"))
	r.Register("WATER.UPDATE", "water", h("WATER.UPDATE"))
	r.Register("WATER.SETWAVEHEIGHT", "water", h("WATER.SETWAVEHEIGHT"))
	r.Register("WATER.GETWAVEY", "water", h("WATER.GETWAVEY"))
	r.Register("WATER.GETDEPTH", "water", h("WATER.GETDEPTH"))
	r.Register("WATER.ISUNDER", "water", h("WATER.ISUNDER"))
	r.Register("WATER.SETSHALLOWCOLOR", "water", h("WATER.SETSHALLOWCOLOR"))
	r.Register("WATER.SETDEEPCOLOR", "water", h("WATER.SETDEEPCOLOR"))
}
