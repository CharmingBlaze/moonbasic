package raylibpurego

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

// LibPath returns the conventional sidecar name for the current OS (for documentation/tests).
func LibPath() string {
	switch runtime.GOOS {
	case "windows":
		return "raylib.dll"
	case "darwin":
		return "libraylib.dylib"
	default:
		return "libraylib.so"
	}
}

// LoadResult holds a successfully opened dynamic library handle from purego.
// Callers use Register* helpers to bind symbols.
type LoadResult struct {
	Handle uintptr
}

// Load attempts to dlopen the Raylib shared library from baseName next to the
// running binary (first path tried). Used by spikes and future bindings.
func Load(baseName string) (*LoadResult, error) {
	if baseName == "" {
		baseName = LibPath()
	}
	exe, err := executableDir()
	if err != nil {
		return nil, err
	}
	return LoadFrom(filepath.Join(exe, baseName))
}

// LoadFrom opens an explicit path (tests or RAYLIB_SO_PATH override).
func LoadFrom(dllPath string) (*LoadResult, error) {
	handle, err := loadDynamicLibrary(dllPath)
	if err != nil {
		return nil, fmt.Errorf("raylibpurego: open %q: %w", dllPath, err)
	}
	return &LoadResult{Handle: handle}, nil
}

func executableDir() (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.Dir(exe), nil
}
