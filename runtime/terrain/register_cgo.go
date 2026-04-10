//go:build cgo || (windows && !cgo)

package terrain

import (
	"fmt"

	"moonbasic/runtime"
	mbcamera "moonbasic/runtime/camera"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

func registerTerrain(m *Module, r runtime.Registrar) {
	r.Register("TERRAIN.MAKE", "terrain", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return terrainMake(m, rt, args...) })
	r.Register("TERRAIN.FREE", "terrain", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return terrainFree(m, rt, args...) })
	r.Register("TERRAIN.SETPOS", "terrain", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return terrainSetPos(m, rt, args...) })
	r.Register("TERRAIN.SETCHUNKSIZE", "terrain", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return terrainSetChunkSize(m, rt, args...) })
	r.Register("TERRAIN.FILLPERLIN", "terrain", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return terrainFillPerlin(m, rt, args...) })
	r.Register("TERRAIN.FILLFLAT", "terrain", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return terrainFillFlat(m, rt, args...) })
	r.Register("TERRAIN.GETHEIGHT", "terrain", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return terrainGetHeight(m, rt, args...) })
	r.Register("TERRAIN.GETSLOPE", "terrain", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return terrainGetSlope(m, rt, args...) })
	r.Register("TERRAIN.RAISE", "terrain", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return terrainRaise(m, rt, args...) })
	r.Register("TERRAIN.LOWER", "terrain", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return terrainLower(m, rt, args...) })
	r.Register("TERRAIN.DRAW", "terrain", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return terrainDraw(m, rt, args...) })
	r.Register("CHUNK.GENERATE", "chunk", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return chunkGenerate(m, rt, args...) })
	r.Register("CHUNK.COUNT", "chunk", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return chunkCount(m, rt, args...) })
	r.Register("CHUNK.SETRANGE", "chunk", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return chunkSetRange(m, rt, args...) })
	r.Register("CHUNK.ISLOADED", "chunk", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return chunkIsLoaded(m, rt, args...) })
	registerTerrainExtended(m, r)
	registerTerrainBlitzAliases(m, r)
	registerTerrainApply(m, r)
}

func castTerrain(m *Module, h heap.Handle) (*TerrainObject, error) {
	return heap.Cast[*TerrainObject](m.h, h)
}

func terrainMake(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("TERRAIN.MAKE: heap not bound")
	}
	if len(args) < 2 || len(args) > 3 {
		return value.Nil, fmt.Errorf("TERRAIN.MAKE expects worldW, worldH [, cellSize#]")
	}
	ww, err := rt.ArgInt(args, 0)
	if err != nil {
		return value.Nil, err
	}
	wh, err := rt.ArgInt(args, 1)
	if err != nil {
		return value.Nil, err
	}
	if ww < 2 || wh < 2 {
		return value.Nil, fmt.Errorf("TERRAIN.MAKE: dimensions must be >= 2")
	}
	cellSize := float32(1)
	if len(args) == 3 {
		cs, err := rt.ArgFloat(args, 2)
		if err != nil {
			return value.Nil, err
		}
		cellSize = float32(cs)
		if cellSize <= 0 {
			return value.Nil, fmt.Errorf("TERRAIN.MAKE: cellSize must be > 0")
		}
	}
	cs := 64
	cw := (int(ww) + cs - 1) / cs
	ch := (int(wh) + cs - 1) / cs
	t := &TerrainObject{
		WorldW:        int(ww),
		WorldH:        int(wh),
		CellSize:      cellSize,
		ChunkSize:     cs,
		Heights:       make([]float32, int(ww)*int(wh)),
		ChunkW:        cw,
		ChunkH:        ch,
		Chunks:        make([]chunkSlot, cw*ch),
		StreamEnabled: true,
		LoadDist:      400,
		UnloadDist:    600,
		MaxHeight:     100,
		ScaleX:        1,
		ScaleY:        1,
		ScaleZ:        1,
		DetailFactor:  1,
	}
	id, err := m.h.Alloc(t)
	if err != nil {
		return value.Nil, err
	}
	m.active = id
	return value.FromHandle(int32(id)), nil
}

