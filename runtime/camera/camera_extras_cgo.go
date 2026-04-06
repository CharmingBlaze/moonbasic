//go:build cgo || (windows && !cgo)

package mbcamera

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/runtime"
	"moonbasic/runtime/mbmatrix"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

func (m *Module) registerCameraExtras(reg runtime.Registrar) {
	reg.Register("CAMERA.GETPOS", "camera", runtime.AdaptLegacy(m.camGetPos))
	reg.Register("CAMERA.GETTARGET", "camera", runtime.AdaptLegacy(m.camGetTarget))
	reg.Register("CAMERA.SETUP", "camera", runtime.AdaptLegacy(m.camSetUp))
	reg.Register("CAMERA.FREE", "camera", runtime.AdaptLegacy(m.camFree))
}

func (m *Module) camGetPos(args []value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("CAMERA.GETPOS: heap not bound")
	}
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("CAMERA.GETPOS expects camera handle")
	}
	h, ok := argHandle(args[0])
	if !ok {
		return value.Nil, fmt.Errorf("CAMERA.GETPOS: invalid handle")
	}
	o, err := heap.Cast[*camObj](m.h, h)
	if err != nil {
		return value.Nil, err
	}
	p := o.cam.Position
	return mbmatrix.AllocVec3Value(m.h, p.X, p.Y, p.Z)
}

func (m *Module) camGetTarget(args []value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("CAMERA.GETTARGET: heap not bound")
	}
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("CAMERA.GETTARGET expects camera handle")
	}
	h, ok := argHandle(args[0])
	if !ok {
		return value.Nil, fmt.Errorf("CAMERA.GETTARGET: invalid handle")
	}
	o, err := heap.Cast[*camObj](m.h, h)
	if err != nil {
		return value.Nil, err
	}
	t := o.cam.Target
	return mbmatrix.AllocVec3Value(m.h, t.X, t.Y, t.Z)
}

func (m *Module) camSetUp(args []value.Value) (value.Value, error) {
	if len(args) != 4 {
		return value.Nil, fmt.Errorf("CAMERA.SETUP expects (handle, ux, uy, uz)")
	}
	h, ok := argHandle(args[0])
	if !ok {
		return value.Nil, fmt.Errorf("CAMERA.SETUP: invalid handle")
	}
	o, err := heap.Cast[*camObj](m.h, h)
	if err != nil {
		return value.Nil, err
	}
	ux, ok1 := argF(args[1])
	uy, ok2 := argF(args[2])
	uz, ok3 := argF(args[3])
	if !ok1 || !ok2 || !ok3 {
		return value.Nil, fmt.Errorf("CAMERA.SETUP: up vector must be numeric")
	}
	o.cam.Up = rl.Vector3{X: ux, Y: uy, Z: uz}
	return value.Nil, nil
}

func (m *Module) camFree(args []value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("CAMERA.FREE: heap not bound")
	}
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("CAMERA.FREE expects camera handle")
	}
	if err := m.h.Free(heap.Handle(args[0].IVal)); err != nil {
		return value.Nil, err
	}
	return value.Nil, nil
}
