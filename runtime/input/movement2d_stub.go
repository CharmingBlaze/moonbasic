//go:build !cgo && !windows

package input

import (
	"fmt"

	"moonbasic/runtime"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

func (m *Module) movement2D(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if m.h == nil {
		return value.Nil, fmt.Errorf("INPUT.MOVEMENT2D: heap not bound")
	}
	if len(args) != 4 {
		return value.Nil, fmt.Errorf("INPUT.MOVEMENT2D expects 4 arguments (keyBack, keyForward, keyLeft, keyRight)")
	}
	arr, err := heap.NewArray([]int64{2})
	if err != nil {
		return value.Nil, err
	}
	id, err := m.h.Alloc(arr)
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}
