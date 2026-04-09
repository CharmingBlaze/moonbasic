package semantic

import (
	"fmt"
	"strings"

	"moonbasic/compiler/ast"
	"moonbasic/compiler/builtinmanifest"
	"moonbasic/compiler/errors"
)

// Analyzer performs constant folding and static checks after parsing.
type Analyzer struct {
	File  string
	Lines []string
	Table *builtinmanifest.Table
	Fold  bool

	// Static Analysis: caller -> set of callee names
	CallGraph   map[string]map[string]bool
	currentFunc string

	// Milestone 6: User-defined Types
	Types map[string]*ast.TypeDef

	funcNames map[string]bool // user FUNCTION names (uppercase)

	// Scopes for tracking assigned variables (Implicit Declaration/First-Assignment)
	// scopes[0] is always global.
	scopes []map[string]bool
}

// DefaultAnalyzer uses the built-in command manifest and enables folding.
func DefaultAnalyzer(file string, lines []string) *Analyzer {
	return &Analyzer{
		File:  file,
		Lines: lines,
		Table: builtinmanifest.Default(),
		Fold:  true,
	}
}

// Run folds constants (if enabled) and type-checks namespace built-in calls.
func (a *Analyzer) Run(prog *ast.Program) error {
	a.CallGraph = make(map[string]map[string]bool)
	a.Types = make(map[string]*ast.TypeDef)
	a.currentFunc = "<MAIN>"

	if a.Fold {
		FoldConstants(prog)
	}
	a.scopes = []map[string]bool{make(map[string]bool)} // Global scope
	a.seedBuiltinConstants()
	return a.checkProgram(prog)
}

func (a *Analyzer) seedBuiltinConstants() {
	// Subset of key constants from runtime/keyglobals.go
	keys := []string{
		"KEY_ESCAPE", "KEY_SPACE", "KEY_W", "KEY_A", "KEY_S", "KEY_D",
		"KEY_Q", "KEY_E", "KEY_I", "KEY_K", "KEY_LEFT", "KEY_RIGHT", "KEY_UP", "KEY_DOWN",
		"KEY_1", "KEY_2", "KEY_3", "KEY_4", "KEY_5", "KEY_6",
		"KEY_F1", "KEY_F2", "KEY_F3", "KEY_F4", "KEY_F5", "KEY_F6", "KEY_F7", "KEY_F8", "KEY_F9", "KEY_F10", "KEY_F11", "KEY_F12",
		"GAMEPAD_AXIS_LEFT_X", "GAMEPAD_AXIS_LEFT_Y",
		"GAMEPAD_AXIS_RIGHT_X", "GAMEPAD_AXIS_RIGHT_Y",

		"GAMEPAD_BUTTON_RIGHT_FACE_DOWN", "GAMEPAD_BUTTON_RIGHT_FACE_RIGHT",
		"GAMEPAD_BUTTON_RIGHT_FACE_LEFT", "GAMEPAD_BUTTON_RIGHT_FACE_UP",
		"GAMEPAD_BUTTON_LEFT_FACE_UP", "GAMEPAD_BUTTON_LEFT_FACE_DOWN", "GAMEPAD_BUTTON_LEFT_FACE_LEFT", "GAMEPAD_BUTTON_LEFT_FACE_RIGHT",
	}
	for _, k := range keys {
		a.scopes[0][k] = true
	}
}

func (a *Analyzer) lineText(line int) string {
	if line < 1 || line > len(a.Lines) {
		return ""
	}
	return a.Lines[line-1]
}

func (a *Analyzer) typeError(line, col int, msg, hint string) error {
	return errors.NewTypeError(a.File, line, col, msg, a.lineText(line), hint)
}

