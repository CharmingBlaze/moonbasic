package codegen

import (
	"moonbasic/compiler/ast"
	"moonbasic/vm/opcode"
)

func (g *CodeGen) emitSelect(ch *opcode.Chunk, n *ast.SelectNode) {
	g.nextReg = g.baseReg
	// Evaluate the selector expression
	selReg := g.emitExpr(ch, n.Expr)
	
	// We need to keep selReg stable across all cases.
	// We'll move it to a dedicated register at the start of our temporary pool.
	selector := g.allocReg()
	ch.Emit(opcode.OpMove, selector, selReg, 0, 0, n.Line)
	
	// Shift baseReg so statement emission within cases doesn't overwrite our selector
	oldBase := g.baseReg
	g.baseReg = g.nextReg

	var endJumps []int
	for _, c := range n.Cases {
		g.nextReg = g.baseReg
		caseValReg := g.emitExpr(ch, c.Value)
		resReg := g.allocReg()
		ch.Emit(opcode.OpEq, resReg, selector, caseValReg, 0, n.Line)
		
		skip := ch.Emit(opcode.OpJumpIfFalse, 0, resReg, 0, 0, n.Line)
		
		// Case body
		for _, st := range c.Body {
			g.emitStmt(ch, st)
		}
		endJumps = append(endJumps, ch.Emit(opcode.OpJump, 0, 0, 0, 0, n.Line))
		ch.Instructions[skip].Operand = int32(len(ch.Instructions))
	}
	
	for _, st := range n.Default {
		g.emitStmt(ch, st)
	}
	
	endIP := int32(len(ch.Instructions))
	for _, j := range endJumps {
		ch.Instructions[j].Operand = endIP
	}
	
	g.baseReg = oldBase
	g.nextReg = g.baseReg
}
