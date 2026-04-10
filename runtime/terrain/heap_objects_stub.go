//go:build !cgo && !windows

package terrain

import (
	"math"

	"moonbasic/vm/heap"
)

// TerrainObject is a heightfield with optional chunk meshes for rendering.
// This variant has no Raylib GPU fields — TERRAIN.* is stubbed when CGO is off on Linux.
type TerrainObject struct {
	WorldW    int
	WorldH    int
	CellSize  float32
	ChunkSize int
	Heights   []float32

	ChunkW int
	ChunkH int

	PX, PY, PZ float32

	Chunks []chunkSlot

	StreamEnabled bool
	LoadDist      float32
	UnloadDist    float32
	CenterX       float32
	CenterZ       float32

	MaxHeight float32

	ScaleX, ScaleY, ScaleZ float32
	DetailFactor           float32

	freed bool
}

type chunkSlot struct {
	Loaded     bool
	Dirty      bool
	CX, CZ     int
	LastUpload int64
	MinH, MaxH float32
	BoundsValid bool
}

func (t *TerrainObject) TypeName() string { return "Terrain" }
func (t *TerrainObject) TypeTag() uint16  { return heap.TagTerrain }

func (t *TerrainObject) Free() {
	if t.freed {
		return
	}
	t.Chunks = nil
	t.Heights = nil
	t.freed = true
}

func idx2(t *TerrainObject, cx, cz int) int {
	return cz*t.ChunkW + cx
}

func cellIndex(t *TerrainObject, wx, wz int) int {
	return wz*t.WorldW + wx
}

func (t *TerrainObject) heightAtCell(wx, wz int) float32 {
	if wx < 0 || wz < 0 || wx >= t.WorldW || wz >= t.WorldH {
		return 0
	}
	return t.Heights[cellIndex(t, wx, wz)]
}

func (t *TerrainObject) scaleXEff() float32 {
	if t.ScaleX <= 0 {
		return 1
	}
	return t.ScaleX
}
func (t *TerrainObject) scaleYEff() float32 {
	if t.ScaleY <= 0 {
		return 1
	}
	return t.ScaleY
}
func (t *TerrainObject) scaleZEff() float32 {
	if t.ScaleZ <= 0 {
		return 1
	}
	return t.ScaleZ
}

func (t *TerrainObject) HeightWorld(x, z float32) float32 {
	if t.WorldW < 2 || t.WorldH < 2 {
		return 0
	}
	sx := t.scaleXEff()
	sz := t.scaleZEff()
	sy := t.scaleYEff()
	lx := (x - t.PX) / (t.CellSize * sx)
	lz := (z - t.PZ) / (t.CellSize * sz)
	if lx < 0 || lz < 0 || lx >= float32(t.WorldW-1) || lz >= float32(t.WorldH-1) {
		return 0
	}
	x0 := int(lx)
	z0 := int(lz)
	fx := lx - float32(x0)
	fz := lz - float32(z0)
	h00 := t.heightAtCell(x0, z0)
	h10 := t.heightAtCell(x0+1, z0)
	h01 := t.heightAtCell(x0, z0+1)
	h11 := t.heightAtCell(x0+1, z0+1)
	a := h00*(1-fx) + h10*fx
	b := h01*(1-fx) + h11*fx
	raw := a*(1-fz) + b*fz
	return raw*sy + t.PY
}

func max32StubTerrain(a, b float32) float32 {
	if a > b {
		return a
	}
	return b
}

// GridXZ matches the raylib TerrainObject implementation (fractional grid coords).
func (t *TerrainObject) GridXZ(x, z float32) (lx, lz float32, ok bool) {
	if t.WorldW < 2 || t.WorldH < 2 {
		return 0, 0, false
	}
	sx := t.scaleXEff()
	sz := t.scaleZEff()
	lx = (x - t.PX) / (t.CellSize * sx)
	lz = (z - t.PZ) / (t.CellSize * sz)
	if lx < 0 || lz < 0 || lx >= float32(t.WorldW-1) || lz >= float32(t.WorldH-1) {
		return lx, lz, false
	}
	return lx, lz, true
}

func (t *TerrainObject) SlopeDeg(x, z float32) float32 {
	d := t.CellSize * 0.5 * max32StubTerrain(t.scaleXEff(), t.scaleZEff())
	if d < 1e-4 {
		d = 1
	}
	hL := t.HeightWorld(x-d, z)
	hR := t.HeightWorld(x+d, z)
	hD := t.HeightWorld(x, z-d)
	hU := t.HeightWorld(x, z+d)
	dhdx := (hR - hL) / (2 * d)
	dhdz := (hU - hD) / (2 * d)
	grad := math.Sqrt(float64(dhdx*dhdx + dhdz*dhdz))
	return float32(math.Atan(grad) * 180 / math.Pi)
}
