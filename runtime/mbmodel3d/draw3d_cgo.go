//go:build cgo || (windows && !cgo)

package mbmodel3d

import (
	"sync"
	"unsafe"

	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/runtime/mblight"
	"moonbasic/vm/heap"
)

var (
	draw3dMu       sync.Mutex
	inCamera3D     bool
	activeCamPos   rl.Vector3
	activeCamera   rl.Camera3D
	activeCameraOK bool
	shadowMapSize  int32 = 2048
	shadowRT       rl.RenderTexture2D
	shadowRTReady  bool
	depthDrawMat   rl.Material
	depthMatInited bool
	lightSpaceVP   rl.Matrix

	deferredMeshes  []deferredMeshRec
	deferredModels  []heap.Handle
	pendingInstDraw []instancedDrawRec
)

type deferredMeshRec struct {
	meshH, matH heap.Handle
	mtx         rl.Matrix
}

type instancedDrawRec struct {
	instH    heap.Handle
	lodMeshH heap.Handle // optional; 0 = use instanced model mesh only
	lodDist  float32     // distance threshold for LOD mesh (INSTANCE.DRAWLOD / deferred)
}

func MarkCamera3DBegin(camX, camY, camZ float32) {
	draw3dMu.Lock()
	defer draw3dMu.Unlock()
	inCamera3D = true
	activeCamPos = rl.Vector3{X: camX, Y: camY, Z: camZ}
}

// StoreActiveCamera3D saves the full camera for billboards (particles). CGO camera.BEGIN calls this.
func StoreActiveCamera3D(cam rl.Camera3D) {
	draw3dMu.Lock()
	defer draw3dMu.Unlock()
	activeCamera = cam
	activeCameraOK = true
}

// ActiveCamera3D returns the camera from the current CAMERA.BEGIN block, if any.
func ActiveCamera3D() (rl.Camera3D, bool) {
	draw3dMu.Lock()
	defer draw3dMu.Unlock()
	return activeCamera, activeCameraOK
}

// ViewerPositionForRendering returns the active 3D camera position and whether BeginMode3D is active.
func ViewerPositionForRendering() (pos rl.Vector3, in3D bool) {
	draw3dMu.Lock()
	defer draw3dMu.Unlock()
	return activeCamPos, inCamera3D
}

func MarkCamera3DEnd() {
	draw3dMu.Lock()
	defer draw3dMu.Unlock()
	inCamera3D = false
	activeCameraOK = false
}

// FlushDeferred3D draws queued mesh/model/instanced passes (shadow depth + color).
func FlushDeferred3D(h *heap.Store) {
	if h == nil {
		return
	}
	draw3dMu.Lock()
	meshes := deferredMeshes
	models := deferredModels
	inst := pendingInstDraw
	deferredMeshes = nil
	deferredModels = nil
	pendingInstDraw = nil
	draw3dMu.Unlock()

	if len(meshes) == 0 && len(models) == 0 && len(inst) == 0 {
		return
	}

	sh := mblight.ShadowCasterHandle()
	shadowOn := sh != 0 && shadowMapSize >= 256

	if shadowOn {
		ensureShadowResources()
		renderShadowPassDepth(h, meshes, models, inst)
	}

	for _, rec := range meshes {
		mo, err := heap.Cast[*meshObj](h, rec.meshH)
		if err != nil {
			continue
		}
		mato, err := heap.Cast[*materialObj](h, rec.matH)
		if err != nil {
			continue
		}
		if mato.pbr && shadowOn {
			bindPBRDrawState(&mato.mat, shadowOn)
		}
		rl.DrawMesh(mo.m, mato.mat, rec.mtx)
		if mato.pbr && shadowOn {
			clearShadowMapSlot(&mato.mat)
		}
	}

	for _, mh := range models {
		drawDeferredModelColorByHandle(h, mh, shadowOn)
	}

	for _, id := range inst {
		io, err := heap.Cast[*instancedModelObj](h, id.instH)
		if err != nil {
			continue
		}
		var lod *meshObj
		if id.lodMeshH != 0 {
			lod, _ = heap.Cast[*meshObj](h, id.lodMeshH)
		}
		drawInstancedRaster(io, lod, id.lodDist, shadowOn)
	}
}

func drawDeferredModelColorByHandle(h *heap.Store, mh heap.Handle, shadowOn bool) {
	if mo, err := heap.Cast[*modelObj](h, mh); err == nil {
		drawDeferredModelsColor(mo, shadowOn)
		return
	}
	if lo, err := heap.Cast[*lodModelObj](h, mh); err == nil {
		drawLODModelColor(lo, shadowOn)
	}
}

