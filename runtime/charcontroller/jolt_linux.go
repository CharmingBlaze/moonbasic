//go:build linux && cgo

package mbcharcontroller

import (
	"fmt"
	"sync"

	"github.com/bbitechnologies/jolt-go/jolt"

	"moonbasic/runtime"
	mbphysics3d "moonbasic/runtime/physics3d"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

// Handles still alive at shutdown must be freed before PHYSICS3D.STOP destroys Jolt.
var charTrackMu sync.Mutex
var charTracked = map[heap.Handle]struct{}{}

func trackChar(h heap.Handle) {
	charTrackMu.Lock()
	charTracked[h] = struct{}{}
	charTrackMu.Unlock()
}

func untrackChar(h heap.Handle) {
	charTrackMu.Lock()
	delete(charTracked, h)
	charTrackMu.Unlock()
}

const charDt = float32(1.0 / 60.0)

type charObj struct {
	cv *jolt.CharacterVirtual
}

func (c *charObj) TypeName() string { return "CharController" }

func (c *charObj) TypeTag() uint16 { return heap.TagCharController }

func (c *charObj) Free() {
	if c.cv != nil {
		c.cv.Destroy()
		c.cv = nil
	}
}

func registerCharControllerCommands(m *Module, reg runtime.Registrar) {
	reg.Register("CHARCONTROLLER.MAKE", "charcontroller", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		return ccMake(m, args)
	}))
	reg.Register("CHARCONTROLLER.SETPOS", "charcontroller", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		return ccSetPos(m, args)
	}))
	reg.Register("CHARCONTROLLER.GETPOS", "charcontroller", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		return ccGetPos(m, args)
	}))
	reg.Register("CHARCONTROLLER.MOVE", "charcontroller", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		return ccMove(m, args)
	}))
	reg.Register("CHARCONTROLLER.ISGROUNDED", "charcontroller", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		return ccIsGrounded(m, args)
	}))
	reg.Register("CHARCONTROLLER.X", "charcontroller", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		return ccAxis(m, args, 0)
	}))
	reg.Register("CHARCONTROLLER.Y", "charcontroller", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		return ccAxis(m, args, 1)
	}))
	reg.Register("CHARCONTROLLER.Z", "charcontroller", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		return ccAxis(m, args, 2)
	}))
	reg.Register("CHARCONTROLLER.FREE", "charcontroller", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		return ccFree(m, args)
	}))
}

func shutdownCharController(m *Module) {
	if m.h == nil {
		return
	}
	charTrackMu.Lock()
	hs := make([]heap.Handle, 0, len(charTracked))
	for h := range charTracked {
		hs = append(hs, h)
	}
	charTracked = make(map[heap.Handle]struct{})
	charTrackMu.Unlock()
	for _, h := range hs {
		m.h.Free(h)
	}
}

func ccMake(m *Module, args []value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, fmt.Errorf("CHARCONTROLLER.MAKE: heap not bound")
	}
	if len(args) < 5 {
		return value.Nil, fmt.Errorf("CHARCONTROLLER.MAKE: need radius, height, x, y, z")
	}
	radius, ok := args[0].ToFloat()
	if !ok {
		return value.Nil, fmt.Errorf("CHARCONTROLLER.MAKE: radius must be numeric")
	}
	height, ok := args[1].ToFloat()
	if !ok {
		return value.Nil, fmt.Errorf("CHARCONTROLLER.MAKE: height must be numeric")
	}
	x, ok := args[2].ToFloat()
	if !ok {
		return value.Nil, fmt.Errorf("CHARCONTROLLER.MAKE: x must be numeric")
	}
	y, ok := args[3].ToFloat()
	if !ok {
		return value.Nil, fmt.Errorf("CHARCONTROLLER.MAKE: y must be numeric")
	}
	z, ok := args[4].ToFloat()
	if !ok {
		return value.Nil, fmt.Errorf("CHARCONTROLLER.MAKE: z must be numeric")
	}
	ps := mbphysics3d.ActiveJoltPhysics()
	if ps == nil {
		return value.Nil, fmt.Errorf("CHARCONTROLLER.MAKE: PHYSICS3D not started")
	}
	fr, fh := float32(radius), float32(height)
	hh := fh/2 - fr
	if hh < 0.05 {
		hh = 0.05
	}
	capsule := jolt.CreateCapsule(hh, fr)
	settings := jolt.NewCharacterVirtualSettings(capsule)
	settings.MaxSlopeAngle = jolt.DegreesToRadians(45)
	cv := ps.CreateCharacterVirtual(settings, jolt.Vec3{X: float32(x), Y: float32(y), Z: float32(z)})
	h, err := m.h.Alloc(&charObj{cv: cv})
	if err != nil {
		if cv != nil {
			cv.Destroy()
		}
		return value.Nil, err
	}
	trackChar(h)
	return value.FromHandle(h), nil
}

