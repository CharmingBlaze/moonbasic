//go:build cgo || (windows && !cgo)

package input

import (
	"fmt"

	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

// movement 2D implementation below

// movement2D returns a 2-float array [forwardAxis, strafeAxis] from two Input.Axis pairs.
// Caller should ERASE the handle when done. Same as Axis(keyBack,keyForward) and Axis(keyLeft,keyRight).
func (m *Module) movement2D(args []value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, fmt.Errorf("INPUT.MOVEMENT2D: heap not bound")
	}
	if len(args) != 4 {
		return value.Nil, fmt.Errorf("INPUT.MOVEMENT2D expects 4 arguments (keyBack, keyForward, keyLeft, keyRight)")
	}
	f, err := axisValueFromKeys(args[0], args[1])
	if err != nil {
		return value.Nil, err
	}
	s, err := axisValueFromKeys(args[2], args[3])
	if err != nil {
		return value.Nil, err
	}
	arr, err := heap.NewArray([]int64{2})
	if err != nil {
		return value.Nil, err
	}
	_ = arr.Set([]int64{0}, f)
	_ = arr.Set([]int64{1}, s)
	id, err := m.h.Alloc(arr)
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}
