//go:build cgo || (windows && !cgo)

package mbnav

import (
	"fmt"
	"math"

	"moonbasic/runtime"
	"moonbasic/runtime/mbmatrix"
	"moonbasic/runtime/mbmodel3d"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

func (m *Module) Register(reg runtime.Registrar) {
	reg.Register("NAV.MAKE", "nav", m.navMake)
	reg.Register("NAV.FREE", "nav", m.navFree)
	reg.Register("NAV.SETGRID", "nav", m.navSetGrid)
	reg.Register("NAV.ADDTERRAIN", "nav", m.navAddTerrain)
	reg.Register("NAV.ADDOBSTACLE", "nav", m.navAddObstacle)
	reg.Register("NAV.BUILD", "nav", m.navBuild)
	reg.Register("NAV.FINDPATH", "nav", m.navFindPath)

	reg.Register("PATH.ISVALID", "path", m.pathIsValid)
	reg.Register("PATH.NODECOUNT", "path", m.pathNodeCount)
	reg.Register("PATH.NODEX", "path", m.pathNodeX)
	reg.Register("PATH.NODEY", "path", m.pathNodeY)
	reg.Register("PATH.NODEZ", "path", m.pathNodeZ)
	reg.Register("PATH.FREE", "path", m.pathFree)

	reg.Register("NAVAGENT.MAKE", "navagent", m.agentMake)
	reg.Register("NAVAGENT.FREE", "navagent", m.agentFree)
	reg.Register("NAVAGENT.SETPOS", "navagent", m.agentSetPos)
	reg.Register("NAVAGENT.SETSPEED", "navagent", m.agentSetSpeed)
	reg.Register("NAVAGENT.SETMAXFORCE", "navagent", m.agentSetMaxForce)
	reg.Register("NAVAGENT.APPLYFORCE", "navagent", m.agentApplyForce)
	reg.Register("NAVAGENT.MOVETO", "navagent", m.agentMoveTo)
	reg.Register("NAVAGENT.UPDATE", "navagent", m.agentUpdate)
	reg.Register("NAVAGENT.ISATDESTINATION", "navagent", m.agentIsAtDestination)
	reg.Register("NAVAGENT.X", "navagent", m.agentX)
	reg.Register("NAVAGENT.Y", "navagent", m.agentY)
	reg.Register("NAVAGENT.Z", "navagent", m.agentZ)

	reg.Register("STEER.GROUPMAKE", "steer", m.steerGroupMake)
	reg.Register("STEER.GROUPADD", "steer", m.steerGroupAdd)
	reg.Register("STEER.GROUPCLEAR", "steer", m.steerGroupClear)
	reg.Register("STEER.SEEK", "steer", m.steerSeek)
	reg.Register("STEER.FLEE", "steer", m.steerFlee)
	reg.Register("STEER.ARRIVE", "steer", m.steerArrive)
	reg.Register("STEER.WANDER", "steer", m.steerWander)
	reg.Register("STEER.FLOCK", "steer", m.steerFlock)
	reg.Register("STEER.AVOIDOBSTACLES", "steer", m.steerAvoidObstacles)
	reg.Register("STEER.FOLLOWPATH", "steer", m.steerFollowPath)

	reg.Register("BTREE.MAKE", "btree", m.btMake)
	reg.Register("BTREE.FREE", "btree", m.btFree)
	reg.Register("BTREE.SEQUENCE", "btree", m.btSequence)
	reg.Register("BTREE.ADDCONDITION", "btree", m.btAddCondition)
	reg.Register("BTREE.ADDACTION", "btree", m.btAddAction)
	reg.Register("BTREE.RUN", "btree", m.btRun)
}

func (m *Module) Shutdown() {}

func (m *Module) requireHeap(rt *runtime.Runtime) (*heap.Store, error) {
	if rt != nil && rt.Heap != nil {
		return rt.Heap, nil
	}
	if m.h != nil {
		return m.h, nil
	}
	return nil, fmt.Errorf("nav/ai: heap not bound")
}

func argF64(v value.Value) (float64, bool) {
	if f, ok := v.ToFloat(); ok {
		return f, true
	}
	if i, ok := v.ToInt(); ok {
		return float64(i), true
	}
	return 0, false
}

func valueTruthy(v value.Value) bool {
	if v.Kind == value.KindBool {
		return v.IVal != 0
	}
	if i, ok := v.ToInt(); ok {
		return i != 0
	}
	if f, ok := v.ToFloat(); ok {
		return f != 0
	}
	return false
}

func (m *Module) navMake(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	h, err := m.requireHeap(rt)
	if err != nil {
		return value.Nil, err
	}
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("NAV.MAKE expects 0 arguments")
	}
	gw, gh := 64, 64
	n := &navObj{
		cell:    1,
		ox:      0,
		oz:      0,
		gw:      gw,
		gh:      gh,
		blocked: make([]bool, gw*gh),
		groundY: make([]float32, gw*gh),
		built:   false,
	}
	for i := range n.blocked {
		n.blocked[i] = false
	}
	id, err := h.Alloc(n)
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}

