package codegen

import (
	"errors"
	"fmt"
	"strings"

	"moonbasic/compiler/ast"
	"moonbasic/compiler/symtable"
	"moonbasic/vm/opcode"
)

func (g *CodeGen) forLoopCompareOp(step ast.Expr) (opcode.OpCode, error) {
	if step == nil {
		return opcode.OpLte, nil
	}
	switch s := step.(type) {
	case *ast.IntLitNode:
		if s.Value < 0 {
			return opcode.OpGte, nil
		}
		return opcode.OpLte, nil
	case *ast.FloatLitNode:
		if s.Value < 0 {
			return opcode.OpGte, nil
		}
		return opcode.OpLte, nil
	default:
		return 0, errors.New("[moonBASIC] Error: dynamic STEP not supported in v1.x — use a constant STEP value")
	}
}

// emitStmt translates a statement AST node into bytecode.
func (g *CodeGen) emitStmt(ch *opcode.Chunk, s ast.Stmt) {
	if g.err != nil {
		return
	}
	switch n := s.(type) {
	case *ast.AssignNode:
		if n.Global {
			g.Symbols.DefineGlobalVar(n.Name)
		}
		g.emitExpr(ch, n.Expr)
		var sym *symtable.Symbol
		if n.Global {
			sym = g.Symbols.Resolve(n.Name)
		} else {
			sym = g.resolveOrDefineAssignTarget(n.Name)
		}
		if sym != nil && (sym.Kind == symtable.Local || sym.Kind == symtable.Param) {
			ch.Emit(opcode.OpStoreLocal, int32(sym.Slot), 0, n.Line)
		} else if sym != nil && sym.Kind == symtable.Static {
			k := ch.AddName(sym.StaticKey)
			ch.Emit(opcode.OpStoreGlobal, k, 0, n.Line)
		} else {
			idx := ch.AddName(n.Name)
			ch.Emit(opcode.OpStoreGlobal, idx, 0, n.Line)
		}
		ch.Emit(opcode.OpPop, 0, 0, n.Line)

	case *ast.NamespaceCallStmt:
		g.emitNamespaceCallStmt(ch, n)

	case *ast.HandleCallStmt:
		g.emitHandleCallStmt(ch, n)

	case *ast.FieldAssignNode:
		g.emitFieldAssign(ch, n)

	case *ast.IndexFieldAssignNode:
		g.emitIndexFieldAssign(ch, n)

	case *ast.DeleteStmt:
		g.emitExpr(ch, n.Expr) // Push handle
		ch.Emit(opcode.OpDelete, 0, 0, n.Line)

	case *ast.CallStmtNode:
		g.emitCallStmt(ch, n)

	case *ast.IfNode:
		g.emitIf(ch, n)

	case *ast.WhileNode:
		g.emitWhile(ch, n)

	case *ast.ForNode:
		g.emitFor(ch, n)

	case *ast.RepeatNode:
		g.emitRepeat(ch, n)

	case *ast.DoLoopNode:
		g.emitDoLoop(ch, n)

	case *ast.ExitStmt:
		g.emitExitStmt(ch, n)

	case *ast.ContinueStmt:
		g.emitContinueStmt(ch, n)

	case *ast.SelectNode:
		g.emitSelect(ch, n)

	case *ast.DimNode:
		g.emitDim(ch, n)

	case *ast.LocalDeclNode:
		g.Symbols.DefineLocal(n.Name)
		if n.Init != nil {
			g.emitExpr(ch, n.Init)
		} else {
			ch.Emit(opcode.OpPushNull, 0, 0, n.Line)
		}
		sym := g.Symbols.Resolve(n.Name)
		if sym != nil && sym.Kind == symtable.Local {
			ch.Emit(opcode.OpStoreLocal, int32(sym.Slot), 0, n.Line)
		} else {
			g.codegenError(n.Line, n.Col, "LOCAL internal error", "")
			return
		}
		ch.Emit(opcode.OpPop, 0, 0, n.Line)

	case *ast.ConstDeclNode:
		g.emitExpr(ch, n.Expr)
		idx := ch.AddName(n.Name)
		ch.Emit(opcode.OpStoreGlobal, idx, 0, n.Line)
		ch.Emit(opcode.OpPop, 0, 0, n.Line)

	case *ast.StaticDeclNode:
		sym := g.Symbols.DefineStatic(g.funcName, n.Name)
		if n.Init != nil {
			g.emitExpr(ch, n.Init)
		} else {
			ch.Emit(opcode.OpPushNull, 0, 0, n.Line)
		}
		k := ch.AddName(sym.StaticKey)
		ch.Emit(opcode.OpStoreGlobal, k, 0, n.Line)
		ch.Emit(opcode.OpPop, 0, 0, n.Line)

	case *ast.SwapStmt:
		g.emitLoadNamed(ch, n.A, n.Line)
		g.emitLoadNamed(ch, n.B, n.Line)
		ch.Emit(opcode.OpSwap, 0, 0, n.Line)
		g.emitStoreNamed(ch, n.B, n.Line)
		ch.Emit(opcode.OpPop, 0, 0, n.Line)
		g.emitStoreNamed(ch, n.A, n.Line)
		ch.Emit(opcode.OpPop, 0, 0, n.Line)

	case *ast.EraseStmt:
		if strings.EqualFold(n.Name, "ALL") {
			ch.Emit(opcode.OpEraseAll, 0, 0, n.Line)
			break
		}
		g.emitLoadNamed(ch, n.Name, n.Line)
		ch.Emit(opcode.OpDelete, 0, 0, n.Line)
		ch.Emit(opcode.OpPushNull, 0, 0, n.Line)
		g.emitStoreNamed(ch, n.Name, n.Line)
		ch.Emit(opcode.OpPop, 0, 0, n.Line)

	case *ast.IndexAssignNode:
		g.emitLoadNamed(ch, n.Array, n.Line)
		for _, ix := range n.Index {
			g.emitExpr(ch, ix)
		}
		g.emitExpr(ch, n.Expr)
		ch.Emit(opcode.OpArraySet, int32(len(n.Index)), 0, n.Line)
		ch.Emit(opcode.OpPop, 0, 0, n.Line)

	case *ast.ReturnNode:
		if n.Expr != nil {
			g.emitExpr(ch, n.Expr)
			ch.Emit(opcode.OpReturn, 1, 0, n.Line)
		} else {
			ch.Emit(opcode.OpReturnVoid, 0, 0, n.Line)
		}

	case *ast.EndProgramStmt:
		ch.Emit(opcode.OpHalt, 0, 0, n.Line)

	default:
		g.codegenError(1, 1, fmt.Sprintf("unsupported statement for codegen: %T", s),
			"This statement is not yet implemented in the bytecode backend.")
	}
}

