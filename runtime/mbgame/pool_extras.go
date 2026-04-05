package mbgame

import (
	"fmt"

	"moonbasic/runtime"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

func (m *Module) registerPoolExtras(r runtime.Registrar) {
	r.Register("GAME.MAKEFLOATARRAY", "game", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		_ = rt
		if m.h == nil {
			return value.Nil, fmt.Errorf("GAME.MAKEFLOATARRAY: heap not bound")
		}
		if len(args) != 1 {
			return value.Nil, fmt.Errorf("GAME.MAKEFLOATARRAY expects 1 argument (length)")
		}
		n, ok := argI(args[0])
		if !ok || n < 1 {
			return value.Nil, fmt.Errorf("GAME.MAKEFLOATARRAY: length must be a positive integer")
		}
		a, err := heap.NewArray([]int64{n})
		if err != nil {
			return value.Nil, err
		}
		h, err := m.h.Alloc(a)
		if err != nil {
			return value.Nil, err
		}
		return value.FromHandle(h), nil
	})
}
