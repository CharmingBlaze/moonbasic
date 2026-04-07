// Package heap implements the moonBASIC handle-based resource management system.
// Everything is a handle — models, textures, entities, and bodies are all
// registered here and referenced by an int32 ID.
//
// Handles are opaque int32s: high 16 bits = generation (incremented on Free),
// low 16 bits = slot index. stale handles after Free do not match Get.
package heap

import (
	"fmt"
	"sync"
)

// Handle is an opaque integer index. 0 is always invalid.
type Handle = int32

// Entry is one slot in the handle table (8-byte aligned).
type Entry struct {
	Obj      HeapObject
	TypeTag  uint16 // type safety — catch wrong-type use
	RefCount uint16 // future: ref counting for shared objects
	Gen      uint16 // generation counter — catch use-after-free
	_        uint16 // padding
}

// HeapObject is the interface for any resource stored in the Heap.
//
// Memory safety contract (engine rules):
//   - Free must be idempotent: a second call is a silent no-op. Use [ReleaseOnce] around native
//     Unload/Destroy/Close paths (Raylib, Jolt, Box2D, ENet, OS).
//   - Every externally allocated resource must be registered in the [Store] before its handle
//     is returned to bytecode; [Store.FreeAll] on shutdown must release everything with zero leaks.
//   - TypeTag must be unique per class so [Cast] produces a clear error on wrong-type handles.
//   - Where ownership is layered (e.g. shape before body), document order on the concrete type;
//     parents must not double-free children already freed explicitly.
//   - Raylib CGO is expected on the main OS thread ([runtime.LockOSThread] in main); do not call
//     Free from other goroutines unless the concrete type documents thread safety.
type HeapObject interface {
	Free()            // Unload/Destroy the underlying resource (idempotent)
	TypeName() string // For error messages (e.g., "Texture", "Model")
	TypeTag() uint16  // Unique ID for the type (must match Entry.TypeTag); see heap_tags.go
}

// Store handles the allocation and lookup of resource handles.
type Store struct {
	mu        sync.RWMutex
	entries   []Entry
	free      []uint16
	next      uint16 // next slot index to allocate if free list empty
	strings   []string
	stringMap map[string]int32
}

// New creates a new handle store with a pre-warmed entries slice.
func New() *Store {
	return &Store{
		entries:   make([]Entry, 4096),
		free:      make([]uint16, 0, 1024),
		next:      1, // 0 is reserved / invalid slot
		// String table is seeded from Program.StringTable in VM.Execute so bytecode
		// string indices match GetString; see SeedProgramStrings.
		strings:   nil,
		stringMap: make(map[string]int32),
	}
}

func encodeHandle(slot, gen uint16) Handle {
	return int32((uint32(gen) << 16) | uint32(slot))
}

func decodeHandle(h Handle) (slot, gen uint16) {
	u := uint32(h)
	return uint16(u), uint16(u >> 16)
}

// Alloc registers a new object and returns its unique handle.
func (s *Store) Alloc(obj HeapObject) (Handle, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var slot uint16
	if n := len(s.free); n > 0 {
		slot = s.free[n-1]
		s.free = s.free[:n-1]
	} else {
		slot = s.next
		if int(slot) >= MaxSlots {
			return 0, fmt.Errorf("heap: maximum handle capacity reached (%d)", MaxSlots)
		}
		s.next++

		// Grow entries if needed
		if int(slot) >= len(s.entries) {
			newCap := len(s.entries) * 2
			if newCap > MaxSlots {
				newCap = MaxSlots
			}
			newEntries := make([]Entry, newCap)
			copy(newEntries, s.entries)
			s.entries = newEntries
		}
	}

	e := &s.entries[slot]
	if e.Obj != nil {
		panic("heap: internal error: reused slot not empty")
	}

	e.Obj = obj
	e.TypeTag = obj.TypeTag()
	// Gen is preserved from the last Free() increment
	return encodeHandle(slot, e.Gen), nil
}

// Get retrieves an object by handle. Checks generation to detect use-after-free.
func (s *Store) Get(h Handle) (HeapObject, bool) {
	slot, gen := decodeHandle(h)
	if h == 0 || int(slot) >= len(s.entries) {
		return nil, false
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	e := &s.entries[slot]
	if e.Obj == nil || e.Gen != gen {
		return nil, false
	}
	return e.Obj, true
}

// Free explicitly releases and removes an object from the heap.
func (s *Store) Free(h Handle) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.freeLocked(h)
}

