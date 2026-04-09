package mbmem

import (
	"encoding/binary"
	"fmt"
	"math"

	"moonbasic/vm/value"
)

func (m *Module) memResize(args []value.Value) (value.Value, error) {
	o, err := m.getMem(args, 0, "MEM.RESIZE")
	if err != nil {
		return value.Nil, err
	}
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("MEM.RESIZE expects (mem, newSize)")
	}
	n, err := argSize(args[1], "MEM.RESIZE")
	if err != nil {
		return value.Nil, err
	}
	if n == len(o.b) {
		return value.Nil, nil
	}
	nb := make([]byte, n)
	copy(nb, o.b)
	o.b = nb
	return value.Nil, nil
}

func (m *Module) memGetDouble(args []value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("MEM.GETDOUBLE expects (mem, offset)")
	}
	o, err := m.getMem(args, 0, "MEM.GETDOUBLE")
	if err != nil {
		return value.Nil, err
	}
	off, err := argOffset(args[1], "MEM.GETDOUBLE")
	if err != nil {
		return value.Nil, err
	}
	if off%8 != 0 {
		return value.Nil, fmt.Errorf("MEM.GETDOUBLE: offset should be 8-byte aligned (got %d)", off)
	}
	if off+8 > len(o.b) {
		return value.Nil, fmt.Errorf("MEM.GETDOUBLE: range out of bounds")
	}
	u := binary.LittleEndian.Uint64(o.b[off : off+8])
	return value.FromFloat(math.Float64frombits(u)), nil
}

func (m *Module) memSetDouble(args []value.Value) (value.Value, error) {
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("MEM.SETDOUBLE expects (mem, offset, value)")
	}
	o, err := m.getMem(args, 0, "MEM.SETDOUBLE")
	if err != nil {
		return value.Nil, err
	}
	off, err := argOffset(args[1], "MEM.SETDOUBLE")
	if err != nil {
		return value.Nil, err
	}
	if off%8 != 0 {
		return value.Nil, fmt.Errorf("MEM.SETDOUBLE: offset should be 8-byte aligned (got %d)", off)
	}
	x, ok := args[2].ToFloat()
	if !ok {
		if i, ok2 := args[2].ToInt(); ok2 {
			x = float64(i)
			ok = true
		}
	}
	if !ok {
		return value.Nil, fmt.Errorf("MEM.SETDOUBLE: value must be numeric")
	}
	if off+8 > len(o.b) {
		return value.Nil, fmt.Errorf("MEM.SETDOUBLE: range out of bounds")
	}
	binary.LittleEndian.PutUint64(o.b[off:off+8], math.Float64bits(x))
	return value.Nil, nil
}
