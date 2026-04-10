//go:build cgo || (windows && !cgo)

package mbnav

import (
	"fmt"
	"math"

	"moonbasic/runtime"
	mbtime "moonbasic/runtime/time"
	"moonbasic/runtime/terrain"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

// navBakeTerrain builds a coarse walkability grid from a heightmap terrain (NAV.BAKE).
func (m *Module) navBakeTerrain(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	h, err := m.requireHeap(rt)
	if err != nil {
		return value.Nil, err
	}
	if len(args) != 3 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("NAV.BAKE expects (terrain, agentRadius#, maxSlopeDegrees#)")
	}
	th := heap.Handle(args[0].IVal)
	agentR, ok1 := argF64(args[1])
	maxSlope, ok2 := argF64(args[2])
	if !ok1 || !ok2 || agentR <= 0 || maxSlope < 0 {
		return value.Nil, fmt.Errorf("NAV.BAKE: agentRadius and maxSlope must be valid numbers (radius > 0)")
	}
	wx0, wz0, wx1, wz1, okb := terrain.WorldXZBounds(h, th)
	if !okb {
		return value.Nil, fmt.Errorf("NAV.BAKE: invalid terrain handle")
	}

	cell := math.Max(agentR*2, 0.5)
	gw := int(math.Ceil((wx1-wx0)/cell)) + 1
	gh := int(math.Ceil((wz1-wz0)/cell)) + 1
	const maxDim = 512
	if gw > maxDim || gh > maxDim {
		scale := math.Max(float64(gw)/maxDim, float64(gh)/maxDim)
		cell *= scale
		gw = int(math.Ceil((wx1-wx0)/cell)) + 1
		gh = int(math.Ceil((wz1-wz0)/cell)) + 1
	}
	if gw < 2 || gh < 2 {
		return value.Nil, fmt.Errorf("NAV.BAKE: computed grid too small")
	}

	n := &navObj{
		cell:    cell,
		ox:      wx0,
		oz:      wz0,
		gw:      gw,
		gh:      gh,
		blocked: make([]bool, gw*gh),
		groundY: make([]float32, gw*gh),
		built:   true,
	}

	for iz := 0; iz < gh; iz++ {
		for ix := 0; ix < gw; ix++ {
			cx, cz := n.cellCenter(ix, iz)
			_, _, ok := terrain.GridXZPublic(h, th, float32(cx), float32(cz))
			idx := n.idx(ix, iz)
			if !ok {
				n.blocked[idx] = true
				n.groundY[idx] = -1000
				continue
			}
			sd := float64(terrain.SlopeDegPublic(h, th, float32(cx), float32(cz)))
			if sd > maxSlope {
				n.blocked[idx] = true
				n.groundY[idx] = -1000
				continue
			}
			n.blocked[idx] = false
			n.groundY[idx] = terrain.HeightWorldPublic(h, th, float32(cx), float32(cz))
		}
	}

	navID, err := h.Alloc(n)
	if err != nil {
		return value.Nil, err
	}

	m.mu.Lock()
	if old, ok := m.terrainNav[th]; ok {
		_ = h.Free(old)
	}
	m.terrainNav[th] = navID
	m.mu.Unlock()

	return value.FromHandle(int32(navID)), nil
}

func (m *Module) navGetPathTerrain(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	h, err := m.requireHeap(rt)
	if err != nil {
		return value.Nil, err
	}
	if len(args) != 5 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("NAV.GETPATH expects (terrain, startX#, startZ#, endX#, endZ#)")
	}
	th := heap.Handle(args[0].IVal)
	sx, ok1 := argF64(args[1])
	sz, ok2 := argF64(args[2])
	ex, ok3 := argF64(args[3])
	ez, ok4 := argF64(args[4])
	if !ok1 || !ok2 || !ok3 || !ok4 {
		return value.Nil, fmt.Errorf("NAV.GETPATH: coordinates must be numeric")
	}

	m.mu.Lock()
	nh, ok := m.terrainNav[th]
	m.mu.Unlock()
	if !ok {
		return value.Nil, fmt.Errorf("NAV.GETPATH: call NAV.BAKE on this terrain first")
	}
	n, err := heap.Cast[*navObj](h, nh)
	if err != nil {
		return value.Nil, err
	}
	sy := float64(terrain.HeightWorldPublic(h, th, float32(sx), float32(sz)))
	ty := float64(terrain.HeightWorldPublic(h, th, float32(ex), float32(ez)))
	p := findPathNav(n, sx, sy, sz, ex, ty, ez)
	pid, err := h.Alloc(p)
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(int32(pid)), nil
}

