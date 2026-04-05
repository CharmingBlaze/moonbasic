//go:build cgo

package terrain

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)
func (t *TerrainObject) rebuildChunkMesh(cx, cz int) {
	cs := t.ChunkSize
	if cs < 2 {
		return
	}
	idx := idx2(t, cx, cz)
	if idx < 0 || idx >= len(t.Chunks) {
		return
	}
	ch := &t.Chunks[idx]

	// Bounds for this chunk in cell coordinates
	x0 := cx * cs
	z0 := cz * cs
	x1 := x0 + cs
	z1 := z0 + cs
	if x1 > t.WorldW {
		x1 = t.WorldW
	}
	if z1 > t.WorldH {
		z1 = t.WorldH
	}
	w := x1 - x0 + 1
	h := z1 - z0 + 1
	if w < 2 || h < 2 {
		return
	}

	minH := t.heightAtCell(x0, z0)
	maxH := minH
	for z := 0; z < h; z++ {
		for x := 0; x < w; x++ {
			v := t.heightAtCell(x0+x, z0+z)
			if v < minH {
				minH = v
			}
			if v > maxH {
				maxH = v
			}
		}
	}
	dh := maxH - minH
	if dh < 1e-4 {
		dh = 1
	}
	ch.MinH = minH
	ch.MaxH = maxH
	ch.BoundsValid = true

	im := rl.GenImageColor(int(w), int(h), rl.Color{R: 0, G: 0, B: 0, A: 0})
	for z := 0; z < h; z++ {
		for x := 0; x < w; x++ {
			v := t.heightAtCell(x0+x, z0+z)
			nv := (v - minH) / dh
			if nv < 0 {
				nv = 0
			}
			if nv > 1 {
				nv = 1
			}
			b := uint8(nv * 255)
			rl.ImageDrawPixel(im, int32(x), int32(z), rl.Color{R: b, G: b, B: b, A: 255})
		}
	}

	sizeX := float32(w-1) * t.CellSize
	sizeZ := float32(h-1) * t.CellSize
	sizeY := dh

	if ch.Loaded {
		rl.UnloadMaterial(ch.Mat)
		rl.UnloadMesh(&ch.Mesh)
		ch.Loaded = false
	}
	ch.Mesh = rl.GenMeshHeightmap(*im, rl.NewVector3(sizeX, sizeY, sizeZ))
	rl.UnloadImage(im)
	ch.Mat = rl.LoadMaterialDefault()
	ch.Loaded = true
	ch.Dirty = false
}

func (t *TerrainObject) drawChunk(cx, cz int) {
	idx := idx2(t, cx, cz)
	if idx < 0 || idx >= len(t.Chunks) {
		return
	}
	ch := &t.Chunks[idx]
	if !ch.Loaded {
		return
	}
	x0 := cx * t.ChunkSize
	z0 := cz * t.ChunkSize
	cs := t.ChunkSize
	x1 := x0 + cs
	z1 := z0 + cs
	if x1 > t.WorldW {
		x1 = t.WorldW
	}
	if z1 > t.WorldH {
		z1 = t.WorldH
	}
	var minH, maxH float32
	if ch.BoundsValid {
		minH, maxH = ch.MinH, ch.MaxH
	} else {
		minH = t.heightAtCell(x0, z0)
		maxH = minH
		for z := z0; z <= z1; z++ {
			for x := x0; x <= x1; x++ {
				v := t.heightAtCell(x, z)
				if v < minH {
					minH = v
				}
				if v > maxH {
					maxH = v
				}
			}
		}
		ch.MinH, ch.MaxH = minH, maxH
		ch.BoundsValid = true
	}
	dh := maxH - minH
	if dh < 1e-4 {
		dh = 1
	}
	tx := t.PX + float32(x0)*t.CellSize
	ty := t.PY + minH
	tz := t.PZ + float32(z0)*t.CellSize
	mat := rl.MatrixTranslate(tx, ty, tz)
	rl.DrawMesh(ch.Mesh, ch.Mat, mat)
}
