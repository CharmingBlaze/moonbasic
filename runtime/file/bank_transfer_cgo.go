//go:build cgo || (windows && !cgo)

package mbfile

// ReadBank / WriteBank bridge heap file handles to raw MemObj banks (see mbmem).

import (
	"fmt"
	"io"

	"moonbasic/runtime"
	mbmem "moonbasic/runtime/mbmem"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

func (m *Module) registerBankTransfer(r runtime.Registrar) {
	r.Register("ReadBank", "file", m.readBank)
	r.Register("WriteBank", "file", m.writeBank)
}

// ReadBank(bank, file, offset, count) reads count bytes from the file (sequential read position) into the bank at offset.
func (m *Module) readBank(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if m.h == nil {
		return value.Nil, runtime.Errorf("ReadBank: heap not bound")
	}
	if len(args) != 4 {
		return value.Nil, fmt.Errorf("ReadBank expects (bank, file, offset, count)")
	}
	if args[0].Kind != value.KindHandle || args[1].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("ReadBank: bank and file must be handles")
	}
	mo, err := heap.Cast[*mbmem.MemObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	fo, err := heap.Cast[*fileObj](m.h, heap.Handle(args[1].IVal))
	if err != nil {
		return value.Nil, err
	}
	if fo.rd == nil {
		return value.Nil, fmt.Errorf("ReadBank: file not open for reading")
	}
	off, err := argBankOffset(args[2])
	if err != nil {
		return value.Nil, err
	}
	cnt, err := argBankSize(args[3])
	if err != nil {
		return value.Nil, err
	}
	if cnt == 0 {
		return value.Nil, nil
	}
	buf := mo.Bytes()
	if off+cnt > len(buf) {
		return value.Nil, fmt.Errorf("ReadBank: bank range out of bounds")
	}
	dst := buf[off : off+cnt]
	if _, err := io.ReadFull(fo.rd, dst); err != nil {
		return value.Nil, err
	}
	return value.Nil, nil
}

// WriteBank(bank, file, offset, count) writes count bytes from the bank at offset to the file.
func (m *Module) writeBank(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if m.h == nil {
		return value.Nil, runtime.Errorf("WriteBank: heap not bound")
	}
	if len(args) != 4 {
		return value.Nil, fmt.Errorf("WriteBank expects (bank, file, offset, count)")
	}
	if args[0].Kind != value.KindHandle || args[1].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("WriteBank: bank and file must be handles")
	}
	mo, err := heap.Cast[*mbmem.MemObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	fo, err := heap.Cast[*fileObj](m.h, heap.Handle(args[1].IVal))
	if err != nil {
		return value.Nil, err
	}
	if fo.wr == nil {
		return value.Nil, fmt.Errorf("WriteBank: file not open for writing")
	}
	off, err := argBankOffset(args[2])
	if err != nil {
		return value.Nil, err
	}
	cnt, err := argBankSize(args[3])
	if err != nil {
		return value.Nil, err
	}
	if cnt == 0 {
		return value.Nil, nil
	}
	buf := mo.Bytes()
	if off+cnt > len(buf) {
		return value.Nil, fmt.Errorf("WriteBank: bank range out of bounds")
	}
	src := buf[off : off+cnt]
	if _, err := fo.wr.Write(src); err != nil {
		return value.Nil, err
	}
	return value.Nil, nil
}

func argBankOffset(v value.Value) (int, error) {
	i, ok := v.ToInt()
	if !ok {
		if f, okf := v.ToFloat(); okf {
			i = int64(f)
			ok = true
		}
	}
	if !ok || i < 0 {
		return 0, fmt.Errorf("offset must be non-negative numeric")
	}
	return int(i), nil
}

func argBankSize(v value.Value) (int, error) {
	i, ok := v.ToInt()
	if !ok {
		if f, okf := v.ToFloat(); okf {
			i = int64(f)
			ok = true
		}
	}
	if !ok || i < 0 {
		return 0, fmt.Errorf("count must be non-negative numeric")
	}
	return int(i), nil
}
