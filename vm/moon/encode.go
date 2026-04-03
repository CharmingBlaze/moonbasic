package moon

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"

	"moonbasic/vm/opcode"
)

func writeProgram(w io.Writer, prog *opcode.Program) error {
	if err := writeStringTable(w, prog.StringTable); err != nil {
		return err
	}
	if err := writeChunk(w, prog.Main); err != nil {
		return err
	}
	names := sortedFuncNames(prog.Functions)
	if err := binaryWriteU32(w, uint32(len(names))); err != nil {
		return err
	}
	for _, name := range names {
		ch := prog.Functions[name]
		if err := writeString(w, name); err != nil {
			return err
		}
		if err := writeChunk(w, ch); err != nil {
			return err
		}
	}
	nTypes := sortedTypeNames(prog.Types)
	if err := binaryWriteU32(w, uint32(len(nTypes))); err != nil {
		return err
	}
	for _, tn := range nTypes {
		td := prog.Types[tn]
		if err := writeString(w, td.Name); err != nil {
			return err
		}
		if err := binaryWriteU32(w, uint32(len(td.Fields))); err != nil {
			return err
		}
		for _, f := range td.Fields {
			if err := writeString(w, f); err != nil {
				return err
			}
		}
	}
	return nil
}

func writeStringTable(w io.Writer, tab []string) error {
	if err := binaryWriteU32(w, uint32(len(tab))); err != nil {
		return err
	}
	for _, s := range tab {
		if err := writeString(w, s); err != nil {
			return err
		}
	}
	return nil
}

func writeChunk(w io.Writer, c *opcode.Chunk) error {
	if c == nil {
		return fmt.Errorf("moon: nil chunk")
	}
	if err := writeString(w, c.Name); err != nil {
		return err
	}
	n := len(c.Instructions)
	if err := binaryWriteU32(w, uint32(n)); err != nil {
		return err
	}
	for i := 0; i < n; i++ {
		in := c.Instructions[i]
		if err := binaryWriteU8(w, byte(in.Op)); err != nil {
			return err
		}
		if err := binaryWriteU8(w, in.Flags); err != nil {
			return err
		}
		var pad [2]byte
		if _, err := w.Write(pad[:]); err != nil {
			return err
		}
		if err := binaryWriteI32(w, in.Operand); err != nil {
			return err
		}
		var line int32
		if i < len(c.SourceLines) {
			line = c.SourceLines[i]
		}
		if err := binaryWriteI32(w, line); err != nil {
			return err
		}
	}
	if err := writeInt64Slice(w, c.IntConsts); err != nil {
		return err
	}
	if err := writeFloat64Slice(w, c.FloatConsts); err != nil {
		return err
	}
	if err := binaryWriteU32(w, uint32(len(c.Names))); err != nil {
		return err
	}
	for _, s := range c.Names {
		if err := writeString(w, s); err != nil {
			return err
		}
	}
	return nil
}

func writeString(w io.Writer, s string) error {
	b := []byte(s)
	if len(b) > 1<<30 {
		return fmt.Errorf("moon: string too long")
	}
	if err := binaryWriteU32(w, uint32(len(b))); err != nil {
		return err
	}
	_, err := w.Write(b)
	return err
}

func writeInt64Slice(w io.Writer, xs []int64) error {
	if err := binaryWriteU32(w, uint32(len(xs))); err != nil {
		return err
	}
	for _, x := range xs {
		var buf [8]byte
		binary.BigEndian.PutUint64(buf[:], uint64(x))
		if _, err := w.Write(buf[:]); err != nil {
			return err
		}
	}
	return nil
}

func writeFloat64Slice(w io.Writer, xs []float64) error {
	if err := binaryWriteU32(w, uint32(len(xs))); err != nil {
		return err
	}
	for _, x := range xs {
		var buf [8]byte
		binary.BigEndian.PutUint64(buf[:], math.Float64bits(x))
		if _, err := w.Write(buf[:]); err != nil {
			return err
		}
	}
	return nil
}