func ccSetPos(m *Module, args []value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, fmt.Errorf("CHARCONTROLLER.SETPOS: heap not bound")
	}
	if len(args) < 4 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("CHARCONTROLLER.SETPOS: need handle, x, y, z")
	}
	co, err := heap.Cast[*charObj](m.h, heap.Handle(args[0].IVal))
	if err != nil || co.cv == nil {
		return value.Nil, fmt.Errorf("CHARCONTROLLER.SETPOS: invalid handle")
	}
	x, ok := args[1].ToFloat()
	if !ok {
		return value.Nil, fmt.Errorf("CHARCONTROLLER.SETPOS: x must be numeric")
	}
	y, ok := args[2].ToFloat()
	if !ok {
		return value.Nil, fmt.Errorf("CHARCONTROLLER.SETPOS: y must be numeric")
	}
	z, ok := args[3].ToFloat()
	if !ok {
		return value.Nil, fmt.Errorf("CHARCONTROLLER.SETPOS: z must be numeric")
	}
	co.cv.SetPosition(jolt.Vec3{X: float32(x), Y: float32(y), Z: float32(z)})
	g := mbphysics3d.GravityVec()
	co.cv.ExtendedUpdate(charDt, g)
	return value.Nil, nil
}

func ccMove(m *Module, args []value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, fmt.Errorf("CHARCONTROLLER.MOVE: heap not bound")
	}
	if len(args) < 4 {
		return value.Nil, fmt.Errorf("CHARCONTROLLER.MOVE: need handle, dx, dy, dz")
	}
	if args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("CHARCONTROLLER.MOVE: first arg must be handle")
	}
	co, err := heap.Cast[*charObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	if co.cv == nil {
		return value.Nil, fmt.Errorf("CHARCONTROLLER.MOVE: invalid handle")
	}
	dx, ok := args[1].ToFloat()
	if !ok {
		return value.Nil, fmt.Errorf("CHARCONTROLLER.MOVE: dx must be numeric")
	}
	dy, ok := args[2].ToFloat()
	if !ok {
		return value.Nil, fmt.Errorf("CHARCONTROLLER.MOVE: dy must be numeric")
	}
	dz, ok := args[3].ToFloat()
	if !ok {
		return value.Nil, fmt.Errorf("CHARCONTROLLER.MOVE: dz must be numeric")
	}
	p := co.cv.GetPosition()
	co.cv.SetPosition(jolt.Vec3{
		X: p.X + float32(dx),
		Y: p.Y + float32(dy),
		Z: p.Z + float32(dz),
	})
	g := mbphysics3d.GravityVec()
	co.cv.ExtendedUpdate(charDt, g)
	return value.Nil, nil
}

func ccAxis(m *Module, args []value.Value, axis int) (value.Value, error) {
	if m.h == nil {
		return value.Nil, fmt.Errorf("CHARCONTROLLER axis: heap not bound")
	}
	if len(args) < 1 || args[0].Kind != value.KindHandle {
		return value.FromFloat(0), nil
	}
	co, err := heap.Cast[*charObj](m.h, heap.Handle(args[0].IVal))
	if err != nil || co.cv == nil {
		return value.FromFloat(0), nil
	}
	p := co.cv.GetPosition()
	switch axis {
	case 0:
		return value.FromFloat(float64(p.X)), nil
	case 1:
		return value.FromFloat(float64(p.Y)), nil
	default:
		return value.FromFloat(float64(p.Z)), nil
	}
}

func ccIsGrounded(m *Module, args []value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, fmt.Errorf("CHARCONTROLLER.ISGROUNDED: heap not bound")
	}
	if len(args) < 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("CHARCONTROLLER.ISGROUNDED: need handle")
	}
	co, err := heap.Cast[*charObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	if co.cv == nil {
		return value.Nil, fmt.Errorf("CHARCONTROLLER.ISGROUNDED: invalid handle")
	}
	if co.cv.IsSupported() {
		return value.FromBool(true), nil
	}
	return value.FromBool(false), nil
}

func ccGetPos(m *Module, args []value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, fmt.Errorf("CHARCONTROLLER.GETPOS: heap not bound")
	}
	if len(args) < 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("CHARCONTROLLER.GETPOS: need handle")
	}
	co, err := heap.Cast[*charObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	if co.cv == nil {
		return value.Nil, fmt.Errorf("CHARCONTROLLER.GETPOS: invalid handle")
	}
	p := co.cv.GetPosition()
	arr, err := heap.NewArray([]int64{3})
	if err != nil {
		return value.Nil, err
	}
	_ = arr.Set([]int64{0}, float64(p.X))
	_ = arr.Set([]int64{1}, float64(p.Y))
	_ = arr.Set([]int64{2}, float64(p.Z))
	ah, err := m.h.Alloc(arr)
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(ah), nil
}

func ccFree(m *Module, args []value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, fmt.Errorf("CHARCONTROLLER.FREE: heap not bound")
	}
	if len(args) < 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("CHARCONTROLLER.FREE: need handle")
	}
	hid := heap.Handle(args[0].IVal)
	untrackChar(hid)
	if err := m.h.Free(hid); err != nil {
		return value.Nil, err
	}
	return value.Nil, nil
}
