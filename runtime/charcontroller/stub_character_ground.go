//go:build !linux || !cgo

package mbcharcontroller

import "moonbasic/vm/heap"

// CharacterGroundNormal is only available with Jolt (linux+cgo build).
func (m *Module) CharacterGroundNormal(h heap.Handle) (nx, ny, nz float64, ok bool) {
	_ = m
	_ = h
	return 0, 0, 0, false
}
