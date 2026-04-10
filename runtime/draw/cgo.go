//go:build cgo || (windows && !cgo)

package mbdraw

import (
	"moonbasic/runtime"
	"moonbasic/vm/value"
)

// Register implements runtime.Module.
func (m *Module) Register(r runtime.Registrar) {
	registerShapeCmds(m, r)
	registerTextureCmds(m, r)
	registerLineCmds(m, r)
	registerTextCmds(m, r)
	registerDebugPrint(m, r)
	registerAdvancedCmds(m, r)
	registerDraw3DCmds(m, r)
	registerPrim3DWrappers(m, r)
	registerPrim2DWrappers(m, r)
	registerTextTextureObjs(m, r)
	registerTextureAdvWrappers(m, r)
	registerTextExObj(m, r)
	registerDrawNamespaceAliases(m, r)
	registerCircleExtraCmds(m, r)
	registerDrawHelperCmds(m, r)
}

// Shutdown implements runtime.Module.
func (m *Module) Shutdown() {}

func argInt(v value.Value) (int32, bool) {
	if i, ok := v.ToInt(); ok {
		return int32(i), true
	}
	if f, ok := v.ToFloat(); ok {
		return int32(f), true
	}
	return 0, false
}

func argFloat(v value.Value) (float32, bool) {
	if f, ok := v.ToFloat(); ok {
		return float32(f), true
	}
	if i, ok := v.ToInt(); ok {
		return float32(i), true
	}
	return 0, false
}
