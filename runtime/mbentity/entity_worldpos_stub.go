//go:build !cgo && !windows

package mbentity

import (
	"moonbasic/vm/heap"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// WorldPosFromEntityHandle is unavailable without the full entity runtime.
func (m *Module) WorldPosFromEntityHandle(_ heap.Handle) (rl.Vector3, bool) {
	return rl.Vector3{}, false
}
