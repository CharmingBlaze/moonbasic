//go:build cgo || (windows && !cgo)

package mblight

import (
	"sync"

	"moonbasic/vm/heap"
)

var (
	followMu          sync.Mutex
	pointFollowLights []heap.Handle
	worldPosGetter    func(int64) (float32, float32, float32, bool)
)

// SetLightFollowWorldPosGetter is called from mbentity so point lights parented to entity# get world positions each frame.
func SetLightFollowWorldPosGetter(f func(int64) (float32, float32, float32, bool)) {
	followMu.Lock()
	worldPosGetter = f
	followMu.Unlock()
}

func registerPointFollow(h heap.Handle) {
	followMu.Lock()
	pointFollowLights = append(pointFollowLights, h)
	followMu.Unlock()
}

func unregisterPointFollow(h heap.Handle) {
	followMu.Lock()
	defer followMu.Unlock()
	for i, v := range pointFollowLights {
		if v == h {
			pointFollowLights = append(pointFollowLights[:i], pointFollowLights[i+1:]...)
			return
		}
	}
}

// SyncPointFollowLights updates LIGHT positions for point lights with parentEntID set.
func SyncPointFollowLights(h *heap.Store) {
	if h == nil {
		return
	}
	followMu.Lock()
	get := worldPosGetter
	list := append([]heap.Handle(nil), pointFollowLights...)
	followMu.Unlock()
	if get == nil {
		return
	}
	for _, hid := range list {
		o, err := heap.Cast[*lightObj](h, hid)
		if err != nil {
			continue
		}
		if o.kind != "point" || o.parentEntID < 1 {
			continue
		}
		px, py, pz, ok := get(o.parentEntID)
		if !ok {
			continue
		}
		o.posX, o.posY, o.posZ = px, py, pz
	}
}
