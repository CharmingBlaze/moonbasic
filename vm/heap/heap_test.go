package heap

import (
	"testing"

	"moonbasic/vm/opcode"
)
type dummyObj struct{}

func (dummyObj) Free()            {}
func (dummyObj) TypeName() string { return "dummy" }
func (dummyObj) TypeTag() uint16  { return TagInstance }

func TestStaleHandleAfterFree(t *testing.T) {
	s := New()
	h, _ := s.Alloc(dummyObj{})
	if _, ok := s.Get(h); !ok {
		t.Fatal("expected live handle")
	}
	_ = s.Free(h)
	if _, ok := s.Get(h); ok {
		t.Fatal("expected stale handle after free")
	}
	h2, _ := s.Alloc(dummyObj{})
	if h2 == h {
		t.Fatal("expected new encoding after slot reuse")
	}
	if _, ok := s.Get(h2); !ok {
		t.Fatal("expected new handle valid")
	}
}

func TestSeedProgramStringsBytecodeIndices(t *testing.T) {
	s := New()
	p := opcode.NewProgram()
	if p.InternString("hello") != 0 || p.InternString("ll") != 1 {
		t.Fatal("expected program pool indices 0 and 1")
	}
	s.SeedProgramStrings(p.StringTable)
	if a, ok := s.GetString(0); !ok || a != "hello" {
		t.Fatalf("GetString(0): ok=%v got %q want hello", ok, a)
	}
	if b, ok := s.GetString(1); !ok || b != "ll" {
		t.Fatalf("GetString(1): ok=%v got %q want ll", ok, b)
	}
	// Runtime Intern continues after the program table.
	if idx := s.Intern("extra"); idx != 2 {
		t.Fatalf("Intern extra: got index %d want 2", idx)
	}
}
