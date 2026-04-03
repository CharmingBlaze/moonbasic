package opcode

import "testing"

func TestChunkEmit(t *testing.T) {
	c := NewChunk("TEST")
	idx := c.AddInt(42)
	if idx != 0 {
		t.Fatalf("expected 0, got %d", idx)
	}
	if c.IntConsts[0] != 42 {
		t.Fatalf("expected 42, got %d", c.IntConsts[0])
	}
	c.Emit(OpPushInt, idx, 0, 100)
	if len(c.Instructions) != 1 {
		t.Fatal("expected 1 instr")
	}
	if c.Instructions[0].Op != OpPushInt {
		t.Fatal("wrong opcode")
	}
	if c.SourceLines[0] != 100 {
		t.Fatal("wrong line")
	}
}
