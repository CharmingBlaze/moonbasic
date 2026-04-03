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

// MaxSlots is the maximum number of distinct heap slots (0 is invalid).
const MaxSlots = 65535

const (
	TagNone uint16 = iota
	TagInstance
	TagArray
	TagSprite
	TagTexture
	TagFont
	TagCamera
	TagFile
	TagJSON
	TagHost
	TagPeer
	TagEvent
	TagPhysicsBody
	TagPhysicsBuilder
	TagCharController
	TagAutomationList
	TagImage
	TagMesh
	TagMaterial
	TagModel
	TagShader
	TagMatrix
	TagVec2
	TagVec3
	TagRay
	TagBBox
	TagBSphere
	TagAudioStream
	TagWave
	TagSound
	TagColor
	TagMem
	TagRng
	TagStringList
	TagPhysics2D
	TagBody2D
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
// It is responsible for releasing its own native resources (raylib, jolt, etc).
type HeapObject interface {
	Free()            // Unload/Destroy the underlying resource
	TypeName() string // For error messages (e.g., "Texture", "Model")
	TypeTag() uint16  // Unique ID for the type (must match Entry.TypeTag)
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
		next:      1,            // 0 is reserved / invalid slot
		strings:   []string{""}, // index 0 is always empty string
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
	slot, gen := decodeHandle(h)
	if h == 0 || int(slot) >= len(s.entries) {
		return fmt.Errorf("heap: invalid handle %d", h)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	e := &s.entries[slot]
	if e.Obj == nil || e.Gen != gen {
		return fmt.Errorf("heap: handle %d is stale or already freed", h)
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

// Cast is a helper to retrieve an object and cast it to a specific type.
func Cast[T HeapObject](s *Store, h Handle) (T, error) {
	var zero T
	obj, ok := s.Get(h)
	if !ok {
		return zero, fmt.Errorf("heap: invalid or stale handle %d", h)
	}
	typed, ok := obj.(T)
	if !ok {
		// obj.TypeName() is safe because Get() returned ok (obj != nil).
		// We use %T on the zero value to report the expected type without calling methods on it.
		return zero, fmt.Errorf("heap: handle %d is %s, but expected type %T", h, obj.TypeName(), zero)
	}
	return typed, nil
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
