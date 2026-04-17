package mbgame

import (
	"fmt"
	"math"

	"moonbasic/runtime"
	"moonbasic/runtime/mbmatrix"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

// player2DObj holds XZ position and optional axis-aligned bounds for Blitz-style helpers.
// Ownership: GAME / PLAYER2D.* only — free with PLAYER2D.FREE or ERASE.
type player2DObj struct {
	X, Z      float64
	MinX      float64
	MaxX      float64
	MinZ      float64
	MaxZ      float64
	hasBounds bool
	release   heap.ReleaseOnce
}

func (o *player2DObj) Free()            { o.release.Do(func() {}) }
func (o *player2DObj) TypeName() string { return "Player2D" }
func (o *player2DObj) TypeTag() uint16  { return heap.TagPlayer2D }

func (m *Module) registerPlayer2D(r runtime.Registrar) {
	r.Register("PLAYER2D.CREATE", "game", m.player2DMake)
	r.Register("PLAYER2D.MAKE", "game", m.player2DMake)
	r.Register("PLAYER2D.FREE", "game", m.player2DFree)
	r.Register("PLAYER2D.MOVE", "game", m.player2DMove)
	r.Register("PLAYER2D.CLAMP", "game", m.player2DClamp)
	r.Register("PLAYER2D.KEEPINBOUNDS", "game", m.player2DKeepInBounds)
	r.Register("PLAYER2D.GETX", "game", m.player2DGetX)
	r.Register("PLAYER2D.GETZ", "game", m.player2DGetZ)
	r.Register("PLAYER2D.GETPOS", "game", m.player2DGetPos)
	r.Register("PLAYER2D.SETPOS", "game", m.player2DSetPos)
	// English / Blitz-style flat names (same behavior as PLAYER2D.*)
	r.Register("MOVEENTITY2D", "game", m.player2DMove)
	r.Register("CLAMPENTITY2D", "game", m.player2DClamp)
	r.Register("MOVEPLAYER", "game", m.player2DMove)
	r.Register("KEEPPLAYERINBOUNDS", "game", m.player2DKeepInBounds)
}

func (m *Module) player2DMake(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("PLAYER2D.MAKE expects 0 arguments")
	}
	if m.h == nil {
		return value.Nil, fmt.Errorf("PLAYER2D.MAKE: heap not bound")
	}
	id, err := m.h.Alloc(&player2DObj{})
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}

func (m *Module) player2DFree(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("PLAYER2D.FREE expects (handle)")
	}
	if _, err := heap.Cast[*player2DObj](m.h, heap.Handle(args[0].IVal)); err != nil {
		return value.Nil, err
	}
	return value.Nil, rt.Heap.Free(heap.Handle(args[0].IVal))
}

func argPlayer2D(h *heap.Store, v value.Value) (*player2DObj, error) {
	if v.Kind != value.KindHandle {
		return nil, fmt.Errorf("expected Player2D handle")
	}
	o, err := heap.Cast[*player2DObj](h, heap.Handle(v.IVal))
	if err != nil {
		return nil, err
	}
	return o, nil
}

func (m *Module) player2DMove(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) != 6 {
		return value.Nil, fmt.Errorf("PLAYER2D.MOVE expects (player, camYaw, forward, strafe, speed, dt)")
	}
	o, err := argPlayer2D(m.h, args[0])
	if err != nil {
		return value.Nil, err
	}
	yaw, ok1 := args[1].ToFloat()
	fwd, ok2 := args[2].ToFloat()
	sf, ok3 := args[3].ToFloat()
	spd, ok4 := args[4].ToFloat()
	dt, ok5 := args[5].ToFloat()
	if !ok1 || !ok2 || !ok3 || !ok4 || !ok5 {
		return value.Nil, fmt.Errorf("PLAYER2D.MOVE: numeric arguments required")
	}
	dx := (-math.Sin(yaw)*fwd + math.Cos(yaw)*sf) * spd * dt
	dz := (-math.Cos(yaw)*fwd + (-math.Sin(yaw))*sf) * spd * dt
	o.X += dx
	o.Z += dz
	return value.Nil, nil
}

func (m *Module) player2DClamp(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) != 5 {
		return value.Nil, fmt.Errorf("PLAYER2D.CLAMP expects (player, minX, maxX, minZ, maxZ)")
	}
	o, err := argPlayer2D(m.h, args[0])
	if err != nil {
		return value.Nil, err
	}
	minX, ok1 := args[1].ToFloat()
	maxX, ok2 := args[2].ToFloat()
	minZ, ok3 := args[3].ToFloat()
	maxZ, ok4 := args[4].ToFloat()
	if !ok1 || !ok2 || !ok3 || !ok4 {
		return value.Nil, fmt.Errorf("PLAYER2D.CLAMP: bounds must be numeric")
	}
	o.MinX, o.MaxX = minX, maxX
	o.MinZ, o.MaxZ = minZ, maxZ
	o.hasBounds = true
	o.X = clampF64(o.X, minX, maxX)
	o.Z = clampF64(o.Z, minZ, maxZ)
	return value.Nil, nil
}

func (m *Module) player2DKeepInBounds(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("PLAYER2D.KEEPINBOUNDS expects (player)")
	}
	o, err := argPlayer2D(m.h, args[0])
	if err != nil {
		return value.Nil, err
	}
	if !o.hasBounds {
		return value.Nil, nil
	}
	o.X = clampF64(o.X, o.MinX, o.MaxX)
	o.Z = clampF64(o.Z, o.MinZ, o.MaxZ)
	return value.Nil, nil
}

func clampF64(v, lo, hi float64) float64 {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

func (m *Module) player2DGetX(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("PLAYER2D.GETX expects (player)")
	}
	o, err := argPlayer2D(m.h, args[0])
	if err != nil {
		return value.Nil, err
	}
	return value.FromFloat(o.X), nil
}

func (m *Module) player2DGetZ(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("PLAYER2D.GETZ expects (player)")
	}
	o, err := argPlayer2D(m.h, args[0])
	if err != nil {
		return value.Nil, err
	}
	return value.FromFloat(o.Z), nil
}

func (m *Module) player2DGetPos(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("PLAYER2D.GETPOS expects (player)")
	}
	o, err := argPlayer2D(m.h, args[0])
	if err != nil {
		return value.Nil, err
	}
	return mbmatrix.AllocVec2Value(m.h, float32(o.X), float32(o.Z))
}

func (m *Module) player2DSetPos(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("PLAYER2D.SETPOS expects (player, x, z)")
	}
	o, err := argPlayer2D(m.h, args[0])
	if err != nil {
		return value.Nil, err
	}
	x, ok1 := args[1].ToFloat()
	z, ok2 := args[2].ToFloat()
	if !ok1 || !ok2 {
		return value.Nil, fmt.Errorf("PLAYER2D.SETPOS: x and z must be numeric")
	}
	o.X, o.Z = x, z
	return value.Nil, nil
}
