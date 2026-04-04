//go:build cgo

package mbfile

import (
	"fmt"

	"moonbasic/runtime"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

func (m *Module) writeOpenForWrite(rt *runtime.Runtime, args []value.Value) (*fileObj, string, error) {
	if len(args) != 2 || args[0].Kind != value.KindHandle || args[1].Kind != value.KindString {
		return nil, "", runtime.Errorf("expects (handle, text$)")
	}
	o, err := heap.Cast[*fileObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return nil, "", err
	}
	if o.wr == nil {
		return nil, "", fmt.Errorf("file not open for writing")
	}
	text, err := rt.ArgString(args, 1)
	if err != nil {
		return nil, "", err
	}
	return o, text, nil
}

func (m *Module) fileWriteRaw(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	o, text, err := m.writeOpenForWrite(rt, args)
	if err != nil {
		return value.Nil, err
	}
	_, err = o.wr.WriteString(text)
	return value.Nil, err
}

func (m *Module) fileWriteLine(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	o, text, err := m.writeOpenForWrite(rt, args)
	if err != nil {
		return value.Nil, err
	}
	_, err = o.wr.WriteString(text + "\n")
	return value.Nil, err
}
