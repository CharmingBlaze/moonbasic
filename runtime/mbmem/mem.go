package mbmem

import (
	"encoding/binary"
	"fmt"
	"math"

	"moonbasic/runtime"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

const maxMemBlock = 256 << 20 // 256 MiB cap per block

type memObj struct {
	b []byte
}

func (o *memObj) TypeName() string { return "Mem" }

func (o *memObj) TypeTag() uint16 { return heap.TagMem }

func (o *memObj) Free() { o.b = nil }

func (m *Module) requireHeap() error {
	if m.h == nil {
		return runtime.Errorf("MEM.* builtins: heap not bound")
	}
	return nil
}

func (m *Module) getMem(args []value.Value, ix int, op string) (*memObj, error) {
	if err := m.requireHeap(); err != nil {
		return nil, err
	}
	if ix >= len(args) || args[ix].Kind != value.KindHandle {
		return nil, fmt.Errorf("%s: argument %d must be mem handle", op, ix+1)
	}
	return heap.Cast[*memObj](m.h, heap.Handle(args[ix].IVal))
}

func argInt64(v value.Value) (int64, bool) {
	if i, ok := v.ToInt(); ok {
		return i, true
	}
	if f, ok := v.ToFloat(); ok {
		return int64(f), true
	}
	return 0, false
}

func argSize(v value.Value, op string) (int, error) {
	n, ok := argInt64(v)
	if !ok {
		return 0, fmt.Errorf("%s: size must be numeric", op)
	}
	if n < 0 {
		return 0, fmt.Errorf("%s: size must be non-negative", op)
	}
	if n > maxMemBlock {
		return 0, fmt.Errorf("%s: size exceeds limit (%d bytes)", op, maxMemBlock)
	}
	return int(n), nil
}

func argOffset(v value.Value, op string) (int, error) {
	n, ok := argInt64(v)
	if !ok {
		return 0, fmt.Errorf("%s: offset must be numeric", op)
	}
	if n < 0 {
		return 0, fmt.Errorf("%s: offset must be non-negative", op)
	}
	if n > maxMemBlock {
		return 0, fmt.Errorf("%s: offset out of range", op)
	}
	return int(n), nil
}

// Register implements runtime.Module.
func (m *Module) Register(r runtime.Registrar) {
	r.Register("MEM.MAKE", "mem", runtime.AdaptLegacy(m.memMake))
	r.Register("MEM.FREE", "mem", runtime.AdaptLegacy(m.memFree))
	r.Register("MEM.SIZE", "mem", runtime.AdaptLegacy(m.memSize))
	r.Register("MEM.CLEAR", "mem", runtime.AdaptLegacy(m.memClear))
	r.Register("MEM.COPY", "mem", runtime.AdaptLegacy(m.memCopy))
	r.Register("MEM.GETBYTE", "mem", runtime.AdaptLegacy(m.memGetByte))
	r.Register("MEM.GETWORD", "mem", runtime.AdaptLegacy(m.memGetWord))
	r.Register("MEM.GETDWORD", "mem", runtime.AdaptLegacy(m.memGetDword))
	r.Register("MEM.GETFLOAT", "mem", runtime.AdaptLegacy(m.memGetFloat))
	r.Register("MEM.GETSTRING", "mem", m.memGetString)
	r.Register("MEM.SETBYTE", "mem", runtime.AdaptLegacy(m.memSetByte))
	r.Register("MEM.SETWORD", "mem", runtime.AdaptLegacy(m.memSetWord))
	r.Register("MEM.SETDWORD", "mem", runtime.AdaptLegacy(m.memSetDword))
	r.Register("MEM.SETFLOAT", "mem", runtime.AdaptLegacy(m.memSetFloat))
	r.Register("MEM.SETSTRING", "mem", m.memSetString)
}

// Shutdown implements runtime.Module.
func (m *Module) Shutdown() {}

func (m *Module) memMake(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("MEM.MAKE expects 1 argument (size)")
	}
	n, err := argSize(args[0], "MEM.MAKE")
	if err != nil {
		return value.Nil, err
	}
	buf := make([]byte, n)
	id, err := m.h.Alloc(&memObj{b: buf})
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}

func (m *Module) memFree(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("MEM.FREE expects mem handle")
	}
	if err := m.h.Free(heap.Handle(args[0].IVal)); err != nil {
		return value.Nil, err
	}
	return value.Nil, nil
}

func (m *Module) memSize(args []value.Value) (value.Value, error) {
	o, err := m.getMem(args, 0, "MEM.SIZE")
	if err != nil {
		return value.Nil, err
	}
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("MEM.SIZE expects mem handle")
	}
	return value.FromInt(int64(len(o.b))), nil
}

