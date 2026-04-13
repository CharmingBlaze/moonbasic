//go:build !cgo && !windows

package mbentity

import (
	"moonbasic/vm/heap"
	"moonbasic/vm/value"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// WorldPosFromEntityHandle is unavailable without the full entity runtime.
func (m *Module) WorldPosFromEntityHandle(_ heap.Handle) (rl.Vector3, bool) {
	return rl.Vector3{}, false
}

// GetWorldPosByID is unavailable without the full entity runtime.
func (m *Module) GetWorldPosByID(_ int) (rl.Vector3, bool) {
	return rl.Vector3{}, false
}

// ResolveEntityID resolves a numeric entity id only (no EntityRef handles without CGO).
func (m *Module) ResolveEntityID(v value.Value) (int64, bool) {
	id, ok := v.ToInt()
	if !ok || id < 1 {
		return 0, false
	}
	return id, true
}
