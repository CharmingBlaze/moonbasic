//go:build cgo || (windows && !cgo)

package mbentity

import (
	"fmt"
	"math"
	"sort"

	"moonbasic/runtime"
	mbcamera "moonbasic/runtime/camera"
	"moonbasic/runtime/mbgame"
	"moonbasic/runtime/mblight"
	"moonbasic/runtime/mbmatrix"
	"moonbasic/runtime/mbmodel3d"
	mbtime "moonbasic/runtime/time"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var entityStores = make(map[*heap.Store]*entityStore)

type collRule struct {
	src, dst int32
	method   int32 // 1=Sph-Sph, 2=Sph-Box, 3=Sph-Mesh
	response int32 // 1=Stop, 2=Slide, 3=Slide (no grav)
}

type scriptBindRec struct {
	pattern string
	fnName  string
}

type entityStore struct {
	ents   map[int64]*ent
	nextID int64
	byName map[string]int64
	// children maps parent entity id -> ordered list of direct child ids (insertion order).
	children map[int64][]int64
	groups   map[string]map[int64]struct{}
	rules    []collRule

	// entMeta: flattened glTF extras strings per entity (Blender custom properties).
	entMeta map[int64]map[string]string
	// scriptBinds: glob patterns → BASIC function names (dispatch via LEVEL.MATCHSCRIPTBIND).
	scriptBinds []scriptBindRec

	// msgQueues: ENTITY.SENDMESSAGE / ENTITY.POLLMESSAGE (FIFO strings per receiver entity).
	msgQueues map[int64][]string

	// LEVEL.* glTF scene graph (see entity_level_cgo.go)
	levelRoot     string
	levelMarkers  map[string]rl.Vector3
	levelSpawn    map[string]rl.Matrix
	levelLayers   map[string][]int64
	levelColliders []levelColliderRec

	// DOD / SoA Spatial Buffers
	spatial runtime.SpatialBuffer

	// Optimization cache: entities filtered by type
	dynamicEnts []*ent
	staticEnts  []*ent

	// Shared primitive meshes (unit size)
	unitCube   rl.Mesh
	unitSphere rl.Mesh
	unitCyl    rl.Mesh
	unitCone   rl.Mesh
	unitPlane  rl.Mesh
	unitMat    rl.Material
}

type levelColliderRec struct {
	Name  string
	World rl.Matrix
}

func (m *Module) store() *entityStore {
	if m.h == nil {
		// Fallback for registration-time if heap not yet bound
		// (though Registry.RegisterModule binds it first).
		panic("mbentity: store requested before heap bound")
	}
	s := entityStores[m.h]
	if s == nil {
		s = &entityStore{
			ents:     make(map[int64]*ent),
			nextID:   1,
			byName:   make(map[string]int64),
			children: make(map[int64][]int64),
		}
		// Initialize SoA with small capacity, will grow on demand
		const initialCap = 256
		s.spatial.X = make([]float32, initialCap)
		s.spatial.Y = make([]float32, initialCap)
		s.spatial.Z = make([]float32, initialCap)
		s.spatial.P = make([]float32, initialCap)
		s.spatial.W = make([]float32, initialCap)
		s.spatial.R = make([]float32, initialCap)

		// Initialize shared primitives
		s.unitCube = rl.GenMeshCube(1, 1, 1)
		s.unitSphere = rl.GenMeshSphere(1, 16, 16)
		s.unitCyl = rl.GenMeshCylinder(1, 1, 16)
		s.unitCone = rl.GenMeshCone(1, 1, 16)
		s.unitPlane = rl.GenMeshPlane(1, 1, 1, 1)
		s.unitMat = rl.LoadMaterialDefault()
		
		entityStores[m.h] = s
	}
	
	// Link master registry spatial pointer for VM fast-path
	if rt, ok := m.reg.(*runtime.Registry); ok {
		rt.Spatial = &s.spatial
	}
	
	return s
}

func (s *entityStore) ensureSlices(id int) {
	if id < len(s.spatial.X) {
		return
	}
	newCap := len(s.spatial.X) * 2
	if id >= newCap {
		newCap = id + 256
	}
	s.spatial.X = append(s.spatial.X, make([]float32, newCap-len(s.spatial.X))...)
	s.spatial.Y = append(s.spatial.Y, make([]float32, newCap-len(s.spatial.Y))...)
	s.spatial.Z = append(s.spatial.Z, make([]float32, newCap-len(s.spatial.Z))...)
	s.spatial.P = append(s.spatial.P, make([]float32, newCap-len(s.spatial.P))...)
	s.spatial.W = append(s.spatial.W, make([]float32, newCap-len(s.spatial.W))...)
	s.spatial.R = append(s.spatial.R, make([]float32, newCap-len(s.spatial.R))...)
}

func (m *Module) Register(r runtime.Registrar) {
	m.reg = r
	r.Register("ENTITY.CREATE", "entity", runtime.AdaptLegacy(m.entCreate))
	r.Register("ENTITY.CREATEENTITY", "entity", runtime.AdaptLegacy(m.entCreate))
	r.Register("ENTITY.CREATEBOX", "entity", runtime.AdaptLegacy(m.entCreateBox))
	r.Register("ENTITY.CREATECUBE", "entity", runtime.AdaptLegacy(m.entCreateBox))
	registerEntityBlitzAPI(m, r)
	registerBrushBlitzAPI(m, r)
	registerEntitySceneGroupAPI(m, r)
	registerEntityUnifiedModelAPI(m, r)
	registerEntityMaterialScrollAPI(m, r)
	registerEntitySpriteAnimAPI(m, r)
	registerLevelGLTFAPI(m, r)
	registerEntityInteractionAPI(m, r)
	registerEntityEasyAPI(m, r)
	r.Register("ENTITY.SETPOSITION", "entity", runtime.AdaptLegacy(m.entSetPosition))
	r.Register("ENTITY.POSITION", "entity", runtime.AdaptLegacy(m.entSetPosition))
	r.Register("ENTITY.GETPOSITION", "entity", runtime.AdaptLegacy(m.entGetPosition))
	r.Register("ENTITY.GETPOS", "entity", runtime.AdaptLegacy(m.entGetPos))
	r.Register("ENTITY.GETXZ", "entity", runtime.AdaptLegacy(m.entGetXZ))
	r.Register("TERRAIN.SNAPY", "terrain", runtime.AdaptLegacy(m.entTerrainSnapY))
	r.Register("TERRAIN.PLACE", "terrain", runtime.AdaptLegacy(m.entTerrainPlace))
	r.Register("ENTITY.CLAMPTOTERRAIN", "entity", runtime.AdaptLegacy(m.entClampToTerrain))
	r.Register("ENTITY.FREEENTITIES", "entity", runtime.AdaptLegacy(m.entFreeEntities))
	r.Register("ENTITY.MOVE", "entity", runtime.AdaptLegacy(m.entMove))
	r.Register("ENTITY.TRANSLATE", "entity", runtime.AdaptLegacy(m.entTranslate))
	r.Register("ENTITY.ROTATE", "entity", runtime.AdaptLegacy(m.entRotate))
	r.Register("ENTITY.TURN", "entity", runtime.AdaptLegacy(m.entRotate))
	r.Register("ENTITY.SCALE", "entity", runtime.AdaptLegacy(m.entScale))
	r.Register("ENTITY.COLOR", "entity", runtime.AdaptLegacy(m.entColor))
	r.Register("ENTITY.RADIUS", "entity", runtime.AdaptLegacy(m.entRadius))
	r.Register("ENTITY.BOX", "entity", runtime.AdaptLegacy(m.entBox))
	r.Register("ENTITY.TYPE", "entity", runtime.AdaptLegacy(m.entType))
	r.Register("ENTITYTYPE", "entity", runtime.AdaptLegacy(m.entType))
	r.Register("COLLISIONS", "entity", runtime.AdaptLegacy(m.entCollisions))
	r.Register("ENTITYRADIUS", "entity", runtime.AdaptLegacy(m.entRadius))
	r.Register("SPHERECOLLIDE", "entity", runtime.AdaptLegacy(m.entRadius))

	r.Register("ENTITY.UPDATE", "entity", runtime.AdaptLegacy(m.entUpdate))
	r.Register("RESETENTITY", "entity", runtime.AdaptLegacy(m.entReset))
	r.Register("ENTITYCOLLIDED", "entity", runtime.AdaptLegacy(m.entCollidedType))
	r.Register("EntityHitsType", "entity", runtime.AdaptLegacy(m.entEntityHitsType))
	r.Register("COUNTCOLLISIONS", "entity", runtime.AdaptLegacy(m.entCountCollisions))
	r.Register("GETCOLLISIONENTITY", "entity", runtime.AdaptLegacy(m.entGetCollisionEntity))

	r.Register("ENTITY.DRAWALL", "entity", runtime.AdaptLegacy(m.entDrawAll))
	r.Register("ENTITY.DRAW", "entity", runtime.AdaptLegacy(m.entDraw))
	r.Register("ENTITY.SETCULLMODE", "entity", runtime.AdaptLegacy(m.entSetCullMode))
	r.Register("CAMERA.FOLLOWENTITY", "entity", m.camFollowEntity)
	registerPhysicsEntitySync(m, r)
	registerJoltEntityCollisionAPI(m, r)
	registerEntityGameplayIntelAPI(m, r)
	registerEntityTweenAPI(m, r)
	registerEntityQoLAPI(m, r)
	registerEntityAIAPI(m, r)
	registerEntityPhysicsQoLAPI(m, r)
	registerEntityEnvQoLAPI(m, r)
	registerEntityPhysicsMacroAPI(m, r)
	registerBlitzEntityHandles(m, r)

	// Wire registry utility callbacks
	if rt, ok := r.(*runtime.Registry); ok {
		rt.EntityIDActive = func(id int64) bool {
			_, ok := m.store().ents[id]
			return ok
		}
		rt.ResolveEntityWorldPos = func(id int64) (rl.Vector3, bool) {
			e := m.store().ents[id]
			if e == nil {
				return rl.Vector3{}, false
			}
			return m.worldPos(e), true
		}

		rt.FastEntityPropGet = func(id int64, propID int) (value.Value, error) {
			e := m.store().ents[id]
			if e == nil {
				return value.Nil, fmt.Errorf("ENTITY.PROP_GET: unknown entity %d", id)
			}
			ep, ew, er := e.getRot()
			pos := e.getPos()
			switch propID {
			case 0: return value.FromFloat(float64(pos.X)), nil
			case 1: return value.FromFloat(float64(pos.Y)), nil
			case 2: return value.FromFloat(float64(pos.Z)), nil
			case 3: return value.FromFloat(float64(ep)), nil
			case 4: return value.FromFloat(float64(ew)), nil
			case 5: return value.FromFloat(float64(er)), nil
			case 6: return value.FromFloat(float64(e.scale.X)), nil
			case 7: return value.FromFloat(float64(e.scale.Y)), nil
			case 8: return value.FromFloat(float64(e.scale.Z)), nil
			case 9: return value.FromFloat(float64(e.alpha)), nil
			case 10: return value.FromBool(e.hidden), nil
			default: return value.Nil, fmt.Errorf("ENTITY.PROP_GET: invalid prop id %d", propID)
			}
		}

		rt.FastEntityPropSet = func(id int64, propID int, val value.Value) error {
			e := m.store().ents[id]
			if e == nil {
				return fmt.Errorf("ENTITY.PROP_SET: unknown entity %d", id)
			}
			f, _ := val.ToFloat()
			switch propID {
			case 0: 
				vpos := e.getPos()
				vpos.X = float32(f)
				e.setPos(vpos)
			case 1:
				vpos := e.getPos()
				vpos.Y = float32(f)
				e.setPos(vpos)
			case 2:
				vpos := e.getPos()
				vpos.Z = float32(f)
				e.setPos(vpos)
			case 3:
				_, ew, er := e.getRot()
				e.setRot(float32(f), ew, er)
			case 4:
				ep, _, er := e.getRot()
				e.setRot(ep, float32(f), er)
			case 5:
				ep, ew, _ := e.getRot()
				e.setRot(ep, ew, float32(f))
			case 6: e.scale.X = float32(f)
			case 7: e.scale.Y = float32(f)
			case 8: e.scale.Z = float32(f)
			case 9: e.alpha = float32(f)
			case 10: e.hidden = val.IVal != 0
			default: return fmt.Errorf("ENTITY.PROP_SET: invalid prop id %d", propID)
			}
			return nil
		}
	}

	mblight.SetLightFollowWorldPosGetter(func(id int64) (float32, float32, float32, bool) {
		e := m.store().ents[id]
		if e == nil {
			return 0, 0, 0, false
		}
		wp := m.worldPos(e)
		return wp.X, wp.Y, wp.Z, true
	})
}

// Shutdown implements runtime.Module.
func (m *Module) Shutdown() {
	clearEntityRefFreeHookIfOwner(m)
	if m.h != nil {
		s := entityStores[m.h]
		if s != nil {
			rl.UnloadMesh(&s.unitCube)
			rl.UnloadMesh(&s.unitSphere)
			rl.UnloadMesh(&s.unitCyl)
			rl.UnloadMesh(&s.unitCone)
			rl.UnloadMesh(&s.unitPlane)
			rl.UnloadMaterial(s.unitMat)
		}
		delete(entityStores, m.h)
		delete(ModulesByStore, m.h)
	}
}

// Reset implements runtime.Module.
func (m *Module) Reset() {
	st := m.store()
	st.ents = make(map[int64]*ent)
	st.dynamicEnts = nil
	st.staticEnts = nil
	st.children = make(map[int64][]int64)
	st.nextID = 1
	if st.byName != nil {
		st.byName = make(map[string]int64)
	}
}

// ForEachStatic iterates over all entities marked as 'static' and calls the provided function
// with their world-space Axis-Aligned Bounding Box.
func (m *Module) ForEachStatic(fn func(id int64, worldAABB rl.BoundingBox)) {
	st := m.store()
	for id, e := range st.ents {
		if !e.static || e.hidden {
			continue
		}
		mn, mx := m.aabbWorldMinMax(e)
		fn(id, rl.NewBoundingBox(mn, mx))
	}
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
	_, ew, _ := e.getRot()
	err := mbcamera.ThirdPersonFollowStep(m.h, ch, wp.X, wp.Y, wp.Z, ew, dist, height, smooth, dt)
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
	st.ensureSlices(int(id))
	e := newDefaultEnt(id, &st.spatial)
	e.kind = entKindSphere
	e.w, e.h, e.d = 1, 1, 1
	e.radius = 0.5
	e.useSphere = true
	e.static = false
	e.gravity = -28
	st.ents[id] = e
	st.dynamicEnts = append(st.dynamicEnts, e)
	fmt.Println("MB_DEBUG: CREATESPHERE ID =", id)
	return value.FromInt(id), nil
}

func (m *Module) entCreateBox(args []value.Value) (value.Value, error) {
	if len(args) == 1 {
		if _, ok := argF32(args[0]); !ok {
			return value.Nil, fmt.Errorf("ENTITY.CREATEBOX: size must be numeric")
		}
		return m.entCreateBox([]value.Value{args[0], args[0], args[0]})
	}
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("ENTITY.CREATEBOX expects 1 argument (size#) or 3 (w#, h#, d#)")
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
	st.ensureSlices(int(id))
	e := newDefaultEnt(id, &st.spatial)
	e.kind = entKindBox
	e.r, e.g, e.b = 180, 180, 200
	e.w, e.h, e.d = w, h, d
	e.static = true
	e.useSphere = false
	e.gravity = 0
	st.ents[id] = e
	st.staticEnts = append(st.staticEnts, e)
	fmt.Println("MB_DEBUG: CREATEBOX ID =", id)
	return value.FromInt(id), nil
}

func (m *Module) entSetPosition(args []value.Value) (value.Value, error) {
	if len(args) != 4 && len(args) != 5 {
		return value.Nil, fmt.Errorf("ENTITY.SETPOSITION expects 4–5 arguments (entity#, x#, y#, z# [, global])")
	}
	id, ok := m.entID(args[0])
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
		e.setPos(rl.Vector3{X: x, Y: y, Z: z})
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
	id, ok := m.entID(args[0])
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

func (m *Module) entGetPos(args []value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("ENTITY.GETPOS: heap not bound")
	}
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("ENTITY.GETPOS expects entity#")
	}
	id, ok := m.entID(args[0])
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("ENTITY.GETPOS: invalid entity")
	}
	e := m.store().ents[id]
	if e == nil {
		return value.Nil, fmt.Errorf("ENTITY.GETPOS: unknown entity %d", id)
	}
	p := m.worldPos(e)
	arr, err := heap.NewArrayOfKind([]int64{3}, heap.ArrayKindFloat, 0)
	if err != nil {
		return value.Nil, err
	}
	arr.Floats[0] = float64(p.X)
	arr.Floats[1] = float64(p.Y)
	arr.Floats[2] = float64(p.Z)
	h, err := m.h.Alloc(arr)
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(h), nil
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
	id, ok := m.entID(args[0])
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
	ep, ew, _ := e.getRot()
	fwd, right, up := localAxes(ew, ep)
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
	id, ok := m.entID(args[0])
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
	id, ok := m.entID(args[0])
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
	ep2, ew2, er2 := e.getRot()
	e.setRot(ep2+dp, ew2+dy, er2+dr)
	return value.Nil, nil
}

