package mbnav

import (
	prioheap "container/heap"
	"math"

	"moonbasic/vm/heap"
)

const (
	btKindSeq = iota
	btKindCond
	btKindAct
)

type navObj struct {
	cell    float64
	ox, oz  float64
	gw, gh  int
	blocked []bool
	groundY []float32
	built   bool
}

func (n *navObj) TypeName() string { return "Nav" }

func (n *navObj) TypeTag() uint16 { return heap.TagNav }

func (n *navObj) Free() {}

func (n *navObj) idx(ix, iz int) int { return iz*n.gw + ix }

func (n *navObj) containsCell(ix, iz int) bool {
	return ix >= 0 && iz >= 0 && ix < n.gw && iz < n.gh
}

func (n *navObj) worldToCell(wx, wz float64) (ix, iz int, ok bool) {
	if n.cell <= 0 || n.gw <= 0 || n.gh <= 0 {
		return 0, 0, false
	}
	ix = int(math.Floor((wx - n.ox) / n.cell))
	iz = int(math.Floor((wz - n.oz) / n.cell))
	if !n.containsCell(ix, iz) {
		return ix, iz, false
	}
	return ix, iz, true
}

func (n *navObj) cellCenter(ix, iz int) (x, z float64) {
	return n.ox + (float64(ix)+0.5)*n.cell, n.oz + (float64(iz)+0.5)*n.cell
}

func (n *navObj) heightAt(ix, iz int) float32 {
	if !n.containsCell(ix, iz) {
		return 0
	}
	return n.groundY[n.idx(ix, iz)]
}

// setBlockedRect blocks [ix0,ix1] x [iz0,iz1] inclusive in grid coords (clamped).
func (n *navObj) setBlockedRect(ix0, iz0, ix1, iz1 int, v bool) {
	if ix0 > ix1 {
		ix0, ix1 = ix1, ix0
	}
	if iz0 > iz1 {
		iz0, iz1 = iz1, iz0
	}
	for iz := iz0; iz <= iz1; iz++ {
		for ix := ix0; ix <= ix1; ix++ {
			if n.containsCell(ix, iz) {
				n.blocked[n.idx(ix, iz)] = v
			}
		}
	}
}

// setOpenRect sets walkable (unblocked) for XZ world AABB footprint.
func (n *navObj) setOpenRect(wx0, wz0, wx1, wz1 float64) {
	if wx0 > wx1 {
		wx0, wx1 = wx1, wx0
	}
	if wz0 > wz1 {
		wz0, wz1 = wz1, wz0
	}
	ix0, iz0, _ := n.worldToCell(wx0, wz0)
	ix1, iz1, _ := n.worldToCell(wx1, wz1)
	ix0 = clampInt(ix0, 0, n.gw-1)
	ix1 = clampInt(ix1, 0, n.gw-1)
	iz0 = clampInt(iz0, 0, n.gh-1)
	iz1 = clampInt(iz1, 0, n.gh-1)
	if ix0 > ix1 {
		ix0, ix1 = ix1, ix0
	}
	if iz0 > iz1 {
		iz0, iz1 = iz1, iz0
	}
	n.setBlockedRect(ix0, iz0, ix1, iz1, false)
}

func (n *navObj) setGroundYRect(wx0, wz0, wx1, wz1 float64, y float32) {
	if wx0 > wx1 {
		wx0, wx1 = wx1, wx0
	}
	if wz0 > wz1 {
		wz0, wz1 = wz1, wz0
	}
	ix0, iz0, ok0 := n.worldToCell(wx0, wz0)
	ix1, iz1, ok1 := n.worldToCell(wx1, wz1)
	if !ok0 || !ok1 {
		// clamp to grid
		ix0 = clampInt(ix0, 0, n.gw-1)
		ix1 = clampInt(ix1, 0, n.gw-1)
		iz0 = clampInt(iz0, 0, n.gh-1)
		iz1 = clampInt(iz1, 0, n.gh-1)
	}
	if ix0 > ix1 {
		ix0, ix1 = ix1, ix0
	}
	if iz0 > iz1 {
		iz0, iz1 = iz1, iz0
	}
	for iz := iz0; iz <= iz1; iz++ {
		for ix := ix0; ix <= ix1; ix++ {
			if n.containsCell(ix, iz) {
				n.groundY[n.idx(ix, iz)] = y
			}
		}
	}
}

