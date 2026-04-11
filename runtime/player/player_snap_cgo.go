//go:build cgo || (windows && !cgo)

package player

import (
	"fmt"

	"moonbasic/runtime"
	"moonbasic/runtime/terrain"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

func registerPlayerTerrainCommands(m *Module, reg runtime.Registrar) {
	reg.Register("PLAYER.SNAPTOGROUND", "player", m.playerSnapToGround)
	reg.Register("PLAYER.ISSWIMMING", "player", m.playerIsSwimming)
}

func (m *Module) playerSnapToGround(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if m.h == nil || m.ent == nil {
		return value.Nil, fmt.Errorf("PLAYER.SNAPTOGROUND: heap/entity not bound")
	}
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("PLAYER.SNAPTOGROUND expects (entity, terrain, offset)")
	}
	id, ok := args[0].ToInt()
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("PLAYER.SNAPTOGROUND: invalid entity")
	}
	if args[1].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("PLAYER.SNAPTOGROUND: terrain must be a handle")
	}
	th := heap.Handle(args[1].IVal)
	off, err := rt.ArgFloat(args, 2)
	if err != nil {
		return value.Nil, err
	}
	wx, wy, wz, ok := m.ent.PlayerBridgeWorldPos(id)
	if !ok {
		return value.Nil, fmt.Errorf("PLAYER.SNAPTOGROUND: unknown entity")
	}
	_ = wy
	hy := terrain.HeightWorldPublic(m.h, th, float32(wx), float32(wz)) + float32(off)
	_ = m.ent.PlayerBridgeSetWorldPos(id, float32(wx), hy, float32(wz))
	// Jolt capsule sync when available (Linux+Jolt PLAYER.CREATE path).
	playerSnapSyncCharacter(m, id, wx, float64(hy), wz)
	return value.Nil, nil
}

func (m *Module) playerIsSwimming(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if m.h == nil || m.ent == nil {
		return value.FromBool(false), nil
	}
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("PLAYER.ISSWIMMING expects (entity)")
	}
	id, ok := args[0].ToInt()
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("PLAYER.ISSWIMMING: invalid entity")
	}
	x, y, z, ok := m.ent.PlayerBridgeWorldPos(id)
	if !ok {
		return value.Nil, fmt.Errorf("PLAYER.ISSWIMMING: unknown entity")
	}
	if m.water == nil {
		return value.FromBool(false), nil
	}
	return value.FromBool(m.water.PointInWaterVolume(float32(x), float32(y), float32(z))), nil
}
