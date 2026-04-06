//go:build cgo || (windows && !cgo)

package mbentity

import (
	"fmt"
	"math"
	"sort"

	mbcamera "moonbasic/runtime/camera"
	"moonbasic/runtime"
	"moonbasic/runtime/mbgame"
	"moonbasic/runtime/mbmatrix"
	mbtime "moonbasic/runtime/time"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var entityStores = make(map[*Module]*entityStore)

type entityStore struct {
	ents    map[int64]*ent
	nextID  int64
	byName  map[string]int64
	groups  map[string]map[int64]struct{}
}

func (m *Module) store() *entityStore {
	s := entityStores[m]
	if s == nil {
		s = &entityStore{
			ents:   make(map[int64]*ent),
			nextID: 1,
			byName: make(map[string]int64),
		}
		entityStores[m] = s
	}
	return s
}

func (m *Module) Register(r runtime.Registrar) {
	r.Register("ENTITY.CREATE", "entity", runtime.AdaptLegacy(m.entCreate))
	r.Register("ENTITY.CREATEENTITY", "entity", runtime.AdaptLegacy(m.entCreate))
	r.Register("ENTITY.CREATEBOX", "entity", runtime.AdaptLegacy(m.entCreateBox))
	r.Register("ENTITY.CREATECUBE", "entity", runtime.AdaptLegacy(m.entCreateBox))
	registerEntityBlitzAPI(m, r)
	registerEntitySceneGroupAPI(m, r)
	r.Register("ENTITY.SETPOSITION", "entity", runtime.AdaptLegacy(m.entSetPosition))
	r.Register("ENTITY.GETPOSITION", "entity", runtime.AdaptLegacy(m.entGetPosition))
	r.Register("ENTITY.MOVE", "entity", runtime.AdaptLegacy(m.entMove))
	r.Register("ENTITY.TRANSLATE", "entity", runtime.AdaptLegacy(m.entTranslate))
	r.Register("ENTITY.ROTATE", "entity", runtime.AdaptLegacy(m.entRotate))
	r.Register("ENTITY.SCALE", "entity", runtime.AdaptLegacy(m.entScale))
	r.Register("ENTITY.COLOR", "entity", runtime.AdaptLegacy(m.entColor))
	r.Register("ENTITY.RADIUS", "entity", runtime.AdaptLegacy(m.entRadius))
	r.Register("ENTITY.BOX", "entity", runtime.AdaptLegacy(m.entBox))
	r.Register("ENTITY.COLLIDED", "entity", runtime.AdaptLegacy(m.entCollided))
	r.Register("ENTITY.COLLISIONOTHER", "entity", runtime.AdaptLegacy(m.entCollisionOther))
	r.Register("ENTITY.FLOOR", "entity", runtime.AdaptLegacy(m.entFloor))
	r.Register("ENTITY.SETGRAVITY", "entity", runtime.AdaptLegacy(m.entSetGravity))
	r.Register("ENTITY.JUMP", "entity", runtime.AdaptLegacy(m.entJump))
	r.Register("ENTITY.UPDATE", "entity", runtime.AdaptLegacy(m.entUpdate))
	r.Register("ENTITY.DRAWALL", "entity", runtime.AdaptLegacy(m.entDrawAll))
	r.Register("CAMERA.FOLLOWENTITY", "entity", m.camFollowEntity)
	registerBlitzEntityHandles(m, r)
}

// Shutdown implements runtime.Module.
func (m *Module) Shutdown() {
	clearEntityRefFreeHookIfOwner(m)
	delete(entityStores, m)
}

func (m *Module) camFollowEntity(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 5 {
		return value.Nil, fmt.Errorf("CAMERA.FOLLOWENTITY expects 5 arguments (camera, entity#, dist#, height#, smooth#)")
	}
	ch, ok := argHandle(args[0])
	if !ok {
		return value.Nil, fmt.Errorf("CAMERA.FOLLOWENTITY: invalid camera handle")
	}
	eid, ok := args[1].ToInt()
	if !ok || eid < 1 {
		return value.Nil, fmt.Errorf("CAMERA.FOLLOWENTITY: invalid entity id")
	}
	dist, ok1 := argF32(args[2])
	height, ok2 := argF32(args[3])
	smooth, ok3 := argF32(args[4])
	if !ok1 || !ok2 || !ok3 {
		return value.Nil, fmt.Errorf("CAMERA.FOLLOWENTITY: numeric arguments required")
	}
	e := m.store().ents[eid]
	if e == nil {
		return value.Nil, fmt.Errorf("CAMERA.FOLLOWENTITY: unknown entity %d", eid)
	}
	dt := mbtime.DeltaSeconds(rt)
	if dt <= 0 {
		dt = 1.0 / 60.0
	}
	wp := m.worldPos(e)
	err := mbcamera.ThirdPersonFollowStep(m.h, ch, wp.X, wp.Y, wp.Z, e.yaw, dist, height, smooth, dt)
	if err != nil {
		return value.Nil, err
	}
	return value.Nil, nil
}

func (m *Module) entCreate(args []value.Value) (value.Value, error) {
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("ENTITY.CREATE expects 0 arguments")
	}
	st := m.store()
	id := st.nextID
	st.nextID++
	e := newDefaultEnt(id)
	e.kind = entKindSphere
	e.w, e.h, e.d = 1, 1, 1
	e.radius = 0.5
	e.useSphere = true
	e.static = false
	e.gravity = -28
	st.ents[id] = e
	return value.FromInt(id), nil
}

