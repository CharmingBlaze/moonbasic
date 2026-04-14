//go:build (!linux && !windows) || !cgo

package mbphysics3d

import (
	"moonbasic/vm/heap"
)

// RegisterEntityBufferLink records that an entity uses a physics matrix buffer index (no-op without Jolt).
func RegisterEntityBufferLink(entityID int64, bufIdx int) {}

// UnregisterEntityCollision removes collision bookkeeping for an entity (no-op without Jolt).
func UnregisterEntityCollision(entityID int64) {}

func registerBufferBodyForCollision(bufIdx int, bodyHandle heap.Handle) {}

func unregisterBufferBodyForCollision(bufIdx int) {}

func collectContactsAfterStep() {}

// Frame collision accessors (stubs).
func PairCollidedThisFrame(a, b int64) (hit ContactHitData, ok bool) { return ContactHitData{}, false }

func CountCollisionsForEntity(e int64) int { return 0 }

func LastCollisionData() ContactHitData { return ContactHitData{} }

// EntityIDForBodyHandle is a stub (no Jolt body map).
func EntityIDForBodyHandle(bodyH heap.Handle) (int64, bool) { return 0, false }

// ReviveEntity is a no-op on stub builds.
func ReviveEntity(id int64) {}

// ContactHitData mirrors linux; exported for stub symmetry.
type ContactHitData struct {
	NX, NY, NZ float64
	PX, PY, PZ float64
	Force      float64
}
