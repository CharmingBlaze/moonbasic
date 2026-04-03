package symtable

import "testing"

func TestScopeVar(t *testing.T) {
	tab := New()
	tab.PushScope()
	tab.DefineParam("A")
	tab.DefineLocal("X")
	if !tab.IsVar("A") || !tab.IsVar("X") {
		t.Fatal()
	}
	tab.PopScope()
	if tab.IsVar("X") {
		t.Fatal("local should be gone")
	}
}

func TestGlobalVar(t *testing.T) {
	tab := New()
	tab.DefineGlobalVar("G")
	if !tab.IsVar("G") {
		t.Fatal()
	}
}
