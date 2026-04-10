//go:build (cgo || (windows && !cgo)) && (!windows || !gopls_stub)

package terrain

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"

	mbmatrix "moonbasic/runtime/mbmatrix"
	"moonbasic/runtime"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

func registerTerrainExtended(m *Module, r runtime.Registrar) {
	r.Register("TERRAIN.LOAD", "terrain", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return terrainLoad(m, rt, args...) })
	r.Register("TERRAIN.GETNORMAL", "terrain", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return terrainGetNormal(m, rt, args...) })
	r.Register("TERRAIN.SETSCALE", "terrain", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return terrainSetScale(m, rt, args...) })
	r.Register("TERRAIN.GETSPLAT", "terrain", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return terrainGetSplat(m, rt, args...) })
	r.Register("TERRAIN.RAYCAST", "terrain", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return terrainRaycast(m, rt, args...) })
	r.Register("TERRAIN.SETDETAIL", "terrain", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return terrainSetDetail(m, rt, args...) })
	r.Register("Terrain.Load", "terrain", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) { return terrainLoad(m, rt, args...) })
}

func terrainLoad(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("TERRAIN.LOAD: heap not bound")
	}
	if len(args) < 1 || len(args) > 2 || args[0].Kind != value.KindString {
		return value.Nil, fmt.Errorf("TERRAIN.LOAD expects (heightmapPath$ [, diffusePath$])")
	}
	hm, err := rt.ArgString(args, 0)
	if err != nil {
		return value.Nil, err
	}
	diffuse := ""
	if len(args) == 2 {
		if args[1].Kind != value.KindString {
			return value.Nil, fmt.Errorf("TERRAIN.LOAD: diffuse path must be a string")
		}
		diffuse, err = rt.ArgString(args, 1)
		if err != nil {
			return value.Nil, err
		}
	}
	return m.loadTerrainFromPaths(hm, diffuse)
}

// loadTerrainFromPaths builds a terrain from a grayscale heightmap image and optional diffuse/splat texture.
func (m *Module) loadTerrainFromPaths(heightmapPath, diffusePath string) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("terrain load: heap not bound")
	}
	im := rl.LoadImage(heightmapPath)
	if im == nil || im.Data == nil {
		return value.Nil, fmt.Errorf("terrain load: failed to load heightmap %q", heightmapPath)
	}
	defer rl.UnloadImage(im)
	w := int(im.Width)
	h := int(im.Height)
	if w < 2 || h < 2 {
		return value.Nil, fmt.Errorf("terrain load: image too small")
	}
	cellSize := float32(1)
	cs := 64
	cw := (w + cs - 1) / cs
	ch := (h + cs - 1) / cs
	t := &TerrainObject{
		WorldW:        w,
		WorldH:        h,
		CellSize:      cellSize,
		ChunkSize:     cs,
		ChunkW:        cw,
		ChunkH:        ch,
		Heights:       make([]float32, w*h),
		Chunks:        make([]chunkSlot, cw*ch),
		StreamEnabled: true,
		LoadDist:      400,
		UnloadDist:    600,
		MaxHeight:     1,
		ScaleX:        1,
		ScaleY:        1,
		ScaleZ:        1,
		DetailFactor:  1,
	}
	var maxH float32
	for z := 0; z < h; z++ {
		for x := 0; x < w; x++ {
			c := rl.GetImageColor(*im, int32(x), int32(z))
			v := (float32(c.R) + float32(c.G) + float32(c.B)) / (3.0 * 255.0)
			hi := v * 100.0
			t.Heights[z*w+x] = hi
			if hi > maxH {
				maxH = hi
			}
		}
	}
	t.MaxHeight = maxH

	if diffusePath != "" {
		dim := rl.LoadImage(diffusePath)
		if dim != nil && dim.Data != nil {
			t.DiffuseTex = rl.LoadTextureFromImage(dim)
			t.DiffuseLoaded = true
			t.SplatImg = rl.ImageCopy(dim)
			rl.UnloadImage(dim)
		}
	}

	id, err := m.h.Alloc(t)
	if err != nil {
		if t.DiffuseLoaded {
			rl.UnloadTexture(t.DiffuseTex)
			t.DiffuseLoaded = false
		}
		if t.SplatImg != nil {
			rl.UnloadImage(t.SplatImg)
			t.SplatImg = nil
		}
		return value.Nil, err
	}
	m.active = id
	return value.FromHandle(int32(id)), nil
}

