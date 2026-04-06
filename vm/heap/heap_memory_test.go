package heap

import "testing"

type countFreeObj struct {
	n *int
}

func (c *countFreeObj) Free() {
	if c.n != nil {
		*c.n++
	}
}

func (c *countFreeObj) TypeName() string { return "countFree" }

func (c *countFreeObj) TypeTag() uint16 { return TagInstance }

func TestFreeAllInvokesFreeOnEveryLiveObject(t *testing.T) {
	var freed int
	s := New()
	for i := 0; i < 50; i++ {
		_, err := s.Alloc(&countFreeObj{n: &freed})
		if err != nil {
			t.Fatalf("Alloc: %v", err)
		}
	}
	s.FreeAll()
	if freed != 50 {
		t.Fatalf("FreeAll: expected 50 Free calls, got %d", freed)
	}
	if s.Count() != 0 {
		t.Fatalf("FreeAll: expected 0 live objects, got %d", s.Count())
	}
}

func TestHandleZeroIsInvalid(t *testing.T) {
	s := New()
	if _, ok := s.Get(0); ok {
		t.Fatal("handle 0 must be invalid")
	}
}

func TestFreeInvalidHandleTwiceErrors(t *testing.T) {
	s := New()
	h, err := s.Alloc(dummyObj{})
	if err != nil {
		t.Fatal(err)
	}
	if err := s.Free(h); err != nil {
		t.Fatal(err)
	}
	if err := s.Free(h); err == nil {
		t.Fatal("second Free on same handle should error")
	}
}
