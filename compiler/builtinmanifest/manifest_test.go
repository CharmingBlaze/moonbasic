package builtinmanifest

import (
	"strings"
	"testing"
)

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
	if c.Returns != "" {
		t.Fatalf("WINDOW.OPEN returns: got %q want void (empty)", c.Returns)
	}
}

// TestSetPositionDeprecationMatchesSetPosKeys ensures every namespace that declares *.SETPOSITION
// also has *.SETPOS and resolves DeprecationReplacement(ns, "SETPOSITION") to NS.SETPOS.
func TestSetPositionDeprecationMatchesSetPosKeys(t *testing.T) {
	tbl := Default()
	seen := make(map[string]bool)
	for key := range tbl.Commands {
		if !strings.HasSuffix(key, ".SETPOSITION") {
			continue
		}
		dot := strings.IndexByte(key, '.')
		if dot < 0 {
			continue
		}
		ns := key[:dot]
		if seen[ns] {
			continue
		}
		seen[ns] = true
		want := ns + ".SETPOS"
		if _, exists := tbl.Commands[want]; !exists {
			t.Fatalf("manifest has %s.SETPOSITION but no %s", ns, want)
		}
		got, ok := tbl.DeprecationReplacement(ns, "SETPOSITION")
		if !ok || got != want {
			t.Fatalf("DeprecationReplacement(%q, SETPOSITION) = (%q, %v), want (%q, true)", ns, got, ok, want)
		}
	}
}

// TestMakeAliasDeprecationMatchesCanonicalKeys ensures every manifest command whose method
// starts with MAKE and whose first-MAKE→CREATE sibling exists resolves via DeprecationReplacement.
// Covers TIMER, INSTANCE, MODEL instancing, NOISE.*, JOINT.MAKE*, etc., without hand-maintained cases.
func TestMakeAliasDeprecationMatchesCanonicalKeys(t *testing.T) {
	tbl := Default()
	for key := range tbl.Commands {
		dot := strings.IndexByte(key, '.')
		if dot < 0 {
			continue
		}
		ns, method := key[:dot], key[dot+1:]
		if !strings.HasPrefix(method, "MAKE") {
			continue
		}
		replMethod := strings.Replace(method, "MAKE", "CREATE", 1)
		canonical := ns + "." + replMethod
		if _, exists := tbl.Commands[canonical]; !exists {
			continue
		}
		got, ok := tbl.DeprecationReplacement(ns, method)
		if !ok || got != canonical {
			t.Fatalf("DeprecationReplacement(%q, %q) = (%q, %v), want (%q, true)", ns, method, got, ok, canonical)
		}
	}
}
