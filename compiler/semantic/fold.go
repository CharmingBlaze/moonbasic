package semantic

import (
	"math"
	"strconv"

	"moonbasic/compiler/ast"
)

// FoldConstants rewrites the program AST by evaluating constant sub-expressions.
func FoldConstants(prog *ast.Program) {
	for _, s := range prog.Stmts {
		foldStmt(s)
	}
	for _, f := range prog.Functions {
		for _, s := range f.Body {
			foldStmt(s)
		}
	}
}

func foldStmt(s ast.Stmt) {
	switch n := s.(type) {
	case *ast.AssignNode:
		n.Expr = foldExpr(n.Expr)
	case *ast.MultiAssignNode:
		n.Expr = foldExpr(n.Expr)
	case *ast.IndexAssignNode:
		for i := range n.Index {
			n.Index[i] = foldExpr(n.Index[i])
		}
		n.Expr = foldExpr(n.Expr)
	case *ast.IndexFieldAssignNode:
		for i := range n.Index {
			n.Index[i] = foldExpr(n.Index[i])
		}
		n.Expr = foldExpr(n.Expr)
	case *ast.FieldAssignNode:
		n.Expr = foldExpr(n.Expr)
	case *ast.CallStmtNode:
		for i := range n.Args {
			n.Args[i] = foldExpr(n.Args[i])
		}
	case *ast.NamespaceCallStmt:
		for i := range n.Args {
			n.Args[i] = foldExpr(n.Args[i])
		}
	case *ast.HandleCallStmt:
		for i := range n.Args {
			n.Args[i] = foldExpr(n.Args[i])
		}
	case *ast.IfNode:
		n.Cond = foldExpr(n.Cond)
		for _, t := range n.Then {
			foldStmt(t)
		}
		for i := range n.ElseIf {
			n.ElseIf[i].Cond = foldExpr(n.ElseIf[i].Cond)
			for _, t := range n.ElseIf[i].Body {
				foldStmt(t)
			}
		}
		for _, t := range n.Else {
			foldStmt(t)
		}
	case *ast.WhileNode:
		n.Cond = foldExpr(n.Cond)
		for _, t := range n.Body {
			foldStmt(t)
		}
	case *ast.ForNode:
		n.From = foldExpr(n.From)
		n.To = foldExpr(n.To)
		if n.Step != nil {
			n.Step = foldExpr(n.Step)
		}
		for _, t := range n.Body {
			foldStmt(t)
		}
	case *ast.RepeatNode:
		for _, t := range n.Body {
			foldStmt(t)
		}
		n.Condition = foldExpr(n.Condition)
	case *ast.DoLoopNode:
		n.Cond = foldExpr(n.Cond)
		for _, t := range n.Body {
			foldStmt(t)
		}
	case *ast.ExitStmt, *ast.ContinueStmt:
		// no expressions
	case *ast.SelectNode:
		n.Expr = foldExpr(n.Expr)
		for i := range n.Cases {
			n.Cases[i].Value = foldExpr(n.Cases[i].Value)
			for _, t := range n.Cases[i].Body {
				foldStmt(t)
			}
		}
		for _, t := range n.Default {
			foldStmt(t)
		}
	case *ast.ReturnNode:
		if n.Expr != nil {
			n.Expr = foldExpr(n.Expr)
		}
	case *ast.DimNode:
		for i := range n.Dims {
			n.Dims[i] = foldExpr(n.Dims[i])
		}
	case *ast.LocalDeclNode:
		if n.Init != nil {
			n.Init = foldExpr(n.Init)
		}
	case *ast.StaticDeclNode:
		if n.Init != nil {
			n.Init = foldExpr(n.Init)
		}
	case *ast.ConstDeclNode:
		n.Expr = foldExpr(n.Expr)
	case *ast.DeleteStmt:
		n.Expr = foldExpr(n.Expr)
	case *ast.EachStmt:
		for _, t := range n.Body {
			foldStmt(t)
		}
	case *ast.ExprStmt:
		n.Expr = foldExpr(n.Expr)
	}
}

func foldExpr(e ast.Expr) ast.Expr {
	switch n := e.(type) {
	case *ast.UnaryNode:
		n.Expr = foldExpr(n.Expr)
		return foldUnary(n)
	case *ast.BinopNode:
		n.Left = foldExpr(n.Left)
		n.Right = foldExpr(n.Right)
		return foldBinop(n)
	case *ast.GroupedExpr:
		n.Inner = foldExpr(n.Inner)
		if lit, ok := n.Inner.(*ast.IntLitNode); ok {
			return lit
		}
		if lit, ok := n.Inner.(*ast.FloatLitNode); ok {
			return lit
		}
		return n
	case *ast.CallExprNode:
		for i := range n.Args {
			n.Args[i] = foldExpr(n.Args[i])
		}
		return n
	case *ast.NamespaceCallExpr:
		for i := range n.Args {
			n.Args[i] = foldExpr(n.Args[i])
		}
		return n
	case *ast.HandleCallExpr:
		for i := range n.Args {
			n.Args[i] = foldExpr(n.Args[i])
		}
		return n
	case *ast.IndexExpr:
		n.Base = foldExpr(n.Base)
		for i := range n.Index {
			n.Index[i] = foldExpr(n.Index[i])
		}
		return n
	case *ast.IndexFieldExpr:
		for i := range n.Index {
			n.Index[i] = foldExpr(n.Index[i])
		}
		return n
	default:
		return e
	}
}