func (a *Analyzer) checkProgram(prog *ast.Program) error {
	a.funcNames = make(map[string]bool)
	for _, f := range prog.Functions {
		a.funcNames[f.Name] = true
	}

	// 0. Register Types (Pass 0)
	for _, t := range prog.Types {
		if _, exists := a.Types[t.Name]; exists {
			return a.typeError(t.Line, t.Col, fmt.Sprintf("duplicate type definition %s", t.Name), "Remove or rename the duplicate type.")
		}
		a.Types[t.Name] = t
		// Verify fields are unique within type
		seen := make(map[string]bool)
		for _, f := range t.Fields {
			if seen[f] {
				return a.typeError(t.Line, t.Col, fmt.Sprintf("duplicate field %s in type %s", f, t.Name), "Ensure field names within a TYPE are unique.")
			}
			seen[f] = true
		}
	}

	a.currentFunc = "<MAIN>"
	for _, s := range prog.Stmts {
		if err := a.checkStmt(s); err != nil {
			return err
		}
	}

	// Functions
	for _, f := range prog.Functions {
		a.currentFunc = f.Name
		a.pushScope()
		// Parameters are implicitly assigned
		for _, p := range f.Params {
			a.assign(p.Name)
		}
		for _, s := range f.Body {
			if err := a.checkStmt(s); err != nil {
				return err
			}
		}
		a.popScope()
	}

	a.currentFunc = "<MAIN>"
	return nil
}

func (a *Analyzer) pushScope() {
	a.scopes = append(a.scopes, make(map[string]bool))
}

func (a *Analyzer) popScope() {
	if len(a.scopes) > 1 {
		a.scopes = a.scopes[:len(a.scopes)-1]
	}
}

func (a *Analyzer) assign(name string) {
	name = strings.ToUpper(name)
	a.scopes[len(a.scopes)-1][name] = true
}

func (a *Analyzer) isAssigned(name string) bool {
	name = strings.ToUpper(name)
	for i := len(a.scopes) - 1; i >= 0; i-- {
		if a.scopes[i][name] {
			return true
		}
	}
	return false
}

func (a *Analyzer) checkStmt(s ast.Stmt) error {
	if ns, ok := s.(*ast.NamespaceCallStmt); ok {
		return a.checkNamespaceCall(ns.NS, ns.Method, ns.Args, ns.Line, ns.Col)
	}
	return a.walkStmtExprs(s)
}

