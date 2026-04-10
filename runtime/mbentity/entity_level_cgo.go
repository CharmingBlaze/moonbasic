//go:build cgo || (windows && !cgo)

package mbentity

import (
	"fmt"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/qmuntal/gltf"

	"moonbasic/runtime"
	"moonbasic/runtime/mbmatrix"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func registerLevelGLTFAPI(m *Module, r runtime.Registrar) {
	r.Register("LEVEL.SETROOT", "entity", m.levelSetRoot)
	r.Register("LEVEL.LOAD", "entity", m.levelLoadGLTF)
	r.Register("LEVEL.FINDENTITY", "entity", m.levelFindEntity)
	r.Register("LEVEL.GETMARKER", "entity", runtime.AdaptLegacy(m.levelGetMarker))
	r.Register("LEVEL.GETSPAWN", "entity", m.levelGetSpawn)
	r.Register("LEVEL.SHOWLAYER", "entity", m.levelShowLayer)
	r.Register("LEVEL.APPLYPHYSICS", "entity", runtime.AdaptLegacy(m.levelApplyPhysics))
	r.Register("LEVEL.SYNCLIGHTS", "entity", runtime.AdaptLegacy(m.levelSyncLights))
	r.Register("PHYSICS.AUTOCREATE", "entity", runtime.AdaptLegacy(m.physicsAutoCreate))
	r.Register("ENTITY.SETSTATIC", "entity", runtime.AdaptLegacy(m.entSetStatic))
	r.Register("ENTITY.SETTRIGGER", "entity", runtime.AdaptLegacy(m.entSetTrigger))
	r.Register("ENTITY.INSTANCE", "entity", runtime.AdaptLegacy(m.entInstanceStub))
	r.Register("LEVEL.BINDSCRIPT", "entity", m.levelBindScript)
	r.Register("LEVEL.MATCHSCRIPTBIND", "entity", m.levelMatchScriptBind)
	r.Register("LEVEL.LOADSKYBOX", "entity", m.levelLoadSkybox)
	r.Register("LEVEL.OPTIMIZE", "entity", runtime.AdaptLegacy(m.levelOptimizeStub))
	r.Register("TRIGGER.CREATEFROMENTITY", "entity", runtime.AdaptLegacy(m.triggerCreateFromEntityStub))
}

func (m *Module) clearLevelState() {
	st := m.store()
	st.levelRoot = ""
	st.levelMarkers = nil
	st.levelSpawn = nil
	st.levelLayers = nil
	st.levelColliders = nil
}

func (m *Module) levelSetRoot(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Kind != value.KindString {
		return value.Nil, fmt.Errorf("LEVEL.SETROOT expects (path$)")
	}
	p, err := rt.ArgString(args, 0)
	if err != nil {
		return value.Nil, err
	}
	p = strings.TrimSpace(p)
	m.store().levelRoot = p
	return value.Nil, nil
}

func resolveLevelPath(root, file string) string {
	file = filepath.Clean(file)
	if filepath.IsAbs(file) {
		return file
	}
	if root == "" {
		return file
	}
	return filepath.Join(filepath.Clean(root), file)
}

func gltfMatToRL(m [16]float64) rl.Matrix {
	return rl.NewMatrix(
		float32(m[0]), float32(m[4]), float32(m[8]), float32(m[12]),
		float32(m[1]), float32(m[5]), float32(m[9]), float32(m[13]),
		float32(m[2]), float32(m[6]), float32(m[10]), float32(m[14]),
		float32(m[3]), float32(m[7]), float32(m[11]), float32(m[15]),
	)
}

func matWorldPos(wm rl.Matrix) rl.Vector3 {
	return rl.Vector3{X: wm.M12, Y: wm.M13, Z: wm.M14}
}

func sceneRootNodes(doc *gltf.Document) []int {
	var si int
	if doc.Scene != nil {
		si = *doc.Scene
	} else if len(doc.Scenes) > 0 {
		si = 0
	} else {
		return nil
	}
	if si < 0 || si >= len(doc.Scenes) || doc.Scenes[si] == nil {
		return nil
	}
	return doc.Scenes[si].Nodes
}

func computeNodeWorlds(doc *gltf.Document) []rl.Matrix {
	n := len(doc.Nodes)
	if n == 0 {
		return nil
	}
	local := make([]rl.Matrix, n)
	for i := range doc.Nodes {
		if doc.Nodes[i] == nil {
			local[i] = rl.MatrixIdentity()
			continue
		}
		local[i] = gltfMatToRL(doc.Nodes[i].MatrixOrDefault())
	}
	worlds := make([]rl.Matrix, n)
	visited := make([]bool, n)
	var walk func(int, rl.Matrix)
	walk = func(idx int, parent rl.Matrix) {
		if idx < 0 || idx >= n || doc.Nodes[idx] == nil || visited[idx] {
			return
		}
		visited[idx] = true
		worlds[idx] = rl.MatrixMultiply(parent, local[idx])
		for _, c := range doc.Nodes[idx].Children {
			walk(c, worlds[idx])
		}
	}
	ident := rl.MatrixIdentity()
	for _, r := range sceneRootNodes(doc) {
		walk(r, ident)
	}
	for i := range visited {
		if !visited[i] && doc.Nodes[i] != nil {
			worlds[i] = local[i]
		}
	}
	return worlds
}

func (m *Module) levelLoadGLTF(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Kind != value.KindString {
		return value.Nil, fmt.Errorf("LEVEL.LOAD expects (path$)")
	}
	path, err := rt.ArgString(args, 0)
	if err != nil {
		return value.Nil, err
	}
	path = strings.TrimSpace(path)
	if path == "" {
		return value.Nil, fmt.Errorf("LEVEL.LOAD: path required")
	}
	abs := resolveLevelPath(m.store().levelRoot, path)
	doc, err := gltf.Open(abs)
	if err != nil {
		return value.Nil, fmt.Errorf("LEVEL.LOAD: %w", err)
	}
	st := m.store()
	st.levelMarkers = make(map[string]rl.Vector3)
	st.levelSpawn = make(map[string]rl.Matrix)
	st.levelLayers = make(map[string][]int64)
	st.levelColliders = nil

	worlds := computeNodeWorlds(doc)
	if len(worlds) == 0 {
		return value.Nil, fmt.Errorf("LEVEL.LOAD: no nodes")
	}

	var firstVisual int = -1
	for i, node := range doc.Nodes {
		if node == nil || node.Mesh == nil {
			continue
		}
		name := strings.TrimSpace(node.Name)
		if name != "" && strings.HasPrefix(strings.ToUpper(name), "COL_") {
			continue
		}
		firstVisual = i
		break
	}
	if firstVisual < 0 {
		for i, node := range doc.Nodes {
			if node != nil && node.Mesh != nil {
				firstVisual = i
				break
			}
		}
	}

	for i, node := range doc.Nodes {
		if node == nil {
			continue
		}
		name := strings.TrimSpace(node.Name)
		if name == "" || i >= len(worlds) {
			continue
		}
		key := groupKey(name)
		wm := worlds[i]
		st.levelSpawn[key] = wm
		st.levelMarkers[key] = matWorldPos(wm)
		if node.Mesh == nil {
			continue
		}
		if strings.HasPrefix(strings.ToUpper(name), "COL_") {
			st.levelColliders = append(st.levelColliders, levelColliderRec{Name: name, World: wm})
		}
	}

	idx := m.h.Intern(abs)
	v, err := m.entLoadMesh(rt, value.Value{Kind: value.KindString, IVal: int64(idx)})
	if err != nil {
		return value.Nil, err
	}
	eid, ok := v.ToInt()
	if !ok || eid < 1 {
		return v, nil
	}
	e := st.ents[eid]
	if e == nil {
		return value.Nil, fmt.Errorf("LEVEL.LOAD: internal entity")
	}
	if firstVisual >= 0 && firstVisual < len(worlds) {
		wm := worlds[firstVisual]
		var t rl.Vector3
		var q rl.Quaternion
		var s rl.Vector3
		wm.Decompose(&t, &q, &s)
		eu := rl.QuaternionToEuler(q)
		e.setPos(t)
		e.setRot(eu.Y, eu.Z, eu.X)
		if s.X != 0 || s.Y != 0 || s.Z != 0 {
			e.scale = s
		}
		n := doc.Nodes[firstVisual]
		if n != nil && strings.TrimSpace(n.Name) != "" {
			nm := strings.TrimSpace(n.Name)
			_, _ = m.entSetName(rt, value.FromInt(eid), value.Value{Kind: value.KindString, IVal: int64(m.h.Intern(nm))})
		}
	}
	if firstVisual >= 0 {
		n := doc.Nodes[firstVisual]
		if n != nil {
			name := strings.TrimSpace(n.Name)
			if name != "" && strings.HasPrefix(strings.ToUpper(name), "COL_") {
				e.hidden = true
			}
		}
	}

	if firstVisual >= 0 && firstVisual < len(doc.Nodes) && doc.Nodes[firstVisual] != nil {
		n := doc.Nodes[firstVisual]
		layer := extrasLayer(n.Extras)
		if layer != "" {
			lk := groupKey(layer)
			st.levelLayers[lk] = append(st.levelLayers[lk], eid)
		}
		if n.Extras != nil {
			meta := make(map[string]string)
			flattenExtras(n.Extras, "", meta)
			if len(meta) > 0 {
				if st.entMeta == nil {
					st.entMeta = make(map[int64]map[string]string)
				}
				st.entMeta[eid] = meta
			}
			if t, ok := meta["tag"]; ok {
				e.getExt().blenderTag = strings.TrimSpace(t)
			}
		}
	}

	return value.FromInt(eid), nil
}

func flattenExtras(ex any, prefix string, out map[string]string) {
	m, ok := ex.(map[string]any)
	if !ok {
		return
	}
	for k, v := range m {
		full := k
		if prefix != "" {
			full = prefix + "." + k
		}
		switch vv := v.(type) {
		case string:
			out[full] = vv
		case float64:
			out[full] = strconv.FormatFloat(vv, 'g', -1, 64)
		case bool:
			out[full] = strconv.FormatBool(vv)
		case map[string]any:
			flattenExtras(vv, full, out)
		default:
			out[full] = fmt.Sprintf("%v", vv)
		}
	}
}

func (m *Module) levelLoadSkybox(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Kind != value.KindString {
		return value.Nil, fmt.Errorf("LEVEL.LOADSKYBOX expects (hdrPath$)")
	}
	if m.tex == nil {
		return value.Nil, fmt.Errorf("LEVEL.LOADSKYBOX: texture module not wired (internal error)")
	}
	path, err := rt.ArgString(args, 0)
	if err != nil {
		return value.Nil, err
	}
	path = strings.TrimSpace(path)
	if path == "" {
		return value.Nil, fmt.Errorf("LEVEL.LOADSKYBOX: path required")
	}
	abs := resolveLevelPath(m.store().levelRoot, path)
	h, err := m.tex.TexLoadPath(abs, 1)
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(h), nil
}

func (m *Module) levelBindScript(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 2 || args[0].Kind != value.KindString || args[1].Kind != value.KindString {
		return value.Nil, fmt.Errorf("LEVEL.BINDSCRIPT expects (entityNamePattern$, functionName$)")
	}
	pat, err := rt.ArgString(args, 0)
	if err != nil {
		return value.Nil, err
	}
	fnName, err := rt.ArgString(args, 1)
	if err != nil {
		return value.Nil, err
	}
	pat = strings.TrimSpace(pat)
	fnName = strings.TrimSpace(fnName)
	if pat == "" || fnName == "" {
		return value.Nil, fmt.Errorf("LEVEL.BINDSCRIPT: pattern and function name must be non-empty")
	}
	st := m.store()
	st.scriptBinds = append(st.scriptBinds, scriptBindRec{pattern: pat, fnName: fnName})
	return value.Nil, nil
}

func (m *Module) levelMatchScriptBind(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Kind != value.KindString {
		return value.Nil, fmt.Errorf("LEVEL.MATCHSCRIPTBIND expects (objectName$)")
	}
	name, err := rt.ArgString(args, 0)
	if err != nil {
		return value.Nil, err
	}
	name = strings.TrimSpace(name)
	st := m.store()
	for _, b := range st.scriptBinds {
		if matchLevelGlob(b.pattern, name) {
			return rt.RetString(b.fnName), nil
		}
	}
	return rt.RetString(""), nil
}

func matchLevelGlob(pat, s string) bool {
	pat = strings.TrimSpace(pat)
	s = strings.TrimSpace(s)
	if pat == "" || s == "" {
		return false
	}
	ok, _ := path.Match(strings.ToUpper(pat), strings.ToUpper(s))
	return ok
}

func (m *Module) levelOptimizeStub(args []value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("LEVEL.OPTIMIZE expects (entityOrScene#)")
	}
	_, _ = m.entID(args[0])
	return value.Nil, fmt.Errorf("LEVEL.OPTIMIZE: static mesh merging / batching not implemented yet — use MODEL.MAKEINSTANCED for draw batches")
}

