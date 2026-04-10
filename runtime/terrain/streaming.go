//go:build cgo || (windows && !cgo)

package terrain

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func chunkDistanceMeters(t *TerrainObject, cx, cz int) float64 {
	cs := t.ChunkSize
	sx := float64(t.scaleXEff())
	sz := float64(t.scaleZEff())
	ccx := float64(t.PX) + (float64(cx*cs)+float64(cs)*0.5)*float64(t.CellSize)*sx
	ccz := float64(t.PZ) + (float64(cz*cs)+float64(cs)*0.5)*float64(t.CellSize)*sz
	dx := float64(t.CenterX) - ccx
	dz := float64(t.CenterZ) - ccz
	return math.Sqrt(dx*dx + dz*dz)
}

// TickStreaming loads/unloads chunk meshes based on center and distances.
func (t *TerrainObject) TickStreaming() {
	if t.freed || !t.StreamEnabled {
		return
	}
	if t.LoadDist <= 0 {
		t.LoadDist = 400
	}
	if t.UnloadDist <= 0 {
		t.UnloadDist = 600
	}
	for cz := 0; cz < t.ChunkH; cz++ {
		for cx := 0; cx < t.ChunkW; cx++ {
			idx := idx2(t, cx, cz)
			ch := &t.Chunks[idx]
			d := chunkDistanceMeters(t, cx, cz)
			if d <= float64(t.LoadDist) {
				if !ch.Loaded || ch.Dirty {
					t.rebuildChunkMesh(cx, cz)
				}
			} else if d >= float64(t.UnloadDist) && ch.Loaded {
				rl.UnloadMaterial(ch.Mat)
				rl.UnloadMesh(&ch.Mesh)
				ch.Loaded = false
			}
		}
	}
}

// PreloadChunksAround loads every chunk within chunkRadius of the chunk containing (CenterX, CenterZ).
func (t *TerrainObject) PreloadChunksAround(chunkRadius int) {
	if t.freed || chunkRadius < 0 {
		return
	}
	// Find chunk under center
	lx := (t.CenterX - t.PX) / t.CellSize
	lz := (t.CenterZ - t.PZ) / t.CellSize
	if lx < 0 {
		lx = 0
	}
	if lz < 0 {
		lz = 0
	}
	ccx := int(lx) / t.ChunkSize
	ccz := int(lz) / t.ChunkSize
	for dz := -chunkRadius; dz <= chunkRadius; dz++ {
		for dx := -chunkRadius; dx <= chunkRadius; dx++ {
			cx := ccx + dx
			cz := ccz + dz
			if cx < 0 || cz < 0 || cx >= t.ChunkW || cz >= t.ChunkH {
				continue
			}
			t.rebuildChunkMesh(cx, cz)
		}
	}
}

func (t *TerrainObject) loadedChunkCount() int {
	n := 0
	for i := range t.Chunks {
		if t.Chunks[i].Loaded {
			n++
		}
	}
	return n
}
