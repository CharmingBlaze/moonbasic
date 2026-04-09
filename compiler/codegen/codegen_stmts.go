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
		g.nextReg = g.baseReg // reset temporaries
		if n.Global {
			g.Symbols.DefineGlobalVar(n.Name)
		}
		valReg := g.emitExpr(ch, n.Expr)
		var sym *symtable.Symbol
		if n.Global {
			sym = g.Symbols.Resolve(n.Name)
		} else {
			sym = g.resolveOrDefineAssignTarget(n.Name)
		}
		if sym != nil && (sym.Kind == symtable.Local || sym.Kind == symtable.Param) {
			// Move value to local register
			ch.Emit(opcode.OpMove, uint8(sym.Slot), valReg, 0, 0, n.Line)
		} else if sym != nil && sym.Kind == symtable.Static {
			k := ch.AddName(sym.StaticKey)
			ch.Emit(opcode.OpStoreGlobal, 0, valReg, 0, k, n.Line)
		} else {
			idx := ch.AddName(n.Name)
			ch.Emit(opcode.OpStoreGlobal, 0, valReg, 0, idx, n.Line)
		}

	case *ast.NamespaceCallStmt:
		g.emitNamespaceCallStmt(ch, n)

	case *ast.HandleCallStmt:
		g.emitHandleCallStmt(ch, n)

	case *ast.FieldAssignNode:
		g.emitFieldAssign(ch, n)

	case *ast.IndexFieldAssignNode:
		g.emitIndexFieldAssign(ch, n)

	case *ast.DeleteStmt:
		g.nextReg = g.baseReg
		hReg := g.emitExpr(ch, n.Expr)
		ch.Emit(opcode.OpDelete, 0, hReg, 0, 0, n.Line)

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
		g.nextReg = g.baseReg
		g.Symbols.DefineLocal(n.Name)
		var valReg uint8
		if n.Init != nil {
			valReg = g.emitExpr(ch, n.Init)
		} else {
			valReg = g.allocReg()
			ch.Emit(opcode.OpPushNull, valReg, 0, 0, 0, n.Line)
		}
		sym := g.Symbols.Resolve(n.Name)
		if sym != nil && sym.Kind == symtable.Local {
			ch.Emit(opcode.OpMove, uint8(sym.Slot), valReg, 0, 0, n.Line)
		} else {
			g.codegenError(n.Line, n.Col, "LOCAL internal error", "")
			return
		}
		g.nextReg = g.baseReg

	case *ast.ConstDeclNode:
		g.nextReg = g.baseReg
		valReg := g.emitExpr(ch, n.Expr)
		idx := ch.AddName(n.Name)
		ch.Emit(opcode.OpStoreGlobal, 0, valReg, 0, idx, n.Line)
		g.nextReg = g.baseReg

	case *ast.StaticDeclNode:
		g.nextReg = g.baseReg
		sym := g.Symbols.DefineStatic(g.funcName, n.Name)
		var valReg uint8
		if n.Init != nil {
			valReg = g.emitExpr(ch, n.Init)
		} else {
			valReg = g.allocReg()
			ch.Emit(opcode.OpPushNull, valReg, 0, 0, 0, n.Line)
		}
		k := ch.AddName(sym.StaticKey)
		ch.Emit(opcode.OpStoreGlobal, 0, valReg, 0, k, n.Line)
		g.nextReg = g.baseReg

	case *ast.SwapStmt:
		g.nextReg = g.baseReg
		// In a register VM, SWAP usually applies to registers.
		// For variables, we resolve their slots and emit OpSwap.
		symA := g.Symbols.Resolve(n.A)
		symB := g.Symbols.Resolve(n.B)
		
		// This is tricky if they are mixed local/global.
		// For now assume locals.
		if symA != nil && symB != nil && (symA.Kind == symtable.Local || symA.Kind == symtable.Param) && (symB.Kind == symtable.Local || symB.Kind == symtable.Param) {
			ch.Emit(opcode.OpSwap, 0, uint8(symA.Slot), uint8(symB.Slot), 0, n.Line)
		}

	case *ast.EraseStmt:
		g.nextReg = g.baseReg
		if strings.EqualFold(n.Name, "ALL") {
			ch.Emit(opcode.OpEraseAll, 0, 0, 0, 0, n.Line)
			break
		}
		r := g.emitLoadNamed(ch, n.Name, n.Line)
		ch.Emit(opcode.OpDelete, 0, r, 0, 0, n.Line)
		
		// Set it to null
		nullReg := g.allocReg()
		ch.Emit(opcode.OpPushNull, nullReg, 0, 0, 0, n.Line)
		g.emitStoreNamed(ch, n.Name, n.Line, nullReg)
		g.nextReg = g.baseReg

	case *ast.IndexAssignNode:
		g.nextReg = g.baseReg
		arrReg := g.emitExpr(ch, &ast.IdentNode{Name: n.Array, Line: n.Line})
		dimStart := g.nextReg
		for _, ix := range n.Index {
			g.emitExpr(ch, ix)
		}
		valReg := g.emitExpr(ch, n.Expr)
		ch.Emit(opcode.OpArraySet, valReg, arrReg, dimStart, int32(len(n.Index)), n.Line)

	case *ast.ReturnNode:
		g.nextReg = g.baseReg
		if n.Expr != nil {
			r := g.emitExpr(ch, n.Expr)
			ch.Emit(opcode.OpReturn, 0, r, 0, 1, n.Line)
		} else {
			ch.Emit(opcode.OpReturnVoid, 0, 0, 0, 0, n.Line)
		}

	case *ast.EndProgramStmt:
		ch.Emit(opcode.OpHalt, 0, 0, 0, 0, n.Line)

	default:
		g.codegenError(1, 1, fmt.Sprintf("unsupported statement for codegen: %T", s),
			"This statement is not yet implemented in the bytecode backend.")
	}
}

