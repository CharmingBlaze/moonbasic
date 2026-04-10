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

// SlopeDegPublic returns approximate terrain slope in degrees at (x,z), or 0 if the handle is not terrain.
func SlopeDegPublic(h *heap.Store, id heap.Handle, x, z float32) float32 {
	o, err := heap.Cast[*TerrainObject](h, id)
	if err != nil {
		return 0
	}
	return o.SlopeDeg(x, z)
}

// GridXZPublic returns fractional grid coordinates; ok is false when (x,z) is outside the heightfield.
func GridXZPublic(h *heap.Store, id heap.Handle, x, z float32) (lx, lz float32, ok bool) {
	o, err := heap.Cast[*TerrainObject](h, id)
	if err != nil {
		return 0, 0, false
	}
	return o.GridXZ(x, z)
}

// WorldXZBounds returns the axis-aligned XZ extent of the heightfield (inclusive corners of the sample grid).
func WorldXZBounds(h *heap.Store, id heap.Handle) (x0, z0, x1, z1 float64, ok bool) {
	o, err := heap.Cast[*TerrainObject](h, id)
	if err != nil || o.WorldW < 2 || o.WorldH < 2 {
		return 0, 0, 0, 0, false
	}
	sx := float64(o.CellSize * o.scaleXEff())
	sz := float64(o.CellSize * o.scaleZEff())
	x0 = float64(o.PX)
	z0 = float64(o.PZ)
	x1 = x0 + sx*float64(o.WorldW-1)
	z1 = z0 + sz*float64(o.WorldH-1)
	return x0, z0, x1, z1, true
}
