//go:build cgo || (windows && !cgo)

package mbentity

import (
	"fmt"

	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

// entFreeEntities frees every entity id stored in a 1D (or flat) numeric array (e.g. DIM entity handles).
func (m *Module) entFreeEntities(args []value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, fmt.Errorf("ENTITY.FREEENTITIES: heap not bound")
	}
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("ENTITY.FREEENTITIES expects array handle")
	}
	arr, err := heap.Cast[*heap.Array](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, fmt.Errorf("ENTITY.FREEENTITIES: %w", err)
	}
	if arr.Kind != heap.ArrayKindFloat && arr.Kind != heap.ArrayKindBool {
		return value.Nil, fmt.Errorf("ENTITY.FREEENTITIES: array must be numeric (entity# per cell)")
	}
	for _, v := range arr.Floats {
		id := int64(v)
		if id >= 1 {
			m.purgeEntityByID(id)
		}
	}
	return value.Nil, nil
}
