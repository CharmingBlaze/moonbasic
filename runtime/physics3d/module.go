// Package mbphysics3d registers PHYSICS3D.* and BODY3D.* (Jolt on linux+cgo; stubs elsewhere).
package mbphysics3d

import (
	"sync"

	"github.com/tetratelabs/wazero/api"

	"moonbasic/internal/joltwasm"
	"moonbasic/runtime"
	"moonbasic/vm"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// ModulesByStore maps VM heaps to their physics module (for cross-package helpers).
var ModulesByStore = make(map[*heap.Store]*Module)

// GetModule returns the physics module associated with the given heap store, or nil if none.
func GetModule(h *heap.Store) *Module {
	return ModulesByStore[h]
}

// UserInvokerSetter allows the pipeline to wire VM.CallUserFunction for collision callbacks.
type UserInvokerSetter interface {
	SetUserInvoker(fn func(string, []value.Value) (value.Value, error))
}

// Module is the 3D physics orchestrator.
type Module struct {
	h           *heap.Store
	invoke      func(string, []value.Value) (value.Value, error)
	accumulator float64
	fixedStep   float64

	vmMu   sync.Mutex
	vmRef  *vm.VM
	wasmMu sync.Mutex
	// wasmPhysicsView is optional guest linear memory for Jolt WASM SoA readback (see joltwasm.StateView).
	wasmPhysicsView joltwasm.StateView

	// meshLookup is a bridge hook to get triangle data from ENTITY models for LEVEL.STATIC.
	meshLookup func(int64) []rl.Mesh
	// xformLookup is a bridge hook to get world position and rotation for VEHICLE integration.
	xformLookup func(int64) (pos rl.Vector3, yaw float32, ok bool)
	// xformUpdate is a bridge hook to set entity position.
	xformUpdate func(int64, rl.Vector3)

	vehicles map[int64]*Vehicle
}

// NewModule creates the physics3d module.
func NewModule() *Module {
	return &Module{
		vehicles: make(map[int64]*Vehicle),
	}
}

func (m *Module) BindHeap(h *heap.Store) {
	m.h = h
	ModulesByStore[h] = m
}

// SetUserInvoker wires callback dispatch from PHYSICS3D.PROCESSCOLLISIONS.
func (m *Module) SetUserInvoker(fn func(string, []value.Value) (value.Value, error)) {
	m.invoke = fn
}

// SetVM wires the active bytecode VM for WASM physics sync ([PHYSICS3D.SYNCWASMTOPHYSREGS] and automatic sync in STEP when a view is bound).
func (m *Module) SetVM(v *vm.VM) {
	m.vmMu.Lock()
	defer m.vmMu.Unlock()
	m.vmRef = v
}

// BindWasmPhysicsView sets the guest memory slice used by [joltwasm.UpdateVMPhysics]. Call when WASM linear memory or export bounds change.
func (m *Module) BindWasmPhysicsView(mem api.Memory, baseOffset, totalBytes uint32) {
	m.wasmMu.Lock()
	defer m.wasmMu.Unlock()
	m.wasmPhysicsView = joltwasm.StateView{Mem: mem, BaseOffset: baseOffset, TotalBytes: totalBytes}
}

// SetMeshLookup wires the entity mesh retrieval for LEVEL.STATIC collision generation.
func (m *Module) SetMeshLookup(fn func(int64) []rl.Mesh) {
	m.meshLookup = fn
}

// SetVehicleHooks wires the entity transform bridges for the VEHICLE system.
func (m *Module) SetVehicleHooks(lookup func(int64) (rl.Vector3, float32, bool), update func(int64, rl.Vector3)) {
	m.xformLookup = lookup
	m.xformUpdate = update
}

// SyncWasmPhysicsAfterStep copies guest SoA floats into [vm.VM.PhysicsScratch] and applies them to registers R0.. after each PHYSICS3D.STEP when a WASM view is bound.
func (m *Module) SyncWasmPhysicsAfterStep() {
	m.vmMu.Lock()
	v := m.vmRef
	m.vmMu.Unlock()
	if v == nil {
		return
	}
	m.wasmMu.Lock()
	view := m.wasmPhysicsView
	m.wasmMu.Unlock()
	if view.Mem == nil {
		return
	}
	joltwasm.UpdateVMPhysics(v, view)
	n := len(v.PhysicsScratch)
	if n > 256 {
		n = 256
	}
	_ = v.ExecSyncPhysics(0, n)
}

// Register implements runtime.Module.
func (m *Module) Register(reg runtime.Registrar) {
	reg.Register("VEHICLE.CREATE", "physics3d", runtime.AdaptLegacy(m.VHCreate))
	reg.Register("VEHICLE.SETWHEEL", "physics3d", runtime.AdaptLegacy(m.VHSetWheel))
	reg.Register("VEHICLE.CONTROL", "physics3d", runtime.AdaptLegacy(m.VHControl))
	reg.Register("VEHICLE.STEP", "physics3d", runtime.AdaptLegacy(m.VHStep))
	registerAeroCommands(m, reg)
	registerPhysics3DCommands(m, reg)
	registerBuoyancyCommands(m, reg)
}

// Shutdown implements runtime.Module.
func (m *Module) Shutdown() {
	if m.h != nil {
		delete(ModulesByStore, m.h)
	}
	shutdownPhysics3D(m)
}

func (m *Module) Reset() {}

