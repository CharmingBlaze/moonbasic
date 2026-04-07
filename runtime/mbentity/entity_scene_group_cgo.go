//go:build cgo || (windows && !cgo)

package mbentity

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	mbcamera "moonbasic/runtime/camera"
	"moonbasic/runtime"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func registerEntitySceneGroupAPI(m *Module, r runtime.Registrar) {
	r.Register("ENTITY.GROUPCREATE", "entity", m.entGroupCreate)
	r.Register("ENTITY.GROUPADD", "entity", m.entGroupAdd)
	r.Register("ENTITY.GROUPREMOVE", "entity", m.entGroupRemove)
	r.Register("ENTITY.ENTITIESINGROUP", "entity", m.entEntitiesInGroup)
	r.Register("ENTITY.ENTITIESINRADIUS", "entity", runtime.AdaptLegacy(m.entEntitiesInRadius))
	r.Register("ENTITY.ENTITIESINBOX", "entity", runtime.AdaptLegacy(m.entEntitiesInBox))
	r.Register("ENTITY.CLEARSCENE", "entity", runtime.AdaptLegacy(m.entClearScene))
	r.Register("ENTITY.SAVESCENE", "entity", m.entSaveScene)
	r.Register("ENTITY.LOADSCENE", "entity", m.entLoadScene)
	r.Register("SCENE.CLEARSCENE", "entity", runtime.AdaptLegacy(m.entClearScene))
	r.Register("SCENE.SAVESCENE", "entity", m.entSaveScene)
	r.Register("SCENE.LOADSCENE", "entity", m.entLoadScene)
	r.Register("CAMERA.SETTARGETENTITY", "entity", m.camSetTargetEntity)
	r.Register("CAMERA.CAMERAFOLLOW", "entity", m.camFollowEntity)
}

func groupKey(name string) string { return strings.ToUpper(strings.TrimSpace(name)) }

func (m *Module) entGroupCreate(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Kind != value.KindString {
		return value.Nil, fmt.Errorf("ENTITY.GROUPCREATE expects (name$)")
	}
	name, err := rt.ArgString(args, 0)
	if err != nil {
		return value.Nil, err
	}
	k := groupKey(name)
	if k == "" {
		return value.Nil, fmt.Errorf("ENTITY.GROUPCREATE: name must be non-empty")
	}
	st := m.store()
	if st.groups == nil {
		st.groups = make(map[string]map[int64]struct{})
	}
	if st.groups[k] == nil {
		st.groups[k] = make(map[int64]struct{})
	}
	return value.Nil, nil
}

func (m *Module) entGroupAdd(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("ENTITY.GROUPADD expects (groupName$, entity#)")
	}
	if args[0].Kind != value.KindString {
		return value.Nil, fmt.Errorf("ENTITY.GROUPADD: group name must be string")
	}
	gname, ok := rt.Heap.GetString(int32(args[0].IVal))
	if !ok {
		return value.Nil, fmt.Errorf("ENTITY.GROUPADD: invalid string")
	}
	eid, ok2 := args[1].ToInt()
	if !ok2 || eid < 1 {
		return value.Nil, fmt.Errorf("ENTITY.GROUPADD: invalid entity")
	}
	st := m.store()
	if st.ents[eid] == nil {
		return value.Nil, fmt.Errorf("ENTITY.GROUPADD: unknown entity")
	}
	k := groupKey(gname)
	if st.groups == nil {
		st.groups = make(map[string]map[int64]struct{})
	}
	if st.groups[k] == nil {
		st.groups[k] = make(map[int64]struct{})
	}
	st.groups[k][eid] = struct{}{}
	return value.Nil, nil
}

