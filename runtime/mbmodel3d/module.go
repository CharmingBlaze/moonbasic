// Package mbmodel3d implements MESH.*, MATERIAL.*, MODEL.*, and SHADER.LOAD (Raylib 3D) when CGO is enabled.
package mbmodel3d

import "moonbasic/vm/heap"

// Module registers 3D mesh/material/model/shader builtins.
type Module struct {
	h *heap.Store
}

// NewModule creates the module.
func NewModule() *Module { return &Module{} }

// BindHeap binds the VM heap before Register (implements runtime.HeapAware).
func (m *Module) BindHeap(h *heap.Store) { m.h = h }