func (g *CodeGen) emitLoadNamed(ch *opcode.Chunk, name string, line int) {
	sym := g.Symbols.Resolve(name)
	if sym != nil && (sym.Kind == symtable.Local || sym.Kind == symtable.Param) {
		ch.Emit(opcode.OpLoadLocal, int32(sym.Slot), 0, line)
		return
	}
	if sym != nil && sym.Kind == symtable.Static {
		k := ch.AddName(sym.StaticKey)
		ch.Emit(opcode.OpLoadGlobal, k, 0, line)
		return
	}
	idx := ch.AddName(name)
	ch.Emit(opcode.OpLoadGlobal, idx, 0, line)
}

func (g *CodeGen) emitStoreNamed(ch *opcode.Chunk, name string, line int) {
	sym := g.Symbols.Resolve(name)
	if sym != nil && (sym.Kind == symtable.Local || sym.Kind == symtable.Param) {
		ch.Emit(opcode.OpStoreLocal, int32(sym.Slot), 0, line)
		return
	}
	if sym != nil && sym.Kind == symtable.Static {
		k := ch.AddName(sym.StaticKey)
		ch.Emit(opcode.OpStoreGlobal, k, 0, line)
		return
	}
	idx := ch.AddName(name)
	ch.Emit(opcode.OpStoreGlobal, idx, 0, line)
}

func (g *CodeGen) emitIf(ch *opcode.Chunk, n *ast.IfNode) {
	g.emitExpr(ch, n.Cond)

	// Placeholder jump — we'll patch it later
	jumpIdx := ch.Emit(opcode.OpJumpIfFalse, 0, 0, n.Line)

	// THEN block
	for _, st := range n.Then {
		g.emitStmt(ch, st)
	}

	// Patch the jump to point here (end of IF)
	ch.Instructions[jumpIdx].Operand = int32(len(ch.Instructions))
}

func (g *CodeGen) emitWhile(ch *opcode.Chunk, n *ast.WhileNode) {
	startIdx := int32(len(ch.Instructions))
	g.emitExpr(ch, n.Cond)
	exitJump := ch.Emit(opcode.OpJumpIfFalse, 0, 0, n.Line)
	g.beginLoop("while", startIdx)
	for _, st := range n.Body {
		g.emitStmt(ch, st)
	}
	ch.Emit(opcode.OpJump, startIdx, 0, n.Line)
	leave := len(ch.Instructions)
	g.endLoop(ch, leave)
	ch.Instructions[exitJump].Operand = int32(leave)
}