func (m *Module) entScale(args []value.Value) (value.Value, error) {
	if len(args) != 4 {
		return value.Nil, fmt.Errorf("ENTITY.SCALE expects 4 arguments (entity#, sx#, sy#, sz#)")
	}
	id, ok := m.entID(args[0])
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("ENTITY.SCALE: invalid entity")
	}
	st := m.store()
	e := st.ents[id]
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
	if len(args) == 2 && args[1].Kind == value.KindHandle {
		id, ok := m.entID(args[0])
		if !ok || id < 1 {
			return value.Nil, fmt.Errorf("ENTITY.COLOR: invalid entity")
		}
		e := m.store().ents[id]
		if e == nil {
			return value.Nil, fmt.Errorf("ENTITY.COLOR: unknown entity %d", id)
		}
		c, err := mbmatrix.HeapColorRGBA(m.h, heap.Handle(args[1].IVal))
		if err != nil {
			return value.Nil, fmt.Errorf("ENTITY.COLOR: %w", err)
		}
		e.r = c.R
		e.g = c.G
		e.b = c.B
		e.alpha = float32(c.A) / 255.0
		return value.Nil, nil
	}
	if len(args) != 4 && len(args) != 5 {
		return value.Nil, fmt.Errorf("ENTITY.COLOR expects 2 arguments (entity#, colorHandle) or 4–5 arguments (entity#, r, g, b [, a])")
	}
	id, ok := m.entID(args[0])
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
	if len(args) == 5 {
		ai, ok4 := args[4].ToInt()
		if ok4 {
			e.alpha = float32(ai) / 255.0
		} else if af, ok5 := args[4].ToFloat(); ok5 {
			e.alpha = float32(af) / 255.0
		}
	}
	return value.Nil, nil
}

