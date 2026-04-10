package mbgrid

import (
	"moonbasic/vm/heap"
)

// gridObj is a flat tactical grid over the XZ plane (optionally with per-cell Y from terrain).
type gridObj struct {
	gw, gd     int
	cell       float64
	ox, oz     float64
	cells      []int32
	cellY      []float32 // optional, len gw*gd after FollowTerrain
	terrainRef heap.Handle
	// occupants maps linear index -> entity id (one slot; extend later)
	occupant []int64
}

func (g *gridObj) TypeName() string { return "Grid" }

func (g *gridObj) TypeTag() uint16 { return heap.TagTacticalGrid }

func (g *gridObj) Free() {}

func (g *gridObj) idx(ix, iz int) int { return iz*g.gw + ix }

func (g *gridObj) contains(ix, iz int) bool {
	return ix >= 0 && iz >= 0 && ix < g.gw && iz < g.gd
}
