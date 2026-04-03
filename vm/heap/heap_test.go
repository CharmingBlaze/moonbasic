package heap

import "testing"

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
