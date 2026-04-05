package mblight

import (
	"moonbasic/vm/heap"
)

// ShadowCasterHandle returns the heap handle of the light marked for shadow maps, or 0.
func ShadowCasterHandle() heap.Handle {
	shadowMu.Lock()
	defer shadowMu.Unlock()
	return shadowCasterHandle
}

// LightDirection returns the normalized light travel direction for a light handle.
func LightDirection(hs *heap.Store, h heap.Handle) (x, y, z float32, ok bool) {
	if hs == nil || h == 0 {
		return 0, -1, 0, false
	}
	o, err := heap.Cast[*lightObj](hs, h)
	if err != nil || !o.enabled {
		return 0, -1, 0, false
	}
	return o.dirX, o.dirY, o.dirZ, true
}

// LightDiffuse scales RGB by intensity for shading (PBR sun color).
func LightDiffuse(hs *heap.Store, h heap.Handle) (r, g, b float32) {
	if hs == nil || h == 0 {
		return 1, 1, 1
	}
	o, err := heap.Cast[*lightObj](hs, h)
	if err != nil || !o.enabled {
		return 0, 0, 0
	}
	return o.r * o.intensity * o.colA, o.g * o.intensity * o.colA, o.b * o.intensity * o.colA
}

// LightShadowTarget returns the world-space point the shadow ortho camera looks at.
func LightShadowTarget(hs *heap.Store, h heap.Handle) (x, y, z float32, ok bool) {
	if hs == nil || h == 0 {
		return 0, 2, 0, false
	}
	o, err := heap.Cast[*lightObj](hs, h)
	if err != nil {
		return 0, 2, 0, false
	}
	return o.targetX, o.targetY, o.targetZ, true
}

// LightShadowBiasK returns the shadow depth bias multiplier for PBR sampling.
func LightShadowBiasK(hs *heap.Store, h heap.Handle) float32 {
	if hs == nil || h == 0 {
		return 1
	}
	o, err := heap.Cast[*lightObj](hs, h)
	if err != nil {
		return 1
	}
	if o.shadowBiasK <= 0 {
		return 1
	}
	return o.shadowBiasK
}
