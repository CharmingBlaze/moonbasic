//go:build cgo || (windows && !cgo)

package mbentity

import (
	"fmt"
	"math"
	"path"
	"strings"
	"unsafe"

	"moonbasic/runtime"
	"moonbasic/runtime/mbmodel3d"
	"moonbasic/runtime/texture"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func registerEntityUnifiedModelAPI(m *Module, r runtime.Registrar) {
	r.Register("ENTITY.LOADANIMATIONS", "entity", m.entLoadAnimationsExternal)
	r.Register("ENTITY.PLAY", "entity", runtime.AdaptLegacy(m.entPlay))
	r.Register("ENTITY.PLAYNAME", "entity", m.entPlayName)
	r.Register("ENTITY.STOPANIM", "entity", runtime.AdaptLegacy(m.entStopAnim))
	r.Register("ENTITY.SETANIMFRAME", "entity", runtime.AdaptLegacy(m.entSetAnimFrame))
	r.Register("ENTITY.SETANIMSPEED", "entity", runtime.AdaptLegacy(m.entSetAnimSpeedOnly))
	r.Register("ENTITY.SETANIMLOOP", "entity", runtime.AdaptLegacy(m.entSetAnimLoop))
	r.Register("ENTITY.ISPLAYING", "entity", runtime.AdaptLegacy(m.entIsPlayingAnim))
	r.Register("ENTITY.CROSSFADE", "entity", runtime.AdaptLegacy(m.entCrossfade))
	r.Register("ENTITY.TRANSITION", "entity", m.entTransition)
	r.Register("ENTITY.GETBONEPOS", "entity", runtime.AdaptLegacy(m.entGetBonePos))
	r.Register("ENTITY.GETBONEROT", "entity", runtime.AdaptLegacy(m.entGetBoneRot))
	r.Register("ENTITY.SETTEXTUREMAP", "entity", runtime.AdaptLegacy(m.entSetTextureMap))
	r.Register("MATERIAL.BULKASSIGN", "entity", runtime.AdaptLegacy(m.materialBulkAssign))
	r.Register("ENTITY.GETMETADATA", "entity", m.entGetMetadata)
	r.Register("ENTITY.SETSHADER", "entity", runtime.AdaptLegacy(m.entSetShaderModel))
	r.Register("ENTITY.GETBOUNDS", "entity", runtime.AdaptLegacy(m.entGetModelBounds))
	r.Register("ENTITY.RAYHIT", "entity", runtime.AdaptLegacy(m.entRayHit))
	r.Register("ENTITY.POINTAT", "entity", runtime.AdaptLegacy(m.entPointEntity))
	r.Register("ENTITY.ANIMNAME", "entity", m.entAnimNameAt)
	r.Register("ENTITY.ANIMNAME$", "entity", m.entAnimNameAt)
	r.Register("ENTITY.CURRENTANIM", "entity", m.entCurrentAnimName)
	r.Register("ENTITY.CURRENTANIM$", "entity", m.entCurrentAnimName)
}

func (m *Module) entLoadAnimationsExternal(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 2 || args[1].Kind != value.KindString {
		return value.Nil, fmt.Errorf("ENTITY.LOADANIMATIONS expects (entity, path)")
	}
	id, ok := m.entID(args[0])
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("invalid entity")
	}
	e := m.store().ents[id]
	if e == nil || !e.hasRLModel {
		return value.Nil, fmt.Errorf("ENTITY.LOADANIMATIONS: entity has no model")
	}
	path, err := rt.ArgString(args, 1)
	if err != nil {
		return value.Nil, err
	}
	path = strings.TrimSpace(path)
	if path == "" {
		return value.Nil, fmt.Errorf("ENTITY.LOADANIMATIONS: path required")
	}
	ext := e.getExt()
	if len(ext.modelAnims) > 0 {
		rl.UnloadModelAnimations(ext.modelAnims)
		ext.modelAnims = nil
	}
	ext.modelAnims = rl.LoadModelAnimations(path)
	ext.animIndex = 0
	ext.animTime = 0
	if len(ext.modelAnims) > 0 {
		ext.animLen = float32(ext.modelAnims[0].FrameCount)
		rl.UpdateModelAnimation(e.rlModel, ext.modelAnims[0], 0)
		rl.UpdateModelAnimationBones(e.rlModel, ext.modelAnims[0], 0)
	}
	return value.Nil, nil
}