func (a *Analyzer) walkStmtExprs(s ast.Stmt) error {
	switch n := s.(type) {
	case *ast.AssignNode:
		if err := a.checkExprCalls(n.Expr); err != nil {
			return err
		}
		a.assign(n.Name)
		return nil
	case *ast.IndexAssignNode:
		if !a.isAssigned(n.Array) {
			return a.typeError(n.Line, n.Col, fmt.Sprintf("use of unassigned variable %s", n.Array), "Assign an array to the variable before subscripting.")
		}
		for _, e := range n.Index {
			if err := a.checkExprCalls(e); err != nil {
				return err
			}
		}
		return a.checkExprCalls(n.Expr)
	case *ast.IndexFieldAssignNode:
		if !a.isAssigned(n.Array) {
			return a.typeError(n.Line, n.Col, fmt.Sprintf("use of unassigned variable %s", n.Array), "Assign an array to the variable before subscripting.")
		}
		for _, e := range n.Index {
			if err := a.checkExprCalls(e); err != nil {
				return err
			}
		}
		return a.checkExprCalls(n.Expr)
	case *ast.FieldAssignNode:
		// Object.Field = expr
		if !a.isAssigned(n.Object) {
			return a.typeError(n.Line, n.Col, fmt.Sprintf("use of unassigned variable %s", n.Object), "Assign a value to the variable before accessing its fields.")
		}
		return a.checkExprCalls(n.Expr)
	case *ast.CallStmtNode:
		for _, e := range n.Args {
			if err := a.checkExprCalls(e); err != nil {
				return err
			}
		}
	case *ast.HandleCallStmt:
		if !a.isAssigned(n.Receiver) {
			return a.typeError(n.Line, n.Col, fmt.Sprintf("use of unassigned variable %s", n.Receiver), "Assign a value to the variable before calling methods on it.")
		}
		for _, e := range n.Args {
			if err := a.checkExprCalls(e); err != nil {
				return err
			}
		}
	case *ast.IfNode:
		if err := a.checkExprCalls(n.Cond); err != nil {
			return err
		}
		for _, t := range n.Then {
			if err := a.checkStmt(t); err != nil {
				return err
			}
		}
		for _, ei := range n.ElseIf {
			if err := a.checkExprCalls(ei.Cond); err != nil {
				return err
			}
			for _, t := range ei.Body {
				if err := a.checkStmt(t); err != nil {
					return err
				}
			}
		}
		for _, t := range n.Else {
			if err := a.checkStmt(t); err != nil {
				return err
			}
		}
	case *ast.WhileNode:
		if err := a.checkExprCalls(n.Cond); err != nil {
			return err
		}
		for _, t := range n.Body {
			if err := a.checkStmt(t); err != nil {
				return err
			}
		}
	case *ast.ForNode:
		a.assign(n.Var) // Iterator is assigned
		for _, e := range []ast.Expr{n.From, n.To} {
			if err := a.checkExprCalls(e); err != nil {
				return err
			}
		}
		if n.Step != nil {
			if err := a.checkExprCalls(n.Step); err != nil {
				return err
			}
		}
		for _, t := range n.Body {
			if err := a.checkStmt(t); err != nil {
				return err
			}
		}
	case *ast.RepeatNode:
		for _, t := range n.Body {
			if err := a.checkStmt(t); err != nil {
				return err
			}
		}
		if err := a.checkExprCalls(n.Condition); err != nil {
			return err
		}
	case *ast.DoLoopNode:
		if err := a.checkExprCalls(n.Cond); err != nil {
			return err
		}
		for _, t := range n.Body {
			if err := a.checkStmt(t); err != nil {
				return err
			}
		}
	case *ast.ExitStmt, *ast.ContinueStmt:
		return nil
	case *ast.SelectNode:
		if err := a.checkExprCalls(n.Expr); err != nil {
			return err
		}
		for _, c := range n.Cases {
			if err := a.checkExprCalls(c.Value); err != nil {
				return err
			}
			for _, t := range c.Body {
				if err := a.checkStmt(t); err != nil {
					return err
				}
			}
		}
		for _, t := range n.Default {
			if err := a.checkStmt(t); err != nil {
				return err
			}
		}
	case *ast.ReturnNode:
		if n.Expr != nil {
			return a.checkExprCalls(n.Expr)
		}
	case *ast.DimNode:
		for _, e := range n.Dims {
			if err := a.checkExprCalls(e); err != nil {
				return err
			}
		}
		a.assign(n.Name)
		return nil
	case *ast.ConstDeclNode:
		if a.currentFunc != "<MAIN>" {
			return a.typeError(n.Line, n.Col, "CONST is only allowed at module scope", "Move CONST to the top-level program, outside any FUNCTION.")
		}
		if err := a.checkExprCalls(n.Expr); err != nil {
			return err
		}
		a.assign(n.Name)
		return nil
	case *ast.StaticDeclNode:
		if n.Init != nil {
			if err := a.checkExprCalls(n.Init); err != nil {
				return err
			}
		}
		a.assign(n.Name)
		return nil
	case *ast.SwapStmt, *ast.EraseStmt:
		return nil
	case *ast.LocalDeclNode:
		if n.Init != nil {
			if err := a.checkExprCalls(n.Init); err != nil {
				return err
			}
		}
		a.assign(n.Name)
		return nil
	case *ast.DeleteStmt:
		return a.checkExprCalls(n.Expr)
	case *ast.EachStmt:
		for _, t := range n.Body {
			if err := a.checkStmt(t); err != nil {
				return err
			}
		}
	case *ast.ExprStmt:
		return a.checkExprCalls(n.Expr)
	}
	return nil
}

