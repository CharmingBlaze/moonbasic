package runtime

import (
	"moonbasic/drivers/video/null"
	"moonbasic/hal"
	"moonbasic/vm/heap"
)

// NewRegistryHeadless returns a registry backed by the null HAL driver (no GPU / no Raylib).
// Use from tests and tooling that only need InitCore and the VM.
func NewRegistryHeadless(h *heap.Store) *Registry {
	d := null.NewDriver()
	return NewRegistry(h, hal.Driver{
		Video:  d,
		Input:  d,
		System: d,
	})
}
