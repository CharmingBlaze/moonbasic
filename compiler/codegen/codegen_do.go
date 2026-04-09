package codegen

import (
	"moonbasic/compiler/ast"
	"moonbasic/vm/opcode"
)

func (g *CodeGen) emitDoLoop(ch *opcode.Chunk, n *ast.DoLoopNode) {
	switch n.Kind {
	case ast.DoPreWhile:
		condStart := int32(len(ch.Instructions))
		g.nextReg = g.baseReg
		condReg := g.emitExpr(ch, n.Cond)
		// OpJumpIfFalse: dst=0, srcA=condReg
		exitJ := ch.Emit(opcode.OpJumpIfFalse, 0, condReg, 0, 0, n.Line)
		g.nextReg = g.baseReg
		
		g.beginLoop("do", condStart)
		for _, st := range n.Body {
			g.emitStmt(ch, st)
		}
		ch.Emit(opcode.OpJump, 0, 0, 0, condStart, n.Line)
		leave := len(ch.Instructions)
		g.endLoop(ch, leave)
		ch.Instructions[exitJ].Operand = int32(leave)
		
	case ast.DoPostWhile:
		bodyStart := int32(len(ch.Instructions))
		g.beginLoop("do", bodyStart)
		for _, st := range n.Body {
			g.emitStmt(ch, st)
		}
		
		g.nextReg = g.baseReg
		condReg := g.emitExpr(ch, n.Cond)
		// Loop while true: Jump if true
		ch.Emit(opcode.OpJumpIfTrue, 0, condReg, 0, bodyStart, n.Line)
		
		leave := len(ch.Instructions)
		g.endLoop(ch, leave)
		
	case ast.DoPostUntil:
		bodyStart := int32(len(ch.Instructions))
		g.beginLoop("do", bodyStart)
		for _, st := range n.Body {
			g.emitStmt(ch, st)
		}
		
		g.nextReg = g.baseReg
		condReg := g.emitExpr(ch, n.Cond)
		// Loop UNTIL true: Jump if false back to start
		ch.Emit(opcode.OpJumpIfFalse, 0, condReg, 0, bodyStart, n.Line)
		
		leave := len(ch.Instructions)
		g.endLoop(ch, leave)
	}
}
