//go:build cgo || (windows && !cgo)

package terrain

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

func (t *TerrainObject) detailStep() int {
	d := t.DetailFactor
	if d <= 0 || d > 1 {
		d = 1
	}
	step := int(1.0/float64(d) + 0.5)
	if step < 1 {
		step = 1
	}
	return step
}

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

	step := t.detailStep()
	ws := (w + step - 1) / step
	hs := (h + step - 1) / step
	if ws < 2 {
		ws = 2
	}
	if hs < 2 {
		hs = 2
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

	sy := t.scaleYEff()
	im := rl.GenImageColor(int(ws), int(hs), rl.Color{R: 0, G: 0, B: 0, A: 0})
	for z := 0; z < hs; z++ {
		for x := 0; x < ws; x++ {
			sx := x * step
			sz := z * step
			if sx >= w {
				sx = w - 1
			}
			if sz >= h {
				sz = h - 1
			}
			v := t.heightAtCell(x0+sx, z0+sz)
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

	sx := t.scaleXEff()
	sz := t.scaleZEff()
	sizeX := float32(ws-1) * t.CellSize * float32(step) * sx
	sizeZ := float32(hs-1) * t.CellSize * float32(step) * sz
	sizeY := dh * sy

	if ch.Loaded {
		rl.UnloadMaterial(ch.Mat)
		rl.UnloadMesh(&ch.Mesh)
		ch.Loaded = false
	}
	ch.Mesh = rl.GenMeshHeightmap(*im, rl.NewVector3(sizeX, sizeY, sizeZ))
	rl.UnloadImage(im)
	ch.Mat = rl.LoadMaterialDefault()
	if t.DiffuseLoaded {
		rl.SetMaterialTexture(&ch.Mat, rl.MapAlbedo, t.DiffuseTex)
	}
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
	sx := t.scaleXEff()
	sz := t.scaleZEff()
	sy := t.scaleYEff()
	tx := t.PX + float32(x0)*t.CellSize*sx
	ty := t.PY + minH*sy
	tz := t.PZ + float32(z0)*t.CellSize*sz
	mat := rl.MatrixTranslate(tx, ty, tz)
	rl.DrawMesh(ch.Mesh, ch.Mat, mat)
}
