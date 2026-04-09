package joltwasm

import (
	"encoding/binary"
	"unsafe"

	"github.com/tetratelabs/wazero/api"
)

// PhysicsStateHeader is the guest-visible prefix of the SoA physics buffer (little-endian).
// The C++/WASM side should place this at [BaseOffset, BaseOffset+physicsStateHeaderSize) and
// pack body transforms as float32 immediately after (layout is engine-specific).
//
// Version is bumped by the guest after a full write so the host can detect torn reads.
type PhysicsStateHeader struct {
	Version   uint32 // monotonically increasing after each complete write
	BodyCount uint32 // number of bodies described in the following float block
	Stride    uint32 // byte stride per body in the SoA block (e.g. 7*4 for pos+quat)
	Reserved  uint32
}

const physicsStateHeaderSize = 16

// StateView aliases guest linear memory for readback without per-frame allocation.
// Use Memory.Read only; do not copy in the hot path unless you must detach from Wasm memory.
type StateView struct {
	Mem        api.Memory
	BaseOffset uint32 // includes header + packed floats
	TotalBytes uint32 // header + bodyCount*stride (or full export size)
}

// Header decodes the header words from linear memory.
func (v StateView) Header() (PhysicsStateHeader, bool) {
	if v.Mem == nil || v.TotalBytes < physicsStateHeaderSize {
		return PhysicsStateHeader{}, false
	}
	buf, ok := v.Mem.Read(v.BaseOffset, physicsStateHeaderSize)
	if !ok || len(buf) < physicsStateHeaderSize {
		return PhysicsStateHeader{}, false
	}
	return PhysicsStateHeader{
		Version:   binary.LittleEndian.Uint32(buf[0:4]),
		BodyCount: binary.LittleEndian.Uint32(buf[4:8]),
		Stride:    binary.LittleEndian.Uint32(buf[8:12]),
		Reserved:  binary.LittleEndian.Uint32(buf[12:16]),
	}, true
}

// RawBytes returns a slice view into Wasm memory covering the whole block [BaseOffset, BaseOffset+TotalBytes).
// This allocates no heap per call beyond the slice header; data aliases guest memory.
func (v StateView) RawBytes() ([]byte, bool) {
	if v.Mem == nil || v.TotalBytes == 0 {
		return nil, false
	}
	return v.Mem.Read(v.BaseOffset, v.TotalBytes)
}

// FloatsView reinterprets the byte view as float32 values without copying.
// len(bytes) must be a multiple of 4.
func FloatsView(bytes []byte) []float32 {
	if len(bytes) < 4 || len(bytes)%4 != 0 {
		return nil
	}
	n := len(bytes) / 4
	return unsafe.Slice((*float32)(unsafe.Pointer(&bytes[0])), n)
}

// FloatsAfterHeader returns float32 view of payload after the fixed header (no alloc).
func (v StateView) FloatsAfterHeader() ([]float32, bool) {
	raw, ok := v.RawBytes()
	if !ok || len(raw) <= physicsStateHeaderSize {
		return nil, false
	}
	return FloatsView(raw[physicsStateHeaderSize:]), true
}
