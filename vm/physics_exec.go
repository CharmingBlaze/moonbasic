package vm

import (
	"fmt"

	"moonbasic/vm/opcode"
	"moonbasic/vm/value"
)

// ExecSyncPhysics applies [VM.PhysicsScratch] to registers R[dst..dst+count-1] (same semantics as [opcode.OpSyncPhysics]).
func (v *VM) ExecSyncPhysics(dst uint8, count int) error {
	if count < 0 || count > 256 {
		return fmt.Errorf("ExecSyncPhysics: invalid count %d", count)
	}
	base := int(dst)
	if base+count > 256 {
		return fmt.Errorf("ExecSyncPhysics: register range overflow")
	}
	for j := 0; j < count; j++ {
		var f float64
		if j < len(v.PhysicsScratch) {
			f = v.PhysicsScratch[j]
		}
		v.setReg(uint8(base+j), value.FromFloat(f))
	}
	return nil
}

// RunSyncPhysicsOpcode runs one [opcode.OpSyncPhysics] instruction (operand=count, Dst=base register).
func (v *VM) RunSyncPhysicsOpcode(i opcode.Instruction) error {
	return v.ExecSyncPhysics(i.Dst, int(i.Operand))
}
