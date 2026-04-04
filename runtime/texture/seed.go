package texture

import "moonbasic/vm/value"

// SeedTextureGlobals installs Raylib texture filter/wrap enum values for TEXTURE.SETFILTER / SETWRAP.
func SeedTextureGlobals(globals map[string]value.Value) {
	globals["FILTER_POINT"] = value.FromInt(0)
	globals["FILTER_BILINEAR"] = value.FromInt(1)
	globals["FILTER_TRILINEAR"] = value.FromInt(2)
	globals["FILTER_ANISOTROPIC_4X"] = value.FromInt(3)
	globals["FILTER_ANISOTROPIC_8X"] = value.FromInt(4)
	globals["FILTER_ANISOTROPIC_16X"] = value.FromInt(5)
	globals["WRAP_REPEAT"] = value.FromInt(0)
	globals["WRAP_CLAMP"] = value.FromInt(1)
	globals["WRAP_MIRROR_REPEAT"] = value.FromInt(2)
	globals["WRAP_MIRROR_CLAMP"] = value.FromInt(3)
}
