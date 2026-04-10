//go:build cgo || (windows && !cgo)

package mbgrid

import (
	"fmt"

	"moonbasic/runtime"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

func (m *Module) gridSnap(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if m.ent == nil {
		return value.Nil, fmt.Errorf("GRID.SNAP: entity module not bound")
	}
	h, err := m.requireHeap(rt)
	if err != nil {
		return value.Nil, err
	}
	if len(args) != 4 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("GRID.SNAP expects (grid, entity#, ix#, iz#)")
	}
	g, err := heap.Cast[*gridObj](h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	eid, ok0 := args[1].ToInt()
	ix, ok1 := args[2].ToInt()
	iz, ok2 := args[3].ToInt()
	if !ok0 || !ok1 || !ok2 || eid < 1 {
		return value.Nil, fmt.Errorf("GRID.SNAP: invalid arguments")
	}
	if !g.contains(int(ix), int(iz)) {
		return value.Nil, fmt.Errorf("GRID.SNAP: cell out of bounds")
	}
	cx := float32(g.ox + (float64(ix)+0.5)*g.cell)
	cz := float32(g.oz + (float64(iz)+0.5)*g.cell)
	cy := float32(0)
	if len(g.cellY) == len(g.cells) {
		cy = g.cellY[g.idx(int(ix), int(iz))]
	}
	m.ent.PlayerBridgeSetWorldPos(eid, cx, cy, cz)
	return value.Nil, nil
}
