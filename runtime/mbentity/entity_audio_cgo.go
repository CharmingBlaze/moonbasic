//go:build cgo || (windows && !cgo)

package mbentity

import (
	"fmt"

	mbaudio "moonbasic/runtime/audio"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

func (m *Module) entEmitSound(args []value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("EmitSound expects (soundHandle, entity#)")
	}
	if args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("EmitSound: sound must be handle")
	}
	id, ok := m.entID(args[1])
	if !ok || id < 1 {
		return value.Nil, fmt.Errorf("EmitSound: invalid entity")
	}
	e := m.store().ents[id]
	if e == nil {
		return value.Nil, fmt.Errorf("EmitSound: unknown entity")
	}
	if m.h == nil {
		return value.Nil, fmt.Errorf("EmitSound: heap not bound")
	}
	wp := m.worldPos(e)
	sh := heap.Handle(args[0].IVal)
	if err := mbaudio.PlaySpatial(m.h, sh, wp.X, wp.Y, wp.Z); err != nil {
		return value.Nil, err
	}
	return value.Nil, nil
}