func (m *Module) navFree(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	h, err := m.requireHeap(rt)
	if err != nil {
		return value.Nil, err
	}
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("NAV.FREE expects nav handle")
	}
	_ = h.Free(heap.Handle(args[0].IVal))
	return value.Nil, nil
}

func (m *Module) navSetGrid(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	h, err := m.requireHeap(rt)
	if err != nil {
		return value.Nil, err
	}
	if len(args) != 6 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("NAV.SETGRID expects (nav, gw, gh, cellSize, originX, originZ)")
	}
	n, err := heap.Cast[*navObj](h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	gw64, ok1 := args[1].ToInt()
	gh64, ok2 := args[2].ToInt()
	cs, ok3 := argF64(args[3])
	ox, ok4 := argF64(args[4])
	oz, ok5 := argF64(args[5])
	if !ok1 || !ok2 || !ok3 || !ok4 || !ok5 {
		return value.Nil, fmt.Errorf("NAV.SETGRID: invalid numeric arguments")
	}
	gw, gh := int(gw64), int(gh64)
	if gw < 1 || gh < 1 || gw > 4096 || gh > 4096 {
		return value.Nil, fmt.Errorf("NAV.SETGRID: gw/gh out of range")
	}
	if cs <= 0 {
		return value.Nil, fmt.Errorf("NAV.SETGRID: cell size must be > 0")
	}
	n.gw, n.gh = gw, gh
	n.cell = cs
	n.ox, n.oz = ox, oz
	n.blocked = make([]bool, gw*gh)
	n.groundY = make([]float32, gw*gh)
	n.built = false
	return value.Nil, nil
}

func (m *Module) navAddTerrain(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	h, err := m.requireHeap(rt)
	if err != nil {
		return value.Nil, err
	}
	if len(args) != 2 || args[0].Kind != value.KindHandle || args[1].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("NAV.ADDTERRAIN expects (nav, modelHandle)")
	}
	n, err := heap.Cast[*navObj](h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	bb, err := mbmodel3d.ModelBoundingBoxForNav(h, heap.Handle(args[1].IVal))
	if err != nil {
		return value.Nil, fmt.Errorf("NAV.ADDTERRAIN: %w", err)
	}
	n.setOpenRect(float64(bb.Min.X), float64(bb.Min.Z), float64(bb.Max.X), float64(bb.Max.Z))
	n.setGroundYRect(float64(bb.Min.X), float64(bb.Min.Z), float64(bb.Max.X), float64(bb.Max.Z), bb.Min.Y)
	return value.Nil, nil
}

func (m *Module) navAddObstacle(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	h, err := m.requireHeap(rt)
	if err != nil {
		return value.Nil, err
	}
	if len(args) != 2 || args[0].Kind != value.KindHandle || args[1].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("NAV.ADDOBSTACLE expects (nav, modelHandle)")
	}
	n, err := heap.Cast[*navObj](h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	bb, err := mbmodel3d.ModelBoundingBoxForNav(h, heap.Handle(args[1].IVal))
	if err != nil {
		return value.Nil, fmt.Errorf("NAV.ADDOBSTACLE: %w", err)
	}
	n.setBlockedRect(
		clampInt(int(math.Floor((float64(bb.Min.X)-n.ox)/n.cell)), 0, n.gw-1),
		clampInt(int(math.Floor((float64(bb.Min.Z)-n.oz)/n.cell)), 0, n.gh-1),
		clampInt(int(math.Floor((float64(bb.Max.X)-n.ox)/n.cell)), 0, n.gw-1),
		clampInt(int(math.Floor((float64(bb.Max.Z)-n.oz)/n.cell)), 0, n.gh-1),
		true,
	)
	return value.Nil, nil
}

func (m *Module) navBuild(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	h, err := m.requireHeap(rt)
	if err != nil {
		return value.Nil, err
	}
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("NAV.BUILD expects nav handle")
	}
	n, err := heap.Cast[*navObj](h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	n.built = true
	return value.Nil, nil
}

