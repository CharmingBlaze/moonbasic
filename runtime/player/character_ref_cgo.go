//go:build (linux || windows) && cgo

package player

import (
	"fmt"
	"math"

	mbcamera "moonbasic/runtime/camera"
	mbmatrix "moonbasic/runtime/mbmatrix"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

func (m *Module) playerCharacterCreate(args []value.Value) (value.Value, error) {
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

func (m *Module) charRefSetVel(args []value.Value) (value.Value, error) {
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

func (m *Module) charRefUpdate(args []value.Value) (value.Value, error) {
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

func (m *Module) charRefJump(args []value.Value) (value.Value, error) {
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

func (m *Module) charRefMoveWithCam(args []value.Value) (value.Value, error) {
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

	// CharacterMoveXZVelocity handles its own integration if we pass a dt, but typically called per-frame.
	// We'll use a fixed dt or fetch from time if we had rt.
	dt := 1.0 / 60.0
	if err := m.char.CharacterMoveXZVelocity(ch, vx, vz, dt); err != nil {
		return value.Nil, err
	}
	x, y, z, ok := m.char.CharacterPosition(ch)
	if ok {
		_ = m.ent.PlayerBridgeSetWorldPos(obj.id, float32(x), float32(y), float32(z))
	}
	return value.Nil, nil
}

func (m *Module) charRefSetMaxSlope(args []value.Value) (value.Value, error) {
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

func (m *Module) charRefSetStepHeight(args []value.Value) (value.Value, error) {
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

func (m *Module) charRefIsGrounded(args []value.Value) (value.Value, error) {
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

func (m *Module) charRefSetPos(args []value.Value) (value.Value, error) {
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

func (m *Module) charRefGetPos(args []value.Value) (value.Value, error) {
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

func (m *Module) charRefFree(args []value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("CHARACTERREF.FREE expects (handle)")
	}
	h := heap.Handle(args[0].IVal)
	if err := m.h.Free(h); err != nil {
		return value.Nil, err
	}
	return value.Nil, nil
}

func (m *Module) charRefGetGroundState(args []value.Value) (value.Value, error) {
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

func (m *Module) charRefAddVel(args []value.Value) (value.Value, error) {
	if len(args) != 4 {
		return value.Nil, fmt.Errorf("CHARACTERREF.ADDVELOCITY expects (handle, vx#, vy#, vz#)")
	}
	obj, err := heap.Cast[*charRefHeapObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	ch, ok := m.entToChar[obj.id]
	if !ok || m.char == nil {
		return value.Nil, fmt.Errorf("CHARACTER: no KCC for entity")
	}
	vx, _ := args[1].ToFloat()
	vy, _ := args[2].ToFloat()
	vz, _ := args[3].ToFloat()
	cvx, cvy, cvz, _ := m.char.CharacterLinearVelocity(ch)
	if err := m.char.SetCharacterLinearVelocity(ch, cvx+vx, cvy+vy, cvz+vz); err != nil {
		return value.Nil, err
	}
	return value.Nil, nil
}

func (m *Module) charRefSetSnapDist(args []value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("CHARACTERREF.SETSNAPDISTANCE expects (handle, dist#)")
	}
	obj, err := heap.Cast[*charRefHeapObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	ch, ok := m.entToChar[obj.id]
	if !ok || m.char == nil {
		return value.Nil, fmt.Errorf("CHARACTER: no KCC for entity")
	}
	d, _ := args[1].ToFloat()
	if err := m.char.SetCharacterStickToFloorDown(ch, float32(d)); err != nil {
		return value.Nil, err
	}
	return value.Nil, nil
}

func (m *Module) charRefGetSpeed(args []value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("CHARACTERREF.GETSPEED expects (handle)")
	}
	obj, err := heap.Cast[*charRefHeapObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	ch, ok := m.entToChar[obj.id]
	if !ok || m.char == nil {
		return value.FromFloat(0), nil
	}
	vx, vy, vz, ok := m.char.CharacterLinearVelocity(ch)
	if !ok {
		return value.FromFloat(0), nil
	}
	sp := math.Sqrt(vx*vx + vy*vy + vz*vz)
	return value.FromFloat(sp), nil
}

func (m *Module) charRefIsMoving(args []value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("CHARACTERREF.ISMOVING expects (handle)")
	}
	obj, err := heap.Cast[*charRefHeapObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	ch, ok := m.entToChar[obj.id]
	if !ok || m.char == nil {
		return value.FromBool(false), nil
	}
	vx, _, vz, ok := m.char.CharacterLinearVelocity(ch)
	if !ok {
		return value.FromBool(false), nil
	}
	hs := math.Hypot(vx, vz)
	return value.FromBool(hs > 0.05), nil
}

func (m *Module) charRefSetFriction(args []value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("CHARACTERREF.SETFRICTION expects (handle, friction#)")
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
	if err := m.char.SetCharacterFriction(ch, f); err != nil {
		return value.Nil, err
	}
	return value.Nil, nil
}

func (m *Module) charRefSetPadding(args []value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("CHARACTERREF.SETPADDING expects (handle, padding#)")
	}
	obj, err := heap.Cast[*charRefHeapObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	ch, ok := m.entToChar[obj.id]
	if !ok || m.char == nil {
		return value.Nil, fmt.Errorf("CHARACTER: no KCC for entity")
	}
	p, _ := args[1].ToFloat()
	newH, err := m.char.SetCharacterPadding(ch, float32(p))
	if err != nil {
		return value.Nil, err
	}
	m.entToChar[obj.id] = newH
	return value.Nil, nil
}

func (m *Module) charRefSetBounce(args []value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("CHARACTERREF.SETBOUNCE expects (handle, bounce#)")
	}
	obj, err := heap.Cast[*charRefHeapObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	ch, ok := m.entToChar[obj.id]
	if !ok || m.char == nil {
		return value.Nil, fmt.Errorf("CHARACTER: no KCC for entity")
	}
	b, _ := args[1].ToFloat()
	if err := m.char.SetCharacterRestitution(ch, b); err != nil {
		return value.Nil, err
	}
	return value.Nil, nil
}

func (m *Module) charRefSetGravityScale(args []value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("CHARACTERREF.SETGRAVITYSCALE expects (handle, scale#)")
	}
	obj, err := heap.Cast[*charRefHeapObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	ch, ok := m.entToChar[obj.id]
	if !ok || m.char == nil {
		return value.Nil, fmt.Errorf("CHARACTER: no KCC for entity")
	}
	s, _ := args[1].ToFloat()
	if err := m.char.SetCharacterGravityScale(ch, s); err != nil {
		return value.Nil, err
	}
	return value.Nil, nil
}