func (m *Module) memClear(args []value.Value) (value.Value, error) {
	o, err := m.getMem(args, 0, "MEM.CLEAR")
	if err != nil {
		return value.Nil, err
	}
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("MEM.CLEAR expects mem handle")
	}
	clear(o.b)
	return value.Nil, nil
}

func (m *Module) memCopy(args []value.Value) (value.Value, error) {
	if len(args) != 5 {
		return value.Nil, fmt.Errorf("MEM.COPY expects (src, dst, srcOff, dstOff, size)")
	}
	src, err := m.getMem(args, 0, "MEM.COPY")
	if err != nil {
		return value.Nil, err
	}
	dst, err := m.getMem(args, 1, "MEM.COPY")
	if err != nil {
		return value.Nil, err
	}
	srcOff, err := argOffset(args[2], "MEM.COPY")
	if err != nil {
		return value.Nil, err
	}
	dstOff, err := argOffset(args[3], "MEM.COPY")
	if err != nil {
		return value.Nil, err
	}
	n, err := argSize(args[4], "MEM.COPY")
	if err != nil {
		return value.Nil, err
	}
	if n == 0 {
		return value.Nil, nil
	}
	if srcOff+n > len(src.b) || dstOff+n > len(dst.b) {
		return value.Nil, fmt.Errorf("MEM.COPY: range out of bounds")
	}
	copy(dst.b[dstOff:dstOff+n], src.b[srcOff:srcOff+n])
	return value.Nil, nil
}

func (m *Module) memGetByte(args []value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("MEM.GETBYTE expects (mem, offset)")
	}
	o, err := m.getMem(args, 0, "MEM.GETBYTE")
	if err != nil {
		return value.Nil, err
	}
	off, err := argOffset(args[1], "MEM.GETBYTE")
	if err != nil {
		return value.Nil, err
	}
	if off >= len(o.b) {
		return value.Nil, fmt.Errorf("MEM.GETBYTE: offset out of bounds")
	}
	return value.FromInt(int64(o.b[off])), nil
}

func (m *Module) memGetWord(args []value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("MEM.GETWORD expects (mem, offset)")
	}
	o, err := m.getMem(args, 0, "MEM.GETWORD")
	if err != nil {
		return value.Nil, err
	}
	off, err := argOffset(args[1], "MEM.GETWORD")
	if err != nil {
		return value.Nil, err
	}
	if off+2 > len(o.b) {
		return value.Nil, fmt.Errorf("MEM.GETWORD: range out of bounds")
	}
	v := binary.LittleEndian.Uint16(o.b[off : off+2])
	return value.FromInt(int64(v)), nil
}

func (m *Module) memGetDword(args []value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("MEM.GETDWORD expects (mem, offset)")
	}
	o, err := m.getMem(args, 0, "MEM.GETDWORD")
	if err != nil {
		return value.Nil, err
	}
	off, err := argOffset(args[1], "MEM.GETDWORD")
	if err != nil {
		return value.Nil, err
	}
	if off+4 > len(o.b) {
		return value.Nil, fmt.Errorf("MEM.GETDWORD: range out of bounds")
	}
	v := binary.LittleEndian.Uint32(o.b[off : off+4])
	return value.FromInt(int64(v)), nil
}

func (m *Module) memGetFloat(args []value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("MEM.GETFLOAT expects (mem, offset)")
	}
	o, err := m.getMem(args, 0, "MEM.GETFLOAT")
	if err != nil {
		return value.Nil, err
	}
	off, err := argOffset(args[1], "MEM.GETFLOAT")
	if err != nil {
		return value.Nil, err
	}
	if off+4 > len(o.b) {
		return value.Nil, fmt.Errorf("MEM.GETFLOAT: range out of bounds")
	}
	u := binary.LittleEndian.Uint32(o.b[off : off+4])
	f := math.Float32frombits(u)
	return value.FromFloat(float64(f)), nil
}

func (m *Module) memGetString(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("MEM.GETSTRING expects (mem, offset)")
	}
	o, err := m.getMem(args, 0, "MEM.GETSTRING")
	if err != nil {
		return value.Nil, err
	}
	off, err := argOffset(args[1], "MEM.GETSTRING")
	if err != nil {
		return value.Nil, err
	}
	if off > len(o.b) {
		return value.Nil, fmt.Errorf("MEM.GETSTRING: offset out of bounds")
	}
	end := off
	for end < len(o.b) && o.b[end] != 0 {
		end++
	}
	s := string(o.b[off:end])
	return rt.RetString(s), nil
}