func (m *Module) entCreateBox(args []value.Value) (value.Value, error) {
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("ENTITY.CREATEBOX expects 3 arguments (w#, h#, d#)")
	}
	w, ok1 := argF32(args[0])
	h, ok2 := argF32(args[1])
	d, ok3 := argF32(args[2])
	if !ok1 || !ok2 || !ok3 {
		return value.Nil, fmt.Errorf("ENTITY.CREATEBOX: dimensions must be numeric")
	}
	st := m.store()
	id := st.nextID
	st.nextID++
	e := newDefaultEnt(id)
	e.kind = entKindBox
	e.r, e.g, e.b = 180, 180, 200
	e.w, e.h, e.d = w, h, d
	e.static = true
	e.useSphere = false
	e.gravity = 0
	st.ents[id] = e
	return value.FromInt(id), nil
}

func (m *Module) entSetPosition(args []value.Value) (value.Value, error) {
	if len(args) != 4 && len(args) != 5 {
		return value.Nil, fmt.Errorf("ENTITY.SETPOSITION expects 4–5 arguments (entity#, x#, y#, z# [, global])")
	}
	id, ok := args[0].ToInt()
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("ENTITY.SETPOSITION: invalid entity")
	}
	e := m.store().ents[id]
	if e == nil {
		return value.Nil, fmt.Errorf("ENTITY.SETPOSITION: unknown entity %d", id)
	}
	x, ok1 := argF32(args[1])
	y, ok2 := argF32(args[2])
	z, ok3 := argF32(args[3])
	if !ok1 || !ok2 || !ok3 {
		return value.Nil, fmt.Errorf("ENTITY.SETPOSITION: position must be numeric")
	}
	global := false
	if len(args) == 5 {
		switch args[4].Kind {
		case value.KindBool:
			global = args[4].IVal != 0
		case value.KindInt:
			global = args[4].IVal != 0
		default:
			return value.Nil, fmt.Errorf("ENTITY.SETPOSITION: global must be TRUE/FALSE or 0/1")
		}
	}
	if global {
		m.setLocalFromWorld(e, x, y, z)
	} else {
		e.pos = rl.Vector3{X: x, Y: y, Z: z}
	}
	return value.Nil, nil
}