func (m *Module) navFindPath(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	h, err := m.requireHeap(rt)
	if err != nil {
		return value.Nil, err
	}
	if len(args) != 7 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("NAV.FINDPATH expects (nav, sx, sy, sz, tx, ty, tz)")
	}
	n, err := heap.Cast[*navObj](h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	var sx, sy, sz, tx, ty, tz float64
	var ok bool
	sx, ok = argF64(args[1])
	if !ok {
		return value.Nil, fmt.Errorf("NAV.FINDPATH: sx invalid")
	}
	sy, ok = argF64(args[2])
	if !ok {
		return value.Nil, fmt.Errorf("NAV.FINDPATH: sy invalid")
	}
	sz, ok = argF64(args[3])
	if !ok {
		return value.Nil, fmt.Errorf("NAV.FINDPATH: sz invalid")
	}
	tx, ok = argF64(args[4])
	if !ok {
		return value.Nil, fmt.Errorf("NAV.FINDPATH: tx invalid")
	}
	ty, ok = argF64(args[5])
	if !ok {
		return value.Nil, fmt.Errorf("NAV.FINDPATH: ty invalid")
	}
	tz, ok = argF64(args[6])
	if !ok {
		return value.Nil, fmt.Errorf("NAV.FINDPATH: tz invalid")
	}
	p := findPathNav(n, sx, sy, sz, tx, ty, tz)
	id, err := h.Alloc(p)
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}

func (m *Module) pathIsValid(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	h, err := m.requireHeap(rt)
	if err != nil {
		return value.Nil, err
	}
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("PATH.ISVALID expects path handle")
	}
	p, err := heap.Cast[*pathObj](h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	return value.FromBool(p.valid), nil
}

func (m *Module) pathNodeCount(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	h, err := m.requireHeap(rt)
	if err != nil {
		return value.Nil, err
	}
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("PATH.NODECOUNT expects path handle")
	}
	p, err := heap.Cast[*pathObj](h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	return value.FromInt(int64(len(p.pts))), nil
}

func (m *Module) pathNodeCoord(rt *runtime.Runtime, args []value.Value, axis int) (value.Value, error) {
	h, err := m.requireHeap(rt)
	if err != nil {
		return value.Nil, err
	}
	if len(args) != 2 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("PATH.NODEX/Y/Z expects (path, index)")
	}
	p, err := heap.Cast[*pathObj](h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	ix, ok := args[1].ToInt()
	if !ok {
		if f, okf := args[1].ToFloat(); okf {
			ix = int64(f)
			ok = true
		}
	}
	if !ok || ix < 0 || int(ix) >= len(p.pts) {
		return value.Nil, fmt.Errorf("PATH: invalid node index")
	}
	pt := p.pts[ix]
	switch axis {
	case 0:
		return value.FromFloat(pt.x), nil
	case 1:
		return value.FromFloat(pt.y), nil
	default:
		return value.FromFloat(pt.z), nil
	}
}

func (m *Module) pathNodeX(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	return m.pathNodeCoord(rt, args, 0)
}
func (m *Module) pathNodeY(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	return m.pathNodeCoord(rt, args, 1)
}
func (m *Module) pathNodeZ(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	return m.pathNodeCoord(rt, args, 2)
}

func (m *Module) pathFree(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	h, err := m.requireHeap(rt)
	if err != nil {
		return value.Nil, err
	}
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("PATH.FREE expects path handle")
	}
	_ = h.Free(heap.Handle(args[0].IVal))
	return value.Nil, nil
}

func (m *Module) agentMake(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	h, err := m.requireHeap(rt)
	if err != nil {
		return value.Nil, err
	}
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("NAVAGENT.MAKE expects nav handle")
	}
	if _, err := heap.Cast[*navObj](h, heap.Handle(args[0].IVal)); err != nil {
		return value.Nil, err
	}
	a := &navAgentObj{
		navH:      heap.Handle(args[0].IVal),
		speed:     5,
		maxForce:  40,
		arriveEps: 0.2,
	}
	id, err := h.Alloc(a)
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}

func (m *Module) agentFree(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	h, err := m.requireHeap(rt)
	if err != nil {
		return value.Nil, err
	}
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("NAVAGENT.FREE expects agent handle")
	}
	_ = h.Free(heap.Handle(args[0].IVal))
	return value.Nil, nil
}

