//go:build (cgo || (windows && !cgo)) && (!windows || !gopls_stub)

package terrain

import (
	"fmt"

	"moonbasic/runtime"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

func registerTerrainBlitzAliases(m *Module, r runtime.Registrar) {
	r.Register("TerrainHeight", "terrain", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		return terrainGetHeight(m, rt, args...)
	})
	r.Register("ModifyTerrain", "terrain", m.terrainModifyHeight)
	r.Register("TerrainX", "terrain", m.terrainWorldToGridX)
	r.Register("TerrainZ", "terrain", m.terrainWorldToGridZ)
	r.Register("TerrainSize", "terrain", m.terrainSizeDims)
	r.Register("LoadTerrain", "terrain", m.terrainLoadHeightmap)
	r.Register("TerrainDetail", "terrain", m.terrainDetailStub)
	r.Register("TerrainShading", "terrain", m.terrainShadingStub)
}

// terrainModifyHeight sets one height sample near (x,z) — (terrain, x#, z#, height#, realtime#); realtime reserved.
func (m *Module) terrainModifyHeight(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) < 4 || len(args) > 5 {
		return value.Nil, fmt.Errorf("ModifyTerrain expects (terrain, x#, z#, height# [, realtime#])")
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
	ht, err := rt.ArgFloat(args, 3)
	if err != nil {
		return value.Nil, err
	}
	lx := int((float32(wx) - obj.PX) / (obj.CellSize * obj.scaleXEff()))
	lz := int((float32(wz) - obj.PZ) / (obj.CellSize * obj.scaleZEff()))
	if lx < 0 || lz < 0 || lx >= obj.WorldW || lz >= obj.WorldH {
		return value.Nil, fmt.Errorf("ModifyTerrain: out of grid")
	}
	sy := obj.scaleYEff()
	obj.Heights[lz*obj.WorldW+lx] = (float32(ht) - obj.PY) / sy
	raw := obj.Heights[lz*obj.WorldW+lx]
	obj.MaxHeight = maxFloatTerrain(obj.MaxHeight, raw)
	for i := range obj.Chunks {
		obj.Chunks[i].Dirty = true
	}
	return value.Nil, nil
}

func maxFloatTerrain(a, b float32) float32 {
	if b > a {
		return b
	}
	return a
}

func (m *Module) terrainWorldToGridX(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 4 {
		return value.Nil, fmt.Errorf("TerrainX expects (terrain, worldX#, worldY#, worldZ#)")
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
	_, _ = rt.ArgFloat(args, 2)
	wz, err := rt.ArgFloat(args, 3)
	if err != nil {
		return value.Nil, err
	}
	lx := (float32(wx) - obj.PX) / (obj.CellSize * obj.scaleXEff())
	_ = wz
	return value.FromFloat(float64(lx)), nil
}

func (m *Module) terrainWorldToGridZ(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 4 {
		return value.Nil, fmt.Errorf("TerrainZ expects (terrain, worldX#, worldY#, worldZ#)")
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
	_, _ = rt.ArgFloat(args, 2)
	wz, err := rt.ArgFloat(args, 3)
	if err != nil {
		return value.Nil, err
	}
	_ = wx
	lz := (float32(wz) - obj.PZ) / (obj.CellSize * obj.scaleZEff())
	return value.FromFloat(float64(lz)), nil
}

func (m *Module) terrainSizeDims(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("TerrainSize expects (terrain)")
	}
	h, err := rt.ArgHandle(args, 0)
	if err != nil {
		return value.Nil, err
	}
	obj, err := castTerrain(m, heap.Handle(h))
	if err != nil {
		return value.Nil, err
	}
	if m.h == nil {
		return value.Nil, runtime.Errorf("TerrainSize: heap not bound")
	}
	arr, err := heap.NewArray([]int64{2})
	if err != nil {
		return value.Nil, err
	}
	_ = arr.Set([]int64{0}, float64(obj.WorldW))
	_ = arr.Set([]int64{1}, float64(obj.WorldH))
	id, err := m.h.Alloc(arr)
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}

// LoadTerrain(path$ [, parent#]) — greyscale heightmap → terrain; parent reserved (entity parenting not wired).
func (m *Module) terrainLoadHeightmap(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("LoadTerrain: heap not bound")
	}
	if len(args) < 1 || len(args) > 2 || args[0].Kind != value.KindString {
		return value.Nil, fmt.Errorf("LoadTerrain expects (path$ [, parent#])")
	}
	path, err := rt.ArgString(args, 0)
	if err != nil {
		return value.Nil, err
	}
	if len(args) == 2 {
		_, _ = args[1].ToInt()
	}
	return m.loadTerrainFromPaths(path, "")
}

func (m *Module) terrainDetailStub(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) < 2 || len(args) > 3 {
		return value.Nil, fmt.Errorf("TerrainDetail expects (terrain, detailLevel# [, morph#])")
	}
	return terrainSetDetail(m, rt, args[0], args[1])
}

func (m *Module) terrainShadingStub(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("TerrainShading expects (terrain, state#)")
	}
	_, _ = rt, args
	return value.Nil, nil
}
