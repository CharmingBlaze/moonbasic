package mathmod

import (
	"fmt"
	"math"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

func (m *Module) registerMovement(r runtime.Registrar) {
	r.Register("MOVEX", "math", m.moveX)
	r.Register("MOVEZ", "math", m.moveZ)
	r.Register("MATH.NEWX", "math", m.mathNewX)
	r.Register("MATH.NEWZ", "math", m.mathNewZ)
	// Bundles MOVEX/MOVEZ with * speed * dt (same as MOVEX(...)*speed*dt).
	r.Register("MOVESTEPX", "math", m.moveStepX)
	r.Register("MOVESTEPZ", "math", m.moveStepZ)
}

// mathNewX / mathNewZ: world position after moving `distance` along the forward axis at `yaw`
// (radians), matching MOVEX(yaw,1,0) and MOVEZ(yaw,1,0) — see MATH.md.
func (m *Module) mathNewX(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = m
	_ = rt
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("MATH.NEWX expects 3 arguments (currentX#, yawRadians#, distance#)")
	}
	cx, ok1 := args[0].ToFloat()
	yaw, ok2 := args[1].ToFloat()
	dist, ok3 := args[2].ToFloat()
	if !ok1 || !ok2 || !ok3 {
		return value.Nil, fmt.Errorf("MATH.NEWX: numeric arguments required")
	}
	nx := cx + (-math.Sin(yaw))*dist
	return value.FromFloat(nx), nil
}

func (m *Module) mathNewZ(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = m
	_ = rt
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("MATH.NEWZ expects 3 arguments (currentZ#, yawRadians#, distance#)")
	}
	cz, ok1 := args[0].ToFloat()
	yaw, ok2 := args[1].ToFloat()
	dist, ok3 := args[2].ToFloat()
	if !ok1 || !ok2 || !ok3 {
		return value.Nil, fmt.Errorf("MATH.NEWZ: numeric arguments required")
	}
	nz := cz + (-math.Cos(yaw))*dist
	return value.FromFloat(nz), nil
}

func (m *Module) moveX(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("MOVEX expects 3 arguments (yaw#, forward#, strafe#)")
	}
	yaw, ok1 := args[0].ToFloat()
	fwd, ok2 := args[1].ToFloat()
	sf, ok3 := args[2].ToFloat()
	if !ok1 || !ok2 || !ok3 {
		return value.Nil, fmt.Errorf("MOVEX: numeric arguments required")
	}
	x := -math.Sin(yaw)*fwd + math.Cos(yaw)*sf
	return value.FromFloat(x), nil
}

func (m *Module) moveZ(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("MOVEZ expects 3 arguments (yaw#, forward#, strafe#)")
	}
	yaw, ok1 := args[0].ToFloat()
	fwd, ok2 := args[1].ToFloat()
	sf, ok3 := args[2].ToFloat()
	if !ok1 || !ok2 || !ok3 {
		return value.Nil, fmt.Errorf("MOVEZ: numeric arguments required")
	}
	z := -math.Cos(yaw)*fwd + (-math.Sin(yaw))*sf
	return value.FromFloat(z), nil
}

func (m *Module) moveStepX(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = m
	if len(args) != 5 {
		return value.Nil, fmt.Errorf("MOVESTEPX expects 5 arguments (yaw#, forward#, strafe#, speed#, dt#)")
	}
	yaw, ok1 := args[0].ToFloat()
	fwd, ok2 := args[1].ToFloat()
	sf, ok3 := args[2].ToFloat()
	spd, ok4 := args[3].ToFloat()
	dt, ok5 := args[4].ToFloat()
	if !ok1 || !ok2 || !ok3 || !ok4 || !ok5 {
		return value.Nil, fmt.Errorf("MOVESTEPX: numeric arguments required")
	}
	x := (-math.Sin(yaw)*fwd + math.Cos(yaw)*sf) * spd * dt
	return value.FromFloat(x), nil
}

func (m *Module) moveStepZ(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = m
	if len(args) != 5 {
		return value.Nil, fmt.Errorf("MOVESTEPZ expects 5 arguments (yaw#, forward#, strafe#, speed#, dt#)")
	}
	yaw, ok1 := args[0].ToFloat()
	fwd, ok2 := args[1].ToFloat()
	sf, ok3 := args[2].ToFloat()
	spd, ok4 := args[3].ToFloat()
	dt, ok5 := args[4].ToFloat()
	if !ok1 || !ok2 || !ok3 || !ok4 || !ok5 {
		return value.Nil, fmt.Errorf("MOVESTEPZ: numeric arguments required")
	}
	z := (-math.Cos(yaw)*fwd + (-math.Sin(yaw))*sf) * spd * dt
	return value.FromFloat(z), nil
}
