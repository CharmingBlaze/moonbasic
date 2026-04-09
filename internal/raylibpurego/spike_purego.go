//go:build darwin || freebsd || linux || windows

package raylibpurego

import (
	"fmt"

	"github.com/ebitengine/purego"
)

// RegisterGetFrameTime binds Raylib GetFrameTime from an opened library handle.
// Symbol name matches raylib C API (no stdcall decoration in export).
func RegisterGetFrameTime(lib *LoadResult, out *func() float32) error {
	if lib == nil || out == nil {
		return fmt.Errorf("raylibpurego: nil argument")
	}
	purego.RegisterLibFunc(out, lib.Handle, "GetFrameTime")
	return nil
}
