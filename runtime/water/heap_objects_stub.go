//go:build !cgo && !windows

package water

import "moonbasic/vm/heap"

// WaterObject stub — WATER.* is unavailable without CGO on this platform.
type WaterObject struct {
	Width    float32
	Depth    float32
	PX, PY, PZ float32
	WaveT    float32
	WaveAmp  float32
	WaveFreq float32
	BedY     float32
	freed    bool
}

func (w *WaterObject) TypeName() string { return "Water" }
func (w *WaterObject) TypeTag() uint16  { return heap.TagWater }

func (w *WaterObject) Free() {
	if w.freed {
		return
	}
	w.freed = true
}
