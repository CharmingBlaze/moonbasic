package semantic

import (
	"strings"

	"moonbasic/compiler/ast"
	"moonbasic/compiler/builtinmanifest"
)

// inferKind returns a static kind for an expression after folding.
func inferKind(e ast.Expr) builtinmanifest.ArgKind {
	switch n := e.(type) {
	case *ast.IntLitNode:
		return builtinmanifest.Int
	case *ast.FloatLitNode:
		return builtinmanifest.Float
	case *ast.StringLitNode:
		return builtinmanifest.String
	case *ast.BoolLitNode:
		return builtinmanifest.Bool
	case *ast.NullLitNode:
		return builtinmanifest.Any
	case *ast.IdentNode:
		return kindFromIdentName(n.Name)
	case *ast.NewNode:
		return builtinmanifest.Handle
	case *ast.GroupedExpr:
		return inferKind(n.Inner)
	default:
		return builtinmanifest.Any
	}
}

func kindFromIdentName(name string) builtinmanifest.ArgKind {
	if name == "" {
		return builtinmanifest.Any
	}
	switch name[len(name)-1] {
	case '#':
		return builtinmanifest.Float
	case '$':
		return builtinmanifest.String
	case '?':
		return builtinmanifest.Bool
	default:
		return builtinmanifest.Any
	}
}

func compatible(want, got builtinmanifest.ArgKind) bool {
	if want == builtinmanifest.Any || got == builtinmanifest.Any {
		return true
	}
	if want == got {
		return true
	}
	// Numeric coercion
	if isNumeric(want) && isNumeric(got) {
		return true
	}
	// Handles are integer IDs at runtime; unknown locals are Any — already handled
	if want == builtinmanifest.Handle && (got == builtinmanifest.Int || got == builtinmanifest.Handle) {
		return true
	}
	return false
}

func isNumeric(k builtinmanifest.ArgKind) bool {
	return k == builtinmanifest.Int || k == builtinmanifest.Float
}

func kindName(k builtinmanifest.ArgKind) string {
	switch k {
	case builtinmanifest.Int:
		return "INT"
	case builtinmanifest.Float:
		return "FLOAT"
	case builtinmanifest.String:
		return "STRING"
	case builtinmanifest.Bool:
		return "BOOL"
	case builtinmanifest.Handle:
		return "HANDLE"
	default:
		return "ANY"
	}
}

func formatGotKind(e ast.Expr) string {
	return strings.ToLower(kindName(inferKind(e)))
}