func (m *Module) entPlay(args []value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("ENTITY.PLAY expects (entity, animIndex)")
	}
	id, ok := m.entID(args[0])
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("invalid entity")
	}
	e := m.store().ents[id]
	if e == nil {
		return value.Nil, fmt.Errorf("unknown entity")
	}
	ext := e.getExt()
	idx, ok := args[1].ToInt()
	if !ok || idx < 0 || int(idx) >= len(ext.modelAnims) {
		return value.Nil, fmt.Errorf("ENTITY.PLAY: invalid animation index")
	}
	ext.animIndex = int32(idx)
	ext.animTime = 0
	ext.animSpeed = 1
	return value.Nil, nil
}

func (m *Module) entPlayName(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 2 || args[1].Kind != value.KindString {
		return value.Nil, fmt.Errorf("ENTITY.PLAYNAME expects (entity, name)")
	}
	id, ok := m.entID(args[0])
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("invalid entity")
	}
	e := m.store().ents[id]
	if e == nil {
		return value.Nil, fmt.Errorf("unknown entity")
	}
	name, err := rt.ArgString(args, 1)
	if err != nil {
		return value.Nil, err
	}
	name = strings.TrimSpace(name)
	if name == "" {
		return value.Nil, fmt.Errorf("ENTITY.PLAYNAME: name required")
	}
	ext := e.getExt()
	for i := range ext.modelAnims {
		if strings.EqualFold(ext.modelAnims[i].GetName(), name) {
			ext.animIndex = int32(i)
			ext.animTime = 0
			ext.animSpeed = 1
			return value.Nil, nil
		}
	}
	return value.Nil, fmt.Errorf("ENTITY.PLAYNAME: no animation %q", name)
}

func (m *Module) entStopAnim(args []value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("ENTITY.STOPANIM expects entity")
	}
	id, ok := m.entID(args[0])
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("invalid entity")
	}
	e := m.store().ents[id]
	if e == nil {
		return value.Nil, fmt.Errorf("unknown entity")
	}
	e.getExt().animSpeed = 0
	return value.Nil, nil
}

func (m *Module) entSetAnimFrame(args []value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("ENTITY.SETANIMFRAME expects (entity, frame)")
	}
	id, ok := m.entID(args[0])
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("invalid entity")
	}
	e := m.store().ents[id]
	if e == nil {
		return value.Nil, fmt.Errorf("unknown entity")
	}
	fr, ok := argF32(args[1])
	if !ok {
		return value.Nil, fmt.Errorf("frame must be numeric")
	}
	ext := e.getExt()
	if len(ext.modelAnims) == 0 {
		return value.Nil, nil
	}
	ai := ext.animIndex
	if ai < 0 || int(ai) >= len(ext.modelAnims) {
		ai = 0
	}
	fc := ext.modelAnims[ai].FrameCount
	if fc <= 0 {
		return value.Nil, nil
	}
	if fr < 0 {
		fr = 0
	}
	if fr >= float32(fc) {
		fr = float32(fc) - 1
	}
	ext.animTime = fr
	return value.Nil, nil
}

func (m *Module) entSetAnimSpeedOnly(args []value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("ENTITY.SETANIMSPEED expects (entity, speed)")
	}
	id, ok := m.entID(args[0])
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("invalid entity")
	}
	e := m.store().ents[id]
	if e == nil {
		return value.Nil, fmt.Errorf("unknown entity")
	}
	s, ok := argF32(args[1])
	if !ok {
		return value.Nil, fmt.Errorf("speed must be numeric")
	}
	e.getExt().animSpeed = s
	return value.Nil, nil
}

func (m *Module) entSetAnimLoop(args []value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("ENTITY.SETANIMLOOP expects (entity, loop)")
	}
	id, ok := m.entID(args[0])
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("invalid entity")
	}
	e := m.store().ents[id]
	if e == nil {
		return value.Nil, fmt.Errorf("unknown entity")
	}
	loop, ok := argBool(args[1])
	if !ok {
		return value.Nil, fmt.Errorf("loop must be bool or numeric")
	}
	ext := e.getExt()
	if loop {
		ext.animMode = 0
	} else {
		ext.animMode = 3
	}
	return value.Nil, nil
}

