package raylibpurego

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLibPathNonEmpty(t *testing.T) {
	if LibPath() == "" {
		t.Fatal("LibPath empty")
	}
}

func TestLoadFromMissingDLL(t *testing.T) {
	_, err := LoadFrom(filepath.Join(t.TempDir(), "nonexistent_raylib.dll"))
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestLoadSidecarIntegration(t *testing.T) {
	p := os.Getenv("RAYLIB_SO_PATH")
	if p == "" {
		t.Skip("set RAYLIB_SO_PATH to a raylib shared library to test LoadFrom + GetFrameTime binding")
	}
	lib, err := LoadFrom(p)
	if err != nil {
		t.Fatal(err)
	}
	var getFrameTime func() float32
	if err := RegisterGetFrameTime(lib, &getFrameTime); err != nil {
		t.Fatal(err)
	}
	_ = getFrameTime
}
