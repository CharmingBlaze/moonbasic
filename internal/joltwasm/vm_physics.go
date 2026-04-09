package joltwasm

// PhysicsScratchSink receives packed float32 physics data for [opcode.OpSyncPhysics] (implemented by *vm.VM).
// Kept as an interface so this package does not import moonbasic/vm (which would pull runtime → raylib init).
type PhysicsScratchSink interface {
	SyncPhysicsFromFloat32View([]float32)
}

// UpdateVMPhysics reads the guest SoA float block after [PhysicsStateHeader] and passes it to the VM scratch buffer.
// Call each frame after WASM memory may have grown so [StateView.Mem.Read] returns a fresh view.
func UpdateVMPhysics(sink PhysicsScratchSink, view StateView) bool {
	floats, ok := view.FloatsAfterHeader()
	if !ok {
		sink.SyncPhysicsFromFloat32View(nil)
		return false
	}
	sink.SyncPhysicsFromFloat32View(floats)
	return true
}
