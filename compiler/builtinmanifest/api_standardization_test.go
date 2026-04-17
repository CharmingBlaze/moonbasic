package builtinmanifest

import (
	"fmt"
	"strings"
	"testing"
)

func TestCreateHasMatchingMakeDeprecatedAlias(t *testing.T) {
	tbl := Default()
	for key, ovs := range tbl.Commands {
		if !strings.Contains(key, ".CREATE") {
			continue
		}
		makeKey := strings.Replace(key, ".CREATE", ".MAKE", 1)
		makeOVs := tbl.Commands[makeKey]
		if len(makeOVs) == 0 {
			t.Fatalf("missing deprecated MAKE alias for %s (expected %s)", key, makeKey)
		}
		for _, c := range ovs {
			found := false
			for _, a := range makeOVs {
				if len(a.Args) != len(c.Args) {
					continue
				}
				found = true
				if !strings.Contains(strings.ToUpper(a.Desc), "DEPRECATED") {
					t.Fatalf("%s arity %d alias must include deprecation text", makeKey, len(a.Args))
				}
				break
			}
			if !found {
				t.Fatalf("missing %s overload matching %s arity %d", makeKey, key, len(c.Args))
			}
		}
	}
}

// TestNoDuplicateAudioManifestArgSignatures guards the AUDIO namespace against redundant
// identical overload rows (same arg-kind list twice), which breaks arity-based overload selection.
func TestNoDuplicateAudioManifestArgSignatures(t *testing.T) {
	tbl := Default()
	for key, ovs := range tbl.Commands {
		if !strings.HasPrefix(key, "AUDIO.") {
			continue
		}
		seen := make(map[string]bool)
		for _, c := range ovs {
			sig := fmt.Sprintf("%v", c.Args)
			if seen[sig] {
				t.Fatalf("duplicate AUDIO manifest overload for %s with arg kinds %v", key, c.Args)
			}
			seen[sig] = true
		}
	}
}

func TestSetPosHasSetPositionAlias(t *testing.T) {
	tbl := Default()
	for key, ovs := range tbl.Commands {
		if !strings.HasSuffix(key, ".SETPOS") || strings.Contains(key, ".SETPOSITION") {
			continue
		}
		aliasKey := strings.Replace(key, ".SETPOS", ".SETPOSITION", 1)
		aliasOVs := tbl.Commands[aliasKey]
		if len(aliasOVs) == 0 {
			t.Fatalf("missing SETPOSITION alias for %s (expected %s)", key, aliasKey)
		}
		for _, c := range ovs {
			found := false
			for _, a := range aliasOVs {
				if len(a.Args) == len(c.Args) {
					found = true
					break
				}
			}
			if !found {
				t.Fatalf("missing %s overload matching %s arity %d", aliasKey, key, len(c.Args))
			}
		}
	}
}

func TestPhysics3DUpdateMirrorsStepInManifest(t *testing.T) {
	tbl := Default()
	stepOVs := tbl.Commands["PHYSICS3D.STEP"]
	upOVs := tbl.Commands["PHYSICS3D.UPDATE"]
	if len(stepOVs) == 0 {
		t.Fatal("missing PHYSICS3D.STEP in manifest")
	}
	if len(upOVs) == 0 {
		t.Fatal("missing PHYSICS3D.UPDATE in manifest")
	}
	if len(stepOVs) != len(upOVs) {
		t.Fatalf("PHYSICS3D.UPDATE overload count %d != PHYSICS3D.STEP %d", len(upOVs), len(stepOVs))
	}
	for i := range stepOVs {
		if len(stepOVs[i].Args) != len(upOVs[i].Args) {
			t.Fatalf("overload %d: STEP args %v vs UPDATE args %v", i, stepOVs[i].Args, upOVs[i].Args)
		}
	}
}

