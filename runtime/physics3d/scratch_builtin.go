package mbphysics3d

import (
	"fmt"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

// phGetScratchFloat exposes VM.PhysicsScratch[index] after WASM sync (see PHYSICS3D.SYNCWASMTOPHYSREGS / SyncWasmPhysicsAfterStep).
func (m *Module) phGetScratchFloat(args []value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("PHYSICS3D.GETSCRATCHFLOAT expects 1 argument (index)")
	}
	ii, ok := args[0].ToInt()
	if !ok {
		if f, okf := args[0].ToFloat(); okf {
			ii = int64(f)
			ok = true
		}
	}
	if !ok {
		return value.Nil, fmt.Errorf("PHYSICS3D.GETSCRATCHFLOAT: index must be numeric")
	}
	idx := int(ii)
	if idx < 0 {
		return value.Nil, fmt.Errorf("PHYSICS3D.GETSCRATCHFLOAT: index must be non-negative")
	}
	m.vmMu.Lock()
	v := m.vmRef
	m.vmMu.Unlock()
	if v == nil {
		return value.Nil, runtime.Errorf("PHYSICS3D.GETSCRATCHFLOAT: VM not bound")
	}
	if idx >= len(v.PhysicsScratch) {
		return value.FromFloat(0), nil
	}
	return value.FromFloat(v.PhysicsScratch[idx]), nil
}