func drawDeferredModelsColor(mo *modelObj, shadowOn bool) {
	meshes := mo.model.GetMeshes()
	mats := mo.model.GetMaterials()
	mm := unsafe.Slice(mo.model.MeshMaterial, mo.model.MeshCount)
	for mi := int32(0); mi < mo.model.MeshCount; mi++ {
		mid := mm[mi]
		mat := mats[mid]
		mesh := meshes[mi]
		if shadowOn && shaderHasUniform(mat.Shader, "shadowEnabled") {
			rl.SetMaterialTexture(&mat, rl.MapBrdf, shadowRT.Depth)
			applyPBRUniformsIfAny(&mat, shadowOn)
		}
		rl.DrawMesh(mesh, mat, mo.model.Transform)
		if shadowOn && shaderHasUniform(mat.Shader, "shadowEnabled") {
			rl.SetMaterialTexture(&mat, rl.MapBrdf, rl.Texture2D{})
		}
	}
}

func drawLODModelColor(lo *lodModelObj, shadowOn bool) {
	li := lo.pickLOD(activeCamPos)
	if li < 0 {
		return
	}
	mod := &lo.models[li]
	saved := mod.Transform
	mod.Transform = lo.transform
	defer func() { mod.Transform = saved }()

	meshes := mod.GetMeshes()
	mats := mod.GetMaterials()
	mm := unsafe.Slice(mod.MeshMaterial, mod.MeshCount)
	for mi := int32(0); mi < mod.MeshCount; mi++ {
		mid := mm[mi]
		mat := mats[mid]
		mesh := meshes[mi]
		if shadowOn && shaderHasUniform(mat.Shader, "shadowEnabled") {
			rl.SetMaterialTexture(&mat, rl.MapBrdf, shadowRT.Depth)
			applyPBRUniformsIfAny(&mat, shadowOn)
		}
		rl.DrawMesh(mesh, mat, mod.Transform)
		if shadowOn && shaderHasUniform(mat.Shader, "shadowEnabled") {
			rl.SetMaterialTexture(&mat, rl.MapBrdf, rl.Texture2D{})
		}
	}
}

func drawDeferredModelShadowDepth(h *heap.Store, mh heap.Handle) {
	if mo, err := heap.Cast[*modelObj](h, mh); err == nil {
		meshes := mo.model.GetMeshes()
		for mi := int32(0); mi < mo.model.MeshCount; mi++ {
			rl.DrawMesh(meshes[mi], depthDrawMat, mo.model.Transform)
		}
		return
	}
	if lo, err := heap.Cast[*lodModelObj](h, mh); err == nil {
		li := lo.pickLOD(activeCamPos)
		if li < 0 {
			return
		}
		mod := &lo.models[li]
		meshes := mod.GetMeshes()
		for mi := int32(0); mi < mod.MeshCount; mi++ {
			rl.DrawMesh(meshes[mi], depthDrawMat, lo.transform)
		}
	}
}

