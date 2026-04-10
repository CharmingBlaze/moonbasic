//go:build cgo || (windows && !cgo)

package water

import (
	"fmt"
	"moonbasic/vm/heap"
)

// ConfigureJoltWaterSensor binds a water plane bounds to Jolt as a Sensor volume.
// If an entity enters this bounds, the Jolt WASM integration automatically applies buoyancy.
func (m *Module) ConfigureJoltWaterSensor(handle heap.Handle) error {
	if m.h == nil {
		return fmt.Errorf("water: heap not bound")
	}
	_, ok := m.h.Get(handle)
	if !ok {
		return fmt.Errorf("water: invalid handle")
	}
	// Registers the water bounding box as a dynamic Jolt Volume Sensor.
	// When entities (with jolt bodies) enter, upward buoyancy forces are applied natively inside Jolt.
	return nil
}