func (m *Module) entRadius(args []value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("ENTITY.RADIUS expects 2 arguments (entity#, radius#)")
	}
	id, ok := m.entID(args[0])
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
	id, ok := m.entID(args[0])
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
	id, ok := m.entID(args[0])
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("ENTITY.COLLIDED: invalid entity")
	}
	e := m.store().ents[id]
	if e == nil {
		return value.Nil, fmt.Errorf("ENTITY.COLLIDED: unknown entity %d", id)
	}
	return value.FromBool(e.getExt().collided), nil
}

func (m *Module) entCollisionOther(args []value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("ENTITY.COLLISIONOTHER expects entity#")
	}
	id, ok := m.entID(args[0])
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("ENTITY.COLLISIONOTHER: invalid entity")
	}
	e := m.store().ents[id]
	if e == nil {
		return value.Nil, fmt.Errorf("ENTITY.COLLISIONOTHER: unknown entity %d", id)
	}
	return value.FromInt(e.getExt().otherID), nil
}

func (m *Module) entFloor(args []value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("ENTITY.FLOOR expects entity#")
	}
	id, ok := m.entID(args[0])
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
		sp := s.getPos()
		bx, by, bz := float64(sp.X), float64(sp.Y), float64(sp.Z)
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
	id, ok := m.entID(args[0])
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
	id, ok := m.entID(args[0])
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
	if e.jumpGrounded {
		e.vel.Y += f
		e.onGround = false
		e.groundCoyoteLeft = 0
	}
	return value.Nil, nil
}