func (m *Module) triggerCreateFromEntityStub(args []value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("TRIGGER.CREATEFROMENTITY expects (entity#)")
	}
	_, _ = m.entID(args[0])
	return value.Nil, fmt.Errorf("TRIGGER.CREATEFROMENTITY: Jolt sensor from mesh not exposed yet — use ENTITY.SETTRIGGER when available")
}

func extrasLayer(extras any) string {
	if extras == nil {
		return ""
	}
	switch m := extras.(type) {
	case map[string]any:
		if v, ok := m["layer"]; ok {
			switch s := v.(type) {
			case string:
				return strings.TrimSpace(s)
			case float64:
				return fmt.Sprintf("%v", s)
			}
		}
	}
	return ""
}

func (m *Module) levelFindEntity(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	return m.entFind(rt, args...)
}

func (m *Module) levelGetMarker(args []value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("LEVEL.GETMARKER: heap not bound")
	}
	if len(args) != 1 || args[0].Kind != value.KindString {
		return value.Nil, fmt.Errorf("LEVEL.GETMARKER expects (name$)")
	}
	name, ok := m.h.GetString(int32(args[0].IVal))
	if !ok {
		return value.Nil, fmt.Errorf("invalid name")
	}
	k := groupKey(name)
	p, ok := m.store().levelMarkers[k]
	if !ok {
		return value.Nil, fmt.Errorf("LEVEL.GETMARKER: unknown marker %q", name)
	}
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

