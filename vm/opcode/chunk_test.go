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
	c.Emit(OpPushInt, 0, 0, 0, idx, 100)
	if len(c.Instructions) != 1 {
		t.Fatal("expected 1 instr")
	}
	if len(c.ArrayDebugName) != 1 || c.ArrayDebugName[0] != -1 {
		t.Fatalf("ArrayDebugName: got %v", c.ArrayDebugName)
	}
	if c.Instructions[0].Op != OpPushInt {
		t.Fatal("wrong opcode")
	}
	if c.SourceLines[0] != 100 {
		t.Fatal("wrong line")
	}
}
