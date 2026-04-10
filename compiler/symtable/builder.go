// Package symtable provides the two-pass symbol table builder for implicit variable declaration.
// This is the foundation of the modern syntax - no VAR required.
package symtable

import (
	"moonbasic/compiler/ast"
	"moonbasic/compiler/types"
	"strings"
)

// Builder performs two-pass analysis to collect all variable declarations.
// Pass 1: Collect all assignments as implicit declarations.
// Pass 2: Collect function-local variables and resolve types.
type Builder struct {
	Symbols *Table
	// Track which variables were implicitly declared
	implicitGlobals map[string]bool
	implicitLocals  map[string]map[string]bool // funcName -> varName
}

// NewBuilder creates a symbol table builder.
func NewBuilder() *Builder {
	return &Builder{
		Symbols:         New(),
		implicitGlobals: make(map[string]bool),
		implicitLocals:  make(map[string]map[string]bool),
	}
}

// Build performs two-pass symbol collection from an AST.
// This is the entry point for implicit declaration support.
func (b *Builder) Build(prog *ast.Program) *Table {
	// Pass 1: Pre-declare functions and types (for forward references)
	b.predeclareFunctionsAndTypes(prog)

	// Pass 2: Collect global variable declarations
	b.collectGlobals(prog)

	// Pass 3: Collect function-local variables
	for _, fn := range prog.Functions {
		b.collectFunctionLocals(fn)
	}

	return b.Symbols
}

// predeclareFunctionsAndTypes scans for FUNCTION and TYPE definitions.
func (b *Builder) predeclareFunctionsAndTypes(prog *ast.Program) {
	for _, fn := range prog.Functions {
		b.Symbols.DefineFunction(fn.Name)
	}
	for _, td := range prog.Types {
		b.Symbols.DefineType(td.Name)
	}
}

// collectGlobals finds all variable assignments at program scope.
// In modern syntax, first assignment declares the variable.
func (b *Builder) collectGlobals(prog *ast.Program) {
	for _, stmt := range prog.Stmts {
		b.collectGlobalFromStmt(stmt)
	}
}

// collectGlobalFromStmt extracts global declarations from a statement.
func (b *Builder) collectGlobalFromStmt(stmt ast.Stmt) {
	switch s := stmt.(type) {
	case *ast.AssignNode:
		// First assignment declares the variable
		name := strings.ToUpper(s.Name)
		if !b.Symbols.IsVar(name) {
			// Infer type from suffix or expression
			varType := b.inferType(s.Expr, s.Name)
			sym := b.Symbols.DefineGlobalVar(s.Name)
			sym.Type = varType
			sym.Persistent = true
			b.implicitGlobals[name] = true
		}
	case *ast.MultiAssignNode:
		for _, nm := range s.Names {
			name := strings.ToUpper(nm)
			if !b.Symbols.IsVar(name) {
				sym := b.Symbols.DefineGlobalVar(nm)
				sym.Type = b.inferType(s.Expr, nm)
				sym.Persistent = true
				b.implicitGlobals[name] = true
			}
		}

	case *ast.DimNode:
		// DIM declares an array
		name := strings.ToUpper(s.Name)
		if !b.Symbols.IsVar(name) {
			sym := b.Symbols.DefineGlobalVar(s.Name)
			sym.Type = types.Array
			sym.Persistent = true
			b.implicitGlobals[name] = true
		}

	case *ast.LocalDeclNode:
		// LOCAL at global scope becomes a global (for compatibility)
		name := strings.ToUpper(s.Name)
		if !b.Symbols.IsVar(name) {
			varType := types.FromSuffix(s.Name)
			sym := b.Symbols.DefineGlobalVar(s.Name)
			sym.Type = varType
		}

	case *ast.ConstDeclNode:
		// CONST declares a constant
		sym := b.Symbols.DefineConst(s.Name)
		sym.Type = b.inferType(s.Expr, s.Name)

	case *ast.IfNode:
		// Recursively check IF/THEN/ELSE for assignments
		for _, st := range s.Then {
			b.collectGlobalFromStmt(st)
		}
		for _, st := range s.Else {
			b.collectGlobalFromStmt(st)
		}

	case *ast.WhileNode:
		for _, st := range s.Body {
			b.collectGlobalFromStmt(st)
		}

	case *ast.ForNode:
		for _, st := range s.Body {
			b.collectGlobalFromStmt(st)
		}
		// The loop variable itself is implicitly declared
		loopVar := strings.ToUpper(s.Var)
		if !b.Symbols.IsVar(loopVar) {
			sym := b.Symbols.DefineGlobalVar(s.Var)
			sym.Type = types.Int // FOR loop variables are always int
			sym.Persistent = true
			b.implicitGlobals[loopVar] = true
		}

	case *ast.RepeatNode:
		for _, st := range s.Body {
			b.collectGlobalFromStmt(st)
		}

	case *ast.DoLoopNode:
		for _, st := range s.Body {
			b.collectGlobalFromStmt(st)
		}

	case *ast.SelectNode:
		for _, c := range s.Cases {
			for _, st := range c.Body {
				b.collectGlobalFromStmt(st)
			}
		}
		for _, st := range s.Default {
			b.collectGlobalFromStmt(st)
		}
	}
}

// collectFunctionLocals finds all variable assignments within a function.
func (b *Builder) collectFunctionLocals(fn *ast.FunctionDef) {
	funcName := strings.ToUpper(fn.Name)
	b.implicitLocals[funcName] = make(map[string]bool)

	// Enter function scope
	b.Symbols.PushScope()

	// Define parameters
	for _, p := range fn.Params {
		sym := b.Symbols.DefineParam(p.Name)
		sym.Type = types.FromSuffix(p.Name)
	}

	// Collect locals from function body
	for _, stmt := range fn.Body {
		b.collectLocalFromStmt(stmt, funcName)
	}

	// Pop function scope
	b.Symbols.PopScope()
}

