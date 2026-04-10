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
)

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
}

// NewModule creates the physics3d module.
func NewModule() *Module { return &Module{} }

// BindHeap implements runtime.HeapAware.
func (m *Module) BindHeap(h *heap.Store) { m.h = h }

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
	registerPhysics3DCommands(m, reg)
	registerBuoyancyCommands(m, reg)
}

// Shutdown implements runtime.Module.
func (m *Module) Shutdown() {
	shutdownPhysics3D(m)
}

func (m *Module) Reset() {}

