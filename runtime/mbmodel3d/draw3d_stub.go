//go:build !cgo && !windows

package mbmodel3d

import "moonbasic/vm/heap"

// SetGlobalHeapGetter is a no-op without CGO raylib.
func SetGlobalHeapGetter(fn func() *heap.Store) {}

func MarkCamera3DBegin(camX, camY, camZ float32) {}

func MarkCamera3DEnd() {}

func FlushDeferred3D(_ *heap.Store) {}

func InCamera3D() bool { return false }

func shadowDeferActive() bool { return false }
