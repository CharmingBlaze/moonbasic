//go:build !cgo

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
