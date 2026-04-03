//go:build cgo

package mbaudio

import (
	"fmt"
	"math"

	"moonbasic/runtime"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

func (m *Module) requireHeap() error {
	if m.h == nil {
		return runtime.Errorf("audio builtins: heap not bound")
	}
	return nil
}

func argInt32(v value.Value) (int32, bool) {
	if i, ok := v.ToInt(); ok {
		return int32(i), true
	}
	if f, ok := v.ToFloat(); ok {
		return int32(f), true
	}
	return 0, false
}

func argUInt32(v value.Value) (uint32, bool) {
	i, ok := argInt32(v)
	if !ok || i < 0 {
		return 0, false
	}
	return uint32(i), true
}

func argFloat32(v value.Value) (float32, bool) {
	if f, ok := v.ToFloat(); ok {
		return float32(f), true
	}
	if i, ok := v.ToInt(); ok {
		return float32(i), true
	}
	return 0, false
}

// pcmSliceFrom1DArray builds a buffer for UpdateAudioStream from a 1-D heap array.
// sampleBits is 8, 16, or 32 per the stream's SampleSize.
func pcmSliceFrom1DArray(arr *heap.Array, sampleBits uint32, op string) (any, error) {
	if len(arr.Dims) != 1 {
		return nil, fmt.Errorf("%s: PCM data must be a 1-D array", op)
	}
	if arr.Kind != heap.ArrayKindFloat {
		return nil, fmt.Errorf("%s: PCM data must be a numeric (#) array", op)
	}
	n := len(arr.Floats)
	switch sampleBits {
	case 32:
		buf := make([]float32, n)
		for i, f := range arr.Floats {
			buf[i] = float32(f)
		}
		return buf, nil
	case 16:
		buf := make([]int16, n)
		for i, f := range arr.Floats {
			x := f * 32767.0
			if x > 32767 {
				x = 32767
			}
			if x < -32768 {
				x = -32768
			}
			buf[i] = int16(math.Round(x))
		}
		return buf, nil
	case 8:
		return nil, fmt.Errorf("%s: 8-bit streams are not supported for array upload; use sampleRate with bitDepth 16 or 32", op)
	default:
		return nil, fmt.Errorf("%s: unsupported stream bit depth %d (use 8, 16, or 32)", op, sampleBits)
	}
}