func (m *Module) getAgent(h *heap.Store, v value.Value, op string) (*navAgentObj, error) {
	if v.Kind != value.KindHandle {
		return nil, fmt.Errorf("%s: expected agent handle", op)
	}
	return heap.Cast[*navAgentObj](h, heap.Handle(v.IVal))
}

func (m *Module) getNav(h *heap.Store, v value.Value, op string) (*navObj, error) {
	if v.Kind != value.KindHandle {
		return nil, fmt.Errorf("%s: expected nav handle", op)
	}
	return heap.Cast[*navObj](h, heap.Handle(v.IVal))
}

func (m *Module) agentSetPos(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	h, err := m.requireHeap(rt)
	if err != nil {
		return value.Nil, err
	}
	if len(args) != 4 {
		return value.Nil, fmt.Errorf("NAVAGENT.SETPOS expects (agent, x, y, z)")
	}
	a, err := m.getAgent(h, args[0], "NAVAGENT.SETPOS")
	if err != nil {
		return value.Nil, err
	}
	x, ok1 := argF64(args[1])
	y, ok2 := argF64(args[2])
	z, ok3 := argF64(args[3])
	if !ok1 || !ok2 || !ok3 {
		return value.Nil, fmt.Errorf("NAVAGENT.SETPOS: x,y,z must be numeric")
	}
	a.x, a.y, a.z = x, y, z
	return value.Nil, nil
}

func (m *Module) agentSetSpeed(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	h, err := m.requireHeap(rt)
	if err != nil {
		return value.Nil, err
	}
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("NAVAGENT.SETSPEED expects (agent, speed)")
	}
	a, err := m.getAgent(h, args[0], "NAVAGENT.SETSPEED")
	if err != nil {
		return value.Nil, err
	}
	s, ok := argF64(args[1])
	if !ok || s < 0 {
		return value.Nil, fmt.Errorf("NAVAGENT.SETSPEED: speed must be >= 0")
	}
	a.speed = s
	return value.Nil, nil
}

func (m *Module) agentSetMaxForce(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	h, err := m.requireHeap(rt)
	if err != nil {
		return value.Nil, err
	}
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("NAVAGENT.SETMAXFORCE expects (agent, maxForce)")
	}
	a, err := m.getAgent(h, args[0], "NAVAGENT.SETMAXFORCE")
	if err != nil {
		return value.Nil, err
	}
	f, ok := argF64(args[1])
	if !ok || f < 0 {
		return value.Nil, fmt.Errorf("NAVAGENT.SETMAXFORCE: invalid")
	}
	a.maxForce = f
	return value.Nil, nil
}

func (m *Module) agentApplyForce(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	h, err := m.requireHeap(rt)
	if err != nil {
		return value.Nil, err
	}
	if len(args) != 4 {
		return value.Nil, fmt.Errorf("NAVAGENT.APPLYFORCE expects (agent, fx, fy, fz)")
	}
	a, err := m.getAgent(h, args[0], "NAVAGENT.APPLYFORCE")
	if err != nil {
		return value.Nil, err
	}
	fx, ok1 := argF64(args[1])
	fy, ok2 := argF64(args[2])
	fz, ok3 := argF64(args[3])
	if !ok1 || !ok2 || !ok3 {
		return value.Nil, fmt.Errorf("NAVAGENT.APPLYFORCE: forces must be numeric")
	}
	fmag := math.Sqrt(fx*fx + fy*fy + fz*fz)
	if fmag > a.maxForce && fmag > 1e-8 {
		s := a.maxForce / fmag
		fx *= s
		fy *= s
		fz *= s
	}
	a.vx += fx
	a.vy += fy
	a.vz += fz
	vmag := math.Sqrt(a.vx*a.vx + a.vy*a.vy + a.vz*a.vz)
	if a.speed > 0 && vmag > a.speed && vmag > 1e-8 {
		s := a.speed / vmag
		a.vx *= s
		a.vy *= s
		a.vz *= s
	}
	return value.Nil, nil
}

