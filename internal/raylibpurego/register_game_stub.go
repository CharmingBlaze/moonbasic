//go:build !(darwin || freebsd || linux || windows)

package raylibpurego

import "fmt"

// RegisterGame is unavailable: ebitengine/purego only provides RegisterLibFunc on darwin, freebsd, linux, and windows.
func RegisterGame(lib *LoadResult, g *Game) error {
	return fmt.Errorf("raylibpurego: RegisterGame is not supported on this GOOS (purego RegisterLibFunc unavailable)")
}
