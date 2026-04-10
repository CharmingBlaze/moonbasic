package codegen

import (
	"fmt"
	"strings"

	"moonbasic/compiler/ast"
	"moonbasic/compiler/symtable"
	"moonbasic/vm/opcode"
)

// emitExpr translates an expression AST node into bytecode and returns the register index.
func (g *CodeGen) emitExpr(ch *opcode.Chunk, e ast.Expr) uint8 {
	if g.err != nil {
		return 0
	}
	switch n := e.(type) {
	case *ast.IntLitNode:
		idx := ch.AddInt(n.Value)
		r := g.allocReg()
		ch.Emit(opcode.OpPushInt, r, 0, 0, idx, n.Line)
		return r

	case *ast.FloatLitNode:
		idx := ch.AddFloat(n.Value)
		r := g.allocReg()
		ch.Emit(opcode.OpPushFloat, r, 0, 0, idx, n.Line)
		return r

	case *ast.StringLitNode:
		idx := g.Prog.InternString(n.Value)
		r := g.allocReg()
		ch.Emit(opcode.OpPushString, r, 0, 0, idx, n.Line)
		return r

	case *ast.BoolLitNode:
		v := int32(0)
		if n.Value {
			v = 1
		}
		r := g.allocReg()
		ch.Emit(opcode.OpPushBool, r, 0, 0, v, n.Line)
		return r

	case *ast.NullLitNode:
		r := g.allocReg()
		ch.Emit(opcode.OpPushNull, r, 0, 0, 0, n.Line)
		return r

	case *ast.BinopNode:
		L := g.emitExpr(ch, n.Left)
		R := g.emitExpr(ch, n.Right)
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
			return 0
		}
		// Result goes into L (reuse register)
		ch.Emit(op, L, L, R, 0, n.Line)
		// R is no longer needed
		g.nextReg = R
		return L

	case *ast.IdentNode:
		sym := g.Symbols.Resolve(n.Name)
		if sym != nil && (sym.Kind == symtable.Local || sym.Kind == symtable.Param) {
			// If it's a local, it's already in a register. We return that register.
			// BE CAREFUL: Binops reuse the register. If we return a local register, 
			// the next operation might overwrite it!
			// We MUST Move it to a temporary register if we want to protect the local.
			r := g.allocReg()
			// Need an OpMove.
			ch.Emit(opcode.OpMove, r, uint8(sym.Slot), 0, 0, n.Line)
			return r
		} else if sym != nil && sym.Kind == symtable.Static {
			k := ch.AddName(sym.StaticKey)
			r := g.allocReg()
			ch.Emit(opcode.OpLoadGlobal, r, 0, 0, k, n.Line)
			return r
		} else {
			idx := ch.AddName(n.Name)
			r := g.allocReg()
			ch.Emit(opcode.OpLoadGlobal, r, 0, 0, idx, n.Line)
			return r
		}

	case *ast.CallExprNode:
		// Constructors
		if td, ok := g.Prog.Types[strings.ToUpper(n.Name)]; ok && g.Prog.Functions[strings.ToUpper(n.Name)] == nil {
			if len(n.Args) != len(td.Fields) {
				g.codegenError(n.Line, n.Col, fmt.Sprintf("type %s constructor: wrong field count", n.Name), "")
				return 0
			}
			argStart := g.emitArgsStable(ch, n.Args, n.Line)
			idx := ch.AddName(strings.ToUpper(n.Name))
			dst := g.allocReg()
			// OpNewFilled Dst, FieldCount, ArgStart, NameIdx
			ch.Emit(opcode.OpNewFilled, dst, uint8(len(n.Args)), argStart, idx, n.Line)
			g.nextReg = argStart // Free temporary registers used for args
			g.allocReg() // Re-allocate dst in the pool
			return dst
		}
		
		// Normal Calls
		argStart := g.emitArgsStable(ch, n.Args, n.Line)
		idx := ch.AddName(n.Name)
		op := opcode.OpCallBuiltin
		if _, ok := g.Prog.Functions[n.Name]; ok {
			op = opcode.OpCallUser
		}
		
		dst := g.allocReg()
		// Encode (ArgCount << 24 | NameIdx) into Operand
		operand := (int32(len(n.Args)) << 24) | idx
		ch.Emit(op, dst, 0, argStart, operand, n.Line)
		
		g.nextReg = argStart
		g.allocReg() // keep dst
		return dst

	case *ast.HandleCallExpr:
		// Evaluate receiver into a register
		recReg := g.emitExpr(ch, &ast.IdentNode{Name: n.Receiver, Line: n.Line, Col: n.Col})
		
		argStart := g.emitArgsStable(ch, n.Args, n.Line)
		midx := ch.AddName(n.Method)
		dst := g.allocReg()
		
		operand := (int32(len(n.Args)) << 24) | midx
		ch.Emit(opcode.OpCallHandle, dst, recReg, argStart, operand, n.Line)
		
		g.nextReg = recReg // and free everything after receiver
		g.allocReg() // keep dst
		return dst

	case *ast.NewNode:
		idx := ch.AddName(strings.ToUpper(n.TypeName))
		dst := g.allocReg()
		ch.Emit(opcode.OpNew, dst, 0, 0, idx, n.Line)
		return dst

	case *ast.FieldAccessNode:
		if strings.EqualFold(n.Field, "length") {
			recReg := g.emitExpr(ch, &ast.IdentNode{Name: n.Object, Line: n.Line, Col: n.Col})
			dst := g.allocReg()
			ch.Emit(opcode.OpArrayLen, dst, recReg, 0, 0, n.Line)
			g.nextReg = recReg
			g.allocReg()
			return dst
		}
		recReg := g.emitExpr(ch, &ast.IdentNode{Name: n.Object, Line: n.Line, Col: n.Col})
		fidx := ch.AddName(strings.ToUpper(n.Field))
		dst := g.allocReg()
		ch.Emit(opcode.OpFieldGet, dst, recReg, 0, fidx, n.Line)

		g.nextReg = recReg
		g.allocReg() // keep dst
		return dst

	case *ast.UnaryNode:
		valReg := g.emitExpr(ch, n.Expr)
		dst := g.allocReg()
		switch n.Op {
		case "-":
			ch.Emit(opcode.OpNeg, dst, valReg, 0, 0, n.Line)
		case "NOT":
			ch.Emit(opcode.OpNot, dst, valReg, 0, 0, n.Line)
		}
		g.nextReg = valReg
		g.allocReg()
		return dst

	case *ast.IndexExpr:
		baseReg := g.emitExpr(ch, n.Base)
		dimStart := g.nextReg
		for _, idx := range n.Index {
			g.emitExpr(ch, idx)
		}
		dst := g.allocReg()
		ch.Emit(opcode.OpArrayGet, dst, baseReg, dimStart, int32(len(n.Index)), n.Line)
		
		g.nextReg = baseReg
		g.allocReg()
		return dst

	case *ast.GroupedExpr:
		return g.emitExpr(ch, n.Inner)

	case *ast.NamespaceCallExpr:
		// NS.METHOD(...)
		argStart := g.emitArgsStable(ch, n.Args, n.Line)
		idx := ch.AddName(n.NS + "." + n.Method)
		dst := g.allocReg()
		
		operand := (int32(len(n.Args)) << 24) | idx
		ch.Emit(opcode.OpCallBuiltin, dst, 0, argStart, operand, n.Line)
		
		g.nextReg = argStart
		g.allocReg()
		return dst

	case *ast.IndexFieldExpr:
		// Array(idx).Field
		// 1. Get array handle
		arrReg := g.emitExpr(ch, &ast.IdentNode{Name: n.Array, Line: n.Line})
		// 2. Index into it
		dimStart := g.nextReg
		for _, ix := range n.Index {
			g.emitExpr(ch, ix)
		}
		objReg := g.allocReg()
		ch.Emit(opcode.OpArrayGet, objReg, arrReg, dimStart, int32(len(n.Index)), n.Line)
		// 3. Get field
		fidx := ch.AddName(strings.ToUpper(n.Field))
		dst := g.allocReg()
		ch.Emit(opcode.OpFieldGet, dst, objReg, 0, fidx, n.Line)
		
		g.nextReg = arrReg
		g.allocReg()
		return dst

	default:
		g.codegenError(1, 1, fmt.Sprintf("unsupported expression: %T", e),
			"This expression type is not yet implemented in the bytecode backend.")
		return 0
	}
}