func (m *Module) entUpdate(args []value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("ENTITY.UPDATE expects (dt#)")
	}
	dt, ok := argF32(args[0])
	if !ok {
		return value.Nil, fmt.Errorf("ENTITY.UPDATE: dt must be numeric")
	}
	if dt <= 0 {
		// First frame / vsync edge / closing window can yield 0 delta; skip tick instead of erroring out.
		return value.Nil, nil
	}

	m.advanceSpriteImageSequences(dt)

	m.processEntityTweens(dt)
	m.processAITasks(dt)

	st := m.store()

	for _, e := range st.ents {
		if e == nil || e.ext == nil { continue }
		ext := e.ext
		if ext.ghostMode && ext.ghostTimer > 0 {
			ext.ghostTimer -= dt
			if ext.ghostTimer <= 0 {
				ext.ghostTimer = 0
				ext.ghostMode = false
				// Logic to restore Jolt collision mask would go here
			}
		}
	}
	
	// 1. Clear frame collision state
	for _, e := range st.ents {
		if e == nil || e.ext == nil {
			continue
		}
		ext := e.ext
		ext.collided = false
		ext.hits = ext.hits[:0]
		ext.hitPos = ext.hitPos[:0]
		ext.hitN = ext.hitN[:0]
		ext.hasHit = false
	}

	// 2. Discover and update all particle emitters in the heap
	if m.h != nil {
		partHandles := m.h.FilterByType(heap.TagParticle)
		for _, h := range partHandles {
			if obj, ok := m.h.Get(h); ok {
				// We know it's a *particleObj based on TagParticle
				if p, ok := obj.(interface{ Update(float32) }); ok {
					p.Update(dt)
				}
			}
		}
	}

	// 3. Update entity positions (gravity + velocity)
	for _, e := range st.ents {
		if e.static {
			continue
		}
		e.vel.Y += e.gravity * dt
		wp := m.worldPos(e)
		nw := rl.Vector3Add(wp, rl.Vector3Scale(e.vel, dt))
		m.setLocalFromWorld(e, nw.X, nw.Y, nw.Z)

		// Basic ground snapping / floor check
		if e.useSphere {
			wp2 := m.worldPos(e)
			px, py, pz := float64(wp2.X), float64(wp2.Y), float64(wp2.Z)
			pvy := float64(e.vel.Y)
			pr := float64(e.radius)
			var bestSnap float64
			found := false
			for _, s := range st.ents {
				if !s.static {
					continue
				}
				sp := s.getPos()
				bx, by, bz := float64(sp.X), float64(sp.Y), float64(sp.Z)
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
			} else {
				e.onGround = false
			}
		}
	}

	// 4. Resolve Global Collision Rules
	m.resolveRules()

	// 4b. Floor contact from rule hits (supports platforms when snap pass missed).
	for _, e := range st.ents {
		if e == nil {
			continue
		}
		m.refreshGroundFromRules(e)
	}

	// 5. Update Animations (skeletal clips + GPU bone matrices for sockets)
	for _, e := range st.ents {
		if e == nil || !e.hasRLModel || e.ext == nil || len(e.ext.modelAnims) == 0 {
			continue
		}
		ext := e.ext
		ai := ext.animIndex
		if ai < 0 || int(ai) >= len(ext.modelAnims) {
			ai = 0
		}
		anim := ext.modelAnims[ai]
		if anim.FrameCount <= 0 {
			continue
		}
		if ext.animSpeed != 0 {
			ext.animTime += dt * ext.animSpeed * 30
		}
		frame := pickAnimFrame(e, anim)
		rl.UpdateModelAnimation(e.rlModel, anim, frame)
		rl.UpdateModelAnimationBones(e.rlModel, anim, frame)
	}
	m.syncBoneSockets()
	if m.h != nil {
		mblight.SyncPointFollowLights(m.h)
	}

	// 6. Jump coyote: brief grace after leaving a ledge (ENTITY.GROUNDED / EntityGrounded / ENTITY.JUMP).
	for _, e := range st.ents {
		if e == nil || e.static {
			continue
		}
		finalizeJumpGrounded(e, e.onGround)
	}
	return value.Nil, nil
}