func (m *Module) agentMoveTo(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	h, err := m.requireHeap(rt)
	if err != nil {
		return value.Nil, err
	}
	if len(args) != 4 {
		return value.Nil, fmt.Errorf("NAVAGENT.MOVETO expects (agent, tx, ty, tz)")
	}
	a, err := m.getAgent(h, args[0], "NAVAGENT.MOVETO")
	if err != nil {
		return value.Nil, err
	}
	tx, ok1 := argF64(args[1])
	ty, ok2 := argF64(args[2])
	tz, ok3 := argF64(args[3])
	if !ok1 || !ok2 || !ok3 {
		return value.Nil, fmt.Errorf("NAVAGENT.MOVETO: target must be numeric")
	}
	nav, err := heap.Cast[*navObj](h, a.navH)
	if err != nil {
		return value.Nil, err
	}
	p := findPathNav(nav, a.x, a.y, a.z, tx, ty, tz)
	if !p.valid || len(p.pts) == 0 {
		a.way = nil
		a.wayIdx = 0
		a.hasDest = false
		return value.Nil, nil
	}
	a.way = append([]pathPt(nil), p.pts...)
	a.wayIdx = 0
	a.destX, a.destY, a.destZ = tx, ty, tz
	a.hasDest = true
	a.vx, a.vy, a.vz = 0, 0, 0
	return value.Nil, nil
}

func (m *Module) agentUpdate(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	h, err := m.requireHeap(rt)
	if err != nil {
		return value.Nil, err
	}
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("NAVAGENT.UPDATE expects (agent, dt)")
	}
	a, err := m.getAgent(h, args[0], "NAVAGENT.UPDATE")
	if err != nil {
		return value.Nil, err
	}
	dt, ok := argF64(args[1])
	if !ok || dt < 0 {
		return value.Nil, fmt.Errorf("NAVAGENT.UPDATE: dt must be numeric and >= 0")
	}
	if len(a.way) > 0 && a.wayIdx < len(a.way) {
		t := a.way[a.wayIdx]
		dx := t.x - a.x
		dy := t.y - a.y
		dz := t.z - a.z
		dist := math.Sqrt(dx*dx + dy*dy + dz*dz)
		if dist < a.arriveEps {
			a.wayIdx++
			if a.wayIdx >= len(a.way) {
				a.hasDest = false
				a.way = nil
				a.wayIdx = 0
			}
			return value.Nil, nil
		}
		step := a.speed * dt
		if step > dist {
			step = dist
		}
		if dist > 1e-10 {
			inv := 1 / dist
			a.x += dx * inv * step
			a.y += dy * inv * step
			a.z += dz * inv * step
		}
		return value.Nil, nil
	}
	a.x += a.vx * dt
	a.y += a.vy * dt
	a.z += a.vz * dt
	damp := math.Exp(-2.0 * dt)
	a.vx *= damp
	a.vy *= damp
	a.vz *= damp
	return value.Nil, nil
}

func (m *Module) agentIsAtDestination(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	h, err := m.requireHeap(rt)
	if err != nil {
		return value.Nil, err
	}
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("NAVAGENT.ISATDESTINATION expects agent handle")
	}
	a, err := m.getAgent(h, args[0], "NAVAGENT.ISATDESTINATION")
	if err != nil {
		return value.Nil, err
	}
	if a.hasDest {
		return value.FromBool(false), nil
	}
	return value.FromBool(true), nil
}

func (m *Module) agentX(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	h, err := m.requireHeap(rt)
	if err != nil {
		return value.Nil, err
	}
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("NAVAGENT.X expects agent handle")
	}
	a, err := m.getAgent(h, args[0], "NAVAGENT.X")
	if err != nil {
		return value.Nil, err
	}
	return value.FromFloat(a.x), nil
}

func (m *Module) agentY(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	h, err := m.requireHeap(rt)
	if err != nil {
		return value.Nil, err
	}
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("NAVAGENT.Y expects agent handle")
	}
	a, err := m.getAgent(h, args[0], "NAVAGENT.Y")
	if err != nil {
		return value.Nil, err
	}
	return value.FromFloat(a.y), nil
}

func (m *Module) agentZ(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	h, err := m.requireHeap(rt)
	if err != nil {
		return value.Nil, err
	}
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("NAVAGENT.Z expects agent handle")
	}
	a, err := m.getAgent(h, args[0], "NAVAGENT.Z")
	if err != nil {
		return value.Nil, err
	}
	return value.FromFloat(a.z), nil
}

func (m *Module) steerGroupMake(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	h, err := m.requireHeap(rt)
	if err != nil {
		return value.Nil, err
	}
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("STEER.GROUPMAKE expects 0 arguments")
	}
	g := &steerGroupObj{agents: nil}
	id, err := h.Alloc(g)
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}

func (m *Module) steerGroupAdd(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	h, err := m.requireHeap(rt)
	if err != nil {
		return value.Nil, err
	}
	if len(args) != 2 || args[0].Kind != value.KindHandle || args[1].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("STEER.GROUPADD expects (group, navAgent)")
	}
	g, err := heap.Cast[*steerGroupObj](h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	if _, err := m.getAgent(h, args[1], "STEER.GROUPADD"); err != nil {
		return value.Nil, err
	}
	ah := heap.Handle(args[1].IVal)
	for _, x := range g.agents {
		if x == ah {
			return value.Nil, nil
		}
	}
	g.agents = append(g.agents, ah)
	return value.Nil, nil
}

