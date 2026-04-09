//go:build !cgo && !windows

package mblight

import "moonbasic/vm/heap"

func unregisterPointFollow(_ heap.Handle) {}

func SetLightFollowWorldPosGetter(_ func(int64) (float32, float32, float32, bool)) {}

func SyncPointFollowLights(_ *heap.Store) {}
