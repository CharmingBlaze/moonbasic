package codegen

import (
	"fmt"

	"moonbasic/compiler/ast"
	"moonbasic/vm/opcode"
)

func (g *CodeGen) emitSelect(ch *opcode.Chunk, n *ast.SelectNode) {
	g.selectTmpID++
	tmp := fmt.Sprintf("__SEL_%d", g.selectTmpID)

	g.emitExpr(ch, n.Expr)
	tidx := ch.AddName(tmp)
	ch.Emit(opcode.OpStoreGlobal, tidx, 0, n.Line)
	ch.Emit(opcode.OpPop, 0, 0, n.Line)

	var endJumps []int
	for _, c := range n.Cases {
		ch.Emit(opcode.OpLoadGlobal, tidx, 0, n.Line)
		g.emitExpr(ch, c.Value)
		ch.Emit(opcode.OpEq, 0, 0, n.Line)
		skip := ch.Emit(opcode.OpJumpIfFalse, 0, 0, n.Line)
		for _, st := range c.Body {
			g.emitStmt(ch, st)
		}
		endJumps = append(endJumps, ch.Emit(opcode.OpJump, 0, 0, n.Line))
		ch.Instructions[skip].Operand = int32(len(ch.Instructions))
	}
	for _, st := range n.Default {
		g.emitStmt(ch, st)
	}
	endIP := len(ch.Instructions)
	for _, j := range endJumps {
		ch.Instructions[j].Operand = int32(endIP)
	}
}
