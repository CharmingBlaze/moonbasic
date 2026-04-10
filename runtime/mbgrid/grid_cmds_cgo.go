//go:build cgo || (windows && !cgo)

package mbgrid

import (
	"fmt"
	"math"

	"moonbasic/runtime"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func (m *Module) gridCreate(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	h, err := m.requireHeap(rt)
	if err != nil {
		return value.Nil, err
	}
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("GRID.CREATE expects (width#, depth#, cellSize#)")
	}
	gw64, ok1 := args[0].ToInt()
	gd64, ok2 := args[1].ToInt()
	cs, ok3 := argF64(args[2])
	if !ok1 || !ok2 || !ok3 || gw64 < 1 || gd64 < 1 || cs <= 0 {
		return value.Nil, fmt.Errorf("GRID.CREATE: invalid dimensions")
	}
	gw, gd := int(gw64), int(gd64)
	if gw > 4096 || gd > 4096 {
		return value.Nil, fmt.Errorf("GRID.CREATE: grid too large")
	}
	g := &gridObj{
		gw: gw, gd: gd, cell: cs,
		cells:    make([]int32, gw*gd),
		cellY:    nil,
		occupant: make([]int64, gw*gd),
	}
	for i := range g.cells {
		g.cells[i] = 0
		g.occupant[i] = 0
	}
	id, err := h.Alloc(g)
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(int32(id)), nil
}

func (m *Module) gridFree(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	h, err := m.requireHeap(rt)
	if err != nil {
		return value.Nil, err
	}
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("GRID.FREE expects grid handle")
	}
	_ = h.Free(heap.Handle(args[0].IVal))
	return value.Nil, nil
}

func (m *Module) gridSetCell(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	h, err := m.requireHeap(rt)
	if err != nil {
		return value.Nil, err
	}
	if len(args) != 4 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("GRID.SETCELL expects (grid, ix#, iz#, type#)")
	}
	g, err := heap.Cast[*gridObj](h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	ix, ok1 := args[1].ToInt()
	iz, ok2 := args[2].ToInt()
	tv, ok3 := args[3].ToInt()
	if !ok1 || !ok2 || !ok3 {
		return value.Nil, fmt.Errorf("GRID.SETCELL: numeric ix, iz, type required")
	}
	if !g.contains(int(ix), int(iz)) {
		return value.Nil, fmt.Errorf("GRID.SETCELL: out of bounds")
	}
	g.cells[g.idx(int(ix), int(iz))] = int32(tv)
	return value.Nil, nil
}

func (m *Module) gridGetCell(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	h, err := m.requireHeap(rt)
	if err != nil {
		return value.Nil, err
	}
	if len(args) != 3 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("GRID.GETCELL expects (grid, ix#, iz#)")
	}
	g, err := heap.Cast[*gridObj](h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	ix, ok1 := args[1].ToInt()
	iz, ok2 := args[2].ToInt()
	if !ok1 || !ok2 {
		return value.Nil, fmt.Errorf("GRID.GETCELL: ix, iz must be numeric")
	}
	if !g.contains(int(ix), int(iz)) {
		return value.FromInt(0), nil
	}
	return value.FromInt(int64(g.cells[g.idx(int(ix), int(iz))])), nil
}

func (m *Module) gridWorldToCell(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	h, err := m.requireHeap(rt)
	if err != nil {
		return value.Nil, err
	}
	if len(args) != 3 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("GRID.WORLDTOCELL expects (grid, worldX#, worldZ#)")
	}
	g, err := heap.Cast[*gridObj](h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	wx, ok1 := argF64(args[1])
	wz, ok2 := argF64(args[2])
	if !ok1 || !ok2 {
		return value.Nil, fmt.Errorf("GRID.WORLDTOCELL: world coords must be numeric")
	}
	ix := int(math.Floor((wx - g.ox) / g.cell))
	iz := int(math.Floor((wz - g.oz) / g.cell))
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

func (m *Module) gridDraw(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	h, err := m.requireHeap(rt)
	if err != nil {
		return value.Nil, err
	}
	if len(args) < 2 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("GRID.DRAW expects (grid, r, g, b [, a])")
	}
	g, err := heap.Cast[*gridObj](h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	var ri, gi, bi, ai int64 = 255, 255, 255, 120
	if len(args) >= 4 {
		ri, _ = args[1].ToInt()
		gi, _ = args[2].ToInt()
		bi, _ = args[3].ToInt()
	}
	if len(args) >= 5 {
		ai, _ = args[4].ToInt()
	}
	col := rl.Color{R: uint8(ri), G: uint8(gi), B: uint8(bi), A: uint8(ai)}
	for iz := 0; iz <= g.gd; iz++ {
		z := float32(g.oz + float64(iz)*g.cell)
		x0 := float32(g.ox)
		x1 := float32(g.ox + float64(g.gw)*g.cell)
		y := float32(0.02)
		if len(g.cellY) == len(g.cells) {
			y = g.cellY[g.idx(0, gridMin(iz, g.gd-1))]
		}
		rl.DrawLine3D(rl.Vector3{X: x0, Y: y, Z: z}, rl.Vector3{X: x1, Y: y, Z: z}, col)
	}
	for ix := 0; ix <= g.gw; ix++ {
		x := float32(g.ox + float64(ix)*g.cell)
		z0 := float32(g.oz)
		z1 := float32(g.oz + float64(g.gd)*g.cell)
		y := float32(0.02)
		rl.DrawLine3D(rl.Vector3{X: x, Y: y, Z: z0}, rl.Vector3{X: x, Y: y, Z: z1}, col)
	}
	return value.Nil, nil
}

func gridMin(a, b int) int {
	if a < b {
		return a
	}
	return b
}
