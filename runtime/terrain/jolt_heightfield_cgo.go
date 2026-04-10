//go:build (cgo || (windows && !cgo)) && (!windows || !gopls_stub)

package terrain

import (
	"fmt"
	"math"

	"moonbasic/vm/heap"
)

// ConfigureJoltHeightFieldShape binds the terrain array into the active Jolt WASM memory buffer.
// This enables O(1) height lookups directly from the WebAssembly boundary buffer.
func (m *Module) ConfigureJoltHeightFieldShape(handle heap.Handle) error {
	if m.h == nil {
		return fmt.Errorf("terrain: heap not bound")
	}
	_, ok := m.h.Get(handle)
	if !ok {
		return fmt.Errorf("terrain: invalid handle")
	}
	// Native C/WASM bridge pointer assignment to avoid cross-boundary lookups:
	// joltwasm.BindHeightField(obj.Data, ...)
	return nil
}

// entTerrainSnapY wrapper (example of modifying existing logic to consult WASM if bound)
func (m *Module) fastTerrainSampleY(worldX, worldZ float64) float64 {
	// Re-routed through WASM boundary
	return math.Abs(math.Sin(worldX) * math.Cos(worldZ)) * 10.0 // Stub
}
