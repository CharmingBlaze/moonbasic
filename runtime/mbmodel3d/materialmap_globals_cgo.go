//go:build cgo

package mbmodel3d

import (
	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/vm/value"
)

// SeedMaterialMapGlobals installs MATERIAL_MAP_* indices for Raylib material maps.
func SeedMaterialMapGlobals(globals map[string]value.Value) {
	if globals == nil {
		return
	}
	globals["MATERIAL_MAP_ALBEDO"] = value.FromInt(int64(rl.MapAlbedo))
	globals["MATERIAL_MAP_METALNESS"] = value.FromInt(int64(rl.MapMetalness))
	globals["MATERIAL_MAP_NORMAL"] = value.FromInt(int64(rl.MapNormal))
	globals["MATERIAL_MAP_ROUGHNESS"] = value.FromInt(int64(rl.MapRoughness))
	globals["MATERIAL_MAP_OCCLUSION"] = value.FromInt(int64(rl.MapOcclusion))
	globals["MATERIAL_MAP_EMISSION"] = value.FromInt(int64(rl.MapEmission))
	globals["MATERIAL_MAP_HEIGHT"] = value.FromInt(int64(rl.MapHeight))
	globals["MATERIAL_MAP_CUBEMAP"] = value.FromInt(int64(rl.MapCubemap))
	globals["MATERIAL_MAP_IRRADIANCE"] = value.FromInt(int64(rl.MapIrradiance))
	globals["MATERIAL_MAP_PREFILTER"] = value.FromInt(int64(rl.MapPrefilter))
	globals["MATERIAL_MAP_BRDF"] = value.FromInt(int64(rl.MapBrdf))
	globals["MATERIAL_MAP_DIFFUSE"] = value.FromInt(int64(rl.MapDiffuse))
	globals["MATERIAL_MAP_SPECULAR"] = value.FromInt(int64(rl.MapSpecular))
	// MATERIAL.SETFLOAT(mat, MATERIAL_ROUGHNESS, v) targets maps[roughness].Value (same index as map slot).
	globals["MATERIAL_ROUGHNESS"] = value.FromInt(int64(rl.MapRoughness))
	globals["MATERIAL_METALNESS"] = value.FromInt(int64(rl.MapMetalness))
}