func terrainGetNormal(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("TERRAIN.GETNORMAL: heap not bound")
	}
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("TERRAIN.GETNORMAL expects (terrain, x#, z#)")
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
	nx, ny, nz := obj.NormalWorld(float32(x), float32(z))
	return mbmatrix.AllocVec3Value(m.h, nx, ny, nz)
}

func terrainSetScale(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 4 {
		return value.Nil, fmt.Errorf("TERRAIN.SETSCALE expects (terrain, x#, y#, z#)")
	}
	h, err := rt.ArgHandle(args, 0)
	if err != nil {
		return value.Nil, err
	}
	obj, err := castTerrain(m, heap.Handle(h))
	if err != nil {
		return value.Nil, err
	}
	sx, err := rt.ArgFloat(args, 1)
	if err != nil {
		return value.Nil, err
	}
	sy, err := rt.ArgFloat(args, 2)
	if err != nil {
		return value.Nil, err
	}
	sz, err := rt.ArgFloat(args, 3)
	if err != nil {
		return value.Nil, err
	}
	if sx <= 0 || sy <= 0 || sz <= 0 {
		return value.Nil, fmt.Errorf("TERRAIN.SETSCALE: scales must be > 0")
	}
	obj.ScaleX = float32(sx)
	obj.ScaleY = float32(sy)
	obj.ScaleZ = float32(sz)
	for i := range obj.Chunks {
		obj.Chunks[i].Dirty = true
	}
	return value.Nil, nil
}

func terrainGetSplat(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("TERRAIN.GETSPLAT expects (terrain, x#, z#)")
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
	v := obj.SplatAt(float32(x), float32(z))
	return value.FromInt(int64(v)), nil
}

func terrainRaycast(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("TERRAIN.RAYCAST: heap not bound")
	}
	if len(args) != 7 {
		return value.Nil, fmt.Errorf("TERRAIN.RAYCAST expects (terrain, ox#, oy#, oz#, dx#, dy#, dz#)")
	}
	h, err := rt.ArgHandle(args, 0)
	if err != nil {
		return value.Nil, err
	}
	obj, err := castTerrain(m, heap.Handle(h))
	if err != nil {
		return value.Nil, err
	}
	ox, err := rt.ArgFloat(args, 1)
	if err != nil {
		return value.Nil, err
	}
	oy, err := rt.ArgFloat(args, 2)
	if err != nil {
		return value.Nil, err
	}
	oz, err := rt.ArgFloat(args, 3)
	if err != nil {
		return value.Nil, err
	}
	dx, err := rt.ArgFloat(args, 4)
	if err != nil {
		return value.Nil, err
	}
	dy, err := rt.ArgFloat(args, 5)
	if err != nil {
		return value.Nil, err
	}
	dz, err := rt.ArgFloat(args, 6)
	if err != nil {
		return value.Nil, err
	}
	const maxD = float64(1e5)
	hit, hx, hy, hz := obj.RaycastTerrain(float32(ox), float32(oy), float32(oz), float32(dx), float32(dy), float32(dz), float32(maxD))
	arr, err := heap.NewArrayOfKind([]int64{4}, heap.ArrayKindFloat, 0)
	if err != nil {
		return value.Nil, err
	}
	if hit {
		arr.Floats[0] = 1
	} else {
		arr.Floats[0] = 0
	}
	arr.Floats[1] = float64(hx)
	arr.Floats[2] = float64(hy)
	arr.Floats[3] = float64(hz)
	hid, err := m.h.Alloc(arr)
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(hid), nil
}

func terrainSetDetail(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("TERRAIN.SETDETAIL expects (terrain, density#)")
	}
	h, err := rt.ArgHandle(args, 0)
	if err != nil {
		return value.Nil, err
	}
	obj, err := castTerrain(m, heap.Handle(h))
	if err != nil {
		return value.Nil, err
	}
	d, err := rt.ArgFloat(args, 1)
	if err != nil {
		return value.Nil, err
	}
	df := float32(d)
	if df <= 0 || df > 1 {
		return value.Nil, fmt.Errorf("TERRAIN.SETDETAIL: density must be in (0, 1]")
	}
	obj.DetailFactor = df
	for i := range obj.Chunks {
		obj.Chunks[i].Dirty = true
	}
	return value.Nil, nil
}
