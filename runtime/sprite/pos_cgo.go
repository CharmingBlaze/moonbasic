//go:build cgo || (windows && !cgo)

package mbsprite

import (
	"fmt"

	"moonbasic/vm/heap"
)

// WorldXY returns sprite world position for CAMERA2D.FOLLOW and similar tools.
func WorldXY(s *heap.Store, h heap.Handle) (x, y float32, err error) {
	o, err := heap.Cast[*spriteObj](s, h)
	if err != nil {
		return 0, 0, fmt.Errorf("sprite handle: %w", err)
	}
	return o.x, o.y, nil
}