func terrainFree(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if m.h == nil || len(args) != 1 {
		return value.Nil, fmt.Errorf("TERRAIN.FREE expects terrain handle")
	}
	h, err := rt.ArgHandle(args, 0)
	if err != nil {
		return value.Nil, err
	}
	m.h.Free(heap.Handle(h))
	if m.active == heap.Handle(h) {
		m.active = 0
	}
	return value.Nil, nil
}

func terrainSetPos(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 4 {
		return value.Nil, fmt.Errorf("TERRAIN.SETPOS expects terrain, x#, y#, z#")
	}
	h, err := rt.ArgHandle(args, 0)
	if err != nil {
		return value.Nil, err
	}
	obj, err := castTerrain(m, heap.Handle(h))
	if err != nil {
		return value.Nil, err
	}
	x, err := rt.ArgFloat(args, 1)
	if err != nil {
		return value.Nil, err
	}
	y, err := rt.ArgFloat(args, 2)
	if err != nil {
		return value.Nil, err
	}
	z, err := rt.ArgFloat(args, 3)
	if err != nil {
		return value.Nil, err
	}
	obj.PX = float32(x)
	obj.PY = float32(y)
	obj.PZ = float32(z)
	return value.Nil, nil
}

func terrainSetChunkSize(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("TERRAIN.SETCHUNKSIZE expects terrain, cells")
	}
	h, err := rt.ArgHandle(args, 0)
	if err != nil {
		return value.Nil, err
	}
	obj, err := castTerrain(m, heap.Handle(h))
	if err != nil {
		return value.Nil, err
	}
	n, err := rt.ArgInt(args, 1)
	if err != nil {
		return value.Nil, err
	}
	if n < 8 || n > 256 {
		return value.Nil, fmt.Errorf("TERRAIN.SETCHUNKSIZE: cells must be 8..256")
	}
	for i := range obj.Chunks {
		if obj.Chunks[i].Loaded {
			return value.Nil, fmt.Errorf("TERRAIN.SETCHUNKSIZE: unload terrain first")
		}
	}
	obj.ChunkSize = int(n)
	obj.ChunkW = (obj.WorldW + obj.ChunkSize - 1) / obj.ChunkSize
	obj.ChunkH = (obj.WorldH + obj.ChunkSize - 1) / obj.ChunkSize
	obj.Chunks = make([]chunkSlot, obj.ChunkW*obj.ChunkH)
	return value.Nil, nil
}

func terrainFillPerlin(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("TERRAIN.FILLPERLIN expects terrain, scale#, height#")
	}
	h, err := rt.ArgHandle(args, 0)
	if err != nil {
		return value.Nil, err
	}
	obj, err := castTerrain(m, heap.Handle(h))
	if err != nil {
		return value.Nil, err
	}
	sc, err := rt.ArgFloat(args, 1)
	if err != nil {
		return value.Nil, err
	}
	amp, err := rt.ArgFloat(args, 2)
	if err != nil {
		return value.Nil, err
	}
	fillPerlinHeights(obj.Heights, obj.WorldW, obj.WorldH, sc, float32(amp), 42)
	obj.MaxHeight = float32(amp)
	for i := range obj.Chunks {
		obj.Chunks[i].Dirty = true
	}
	return value.Nil, nil
}

func terrainFillFlat(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("TERRAIN.FILLFLAT expects terrain, height#")
	}
	h, err := rt.ArgHandle(args, 0)
	if err != nil {
		return value.Nil, err
	}
	obj, err := castTerrain(m, heap.Handle(h))
	if err != nil {
		return value.Nil, err
	}
	he, err := rt.ArgFloat(args, 1)
	if err != nil {
		return value.Nil, err
	}
	hf := float32(he)
	for i := range obj.Heights {
		obj.Heights[i] = hf
	}
	obj.MaxHeight = hf
	for i := range obj.Chunks {
		obj.Chunks[i].Dirty = true
	}
	return value.Nil, nil
}