func (m *Module) entGroupRemove(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("ENTITY.GROUPREMOVE expects (groupName$, entity#)")
	}
	if args[0].Kind != value.KindString {
		return value.Nil, fmt.Errorf("ENTITY.GROUPREMOVE: group name must be string")
	}
	gname, ok := rt.Heap.GetString(int32(args[0].IVal))
	if !ok {
		return value.Nil, fmt.Errorf("ENTITY.GROUPREMOVE: invalid string")
	}
	eid, ok2 := args[1].ToInt()
	if !ok2 || eid < 1 {
		return value.Nil, fmt.Errorf("ENTITY.GROUPREMOVE: invalid entity")
	}
	st := m.store()
	k := groupKey(gname)
	if st.groups == nil || st.groups[k] == nil {
		return value.Nil, nil
	}
	delete(st.groups[k], eid)
	return value.Nil, nil
}

func (m *Module) allocFloatArray(vals []float64) (value.Value, error) {
	if m.h == nil {
		return value.Nil, fmt.Errorf("heap not bound")
	}
	if len(vals) == 0 {
		return value.Nil, nil
	}
	n := int64(len(vals))
	arr, err := heap.NewArray([]int64{n})
	if err != nil {
		return value.Nil, err
	}
	for i, v := range vals {
		_ = arr.Set([]int64{int64(i)}, v)
	}
	h, err := m.h.Alloc(arr)
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(h), nil
}

func (m *Module) entEntitiesInGroup(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Kind != value.KindString {
		return value.Nil, fmt.Errorf("ENTITY.ENTITIESINGROUP expects (groupName$)")
	}
	gname, err := rt.ArgString(args, 0)
	if err != nil {
		return value.Nil, err
	}
	k := groupKey(gname)
	st := m.store()
	if st.groups == nil || st.groups[k] == nil {
		return m.allocFloatArray(nil)
	}
	vals := make([]float64, 0, len(st.groups[k]))
	for id := range st.groups[k] {
		if st.ents[id] != nil {
			vals = append(vals, float64(id))
		}
	}
	return m.allocFloatArray(vals)
}

func (m *Module) entEntitiesInRadius(args []value.Value) (value.Value, error) {
	if len(args) != 4 {
		return value.Nil, fmt.Errorf("ENTITY.ENTITIESINRADIUS expects (x#, y#, z#, radius#)")
	}
	cx, ok1 := argF32(args[0])
	cy, ok2 := argF32(args[1])
	cz, ok3 := argF32(args[2])
	rad, ok4 := argF32(args[3])
	if !ok1 || !ok2 || !ok3 || !ok4 || rad < 0 {
		return value.Nil, fmt.Errorf("ENTITY.ENTITIESINRADIUS: numeric center and non-negative radius")
	}
	rad2 := float64(rad) * float64(rad)
	var vals []float64
	for id, e := range m.store().ents {
		if e == nil {
			continue
		}
		wp := m.worldPos(e)
		dx := float64(wp.X - cx)
		dy := float64(wp.Y - cy)
		dz := float64(wp.Z - cz)
		if dx*dx+dy*dy+dz*dz <= rad2 {
			vals = append(vals, float64(id))
		}
	}
	return m.allocFloatArray(vals)
}

func (m *Module) entEntitiesInBox(args []value.Value) (value.Value, error) {
	if len(args) != 6 {
		return value.Nil, fmt.Errorf("ENTITY.ENTITIESINBOX expects (x1#, y1#, z1#, x2#, y2#, z2#)")
	}
	x1, ok1 := argF32(args[0])
	y1, ok2 := argF32(args[1])
	z1, ok3 := argF32(args[2])
	x2, ok4 := argF32(args[3])
	y2, ok5 := argF32(args[4])
	z2, ok6 := argF32(args[5])
	if !ok1 || !ok2 || !ok3 || !ok4 || !ok5 || !ok6 {
		return value.Nil, fmt.Errorf("ENTITY.ENTITIESINBOX: numeric arguments required")
	}
	mnX := float64(x1)
	mxX := float64(x2)
	if mnX > mxX {
		mnX, mxX = mxX, mnX
	}
	mnY := float64(y1)
	mxY := float64(y2)
	if mnY > mxY {
		mnY, mxY = mxY, mnY
	}
	mnZ := float64(z1)
	mxZ := float64(z2)
	if mnZ > mxZ {
		mnZ, mxZ = mxZ, mnZ
	}
	var vals []float64
	for id, e := range m.store().ents {
		if e == nil {
			continue
		}
		wp := m.worldPos(e)
		px, py, pz := float64(wp.X), float64(wp.Y), float64(wp.Z)
		if px >= mnX && px <= mxX && py >= mnY && py <= mxY && pz >= mnZ && pz <= mxZ {
			vals = append(vals, float64(id))
		}
	}
	return m.allocFloatArray(vals)
}