func (g *CodeGen) emitLoadNamed(ch *opcode.Chunk, name string, line int) uint8 {
	sym := g.Symbols.Resolve(name)
	if sym != nil && (sym.Kind == symtable.Local || sym.Kind == symtable.Param) {
		// Already in a register, but we MUST move it to a temporary if it's used in an expression
		// that might overwrite it (like a Binop reusing the register).
		// However, for simplicity now, let's just return the slot and be careful in emitExpr.
		// Wait, emitExpr ALREADY does the OpMove for locals! 
		// So emitLoadNamed is usually used by statements that want the value.
		r := g.allocReg()
		ch.Emit(opcode.OpMove, r, uint8(sym.Slot), 0, 0, line)
		return r
	}
	if sym != nil && sym.Kind == symtable.Static {
		k := ch.AddName(sym.StaticKey)
		r := g.allocReg()
		ch.Emit(opcode.OpLoadGlobal, r, 0, 0, k, line)
		return r
	}
	idx := ch.AddName(name)
	r := g.allocReg()
	ch.Emit(opcode.OpLoadGlobal, r, 0, 0, idx, line)
	return r
}

func (g *CodeGen) emitStoreNamed(ch *opcode.Chunk, name string, line int, src uint8) {
	sym := g.Symbols.Resolve(name)
	if sym != nil && (sym.Kind == symtable.Local || sym.Kind == symtable.Param) {
		ch.Emit(opcode.OpMove, uint8(sym.Slot), src, 0, 0, line)
		return
	}
	if sym != nil && sym.Kind == symtable.Static {
		k := ch.AddName(sym.StaticKey)
		ch.Emit(opcode.OpStoreGlobal, 0, src, 0, k, line)
		return
	}
	idx := ch.AddName(name)
	ch.Emit(opcode.OpStoreGlobal, 0, src, 0, idx, line)
}

func (g *CodeGen) emitIf(ch *opcode.Chunk, n *ast.IfNode) {
	g.nextReg = g.baseReg
	condReg := g.emitExpr(ch, n.Cond)
	// Placeholder jump — we'll patch it later. dst=0, srcA=condReg.
	jumpIdx := ch.Emit(opcode.OpJumpIfFalse, 0, condReg, 0, 0, n.Line)
	g.nextReg = g.baseReg // Clear temporaries after condition

	// THEN block
	for _, st := range n.Then {
		g.emitStmt(ch, st)
	}

	// Patch the jump to point here (end of IF)
	ch.Instructions[jumpIdx].Operand = int32(len(ch.Instructions))
}

func (g *CodeGen) emitWhile(ch *opcode.Chunk, n *ast.WhileNode) {
	startIdx := int32(len(ch.Instructions))
	g.nextReg = g.baseReg
	condReg := g.emitExpr(ch, n.Cond)
	exitJump := ch.Emit(opcode.OpJumpIfFalse, 0, condReg, 0, 0, n.Line)
	g.nextReg = g.baseReg
	
	g.beginLoop("while", startIdx)
	for _, st := range n.Body {
		g.emitStmt(ch, st)
	}
	ch.Emit(opcode.OpJump, 0, 0, 0, startIdx, n.Line)
	leave := len(ch.Instructions)
	g.endLoop(ch, leave)
	ch.Instructions[exitJump].Operand = int32(leave)
}

