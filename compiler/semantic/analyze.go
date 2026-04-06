package semantic

import (
	"fmt"

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
	return a.checkProgram(prog)
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

	// 1. Main
	a.currentFunc = "<MAIN>"
	for _, s := range prog.Stmts {
		if err := a.checkStmt(s); err != nil {
			return err
		}
	}

	// Functions
	for _, f := range prog.Functions {
		a.currentFunc = f.Name
		for _, s := range f.Body {
			if err := a.checkStmt(s); err != nil {
				return err
			}
		}
	}

	a.currentFunc = "<MAIN>"
	return nil
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
		return a.checkExprCalls(n.Expr)
	case *ast.IndexAssignNode:
		for _, e := range n.Index {
			if err := a.checkExprCalls(e); err != nil {
				return err
			}
		}
		return a.checkExprCalls(n.Expr)
	case *ast.IndexFieldAssignNode:
		for _, e := range n.Index {
			if err := a.checkExprCalls(e); err != nil {
				return err
			}
		}
		return a.checkExprCalls(n.Expr)
	case *ast.FieldAssignNode:
		return a.checkExprCalls(n.Expr)
	case *ast.CallStmtNode:
		for _, e := range n.Args {
			if err := a.checkExprCalls(e); err != nil {
				return err
			}
		}
	case *ast.HandleCallStmt:
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
	case *ast.ConstDeclNode:
		if a.currentFunc != "<MAIN>" {
			return a.typeError(n.Line, n.Col, "CONST is only allowed at module scope", "Move CONST to the top-level program, outside any FUNCTION.")
		}
		return a.checkExprCalls(n.Expr)
	case *ast.StaticDeclNode:
		if n.Init != nil {
			return a.checkExprCalls(n.Init)
		}
	case *ast.SwapStmt, *ast.EraseStmt:
		return nil
	case *ast.LocalDeclNode:
		if n.Init != nil {
			return a.checkExprCalls(n.Init)
		}
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
