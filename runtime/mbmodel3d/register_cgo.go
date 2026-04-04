//go:build cgo

package mbmodel3d

import "moonbasic/runtime"

// Register wires MESH.*, MATERIAL.*, MODEL.*, SHADER.LOAD.
func (m *Module) Register(reg runtime.Registrar) {
	registerMeshGen(m, reg)
	registerMeshOps(m, reg)
	registerMaterialCmds(m, reg)
	registerShaderCmds(m, reg)
	registerModelLoad(m, reg)
	registerModelMaterial(m, reg)
	registerModelTextureStages(m, reg)
	registerModelRenderHierarchy(m, reg)
	registerModelTransform(m, reg)
	registerModelInstDraw(m, reg)
	registerModelLOD(m, reg)
}

// Shutdown implements runtime.Module.
func (m *Module) Shutdown() {}
