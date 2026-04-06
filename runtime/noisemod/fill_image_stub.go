//go:build !cgo && !windows

package noisemod

import (
	"fmt"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

func (m *Module) noiseFillImage(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = m
	_ = rt
	_ = args
	return value.Nil, fmt.Errorf("NOISE.FILLIMAGE requires CGO (Raylib image)")
}
