package heap

import (
	"moonbasic/vm/value"
)

// TypedHolder wraps another HeapObject when you need a stable outer type label.
type TypedHolder struct {
	Inner HeapObject
	Label string
}

// Free forwards to Inner.
func (t *TypedHolder) Free() {
	if t.Inner != nil {
		t.Inner.Free()
	}
}

// TypeName returns Label when set, otherwise Inner.TypeName().
func (t *TypedHolder) TypeName() string {
	if t.Label != "" {
		return t.Label
	}
	if t.Inner != nil {
		return t.Inner.TypeName()
	}
	return "unknown"
}

// TypeTag returns the inner object's type tag.
func (t *TypedHolder) TypeTag() uint16 {
	if t.Inner != nil {
		return t.Inner.TypeTag()
	}
	return TagNone
}

// Instance is a user-defined TYPE instance stored in the heap.
type Instance struct {
	Type   string
	Fields map[string]value.Value
}

func NewInstance(typeName string) *Instance {
	return &Instance{
		Type:   typeName,
		Fields: make(map[string]value.Value),
	}
}

func (i *Instance) Free() {
	// Fields contain value.Value which are plain data/strings. 
	// No native resources to free in a basic moonBASIC instance for now.
}

func (i *Instance) TypeName() string {
	return i.Type
}

func (i *Instance) TypeTag() uint16 {
	return TagInstance
}

func (i *Instance) GetField(name string) value.Value {
	if val, ok := i.Fields[name]; ok {
		return val
	}
	return value.Nil
}

func (i *Instance) SetField(name string, v value.Value) {
	i.Fields[name] = v
}
