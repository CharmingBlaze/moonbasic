package driver

import (
	"os"
	"testing"
)

func TestGetDefaultDriverUnknownEnv(t *testing.T) {
	t.Setenv(EnvDriver, "not-a-valid-value")
	sel := GetDefaultDriver()
	if sel.Kind != KindUnavailable {
		t.Fatalf("expected KindUnavailable, got %v", sel.Kind)
	}
}

func TestGetDefaultDriverPuregoOverride(t *testing.T) {
	t.Setenv(EnvDriver, "purego")
	// May succeed (DLL present) or fail; kind must be purego or unavailable, never native_cgo on typical dev machine without cgo tag.
	sel := GetDefaultDriver()
	if sel.Kind != KindPuregoDLL && sel.Kind != KindUnavailable {
		t.Fatalf("unexpected kind %v", sel.Kind)
	}
	_ = os.Getenv(EnvDriver)
}
