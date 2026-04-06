//go:build !cgo && !windows

package mbcamera

import (
	"fmt"

	"moonbasic/vm/heap"
)

// ThirdPersonFollowStep is a no-op stub when Raylib camera natives are unavailable.
func ThirdPersonFollowStep(_ *heap.Store, _ heap.Handle, _, _, _, _, _, _, _ float32, _ float64) error {
	return fmt.Errorf("CAMERA.FOLLOW requires CGO: set CGO_ENABLED=1 and install a C compiler, then rebuild")
}
