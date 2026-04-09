package mbcamera

import (
	"fmt"
	"math"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

// camSmoothExp moves current toward target with critically-damped-style exponential smoothing:
// out = current + (target - current) * (1 - exp(-smoothHz * dt)).
// Use for third-person orbit yaw/pitch (see examples/mario64/main_entities.mb).
func (m *Module) camSmoothExp(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) != 4 {
		return value.Nil, fmt.Errorf("CAMERA.SMOOTHEXP expects (current#, target#, smoothHz#, dt#)")
	}
	cur, ok1 := args[0].ToFloat()
	tgt, ok2 := args[1].ToFloat()
	hz, ok3 := args[2].ToFloat()
	dt, ok4 := args[3].ToFloat()
	if !ok1 || !ok2 || !ok3 || !ok4 {
		return value.Nil, fmt.Errorf("CAMERA.SMOOTHEXP: all arguments must be numeric")
	}
	if dt <= 0 {
		return value.FromFloat(cur), nil
	}
	if hz <= 0 {
		return value.FromFloat(tgt), nil
	}
	alpha := 1.0 - math.Exp(-hz*dt)
	if alpha > 1.0 {
		alpha = 1.0
	}
	return value.FromFloat(cur + (tgt-cur)*alpha), nil
}
