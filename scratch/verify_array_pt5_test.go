package main

import (
	"testing"
	"moonbasic/runtime"
	"moonbasic/runtime/mbarray"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

type mockRegistrar struct {
	cmds map[string]runtime.BuiltinFn
}

func (m *mockRegistrar) Register(name, help string, fn runtime.BuiltinFn) {
	if m.cmds == nil {
		m.cmds = make(map[string]runtime.BuiltinFn)
	}
	m.cmds[name] = fn
}

func TestArrayStandardization(t *testing.T) {
	h := heap.New()
	mod := mbarray.NewModule()
	mod.BindHeap(h)
	
	r := &mockRegistrar{}
	mod.Register(r)
	
	rt := &runtime.Runtime{Heap: h}
	
	// 1. Test ARRAY.CREATE
	createFn, ok := r.cmds["ARRAY.CREATE"]
	if !ok {
		t.Fatal("ARRAY.CREATE not registered")
	}
	
	res, err := createFn(rt, value.FromInt(10))
	if err != nil {
		t.Fatalf("ARRAY.CREATE failed: %v", err)
	}
	if res.Kind != value.KindHandle {
		t.Fatalf("Expected handle, got %v", res.Kind)
	}
	arrHandle := res
	
	// 2. Test ARRAY.FILL (Fluent)
	fillFn, ok := r.cmds["ARRAY.FILL"]
	if !ok {
		t.Fatal("ARRAY.FILL not registered")
	}
	
	res, err = fillFn(rt, arrHandle, value.FromFloat(42.0))
	if err != nil {
		t.Fatalf("ARRAY.FILL failed: %v", err)
	}
	if res.Kind != value.KindHandle || res.IVal != arrHandle.IVal {
		t.Fatalf("ARRAY.FILL should return the same handle for chaining")
	}
	
	// 3. Verify content
	lenFn := r.cmds["ARRAY.LEN"]
	res, _ = lenFn(rt, arrHandle)
	if res.IVal != 10 {
		t.Errorf("Expected len 10, got %d", res.IVal)
	}
	
	t.Log("Array standardization verified successfully")
}
