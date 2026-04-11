package mbgame

import (
	"fmt"
	"math"

	"moonbasic/runtime"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

// moverFacade forwards PLAYER.MOVERELATIVE / MOVESTEP* / LANDBOXES without storing platform state.
type moverFacade struct{ release heap.ReleaseOnce }

func (o *moverFacade) Free() { o.release.Do(func() {}) }
func (o *moverFacade) TypeTag() uint16  { return heap.TagMoverFacade }
func (o *moverFacade) TypeName() string { return "MOVER" }

func (m *Module) registerMoverFacade(r runtime.Registrar) {
	r.Register("MOVER.MOVEXZ", "game", m.moverMoveXZ)
	r.Register("MOVER.MOVESTEPX", "game", m.moverMoveStepX)
	r.Register("MOVER.MOVESTEPZ", "game", m.moverMoveStepZ)
	r.Register("MOVER.LAND", "game", m.moverLand)
	r.Register("MOVER.MOVEREL", "game", m.moverMoveRelArray)
	r.Register("MOVER.FREE", "game", m.moverFree)
	r.Register("MOVER", "game", runtime.AdaptLegacy(m.makeMover))
}

func (m *Module) makeMover(args []value.Value) (value.Value, error) {
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("MOVER expects 0 arguments")
	}
	if m.h == nil {
		return value.Nil, fmt.Errorf("MOVER: heap not bound")
	}
	id, err := m.h.Alloc(&moverFacade{})
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}

func (m *Module) moverMoveXZ(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) != 6 {
		return value.Nil, fmt.Errorf("MOVER.MOVEXZ expects (handle, yaw, forward, strafe, speed, dt)")
	}
	if _, err := heap.Cast[*moverFacade](m.h, heap.Handle(args[0].IVal)); err != nil {
		return value.Nil, err
	}
	yaw, ok1 := args[1].ToFloat()
	fwd, ok2 := args[2].ToFloat()
	sf, ok3 := args[3].ToFloat()
	spd, ok4 := args[4].ToFloat()
	dt, ok5 := args[5].ToFloat()
	if !ok1 || !ok2 || !ok3 || !ok4 || !ok5 {
		return value.Nil, fmt.Errorf("MOVER.MOVEXZ: numeric required")
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

func (m *Module) moverMoveStepX(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 6 {
		return value.Nil, fmt.Errorf("MOVER.MOVESTEPX expects (handle, yaw, forward, strafe, speed, dt)")
	}
	if _, err := heap.Cast[*moverFacade](m.h, heap.Handle(args[0].IVal)); err != nil {
		return value.Nil, err
	}
	return rt.Call("MOVESTEPX", args[1:6])
}

func (m *Module) moverMoveStepZ(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 6 {
		return value.Nil, fmt.Errorf("MOVER.MOVESTEPZ expects (handle, yaw, forward, strafe, speed, dt)")
	}
	if _, err := heap.Cast[*moverFacade](m.h, heap.Handle(args[0].IVal)); err != nil {
		return value.Nil, err
	}
	return rt.Call("MOVESTEPZ", args[1:6])
}

func (m *Module) moverMoveRelArray(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 6 {
		return value.Nil, fmt.Errorf("MOVER.MOVEREL expects (handle, camYaw, forward, strafe, speed, dt)")
	}
	if _, err := heap.Cast[*moverFacade](m.h, heap.Handle(args[0].IVal)); err != nil {
		return value.Nil, err
	}
	return m.playerMoveRelative(rt, args[1], args[2], args[3], args[4], args[5])
}

func (m *Module) moverLand(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 13 {
		return value.Nil, fmt.Errorf("MOVER.LAND expects handle + 12 arguments (same as LANDBOXES)")
	}
	if _, err := heap.Cast[*moverFacade](m.h, heap.Handle(args[0].IVal)); err != nil {
		return value.Nil, err
	}
	return landBoxes(rt, args[1], args[2], args[3], args[4], args[5], args[6], args[7], args[8], args[9], args[10], args[11], args[12])
}

func (m *Module) moverFree(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("MOVER.FREE expects handle")
	}
	if _, err := heap.Cast[*moverFacade](m.h, heap.Handle(args[0].IVal)); err != nil {
		return value.Nil, err
	}
	return value.Nil, rt.Heap.Free(heap.Handle(args[0].IVal))
}