func foldUnary(n *ast.UnaryNode) ast.Expr {
	switch n.Op {
	case "-":
		if iv, ok := n.Expr.(*ast.IntLitNode); ok {
			return &ast.IntLitNode{Value: -iv.Value, Line: n.Line, Col: n.Col}
		}
		if fv, ok := n.Expr.(*ast.FloatLitNode); ok {
			return &ast.FloatLitNode{Value: -fv.Value, Lit: formatFloat(-fv.Value), Line: n.Line, Col: n.Col}
		}
	case "NOT":
		if bv, ok := n.Expr.(*ast.BoolLitNode); ok {
			return &ast.BoolLitNode{Value: !bv.Value, Line: n.Line, Col: n.Col}
		}
	}
	return n
}

func foldBinop(n *ast.BinopNode) ast.Expr {
	// String concatenation
	if n.Op == "+" {
		if ls, ok := n.Left.(*ast.StringLitNode); ok {
			if rs, ok := n.Right.(*ast.StringLitNode); ok {
				return &ast.StringLitNode{Value: ls.Value + rs.Value, Line: n.Line, Col: n.Col}
			}
		}
	}

	li, il := n.Left.(*ast.IntLitNode)
	ri, ir := n.Right.(*ast.IntLitNode)
	lf, fl := n.Left.(*ast.FloatLitNode)
	rf, fr := n.Right.(*ast.FloatLitNode)

	// Pure integer path
	if il && ir {
		if v := foldCompare(n.Op, float64(li.Value), float64(ri.Value)); v != nil {
			return &ast.BoolLitNode{Value: *v, Line: n.Line, Col: n.Col}
		}
		return foldIntBinop(n.Op, li.Value, ri.Value, n.Line, n.Col)
	}
	// At least one float, or int+float
	if fl || fr || (il && fr) || (fl && ir) {
		var a, b float64
		okLeft, okRight := false, false
		if il {
			a, okLeft = float64(li.Value), true
		} else if fl {
			a, okLeft = lf.Value, true
		}
		if ir {
			b, okRight = float64(ri.Value), true
		} else if fr {
			b, okRight = rf.Value, true
		}
		if okLeft && okRight {
			if v := foldCompare(n.Op, a, b); v != nil {
				return &ast.BoolLitNode{Value: *v, Line: n.Line, Col: n.Col}
			}
			return foldFloatBinop(n.Op, a, b, n.Line, n.Col)
		}
	}

	return n
}

func foldIntBinop(op string, a, b int64, line, col int) ast.Expr {
	switch op {
	case "+":
		return &ast.IntLitNode{Value: a + b, Line: line, Col: col}
	case "-":
		return &ast.IntLitNode{Value: a - b, Line: line, Col: col}
	case "*":
		return &ast.IntLitNode{Value: a * b, Line: line, Col: col}
	case "/":
		if b == 0 {
			return &ast.BinopNode{Op: op, Left: &ast.IntLitNode{Value: a, Line: line, Col: col}, Right: &ast.IntLitNode{Value: b, Line: line, Col: col}, Line: line, Col: col}
		}
		return &ast.IntLitNode{Value: a / b, Line: line, Col: col}
	case "MOD":
		if b == 0 {
			return &ast.BinopNode{Op: op, Left: &ast.IntLitNode{Value: a, Line: line, Col: col}, Right: &ast.IntLitNode{Value: b, Line: line, Col: col}, Line: line, Col: col}
		}
		return &ast.IntLitNode{Value: a % b, Line: line, Col: col}
	case "^":
		return &ast.FloatLitNode{Value: math.Pow(float64(a), float64(b)), Lit: formatFloat(math.Pow(float64(a), float64(b))), Line: line, Col: col}
	default:
		return &ast.BinopNode{Op: op, Left: &ast.IntLitNode{Value: a, Line: line, Col: col}, Right: &ast.IntLitNode{Value: b, Line: line, Col: col}, Line: line, Col: col}
	}
}

func foldFloatBinop(op string, a, b float64, line, col int) ast.Expr {
	switch op {
	case "+":
		return &ast.FloatLitNode{Value: a + b, Lit: formatFloat(a + b), Line: line, Col: col}
	case "-":
		return &ast.FloatLitNode{Value: a - b, Lit: formatFloat(a - b), Line: line, Col: col}
	case "*":
		return &ast.FloatLitNode{Value: a * b, Lit: formatFloat(a * b), Line: line, Col: col}
	case "/":
		if b == 0 {
			return &ast.BinopNode{Op: op,
				Left:  &ast.FloatLitNode{Value: a, Lit: formatFloat(a), Line: line, Col: col},
				Right: &ast.FloatLitNode{Value: b, Lit: formatFloat(b), Line: line, Col: col},
				Line:  line, Col: col}
		}
		return &ast.FloatLitNode{Value: a / b, Lit: formatFloat(a / b), Line: line, Col: col}
	case "MOD":
		return &ast.FloatLitNode{Value: math.Mod(a, b), Lit: formatFloat(math.Mod(a, b)), Line: line, Col: col}
	case "^":
		return &ast.FloatLitNode{Value: math.Pow(a, b), Lit: formatFloat(math.Pow(a, b)), Line: line, Col: col}
	default:
		return &ast.BinopNode{Op: op,
			Left:  &ast.FloatLitNode{Value: a, Lit: formatFloat(a), Line: line, Col: col},
			Right: &ast.FloatLitNode{Value: b, Lit: formatFloat(b), Line: line, Col: col},
			Line:  line, Col: col}
	}
}

func foldCompare(op string, a, b float64) *bool {
	var v bool
	switch op {
	case "=":
		v = a == b
	case "<>":
		v = a != b
	case "<":
		v = a < b
	case ">":
		v = a > b
	case "<=":
		v = a <= b
	case ">=":
		v = a >= b
	default:
		return nil
	}
	return &v
}

func formatFloat(f float64) string {
	return strconv.FormatFloat(f, 'g', -1, 64)
}