func (m *Module) memSetByte(args []value.Value) (value.Value, error) {
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("MEM.SETBYTE expects (mem, offset, value)")
	}
	o, err := m.getMem(args, 0, "MEM.SETBYTE")
	if err != nil {
		return value.Nil, err
	}
	off, err := argOffset(args[1], "MEM.SETBYTE")
	if err != nil {
		return value.Nil, err
	}
	v, ok := argInt64(args[2])
	if !ok {
		return value.Nil, fmt.Errorf("MEM.SETBYTE: value must be numeric")
	}
	if v < 0 || v > 255 {
		return value.Nil, fmt.Errorf("MEM.SETBYTE: value must be 0..255")
	}
	if off >= len(o.b) {
		return value.Nil, fmt.Errorf("MEM.SETBYTE: offset out of bounds")
	}
	o.b[off] = byte(v)
	return value.Nil, nil
}

func (m *Module) memSetWord(args []value.Value) (value.Value, error) {
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("MEM.SETWORD expects (mem, offset, value)")
	}
	o, err := m.getMem(args, 0, "MEM.SETWORD")
	if err != nil {
		return value.Nil, err
	}
	off, err := argOffset(args[1], "MEM.SETWORD")
	if err != nil {
		return value.Nil, err
	}
	v, ok := argInt64(args[2])
	if !ok {
		return value.Nil, fmt.Errorf("MEM.SETWORD: value must be numeric")
	}
	if v < 0 || v > 0xffff {
		return value.Nil, fmt.Errorf("MEM.SETWORD: value must be 0..65535")
	}
	if off+2 > len(o.b) {
		return value.Nil, fmt.Errorf("MEM.SETWORD: range out of bounds")
	}
	binary.LittleEndian.PutUint16(o.b[off:off+2], uint16(v))
	return value.Nil, nil
}

func (m *Module) memSetDword(args []value.Value) (value.Value, error) {
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("MEM.SETDWORD expects (mem, offset, value)")
	}
	o, err := m.getMem(args, 0, "MEM.SETDWORD")
	if err != nil {
		return value.Nil, err
	}
	off, err := argOffset(args[1], "MEM.SETDWORD")
	if err != nil {
		return value.Nil, err
	}
	v, ok := argInt64(args[2])
	if !ok {
		return value.Nil, fmt.Errorf("MEM.SETDWORD: value must be numeric")
	}
	if off+4 > len(o.b) {
		return value.Nil, fmt.Errorf("MEM.SETDWORD: range out of bounds")
	}
	// Low 32 bits (two's complement), matching typical C uint32_t / int32_t layouts.
	binary.LittleEndian.PutUint32(o.b[off:off+4], uint32(v))
	return value.Nil, nil
}

func (m *Module) memSetFloat(args []value.Value) (value.Value, error) {
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("MEM.SETFLOAT expects (mem, offset, value)")
	}
	o, err := m.getMem(args, 0, "MEM.SETFLOAT")
	if err != nil {
		return value.Nil, err
	}
	off, err := argOffset(args[1], "MEM.SETFLOAT")
	if err != nil {
		return value.Nil, err
	}
	f, ok := args[2].ToFloat()
	if !ok {
		if i, ok2 := args[2].ToInt(); ok2 {
			f = float64(i)
			ok = true
		}
	}
	if !ok {
		return value.Nil, fmt.Errorf("MEM.SETFLOAT: value must be numeric")
	}
	if off+4 > len(o.b) {
		return value.Nil, fmt.Errorf("MEM.SETFLOAT: range out of bounds")
	}
	u := math.Float32bits(float32(f))
	binary.LittleEndian.PutUint32(o.b[off:off+4], u)
	return value.Nil, nil
}

func (m *Module) memSetString(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("MEM.SETSTRING expects (mem, offset, string)")
	}
	o, err := m.getMem(args, 0, "MEM.SETSTRING")
	if err != nil {
		return value.Nil, err
	}
	off, err := argOffset(args[1], "MEM.SETSTRING")
	if err != nil {
		return value.Nil, err
	}
	if args[2].Kind != value.KindString {
		return value.Nil, fmt.Errorf("MEM.SETSTRING: third argument must be string")
	}
	s, err := rt.ArgString(args, 2)
	if err != nil {
		return value.Nil, err
	}
	raw := []byte(s)
	need := len(raw) + 1
	if off+need > len(o.b) {
		return value.Nil, fmt.Errorf("MEM.SETSTRING: not enough space (need %d bytes incl. NUL)", need)
	}
	copy(o.b[off:], raw)
	o.b[off+len(raw)] = 0
	return value.Nil, nil
}