func argBool(v value.Value) (bool, bool) {
	if v.Kind == value.KindBool {
		return v.IVal != 0, true
	}
	if x, ok := v.ToInt(); ok {
		return x != 0, true
	}
	return false, false
}

func (m *Module) entIsPlayingAnim(args []value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("ENTITY.ISPLAYING expects entity")
	}
	id, ok := m.entID(args[0])
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("invalid entity")
	}
	e := m.store().ents[id]
	if e == nil {
		return value.Nil, fmt.Errorf("unknown entity")
	}
	ext := e.getExt()
	on := len(ext.modelAnims) > 0 && ext.animSpeed != 0
	return value.FromBool(on), nil
}

func (m *Module) entCrossfade(args []value.Value) (value.Value, error) {
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("ENTITY.CROSSFADE expects (entity, nextIndex, duration)")
	}
	id, ok := m.entID(args[0])
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("invalid entity")
	}
	e := m.store().ents[id]
	if e == nil {
		return value.Nil, fmt.Errorf("unknown entity")
	}
	ext := e.getExt()
	idx, ok := args[1].ToInt()
	if !ok || idx < 0 || int(idx) >= len(ext.modelAnims) {
		return value.Nil, fmt.Errorf("ENTITY.CROSSFADE: invalid animation index")
	}
	ext.animIndex = int32(idx)
	ext.animTime = 0
	ext.animSpeed = 1
	return value.Nil, nil
}

func (m *Module) entTransition(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 3 || args[1].Kind != value.KindString {
		return value.Nil, fmt.Errorf("ENTITY.TRANSITION expects (entity, name, duration): duration reserved; same as PLAYNAME")
	}
	_, err := rt.ArgString(args, 1)
	if err != nil {
		return value.Nil, err
	}
	return m.entPlayName(rt, args[0], args[1])
}

func (m *Module) boneWorldMatrix(host *ent, boneName string) (rl.Matrix, error) {
	if host == nil || !host.hasRLModel {
		return rl.MatrixIdentity(), fmt.Errorf("no model")
	}
	name := strings.TrimSpace(boneName)
	if name == "" {
		return rl.MatrixIdentity(), fmt.Errorf("bone name required")
	}
	bones := host.rlModel.GetBones()
	var bi int32 = -1
	for i := range bones {
		if strings.EqualFold(boneNameStr(bones[i].Name), name) {
			bi = int32(i)
			break
		}
	}
	if bi < 0 {
		return rl.MatrixIdentity(), fmt.Errorf("no bone %q", boneName)
	}
	meshes := host.rlModel.GetMeshes()
	if len(meshes) == 0 {
		return rl.MatrixIdentity(), fmt.Errorf("no mesh")
	}
	mesh := meshes[0]
	if mesh.BoneMatrices == nil || int(mesh.BoneCount) <= int(bi) {
		return rl.MatrixIdentity(), fmt.Errorf("bone matrices not ready")
	}
	bm := unsafe.Slice(mesh.BoneMatrices, mesh.BoneCount)[bi]
	hw := m.worldMatrix(host)
	return rl.MatrixMultiply(hw, bm), nil
}

func (m *Module) entGetBonePos(args []value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("ENTITY.GETBONEPOS: heap not bound")
	}
	if len(args) != 2 || args[1].Kind != value.KindString {
		return value.Nil, fmt.Errorf("ENTITY.GETBONEPOS expects (entity, boneName)")
	}
	id, ok := m.entID(args[0])
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("invalid entity")
	}
	e := m.store().ents[id]
	if e == nil {
		return value.Nil, fmt.Errorf("unknown entity")
	}
	if args[1].Kind != value.KindString {
		return value.Nil, fmt.Errorf("bone name must be string")
	}
	name, ok := m.h.GetString(int32(args[1].IVal))
	if !ok {
		return value.Nil, fmt.Errorf("invalid bone name string")
	}
	bw, err := m.boneWorldMatrix(e, name)
	if err != nil {
		return value.Nil, fmt.Errorf("ENTITY.GETBONEPOS: %w", err)
	}
	x, y, z := bw.M12, bw.M13, bw.M14
	arr, err := heap.NewArrayOfKind([]int64{3}, heap.ArrayKindFloat, 0)
	if err != nil {
		return value.Nil, err
	}
	arr.Floats[0] = float64(x)
	arr.Floats[1] = float64(y)
	arr.Floats[2] = float64(z)
	h, err := m.h.Alloc(arr)
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(h), nil
}