func (m *Module) entGetPosition(args []value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("ENTITY.GETPOSITION: heap not bound")
	}
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("ENTITY.GETPOSITION expects entity#")
	}
	id, ok := args[0].ToInt()
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("ENTITY.GETPOSITION: invalid entity")
	}
	e := m.store().ents[id]
	if e == nil {
		return value.Nil, fmt.Errorf("ENTITY.GETPOSITION: unknown entity %d", id)
	}
	p := m.worldPos(e)
	return mbmatrix.AllocVec3Value(m.h, p.X, p.Y, p.Z)
}

func localAxes(yaw, pitch float32) (forward, right, up rl.Vector3) {
	cp := float64(math.Cos(float64(pitch)))
	sp := float64(math.Sin(float64(pitch)))
	sy := float64(math.Sin(float64(yaw)))
	cy := float64(math.Cos(float64(yaw)))
	fx := float32(sy * cp)
	fy := float32(sp)
	fz := float32(cy * cp)
	forward = rl.Vector3Normalize(rl.Vector3{X: fx, Y: fy, Z: fz})
	worldUp := rl.Vector3{X: 0, Y: 1, Z: 0}
	right = rl.Vector3Normalize(rl.Vector3CrossProduct(worldUp, forward))
	if rl.Vector3Length(right) < 1e-6 {
		right = rl.Vector3{X: 1, Y: 0, Z: 0}
	}
	up = rl.Vector3Normalize(rl.Vector3CrossProduct(right, forward))
	return forward, right, up
}

func (m *Module) entMove(args []value.Value) (value.Value, error) {
	if len(args) != 4 {
		return value.Nil, fmt.Errorf("ENTITY.MOVE expects 4 arguments (entity#, forward#, right#, up#)")
	}
	id, ok := args[0].ToInt()
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("ENTITY.MOVE: invalid entity")
	}
	e := m.store().ents[id]
	if e == nil {
		return value.Nil, fmt.Errorf("ENTITY.MOVE: unknown entity %d", id)
	}
	f, ok1 := argF32(args[1])
	rg, ok2 := argF32(args[2])
	u, ok3 := argF32(args[3])
	if !ok1 || !ok2 || !ok3 {
		return value.Nil, fmt.Errorf("ENTITY.MOVE: deltas must be numeric")
	}
	fwd, right, up := localAxes(e.yaw, e.pitch)
	delta := rl.Vector3Add(rl.Vector3Add(rl.Vector3Scale(fwd, f), rl.Vector3Scale(right, rg)), rl.Vector3Scale(up, u))
	wp := m.worldPos(e)
	nw := rl.Vector3Add(wp, delta)
	m.setLocalFromWorld(e, nw.X, nw.Y, nw.Z)
	return value.Nil, nil
}

func (m *Module) entTranslate(args []value.Value) (value.Value, error) {
	if len(args) != 4 {
		return value.Nil, fmt.Errorf("ENTITY.TRANSLATE expects 4 arguments (entity#, dx#, dy#, dz#)")
	}
	id, ok := args[0].ToInt()
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("ENTITY.TRANSLATE: invalid entity")
	}
	e := m.store().ents[id]
	if e == nil {
		return value.Nil, fmt.Errorf("ENTITY.TRANSLATE: unknown entity %d", id)
	}
	dx, ok1 := argF32(args[1])
	dy, ok2 := argF32(args[2])
	dz, ok3 := argF32(args[3])
	if !ok1 || !ok2 || !ok3 {
		return value.Nil, fmt.Errorf("ENTITY.TRANSLATE: deltas must be numeric")
	}
	wp := m.worldPos(e)
	nw := rl.Vector3Add(wp, rl.Vector3{X: dx, Y: dy, Z: dz})
	m.setLocalFromWorld(e, nw.X, nw.Y, nw.Z)
	return value.Nil, nil
}

