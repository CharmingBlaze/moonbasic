//go:build !cgo && !windows

package mbmatrix

import (
	"fmt"

	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

// AllocVec3Value is unavailable without CGO (see VEC3.* stubs).
func AllocVec3Value(h *heap.Store, x, y, z float32) (value.Value, error) {
	_ = h
	_ = x
	_ = y
	_ = z
	return value.Nil, fmt.Errorf("AllocVec3Value: %s", hint)
}

// AllocVec2Value is unavailable without CGO (see VEC2.* stubs).
func AllocVec2Value(h *heap.Store, x, y float32) (value.Value, error) {
	_ = h
	_ = x
	_ = y
	return value.Nil, fmt.Errorf("AllocVec2Value: %s", hint)
}
