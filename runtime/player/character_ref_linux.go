//go:build linux && cgo

package player

import (
	"fmt"

	mbcamera "moonbasic/runtime/camera"
	mbmatrix "moonbasic/runtime/mbmatrix"
	mbtime "moonbasic/runtime/time"
	"moonbasic/runtime"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

func registerCharacterRefCommands(m *Module, reg runtime.Registrar) {
	reg.Register("CHARACTER.CREATE", "player", m.playerCharacterCreate)
	reg.Register("CHARACTERREF.SETLINEARVELOCITY", "player", m.charRefSetVel)
	reg.Register("CHARACTERREF.UPDATE", "player", m.charRefUpdate)
	reg.Register("CHARACTERREF.UPDATEMOVE", "player", m.charRefUpdate)
	reg.Register("CHARACTERREF.JUMP", "player", m.charRefJump)
	reg.Register("CHARACTERREF.MOVEWITHCAMERA", "player", m.charRefMoveWithCam)
	reg.Register("CHARACTERREF.SETMAXSLOPE", "player", m.charRefSetMaxSlope)
	reg.Register("CHARACTERREF.SETSTEPHEIGHT", "player", m.charRefSetStepHeight)
	reg.Register("CHARACTERREF.ISGROUNDED", "player", m.charRefIsGrounded)
	reg.Register("CHARACTERREF.SETPOSITION", "player", m.charRefSetPos)
	reg.Register("CHARACTERREF.GETPOSITION", "player", m.charRefGetPos)
	reg.Register("CHARACTERREF.FREE", "player", m.charRefFree)
	reg.Register("CHARACTERREF.GETGROUNDSTATE", "player", m.charRefGetGroundState)
}

func (m *Module) playerCharacterCreate(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if m.h == nil {
		return value.Nil, fmt.Errorf("CHARACTER.CREATE: heap not bound")
	}
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("CHARACTER.CREATE expects (entity, radius#, height#)")
	}
	id, err := m.playerCreateInternal(args)
	if err != nil {
		return value.Nil, err
	}
	obj := &charRefHeapObj{id: id, m: m}
	hh, err := m.h.Alloc(obj)
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(hh), nil
}

func (m *Module) charRefSetVel(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) != 4 {
		return value.Nil, fmt.Errorf("CHARACTERREF.SETLINEARVELOCITY expects (handle, vx#, vy#, vz#)")
	}
	obj, err := heap.Cast[*charRefHeapObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	ch, ok := m.entToChar[obj.id]
	if !ok || m.char == nil {
		return value.Nil, fmt.Errorf("CHARACTER: call CHARACTER.CREATE first")
	}
	vx, _ := args[1].ToFloat()
	vy, _ := args[2].ToFloat()
	vz, _ := args[3].ToFloat()
	if err := m.char.SetCharacterLinearVelocity(ch, vx, vy, vz); err != nil {
		return value.Nil, err
	}
	return value.Nil, nil
}

func (m *Module) charRefUpdate(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("CHARACTERREF.UPDATE expects (handle, dt#)")
	}
	obj, err := heap.Cast[*charRefHeapObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	ch, ok := m.entToChar[obj.id]
	if !ok || m.char == nil {
		return value.Nil, fmt.Errorf("CHARACTER: no KCC for entity")
	}
	dt, _ := args[1].ToFloat()
	if err := m.char.CharacterIntegrateStep(ch, dt); err != nil {
		return value.Nil, err
	}
	if m.ent != nil {
		x, y, z, ok := m.char.CharacterPosition(ch)
		if ok {
			_ = m.ent.PlayerBridgeSetWorldPos(obj.id, float32(x), float32(y), float32(z))
		}
	}
	return value.Nil, nil
}

func (m *Module) charRefJump(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("CHARACTERREF.JUMP expects (handle, vy#)")
	}
	obj, err := heap.Cast[*charRefHeapObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	ch, ok := m.entToChar[obj.id]
	if !ok || m.char == nil {
		return value.Nil, fmt.Errorf("CHARACTER: no KCC for entity")
	}
	f, _ := args[1].ToFloat()
	vx, _, vz, ok := m.char.CharacterLinearVelocity(ch)
	if !ok {
		return value.Nil, fmt.Errorf("CHARACTERREF.JUMP: internal")
	}
	if err := m.char.SetCharacterLinearVelocity(ch, vx, f, vz); err != nil {
		return value.Nil, err
	}
	return value.Nil, nil
}

func (m *Module) charRefMoveWithCam(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if m.h == nil || m.char == nil || m.ent == nil {
		return value.Nil, fmt.Errorf("CHARACTERREF.MOVEWITHCAMERA: not available")
	}
	if len(args) != 5 {
		return value.Nil, fmt.Errorf("CHARACTERREF.MOVEWITHCAMERA expects (handle, camera, forwardAxis#, strafeAxis#, speed#)")
	}
	obj, err := heap.Cast[*charRefHeapObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	if args[1].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("CHARACTERREF.MOVEWITHCAMERA: camera handle required")
	}
	ch, ok := m.entToChar[obj.id]
	if !ok {
		return value.Nil, fmt.Errorf("CHARACTER: no KCC for entity")
	}
	f, ok1 := args[2].ToFloat()
	s, ok2 := args[3].ToFloat()
	spd, ok3 := args[4].ToFloat()
	if !ok1 || !ok2 || !ok3 {
		return value.Nil, fmt.Errorf("CHARACTERREF.MOVEWITHCAMERA: forward/strafe/speed must be numeric")
	}
	if spd < 0 {
		return value.Nil, fmt.Errorf("CHARACTERREF.MOVEWITHCAMERA: speed must be non-negative")
	}
	camH := heap.Handle(args[1].IVal)
	fwd, right, err := mbcamera.CameraXZWalkBasis(m.h, camH)
	if err != nil {
		return value.Nil, err
	}
	vx := (float64(fwd.X)*f + float64(right.X)*s) * spd
	vz := (float64(fwd.Z)*f + float64(right.Z)*s) * spd
	dt := mbtime.DeltaSeconds(rt)
	if dt <= 0 {
		dt = 1.0 / 60.0
	}
	if err := m.char.CharacterMoveXZVelocity(ch, vx, vz, dt); err != nil {
		return value.Nil, err
	}
	x, y, z, ok := m.char.CharacterPosition(ch)
	if ok {
		_ = m.ent.PlayerBridgeSetWorldPos(obj.id, float32(x), float32(y), float32(z))
	}
	return value.Nil, nil
}