func (m *Module) steerGroupClear(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	h, err := m.requireHeap(rt)
	if err != nil {
		return value.Nil, err
	}
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("STEER.GROUPCLEAR expects group handle")
	}
	g, err := heap.Cast[*steerGroupObj](h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	g.agents = nil
	return value.Nil, nil
}

func agentPos(h *heap.Store, ah heap.Handle) (x, y, z float64, err error) {
	a, err := heap.Cast[*navAgentObj](h, ah)
	if err != nil {
		return 0, 0, 0, err
	}
	return a.x, a.y, a.z, nil
}

func (m *Module) allocSteer(h *heap.Store, fx, fy, fz float64) (value.Value, error) {
	return mbmatrix.AllocVec3Value(h, float32(fx), float32(fy), float32(fz))
}

func (m *Module) steerSeek(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	h, err := m.requireHeap(rt)
	if err != nil {
		return value.Nil, err
	}
	if len(args) != 4 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("STEER.SEEK expects (agent, tx, ty, tz)")
	}
	ax, ay, az, err := agentPos(h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	tx, ok1 := argF64(args[1])
	ty, ok2 := argF64(args[2])
	tz, ok3 := argF64(args[3])
	if !ok1 || !ok2 || !ok3 {
		return value.Nil, fmt.Errorf("STEER.SEEK: target invalid")
	}
	dx, dy, dz := tx-ax, ty-ay, tz-az
	return m.allocSteer(h, dx, dy, dz)
}

func (m *Module) steerFlee(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	h, err := m.requireHeap(rt)
	if err != nil {
		return value.Nil, err
	}
	if len(args) != 4 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("STEER.FLEE expects (agent, tx, ty, tz)")
	}
	ax, ay, az, err := agentPos(h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	tx, ok1 := argF64(args[1])
	ty, ok2 := argF64(args[2])
	tz, ok3 := argF64(args[3])
	if !ok1 || !ok2 || !ok3 {
		return value.Nil, fmt.Errorf("STEER.FLEE: target invalid")
	}
	dx, dy, dz := ax-tx, ay-ty, az-tz
	return m.allocSteer(h, dx, dy, dz)
}

func (m *Module) steerArrive(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	h, err := m.requireHeap(rt)
	if err != nil {
		return value.Nil, err
	}
	if len(args) != 5 {
		return value.Nil, fmt.Errorf("STEER.ARRIVE expects (agent, tx, ty, tz, slowingRadius)")
	}
	a, err := heap.Cast[*navAgentObj](h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	tx, ok1 := argF64(args[1])
	ty, ok2 := argF64(args[2])
	tz, ok3 := argF64(args[3])
	rad, ok4 := argF64(args[4])
	if !ok1 || !ok2 || !ok3 || !ok4 || rad <= 0 {
		return value.Nil, fmt.Errorf("STEER.ARRIVE: invalid arguments")
	}
	dx, dy, dz := tx-a.x, ty-a.y, tz-a.z
	dist := math.Sqrt(dx*dx + dy*dy + dz*dz)
	if dist < 1e-6 {
		return m.allocSteer(h, 0, 0, 0)
	}
	speed := a.speed
	if dist < rad {
		speed *= dist / rad
	}
	inv := 1 / dist
	return m.allocSteer(h, dx*inv*speed, dy*inv*speed, dz*inv*speed)
}

func (m *Module) steerWander(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	h, err := m.requireHeap(rt)
	if err != nil {
		return value.Nil, err
	}
	if len(args) != 3 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("STEER.WANDER expects (agent, speed, jitterRadius)")
	}
	if _, err := heap.Cast[*navAgentObj](h, heap.Handle(args[0].IVal)); err != nil {
		return value.Nil, err
	}
	sp, ok1 := argF64(args[1])
	rad, ok2 := argF64(args[2])
	if !ok1 || !ok2 {
		return value.Nil, fmt.Errorf("STEER.WANDER: invalid numbers")
	}
	// deterministic-ish hash from agent id + time would need frame; use simple pseudo-random from handle
	seed := uint64(args[0].IVal)*0x9E3779B97F4A7C15 + 1
	ax := float64((seed>>0)&0xFFFF)/65535.0*2 - 1
	ay := float64((seed>>16)&0xFFFF)/65535.0*2 - 1
	az := float64((seed>>32)&0xFFFF)/65535.0*2 - 1
	mag := math.Sqrt(ax*ax + ay*ay + az*az)
	if mag < 1e-6 {
		return m.allocSteer(h, 0, sp, 0)
	}
	inv := rad / mag
	return m.allocSteer(h, ax*inv*sp, ay*inv*sp, az*inv*sp)
}