func (m *Module) levelGetSpawn(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("LEVEL.GETSPAWN: heap not bound")
	}
	if len(args) != 1 || args[0].Kind != value.KindString {
		return value.Nil, fmt.Errorf("LEVEL.GETSPAWN expects (name$)")
	}
	name, err := rt.ArgString(args, 0)
	if err != nil {
		return value.Nil, err
	}
	k := groupKey(name)
	wm, ok := m.store().levelSpawn[k]
	if !ok {
		return value.Nil, fmt.Errorf("LEVEL.GETSPAWN: unknown spawn %q", name)
	}
	h, err := mbmatrix.AllocMatrix(m.h, wm)
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(h), nil
}

func (m *Module) levelShowLayer(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("LEVEL.SHOWLAYER expects (layerName$, visible?)")
	}
	if args[0].Kind != value.KindString {
		return value.Nil, fmt.Errorf("layer name must be string")
	}
	lname, err := rt.ArgString(args, 0)
	if err != nil {
		return value.Nil, err
	}
	en, ok := argBool(args[1])
	if !ok {
		return value.Nil, fmt.Errorf("visible must be bool or numeric")
	}
	k := groupKey(lname)
	st := m.store()
	ids := st.levelLayers[k]
	if len(ids) == 0 && st.groups != nil && st.groups[k] != nil {
		for id := range st.groups[k] {
			ids = append(ids, id)
		}
	}
	for _, id := range ids {
		e := st.ents[id]
		if e == nil {
			continue
		}
		e.hidden = !en
	}
	return value.Nil, nil
}

