package terrain

import (
	"moonbasic/vm/heap"
)

// HeightWorldPublic samples bilinear height for any terrain handle (used by scatter/biome).
func HeightWorldPublic(h *heap.Store, id heap.Handle, x, z float32) float32 {
	o, err := heap.Cast[*TerrainObject](h, id)
	if err != nil {
		return 0
	}
	return o.HeightWorld(x, z)
}
