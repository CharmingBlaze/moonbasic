package builtinmanifest

import "testing"

func TestHasArityExact(t *testing.T) {
	tbl := Default()
	if !tbl.HasArityExact("DrawEntities", 0) {
		t.Fatal("DrawEntities should have 0-arg overload")
	}
	if tbl.HasArityExact("DrawEntities", 1) {
		t.Fatal("DrawEntities should not match arity 1")
	}
	if !tbl.HasArityExact("TURNENTITY", 4) {
		t.Fatal("TURNENTITY should have 4-arg overload")
	}
}
