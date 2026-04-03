package moon

import (
	"encoding/binary"
	"io"
	"sort"

	"moonbasic/vm/opcode"
)

func binaryWriteU8(w io.Writer, v byte) error {
	_, err := w.Write([]byte{v})
	return err
}

func binaryWriteU32(w io.Writer, v uint32) error {
	var buf [4]byte
	binary.BigEndian.PutUint32(buf[:], v)
	_, err := w.Write(buf[:])
	return err
}

func binaryWriteI32(w io.Writer, v int32) error {
	return binaryWriteU32(w, uint32(v))
}

func sortedFuncNames(m map[string]*opcode.Chunk) []string {
	names := make([]string, 0, len(m))
	for k := range m {
		names = append(names, k)
	}
	sort.Strings(names)
	return names
}

func sortedTypeNames(m map[string]*opcode.TypeDef) []string {
	names := make([]string, 0, len(m))
	for k := range m {
		names = append(names, k)
	}
	sort.Strings(names)
	return names
}