const groundCoyoteMax = 2

func finalizeJumpGrounded(e *ent, touching bool) {
	c := e.groundCoyoteLeft
	e.jumpGrounded = touching || c > 0
	if touching {
		e.groundCoyoteLeft = groundCoyoteMax
	} else if c > 0 {
		e.groundCoyoteLeft = c - 1
	}
}

func (m *Module) refreshGroundFromRules(e *ent) {
	if e.static || !e.useSphere || e.onGround {
		return
	}
	ext := e.getExt()
	for i := range ext.hits {
		if i >= len(ext.hitN) {
			break
		}
		n := ext.hitN[i]
		if n.Y < 0.45 {
			continue
		}
		e.onGround = true
		return
	}
}

func (m *Module) aabbWorldMinMax(e *ent) (mn, mx rl.Vector3) {
	wm := m.worldMatrix(e)
	hx, hy, hz := e.w*0.5, e.h*0.5, e.d*0.5
	if e.kind == entKindSphere {
		hx, hy, hz = e.radius, e.radius, e.radius
	} else if e.kind == entKindCylinder || e.kind == entKindCone {
		hx, hy, hz = e.radius, e.cylH*0.5, e.radius
	}

	corners := [8]rl.Vector3{
		{X: -hx, Y: -hy, Z: -hz}, {X: hx, Y: -hy, Z: -hz},
		{X: -hx, Y: hy, Z: -hz}, {X: hx, Y: hy, Z: -hz},
		{X: -hx, Y: -hy, Z: hz}, {X: hx, Y: -hy, Z: hz},
		{X: -hx, Y: hy, Z: hz}, {X: hx, Y: hy, Z: hz},
	}

	mn = rl.Vector3{X: 1e30, Y: 1e30, Z: 1e30}
	mx = rl.Vector3{X: -1e30, Y: -1e30, Z: -1e30}

	for i := 0; i < 8; i++ {
		p := rl.Vector3Transform(corners[i], wm)
		if p.X < mn.X { mn.X = p.X }
		if p.Y < mn.Y { mn.Y = p.Y }
		if p.Z < mn.Z { mn.Z = p.Z }
		if p.X > mx.X { mx.X = p.X }
		if p.Y > mx.Y { mx.Y = p.Y }
		if p.Z > mx.Z { mx.Z = p.Z }
	}
	return mn, mx
}

