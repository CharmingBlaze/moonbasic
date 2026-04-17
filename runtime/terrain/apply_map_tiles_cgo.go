//go:build (cgo || (windows && !cgo)) && (!windows || !gopls_stub)

package terrain

import (
	"fmt"

	mbimage "moonbasic/runtime/mbimage"
	mbtilemap "moonbasic/runtime/mbtilemap"
	mbentity "moonbasic/runtime/mbentity"
	"moonbasic/runtime"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func init() {
	mbentity.GetTerrainHeightCb = HeightWorldPublic
}

func registerTerrainApply(m *Module, r runtime.Registrar) {
	r.Register("TERRAIN.APPLYMAP", "terrain", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		return terrainApplyMap(m, rt, args...)
	})
	r.Register("TERRAIN.APPLYTILES", "terrain", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		return terrainApplyTiles(m, rt, args...)
	})
}

// terrainApplyMap drapes a CPU image (splat / diffuse paint) onto the terrain: updates GPU diffuse and splat sampling, then rebuilds loaded chunk meshes.
func terrainApplyMap(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("TERRAIN.APPLYMAP: heap not bound")
	}
	if len(args) != 2 || args[0].Kind != value.KindHandle || args[1].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("TERRAIN.APPLYMAP expects (terrain, imageHandle)")
	}
	t, err := castTerrain(m, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	im, err := mbimage.RayImageForTexture(m.h, heap.Handle(args[1].IVal))
	if err != nil {
		return value.Nil, fmt.Errorf("TERRAIN.APPLYMAP: %w", err)
	}
	if im == nil || im.Data == nil {
		return value.Nil, fmt.Errorf("TERRAIN.APPLYMAP: empty image")
	}

	cpy := rl.ImageCopy(im)
	if cpy == nil || cpy.Data == nil {
		return value.Nil, fmt.Errorf("TERRAIN.APPLYMAP: ImageCopy failed")
	}
	defer rl.UnloadImage(cpy)

	if t.DiffuseLoaded {
		rl.UnloadTexture(t.DiffuseTex)
		t.DiffuseLoaded = false
	}
	t.DiffuseTex = rl.LoadTextureFromImage(cpy)
	t.DiffuseLoaded = true

	if t.SplatImg != nil {
		rl.UnloadImage(t.SplatImg)
		t.SplatImg = nil
	}
	t.SplatImg = rl.ImageCopy(im)

	for i := range t.Chunks {
		t.Chunks[i].Dirty = true
	}
	for cz := 0; cz < t.ChunkH; cz++ {
		for cx := 0; cx < t.ChunkW; cx++ {
			idx := idx2(t, cx, cz)
			if idx >= 0 && idx < len(t.Chunks) && t.Chunks[idx].Loaded {
				t.rebuildChunkMesh(cx, cz)
			}
		}
	}
	return args[0], nil
}

// terrainApplyTiles places a copy of templateEntity at each non-empty tile center on layer layerIndex (default 0).
// World XZ uses tilemap draw size (TILEMAP.SETTILESIZE); Y comes from terrain height at that XZ plus half template height heuristic (template Y preserved offset).
func terrainApplyTiles(m *Module, rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("TERRAIN.APPLYTILES: heap not bound")
	}
	if m.ent == nil {
		return value.Nil, fmt.Errorf("TERRAIN.APPLYTILES: entity module not bound")
	}
	if len(args) < 3 || len(args) > 4 {
		return value.Nil, fmt.Errorf("TERRAIN.APPLYTILES expects (terrain, tilemap, templateEntity# [, layerIndex#])")
	}
	if args[0].Kind != value.KindHandle || args[1].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("TERRAIN.APPLYTILES: terrain and tilemap handles required")
	}
	terr, err := castTerrain(m, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	layer := 0
	if len(args) == 4 {
		li, ok := args[3].ToInt()
		if !ok || li < 0 {
			return value.Nil, fmt.Errorf("TERRAIN.APPLYTILES: layerIndex must be >= 0")
		}
		layer = int(li)
	}
	tw, th, dw, dh, gids, err := mbtilemap.TileLayerGIDs(m.h, heap.Handle(args[1].IVal), layer)
	if err != nil {
		return value.Nil, err
	}
	tid, ok := args[2].ToInt()
	if !ok || tid < 1 {
		return value.Nil, fmt.Errorf("TERRAIN.APPLYTILES: invalid template entity")
	}

	var placed int64
	for y := int32(0); y < th; y++ {
		for x := int32(0); x < tw; x++ {
			iy := int(y)
			ix := int(x)
			if iy >= len(gids) || ix >= len(gids[iy]) {
				continue
			}
			gid := gids[iy][ix]
			if gid == 0 {
				continue
			}
			wx := terr.PX + float32(x)*float32(dw) + float32(dw)*0.5
			wz := terr.PZ + float32(y)*float32(dh) + float32(dh)*0.5
			hy := terr.HeightWorld(wx, wz)
			nid, err := m.ent.DuplicateEntityAt(tid, wx, hy, wz)
			if err != nil {
				return value.Nil, err
			}
			_ = nid
			placed++
		}
	}
	return value.FromInt(placed), nil
}
