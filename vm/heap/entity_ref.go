package heap

// EntityRef wraps a Blitz-style entity id (integer) so scripts can use handle dot-syntax
// (cube.Pos …) while the runtime still stores entities in the entity module.
type EntityRef struct {
	ID int64
}

// EntityFreeHook is set by mbentity to remove the native entity when the heap frees
// an EntityRef (e.g. ERASE ALL) or when the object is explicitly freed.
var EntityFreeHook func(id int64)

func (e *EntityRef) TypeName() string { return "ENTITYREF" }
func (e *EntityRef) GetID() int64    { return e.ID }

func (e *EntityRef) TypeTag() uint16 { return TagEntityRef }

func (e *EntityRef) Free() {
	if EntityFreeHook != nil && e.ID > 0 {
		EntityFreeHook(e.ID)
	}
	e.ID = 0
}
