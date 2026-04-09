// Package builtinmanifest holds compile-time signatures for built-in commands.
// The compiler uses this for semantic analysis without importing runtime.
package builtinmanifest

import (
	"fmt"
	"sort"
	"strings"
)

// ArgKind is a coarse argument type for static checking.
type ArgKind int

const (
	Any ArgKind = iota
	Int
	Float
	String
	Bool
	Handle
)

// Command describes one dotted built-in name and its argument kinds.
// Key is the canonical manifest name (from JSON "key"); optional metadata is for tooling/docs.
type Command struct {
	Key       string
	Args      []ArgKind
	Returns   string `json:"returns,omitempty"`
	Pure      bool   `json:"pure,omitempty"`
	Phase     string `json:"phase,omitempty"`
	Namespace string `json:"namespace,omitempty"`
	Stub      string `json:"stub,omitempty"`
}

// Table maps canonical command keys to one or more overloads (different arities).
type Table struct {
	Commands map[string][]Command
}

// Key builds the lookup key for a namespace and method (already uppercased by lexer/parser).
func Key(ns, method string) string {
	return ns + "." + method
}

// LookupArity returns the overload whose arity matches argc, or false if none.
func (t *Table) LookupArity(ns, method string, argc int) (Command, bool) {
	if t == nil || t.Commands == nil {
		return Command{}, false
	}
	ovs := t.Commands[Key(ns, method)]
	for _, c := range ovs {
		if len(c.Args) == argc {
			return c, true
		}
	}
	return Command{}, false
}

// Has reports whether any overload exists for NS.METHOD.
func (t *Table) Has(ns, method string) bool {
	if t == nil || t.Commands == nil {
		return false
	}
	ovs := t.Commands[Key(ns, method)]
	return len(ovs) > 0
}

// ArityHint lists expected argument counts for NS.METHOD overloads.
func (t *Table) ArityHint(ns, method string) string {
	if t == nil || t.Commands == nil {
		return ""
	}
	ovs := t.Commands[Key(ns, method)]
	if len(ovs) == 0 {
		return ""
	}
	var parts []string
	for _, c := range ovs {
		parts = append(parts, fmt.Sprintf("%d", len(c.Args)))
	}
	return fmt.Sprintf("Overloads expect argument count(s): %s.", strings.Join(parts, ", "))
}

var defaultTable = mustDefaultTable()

// Default returns the embedded JSON manifest (Raylib/Jolt-style engine surface).
func Default() *Table {
	return defaultTable
}

// Keys returns sorted canonical command names (one entry per overload key; duplicates collapsed).
func (t *Table) Keys() []string {
	if t == nil || t.Commands == nil {
		return nil
	}
	keys := make([]string, 0, len(t.Commands))
	for k := range t.Commands {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// HasArityExact reports whether any overload of a global (non-dotted) command name has exactly
// argc parameters. Name is normalized the same way as manifest keys (see NormalizeCommand).
func (t *Table) HasArityExact(globalName string, argc int) bool {
	if t == nil || t.Commands == nil {
		return false
	}
	k := NormalizeCommand(globalName)
	ovs := t.Commands[k]
	for _, c := range ovs {
		if len(c.Args) == argc {
			return true
		}
	}
	return false
}

// FirstOverload returns the first manifest entry for key (for docs/LSP when arity is unknown).
func (t *Table) FirstOverload(key string) (Command, bool) {
	if t == nil || t.Commands == nil {
		return Command{}, false
	}
	ovs := t.Commands[key]
	if len(ovs) == 0 {
		return Command{}, false
	}
	return ovs[0], true
}

// NormalizeCommand applies the same dotted-name rule as runtime (uppercase segments).
func NormalizeCommand(name string) string {
	parts := strings.Split(name, ".")
	for i := range parts {
		parts[i] = strings.ToUpper(strings.TrimSpace(parts[i]))
	}
	return strings.Join(parts, ".")
}
