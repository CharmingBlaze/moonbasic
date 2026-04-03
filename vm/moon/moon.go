// Package moon implements the MOON bytecode container format (.mbc).
//
// Layout:
//
//	0..3   magic "MOON"
//	4..7   version (big-endian u32, e.g. 0x00020000)
//	8..11  flags (reserved, 0)
//	12..15 entry offset (big-endian u32, byte offset of payload from file start)
//	payload: encoded opcode.Program (see encode.go)
package moon

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"

	"moonbasic/vm/opcode"
)

const (
	magic0 = 'M'
	magic1 = 'O'
	magic2 = 'O'
	magic3 = 'N'
)

// Version is the MOON format version (major<<16 | minor).
const Version uint32 = 0x00020000

// VersionV1 is the legacy MOON 1.x format (rejected at load time).
const VersionV1 uint32 = 0x00010000

var (
	// ErrBadMagic is returned when the file does not start with MOON.
	ErrBadMagic = errors.New("moon: bad magic (not a MOON file)")
	// ErrVersion is returned when the format version is unsupported.
	ErrVersion = errors.New("moon: unsupported format version")
)

// HeaderSize is the fixed MOON header length in bytes.
const HeaderSize = 16

// ValidateHeader checks magic and version without decoding the payload.
func ValidateHeader(r io.Reader) (entryOffset uint32, err error) {
	var hdr [HeaderSize]byte
	if _, err := io.ReadFull(r, hdr[:]); err != nil {
		return 0, err
	}
	if hdr[0] != magic0 || hdr[1] != magic1 || hdr[2] != magic2 || hdr[3] != magic3 {
		return 0, ErrBadMagic
	}
	ver := binary.BigEndian.Uint32(hdr[4:8])
	if ver == VersionV1 {
		return 0, fmt.Errorf("%w: MOON 1.x bytecode is unsupported; recompile from .mb with this moonbasic", ErrVersion)
	}
	if ver != Version {
		return 0, fmt.Errorf("%w: got 0x%08x want 0x%08x", ErrVersion, ver, Version)
	}
	// hdr[8:12] flags — reserved
	entryOffset = binary.BigEndian.Uint32(hdr[12:16])
	if entryOffset < HeaderSize {
		return 0, fmt.Errorf("moon: invalid entry offset %d", entryOffset)
	}
	return entryOffset, nil
}

// Encode serializes prog to MOON bytecode.
func Encode(prog *opcode.Program) ([]byte, error) {
	if prog == nil || prog.Main == nil {
		return nil, fmt.Errorf("moon: nil program")
	}
	var payload bytes.Buffer
	if err := writeProgram(&payload, prog); err != nil {
		return nil, err
	}
	pl := payload.Bytes()
	out := make([]byte, HeaderSize+len(pl))
	out[0], out[1], out[2], out[3] = magic0, magic1, magic2, magic3
	binary.BigEndian.PutUint32(out[4:8], Version)
	binary.BigEndian.PutUint32(out[8:12], 0) // flags
	binary.BigEndian.PutUint32(out[12:16], HeaderSize)
	copy(out[HeaderSize:], pl)
	return out, nil
}

// Decode loads a program from MOON bytecode.
func Decode(data []byte) (*opcode.Program, error) {
	entry, err := ValidateHeader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	if int(entry) > len(data) {
		return nil, fmt.Errorf("moon: entry offset past end of file")
	}
	return readProgram(bytes.NewReader(data[entry:]))
}
