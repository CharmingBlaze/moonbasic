package pipeline

import (
	"fmt"

	"moonbasic/compiler/arena"
	"moonbasic/compiler/include"
	"moonbasic/compiler/parser"
	"moonbasic/compiler/symtable"
)

// BuildSymbolTable runs the implicit-declaration symbol builder on parsed source (after INCLUDE expansion).
// It does not run semantic analysis or codegen. Use for LSP outline and tooling when a fast symbol list is needed.
func BuildSymbolTable(name, src string) (*symtable.Table, error) {
	SyncPackageIncludeRoots()
	ar := arena.NewArena()
	defer ar.Reset()
	prog, err := parser.ParseSourceWithArena(name, src, ar)
	if err != nil {
		return nil, err
	}
	prog, err = include.ExpandWithArena(name, prog, ar)
	if err != nil {
		return nil, err
	}
	builder := symtable.NewBuilder()
	return builder.Build(prog), nil
}

// ExportSymbolTableJSON returns JSON from [symtable.Table.ExportJSONWithPath] for the given source
// (includes path, globals with Persistent, funcs, types).
func ExportSymbolTableJSON(name, src string) ([]byte, error) {
	t, err := BuildSymbolTable(name, src)
	if err != nil {
		return nil, err
	}
	return t.ExportJSONWithPath(name)
}

// DocumentSymbols returns a minimal LSP-friendly symbol list (name + kind) from the symbol table.
func DocumentSymbols(name, src string) ([]map[string]any, error) {
	t, err := BuildSymbolTable(name, src)
	if err != nil {
		return nil, err
	}
	out := make([]map[string]any, 0, 32)
	t.ForEachGlobal(func(n string, sym *symtable.Symbol) {
		if sym == nil {
			return
		}
		m := map[string]any{
			"name":   n,
			"kind":   lspSymbolKind(sym.Kind),
			"detail": fmt.Sprintf("%s %s", sym.Kind.String(), sym.Type.String()),
		}
		if sym.Kind == symtable.Var || sym.Kind == symtable.Const {
			m["persistent"] = sym.Persistent
		}
		out = append(out, m)
	})
	return out, nil
}

// lspSymbolKind maps symtable.Kind to LSP SymbolKind (numeric).
func lspSymbolKind(k symtable.Kind) int {
	switch k {
	case symtable.Const:
		return 14 // Constant
	case symtable.Var:
		return 13 // Variable
	case symtable.Local, symtable.Param:
		return 13
	case symtable.Func:
		return 12
	case symtable.TypeSym:
		return 5
	case symtable.Static:
		return 13
	default:
		return 13
	}
}
