//go:build cgo || (windows && !cgo)

package mbphysics3d

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

// This file provides package-level Go functions for internal engine use (e.g. by mbentity).
// These functions delegate to the appropriate Module instance found via GetModule(heap).

// PHWorldSetup is a package-level facade for Module.phWorldSetup.
func PHWorldSetup(m *Module, args []value.Value) (value.Value, error) {
	if m == nil {
		return value.Nil, fmt.Errorf("physics3d: module nil")
	}
	return m.phWorldSetup(args)
}

// BDAddMesh is a package-level facade for Module.bdAddMesh.
func BDAddMesh(h *heap.Store, args []value.Value) (value.Value, error) {
	m := GetModule(h)
	if m == nil {
		return value.Nil, fmt.Errorf("physics3d: module not found for heap")
	}
	return m.bdAddMesh(args)
}

// BDCommit is a package-level facade for Module.bdCommit.
func BDCommit(h *heap.Store, args []value.Value) (value.Value, error) {
	m := GetModule(h)
	if m == nil {
		return value.Nil, fmt.Errorf("physics3d: module not found for heap")
	}
	return m.bdCommit(args)
}

// BDBufferIndex is a package-level facade for Module.bdBufferIndex.
func BDBufferIndex(h *heap.Store, args []value.Value) (value.Value, error) {
	m := GetModule(h)
	if m == nil {
		return value.Nil, fmt.Errorf("physics3d: module not found for heap")
	}
	return m.bdBufferIndex(args)
}

// BDAddSphere is a package-level facade for Module.bdAddSphere.
func BDAddSphere(h *heap.Store, args []value.Value) (value.Value, error) {
	m := GetModule(h)
	if m == nil {
		return value.Nil, fmt.Errorf("physics3d: module not found for heap")
	}
	return m.bdAddSphere(args)
}

// BDAddCapsule is a package-level facade for Module.bdAddCapsule.
func BDAddCapsule(h *heap.Store, args []value.Value) (value.Value, error) {
	m := GetModule(h)
	if m == nil {
		return value.Nil, fmt.Errorf("physics3d: module not found for heap")
	}
	return m.bdAddCapsule(args)
}

// BDAddBox is a package-level facade for Module.bdAddBox.
func BDAddBox(h *heap.Store, args []value.Value) (value.Value, error) {
	m := GetModule(h)
	if m == nil {
		return value.Nil, fmt.Errorf("physics3d: module not found for heap")
	}
	return m.bdAddBox(args)
}

// SetMeshLookupForHeap is a package-level facade for GetModule(h).SetMeshLookup (mbentity / engine bridges).
func SetMeshLookupForHeap(h *heap.Store, fn func(int64) []rl.Mesh) {
	m := GetModule(h)
	if m == nil {
		return
	}
	m.SetMeshLookup(fn)
}

// SetVehicleHooksForHeap is a package-level facade for GetModule(h).SetVehicleHooks.
func SetVehicleHooksForHeap(h *heap.Store, lookup func(int64) (rl.Vector3, float32, bool), update func(int64, rl.Vector3)) {
	m := GetModule(h)
	if m == nil {
		return
	}
	m.SetVehicleHooks(lookup, update)
}
