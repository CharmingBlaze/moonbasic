package mbgame

import (
	"fmt"
	"math"

	"moonbasic/runtime"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

func (m *Module) registerPlayerHelpers(r runtime.Registrar) {
	r.Register("PLAYER.MOVERELATIVE", "game", m.playerMoveRelative)
}

// playerMoveRelative returns a 2-float array handle [deltaX, deltaZ] for camera-relative
// movement on the XZ plane — same as MOVESTEPX/MOVESTEPZ combined.
func (m *Module) playerMoveRelative(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if m.h == nil {
		return value.Nil, fmt.Errorf("PLAYER.MOVERELATIVE: heap not bound")
	}
	if len(args) != 5 {
		return value.Nil, fmt.Errorf("PLAYER.MOVERELATIVE expects 5 arguments (camYaw#, forward#, strafe#, speed#, dt#)")
	}
	yaw, ok1 := args[0].ToFloat()
	fwd, ok2 := args[1].ToFloat()
	sf, ok3 := args[2].ToFloat()
	spd, ok4 := args[3].ToFloat()
	dt, ok5 := args[4].ToFloat()
	if !ok1 || !ok2 || !ok3 || !ok4 || !ok5 {
		return value.Nil, fmt.Errorf("PLAYER.MOVERELATIVE: numeric arguments required")
	}
	dx := (-math.Sin(yaw)*fwd + math.Cos(yaw)*sf) * spd * dt
	dz := (-math.Cos(yaw)*fwd + (-math.Sin(yaw))*sf) * spd * dt
	arr, err := heap.NewArray([]int64{2})
	if err != nil {
		return value.Nil, err
	}
	_ = arr.Set([]int64{0}, dx)
	_ = arr.Set([]int64{1}, dz)
	id, err := m.h.Alloc(arr)
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}
