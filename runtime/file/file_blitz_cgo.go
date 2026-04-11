//go:build cgo || (windows && !cgo)

package mbfile

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"

	"moonbasic/runtime"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

func (m *Module) registerFileBlitz(r runtime.Registrar) {
	r.Register("WriteFile", "file", m.blitzOpenWrite)
	r.Register("ReadFile", "file", m.blitzOpenRead)
	r.Register("CloseFile", "file", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		return m.Run(rt, "FILE.CLOSE", args...)
	})
	r.Register("WriteLine", "file", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		return m.Run(rt, "FILE.WRITELN", args...)
	})
	r.Register("ReadLine", "file", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		return m.Run(rt, "FILE.READLINE", args...)
	})
	r.Register("WriteInt", "file", m.fileWriteIntLE)
	r.Register("ReadInt", "file", m.fileReadIntLE)
	r.Register("WriteFloat", "file", m.fileWriteFloatLE)
	r.Register("ReadFloat", "file", m.fileReadFloatLE)
}

func (m *Module) blitzOpenWrite(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Kind != value.KindString {
		return value.Nil, runtime.Errorf("WriteFile expects (path)")
	}
	return m.Run(rt, "FILE.OPENWRITE", args[0])
}

func (m *Module) blitzOpenRead(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Kind != value.KindString {
		return value.Nil, runtime.Errorf("ReadFile expects (path)")
	}
	return m.Run(rt, "FILE.OPENREAD", args[0])
}

func (m *Module) fileWriteIntLE(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	o, err := m.fileExpectWriter(args)
	if err != nil {
		return value.Nil, err
	}
	v, err := rt.ArgInt(args, 1)
	if err != nil {
		return value.Nil, err
	}
	var buf [4]byte
	binary.LittleEndian.PutUint32(buf[:], uint32(int32(v)))
	if _, err := o.wr.Write(buf[:]); err != nil {
		return value.Nil, err
	}
	return value.Nil, nil
}

func (m *Module) fileReadIntLE(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	o, err := m.fileExpectReader(args)
	if err != nil {
		return value.Nil, err
	}
	var buf [4]byte
	if _, err := io.ReadFull(o.rd, buf[:]); err != nil {
		return value.Nil, err
	}
	u := binary.LittleEndian.Uint32(buf[:])
	return value.FromInt(int64(int32(u))), nil
}

func (m *Module) fileWriteFloatLE(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	o, err := m.fileExpectWriter(args)
	if err != nil {
		return value.Nil, err
	}
	x, err := rt.ArgFloat(args, 1)
	if err != nil {
		return value.Nil, err
	}
	var buf [4]byte
	binary.LittleEndian.PutUint32(buf[:], math.Float32bits(float32(x)))
	if _, err := o.wr.Write(buf[:]); err != nil {
		return value.Nil, err
	}
	return value.Nil, nil
}

func (m *Module) fileReadFloatLE(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	o, err := m.fileExpectReader(args)
	if err != nil {
		return value.Nil, err
	}
	var buf [4]byte
	if _, err := io.ReadFull(o.rd, buf[:]); err != nil {
		return value.Nil, err
	}
	bits := binary.LittleEndian.Uint32(buf[:])
	return value.FromFloat(float64(math.Float32frombits(bits))), nil
}

func (m *Module) fileExpectWriter(args []value.Value) (*fileObj, error) {
	if len(args) != 2 || args[0].Kind != value.KindHandle {
		return nil, fmt.Errorf("expects (fileHandle, value)")
	}
	o, err := heap.Cast[*fileObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return nil, err
	}
	if o.wr == nil {
		return nil, fmt.Errorf("file not open for writing")
	}
	return o, nil
}

func (m *Module) fileExpectReader(args []value.Value) (*fileObj, error) {
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return nil, fmt.Errorf("expects (fileHandle)")
	}
	o, err := heap.Cast[*fileObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return nil, err
	}
	if o.rd == nil {
		return nil, fmt.Errorf("file not open for reading")
	}
	return o, nil
}