func (m *Module) charRefSetMaxSlope(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("CHARACTERREF.SETMAXSLOPE expects (handle, maxSlopeDegrees#)")
	}
	obj, err := heap.Cast[*charRefHeapObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	deg, ok := args[1].ToFloat()
	if !ok || deg <= 0 || deg >= 90 {
		return value.Nil, fmt.Errorf("CHARACTERREF.SETMAXSLOPE: angle must be in (0, 90) degrees")
	}
	ch, ok := m.entToChar[obj.id]
	if !ok || m.char == nil {
		return value.Nil, fmt.Errorf("CHARACTER: no KCC for entity")
	}
	rad, fh := float64(playerCapsuleRadius), float64(playerCapsuleHeight)
	if cr, chh, ok := m.char.CharacterCapsuleDims(ch); ok {
		rad, fh = cr, chh
	}
	newH, err := m.char.RecreateCharacterWithSlope(ch, rad, fh, deg)
	if err != nil {
		return value.Nil, err
	}
	m.entToChar[obj.id] = newH
	return value.Nil, nil
}

func (m *Module) charRefSetStepHeight(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("CHARACTERREF.SETSTEPHEIGHT expects (handle, height#)")
	}
	obj, err := heap.Cast[*charRefHeapObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	h, ok := args[1].ToFloat()
	if !ok || h < 0 {
		return value.Nil, fmt.Errorf("CHARACTERREF.SETSTEPHEIGHT: height must be >= 0")
	}
	ch, ok := m.entToChar[obj.id]
	if !ok || m.char == nil {
		return value.Nil, fmt.Errorf("CHARACTER: no KCC for entity")
	}
	m.stepHeight[obj.id] = h
	if err := m.char.SetCharacterWalkStairsStepUp(ch, float32(h)); err != nil {
		return value.Nil, err
	}
	return value.Nil, nil
}

func (m *Module) charRefIsGrounded(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) < 1 {
		return value.Nil, fmt.Errorf("CHARACTERREF.ISGROUNDED expects (handle)")
	}
	obj, err := heap.Cast[*charRefHeapObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	ch, ok := m.entToChar[obj.id]
	if !ok || m.char == nil {
		return value.FromBool(false), nil
	}
	g, err := m.char.CharacterIsGrounded(ch)
	if err != nil {
		return value.Nil, err
	}
	return value.FromBool(g), nil
}

func (m *Module) charRefSetPos(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) != 4 {
		return value.Nil, fmt.Errorf("CHARACTERREF.SETPOSITION expects (handle, x#, y#, z#)")
	}
	obj, err := heap.Cast[*charRefHeapObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	ch, ok := m.entToChar[obj.id]
	if !ok || m.char == nil {
		return value.Nil, fmt.Errorf("CHARACTER: no KCC for entity")
	}
	x, _ := args[1].ToFloat()
	y, _ := args[2].ToFloat()
	z, _ := args[3].ToFloat()
	if err := m.char.CharacterTeleport(ch, x, y, z); err != nil {
		return value.Nil, err
	}
	if m.ent != nil {
		_ = m.ent.PlayerBridgeSetWorldPos(obj.id, float32(x), float32(y), float32(z))
	}
	return value.Nil, nil
}

func (m *Module) charRefGetPos(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if m.h == nil {
		return value.Nil, fmt.Errorf("CHARACTERREF.GETPOSITION: heap not bound")
	}
	obj, err := heap.Cast[*charRefHeapObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	ch, ok := m.entToChar[obj.id]
	if !ok || m.char == nil {
		return value.Nil, nil
	}
	x, y, z, ok := m.char.CharacterPosition(ch)
	if !ok {
		return value.Nil, nil
	}
	return mbmatrix.AllocVec3Value(m.h, float32(x), float32(y), float32(z))
}

func (m *Module) charRefFree(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("CHARACTERREF.FREE expects (handle)")
	}
	h := heap.Handle(args[0].IVal)
	if err := m.h.Free(h); err != nil {
		return value.Nil, err
	}
	return value.Nil, nil
}

func (m *Module) charRefGetGroundState(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("CHARACTERREF.GETGROUNDSTATE expects (handle)")
	}
	obj, err := heap.Cast[*charRefHeapObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	ch, ok := m.entToChar[obj.id]
	if !ok || m.char == nil {
		return value.FromInt(3), nil
	}
	gi, ok := m.char.CharacterGroundStateInt(ch)
	if !ok {
		return value.FromInt(3), nil
	}
	return value.FromInt(int64(gi)), nil
}
