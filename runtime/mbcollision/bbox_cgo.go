//go:build cgo

package mbcollision

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/runtime"
	"moonbasic/runtime/mbmodel3d"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

func (m *Module) registerBBoxBuiltins(reg runtime.Registrar) {
	reg.Register("BBOX.MAKE", "collision", runtime.AdaptLegacy(m.bboxMake))
	reg.Register("BBOX.FROMMODEL", "collision", runtime.AdaptLegacy(m.bboxFromModel))
	reg.Register("BBOX.CHECK", "collision", runtime.AdaptLegacy(m.bboxCheck))
	reg.Register("BBOX.CHECKSPHERE", "collision", runtime.AdaptLegacy(m.bboxCheckSphere))
	reg.Register("BBOX.FREE", "collision", runtime.AdaptLegacy(m.bboxFree))
}

func (m *Module) bboxMake(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 6 {
		return value.Nil, fmt.Errorf("BBOX.MAKE expects 6 arguments (minX, minY, minZ, maxX, maxY, maxZ)")
	}
	var min, max rl.Vector3
	for i, tgt := range []*rl.Vector3{&min, &max} {
		base := i * 3
		x, okx := argF(args[base])
		y, oky := argF(args[base+1])
		z, okz := argF(args[base+2])
		if !okx || !oky || !okz {
			return value.Nil, fmt.Errorf("BBOX.MAKE: arguments must be numeric")
		}
		tgt.X, tgt.Y, tgt.Z = x, y, z
	}
	o := &bboxObj{box: rl.BoundingBox{Min: min, Max: max}}
	id, err := m.h.Alloc(o)
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}

func (m *Module) bboxFromModel(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("BBOX.FROMMODEL expects model handle")
	}
	model, err := mbmodel3d.ModelRaylib(m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, fmt.Errorf("BBOX.FROMMODEL: %w", err)
	}
	box := rl.GetModelBoundingBox(model)
	o := &bboxObj{box: box}
	id, err := m.h.Alloc(o)
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}

func (m *Module) bboxCheck(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 2 || args[0].Kind != value.KindHandle || args[1].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("BBOX.CHECK expects two bounding box handles")
	}
	a, err := heap.Cast[*bboxObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, fmt.Errorf("BBOX.CHECK: %w", err)
	}
	b, err := heap.Cast[*bboxObj](m.h, heap.Handle(args[1].IVal))
	if err != nil {
		return value.Nil, fmt.Errorf("BBOX.CHECK: %w", err)
	}
	return value.FromBool(rl.CheckCollisionBoxes(a.box, b.box)), nil
}

func (m *Module) bboxCheckSphere(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 5 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("BBOX.CHECKSPHERE expects (box, cx, cy, cz, r)")
	}
	bo, err := heap.Cast[*bboxObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, fmt.Errorf("BBOX.CHECKSPHERE: %w", err)
	}
	cx, ok1 := argF(args[1])
	cy, ok2 := argF(args[2])
	cz, ok3 := argF(args[3])
	rad, ok4 := argF(args[4])
	if !ok1 || !ok2 || !ok3 || !ok4 {
		return value.Nil, fmt.Errorf("BBOX.CHECKSPHERE: center and radius must be numeric")
	}
	center := rl.Vector3{X: cx, Y: cy, Z: cz}
	return value.FromBool(rl.CheckCollisionBoxSphere(bo.box, center, rad)), nil
}

func (m *Module) bboxFree(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("BBOX.FREE expects bounding box handle")
	}
	m.h.Free(heap.Handle(args[0].IVal))
	return value.Nil, nil
}
