//go:build !cgo

package mbmodel3d

import "moonbasic/vm/value"

// SeedMaterialMapGlobals installs MATERIAL_MAP_* indices matching raylib (for semantic check / stub runs).
func SeedMaterialMapGlobals(globals map[string]value.Value) {
	if globals == nil {
		return
	}
	globals["MATERIAL_MAP_ALBEDO"] = value.FromInt(0)
	globals["MATERIAL_MAP_METALNESS"] = value.FromInt(1)
	globals["MATERIAL_MAP_NORMAL"] = value.FromInt(2)
	globals["MATERIAL_MAP_ROUGHNESS"] = value.FromInt(3)
	globals["MATERIAL_MAP_OCCLUSION"] = value.FromInt(4)
	globals["MATERIAL_MAP_EMISSION"] = value.FromInt(5)
	globals["MATERIAL_MAP_HEIGHT"] = value.FromInt(6)
	globals["MATERIAL_MAP_CUBEMAP"] = value.FromInt(8)
	globals["MATERIAL_MAP_IRRADIANCE"] = value.FromInt(9)
	globals["MATERIAL_MAP_PREFILTER"] = value.FromInt(10)
	globals["MATERIAL_MAP_BRDF"] = value.FromInt(11)
	globals["MATERIAL_MAP_DIFFUSE"] = value.FromInt(0)
	globals["MATERIAL_MAP_SPECULAR"] = value.FromInt(1)
	globals["MATERIAL_ROUGHNESS"] = value.FromInt(3)
	globals["MATERIAL_METALNESS"] = value.FromInt(1)
}