func terrainGetHeight(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("TERRAIN.GETHEIGHT expects terrain, x#, z#")
	}
	h, err := rt.ArgHandle(args, 0)
	if err != nil {
		return value.Nil, err
	}
	obj, err := castTerrain(m, heap.Handle(h))
	if err != nil {
		return value.Nil, err
	}
	x, err := rt.ArgFloat(args, 1)
	if err != nil {
		return value.Nil, err
	}
	z, err := rt.ArgFloat(args, 2)
	if err != nil {
		return value.Nil, err
	}
	return value.FromFloat(float64(obj.HeightWorld(float32(x), float32(z)))), nil
}

func terrainGetSlope(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("TERRAIN.GETSLOPE expects terrain, x#, z#")
	}
	h, err := rt.ArgHandle(args, 0)
	if err != nil {
		return value.Nil, err
	}
	obj, err := castTerrain(m, heap.Handle(h))
	if err != nil {
		return value.Nil, err
	}
	x, err := rt.ArgFloat(args, 1)
	if err != nil {
		return value.Nil, err
	}
	z, err := rt.ArgFloat(args, 2)
	if err != nil {
		return value.Nil, err
	}
	return value.FromFloat(float64(obj.SlopeDeg(float32(x), float32(z)))), nil
}

func terrainRaise(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 5 {
		return value.Nil, fmt.Errorf("TERRAIN.RAISE expects terrain, x#, z#, radius#, amount#")
	}
	h, err := rt.ArgHandle(args, 0)
	if err != nil {
		return value.Nil, err
	}
	obj, err := castTerrain(m, heap.Handle(h))
	if err != nil {
		return value.Nil, err
	}
	wx, err := rt.ArgFloat(args, 1)
	if err != nil {
		return value.Nil, err
	}
	wz, err := rt.ArgFloat(args, 2)
	if err != nil {
		return value.Nil, err
	}
	rad, err := rt.ArgFloat(args, 3)
	if err != nil {
		return value.Nil, err
	}
	amt, err := rt.ArgFloat(args, 4)
	if err != nil {
		return value.Nil, err
	}
	brush(obj, float32(wx), float32(wz), float32(rad), float32(amt), 1)
	return value.Nil, nil
}

func terrainLower(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 5 {
		return value.Nil, fmt.Errorf("TERRAIN.LOWER expects terrain, x#, z#, radius#, amount#")
	}
	h, err := rt.ArgHandle(args, 0)
	if err != nil {
		return value.Nil, err
	}
	obj, err := castTerrain(m, heap.Handle(h))
	if err != nil {
		return value.Nil, err
	}
	wx, err := rt.ArgFloat(args, 1)
	if err != nil {
		return value.Nil, err
	}
	wz, err := rt.ArgFloat(args, 2)
	if err != nil {
		return value.Nil, err
	}
	rad, err := rt.ArgFloat(args, 3)
	if err != nil {
		return value.Nil, err
	}
	amt, err := rt.ArgFloat(args, 4)
	if err != nil {
		return value.Nil, err
	}
	brush(obj, float32(wx), float32(wz), float32(rad), float32(amt), -1)
	return value.Nil, nil
}

func brush(obj *TerrainObject, wx, wz, radius, amount float32, sign int) {
	sx := obj.scaleXEff()
	sz := obj.scaleZEff()
	cx0 := int((wx - obj.PX - radius) / (obj.CellSize * sx))
	cx1 := int((wx - obj.PX + radius) / (obj.CellSize * sx))
	cz0 := int((wz - obj.PZ - radius) / (obj.CellSize * sz))
	cz1 := int((wz - obj.PZ + radius) / (obj.CellSize * sz))
	if cx0 < 0 {
		cx0 = 0
	}
	if cz0 < 0 {
		cz0 = 0
	}
	if cx1 >= obj.WorldW {
		cx1 = obj.WorldW - 1
	}
	if cz1 >= obj.WorldH {
		cz1 = obj.WorldH - 1
	}
	r2 := radius * radius
	for z := cz0; z <= cz1; z++ {
		for x := cx0; x <= cx1; x++ {
			px := obj.PX + float32(x)*obj.CellSize*sx
			pz := obj.PZ + float32(z)*obj.CellSize*sz
			dx := px - wx
			dz := pz - wz
			if dx*dx+dz*dz <= r2 {
				i := z*obj.WorldW + x
				obj.Heights[i] += float32(sign) * amount
				if obj.Heights[i] < 0 {
					obj.Heights[i] = 0
				}
			}
		}
	}
	// mark chunks dirty in rect
	ccx0 := cx0 / obj.ChunkSize
	ccx1 := cx1 / obj.ChunkSize
	ccz0 := cz0 / obj.ChunkSize
	ccz1 := cz1 / obj.ChunkSize
	for czz := ccz0; czz <= ccz1; czz++ {
		for cxx := ccx0; cxx <= ccx1; cxx++ {
			if cxx >= 0 && czz >= 0 && cxx < obj.ChunkW && czz < obj.ChunkH {
				obj.Chunks[idx2(obj, cxx, czz)].Dirty = true
			}
		}
	}
}

