//go:build cgo || (windows && !cgo)

package mbcamera

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/vm/heap"
)

// WorldRayFromScreen builds a world-space ray from a camera handle and viewport pixel coordinates.
// Direction components are the **normalized** ray direction; multiply by maxDist for Jolt (length-based cast).
func WorldRayFromScreen(h *heap.Store, camHandle heap.Handle, sx, sy float32) (pos rl.Vector3, dir rl.Vector3, err error) {
	if h == nil {
		return rl.Vector3{}, rl.Vector3{}, fmt.Errorf("WorldRayFromScreen: heap not bound")
	}
	o, e := heap.Cast[*camObj](h, camHandle)
	if e != nil {
		return rl.Vector3{}, rl.Vector3{}, e
	}
	ray := rl.GetScreenToWorldRayEx(rl.Vector2{X: sx, Y: sy}, o.cam, int32(rl.GetRenderWidth()), int32(rl.GetRenderHeight()))
	return ray.Position, ray.Direction, nil
}
