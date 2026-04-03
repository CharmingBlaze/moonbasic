// Package symtable tracks symbols for parsing and code generation.
// Names are stored as lexed (uppercased); suffixes # $ ? are part of the string, so
// X and X# are distinct symbols (Blitz-style).
package symtable

import (
	"strings"
)

// Kind classifies a symbol.
type Kind int

const (
	None Kind = iota
	Var
	Local
	Param
	Func
	TypeSym
	Const
	Static // function-local static; storage key in Symbol.StaticKey
)

// Symbol is one entry in a scope.
type Symbol struct {
	Name      string
	Kind      Kind
	Slot      int    // local/param slot index when applicable
	StaticKey string // globals key for KindStatic (FUNCNAME`VARNAME)
}

// Table is a stack of scopes plus global declarations.
type Table struct {
	globals   map[string]*Symbol
	scopes    []map[string]*Symbol
	funcs     map[string]bool
	types     map[string]bool
	nextLocal int
}

// New builds an empty symbol table.
func New() *Table {
	return &Table{
		globals:   make(map[string]*Symbol),
		scopes:    nil,
		funcs:     make(map[string]bool),
		types:     make(map[string]bool),
		nextLocal: 0,
	}
}

// PredeclareFunction records a function name (forward reference).
func (t *Table) PredeclareFunction(name string) {
	t.funcs[strings.ToUpper(name)] = true
}

// PredeclareType records a user type name.
func (t *Table) PredeclareType(name string) {
	t.types[strings.ToUpper(name)] = true
}

// IsFunction reports whether name is a known user function.
func (t *Table) IsFunction(name string) bool {
	return t.funcs[strings.ToUpper(name)]
}

// IsTypeName reports whether name is a user-defined type.
func (t *Table) IsTypeName(name string) bool {
	return t.types[strings.ToUpper(name)]
}

// IsVar reports whether name refers to a variable in current visibility.
func (t *Table) IsVar(name string) bool {
	name = strings.ToUpper(name)
	for i := len(t.scopes) - 1; i >= 0; i-- {
		if s, ok := t.scopes[i][name]; ok && (s.Kind == Var || s.Kind == Local || s.Kind == Param || s.Kind == Static) {
			return true
		}
	}
	if s, ok := t.globals[name]; ok && (s.Kind == Var || s.Kind == Const) {
		return true
	}
	return false
}

// DefineConst defines a global constant.
func (t *Table) DefineConst(name string) *Symbol {
	name = strings.ToUpper(name)
	s := &Symbol{Name: name, Kind: Const}
	t.globals[name] = s
	return s
}

// DefineGlobalVar defines a global variable (implicit or explicit).
func (t *Table) DefineGlobalVar(name string) *Symbol {
	name = strings.ToUpper(name)
	if s, ok := t.globals[name]; ok {
		return s
	}
	s := &Symbol{Name: name, Kind: Var}
	t.globals[name] = s
	return s
}

// PushScope enters a function or block scope.
func (t *Table) PushScope() {
	t.scopes = append(t.scopes, make(map[string]*Symbol))
	t.nextLocal = 0
}

// PopScope leaves the innermost scope.
func (t *Table) PopScope() {
	if len(t.scopes) == 0 {
		return
	}
	t.scopes = t.scopes[:len(t.scopes)-1]
}

// ScopeDepth is the number of nested local scopes (0 at program top-level).
func (t *Table) ScopeDepth() int {
	return len(t.scopes)
}

// DefineLocal defines a local variable in the current scope and returns its slot.
func (t *Table) DefineLocal(name string) *Symbol {
	name = strings.ToUpper(name)
	scope := t.scopes[len(t.scopes)-1]
	if s, ok := scope[name]; ok {
		return s
	}
	s := &Symbol{Name: name, Kind: Local, Slot: t.nextLocal}
	t.nextLocal++
	scope[name] = s
	return s
}

// DefineStatic declares a STATIC variable inside the current function scope.
func (t *Table) DefineStatic(funcName, varName string) *Symbol {
	varName = strings.ToUpper(varName)
	if len(t.scopes) == 0 {
		return nil
	}
	scope := t.scopes[len(t.scopes)-1]
	if s, ok := scope[varName]; ok {
		return s
	}
	key := strings.ToUpper(funcName + "`" + varName)
	s := &Symbol{Name: varName, Kind: Static, StaticKey: key}
	scope[varName] = s
	return s
}

// DefineParam defines a parameter in the current (function) scope.
func (t *Table) DefineParam(name string) *Symbol {
	name = strings.ToUpper(name)
	scope := t.scopes[len(t.scopes)-1]
	s := &Symbol{Name: name, Kind: Param, Slot: t.nextLocal}
	t.nextLocal++
	scope[name] = s
	return s
}

// DefineFunction registers a function body name in globals scope map as Func.
func (t *Table) DefineFunction(name string) *Symbol {
	name = strings.ToUpper(name)
	s := &Symbol{Name: name, Kind: Func}
	t.globals[name] = s
	t.PredeclareFunction(name)
	return s
}

// DefineType registers a type name.
func (t *Table) DefineType(name string) *Symbol {
	name = strings.ToUpper(name)
	s := &Symbol{Name: name, Kind: TypeSym}
	t.globals[name] = s
	t.PredeclareType(name)
	return s
}

// Resolve finds a symbol by name, searching innermost scope first then globals.
func (t *Table) Resolve(name string) *Symbol {
	name = strings.ToUpper(name)
	for i := len(t.scopes) - 1; i >= 0; i-- {
		if s, ok := t.scopes[i][name]; ok {
			return s
		}
	}
	if s, ok := t.globals[name]; ok {
		return s
	}
	return nil
}

// SlotOf returns the local slot for a resolved local/param, or -1.
func (t *Table) SlotOf(name string) int {
	s := t.Resolve(name)
	if s == nil {
		return -1
	}
	if s.Kind == Local || s.Kind == Param {
		return s.Slot
	}
	return -1
}

// Funcs returns a read-only view of the declared function names.
// Used by the parser to propagate forward declarations into nested scopes.
func (t *Table) Funcs() map[string]bool { return t.funcs }

// Types returns a read-only view of the declared type names.
func (t *Table) Types() map[string]bool { return t.types }
