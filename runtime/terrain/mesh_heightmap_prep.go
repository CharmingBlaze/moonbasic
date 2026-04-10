//go:build (cgo || (windows && !cgo)) && (!windows || !gopls_stub)

package terrain

// chunkHeightSnapshot is a copy of height samples for one chunk cell rectangle.
// Safe to process off the main thread (no TerrainObject access after snapshot).
type chunkHeightSnapshot struct {
	X0, Z0 int
	W, H   int
	// Heights row-major: index lz*W + lx
	Heights []float32

	CellSize     float32
	DetailFactor float32
	ScaleX       float32
	ScaleY       float32
	ScaleZ       float32
}

func detailStepForTerrainDetail(d float32) int {
	d2 := d
	if d2 <= 0 || d2 > 1 {
		d2 = 1
	}
	step := int(1.0/float64(d2) + 0.5)
	if step < 1 {
		step = 1
	}
	return step
}

func (s *chunkHeightSnapshot) heightAtLocal(lx, lz int) float32 {
	if lx < 0 || lz < 0 || lx >= s.W || lz >= s.H {
		return 0
	}
	return s.Heights[lz*s.W+lx]
}

// snapshotChunkHeights copies the height rectangle for chunk (cx,cz) into a detached snapshot.
func (t *TerrainObject) snapshotChunkHeights(cx, cz int) (*chunkHeightSnapshot, bool) {
	cs := t.ChunkSize
	if cs < 2 {
		return nil, false
	}
	idx := idx2(t, cx, cz)
	if idx < 0 || idx >= len(t.Chunks) {
		return nil, false
	}
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
		return nil, false
	}
	heights := make([]float32, w*h)
	for z := 0; z < h; z++ {
		for x := 0; x < w; x++ {
			heights[z*w+x] = t.heightAtCell(x0+x, z0+z)
		}
	}
	return &chunkHeightSnapshot{
		X0: x0, Z0: z0,
		W: w, H: h,
		Heights:      heights,
		CellSize:     t.CellSize,
		DetailFactor: t.DetailFactor,
		ScaleX:       t.scaleXEff(),
		ScaleY:       t.scaleYEff(),
		ScaleZ:       t.scaleZEff(),
	}, true
}

// heightmapPrep holds RGBA pixel data for GenMeshHeightmap plus mesh sizing (CPU-built).
type heightmapPrep struct {
	Pixels []byte
	W, H   int
	// Mesh extent passed to GenMeshHeightmap
	SizeX, SizeY, SizeZ float32
	MinH, MaxH          float32
}

// buildHeightmapPrepFromSnapshot builds grayscale heightmap pixels and bounds (no Raylib calls).
func buildHeightmapPrepFromSnapshot(s *chunkHeightSnapshot) (*heightmapPrep, bool) {
	w := s.W
	h := s.H
	if w < 2 || h < 2 {
		return nil, false
	}
	step := detailStepForTerrainDetail(s.DetailFactor)
	ws := (w + step - 1) / step
	hs := (h + step - 1) / step
	if ws < 2 {
		ws = 2
	}
	if hs < 2 {
		hs = 2
	}

	minH := s.heightAtLocal(0, 0)
	maxH := minH
	for z := 0; z < h; z++ {
		for x := 0; x < w; x++ {
			v := s.heightAtLocal(x, z)
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

	pix := make([]byte, ws*hs*4)
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
			v := s.heightAtLocal(sx, sz)
			nv := (v - minH) / dh
			if nv < 0 {
				nv = 0
			}
			if nv > 1 {
				nv = 1
			}
			b := uint8(nv * 255)
			o := (z*ws + x) * 4
			pix[o] = b
			pix[o+1] = b
			pix[o+2] = b
			pix[o+3] = 255
		}
	}

	sy := s.ScaleY
	if sy <= 0 {
		sy = 1
	}
	sx := s.ScaleX
	if sx <= 0 {
		sx = 1
	}
	sz := s.ScaleZ
	if sz <= 0 {
		sz = 1
	}
	sizeX := float32(ws-1) * s.CellSize * float32(step) * sx
	sizeZ := float32(hs-1) * s.CellSize * float32(step) * sz
	sizeY := dh * sy

	return &heightmapPrep{
		Pixels: pix,
		W:      ws,
		H:      hs,
		SizeX:  sizeX,
		SizeY:  sizeY,
		SizeZ:  sizeZ,
		MinH:   minH,
		MaxH:   maxH,
	}, true
}
