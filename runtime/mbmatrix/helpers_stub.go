//go:build !cgo && !windows

package mbmatrix

import (
	"fmt"

	"moonbasic/hal"
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

// AllocColorValue is unavailable without CGO (see COLOR.* / DECAL stubs).
func AllocColorValue(h *heap.Store, r, g, b, a uint8) (value.Value, error) {
	_ = h
	_ = r
	_ = g
	_ = b
	_ = a
	return value.Nil, fmt.Errorf("AllocColorValue: %s", hint)
}

// Vec2FromHandle is unavailable without CGO (see VEC2.* stubs).
func Vec2FromHandle(s *heap.Store, h heap.Handle) (hal.V2, error) {
	_ = s
	_ = h
	return hal.V2{}, fmt.Errorf("Vec2FromHandle: %s", hint)
}

// Vec3FromHandle is unavailable without CGO (see VEC3.* stubs).
func Vec3FromHandle(s *heap.Store, h heap.Handle) (hal.V3, error) {
	_ = s
	_ = h
	return hal.V3{}, fmt.Errorf("Vec3FromHandle: %s", hint)
}
