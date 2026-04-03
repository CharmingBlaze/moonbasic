package builtinmanifest

import "testing"

func TestDefaultManifestLoads(t *testing.T) {
	tbl := Default()
	if len(tbl.Commands) < 10 {
		t.Fatalf("expected many commands, got %d", len(tbl.Commands))
	}
	if _, ok := tbl.LookupArity("WINDOW", "OPEN", 3); !ok {
		t.Fatal("missing WINDOW.OPEN (3-arg overload)")
	}
}

func TestKeysSorted(t *testing.T) {
	keys := Default().Keys()
	if len(keys) < 2 {
		t.Fatal("no keys")
	}
	for i := 1; i < len(keys); i++ {
		if keys[i] < keys[i-1] {
			t.Fatalf("keys not sorted: %q before %q", keys[i-1], keys[i])
		}
	}
}

func TestManifestOptionalMetadata(t *testing.T) {
	c, ok := Default().LookupArity("WINDOW", "OPEN", 3)
	if !ok {
		t.Fatal("missing WINDOW.OPEN")
	}
	if c.Phase != "init" {
		t.Fatalf("WINDOW.OPEN phase: got %q want init", c.Phase)
	}
	if c.Returns != "bool" {
		t.Fatalf("WINDOW.OPEN returns: got %q want bool", c.Returns)
	}
}