func (g *CodeGen) emitFor(ch *opcode.Chunk, n *ast.ForNode) {
	g.nextReg = g.baseReg
	cmpOp, err := g.forLoopCompareOp(n.Step)
	if err != nil {
		g.codegenError(n.Line, n.Col, err.Error(), "Use a numeric literal STEP (e.g. STEP -2).")
		return
	}

	// 1. Initial Assignment: var = from
	fromReg := g.emitExpr(ch, n.From)
	sym := g.Symbols.Resolve(n.Var)
	if sym != nil && (sym.Kind == symtable.Local || sym.Kind == symtable.Param) {
		ch.Emit(opcode.OpMove, uint8(sym.Slot), fromReg, 0, 0, n.Line)
	} else {
		idx := ch.AddName(n.Var)
		ch.Emit(opcode.OpStoreGlobal, 0, fromReg, 0, idx, n.Line)
	}
	g.nextReg = g.baseReg // var is set, clear temps

	// 2. Loop Header
	startIdx := len(ch.Instructions)

	// 3. Condition: var <= to
	varReg := g.emitLoadNamed(ch, n.Var, n.Line)
	toReg := g.emitExpr(ch, n.To)
	resReg := g.allocReg()
	ch.Emit(cmpOp, resReg, varReg, toReg, 0, n.Line)

	exitJump := ch.Emit(opcode.OpJumpIfFalse, 0, resReg, 0, 0, n.Line)

	g.beginLoop("for", -1)
	g.nextReg = g.baseReg

	// 4. Body
	for _, st := range n.Body {
		g.emitStmt(ch, st)
	}

	contIdx := int32(len(ch.Instructions))
	g.setLoopContinueTarget(contIdx)

	// 5. Increment: var = var + step
	g.nextReg = g.baseReg
	curVarReg := g.emitLoadNamed(ch, n.Var, n.Line)
	var stepReg uint8
	if n.Step != nil {
		stepReg = g.emitExpr(ch, n.Step)
	} else {
		idx := ch.AddInt(1)
		stepReg = g.allocReg()
		ch.Emit(opcode.OpPushInt, stepReg, 0, 0, idx, n.Line)
	}
	
	newVarReg := g.allocReg()
	ch.Emit(opcode.OpAdd, newVarReg, curVarReg, stepReg, 0, n.Line)
	
	if sym != nil && (sym.Kind == symtable.Local || sym.Kind == symtable.Param) {
		ch.Emit(opcode.OpMove, uint8(sym.Slot), newVarReg, 0, 0, n.Line)
	} else {
		idx := ch.AddName(n.Var)
		ch.Emit(opcode.OpStoreGlobal, 0, newVarReg, 0, idx, n.Line)
	}
	g.nextReg = g.baseReg

	// 6. Jump back
	ch.Emit(opcode.OpJump, 0, 0, 0, int32(startIdx), n.Line)

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
	g.nextReg = g.baseReg
	condReg := g.emitExpr(ch, n.Condition)
	// REPEAT UNTIL Cond is true? or while cond is false?
	// MoonBASIC REPEAT UNTIL follows QBasic: loop while condition is FALSE.
	ch.Emit(opcode.OpJumpIfFalse, 0, condReg, 0, startIdx, n.Line)
	leave := len(ch.Instructions)
	g.endLoop(ch, leave)
}

func (g *CodeGen) emitNamespaceCallStmt(ch *opcode.Chunk, n *ast.NamespaceCallStmt) {
	g.nextReg = g.baseReg
	argStart := g.nextReg
	for _, a := range n.Args {
		g.emitExpr(ch, a)
	}

	idx := ch.AddName(n.NS + "." + n.Method)
	// Statements don't use the result, but builtins expect a Dst. We provide temporary.
	dst := g.allocReg()
	operand := (int32(len(n.Args)) << 24) | idx
	ch.Emit(opcode.OpCallBuiltin, dst, 0, argStart, operand, n.Line)
	g.nextReg = g.baseReg
}

func (g *CodeGen) emitIndexFieldAssign(ch *opcode.Chunk, n *ast.IndexFieldAssignNode) {
	g.nextReg = g.baseReg
	arrReg := g.emitLoadNamed(ch, n.Array, n.Line)
	dimStart := g.nextReg
	for _, ix := range n.Index {
		g.emitExpr(ch, ix)
	}
	objReg := g.allocReg()
	ch.Emit(opcode.OpArrayGet, objReg, arrReg, dimStart, int32(len(n.Index)), n.Line)
	
	valReg := g.emitExpr(ch, n.Expr)
	fidx := ch.AddName(n.Field)
	ch.Emit(opcode.OpFieldSet, 0, objReg, valReg, fidx, n.Line)
	g.nextReg = g.baseReg
}

func (g *CodeGen) emitFieldAssign(ch *opcode.Chunk, n *ast.FieldAssignNode) {
	g.nextReg = g.baseReg
	recReg := g.emitLoadNamed(ch, n.Object, n.Line)
	valReg := g.emitExpr(ch, n.Expr)
	fidx := ch.AddName(n.Field)
	ch.Emit(opcode.OpFieldSet, 0, recReg, valReg, fidx, n.Line)
	g.nextReg = g.baseReg
}

func (g *CodeGen) emitHandleCallStmt(ch *opcode.Chunk, n *ast.HandleCallStmt) {
	g.nextReg = g.baseReg
	recReg := g.emitLoadNamed(ch, n.Receiver, n.Line)
	argStart := g.nextReg
	for _, a := range n.Args {
		g.emitExpr(ch, a)
	}

	midx := ch.AddName(n.Method)
	dst := g.allocReg() // discard result
	operand := (int32(len(n.Args)) << 24) | midx
	ch.Emit(opcode.OpCallHandle, dst, recReg, argStart, operand, n.Line)
	g.nextReg = g.baseReg
}
