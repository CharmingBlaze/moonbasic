// Package intern provides a per-compile string table so lexer/parser share
// canonical identifier spellings (pointer-stable equal strings, fewer allocations).
package intern

import (
	"sync"
)

// Table maps equal strings to one canonical instance.
// It is thread-safe for use in parallel compilation phases if required.
type Table struct {
	mu sync.RWMutex
	m  map[string]string
}

// New returns an empty intern table.
func New() *Table {
	return &Table{m: make(map[string]string)}
}

// Intern returns the canonical copy of s.
func (t *Table) Intern(s string) string {
	if t == nil {
		return s
	}

	t.mu.RLock()
	if c, ok := t.m[s]; ok {
		t.mu.RUnlock()
		return c
	}
	t.mu.RUnlock()

	t.mu.Lock()
	defer t.mu.Unlock()

	// Double check
	if c, ok := t.m[s]; ok {
		return c
	}

	t.m[s] = s
	return s
}

// InternBytes returns the canonical copy of the given byte slice.
// It uses a Go compiler optimization to avoid allocation during lookup: m[string(b)].
func (t *Table) InternBytes(b []byte) string {
	if t == nil {
		return string(b)
	}

	t.mu.RLock()
	if c, ok := t.m[string(b)]; ok {
		t.mu.RUnlock()
		return c
	}
	t.mu.RUnlock()

	s := string(b)
	t.mu.Lock()
	defer t.mu.Unlock()

	if c, ok := t.m[s]; ok {
		return c
	}

	t.m[s] = s
	return s
}

// Size returns the number of unique strings in the table.
func (t *Table) Size() int {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return len(t.m)
}