func (m *Module) entGetBoneRot(args []value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("ENTITY.GETBONEROT: heap not bound")
	}
	if len(args) != 2 || args[1].Kind != value.KindString {
		return value.Nil, fmt.Errorf("ENTITY.GETBONEROT expects (entity, boneName)")
	}
	id, ok := m.entID(args[0])
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("invalid entity")
	}
	e := m.store().ents[id]
	if e == nil {
		return value.Nil, fmt.Errorf("unknown entity")
	}
	if args[1].Kind != value.KindString {
		return value.Nil, fmt.Errorf("bone name must be string")
	}
	name, ok := m.h.GetString(int32(args[1].IVal))
	if !ok {
		return value.Nil, fmt.Errorf("invalid bone name string")
	}
	bw, err := m.boneWorldMatrix(e, name)
	if err != nil {
		return value.Nil, fmt.Errorf("ENTITY.GETBONEROT: %w", err)
	}
	q := rl.QuaternionFromMatrix(bw)
	v := rl.QuaternionToEuler(q)
	// Match ENTITY.ENTITYPITCH/YAW/ROLL convention (entity_xform worldEuler)
	pitch, yaw, roll := v.Y, v.Z, v.X
	arr, err := heap.NewArrayOfKind([]int64{3}, heap.ArrayKindFloat, 0)
	if err != nil {
		return value.Nil, err
	}
	arr.Floats[0] = float64(pitch)
	arr.Floats[1] = float64(yaw)
	arr.Floats[2] = float64(roll)
	h, err := m.h.Alloc(arr)
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(h), nil
}

func (m *Module) entSetTextureMap(args []value.Value) (value.Value, error) {
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("ENTITY.SETTEXTUREMAP expects (entity, materialIndex, textureHandle)")
	}
	id, ok := m.entID(args[0])
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("invalid entity")
	}
	e := m.store().ents[id]
	if e == nil || !e.hasRLModel {
		return value.Nil, fmt.Errorf("ENTITY.SETTEXTUREMAP: entity has no model")
	}
	mi, ok := args[1].ToInt()
	if !ok || mi < 0 {
		return value.Nil, fmt.Errorf("invalid material index")
	}
	texH, ok := argHandle(args[2])
	if !ok {
		return value.Nil, fmt.Errorf("texture handle required")
	}
	tex, err := texture.ForBinding(m.h, texH)
	if err != nil {
		return value.Nil, err
	}
	mats := e.rlModel.GetMaterials()
	if int(mi) >= len(mats) {
		return value.Nil, fmt.Errorf("ENTITY.SETTEXTUREMAP: material index out of range")
	}
	rl.SetMaterialTexture(&mats[mi], rl.MapAlbedo, tex)
	return value.Nil, nil
}

func matchMaterialGlob(pat, s string) bool {
	pat = strings.TrimSpace(pat)
	s = strings.TrimSpace(s)
	if pat == "" || s == "" {
		return false
	}
	ok, _ := path.Match(strings.ToUpper(pat), strings.ToUpper(s))
	return ok
}

func (m *Module) materialBulkAssign(args []value.Value) (value.Value, error) {
	if len(args) != 2 && len(args) != 3 {
		return value.Nil, fmt.Errorf("MATERIAL.BULKASSIGN expects (pattern, textureHandle [, materialIndex])")
	}
	if args[0].Kind != value.KindString {
		return value.Nil, fmt.Errorf("pattern must be string")
	}
	pat, ok := m.h.GetString(int32(args[0].IVal))
	if !ok {
		return value.Nil, fmt.Errorf("invalid pattern string")
	}
	texH, ok := argHandle(args[1])
	if !ok {
		return value.Nil, fmt.Errorf("texture handle required")
	}
	mi := int64(0)
	if len(args) == 3 {
		var ok2 bool
		mi, ok2 = args[2].ToInt()
		if !ok2 || mi < 0 {
			return value.Nil, fmt.Errorf("invalid material index")
		}
	}
	tex, err := texture.ForBinding(m.h, texH)
	if err != nil {
		return value.Nil, err
	}
	st := m.store()
	var n int64
	for _, e := range st.ents {
		if e == nil || !e.hasRLModel {
			continue
		}
		ext := e.getExt()
		if !matchMaterialGlob(pat, ext.name) && !matchMaterialGlob(pat, ext.blenderTag) {
			continue
		}
		mats := e.rlModel.GetMaterials()
		if int(mi) >= len(mats) {
			continue
		}
		rl.SetMaterialTexture(&mats[mi], rl.MapAlbedo, tex)
		n++
	}
	return value.FromInt(n), nil
}