func terrainDraw(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("TERRAIN.DRAW expects terrain")
	}
	h, err := rt.ArgHandle(args, 0)
	if err != nil {
		return value.Nil, err
	}
	obj, err := castTerrain(m, heap.Handle(h))
	if err != nil {
		return value.Nil, err
	}
	cs := obj.ChunkSize
	for cz := 0; cz < obj.ChunkH; cz++ {
		for cx := 0; cx < obj.ChunkW; cx++ {
			idx := idx2(obj, cx, cz)
			ch := &obj.Chunks[idx]
			if !ch.Loaded || ch.Mesh.VertexCount == 0 {
				continue
			}
			minH, maxH := ch.MinH, ch.MaxH
			if !ch.BoundsValid {
				x0 := cx * cs
				z0 := cz * cs
				x1 := x0 + cs
				z1 := z0 + cs
				if x1 > obj.WorldW {
					x1 = obj.WorldW
				}
				if z1 > obj.WorldH {
					z1 = obj.WorldH
				}
				minH = obj.heightAtCell(x0, z0)
				maxH = minH
				for z := z0; z <= z1; z++ {
					for x := x0; x <= x1; x++ {
						v := obj.heightAtCell(x, z)
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
			sx := obj.scaleXEff()
			sz := obj.scaleZEff()
			minX := obj.PX + float32(cx*cs)*obj.CellSize*sx
			minZ := obj.PZ + float32(cz*cs)*obj.CellSize*sz
			maxX := minX + float32(cs)*obj.CellSize*sx
			maxZ := minZ + float32(cs)*obj.CellSize*sz
			minY := obj.PY + minH
			maxY := obj.PY + maxH
			centX := (minX + maxX) * 0.5
			centZ := (minZ + maxZ) * 0.5
			if mbcamera.BehindHorizonActive(maxY, centX, centZ) {
				continue
			}
			if !mbcamera.WithinDistanceActive(centX, (minY+maxY)*0.5, centZ) {
				continue
			}
			if !mbcamera.AABBVisibleActive(minX, minY, minZ, maxX, maxY, maxZ) {
				continue
			}
			obj.drawChunk(cx, cz)
		}
	}
	return value.Nil, nil
}

func chunkGenerate(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("CHUNK.GENERATE expects terrain, chunkX, chunkZ")
	}
	h, err := rt.ArgHandle(args, 0)
	if err != nil {
		return value.Nil, err
	}
	obj, err := castTerrain(m, heap.Handle(h))
	if err != nil {
		return value.Nil, err
	}
	cx, err := rt.ArgInt(args, 1)
	if err != nil {
		return value.Nil, err
	}
	cz, err := rt.ArgInt(args, 2)
	if err != nil {
		return value.Nil, err
	}
	if int(cx) < 0 || int(cz) < 0 || int(cx) >= obj.ChunkW || int(cz) >= obj.ChunkH {
		return value.Nil, nil
	}
	obj.rebuildChunkMesh(int(cx), int(cz))
	return value.Nil, nil
}

func chunkCount(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("CHUNK.COUNT expects terrain")
	}
	h, err := rt.ArgHandle(args, 0)
	if err != nil {
		return value.Nil, err
	}
	obj, err := castTerrain(m, heap.Handle(h))
	if err != nil {
		return value.Nil, err
	}
	return value.FromInt(int64(obj.loadedChunkCount())), nil
}