func (m *Module) entRotate(args []value.Value) (value.Value, error) {
	if len(args) != 4 {
		return value.Nil, fmt.Errorf("ENTITY.ROTATE expects 4 arguments (entity#, dpitch#, dyaw#, droll#)")
	}
	id, ok := args[0].ToInt()
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("ENTITY.ROTATE: invalid entity")
	}
	e := m.store().ents[id]
	if e == nil {
		return value.Nil, fmt.Errorf("ENTITY.ROTATE: unknown entity %d", id)
	}
	dp, ok1 := argF32(args[1])
	dy, ok2 := argF32(args[2])
	dr, ok3 := argF32(args[3])
	if !ok1 || !ok2 || !ok3 {
		return value.Nil, fmt.Errorf("ENTITY.ROTATE: angles must be numeric")
	}
	e.pitch += dp
	e.yaw += dy
	e.roll += dr
	return value.Nil, nil
}

func (m *Module) entScale(args []value.Value) (value.Value, error) {
	if len(args) != 4 {
		return value.Nil, fmt.Errorf("ENTITY.SCALE expects 4 arguments (entity#, sx#, sy#, sz#)")
	}
	id, ok := args[0].ToInt()
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("ENTITY.SCALE: invalid entity")
	}
	e := m.store().ents[id]
	if e == nil {
		return value.Nil, fmt.Errorf("ENTITY.SCALE: unknown entity %d", id)
	}
	sx, ok1 := argF32(args[1])
	sy, ok2 := argF32(args[2])
	sz, ok3 := argF32(args[3])
	if !ok1 || !ok2 || !ok3 {
		return value.Nil, fmt.Errorf("ENTITY.SCALE: scale must be numeric")
	}
	e.scale = rl.Vector3{X: sx, Y: sy, Z: sz}
	return value.Nil, nil
}

func (m *Module) entColor(args []value.Value) (value.Value, error) {
	if len(args) != 4 {
		return value.Nil, fmt.Errorf("ENTITY.COLOR expects 4 arguments (entity#, r, g, b)")
	}
	id, ok := args[0].ToInt()
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("ENTITY.COLOR: invalid entity")
	}
	e := m.store().ents[id]
	if e == nil {
		return value.Nil, fmt.Errorf("ENTITY.COLOR: unknown entity %d", id)
	}
	ri, ok1 := args[1].ToInt()
	gi, ok2 := args[2].ToInt()
	bi, ok3 := args[3].ToInt()
	if !ok1 || !ok2 || !ok3 {
		return value.Nil, fmt.Errorf("ENTITY.COLOR: RGB must be integer")
	}
	e.r = uint8(ri)
	e.g = uint8(gi)
	e.b = uint8(bi)
	return value.Nil, nil
}

func (m *Module) entRadius(args []value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("ENTITY.RADIUS expects 2 arguments (entity#, radius#)")
	}
	id, ok := args[0].ToInt()
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("ENTITY.RADIUS: invalid entity")
	}
	e := m.store().ents[id]
	if e == nil {
		return value.Nil, fmt.Errorf("ENTITY.RADIUS: unknown entity %d", id)
	}
	rad, ok1 := argF32(args[1])
	if !ok1 {
		return value.Nil, fmt.Errorf("ENTITY.RADIUS: radius must be numeric")
	}
	e.radius = rad
	e.useSphere = true
	e.static = false
	if e.gravity == 0 {
		e.gravity = -28
	}
	return value.Nil, nil
}

func (m *Module) entBox(args []value.Value) (value.Value, error) {
	if len(args) != 4 {
		return value.Nil, fmt.Errorf("ENTITY.BOX expects 4 arguments (entity#, w#, h#, d#)")
	}
	id, ok := args[0].ToInt()
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("ENTITY.BOX: invalid entity")
	}
	e := m.store().ents[id]
	if e == nil {
		return value.Nil, fmt.Errorf("ENTITY.BOX: unknown entity %d", id)
	}
	w, ok1 := argF32(args[1])
	h, ok2 := argF32(args[2])
	d, ok3 := argF32(args[3])
	if !ok1 || !ok2 || !ok3 {
		return value.Nil, fmt.Errorf("ENTITY.BOX: dimensions must be numeric")
	}
	e.w, e.h, e.d = w, h, d
	e.useSphere = false
	return value.Nil, nil
}