func (m *Module) navIsReachableTerrain(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	h, err := m.requireHeap(rt)
	if err != nil {
		return value.Nil, err
	}
	if len(args) != 5 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("NAV.ISREACHABLE expects (terrain, startX#, startZ#, endX#, endZ#)")
	}
	th := heap.Handle(args[0].IVal)
	sx, ok1 := argF64(args[1])
	sz, ok2 := argF64(args[2])
	ex, ok3 := argF64(args[3])
	ez, ok4 := argF64(args[4])
	if !ok1 || !ok2 || !ok3 || !ok4 {
		return value.Nil, fmt.Errorf("NAV.ISREACHABLE: coordinates must be numeric")
	}

	m.mu.Lock()
	nh, ok := m.terrainNav[th]
	m.mu.Unlock()
	if !ok {
		return value.FromBool(false), nil
	}
	n, err := heap.Cast[*navObj](h, nh)
	if err != nil {
		return value.FromBool(false), nil
	}
	sy := float64(terrain.HeightWorldPublic(h, th, float32(sx), float32(sz)))
	ty := float64(terrain.HeightWorldPublic(h, th, float32(ex), float32(ez)))
	p := findPathNav(n, sx, sy, sz, ex, ty, ez)
	return value.FromBool(p != nil && p.valid && len(p.pts) > 0), nil
}

func (m *Module) enemyFollowPath(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	h, err := m.requireHeap(rt)
	if err != nil {
		return value.Nil, err
	}
	if m.ent == nil {
		return value.Nil, fmt.Errorf("ENEMY.FOLLOWPATH: entity module not bound")
	}
	if len(args) != 3 || args[1].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("ENEMY.FOLLOWPATH expects (entity#, path#, speed#)")
	}
	eid, ok := args[0].ToInt()
	if !ok || eid < 1 {
		return value.Nil, fmt.Errorf("ENEMY.FOLLOWPATH: invalid entity")
	}
	sp, ok := argF64(args[2])
	if !ok || sp <= 0 {
		return value.Nil, fmt.Errorf("ENEMY.FOLLOWPATH: speed must be > 0")
	}
	pathH := heap.Handle(args[1].IVal)
	po, err := heap.Cast[*pathObj](h, pathH)
	if err != nil {
		return value.Nil, err
	}
	if !po.valid || len(po.pts) == 0 {
		return value.Nil, nil
	}

	m.mu.Lock()
	st := m.enemyFollow[eid]
	if st.pathH != pathH {
		st = enemyFollowState{pathH: pathH, idx: 0}
		m.enemyFollow[eid] = st
	}
	idx := st.idx
	if idx >= len(po.pts) {
		m.mu.Unlock()
		return value.Nil, nil
	}
	target := po.pts[idx]
	m.mu.Unlock()

	px, py, pz, pok := m.ent.PlayerBridgeWorldPos(eid)
	if !pok {
		return value.Nil, fmt.Errorf("ENEMY.FOLLOWPATH: unknown entity")
	}
	dt := mbtime.DeltaSeconds(rt)

	dx := target.x - px
	dy := target.y - py
	dz := target.z - pz
	dist := math.Sqrt(dx*dx + dy*dy + dz*dz)
	step := sp * dt
	if dist <= 0.01 {
		m.mu.Lock()
		st := m.enemyFollow[eid]
		if st.pathH == pathH {
			st.idx++
			m.enemyFollow[eid] = st
		}
		m.mu.Unlock()
		return value.Nil, nil
	}
	if dist <= step {
		m.ent.PlayerBridgeSetWorldPos(eid, float32(target.x), float32(target.y), float32(target.z))
		m.mu.Lock()
		st := m.enemyFollow[eid]
		if st.pathH == pathH {
			st.idx++
			m.enemyFollow[eid] = st
		}
		m.mu.Unlock()
		return value.Nil, nil
	}
	t := step / dist
	nx := px + dx*t
	ny := py + dy*t
	nz := pz + dz*t
	m.ent.PlayerBridgeSetWorldPos(eid, float32(nx), float32(ny), float32(nz))
	return value.Nil, nil
}
