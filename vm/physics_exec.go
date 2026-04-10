package vm

import (
	"fmt"

	"moonbasic/vm/opcode"
)

// ExecSyncPhysics writes entity transforms directly into a pre-allocated segment of the WASM memory buffer.
// Currently stubbed out to bypass v.PhysicsScratch copy and map to the generic linear memory slice format: BaseAddress + (EntityID * Stride).
func (v *VM) ExecSyncPhysics(dst uint8, count int) error {
	if count < 0 || count > 256 {
		return fmt.Errorf("ExecSyncPhysics: invalid count %d", count)
	}
	base := int(dst)
	if base+count > 256 {
		return fmt.Errorf("ExecSyncPhysics: register range overflow")
	}
	// Direct memory offsets: BaseAddress + (EntityID * Stride)
	// Handled through natively bound WASM linear memory view.
	return nil
}

// RunSyncPhysicsOpcode runs one [opcode.OpSyncPhysics] instruction (operand=count, Dst=base register).
func (v *VM) RunSyncPhysicsOpcode(i opcode.Instruction) error {
	return v.ExecSyncPhysics(i.Dst, int(i.Operand))
}
