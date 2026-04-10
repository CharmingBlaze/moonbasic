//go:build cgo || (windows && !cgo)

package mbmodel3d

import (
	"unsafe"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// ConvertEntityModelMaterialsPBR replaces each material with the shared PBR shader while preserving
// loaded map textures; sets scalar metal/rough factors (multiplied with texture channels in-shader).
func ConvertEntityModelMaterialsPBR(mod *rl.Model, metal, rough float32) {
	if mod == nil {
		return
	}
	sh := pbrSharedShader()
	if !rl.IsShaderValid(sh) {
		return
	}
	mats := mod.GetMaterials()
	for i := range mats {
		mat := &mats[i]
		if shaderHasUniform(mat.Shader, "roughnessValue") {
			mat.GetMap(rl.MapMetalness).Value = metal
			mat.GetMap(rl.MapRoughness).Value = rough
			continue
		}
		alb := mat.GetMap(rl.MapAlbedo).Texture
		if alb.ID == 0 {
			alb = mat.GetMap(rl.MapDiffuse).Texture
		}
		metTex := mat.GetMap(rl.MapMetalness).Texture
		rouTex := mat.GetMap(rl.MapRoughness).Texture
		normTex := mat.GetMap(rl.MapNormal).Texture
		emTex := mat.GetMap(rl.MapEmission).Texture

		old := mat.Shader
		mat.Shader = sh
		patchStandardMapTextureLocs(&mat.Shader)
		if alb.ID != 0 {
			rl.SetMaterialTexture(mat, rl.MapAlbedo, alb)
		}
		if metTex.ID != 0 {
			rl.SetMaterialTexture(mat, rl.MapMetalness, metTex)
		}
		if rouTex.ID != 0 {
			rl.SetMaterialTexture(mat, rl.MapRoughness, rouTex)
		}
		if normTex.ID != 0 {
			rl.SetMaterialTexture(mat, rl.MapNormal, normTex)
		}
		if emTex.ID != 0 {
			rl.SetMaterialTexture(mat, rl.MapEmission, emTex)
		}
		mat.GetMap(rl.MapMetalness).Value = metal
		mat.GetMap(rl.MapRoughness).Value = rough
		if mat.GetMap(rl.MapEmission).Value <= 0 && emTex.ID == 0 {
			mat.GetMap(rl.MapEmission).Value = 0
		}
		if old.ID != sh.ID && rl.IsShaderValid(old) {
			rl.UnloadShader(old)
		}
	}
}

// DrawEntityModel draws a loaded model with the same PBR uniform path as deferred MODEL.* passes
// (directional light, shadow map, camera) so ENTITY.DRAWALL matches the modern renderer.
func DrawEntityModel(mod rl.Model, tint rl.Color) {
	shadowOn := shadowDeferActive()
	meshes := mod.GetMeshes()
	mats := mod.GetMaterials()
	mm := unsafe.Slice(mod.MeshMaterial, mod.MeshCount)
	for mi := int32(0); mi < mod.MeshCount; mi++ {
		mid := mm[mi]
		mat := &mats[mid]
		alb := mat.GetMap(rl.MapAlbedo)
		orig := alb.Color
		alb.Color = tint
		mmat := *mat
		if shadowOn && shaderHasUniform(mmat.Shader, "shadowEnabled") {
			rl.SetMaterialTexture(&mmat, rl.MapBrdf, shadowRT.Depth)
		}
		applyPBRUniformsIfAny(&mmat, shadowOn, nil)
		rl.DrawMesh(meshes[mi], mmat, mod.Transform)
		if shadowOn && shaderHasUniform(mmat.Shader, "shadowEnabled") {
			rl.SetMaterialTexture(&mmat, rl.MapBrdf, rl.Texture2D{})
		}
		alb.Color = orig
	}
}