func (m *Module) levelApplyPhysics(args []value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("LEVEL.APPLYPHYSICS expects (entity#)")
	}
	_, _ = m.entID(args[0])
	return value.Nil, fmt.Errorf("LEVEL.APPLYPHYSICS: not implemented (map Blender extras to BODY3D.* + COMMIT manually; see PHYSICS3D.md)")
}

func (m *Module) levelSyncLights(args []value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("LEVEL.SYNCLIGHTS expects (toggle?)")
	}
	_, _ = argBool(args[0])
	return value.Nil, fmt.Errorf("LEVEL.SYNCLIGHTS: KHR_lights_punctual → LIGHT.* not wired yet")
}

func (m *Module) physicsAutoCreate(args []value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("PHYSICS.AUTOCREATE expects (entity#)")
	}
	_, _ = m.entID(args[0])
	return value.Nil, fmt.Errorf("PHYSICS.AUTOCREATE: use BODY3D.ADDBOX/ADDMESH from ENTITY.GETBOUNDS; automation not wired yet")
}

func (m *Module) entSetStatic(args []value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("ENTITY.SETSTATIC expects (entity#)")
	}
	_, _ = m.entID(args[0])
	return value.Nil, fmt.Errorf("ENTITY.SETSTATIC: use BODY3D.MAKE(\"STATIC\") + mesh shapes; entity motion flags not exposed here yet")
}

func (m *Module) entSetTrigger(args []value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("ENTITY.SETTRIGGER expects (entity#)")
	}
	_, _ = m.entID(args[0])
	return value.Nil, fmt.Errorf("ENTITY.SETTRIGGER: Jolt sensor shapes not exposed in jolt-go binding yet")
}

func (m *Module) entInstanceStub(args []value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("ENTITY.INSTANCE expects (entity#)")
	}
	_, _ = m.entID(args[0])
	return value.Nil, fmt.Errorf("ENTITY.INSTANCE: not implemented — use MODEL.MAKEINSTANCED for GPU instancing; ENTITY.COPY duplicates VRAM today")
}