func (m *Module) entCollided(args []value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("ENTITY.COLLIDED expects entity#")
	}
	id, ok := args[0].ToInt()
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("ENTITY.COLLIDED: invalid entity")
	}
	e := m.store().ents[id]
	if e == nil {
		return value.Nil, fmt.Errorf("ENTITY.COLLIDED: unknown entity %d", id)
	}
	return value.FromBool(e.collided), nil
}

func (m *Module) entCollisionOther(args []value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("ENTITY.COLLISIONOTHER expects entity#")
	}
	id, ok := args[0].ToInt()
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("ENTITY.COLLISIONOTHER: invalid entity")
	}
	e := m.store().ents[id]
	if e == nil {
		return value.Nil, fmt.Errorf("ENTITY.COLLISIONOTHER: unknown entity %d", id)
	}
	return value.FromInt(e.otherID), nil
}

func (m *Module) entFloor(args []value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("ENTITY.FLOOR expects entity#")
	}
	id, ok := args[0].ToInt()
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("ENTITY.FLOOR: invalid entity")
	}
	e := m.store().ents[id]
	if e == nil {
		return value.Nil, fmt.Errorf("ENTITY.FLOOR: unknown entity %d", id)
	}
	y := m.queryFloorY(e)
	return value.FromFloat(y), nil
}

func (m *Module) queryFloorY(e *ent) float64 {
	var pr float64
	if e.useSphere {
		pr = float64(e.radius)
	} else {
		pr = float64(e.h) * 0.5
	}
	wp := m.worldPos(e)
	px, py, pz := float64(wp.X), float64(wp.Y), float64(wp.Z)
	var best float64
	found := false
	for _, s := range m.store().ents {
		if !s.static {
			continue
		}
		bx, by, bz := float64(s.pos.X), float64(s.pos.Y), float64(s.pos.Z)
		bw, bh, bd := float64(s.w), float64(s.h), float64(s.d)
		top := by + bh*0.5
		halfW := bw*0.5 + pr
		halfD := bd*0.5 + pr
		if math.Abs(px-bx) > halfW || math.Abs(pz-bz) > halfD {
			continue
		}
		feet := py - pr
		if feet >= top-0.25 && feet <= top+3.0 {
			if !found || top > best {
				best = top
				found = true
			}
		}
	}
	if !found {
		return 0
	}
	return best
}

func (m *Module) entSetGravity(args []value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("ENTITY.SETGRAVITY expects 2 arguments (entity#, gravity#)")
	}
	id, ok := args[0].ToInt()
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("ENTITY.SETGRAVITY: invalid entity")
	}
	e := m.store().ents[id]
	if e == nil {
		return value.Nil, fmt.Errorf("ENTITY.SETGRAVITY: unknown entity %d", id)
	}
	g, ok1 := argF32(args[1])
	if !ok1 {
		return value.Nil, fmt.Errorf("ENTITY.SETGRAVITY: gravity must be numeric")
	}
	e.gravity = g
	e.static = false
	return value.Nil, nil
}

func (m *Module) entJump(args []value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("ENTITY.JUMP expects 2 arguments (entity#, force#)")
	}
	id, ok := args[0].ToInt()
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("ENTITY.JUMP: invalid entity")
	}
	e := m.store().ents[id]
	if e == nil {
		return value.Nil, fmt.Errorf("ENTITY.JUMP: unknown entity %d", id)
	}
	f, ok1 := argF32(args[1])
	if !ok1 {
		return value.Nil, fmt.Errorf("ENTITY.JUMP: force must be numeric")
	}
	if e.onGround {
		e.vel.Y += f
		e.onGround = false
	}
	return value.Nil, nil
}

