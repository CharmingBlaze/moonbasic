//go:build cgo || (windows && !cgo)

package mbcollision

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/runtime"
	"moonbasic/runtime/mbmatrix"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

func (m *Module) registerBSphereBuiltins(reg runtime.Registrar) {
	reg.Register("BSPHERE.CREATE", "collision", runtime.AdaptLegacy(m.bsphereMake))
	reg.Register("BSPHERE.MAKE", "collision", runtime.AdaptLegacy(m.bsphereMake))
	reg.Register("BSPHERE.CHECK", "collision", runtime.AdaptLegacy(m.bsphereCheck))
	reg.Register("BSPHERE.CHECKBOX", "collision", runtime.AdaptLegacy(m.bsphereCheckBox))
	reg.Register("BSPHERE.FREE", "collision", runtime.AdaptLegacy(m.bsphereFree))
	reg.Register("BSPHERE.SETPOS", "collision", runtime.AdaptLegacy(m.bsphereSetPos))
	reg.Register("BSPHERE.SETPOSITION", "collision", runtime.AdaptLegacy(m.bsphereSetPos))
	reg.Register("BSPHERE.SETRADIUS", "collision", runtime.AdaptLegacy(m.bsphereSetRadius))
	reg.Register("BSPHERE.GETPOS", "collision", runtime.AdaptLegacy(m.bsphereGetPos))
	reg.Register("BSPHERE.GETRADIUS", "collision", runtime.AdaptLegacy(m.bsphereGetRadius))
}

func (m *Module) bsphereMake(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 4 {
		return value.Nil, fmt.Errorf("BSPHERE.MAKE expects 4 arguments (cx, cy, cz, r)")
	}
	cx, ok1 := argF(args[0])
	cy, ok2 := argF(args[1])
	cz, ok3 := argF(args[2])
	rad, ok4 := argF(args[3])
	if !ok1 || !ok2 || !ok3 || !ok4 {
		return value.Nil, fmt.Errorf("BSPHERE.MAKE: arguments must be numeric")
	}
	o := &bsphereObj{center: rl.Vector3{X: cx, Y: cy, Z: cz}, radius: rad}
	id, err := m.h.Alloc(o)
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}

func (m *Module) bsphereCheck(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 2 || args[0].Kind != value.KindHandle || args[1].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("BSPHERE.CHECK expects two bounding sphere handles")
	}
	a, err := heap.Cast[*bsphereObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, fmt.Errorf("BSPHERE.CHECK: %w", err)
	}
	b, err := heap.Cast[*bsphereObj](m.h, heap.Handle(args[1].IVal))
	if err != nil {
		return value.Nil, fmt.Errorf("BSPHERE.CHECK: %w", err)
	}
	hit := rl.CheckCollisionSpheres(a.center, a.radius, b.center, b.radius)
	return value.FromBool(hit), nil
}

func (m *Module) bsphereCheckBox(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 2 || args[0].Kind != value.KindHandle || args[1].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("BSPHERE.CHECKBOX expects (sphere, box)")
	}
	s, err := heap.Cast[*bsphereObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, fmt.Errorf("BSPHERE.CHECKBOX: %w", err)
	}
	bo, err := heap.Cast[*bboxObj](m.h, heap.Handle(args[1].IVal))
	if err != nil {
		return value.Nil, fmt.Errorf("BSPHERE.CHECKBOX: %w", err)
	}
	hit := rl.CheckCollisionBoxSphere(bo.box, s.center, s.radius)
	return value.FromBool(hit), nil
}

func (m *Module) bsphereFree(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("BSPHERE.FREE expects bounding sphere handle")
	}
	if err := m.h.Free(heap.Handle(args[0].IVal)); err != nil {
		return value.Nil, err
	}
	return value.Nil, nil
}

func (m *Module) bsphereSetPos(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 4 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("BSPHERE.SETPOS expects (handle, x, y, z)")
	}
	o, err := heap.Cast[*bsphereObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, fmt.Errorf("BSPHERE.SETPOS: %w", err)
	}
	x, ok1 := argF(args[1])
	y, ok2 := argF(args[2])
	z, ok3 := argF(args[3])
	if !ok1 || !ok2 || !ok3 {
		return value.Nil, fmt.Errorf("BSPHERE.SETPOS: components must be numeric")
	}
	o.center = rl.Vector3{X: x, Y: y, Z: z}
	return args[0], nil
}

func (m *Module) bsphereSetRadius(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 2 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("BSPHERE.SETRADIUS expects (handle, radius)")
	}
	o, err := heap.Cast[*bsphereObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, fmt.Errorf("BSPHERE.SETRADIUS: %w", err)
	}
	r, ok := argF(args[1])
	if !ok {
		return value.Nil, fmt.Errorf("BSPHERE.SETRADIUS: radius must be numeric")
	}
	o.radius = r
	return args[0], nil
}

func (m *Module) bsphereGetPos(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("BSPHERE.GETPOS expects handle")
	}
	o, err := heap.Cast[*bsphereObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, fmt.Errorf("BSPHERE.GETPOS: %w", err)
	}
	return mbmatrix.AllocVec3Value(m.h, o.center.X, o.center.Y, o.center.Z)
}

func (m *Module) bsphereGetRadius(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("BSPHERE.GETRADIUS expects handle")
	}
	o, err := heap.Cast[*bsphereObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, fmt.Errorf("BSPHERE.GETRADIUS: %w", err)
	}
	return value.FromFloat(float64(o.radius)), nil
}
