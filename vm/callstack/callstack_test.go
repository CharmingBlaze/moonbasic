package callstack

import (
	"testing"
	"moonbasic/vm/opcode"
)

func TestPushPop(t *testing.T) {
	s := New()
	chunk := &opcode.Chunk{Name: "TEST"}
	s.Push(chunk, 42, 10)
	
	f := s.Top()
	if f == nil || f.IP != 42 || f.StackBase != 10 {
		t.Fatal(f)
	}
	
	f2 := s.Pop()
	if f2.IP != 42 {
		t.Fatal(f2)
	}
	
	if s.Depth() != 0 {
		t.Fatal(s.Depth())
	}
}
