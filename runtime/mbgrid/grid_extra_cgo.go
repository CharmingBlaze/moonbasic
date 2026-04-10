//go:build cgo || (windows && !cgo)

package mbgrid

import (
	"fmt"
	"math"

	"moonbasic/runtime"
	mbmodel3d "moonbasic/runtime/mbmodel3d"
	"moonbasic/runtime/terrain"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func (m *Module) gridGetPath(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	h, err := m.requireHeap(rt)
	if err != nil {
		return value.Nil, err
	}
	if len(args) != 5 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("GRID.GETPATH expects (grid, sx#, sz#, ex#, ez#)")
	}
	g, err := heap.Cast[*gridObj](h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	sx, ok1 := argF64(args[1])
	sz, ok2 := argF64(args[2])
	ex, ok3 := argF64(args[3])
	ez, ok4 := argF64(args[4])
	if !ok1 || !ok2 || !ok3 || !ok4 {
		return value.Nil, fmt.Errorf("GRID.GETPATH: coordinates must be numeric")
	}
	six := int(math.Floor((sx - g.ox) / g.cell))
	siz := int(math.Floor((sz - g.oz) / g.cell))
	tix := int(math.Floor((ex - g.ox) / g.cell))
	tiz := int(math.Floor((ez - g.oz) / g.cell))
	chain := findPath(g, six, siz, tix, tiz)
	if chain == nil {
		arr, err := heap.NewArray([]int64{0})
		if err != nil {
			return value.Nil, err
		}
		id, err := h.Alloc(arr)
		if err != nil {
			return value.Nil, err
		}
		return value.FromHandle(id), nil
	}
	n := len(chain) * 2
	dims := []int64{int64(n)}
	arr, err := heap.NewArray(dims)
	if err != nil {
		return value.Nil, err
	}
	for i, ci := range chain {
		ix := ci % g.gw
		iz := ci / g.gw
		_ = arr.Set([]int64{int64(i * 2)}, float64(ix))
		_ = arr.Set([]int64{int64(i*2 + 1)}, float64(iz))
	}
	id, err := h.Alloc(arr)
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}

func (m *Module) gridFollowTerrain(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	h, err := m.requireHeap(rt)
	if err != nil {
		return value.Nil, err
	}
	if len(args) != 2 || args[0].Kind != value.KindHandle || args[1].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("GRID.FOLLOWTERRAIN expects (grid, terrain)")
	}
	g, err := heap.Cast[*gridObj](h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	th := heap.Handle(args[1].IVal)
	g.cellY = make([]float32, len(g.cells))
	g.terrainRef = th
	for iz := 0; iz < g.gd; iz++ {
		for ix := 0; ix < g.gw; ix++ {
			cx := float32(g.ox + (float64(ix)+0.5)*g.cell)
			cz := float32(g.oz + (float64(iz)+0.5)*g.cell)
			g.cellY[g.idx(ix, iz)] = terrain.HeightWorldPublic(h, th, cx, cz)
		}
	}
	return value.Nil, nil
}

func (m *Module) gridPlaceEntity(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	h, err := m.requireHeap(rt)
	if err != nil {
		return value.Nil, err
	}
	if len(args) != 4 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("GRID.PLACEENTITY expects (grid, ix#, iz#, entity#)")
	}
	g, err := heap.Cast[*gridObj](h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	ix, ok1 := args[1].ToInt()
	iz, ok2 := args[2].ToInt()
	eid, ok3 := args[3].ToInt()
	if !ok1 || !ok2 || !ok3 || eid < 1 {
		return value.Nil, fmt.Errorf("GRID.PLACEENTITY: invalid arguments")
	}
	if !g.contains(int(ix), int(iz)) {
		return value.Nil, fmt.Errorf("GRID.PLACEENTITY: out of bounds")
	}
	g.occupant[g.idx(int(ix), int(iz))] = eid
	return value.Nil, nil
}

func (m *Module) gridRaycast(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	h, err := m.requireHeap(rt)
	if err != nil {
		return value.Nil, err
	}
	if len(args) != 3 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("GRID.RAYCAST expects (grid, screenX#, screenY#)")
	}
	g, err := heap.Cast[*gridObj](h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	sx, ok1 := argF64(args[1])
	sy, ok2 := argF64(args[2])
	if !ok1 || !ok2 {
		return value.Nil, fmt.Errorf("GRID.RAYCAST: screen coords must be numeric")
	}
	cam, okc := mbmodel3d.ActiveCamera3D()
	if !okc {
		return value.Nil, fmt.Errorf("GRID.RAYCAST: no active 3D camera")
	}
	rw := float32(rl.GetRenderWidth())
	rh := float32(rl.GetRenderHeight())
	ray := rl.GetScreenToWorldRayEx(rl.Vector2{X: float32(sx), Y: float32(sy)}, cam, int32(rw), int32(rh))
	o := ray.Position
	d := ray.Direction
	mk := func(ix, iz int) (value.Value, error) {
		arr, err := heap.NewArray([]int64{2})
		if err != nil {
			return value.Nil, err
		}
		_ = arr.Set([]int64{0}, float64(ix))
		_ = arr.Set([]int64{1}, float64(iz))
		id, err := h.Alloc(arr)
		if err != nil {
			return value.Nil, err
		}
		return value.FromHandle(id), nil
	}
	if math.Abs(float64(d.Y)) < 1e-5 {
		return mk(-1, -1)
	}
	t := -float64(o.Y) / float64(d.Y)
	if t < 0 {
		return mk(-1, -1)
	}
	px := float64(o.X) + float64(d.X)*t
	pz := float64(o.Z) + float64(d.Z)*t
	ix := int(math.Floor((px - g.ox) / g.cell))
	iz := int(math.Floor((pz - g.oz) / g.cell))
	if !g.contains(ix, iz) {
		return mk(-1, -1)
	}
	return mk(ix, iz)
}

func (m *Module) gridGetNeighbors(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	h, err := m.requireHeap(rt)
	if err != nil {
		return value.Nil, err
	}
	if len(args) != 4 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("GRID.GETNEIGHBORS expects (grid, ix#, iz#, radius#)")
	}
	g, err := heap.Cast[*gridObj](h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	cix, ok1 := args[1].ToInt()
	ciz, ok2 := args[2].ToInt()
	rad, ok3 := args[3].ToInt()
	if !ok1 || !ok2 || !ok3 || rad < 0 {
		return value.Nil, fmt.Errorf("GRID.GETNEIGHBORS: invalid arguments")
	}
	var ids []int64
	r := int(rad)
	for dz := -r; dz <= r; dz++ {
		for dx := -r; dx <= r; dx++ {
			ix := int(cix) + dx
			iz := int(ciz) + dz
			if !g.contains(ix, iz) {
				continue
			}
			eid := g.occupant[g.idx(ix, iz)]
			if eid > 0 {
				ids = append(ids, eid)
			}
		}
	}
	dims := []int64{int64(len(ids))}
	arr, err := heap.NewArray(dims)
	if err != nil {
		return value.Nil, err
	}
	for i, id := range ids {
		_ = arr.Set([]int64{int64(i)}, float64(id))
	}
	idh, err := h.Alloc(arr)
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(idh), nil
}