// freeLocked frees a handle; mutex must be held. Recurses into handle arrays.
func (s *Store) freeLocked(h Handle) error {
	slot, gen := decodeHandle(h)
	if h == 0 || int(slot) >= len(s.entries) {
		return fmt.Errorf("heap: invalid handle %d", h)
	}

	e := &s.entries[slot]
	if e.Obj == nil || e.Gen != gen {
		return fmt.Errorf("heap: handle %d is stale or already freed", h)
	}

	if a, ok := e.Obj.(*Array); ok && a.Kind == ArrayKindHandle && len(a.Handles) > 0 {
		kids := append([]int32(nil), a.Handles...)
		for _, hid := range kids {
			if hid != 0 {
				_ = s.freeLocked(Handle(hid))
			}
		}
	}

	e.Obj.Free()
	e.Obj = nil
	e.Gen++ // Invalidate all existing handles to this slot
	s.free = append(s.free, slot)
	return nil
}

// FreeAll releases all objects currently in the heap. Called on VM shutdown.
func (s *Store) FreeAll() {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i := range s.entries {
		e := &s.entries[i]
		if e.Obj != nil {
			e.Obj.Free()
			e.Obj = nil
			e.Gen++
		}
	}
	s.free = s.free[:0]
	s.next = 1
	// Do not reset strings, they are part of the program's static data
}

// Stats returns usage information for debugging.
type Stats struct {
	LiveCount uint32
	FreeSlots uint32
	PeakSlots uint32
}

func (s *Store) Stats() Stats {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var live uint32
	for i := range s.entries {
		if s.entries[i].Obj != nil {
			live++
		}
	}

	return Stats{
		LiveCount: live,
		FreeSlots: uint32(len(s.free)),
		PeakSlots: uint32(s.next - 1),
	}
}

// Count returns the number of active objects in the heap (for debugging).
func (s *Store) Count() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	n := 0
	for i := range s.entries {
		if s.entries[i].Obj != nil {
			n++
		}
	}
	return n
}

// FilterByType returns handles for all active objects of a specific type.
func (s *Store) FilterByType(tag uint16) []Handle {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var res []Handle
	for i := range s.entries {
		e := &s.entries[i]
		if e.Obj != nil && e.TypeTag == tag {
			res = append(res, encodeHandle(uint16(i), e.Gen))
		}
	}
	return res
}

// Cast is a helper to retrieve an object and cast it to a specific type.
func Cast[T HeapObject](s *Store, h Handle) (T, error) {
	var zero T
	obj, ok := s.Get(h)
	if !ok {
		if h == 0 {
			return zero, fmt.Errorf("null handle (0): no object is assigned\n  Hint: Assign a handle from MAKE/LOAD before use; 0 means uninitialized.")
		}
		return zero, fmt.Errorf("invalid or stale handle %d: slot empty or generation mismatch\n  Hint: The object was freed or the handle is outdated; obtain a new handle and avoid using resources after FREE.", h)
	}
	typed, ok := obj.(T)
	if !ok {
		// obj.TypeName() is safe because Get() returned ok (obj != nil).
		return zero, fmt.Errorf("handle %d is %s, but this operation requires a different resource type\n  Hint: Pass the handle returned by the matching MAKE/LOAD for this API (wrong-type handles often come from reusing a variable).", h, obj.TypeName())
	}
	return typed, nil
}

// SeedProgramStrings replaces the heap string table with a copy of the program's compile-time
// pool. Bytecode PUSH_STRING and KindString IVal use these indices; without seeding, the heap's
// default table would shadow index 0 (and break ArgString / builtins for the first literal).
// Intern appends new strings after this slice.
func (s *Store) SeedProgramStrings(tab []string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(tab) == 0 {
		s.strings = nil
		s.stringMap = make(map[string]int32)
		return
	}
	s.strings = make([]string, len(tab))
	copy(s.strings, tab)
	s.stringMap = make(map[string]int32, len(tab)*2)
	for i, str := range s.strings {
		if _, exists := s.stringMap[str]; !exists {
			s.stringMap[str] = int32(i)
		}
	}
}

// Intern adds a string to the heap's string table if it doesn't exist, and returns its index.
func (s *Store) Intern(str string) int32 {
	s.mu.Lock()
	defer s.mu.Unlock()

	if idx, ok := s.stringMap[str]; ok {
		return idx
	}

	idx := int32(len(s.strings))
	s.strings = append(s.strings, str)
	s.stringMap[str] = idx
	return idx
}

// GetString retrieves a string from the heap's string table by its index.
func (s *Store) GetString(idx int32) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if idx < 0 || int(idx) >= len(s.strings) {
		return "", false
	}
	return s.strings[idx], true
}
