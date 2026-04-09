//go:build !(darwin || freebsd || linux || windows)

package raylibpurego

import "fmt"

// RegisterGetFrameTime is unavailable on this GOOS (see register_game_stub.go).
func RegisterGetFrameTime(lib *LoadResult, out *func() float32) error {
	return fmt.Errorf("raylibpurego: RegisterGetFrameTime is not supported on this GOOS")
}