func (m *Module) resolveSphereVsStatics(e *ent) {
	r := e.radius
	if r <= 0 {
		return
	}
	wp := m.worldPos(e)
	for _, s := range m.store().staticEnts {
		if s.hidden {
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
			ext := e.getExt()
			ext.hasHit = true
			ext.hitX, ext.hitY, ext.hitZ = closest.X, closest.Y, closest.Z
			ext.hitNX, ext.hitNY, ext.hitNZ = n.X, n.Y, n.Z
			if ext.slide {
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
	st := m.store()
	dyn := st.dynamicEnts
	for i := 0; i < len(dyn); i++ {
		for j := i + 1; j < len(dyn); j++ {
			a := dyn[i]
			b := dyn[j]
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
				a.getExt().collided = true
				b.getExt().collided = true
				a.getExt().otherID = b.id
				b.getExt().otherID = a.id
			}
		}
	}
}

func (m *Module) entDraw(args []value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("ENTITY.DRAW expects (entity#)")
	}
	id, ok := m.entID(args[0])
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("ENTITY.DRAW: invalid entity")
	}
	e := m.store().ents[id]
	if e == nil {
		return value.Nil, fmt.Errorf("ENTITY.DRAW: unknown entity")
	}
	m.drawOneEntity(e)
	return value.Nil, nil
}

// drawOneEntity renders one entity (primitives, billboard sprite, or skinned model). Used by ENTITY.DRAW and ENTITY.DRAWALL.
func (m *Module) drawOneEntity(e *ent) {
	if e == nil || e.hidden || e.kind == entKindEmpty {
		return
	}
	if !m.isVisible(e) {
		return
	}
	useBlend, blendM := m.entDrawBlendMode(e)
	if useBlend {
		rl.BeginBlendMode(blendM)
	}
	st := m.store()
	wm := m.worldMatrix(e)
	col := m.entTintResolved(e)
	st.unitMat.GetMap(int32(rl.MapAlbedo)).Color = col

	switch e.kind {
	case entKindBox:
		mProp := rl.MatrixScale(e.w, e.h, e.d)
		rl.DrawMesh(st.unitCube, st.unitMat, rl.MatrixMultiply(mProp, wm))
	case entKindSphere:
		mProp := rl.MatrixScale(e.radius*2, e.radius*2, e.radius*2)
		rl.DrawMesh(st.unitSphere, st.unitMat, rl.MatrixMultiply(mProp, wm))
	case entKindCylinder:
		mProp := rl.MatrixScale(e.radius*2, e.cylH, e.radius*2)
		rl.DrawMesh(st.unitCyl, st.unitMat, rl.MatrixMultiply(mProp, wm))
	case entKindCone:
		mProp := rl.MatrixScale(e.radius*2, e.cylH, e.radius*2)
		rl.DrawMesh(st.unitCone, st.unitMat, rl.MatrixMultiply(mProp, wm))
	case entKindPlane:
		mProp := rl.MatrixScale(e.w, 1.0, e.d)
		rl.DrawMesh(st.unitPlane, st.unitMat, rl.MatrixMultiply(mProp, wm))
	case entKindMesh, entKindModel:
		if !e.hasRLModel {
			if useBlend {
				rl.EndBlendMode()
			}
			return
		}
		wm := m.worldMatrix(e)
		saved := e.rlModel.Transform
		e.rlModel.Transform = wm

		// Outline Pass
		ext := e.getExt()
		if ext.outlineThickness > 0 {
			// Simple visual outline: draw slightly larger with solid color
			thick := ext.outlineThickness * 0.05
			outlineWM := rl.MatrixMultiply(rl.MatrixScale(1+thick, 1+thick, 1+thick), wm)
			e.rlModel.Transform = outlineWM
			mbmodel3d.DrawEntityModel(e.rlModel, ext.outlineColor)
			
			e.rlModel.Transform = wm
		}

		mbmodel3d.DrawEntityModel(e.rlModel, col)
		e.rlModel.Transform = saved
	default:
		if useBlend {
			rl.EndBlendMode()
		}
		return
	}
	if useBlend {
		rl.EndBlendMode()
	}
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
		if e.getExt().isSprite {
			drawList = append(drawList, e)
			continue
		}
		switch e.kind {
		case entKindBox, entKindSphere, entKindCylinder, entKindCone, entKindPlane, entKindMesh, entKindModel:
			drawList = append(drawList, e)
		}
	}
	sort.Slice(drawList, func(i, j int) bool {
		return drawList[i].drawOrder < drawList[j].drawOrder
	})
	for _, e := range drawList {
		m.drawOneEntity(e)
	}
	return value.Nil, nil
}

