package vm

import (
	"testing"

	"moonbasic/runtime"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

type eraseAllTestObj struct {
	freed *int
}

func (o eraseAllTestObj) Free() {
	if o.freed != nil {
		(*o.freed)++
	}
}
func (eraseAllTestObj) TypeName() string { return "eraseAllTestObj" }
func (eraseAllTestObj) TypeTag() uint16  { return heap.TagInstance }

func TestEraseAllHandlesFreesAndClearsSlots(t *testing.T) {
	h := heap.New()
	reg := runtime.NewRegistry(h)
	reg.InitCore()
	v := New(reg, h)
	reg.EraseAllHandlesFn = v.EraseAllHandles
	defer func() { reg.EraseAllHandlesFn = nil }()

	var freed int
	h1, err := h.Alloc(eraseAllTestObj{freed: &freed})
	if err != nil {
		t.Fatal(err)
	}
	v.Globals["CAM"] = value.FromHandle(h1)
	v.Stack = append(v.Stack, value.FromHandle(h1))

	if err := v.EraseAllHandles(); err != nil {
		t.Fatal(err)
	}
	if freed != 1 {
		t.Fatalf("expected 1 Free, got %d", freed)
	}
	if h.Count() != 0 {
		t.Fatalf("expected empty heap, got %d live", h.Count())
	}
	if v.Globals["CAM"].Kind != value.KindNil {
		t.Fatalf("expected global nulled, got kind %v", v.Globals["CAM"].Kind)
	}
	if len(v.Stack) != 1 || v.Stack[0].Kind != value.KindNil {
		t.Fatalf("expected stack slot nulled, got %v", v.Stack)
	}
	reg.Shutdown()
}