func (m *Module) entGetMetadata(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 2 || args[1].Kind != value.KindString {
		return value.Nil, fmt.Errorf("ENTITY.GETMETADATA expects (entity, key)")
	}
	id, ok := m.entID(args[0])
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("invalid entity")
	}
	key, err := rt.ArgString(args, 1)
	if err != nil {
		return value.Nil, err
	}
	key = strings.TrimSpace(key)
	st := m.store()
	if st.entMeta == nil {
		return rt.RetString(""), nil
	}
	row := st.entMeta[id]
	if row == nil {
		return rt.RetString(""), nil
	}
	if v, ok := row[key]; ok {
		return rt.RetString(v), nil
	}
	for k, v := range row {
		if strings.EqualFold(k, key) {
			return rt.RetString(v), nil
		}
	}
	return rt.RetString(""), nil
}

func (m *Module) entSetShaderModel(args []value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("ENTITY.SETSHADER: heap not bound")
	}
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("ENTITY.SETSHADER expects (entity, shaderHandle)")
	}
	id, ok := m.entID(args[0])
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("invalid entity")
	}
	e := m.store().ents[id]
	if e == nil || !e.hasRLModel {
		return value.Nil, fmt.Errorf("ENTITY.SETSHADER: entity has no model")
	}
	shH, ok := argHandle(args[1])
	if !ok {
		return value.Nil, fmt.Errorf("shader handle required")
	}
	sh, err := mbmodel3d.ShaderRaylib(m.h, shH)
	if err != nil {
		return value.Nil, err
	}
	mats := e.rlModel.GetMaterials()
	for i := range mats {
		mats[i].Shader = sh
	}
	return value.Nil, nil
}

func aabbWorldFromModel(bb rl.BoundingBox, mat rl.Matrix) (mn, mx rl.Vector3) {
	corners := [8]rl.Vector3{
		{X: bb.Min.X, Y: bb.Min.Y, Z: bb.Min.Z},
		{X: bb.Max.X, Y: bb.Min.Y, Z: bb.Min.Z},
		{X: bb.Min.X, Y: bb.Max.Y, Z: bb.Min.Z},
		{X: bb.Max.X, Y: bb.Max.Y, Z: bb.Min.Z},
		{X: bb.Min.X, Y: bb.Min.Y, Z: bb.Max.Z},
		{X: bb.Max.X, Y: bb.Min.Y, Z: bb.Max.Z},
		{X: bb.Min.X, Y: bb.Max.Y, Z: bb.Max.Z},
		{X: bb.Max.X, Y: bb.Max.Y, Z: bb.Max.Z},
	}
	p0 := rl.Vector3Transform(corners[0], mat)
	mn = p0
	mx = p0
	for i := 1; i < 8; i++ {
		p := rl.Vector3Transform(corners[i], mat)
		mn.X = float32(math.Min(float64(mn.X), float64(p.X)))
		mn.Y = float32(math.Min(float64(mn.Y), float64(p.Y)))
		mn.Z = float32(math.Min(float64(mn.Z), float64(p.Z)))
		mx.X = float32(math.Max(float64(mx.X), float64(p.X)))
		mx.Y = float32(math.Max(float64(mx.Y), float64(p.Y)))
		mx.Z = float32(math.Max(float64(mx.Z), float64(p.Z)))
	}
	return mn, mx
}

