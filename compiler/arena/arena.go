// Package arena provides a high-performance, slab-based bump allocator for compiler AST nodes.
// Memory is allocated in large slabs to minimize GC pressure and is freed all at once
// by dropping slab references during Reset.
package arena

import (
	"unsafe"
)

// DefaultSlabSize is the default size for memory slabs (4MiB).
const DefaultSlabSize = 4 * 1024 * 1024

// Arena tracks slab-based allocations for a compilation unit.
type Arena struct {
	slabs    [][]byte
	current  []byte
	offset   int
	slabSize int
}

// NewArena returns an empty arena with the default slab size.
func NewArena() *Arena {
	return &Arena{
		slabSize: DefaultSlabSize,
	}
}

// Reset clears all slabs and resets the arena. Call this after a compile job finishes.
func (a *Arena) Reset() {
	if a == nil {
		return
	}
	a.slabs = nil
	a.current = nil
	a.offset = 0
}

// Alloc returns a slice of n bytes from the arena. It panics if system memory is exhausted.
func (a *Arena) Alloc(n int) []byte {
	if a == nil {
		return make([]byte, n)
	}

	// 8-byte alignment (required for common Go pointer-containing structs)
	padding := (8 - (a.offset % 8)) % 8
	total := n + padding

	if a.current == nil || a.offset+total > len(a.current) {
		a.newSlab(total)
		a.offset = 0
	} else {
		a.offset += padding
	}

	start := a.offset
	a.offset += n
	return a.current[start:a.offset]
}

func (a *Arena) newSlab(n int) {
	size := a.slabSize
	if n > size {
		size = n
	}
	slab := make([]byte, size)
	a.slabs = append(a.slabs, slab)
	a.current = slab
}

// AllocT allocates a zeroed value of type T from the arena and returns its pointer.
func AllocT[T any](a *Arena) *T {
	var zero T
	size := int(unsafe.Sizeof(zero))
	if size == 0 {
		return &zero
	}
	b := a.Alloc(size)
	// We use the first byte's pointer. GC tracks the whole slab while these are live.
	return (*T)(unsafe.Pointer(&b[0]))
}

// Make allocates space for a copy of v in the arena and returns a pointer to it.
func Make[T any](a *Arena, v T) *T {
	p := AllocT[T](a)
	*p = v
	return p
}