// collectLocalFromStmt extracts local declarations from a statement.
func (b *Builder) collectLocalFromStmt(stmt ast.Stmt, funcName string) {
	switch s := stmt.(type) {
	case *ast.AssignNode:
		name := strings.ToUpper(s.Name)
		// Check if already defined in current scope chain
		if b.Symbols.Resolve(name) == nil {
			// Implicit local declaration
			varType := b.inferType(s.Expr, s.Name)
			sym := b.Symbols.DefineLocal(s.Name)
			sym.Type = varType
			b.implicitLocals[funcName][name] = true
		}
	case *ast.MultiAssignNode:
		for _, nm := range s.Names {
			name := strings.ToUpper(nm)
			if b.Symbols.Resolve(name) == nil {
				sym := b.Symbols.DefineLocal(nm)
				sym.Type = b.inferType(s.Expr, nm)
				b.implicitLocals[funcName][name] = true
			}
		}

	case *ast.DimNode:
		name := strings.ToUpper(s.Name)
		if b.Symbols.Resolve(name) == nil {
			sym := b.Symbols.DefineLocal(s.Name)
			sym.Type = types.Array
			b.implicitLocals[funcName][name] = true
		}

	case *ast.LocalDeclNode:
		name := strings.ToUpper(s.Name)
		if b.Symbols.Resolve(name) == nil {
			varType := types.FromSuffix(s.Name)
			sym := b.Symbols.DefineLocal(s.Name)
			sym.Type = varType
		}

	case *ast.StaticDeclNode:
		// STATIC variables persist across calls
		sym := b.Symbols.DefineStatic(funcName, s.Name)
		if s.Init != nil {
			sym.Type = b.inferType(s.Init, s.Name)
		} else {
			sym.Type = types.FromSuffix(s.Name)
		}

	case *ast.IfNode:
		for _, st := range s.Then {
			b.collectLocalFromStmt(st, funcName)
		}
		for _, st := range s.Else {
			b.collectLocalFromStmt(st, funcName)
		}

	case *ast.WhileNode:
		for _, st := range s.Body {
			b.collectLocalFromStmt(st, funcName)
		}

	case *ast.ForNode:
		for _, st := range s.Body {
			b.collectLocalFromStmt(st, funcName)
		}
		// Loop variable
		loopVar := strings.ToUpper(s.Var)
		if b.Symbols.Resolve(loopVar) == nil {
			sym := b.Symbols.DefineLocal(s.Var)
			sym.Type = types.Int
			b.implicitLocals[funcName][loopVar] = true
		}

	case *ast.RepeatNode:
		for _, st := range s.Body {
			b.collectLocalFromStmt(st, funcName)
		}

	case *ast.DoLoopNode:
		for _, st := range s.Body {
			b.collectLocalFromStmt(st, funcName)
		}

	case *ast.SelectNode:
		for _, c := range s.Cases {
			for _, st := range c.Body {
				b.collectLocalFromStmt(st, funcName)
			}
		}
		for _, st := range s.Default {
			b.collectLocalFromStmt(st, funcName)
		}
	}
}

// inferType deduces the type from an expression or variable suffix.
func (b *Builder) inferType(expr ast.Expr, varName string) types.Tag {
	// First, check for explicit suffix in variable name
	suffixType := types.FromSuffix(varName)
	if suffixType != types.Int { // Has explicit suffix (# $ ?)
		return suffixType
	}

	// Unsuffixed names: float-first inference for implicit declarations (IR v3 / value union).
	return b.inferFromExprImplicit(expr)
}

// inferFromExprImplicit is used for unsuffixed identifiers (Blitz "plain" names): numeric
// literals and unknowns default to float so codegen uses float immediates and the 24-byte
// value union stays consistent with gameplay math.
func (b *Builder) inferFromExprImplicit(expr ast.Expr) types.Tag {
	if expr == nil {
		return types.Float
	}

	switch e := expr.(type) {
	case *ast.IntLitNode:
		return types.Float

	case *ast.FloatLitNode:
		return types.Float

	case *ast.StringLitNode:
		return types.String

	case *ast.BoolLitNode:
		return types.Bool

	case *ast.BinopNode:
		left := b.inferFromExprImplicit(e.Left)
		right := b.inferFromExprImplicit(e.Right)
		return types.Promote(left, right)

	case *ast.IdentNode:
		sym := b.Symbols.Resolve(e.Name)
		if sym != nil && sym.Type != types.Unknown {
			return sym.Type
		}
		if types.FromSuffix(e.Name) != types.Int {
			return types.FromSuffix(e.Name)
		}
		return types.Float

	case *ast.FieldAccessNode:
		return types.Float

	case *ast.IndexExpr:
		return types.Float

	case *ast.CallExprNode:
		return types.FromSuffix(e.Name)

	case *ast.UnaryNode:
		return b.inferFromExprImplicit(e.Expr)

	default:
		return types.Float
	}
}

// IsImplicitGlobal reports whether a variable was implicitly declared at global scope.
func (b *Builder) IsImplicitGlobal(name string) bool {
	return b.implicitGlobals[strings.ToUpper(name)]
}

// IsImplicitLocal reports whether a variable was implicitly declared in a function.
func (b *Builder) IsImplicitLocal(funcName, varName string) bool {
	funcMap, ok := b.implicitLocals[strings.ToUpper(funcName)]
	if !ok {
		return false
	}
	return funcMap[strings.ToUpper(varName)]
}