func entityRayMeshCollision(ray rl.Ray, model rl.Model) rl.RayCollision {
	best := rl.RayCollision{Hit: false, Distance: float32(math.MaxFloat32)}
	tform := model.Transform
	meshes := model.GetMeshes()
	for i := range meshes {
		col := rl.GetRayCollisionMesh(ray, meshes[i], tform)
		if col.Hit && col.Distance < best.Distance {
			best = col
		}
	}
	if !best.Hit {
		return rl.RayCollision{Hit: false}
	}
	return best
}

func (m *Module) entGetModelBounds(args []value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("ENTITY.GETBOUNDS: heap not bound")
	}
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("ENTITY.GETBOUNDS expects entity")
	}
	id, ok := m.entID(args[0])
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("invalid entity")
	}
	e := m.store().ents[id]
	if e == nil || !e.hasRLModel {
		return value.Nil, fmt.Errorf("ENTITY.GETBOUNDS: entity has no model")
	}
	bb := rl.GetModelBoundingBox(e.rlModel)
	wm := m.worldMatrix(e)
	mn, mx := aabbWorldFromModel(bb, wm)
	arr, err := heap.NewArrayOfKind([]int64{6}, heap.ArrayKindFloat, 0)
	if err != nil {
		return value.Nil, err
	}
	arr.Floats[0] = float64(mn.X)
	arr.Floats[1] = float64(mn.Y)
	arr.Floats[2] = float64(mn.Z)
	arr.Floats[3] = float64(mx.X)
	arr.Floats[4] = float64(mx.Y)
	arr.Floats[5] = float64(mx.Z)
	h, err := m.h.Alloc(arr)
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(h), nil
}

func (m *Module) entRayHit(args []value.Value) (value.Value, error) {
	if len(args) != 7 {
		return value.Nil, fmt.Errorf("ENTITY.RAYHIT expects (entity, ox, oy, oz, dx, dy, dz)")
	}
	id, ok := m.entID(args[0])
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("invalid entity")
	}
	e := m.store().ents[id]
	if e == nil || !e.hasRLModel {
		return value.FromBool(false), nil
	}
	ox, ok1 := argF32(args[1])
	oy, ok2 := argF32(args[2])
	oz, ok3 := argF32(args[3])
	dx, ok4 := argF32(args[4])
	dy, ok5 := argF32(args[5])
	dz, ok6 := argF32(args[6])
	if !ok1 || !ok2 || !ok3 || !ok4 || !ok5 || !ok6 {
		return value.Nil, fmt.Errorf("ENTITY.RAYHIT: numeric ray components required")
	}
	dir := rl.Vector3{X: dx, Y: dy, Z: dz}
	if rl.Vector3Length(dir) < 1e-8 {
		return value.FromBool(false), nil
	}
	dir = rl.Vector3Normalize(dir)
	ray := rl.Ray{Position: rl.Vector3{X: ox, Y: oy, Z: oz}, Direction: dir}
	wm := m.worldMatrix(e)
	saved := e.rlModel.Transform
	e.rlModel.Transform = wm
	col := entityRayMeshCollision(ray, e.rlModel)
	e.rlModel.Transform = saved
	return value.FromBool(col.Hit), nil
}

func (m *Module) entAnimNameAt(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("ENTITY.ANIMNAME expects (entity, idx)")
	}
	id, ok := m.entID(args[0])
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("invalid entity")
	}
	e := m.store().ents[id]
	if e == nil {
		return value.Nil, fmt.Errorf("unknown entity")
	}
	ext := e.getExt()
	idx, ok := args[1].ToInt()
	if !ok || idx < 0 || int(idx) >= len(ext.modelAnims) {
		return value.Nil, fmt.Errorf("ENTITY.ANIMNAME: invalid index")
	}
	return rt.RetString(ext.modelAnims[idx].GetName()), nil
}

func (m *Module) entCurrentAnimName(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("ENTITY.CURRENTANIM expects entity")
	}
	id, ok := m.entID(args[0])
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("invalid entity")
	}
	e := m.store().ents[id]
	if e == nil {
		return value.Nil, fmt.Errorf("unknown entity")
	}
	ext := e.getExt()
	if len(ext.modelAnims) == 0 {
		return rt.RetString(""), nil
	}
	ai := ext.animIndex
	if ai < 0 || int(ai) >= len(ext.modelAnims) {
		ai = 0
	}
	return rt.RetString(ext.modelAnims[ai].GetName()), nil
}