func (m *Module) entClearScene(args []value.Value) (value.Value, error) {
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("ENTITY.CLEARSCENE expects 0 arguments")
	}
	st := m.store()
	ids := make([]int64, 0, len(st.ents))
	for id := range st.ents {
		ids = append(ids, id)
	}
	for _, id := range ids {
		_, _ = m.entFree([]value.Value{value.FromInt(id)})
	}
	st.groups = make(map[string]map[int64]struct{})
	st.nextID = 1
	return value.Nil, nil
}

type sceneFileV1 struct {
	V   int           `json:"v"`
	Ent []sceneEntRec `json:"e"`
}

type sceneEntRec struct {
	K int `json:"k"`

	PX float32 `json:"px,omitempty"`
	PY float32 `json:"py,omitempty"`
	PZ float32 `json:"pz,omitempty"`
	Pitch float32 `json:"pitch,omitempty"`
	Yaw   float32 `json:"yaw,omitempty"`
	Roll  float32 `json:"roll,omitempty"`
	SX float32 `json:"sx,omitempty"`
	SY float32 `json:"sy,omitempty"`
	SZ float32 `json:"sz,omitempty"`

	W float32 `json:"w,omitempty"`
	H float32 `json:"h,omitempty"`
	D float32 `json:"d,omitempty"`
	Rad     float32 `json:"rad,omitempty"`
	CylH    float32 `json:"cylh,omitempty"`
	SegH    int32   `json:"segh,omitempty"`
	SegV    int32   `json:"segv,omitempty"`

	R uint8 `json:"r,omitempty"`
	G uint8 `json:"g,omitempty"`
	B uint8 `json:"b,omitempty"`

	Static    bool    `json:"st,omitempty"`
	UseSphere bool    `json:"us,omitempty"`
	Grav      float32 `json:"grav,omitempty"`
	Path      string  `json:"path,omitempty"`
}

func (m *Module) entSaveScene(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Kind != value.KindString {
		return value.Nil, fmt.Errorf("ENTITY.SAVESCENE expects (path$)")
	}
	path, err := rt.ArgString(args, 0)
	if err != nil {
		return value.Nil, err
	}
	path = strings.TrimSpace(path)
	if path == "" {
		return value.Nil, fmt.Errorf("ENTITY.SAVESCENE: path required")
	}
	st := m.store()
	out := sceneFileV1{V: 1, Ent: make([]sceneEntRec, 0, len(st.ents))}
	for _, e := range st.ents {
		if e == nil {
			continue
		}
		out.Ent = append(out.Ent, sceneEntRec{
			K: int(e.kind),
			PX: e.pos.X, PY: e.pos.Y, PZ: e.pos.Z,
			Pitch: e.pitch, Yaw: e.yaw, Roll: e.roll,
			SX: e.scale.X, SY: e.scale.Y, SZ: e.scale.Z,
			W: e.w, H: e.h, D: e.d,
			Rad: e.radius, CylH: e.cylH, SegH: e.segH, SegV: e.segV,
			R: e.r, G: e.g, B: e.b,
			Static: e.static, UseSphere: e.useSphere, Grav: e.gravity,
			Path: e.loadPath,
		})
	}
	data, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		return value.Nil, err
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		return value.Nil, err
	}
	return value.Nil, nil
}