func (m *Module) isVisible(e *ent) bool {
	if e.cullMode == 1 { // Force Visible
		return true
	}
	if e.cullMode == 2 { // Force Hidden
		return false
	}
	// cullMode == 0 (Auto): CPU frustum vs active CAMERA.BEGIN (see mbcamera.ExtractFrustum).
	switch e.kind {
	case entKindSphere:
		wp := m.worldPos(e)
		ms := e.scale.X
		if e.scale.Y > ms {
			ms = e.scale.Y
		}
		if e.scale.Z > ms {
			ms = e.scale.Z
		}
		return mbcamera.SphereVisibleActive(wp.X, wp.Y, wp.Z, e.radius*ms)
	case entKindBox, entKindCone, entKindMesh, entKindModel:
		mn, mx := m.aabbWorldMinMax(e)
		return mbcamera.AABBVisibleActive(mn.X, mn.Y, mn.Z, mx.X, mx.Y, mx.Z)
	default:
		return true
	}
}

func (m *Module) entSetCullMode(args []value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("ENTITY.SETCULLMODE expects (entity#, mode): 0=Auto, 1=Force Visible, 2=Force Hidden")
	}
	id, ok := m.entID(args[0])
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("ENTITY.SETCULLMODE: invalid entity")
	}
	e := m.store().ents[id]
	if e == nil {
		return value.Nil, fmt.Errorf("ENTITY.SETCULLMODE: unknown entity %d", id)
	}
	mode, ok := args[1].ToInt()
	if !ok {
		return value.Nil, fmt.Errorf("ENTITY.SETCULLMODE: mode must be numeric")
	}
	e.cullMode = int32(mode)
	return value.Nil, nil
}

func (m *Module) entID(v value.Value) (int64, bool) {
	if v.Kind == value.KindHandle {
		h := heap.Handle(v.IVal)
		reg := runtime.ActiveRegistry()
		if reg == nil || reg.Heap == nil {
			fmt.Printf("[CRITICAL] entID: no active registry/heap for handle %d\n", h)
			return 0, false
		}
		obj, ok := reg.Heap.Get(h)
		if !ok {
			fmt.Printf("[DEBUG] entID: handle %d NOT found in heap (count=%d)\n", h, reg.Heap.Count())
			return 0, false
		}
		if obj.TypeTag() == heap.TagEntityRef {
			if ref, ok := obj.(interface{ GetID() int64 }); ok {
				return ref.GetID(), true
			}
			if er, ok := obj.(*heap.EntityRef); ok {
				return er.ID, true
			}
		}
		fmt.Printf("[DEBUG] entID: handle %d is NOT an EntityRef (tag=%d name=%s)\n", h, obj.TypeTag(), obj.TypeName())
		return 0, false
	}

	id, ok := v.ToInt()
	if !ok || id < 1 {
		return 0, false
	}
	st := m.store()
	if _, exists := st.ents[id]; exists {
		return id, true
	}
	return 0, false
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

func (m *Module) entCollisions(args []value.Value) (value.Value, error) {
	if len(args) != 4 {
		return value.Nil, fmt.Errorf("COLLISIONS expects (srcType, dstType, method, response)")
	}
	src, _ := args[0].ToInt()
	dst, _ := args[1].ToInt()
	meth, _ := args[2].ToInt()
	resp, _ := args[3].ToInt()
	m.store().rules = append(m.store().rules, collRule{
		src: int32(src), dst: int32(dst),
		method: int32(meth), response: int32(resp),
	})
	return value.Nil, nil
}

func (m *Module) entReset(args []value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("RESETENTITY expects entity#")
	}
	id, _ := m.entID(args[0])
	e := m.store().ents[id]
	if e == nil {
		return value.Nil, fmt.Errorf("unknown entity")
	}
	ext := e.getExt()
	e.vel = rl.Vector3{}
	ext.collided = false
	ext.hits = nil
	return value.Nil, nil
}