func (m *Module) entUpdate(args []value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("ENTITY.UPDATE expects 1 argument (dt#)")
	}
	dt, ok := argF32(args[0])
	if !ok || dt <= 0 {
		return value.Nil, fmt.Errorf("ENTITY.UPDATE: dt must be positive")
	}
	for _, e := range m.store().ents {
		e.collided = false
		e.otherID = 0
		e.hasHit = false
	}
	for _, e := range m.store().ents {
		if e.static {
			continue
		}
		e.vel.Y += e.gravity * dt
		wp := m.worldPos(e)
		nw := rl.Vector3Add(wp, rl.Vector3Scale(e.vel, dt))
		m.setLocalFromWorld(e, nw.X, nw.Y, nw.Z)

		if e.useSphere {
			m.resolveSphereVsStatics(e)
		} else {
			m.resolveBoxVsStatics(e)
		}
		e.onGround = false
		if e.useSphere {
			wp2 := m.worldPos(e)
			px, py, pz := float64(wp2.X), float64(wp2.Y), float64(wp2.Z)
			pvy := float64(e.vel.Y)
			pr := float64(e.radius)
			var bestSnap float64
			found := false
			for _, s := range m.store().ents {
				if !s.static {
					continue
				}
				bx, by, bz := float64(s.pos.X), float64(s.pos.Y), float64(s.pos.Z)
				bw, bh, bd := float64(s.w), float64(s.h), float64(s.d)
				snap := mbgame.BoxTopLandSnap(px, py, pz, pvy, pr, bx, by, bz, bw, bh, bd)
				if snap != 0 && (!found || snap > bestSnap) {
					bestSnap = snap
					found = true
				}
			}
			if found {
				wp := m.worldPos(e)
				wp.Y = float32(bestSnap)
				m.setLocalFromWorld(e, wp.X, wp.Y, wp.Z)
				if e.vel.Y < 0 {
					e.vel.Y = 0
				}
				e.onGround = true
			}
		}
	}
	m.pairwiseDynamic()
	for _, e := range m.store().ents {
		if len(e.modelAnims) == 0 {
			continue
		}
		ai := e.animIndex
		if ai < 0 || int(ai) >= len(e.modelAnims) {
			ai = 0
		}
		anim := e.modelAnims[ai]
		if anim.FrameCount <= 0 {
			continue
		}
		if e.animSpeed != 0 {
			e.animTime += dt * e.animSpeed * 30
		}
		var frame int32
		if e.animMode == 0 {
			frame = int32(e.animTime) % anim.FrameCount
			if frame < 0 {
				frame += anim.FrameCount
			}
		} else {
			frame = int32(e.animTime)
			if frame >= anim.FrameCount {
				frame = anim.FrameCount - 1
			}
		}
		rl.UpdateModelAnimation(e.rlModel, anim, frame)
	}
	return value.Nil, nil
}

func (m *Module) aabbWorldMinMax(e *ent) (mn, mx rl.Vector3) {
	c := m.worldPos(e)
	hx, hy, hz := e.w*e.scale.X*0.5, e.h*e.scale.Y*0.5, e.d*e.scale.Z*0.5
	mn = rl.Vector3{X: c.X - hx, Y: c.Y - hy, Z: c.Z - hz}
	mx = rl.Vector3{X: c.X + hx, Y: c.Y + hy, Z: c.Z + hz}
	return mn, mx
}