func (a *Analyzer) checkExprCalls(e ast.Expr) error {
	switch n := e.(type) {
	case *ast.NamespaceCallExpr:
		if err := a.checkNamespaceCall(n.NS, n.Method, n.Args, n.Line, n.Col); err != nil {
			return err
		}
	case *ast.BinopNode:
		if err := a.checkExprCalls(n.Left); err != nil {
			return err
		}
		return a.checkExprCalls(n.Right)
	case *ast.UnaryNode:
		return a.checkExprCalls(n.Expr)
	case *ast.GroupedExpr:
		return a.checkExprCalls(n.Inner)
	case *ast.CallExprNode:
		if td, ok := a.Types[n.Name]; ok && !a.funcNames[n.Name] {
			if len(n.Args) != len(td.Fields) {
				return a.typeError(n.Line, n.Col,
					fmt.Sprintf("type %s constructor expects %d arguments, got %d", n.Name, len(td.Fields), len(n.Args)),
					"Pass one value per field in declaration order.")
			}
			for _, arg := range n.Args {
				if err := a.checkExprCalls(arg); err != nil {
					return err
				}
			}
			return nil
		}
		for _, arg := range n.Args {
			if err := a.checkExprCalls(arg); err != nil {
				return err
			}
		}
	case *ast.IndexFieldExpr:
		for _, arg := range n.Index {
			if err := a.checkExprCalls(arg); err != nil {
				return err
			}
		}
	case *ast.HandleCallExpr:
		if !a.isAssigned(n.Receiver) {
			return a.typeError(n.Line, n.Col, fmt.Sprintf("use of unassigned variable %s", n.Receiver), "Assign a value to the variable before calling methods on it.")
		}
		for _, arg := range n.Args {
			if err := a.checkExprCalls(arg); err != nil {
				return err
			}
		}
	case *ast.IndexExpr:
		if err := a.checkExprCalls(n.Base); err != nil {
			return err
		}
		for _, x := range n.Index {
			if err := a.checkExprCalls(x); err != nil {
				return err
			}
		}
	case *ast.IdentNode:
		if !a.isAssigned(n.Name) {
			return a.typeError(n.Line, n.Col, fmt.Sprintf("use of unassigned variable %s", n.Name), "Assign a value to the variable before using it in an expression.")
		}
	case *ast.NewNode:
		if _, exists := a.Types[n.TypeName]; !exists {
			return a.typeError(n.Line, n.Col, fmt.Sprintf("unknown type %s", n.TypeName), "Ensure the type is defined with TYPE ... END TYPE before use.")
		}
	}
	return nil
}

func (a *Analyzer) checkNamespaceCall(ns, method string, args []ast.Expr, line, col int) error {
	for _, arg := range args {
		if err := a.checkExprCalls(arg); err != nil {
			return err
		}
	}

	cmd, ok := a.Table.LookupArity(ns, method, len(args))
	if !ok {
		if a.Table.Has(ns, method) {
			hint := a.Table.ArityHint(ns, method)
			return a.typeError(line, col,
				fmt.Sprintf("%s.%s: no overload matches %d argument(s)", ns, method, len(args)),
				hint)
		}
		msg, hint := unknownCommandMessageAndHint(a.Table, ns, method)
		return a.typeError(line, col, msg, hint)
	}
	key := builtinmanifest.Key(ns, method)

	// Record CallGraph edge
	if _, exists := a.CallGraph[a.currentFunc]; !exists {
		a.CallGraph[a.currentFunc] = make(map[string]bool)
	}
	a.CallGraph[a.currentFunc][key] = true

	if len(args) != len(cmd.Args) {
		return a.typeError(line, col,
			fmt.Sprintf("%s.%s expects %d argument(s), got %d", ns, method, len(cmd.Args), len(args)),
			fmt.Sprintf("Provide %d argument(s) matching the built-in signature.", len(cmd.Args)))
	}
	for i, want := range cmd.Args {
		got := inferKind(args[i])
		if !compatible(want, got) {
			return a.typeError(line, col,
				fmt.Sprintf("%s.%s argument %d: expected %s, got %s", ns, method, i+1, kindName(want), formatGotKind(args[i])),
				"Fix the argument type to match the built-in signature.")
		}
	}
	return nil
}