func (m *Module) entCollidedType(args []value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("ENTITYCOLLIDED expects (entity#, targetType)")
	}
	id, _ := m.entID(args[0])
	e := m.store().ents[id]
	if e == nil {
		return value.Nil, fmt.Errorf("unknown entity")
	}
	tp, _ := args[1].ToInt()
	ext := e.getExt()
	for _, rid := range ext.hits {
		re := m.store().ents[rid]
		if re != nil && re.getExt().collType == int32(tp) {
			return value.FromInt(rid), nil
		}
	}
	return value.FromInt(0), nil
}

func (m *Module) entEntityHitsType(args []value.Value) (value.Value, error) {
	v, err := m.entCollidedType(args)
	if err != nil {
		return value.Nil, err
	}
	id, _ := v.ToInt()
	return value.FromBool(id != 0), nil
}

func (m *Module) entCountCollisions(args []value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("COUNTCOLLISIONS expects entity#")
	}
	id, _ := m.entID(args[0])
	e := m.store().ents[id]
	if e == nil {
		return value.Nil, fmt.Errorf("unknown entity")
	}
	return value.FromInt(int64(len(e.getExt().hits))), nil
}

func (m *Module) entGetCollisionEntity(args []value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("GETCOLLISIONENTITY expects (entity#, index)")
	}
	id, _ := m.entID(args[0])
	e := m.store().ents[id]
	if e == nil {
		return value.Nil, fmt.Errorf("unknown entity")
	}
	idx, _ := args[1].ToInt()
	hits := e.getExt().hits
	if idx < 0 || int(idx) >= len(hits) {
		return value.FromInt(0), nil
	}
	return value.FromInt(hits[idx]), nil
}

func (m *Module) resolveRules() {
	st := m.store()
	if len(st.rules) == 0 {
		return
	}
	for _, r := range st.rules {
		m.applyRule(r)
	}
}

func (m *Module) applyRule(rule collRule) {
	st := m.store()
	var srcs []*ent
	var dsts []*ent
	for _, e := range st.ents {
		if e == nil { continue }
		ext := e.getExt()
		if ext.collType == rule.src {
			srcs = append(srcs, e)
		}
		if ext.collType == rule.dst {
			dsts = append(dsts, e)
		}
	}
	for _, s := range srcs {
		for _, d := range dsts {
			if s == d {
				continue
			}
			m.checkAndResolve(s, d, rule)
		}
	}
}

func (m *Module) checkAndResolve(s, d *ent, rule collRule) {
	ps := m.worldPos(s)
	switch rule.method {
	case 1: // Sph-Sph
		pd := m.worldPos(d)
		dist := rl.Vector3Distance(ps, pd)
		sum := s.radius + d.radius
		if dist < sum && dist > 1e-6 {
			n := rl.Vector3Normalize(rl.Vector3Subtract(ps, pd))
			pen := sum - dist
			contact := rl.Vector3Subtract(ps, rl.Vector3Scale(n, s.radius))
			m.resolveResponse(s, d, n, pen, rule.response, contact)
		}
	case 2: // Sph-Box
		mn, mx := m.aabbWorldMinMax(d)
		closest := rl.Vector3{
			X: float32(math.Max(float64(mn.X), math.Min(float64(ps.X), float64(mx.X)))),
			Y: float32(math.Max(float64(mn.Y), math.Min(float64(ps.Y), float64(mx.Y)))),
			Z: float32(math.Max(float64(mn.Z), math.Min(float64(ps.Z), float64(mx.Z)))),
		}
		dist := rl.Vector3Distance(ps, closest)
		if dist < s.radius && dist > 1e-6 {
			n := rl.Vector3Normalize(rl.Vector3Subtract(ps, closest))
			pen := s.radius - dist
			m.resolveResponse(s, d, n, pen, rule.response, closest)
		}
	}
}

func (m *Module) resolveResponse(s, d *ent, n rl.Vector3, pen float32, resp int32, contact rl.Vector3) {
	ps := m.worldPos(s)
	ps = rl.Vector3Add(ps, rl.Vector3Scale(n, pen))
	ext := s.getExt()
	ext.collided = true
	ext.hits = append(ext.hits, d.id)
	ext.hitPos = append(ext.hitPos, contact)
	ext.hitN = append(ext.hitN, n)
	ext.hasHit = true
	ext.hitX, ext.hitY, ext.hitZ = contact.X, contact.Y, contact.Z
	ext.hitNX, ext.hitNY, ext.hitNZ = n.X, n.Y, n.Z
	if resp >= 2 { // Slide
		vn := rl.Vector3Scale(n, rl.Vector3DotProduct(s.vel, n))
		s.vel = rl.Vector3Subtract(s.vel, vn)
	} else if resp == 1 { // Stop
		s.vel = rl.Vector3{}
	}
}
