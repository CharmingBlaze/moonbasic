package codegen

import (
	"moonbasic/compiler/ast"
	"moonbasic/vm/opcode"
)

func (g *CodeGen) emitDoLoop(ch *opcode.Chunk, n *ast.DoLoopNode) {
	switch n.Kind {
	case ast.DoPreWhile:
		condStart := int32(len(ch.Instructions))
		g.emitExpr(ch, n.Cond)
		exitJ := ch.Emit(opcode.OpJumpIfFalse, 0, 0, n.Line)
		g.beginLoop("do", condStart)
		for _, st := range n.Body {
			g.emitStmt(ch, st)
		}
		ch.Emit(opcode.OpJump, condStart, 0, n.Line)
		leave := len(ch.Instructions)
		g.endLoop(ch, leave)
		ch.Instructions[exitJ].Operand = int32(leave)
	case ast.DoPostWhile:
		bodyStart := int32(len(ch.Instructions))
		g.beginLoop("do", bodyStart)
		for _, st := range n.Body {
			g.emitStmt(ch, st)
		}
		g.emitExpr(ch, n.Cond)
		exitJ := ch.Emit(opcode.OpJumpIfFalse, 0, 0, n.Line)
		ch.Emit(opcode.OpJump, bodyStart, 0, n.Line)
		leave := len(ch.Instructions)
		g.endLoop(ch, leave)
		ch.Instructions[exitJ].Operand = int32(leave)
	case ast.DoPostUntil:
		bodyStart := int32(len(ch.Instructions))
		g.beginLoop("do", bodyStart)
		for _, st := range n.Body {
			g.emitStmt(ch, st)
		}
		g.emitExpr(ch, n.Cond)
		exitJ := ch.Emit(opcode.OpJumpIfTrue, 0, 0, n.Line)
		ch.Emit(opcode.OpJump, bodyStart, 0, n.Line)
		leave := len(ch.Instructions)
		g.endLoop(ch, leave)
		ch.Instructions[exitJ].Operand = int32(leave)
	}
}
