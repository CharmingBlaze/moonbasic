//go:build cgo || (windows && !cgo)

package water

import (
	"math"

	"moonbasic/vm/heap"
)

// EntitySubmergedFraction returns approximately what fraction 0..1 of the entity's vertical extent [mnY,mxY]
// lies below the water surface, for water volumes whose XZ footprint contains (cx,cz).
func EntitySubmergedFraction(h *heap.Store, mnY, mxY, cx, cz float32) float32 {
	if h == nil || mxY <= mnY {
		return 0
	}
	height := mxY - mnY
	var best float32
	for _, wh := range h.FilterByType(heap.TagWater) {
		wo, err := heap.Cast[*WaterObject](h, wh)
		if err != nil {
			continue
		}
		hx0 := wo.PX - wo.Width*0.5
		hx1 := wo.PX + wo.Width*0.5
		hz0 := wo.PZ - wo.Depth*0.5
		hz1 := wo.PZ + wo.Depth*0.5
		if cx < hx0 || cx > hx1 || cz < hz0 || cz > hz1 {
			continue
		}
		bob := float32(math.Sin(float64(wo.WaveT))) * wo.WaveAmp * 0.15
		surf := wo.PY + bob
		top := float32(math.Min(float64(mxY), float64(surf)))
		bot := float32(math.Max(float64(mnY), float64(wo.BedY)))
		ov := top - bot
		if ov <= 0 {
			continue
		}
		frac := ov / height
		if frac > best {
			best = frac
		}
	}
	if best > 1 {
		best = 1
	}
	return best
}