func clampInt(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

type pathObj struct {
	valid bool
	pts   []pathPt
}

type pathPt struct {
	x, y, z float64
}

func (p *pathObj) TypeName() string { return "Path" }

func (p *pathObj) TypeTag() uint16 { return heap.TagPath }

func (p *pathObj) Free() {}

type navAgentObj struct {
	navH     heap.Handle
	x, y, z  float64
	vx, vy, vz float64
	speed    float64
	maxForce float64

	way       []pathPt
	wayIdx    int
	destX, destY, destZ float64
	hasDest   bool
	arriveEps float64
	rotY      float64
	manualRot bool
}

func (a *navAgentObj) TypeName() string { return "NavAgent" }

func (a *navAgentObj) TypeTag() uint16 { return heap.TagNavAgent }

func (a *navAgentObj) Free() {
	a.way = nil
}

// movementDirection returns world-space direction for facing: toward the active
// waypoint while pathing, otherwise the stored velocity (steering integration).
func (a *navAgentObj) movementDirection() (dx, dy, dz float64) {
	if len(a.way) > 0 && a.wayIdx < len(a.way) {
		t := a.way[a.wayIdx]
		return t.x - a.x, t.y - a.y, t.z - a.z
	}
	return a.vx, a.vy, a.vz
}

type steerGroupObj struct {
	agents []heap.Handle
}

func (g *steerGroupObj) TypeName() string { return "SteerGroup" }

func (g *steerGroupObj) TypeTag() uint16 { return heap.TagSteerGroup }

func (g *steerGroupObj) Free() { g.agents = nil }

type btNode struct {
	kind  int
	fn    string
	kids  []*btNode
}

type btObj struct {
	root *btNode
}

func (b *btObj) TypeName() string { return "BTree" }

func (b *btObj) TypeTag() uint16 { return heap.TagBTree }

func (b *btObj) Free() {
	b.root = nil
}

// --- A* ---

type astarNode struct {
	ix, iz int
	g, f   float64
	px, pz int // parent
	open   bool
	closed bool
}

type openPQ []*astarNode

func (h openPQ) Len() int           { return len(h) }
func (h openPQ) Less(i, j int) bool { return h[i].f < h[j].f }
func (h openPQ) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *openPQ) Push(x any) { *h = append(*h, x.(*astarNode)) }

func (h *openPQ) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

func findPathNav(n *navObj, sx, sy, sz, tx, ty, tz float64) *pathObj {
	out := &pathObj{valid: false}
	if n == nil || n.gw <= 0 || n.gh <= 0 || n.cell <= 0 {
		return out
	}
	six, siz, okS := n.worldToCell(sx, sz)
	tix, tiz, okT := n.worldToCell(tx, tz)
	if !okS || !okT {
		return out
	}
	if n.blocked[n.idx(six, siz)] || n.blocked[n.idx(tix, tiz)] {
		return out
	}

	nodes := make([]astarNode, n.gw*n.gh)
	for iz := 0; iz < n.gh; iz++ {
		for ix := 0; ix < n.gw; ix++ {
			i := n.idx(ix, iz)
			nodes[i].ix = ix
			nodes[i].iz = iz
			nodes[i].g = 1e30
			nodes[i].f = 1e30
			nodes[i].px = -1
			nodes[i].pz = -1
		}
	}

	start := &nodes[n.idx(six, siz)]

	hCost := func(ix, iz int) float64 {
		dx := float64(ix - tix)
		dz := float64(iz - tiz)
		return math.Sqrt(dx*dx+dz*dz) * n.cell
	}

	start.g = 0
	start.f = hCost(six, siz)

	var oh openPQ
	prioheap.Init(&oh)
	prioheap.Push(&oh, start)
	start.open = true

	neighbors := [8][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}, {-1, -1}, {1, -1}, {-1, 1}, {1, 1}}
	costs := [8]float64{1, 1, 1, 1, math.Sqrt2, math.Sqrt2, math.Sqrt2, math.Sqrt2}

	for oh.Len() > 0 {
		cur := prioheap.Pop(&oh).(*astarNode)
		cur.open = false
		cur.closed = true
		if cur.ix == tix && cur.iz == tiz {
			// reconstruct
			var chain []pathPt
			for c := cur; ; {
				cx, cz := n.cellCenter(c.ix, c.iz)
				y := float64(n.heightAt(c.ix, c.iz))
				if c.ix == six && c.iz == siz {
					y = sy
				}
				if c.ix == tix && c.iz == tiz {
					y = ty
				}
				chain = append(chain, pathPt{x: cx, y: y, z: cz})
				if c.px < 0 {
					break
				}
				c = &nodes[n.idx(c.px, c.pz)]
			}
			// reverse
			for i, j := 0, len(chain)-1; i < j; i, j = i+1, j-1 {
				chain[i], chain[j] = chain[j], chain[i]
			}
			out.pts = chain
			out.valid = true
			return out
		}

		for ni := range neighbors {
			dx, dz := neighbors[ni][0], neighbors[ni][1]
			nix, niz := cur.ix+dx, cur.iz+dz
			if !n.containsCell(nix, niz) {
				continue
			}
			if n.blocked[n.idx(nix, niz)] {
				continue
			}
			// corner cutting: for diagonal, both adjacent cardinals must be free
			if dx != 0 && dz != 0 {
				if n.blocked[n.idx(cur.ix+dx, cur.iz)] || n.blocked[n.idx(cur.ix, cur.iz+dz)] {
					continue
				}
			}
			nb := &nodes[n.idx(nix, niz)]
			if nb.closed {
				continue
			}
			tent := cur.g + costs[ni]*n.cell
			if tent < nb.g {
				nb.px = cur.ix
				nb.pz = cur.iz
				nb.g = tent
				nb.f = tent + hCost(nix, niz)
				if !nb.open {
					nb.open = true
					prioheap.Push(&oh, nb)
				}
			}
		}
	}
	return out
}
