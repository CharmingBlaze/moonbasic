package mbphysics3d

import (
	"fmt"

	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

// valueVec3FromFloats returns a 3-element numeric array handle (same shape as BODY3D.GETPOS).
func valueVec3FromFloats(h *heap.Store, x, y, z float64) (value.Value, error) {
	if h == nil {
		return value.Nil, fmt.Errorf("heap not bound")
	}
	arr, err := heap.NewArray([]int64{3})
	if err != nil {
		return value.Nil, err
	}
	_ = arr.Set([]int64{0}, x)
	_ = arr.Set([]int64{1}, y)
	_ = arr.Set([]int64{2}, z)
	id, err := h.Alloc(arr)
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}