func (m *Module) entLoadScene(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Kind != value.KindString {
		return value.Nil, fmt.Errorf("ENTITY.LOADSCENE expects (path$)")
	}
	path, err := rt.ArgString(args, 0)
	if err != nil {
		return value.Nil, err
	}
	path = strings.TrimSpace(path)
	data, err := os.ReadFile(path)
	if err != nil {
		return value.Nil, err
	}
	var sf sceneFileV1
	if err := json.Unmarshal(data, &sf); err != nil {
		return value.Nil, err
	}
	if sf.V != 1 {
		return value.Nil, fmt.Errorf("ENTITY.LOADSCENE: unsupported scene version %d", sf.V)
	}
	if _, err := m.entClearScene([]value.Value{}); err != nil {
		return value.Nil, err
	}
	// CGo: if any entity fails to build, free all entities created in this load so Raylib/Jolt
	// resources never leak (LoadModel / LoadModelAnimations / GenMesh are not GC-managed).
	needRollback := true
	defer func() {
		if needRollback {
			_, _ = m.entClearScene([]value.Value{})
		}
	}()
	st := m.store()
	for _, r := range sf.Ent {
		id := st.nextID
		st.nextID++
		e := newDefaultEnt(id)
		e.kind = entKind(r.K)
		e.pos = rl.Vector3{X: r.PX, Y: r.PY, Z: r.PZ}
		e.pitch, e.yaw, e.roll = r.Pitch, r.Yaw, r.Roll
		e.scale = rl.Vector3{X: r.SX, Y: r.SY, Z: r.SZ}
		if e.scale.X == 0 && e.scale.Y == 0 && e.scale.Z == 0 {
			e.scale = rl.Vector3{X: 1, Y: 1, Z: 1}
		}
		e.w, e.h, e.d = r.W, r.H, r.D
		e.radius = r.Rad
		e.cylH = r.CylH
		e.segH, e.segV = r.SegH, r.SegV
		e.r, e.g, e.b = r.R, r.G, r.B
		e.static = r.Static
		e.useSphere = r.UseSphere
		e.gravity = r.Grav
		e.loadPath = r.Path
		if e.loadPath != "" {
			mod := rl.LoadModel(e.loadPath)
			if mod.MeshCount <= 0 {
				return value.Nil, fmt.Errorf("ENTITY.LOADSCENE: no meshes in model %q", e.loadPath)
			}
			e.rlModel = mod
			e.hasRLModel = true
			e.kind = entKindModel
			anims := rl.LoadModelAnimations(e.loadPath)
			if len(anims) > 0 {
				e.modelAnims = anims
				e.animLen = float32(anims[0].FrameCount)
			}
		} else if e.kind == entKindMesh {
			mesh := rl.GenMeshCube(1, 1, 1)
			mod := rl.LoadModelFromMesh(mesh)
			rl.UnloadMesh(&mesh)
			if mod.MeshCount <= 0 {
				return value.Nil, fmt.Errorf("ENTITY.LOADSCENE: internal mesh failed")
			}
			e.rlModel = mod
			e.hasRLModel = true
		}
		st.ents[id] = e
	}
	needRollback = false
	return value.Nil, nil
}

func (m *Module) camSetTargetEntity(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("CAMERA.SETTARGETENTITY expects (camera, entity#)")
	}
	ch, ok := argHandle(args[0])
	if !ok {
		return value.Nil, fmt.Errorf("CAMERA.SETTARGETENTITY: invalid camera")
	}
	eid, ok2 := args[1].ToInt()
	if !ok2 || eid < 1 {
		return value.Nil, fmt.Errorf("CAMERA.SETTARGETENTITY: invalid entity")
	}
	e := m.store().ents[eid]
	if e == nil {
		return value.Nil, fmt.Errorf("CAMERA.SETTARGETENTITY: unknown entity")
	}
	wp := m.worldPos(e)
	if err := mbcamera.ApplySetTarget(m.h, ch, wp.X, wp.Y, wp.Z); err != nil {
		return value.Nil, err
	}
	return value.Nil, nil
}
