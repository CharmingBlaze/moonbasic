//go:build cgo || (windows && !cgo)

package terrain

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/vm/heap"
)

// TerrainObject is a heightfield with optional chunk meshes for rendering.
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

	freed bool
}

// chunkSlot owns GPU mesh and default material for one terrain chunk (rebuildChunkMesh).
type chunkSlot struct {
	Mesh       rl.Mesh
	Mat        rl.Material
	Loaded     bool
	Dirty      bool
	CX, CZ     int
	LastUpload int64
	// MinH/MaxH are heightfield samples (before TerrainObject.PY); set in rebuildChunkMesh.
	MinH, MaxH float32
	BoundsValid bool
}

func (t *TerrainObject) TypeName() string { return "Terrain" }
func (t *TerrainObject) TypeTag() uint16  { return heap.TagTerrain }

// Free releases GPU meshes and height data.
func (t *TerrainObject) Free() {
	if t.freed {
		return
	}
	for i := range t.Chunks {
		ch := &t.Chunks[i]
		if ch.Loaded {
			rl.UnloadMaterial(ch.Mat)
			rl.UnloadMesh(&ch.Mesh)
			ch.Loaded = false
		}
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

// HeightWorld returns bilinear height at world xz (relative to origin px,pz).
func (t *TerrainObject) HeightWorld(x, z float32) float32 {
	if t.WorldW < 2 || t.WorldH < 2 {
		return 0
	}
	lx := (x - t.PX) / t.CellSize
	lz := (z - t.PZ) / t.CellSize
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
	return a*(1-fz) + b*fz + t.PY
}

// SlopeDeg approximate slope angle at world position (degrees).
func (t *TerrainObject) SlopeDeg(x, z float32) float32 {
	d := t.CellSize * 0.5
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