func (m *Module) resolveSphereVsStatics(e *ent) {
	r := e.radius
	if r <= 0 {
		return
	}
	wp := m.worldPos(e)
	for _, s := range m.store().ents {
		if !s.static {
			continue
		}
		smn, smx := m.aabbWorldMinMax(s)
		closest := rl.Vector3{
			X: float32(math.Max(float64(smn.X), math.Min(float64(wp.X), float64(smx.X)))),
			Y: float32(math.Max(float64(smn.Y), math.Min(float64(wp.Y), float64(smx.Y)))),
			Z: float32(math.Max(float64(smn.Z), math.Min(float64(wp.Z), float64(smx.Z)))),
		}
		d := rl.Vector3Distance(wp, closest)
		if d < r && d > 1e-6 {
			n := rl.Vector3Subtract(wp, closest)
			n = rl.Vector3Normalize(n)
			pen := r - d
			nwp := rl.Vector3Add(wp, rl.Vector3Scale(n, pen))
			m.setLocalFromWorld(e, nwp.X, nwp.Y, nwp.Z)
			e.hasHit = true
			e.hitX, e.hitY, e.hitZ = closest.X, closest.Y, closest.Z
			e.hitNX, e.hitNY, e.hitNZ = n.X, n.Y, n.Z
			if e.slide {
				vn := rl.Vector3Scale(n, rl.Vector3DotProduct(e.vel, n))
				e.vel = rl.Vector3Subtract(e.vel, vn)
			}
			dot := n.Y
			if math.Abs(float64(dot)) < 0.4 {
				fr := e.friction
				if fr <= 0 {
					fr = 0.9
				}
				e.vel.X *= fr
				e.vel.Z *= fr
			}
		} else if d <= 1e-6 {
			nwp := wp
			nwp.Y = smx.Y + r + 0.01
			m.setLocalFromWorld(e, nwp.X, nwp.Y, nwp.Z)
		}
	}
}

func (m *Module) resolveBoxVsStatics(e *ent) {
	dmn, dmx := m.aabbWorldMinMax(e)
	for _, s := range m.store().ents {
		if !s.static {
			continue
		}
		smn, smx := m.aabbWorldMinMax(s)
		if dmx.X < smn.X || dmn.X > smx.X || dmx.Y < smn.Y || dmn.Y > smx.Y || dmx.Z < smn.Z || dmn.Z > smx.Z {
			continue
		}
		// minimal penetration axis
		ox := minFloat32(smx.X-dmn.X, dmx.X-smn.X)
		oy := minFloat32(smx.Y-dmn.Y, dmx.Y-smn.Y)
		oz := minFloat32(smx.Z-dmn.Z, dmx.Z-smn.Z)
		wc := m.worldPos(e)
		switch {
		case ox <= oy && ox <= oz:
			nwc := wc
			if wc.X < m.worldPos(s).X {
				nwc.X -= ox
			} else {
				nwc.X += ox
			}
			m.setLocalFromWorld(e, nwc.X, nwc.Y, nwc.Z)
			e.vel.X = 0
		case oy <= ox && oy <= oz:
			nwc := wc
			if wc.Y < m.worldPos(s).Y {
				nwc.Y -= oy
				if e.vel.Y > 0 {
					e.vel.Y = 0
				}
			} else {
				nwc.Y += oy
				if e.vel.Y < 0 {
					e.vel.Y = 0
				}
			}
			m.setLocalFromWorld(e, nwc.X, nwc.Y, nwc.Z)
		default:
			nwc := wc
			if wc.Z < m.worldPos(s).Z {
				nwc.Z -= oz
			} else {
				nwc.Z += oz
			}
			m.setLocalFromWorld(e, nwc.X, nwc.Y, nwc.Z)
			e.vel.Z = 0
		}
	}
}

func minFloat32(a, b float32) float32 {
	if a < b {
		return a
	}
	return b
}