// drawInstancedRaster draws an instanced batch (immediate MODEL.DRAW, INSTANCE.DRAW, or deferred color pass).
// When lodMesh is non-nil and lodDist > 0, uses lodMesh for instances beyond the distance threshold
// (camera to instance centroid). Per-instance colors: uniform tint uses one DrawMeshInstanced; varying
// tints fall back to per-instance DrawMesh (slower).
func drawInstancedRaster(io *instancedModelObj, lodMesh *meshObj, lodDist float32, shadowOn bool) {
	if io == nil {
		return
	}
	if io.shouldCull() {
		return
	}
	if io.meshIdx < 0 || io.meshIdx >= io.model.MeshCount {
		return
	}
	mi := io.meshIdx
	meshes := io.model.GetMeshes()
	mats := io.model.GetMaterials()
	mm := unsafe.Slice(io.model.MeshMaterial, io.model.MeshCount)
	mid := mm[mi]
	mesh := meshes[mi]
	if lodMesh != nil && lodDist > 0 {
		cam, _ := ViewerPositionForRendering()
		if rl.Vector3Distance(io.anchorPos(), cam) > lodDist {
			mesh = lodMesh.m
		}
	}
	mat := mats[mid]
	n := io.count
	if n <= 0 || len(io.transforms) < n {
		return
	}
	shadowed := shadowOn && shaderHasUniform(mat.Shader, "shadowEnabled")
	if shadowed {
		rl.SetMaterialTexture(&mat, rl.MapBrdf, shadowRT.Depth)
		applyPBRUniformsIfAny(&mat, shadowOn)
	}
	if shadowed {
		defer rl.SetMaterialTexture(&mat, rl.MapBrdf, rl.Texture2D{})
	}

	albedoMap := mat.GetMap(int32(rl.MapAlbedo))
	if io.uniformInstanceColors() && io.cr[0] == 255 && io.cg[0] == 255 && io.cb[0] == 255 && io.ca[0] == 255 {
		drawMeshInstancedCompat(mesh, mat, io.transforms[:n], n)
		return
	}
	if io.uniformInstanceColors() {
		saved := albedoMap.Color
		albedoMap.Color = rl.Color{R: uint8(io.cr[0]), G: uint8(io.cg[0]), B: uint8(io.cb[0]), A: uint8(io.ca[0])}
		drawMeshInstancedCompat(mesh, mat, io.transforms[:n], n)
		albedoMap.Color = saved
		return
	}
	for i := 0; i < n; i++ {
		saved := albedoMap.Color
		albedoMap.Color = rl.Color{R: uint8(io.cr[i]), G: uint8(io.cg[i]), B: uint8(io.cb[i]), A: uint8(io.ca[i])}
		rl.DrawMesh(mesh, mat, io.transforms[i])
		albedoMap.Color = saved
	}
}

func renderShadowPassDepth(h *heap.Store, meshes []deferredMeshRec, models []heap.Handle, inst []instancedDrawRec) {
	rl.BeginTextureMode(shadowRT)
	rl.ClearBackground(rl.White)
	lc := lightCamera()
	rl.BeginMode3D(lc)

	for _, rec := range meshes {
		mo, err := heap.Cast[*meshObj](h, rec.meshH)
		if err != nil {
			continue
		}
		rl.DrawMesh(mo.m, depthDrawMat, rec.mtx)
	}
	for _, mh := range models {
		drawDeferredModelShadowDepth(h, mh)
	}
	for _, id := range inst {
		io, err := heap.Cast[*instancedModelObj](h, id.instH)
		if err != nil {
			continue
		}
		if io.shouldCull() {
			continue
		}
		if io.meshIdx < 0 || io.meshIdx >= io.model.MeshCount {
			continue
		}
		mi := io.meshIdx
		n := io.count
		if n <= 0 || len(io.transforms) < n {
			continue
		}
		meshes := io.model.GetMeshes()
		drawMeshInstancedCompat(meshes[mi], depthDrawMat, io.transforms[:n], n)
	}

	rl.EndMode3D()
	rl.EndTextureMode()

	computeLightSpaceMatrix(lc)
}

func lightCamera() rl.Camera3D {
	hs := globalHeapForLight()
	hh := mblight.ShadowCasterHandle()
	dx, dy, dz, ok := mblight.LightDirection(hs, hh)
	if !ok {
		dx, dy, dz = 0, -1, 0
	}
	ext := float32(35)
	center := rl.Vector3{X: 0, Y: 2, Z: 0}
	if tx, ty, tz, ok2 := mblight.LightShadowTarget(hs, hh); ok2 {
		center = rl.Vector3{X: tx, Y: ty, Z: tz}
	}
	eye := rl.Vector3{
		X: center.X - dx*ext,
		Y: center.Y - dy*ext,
		Z: center.Z - dz*ext,
	}
	return rl.Camera3D{
		Position:   eye,
		Target:     center,
		Up:         rl.Vector3{X: 0, Y: 1, Z: 0},
		Fovy:       30,
		Projection: rl.CameraOrthographic,
	}
}

var globalHeapForLight func() *heap.Store = func() *heap.Store { return nil }

// SetGlobalHeapGetter allows mblight lookups during draws (direction / color).
func SetGlobalHeapGetter(fn func() *heap.Store) {
	globalHeapForLight = fn
}

func computeLightSpaceMatrix(cam rl.Camera3D) {
	view := rl.MatrixLookAt(cam.Position, cam.Target, cam.Up)
	proj := rl.MatrixOrtho(-18, 18, -18, 18, 0.5, 80)
	lightSpaceVP = rl.MatrixMultiply(proj, view)
}