func (g *CodeGen) emitFor(ch *opcode.Chunk, n *ast.ForNode) {
	cmpOp, err := g.forLoopCompareOp(n.Step)
	if err != nil {
		g.codegenError(n.Line, n.Col, err.Error(), "Use a numeric literal STEP (e.g. STEP -2).")
		return
	}

	// 1. Initial Assignment: var = from
	g.emitExpr(ch, n.From)
	sym := g.Symbols.Resolve(n.Var)
	if sym != nil && (sym.Kind == symtable.Local || sym.Kind == symtable.Param) {
		ch.Emit(opcode.OpStoreLocal, int32(sym.Slot), 0, n.Line)
	} else {
		idx := ch.AddName(n.Var)
		ch.Emit(opcode.OpStoreGlobal, idx, 0, n.Line)
	}
	ch.Emit(opcode.OpPop, 0, 0, n.Line)

	// 2. Loop Header
	startIdx := len(ch.Instructions)

	// 3. Condition: var <= to (step >= 0) or var >= to (constant step < 0)
	g.emitExpr(ch, &ast.IdentNode{Name: n.Var, Line: n.Line})
	g.emitExpr(ch, n.To)
	ch.Emit(cmpOp, 0, 0, n.Line)

	exitJump := ch.Emit(opcode.OpJumpIfFalse, 0, 0, n.Line)

	g.beginLoop("for", -1)

	// 4. Body
	for _, st := range n.Body {
		g.emitStmt(ch, st)
	}

	contIdx := int32(len(ch.Instructions))
	g.setLoopContinueTarget(contIdx)

	// 5. Increment: var = var + step
	g.emitExpr(ch, &ast.IdentNode{Name: n.Var, Line: n.Line})
	if n.Step != nil {
		g.emitExpr(ch, n.Step)
	} else {
		idx := ch.AddInt(1)
		ch.Emit(opcode.OpPushInt, idx, 0, n.Line)
	}
	ch.Emit(opcode.OpAdd, 0, 0, n.Line)

	if sym != nil && (sym.Kind == symtable.Local || sym.Kind == symtable.Param) {
		ch.Emit(opcode.OpStoreLocal, int32(sym.Slot), 0, n.Line)
	} else {
		idx := ch.AddName(n.Var)
		ch.Emit(opcode.OpStoreGlobal, idx, 0, n.Line)
	}
	ch.Emit(opcode.OpPop, 0, 0, n.Line)

	// 6. Jump back
	ch.Emit(opcode.OpJump, int32(startIdx), 0, n.Line)

	leave := len(ch.Instructions)
	g.endLoop(ch, leave)
	ch.Instructions[exitJump].Operand = int32(leave)
}

func (g *CodeGen) emitRepeat(ch *opcode.Chunk, n *ast.RepeatNode) {
	startIdx := int32(len(ch.Instructions))
	g.beginLoop("repeat", startIdx)
	for _, st := range n.Body {
		g.emitStmt(ch, st)
	}
	g.emitExpr(ch, n.Condition)
	ch.Emit(opcode.OpJumpIfFalse, startIdx, 0, n.Line)
	leave := len(ch.Instructions)
	g.endLoop(ch, leave)
}

func (g *CodeGen) emitNamespaceCallStmt(ch *opcode.Chunk, n *ast.NamespaceCallStmt) {
	// 1. Push arguments
	for _, a := range n.Args {
		g.emitExpr(ch, a)
	}

	// 2. Resolve name "NS.METHOD"
	idx := ch.AddName(n.NS + "." + n.Method)

	// 3. Emit Call
	ac := g.argCountFlags(len(n.Args), n.Line, n.Col)
	if g.err != nil {
		return
	}
	ch.Emit(opcode.OpCallBuiltin, idx, ac, n.Line)

	// 4. Pop return value (engine commands as statements must be clean)
	ch.Emit(opcode.OpPop, 0, 0, n.Line)
}

func (g *CodeGen) emitIndexFieldAssign(ch *opcode.Chunk, n *ast.IndexFieldAssignNode) {
	g.emitLoadNamed(ch, n.Array, n.Line)
	for _, ix := range n.Index {
		g.emitExpr(ch, ix)
	}
	ch.Emit(opcode.OpArrayGet, int32(len(n.Index)), 0, n.Line)
	g.emitExpr(ch, n.Expr)
	fidx := ch.AddName(n.Field)
	ch.Emit(opcode.OpFieldSet, fidx, 0, n.Line)
	ch.Emit(opcode.OpPop, 0, 0, n.Line)
}

func (g *CodeGen) emitFieldAssign(ch *opcode.Chunk, n *ast.FieldAssignNode) {
	// 1. Load Receiver handle
	sym := g.Symbols.Resolve(n.Object)
	if sym != nil && (sym.Kind == symtable.Local || sym.Kind == symtable.Param) {
		ch.Emit(opcode.OpLoadLocal, int32(sym.Slot), 0, n.Line)
	} else {
		idx := ch.AddName(n.Object)
		ch.Emit(opcode.OpLoadGlobal, idx, 0, n.Line)
	}

	// 2. Push Value
	g.emitExpr(ch, n.Expr)

	// 3. Emit FieldSet
	fidx := ch.AddName(n.Field)
	ch.Emit(opcode.OpFieldSet, fidx, 0, n.Line)

	// 4. Pop result (Assignment is a statement)
	ch.Emit(opcode.OpPop, 0, 0, n.Line)
}

func (g *CodeGen) emitHandleCallStmt(ch *opcode.Chunk, n *ast.HandleCallStmt) {
	// 1. Load Receiver handle (self)
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

	// 3. Emit CallHandle
	midx := ch.AddName(n.Method)
	ac := g.argCountFlags(len(n.Args), n.Line, n.Col)
	if g.err != nil {
		return
	}
	ch.Emit(opcode.OpCallHandle, midx, ac, n.Line)

	// 4. Pop return value
	ch.Emit(opcode.OpPop, 0, 0, n.Line)
}