func (m *Module) pairwiseDynamic() {
	ids := make([]int64, 0, len(m.store().ents))
	for id := range m.store().ents {
		ids = append(ids, id)
	}
	for i := 0; i < len(ids); i++ {
		for j := i + 1; j < len(ids); j++ {
			a := m.store().ents[ids[i]]
			b := m.store().ents[ids[j]]
			if a.static || b.static {
				continue
			}
			if !a.useSphere || !b.useSphere {
				continue
			}
			pa := m.worldPos(a)
			pb := m.worldPos(b)
			d := rl.Vector3Distance(pa, pb)
			sum := a.radius + b.radius
			if d < sum && d > 1e-6 {
				n := rl.Vector3Subtract(pa, pb)
				n = rl.Vector3Normalize(n)
				pen := sum - d
				npa := rl.Vector3Add(pa, rl.Vector3Scale(n, pen*0.5))
				npb := rl.Vector3Subtract(pb, rl.Vector3Scale(n, pen*0.5))
				m.setLocalFromWorld(a, npa.X, npa.Y, npa.Z)
				m.setLocalFromWorld(b, npb.X, npb.Y, npb.Z)
				a.collided = true
				b.collided = true
				a.otherID = b.id
				b.otherID = a.id
			}
		}
	}
}

func entTint(e *ent) rl.Color {
	a := e.alpha
	if a < 0 {
		a = 0
	}
	if a > 1 {
		a = 1
	}
	return rl.Color{R: e.r, G: e.g, B: e.b, A: uint8(a * 255)}
}

func (m *Module) entDrawAll(args []value.Value) (value.Value, error) {
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("ENTITY.DRAWALL expects 0 arguments")
	}
	st := m.store()
	drawList := make([]*ent, 0, len(st.ents))
	for _, e := range st.ents {
		if e == nil || e.hidden || e.kind == entKindEmpty {
			continue
		}
		switch e.kind {
		case entKindBox, entKindSphere, entKindCylinder, entKindPlane, entKindMesh, entKindModel:
			drawList = append(drawList, e)
		}
	}
	sort.Slice(drawList, func(i, j int) bool {
		return drawList[i].drawOrder < drawList[j].drawOrder
	})
	for _, e := range drawList {
		wp := m.worldPos(e)
		col := entTint(e)
		switch e.kind {
		case entKindBox:
			rl.DrawCube(wp, e.w*e.scale.X, e.h*e.scale.Y, e.d*e.scale.Z, col)
		case entKindSphere:
			sx, sy, sz := e.scale.X, e.scale.Y, e.scale.Z
			ms := sx
			if sy > ms {
				ms = sy
			}
			if sz > ms {
				ms = sz
			}
			rad := e.radius * ms
			if rad <= 1e-6 {
				rad = 0.01
			}
			rings := e.segH
			slices := e.segV
			if rings < 8 {
				rings = 16
			}
			if slices < 8 {
				slices = 16
			}
			rl.DrawSphereEx(wp, rad, rings, slices, col)
		case entKindCylinder:
			h := e.cylH * e.scale.Y
			rs := e.scale.X
			if e.scale.Z > rs {
				rs = e.scale.Z
			}
			rt := e.radius * rs
			slices := e.segV
			if slices < 3 {
				slices = 16
			}
			rl.DrawCylinder(wp, rt, rt, h, slices, col)
		case entKindPlane:
			sx := e.w * e.scale.X
			sz := e.d * e.scale.Z
			if sx <= 1e-6 {
				sx = 1
			}
			if sz <= 1e-6 {
				sz = 1
			}
			rl.DrawPlane(wp, rl.Vector2{X: sx, Y: sz}, col)
		case entKindMesh, entKindModel:
			if !e.hasRLModel {
				continue
			}
			q := m.worldRotQuat(e)
			var axis rl.Vector3
			var ang float32
			rl.QuaternionToAxisAngle(q, &axis, &ang)
			sc := rl.Vector3{X: e.scale.X, Y: e.scale.Y, Z: e.scale.Z}
			rl.DrawModelEx(e.rlModel, wp, axis, ang, sc, col)
		}
	}
	return value.Nil, nil
}

func argHandle(v value.Value) (heap.Handle, bool) {
	if v.Kind != value.KindHandle {
		return 0, false
	}
	return heap.Handle(v.IVal), true
}

func argF32(v value.Value) (float32, bool) {
	if f, ok := v.ToFloat(); ok {
		return float32(f), true
	}
	if i, ok := v.ToInt(); ok {
		return float32(i), true
	}
	return 0, false
}