func (m *Module) steerFlock(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	h, err := m.requireHeap(rt)
	if err != nil {
		return value.Nil, err
	}
	if len(args) != 5 || args[0].Kind != value.KindHandle || args[1].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("STEER.FLOCK expects (selfAgent, group, cohesion, separation, alignment)")
	}
	self, err := heap.Cast[*navAgentObj](h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	g, err := heap.Cast[*steerGroupObj](h, heap.Handle(args[1].IVal))
	if err != nil {
		return value.Nil, err
	}
	coh, ok1 := argF64(args[2])
	sep, ok2 := argF64(args[3])
	aln, ok3 := argF64(args[4])
	if !ok1 || !ok2 || !ok3 {
		return value.Nil, fmt.Errorf("STEER.FLOCK: weights invalid")
	}
	var cx, cy, cz float64
	var sx, sy, sz float64
	var vx, vy, vz float64
	var nv int
	selfH := heap.Handle(args[0].IVal)
	for _, ah := range g.agents {
		if ah == selfH {
			continue
		}
		ox, oy, oz, e := agentPos(h, ah)
		if e != nil {
			continue
		}
		nv++
		cx += ox
		cy += oy
		cz += oz
		dx := self.x - ox
		dy := self.y - oy
		dz := self.z - oz
		d2 := dx*dx + dy*dy + dz*dz
		if d2 < 1e-8 {
			continue
		}
		inv := 1 / d2
		sx += dx * inv
		sy += dy * inv
		sz += dz * inv
		if oa, e2 := heap.Cast[*navAgentObj](h, ah); e2 == nil {
			vx += oa.vx
			vy += oa.vy
			vz += oa.vz
		}
	}
	if nv == 0 {
		return m.allocSteer(h, 0, 0, 0)
	}
	cx /= float64(nv)
	cy /= float64(nv)
	cz /= float64(nv)
	fx := (cx-self.x)*coh + sx*sep + (vx/float64(nv)-self.vx)*aln
	fy := (cy-self.y)*coh + sy*sep + (vy/float64(nv)-self.vy)*aln
	fz := (cz-self.z)*coh + sz*sep + (vz/float64(nv)-self.vz)*aln
	return m.allocSteer(h, fx, fy, fz)
}

func (m *Module) steerAvoidObstacles(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	h, err := m.requireHeap(rt)
	if err != nil {
		return value.Nil, err
	}
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("STEER.AVOIDOBSTACLES expects (agent, radius)")
	}
	a, err := m.getAgent(h, args[0], "STEER.AVOIDOBSTACLES")
	if err != nil {
		return value.Nil, err
	}
	rad, ok := argF64(args[1])
	if !ok || rad <= 0 {
		return value.Nil, fmt.Errorf("STEER.AVOIDOBSTACLES: radius must be > 0")
	}
	nav, err := heap.Cast[*navObj](h, a.navH)
	if err != nil {
		return value.Nil, err
	}
	ix, iz, okc := nav.worldToCell(a.x, a.z)
	if !okc {
		return m.allocSteer(h, 0, 0, 0)
	}
	ri := int(math.Ceil(rad / nav.cell))
	if ri < 1 {
		ri = 1
	}
	var fx, fz float64
	for dz := -ri; dz <= ri; dz++ {
		for dx := -ri; dx <= ri; dx++ {
			nix, niz := ix+dx, iz+dz
			if !nav.containsCell(nix, niz) {
				continue
			}
			if !nav.blocked[nav.idx(nix, niz)] {
				continue
			}
			cx, cz := nav.cellCenter(nix, niz)
			ddx := a.x - cx
			ddz := a.z - cz
			d2 := ddx*ddx + ddz*ddz
			if d2 < 1e-8 {
				fx += 1
				continue
			}
			inv := 1 / d2
			fx += ddx * inv
			fz += ddz * inv
		}
	}
	return m.allocSteer(h, fx, 0, fz)
}