func chunkSetRange(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("CHUNK.SETRANGE expects terrain, loadDist#, unloadDist#")
	}
	h, err := rt.ArgHandle(args, 0)
	if err != nil {
		return value.Nil, err
	}
	obj, err := castTerrain(m, heap.Handle(h))
	if err != nil {
		return value.Nil, err
	}
	ld, err := rt.ArgFloat(args, 1)
	if err != nil {
		return value.Nil, err
	}
	ud, err := rt.ArgFloat(args, 2)
	if err != nil {
		return value.Nil, err
	}
	if float32(ud) <= float32(ld) {
		return value.Nil, fmt.Errorf("CHUNK.SETRANGE: unload must be > load")
	}
	obj.LoadDist = float32(ld)
	obj.UnloadDist = float32(ud)
	return value.Nil, nil
}

func chunkIsLoaded(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("CHUNK.ISLOADED expects terrain, chunkX, chunkZ")
	}
	h, err := rt.ArgHandle(args, 0)
	if err != nil {
		return value.Nil, err
	}
	obj, err := castTerrain(m, heap.Handle(h))
	if err != nil {
		return value.Nil, err
	}
	cx, err := rt.ArgInt(args, 1)
	if err != nil {
		return value.Nil, err
	}
	cz, err := rt.ArgInt(args, 2)
	if err != nil {
		return value.Nil, err
	}
	if cx < 0 || cz < 0 || int(cx) >= obj.ChunkW || int(cz) >= obj.ChunkH {
		return value.FromBool(false), nil
	}
	return value.FromBool(obj.Chunks[idx2(obj, int(cx), int(cz))].Loaded), nil
}

// TickStreaming updates active terrain streaming (called from WORLD.UPDATE).
func (m *Module) TickStreaming(rt *runtime.Runtime) {
	if m.h == nil || m.active == 0 {
		return
	}
	obj, err := castTerrain(m, m.active)
	if err != nil {
		return
	}
	obj.TickStreaming()
}

// Preload wraps PreloadChunksAround for active terrain.
func (m *Module) Preload(rt *runtime.Runtime, radius int) {
	if m.h == nil || m.active == 0 {
		return
	}
	obj, err := castTerrain(m, m.active)
	if err != nil {
		return
	}
	obj.PreloadChunksAround(radius)
}

// PreloadTerrain preloads chunks for a specific terrain handle.
func (m *Module) PreloadTerrain(h heap.Handle, radius int) error {
	obj, err := castTerrain(m, h)
	if err != nil {
		return err
	}
	obj.PreloadChunksAround(radius)
	return nil
}

// SetCenter sets streaming center on active terrain.
func (m *Module) SetCenter(x, z float32) {
	if m.h == nil || m.active == 0 {
		return
	}
	obj, err := castTerrain(m, m.active)
	if err != nil {
		return
	}
	obj.CenterX = x
	obj.CenterZ = z
}

// SetStreamEnabled toggles streaming on active terrain.
func (m *Module) SetStreamEnabled(on bool) {
	if m.h == nil || m.active == 0 {
		return
	}
	obj, err := castTerrain(m, m.active)
	if err != nil {
		return
	}
	obj.StreamEnabled = on
}

// StatusString returns a short status for WORLD.STATUS.
func (m *Module) StatusString() string {
	if m.h == nil || m.active == 0 {
		return "no active terrain"
	}
	obj, err := castTerrain(m, m.active)
	if err != nil {
		return "invalid terrain"
	}
	n := obj.loadedChunkCount()
	return fmt.Sprintf("Chunks: %d loaded (stream=%v)", n, obj.StreamEnabled)
}

// IsReady returns true when all chunks within load distance are loaded for the given terrain.
func (m *Module) IsReadyTerrain(h heap.Handle) bool {
	if m.h == nil {
		return false
	}
	obj, err := castTerrain(m, h)
	if err != nil {
		return false
	}
	for cz := 0; cz < obj.ChunkH; cz++ {
		for cx := 0; cx < obj.ChunkW; cx++ {
			if chunkDistanceMeters(obj, cx, cz) <= float64(obj.LoadDist) {
				if !obj.Chunks[idx2(obj, cx, cz)].Loaded {
					return false
				}
			}
		}
	}
	return true
}
