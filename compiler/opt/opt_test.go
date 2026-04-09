package opt

import (
	"testing"

	"moonbasic/vm/opcode"
)

func TestThreadJumpChain(t *testing.T) {
	ch := opcode.NewChunk("t")
	// IP 0: jump to 1; IP 1: jump to 2; IP 2: halt
	// Emit(op, dst, srcA, srcB, operand, line)
	ch.Emit(opcode.OpJump, 0, 0, 0, 1, 1)
	ch.Emit(opcode.OpJump, 0, 0, 0, 2, 1)
	ch.Emit(opcode.OpHalt, 0, 0, 0, 0, 1)
	OptimizeChunk(ch)
	if ch.Instructions[0].Op != opcode.OpJump || ch.Instructions[0].Operand != 2 {
		t.Fatalf("first jump should target halt at 2, got %+v", ch.Instructions[0])
	}
}
