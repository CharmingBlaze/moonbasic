package moon

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"

	"moonbasic/vm/opcode"
)

func readProgram(r io.Reader) (*opcode.Program, error) {
	tab, err := readStringTable(r)
	if err != nil {
		return nil, err
	}
	main, err := readChunk(r)
	if err != nil {
		return nil, err
	}
	nFn, err := binaryReadU32(r)
	if err != nil {
		return nil, err
	}
	funcs := make(map[string]*opcode.Chunk)
	for i := uint32(0); i < nFn; i++ {
		name, err := readString(r)
		if err != nil {
			return nil, err
		}
		ch, err := readChunk(r)
		if err != nil {
			return nil, err
		}
		funcs[name] = ch
	}
	nTy, err := binaryReadU32(r)
	if err != nil {
		return nil, err
	}
	types := make(map[string]*opcode.TypeDef)
	for i := uint32(0); i < nTy; i++ {
		tname, err := readString(r)
		if err != nil {
			return nil, err
		}
		nf, err := binaryReadU32(r)
		if err != nil {
			return nil, err
		}
		fields := make([]string, 0, nf)
		for j := uint32(0); j < nf; j++ {
			f, err := readString(r)
			if err != nil {
				return nil, err
			}
			fields = append(fields, f)
		}
		types[tname] = &opcode.TypeDef{Name: tname, Fields: fields}
	}
	return &opcode.Program{StringTable: tab, Main: main, Functions: funcs, Types: types}, nil
}

func readStringTable(r io.Reader) ([]string, error) {
	n, err := binaryReadU32(r)
	if err != nil {
		return nil, err
	}
	out := make([]string, 0, n)
	for i := uint32(0); i < n; i++ {
		s, err := readString(r)
		if err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	return out, nil
}

func readChunk(r io.Reader) (*opcode.Chunk, error) {
	name, err := readString(r)
	if err != nil {
		return nil, err
	}
	n, err := binaryReadU32(r)
	if err != nil {
		return nil, err
	}
	c := opcode.NewChunk(name)
	for i := uint32(0); i < n; i++ {
		opB, err := binaryReadU8(r)
		if err != nil {
			return nil, err
		}
		dst, err := binaryReadU8(r)
		if err != nil {
			return nil, err
		}
		srcA, err := binaryReadU8(r)
		if err != nil {
			return nil, err
		}
		srcB, err := binaryReadU8(r)
		if err != nil {
			return nil, err
		}
		op1, err := binaryReadI32(r)
		if err != nil {
			return nil, err
		}
		line, err := binaryReadI32(r)
		if err != nil {
			return nil, err
		}
		c.Instructions = append(c.Instructions, opcode.Instruction{
			Op:      opcode.OpCode(opB),
			Dst:     dst,
			SrcA:    srcA,
			SrcB:    srcB,
			Operand: op1,
		})
		c.SourceLines = append(c.SourceLines, line)
	}
	c.IntConsts, err = readInt64Slice(r)
	if err != nil {
		return nil, err
	}
	c.FloatConsts, err = readFloat64Slice(r)
	if err != nil {
		return nil, err
	}
	nn, err := binaryReadU32(r)
	if err != nil {
		return nil, err
	}
	for i := uint32(0); i < nn; i++ {
		s, err := readString(r)
		if err != nil {
			return nil, err
		}
		c.Names = append(c.Names, s)
	}
	return c, nil
}

func readString(r io.Reader) (string, error) {
	n, err := binaryReadU32(r)
	if err != nil {
		return "", err
	}
	if n > 1<<30 {
		return "", fmt.Errorf("moon: string length too large")
	}
	buf := make([]byte, n)
	if _, err := io.ReadFull(r, buf); err != nil {
		return "", err
	}
	return string(buf), nil
}

func readInt64Slice(r io.Reader) ([]int64, error) {
	n, err := binaryReadU32(r)
	if err != nil {
		return nil, err
	}
	out := make([]int64, n)
	var buf [8]byte
	for i := uint32(0); i < n; i++ {
		if _, err := io.ReadFull(r, buf[:]); err != nil {
			return nil, err
		}
		out[i] = int64(binary.BigEndian.Uint64(buf[:]))
	}
	return out, nil
}

func readFloat64Slice(r io.Reader) ([]float64, error) {
	n, err := binaryReadU32(r)
	if err != nil {
		return nil, err
	}
	out := make([]float64, n)
	var buf [8]byte
	for i := uint32(0); i < n; i++ {
		if _, err := io.ReadFull(r, buf[:]); err != nil {
			return nil, err
		}
		out[i] = math.Float64frombits(binary.BigEndian.Uint64(buf[:]))
	}
	return out, nil
}

func binaryReadU8(r io.Reader) (byte, error) {
	var b [1]byte
	if _, err := io.ReadFull(r, b[:]); err != nil {
		return 0, err
	}
	return b[0], nil
}

func binaryReadU32(r io.Reader) (uint32, error) {
	var buf [4]byte
	if _, err := io.ReadFull(r, buf[:]); err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint32(buf[:]), nil
}

func binaryReadI32(r io.Reader) (int32, error) {
	u, err := binaryReadU32(r)
	if err != nil {
		return 0, err
	}
	return int32(u), nil
}
