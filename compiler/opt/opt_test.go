package opt

import (
	"testing"

	"moonbasic/vm/opcode"
)

func TestPeepholePushPop(t *testing.T) {
	ch := opcode.NewChunk("t")
	ch.Emit(opcode.OpPushInt, 0, 0, 1)
	ch.Emit(opcode.OpPop, 0, 0, 1)
	ch.Emit(opcode.OpHalt, 0, 0, 1)
	OptimizeChunk(ch)
	if len(ch.Instructions) != 1 || ch.Instructions[0].Op != opcode.OpHalt {
		t.Fatalf("got %+v", ch.Instructions)
	}
}

func TestThreadJumpChain(t *testing.T) {
	ch := opcode.NewChunk("t")
	// IP 0: jump to 1; IP 1: jump to 2; IP 2: halt
	ch.Emit(opcode.OpJump, 1, 0, 1)
	ch.Emit(opcode.OpJump, 2, 0, 1)
	ch.Emit(opcode.OpHalt, 0, 0, 1)
	OptimizeChunk(ch)
	if ch.Instructions[0].Op != opcode.OpJump || ch.Instructions[0].Operand != 2 {
		t.Fatalf("first jump should target halt at 2, got %+v", ch.Instructions[0])
	}
}
