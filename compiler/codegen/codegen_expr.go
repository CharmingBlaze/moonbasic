package codegen

import (
	"fmt"

	"moonbasic/compiler/ast"
	"moonbasic/compiler/symtable"
	"moonbasic/vm/opcode"
)

// emitExpr translates an expression AST node into bytecode.
func (g *CodeGen) emitExpr(ch *opcode.Chunk, e ast.Expr) {
	if g.err != nil {
		return
	}
	switch n := e.(type) {
	case *ast.IntLitNode:
		idx := ch.AddInt(n.Value)
		ch.Emit(opcode.OpPushInt, idx, 0, n.Line)

	case *ast.FloatLitNode:
		idx := ch.AddFloat(n.Value)
		ch.Emit(opcode.OpPushFloat, idx, 0, n.Line)

	case *ast.StringLitNode:
		idx := g.Prog.InternString(n.Value)
		ch.Emit(opcode.OpPushString, idx, 0, n.Line)

	case *ast.BoolLitNode:
		v := int32(0)
		if n.Value {
			v = 1
		}
		ch.Emit(opcode.OpPushBool, v, 0, n.Line)

	case *ast.NullLitNode:
		ch.Emit(opcode.OpPushNull, 0, 0, n.Line)

	case *ast.BinopNode:
		g.emitExpr(ch, n.Left)
		g.emitExpr(ch, n.Right)
		var op opcode.OpCode
		switch n.Op {
		case "+":
			op = opcode.OpAdd
		case "-":
			op = opcode.OpSub
		case "*":
			op = opcode.OpMul
		case "/":
			op = opcode.OpDiv
		case ">":
			op = opcode.OpGt
		case "<":
			op = opcode.OpLt
		case ">=":
			op = opcode.OpGte
		case "<=":
			op = opcode.OpLte
		case "=":
			op = opcode.OpEq
		case "<>":
			op = opcode.OpNeq
		case "AND":
			op = opcode.OpAnd
		case "OR":
			op = opcode.OpOr
		case "XOR":
			op = opcode.OpXor
		default:
			g.codegenError(n.Line, n.Col, fmt.Sprintf("unsupported operator %q", n.Op), "")
			return
		}
		ch.Emit(op, 0, 0, n.Line)

	case *ast.IdentNode:
		sym := g.Symbols.Resolve(n.Name)
		if sym != nil && (sym.Kind == symtable.Local || sym.Kind == symtable.Param) {
			ch.Emit(opcode.OpLoadLocal, int32(sym.Slot), 0, n.Line)
		} else if sym != nil && sym.Kind == symtable.Static {
			k := ch.AddName(sym.StaticKey)
			ch.Emit(opcode.OpLoadGlobal, k, 0, n.Line)
		} else {
			idx := ch.AddName(n.Name)
			ch.Emit(opcode.OpLoadGlobal, idx, 0, n.Line)
		}

	case *ast.CallExprNode:
		if td, ok := g.Prog.Types[n.Name]; ok && g.Prog.Functions[n.Name] == nil {
			if len(n.Args) != len(td.Fields) {
				g.codegenError(n.Line, n.Col, fmt.Sprintf("type %s constructor: wrong field count", n.Name), "")
				return
			}
			for _, a := range n.Args {
				g.emitExpr(ch, a)
			}
			idx := ch.AddName(n.Name)
			ch.Emit(opcode.OpNewFilled, idx, uint8(len(n.Args)), n.Line)
			break
		}
		for _, a := range n.Args {
			g.emitExpr(ch, a)
		}
		idx := ch.AddName(n.Name)
		op := opcode.OpCallBuiltin
		if _, ok := g.Prog.Functions[n.Name]; ok {
			op = opcode.OpCallUser
		}
		ac := g.argCountFlags(len(n.Args), n.Line, n.Col)
		if g.err != nil {
			return
		}
		ch.Emit(op, idx, ac, n.Line)

	case *ast.HandleCallExpr:
		// 1. Push self (the handle)
		sym := g.Symbols.Resolve(n.Receiver)
		if sym != nil && (sym.Kind == symtable.Local || sym.Kind == symtable.Param) {
			ch.Emit(opcode.OpLoadLocal, int32(sym.Slot), 0, n.Line)
		} else {
			idx := ch.AddName(n.Receiver)
			ch.Emit(opcode.OpLoadGlobal, idx, 0, n.Line)
		}
		// 2. Push arguments
		for _, a := range n.Args {
			g.emitExpr(ch, a)
		}
		// 3. Resolve method name and emit CallHandle
		midx := ch.AddName(n.Method)
		ac := g.argCountFlags(len(n.Args), n.Line, n.Col)
		if g.err != nil {
			return
		}
		ch.Emit(opcode.OpCallHandle, midx, ac, n.Line)

	case *ast.NewNode:
		idx := ch.AddName(n.TypeName)
		ch.Emit(opcode.OpNew, idx, 0, n.Line)

	case *ast.FieldAccessNode:
		// 1. Load Receiver handle
		sym := g.Symbols.Resolve(n.Object)
		if sym != nil && (sym.Kind == symtable.Local || sym.Kind == symtable.Param) {
			ch.Emit(opcode.OpLoadLocal, int32(sym.Slot), 0, n.Line)
		} else {
			idx := ch.AddName(n.Object)
			ch.Emit(opcode.OpLoadGlobal, idx, 0, n.Line)
		}
		// 2. Emit FieldGet
		fidx := ch.AddName(n.Field)
		ch.Emit(opcode.OpFieldGet, fidx, 0, n.Line)

	case *ast.IndexFieldExpr:
		g.emitLoadNamed(ch, n.Array, n.Line)
		for _, ix := range n.Index {
			g.emitExpr(ch, ix)
		}
		ch.Emit(opcode.OpArrayGet, int32(len(n.Index)), 0, n.Line)
		fidx := ch.AddName(n.Field)
		ch.Emit(opcode.OpFieldGet, fidx, 0, n.Line)

	case *ast.NamespaceCallExpr:
		g.emitNamespaceCallExpr(ch, n)

	case *ast.GroupedExpr:
		g.emitExpr(ch, n.Inner)

	case *ast.UnaryNode:
		g.emitExpr(ch, n.Expr)
		switch n.Op {
		case "-":
			ch.Emit(opcode.OpNeg, 0, 0, n.Line)
		case "NOT":
			ch.Emit(opcode.OpNot, 0, 0, n.Line)
		}

	case *ast.IndexExpr:
		g.emitExpr(ch, n.Base)
		for _, idx := range n.Index {
			g.emitExpr(ch, idx)
		}
		ch.Emit(opcode.OpArrayGet, int32(len(n.Index)), 0, n.Line)

	default:
		g.codegenError(1, 1, fmt.Sprintf("unsupported expression: %T", e),
			"This expression type is not yet implemented in the bytecode backend.")
	}
}

func (g *CodeGen) emitNamespaceCallExpr(ch *opcode.Chunk, n *ast.NamespaceCallExpr) {
	// 1. Push arguments
	for _, a := range n.Args {
		g.emitExpr(ch, a)
	}

	// 2. Resolve name NS.METHOD
	idx := ch.AddName(n.NS + "." + n.Method)

	// 3. Emit Call
	ac := g.argCountFlags(len(n.Args), n.Line, n.Col)
	if g.err != nil {
		return
	}
	ch.Emit(opcode.OpCallBuiltin, idx, ac, n.Line)
}