func (m *Module) steerFollowPath(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	h, err := m.requireHeap(rt)
	if err != nil {
		return value.Nil, err
	}
	if len(args) != 2 || args[1].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("STEER.FOLLOWPATH expects (agent, path)")
	}
	a, err := m.getAgent(h, args[0], "STEER.FOLLOWPATH")
	if err != nil {
		return value.Nil, err
	}
	p, err := heap.Cast[*pathObj](h, heap.Handle(args[1].IVal))
	if err != nil {
		return value.Nil, err
	}
	if !p.valid || len(p.pts) == 0 {
		return m.allocSteer(h, 0, 0, 0)
	}
	// seek first waypoint ahead of agent (closest forward)
	best := 0
	bestD := 1e30
	for i := range p.pts {
		pt := p.pts[i]
		dx := pt.x - a.x
		dy := pt.y - a.y
		dz := pt.z - a.z
		d := dx*dx + dy*dy + dz*dz
		if d < bestD {
			bestD = d
			best = i
		}
	}
	pt := p.pts[best]
	dx, dy, dz := pt.x-a.x, pt.y-a.y, pt.z-a.z
	return m.allocSteer(h, dx, dy, dz)
}

func (m *Module) btMake(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	h, err := m.requireHeap(rt)
	if err != nil {
		return value.Nil, err
	}
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("BTREE.MAKE expects 0 arguments")
	}
	b := &btObj{root: &btNode{kind: btKindSeq}}
	id, err := h.Alloc(b)
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}

func (m *Module) btFree(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	h, err := m.requireHeap(rt)
	if err != nil {
		return value.Nil, err
	}
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("BTREE.FREE expects btree handle")
	}
	_ = h.Free(heap.Handle(args[0].IVal))
	return value.Nil, nil
}

func (m *Module) btSequence(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("BTREE.SEQUENCE expects btree handle")
	}
	return args[0], nil
}

func (m *Module) btAddCondition(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	h, err := m.requireHeap(rt)
	if err != nil {
		return value.Nil, err
	}
	if len(args) != 2 || args[0].Kind != value.KindHandle || args[1].Kind != value.KindString {
		return value.Nil, fmt.Errorf("BTREE.ADDCONDITION expects (btreeOrSeq, functionName$)")
	}
	b, err := heap.Cast[*btObj](h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	fn, err := rt.ArgString(args, 1)
	if err != nil {
		return value.Nil, err
	}
	if b.root == nil {
		b.root = &btNode{kind: btKindSeq}
	}
	b.root.kids = append(b.root.kids, &btNode{kind: btKindCond, fn: fn})
	return value.Nil, nil
}

func (m *Module) btAddAction(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	h, err := m.requireHeap(rt)
	if err != nil {
		return value.Nil, err
	}
	if len(args) != 2 || args[0].Kind != value.KindHandle || args[1].Kind != value.KindString {
		return value.Nil, fmt.Errorf("BTREE.ADDACTION expects (btreeOrSeq, functionName$)")
	}
	b, err := heap.Cast[*btObj](h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	fn, err := rt.ArgString(args, 1)
	if err != nil {
		return value.Nil, err
	}
	if b.root == nil {
		b.root = &btNode{kind: btKindSeq}
	}
	b.root.kids = append(b.root.kids, &btNode{kind: btKindAct, fn: fn})
	return value.Nil, nil
}

func tickBT(n *btNode, agent value.Value, inv func(string, []value.Value) (value.Value, error)) (bool, error) {
	if n == nil {
		return true, nil
	}
	switch n.kind {
	case btKindSeq:
		for _, ch := range n.kids {
			ok, err := tickBT(ch, agent, inv)
			if err != nil || !ok {
				return ok, err
			}
		}
		return true, nil
	case btKindCond, btKindAct:
		v, err := inv(n.fn, []value.Value{agent})
		if err != nil {
			return false, err
		}
		return valueTruthy(v), nil
	default:
		return true, nil
	}
}

func (m *Module) btRun(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	h, err := m.requireHeap(rt)
	if err != nil {
		return value.Nil, err
	}
	if len(args) != 3 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("BTREE.RUN expects (btree, agentHandle, dt)")
	}
	b, err := heap.Cast[*btObj](h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	if args[1].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("BTREE.RUN: agent must be a handle")
	}
	agent := args[1]
	if _, ok := argF64(args[2]); !ok {
		return value.Nil, fmt.Errorf("BTREE.RUN: dt must be numeric (reserved for future)")
	}
	_, err = tickBT(b.root, agent, func(name string, av []value.Value) (value.Value, error) {
		return m.callUser(name, av)
	})
	if err != nil {
		return value.Nil, err
	}
	return value.Nil, nil
}