func ensureShadowResources() {
	if shadowRTReady && shadowRT.Texture.Width == shadowMapSize {
		if !depthMatInited {
			depthDrawMat = rl.LoadMaterialDefault()
			depthMatInited = true
		}
		return
	}
	if shadowRT.Texture.ID > 0 {
		rl.UnloadRenderTexture(shadowRT)
		shadowRTReady = false
	}
	shadowRT = rl.LoadRenderTexture(shadowMapSize, shadowMapSize)
	shadowRTReady = true
	depthDrawMat = rl.LoadMaterialDefault()
	depthMatInited = true
}

// SetShadowMapResolution resizes the shadow depth target (RENDER.SETSHADOWMAPSIZE).
func SetShadowMapResolution(size int32) {
	draw3dMu.Lock()
	defer draw3dMu.Unlock()
	if size < 256 {
		size = 256
	}
	if size > 8192 {
		size = 8192
	}
	shadowMapSize = size
	if shadowRT.Texture.ID > 0 {
		rl.UnloadRenderTexture(shadowRT)
		shadowRTReady = false
	}
}

func bindPBRDrawState(mat *rl.Material, shadowOn bool) {
	if shadowOn {
		rl.SetMaterialTexture(mat, rl.MapBrdf, shadowRT.Depth)
	}
	applyPBRUniformsIfAny(mat, shadowOn)
}

func clearShadowMapSlot(mat *rl.Material) {
	rl.SetMaterialTexture(mat, rl.MapBrdf, rl.Texture2D{})
}

func shaderHasUniform(sh rl.Shader, name string) bool {
	return rl.GetShaderLocation(sh, name) >= 0
}

func applyPBRUniformsIfAny(mat *rl.Material, shadowOn bool) {
	sh := mat.Shader
	if sh.Locs == nil || rl.GetShaderLocation(sh, "roughnessValue") < 0 {
		return
	}
	defTex := rl.GetTextureIdDefault()
	nmap := mat.GetMap(rl.MapNormal)
	useNorm := int32(0)
	if nmap.Texture.ID != 0 && nmap.Texture.ID != defTex {
		useNorm = 1
	}

	rough := mat.GetMap(rl.MapRoughness).Value
	metal := mat.GetMap(rl.MapMetalness).Value
	if rough <= 0 {
		rough = 1
	}
	if metal <= 0 {
		metal = 1
	}

	hs := globalHeapForLight()
	hh := mblight.ShadowCasterHandle()
	lx, ly, lz, _ := mblight.LightDirection(hs, hh)
	lr, lg, lb := mblight.LightDiffuse(hs, hh)
	ar, ag, ab := sceneAmbientRGB()
	sbk := mblight.LightShadowBiasK(hs, hh)

	setInt := func(name string, v int32) {
		loc := rl.GetShaderLocation(sh, name)
		if loc < 0 {
			return
		}
		rl.SetShaderValue(sh, loc, []float32{float32(v)}, rl.ShaderUniformInt)
	}
	setFloat := func(name string, v float32) {
		loc := rl.GetShaderLocation(sh, name)
		if loc < 0 {
			return
		}
		rl.SetShaderValue(sh, loc, []float32{v}, rl.ShaderUniformFloat)
	}
	setVec3 := func(name string, x, y, z float32) {
		loc := rl.GetShaderLocation(sh, name)
		if loc < 0 {
			return
		}
		rl.SetShaderValue(sh, loc, []float32{x, y, z}, rl.ShaderUniformVec3)
	}

	setFloat("roughnessValue", rough)
	setFloat("metalnessValue", metal)
	setVec3("camPos", activeCamPos.X, activeCamPos.Y, activeCamPos.Z)
	setVec3("lightDir", lx, ly, lz)
	setVec3("lightColor", lr, lg, lb)
	setVec3("ambientColor", ar, ag, ab)
	setFloat("shadowBiasK", sbk)
	setInt("useNormalMap", useNorm)
	setInt("shadowEnabled", boolToInt(shadowOn))
	if shadowOn {
		loc := rl.GetShaderLocation(sh, "lightVP")
		if loc >= 0 {
			rl.SetShaderValueMatrix(sh, loc, lightSpaceVP)
		}
	}
}

func boolToInt(b bool) int32 {
	if b {
		return 1
	}
	return 0
}

func shadowDeferActive() bool {
	return mblight.ShadowCasterHandle() != 0 && shadowMapSize >= 256
}

func InCamera3D() bool {
	draw3dMu.Lock()
	defer draw3dMu.Unlock()
	return inCamera3D
}
