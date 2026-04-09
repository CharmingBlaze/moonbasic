package codegen

import (
	"moonbasic/compiler/ast"
	"moonbasic/vm/opcode"
)

// emitCallStmt translates a call statement (builtin or user function) into bytecode.
func (g *CodeGen) emitCallStmt(ch *opcode.Chunk, n *ast.CallStmtNode) {
	if g.err != nil {
		return
	}
	g.nextReg = g.baseReg
	argStart := g.nextReg
	for _, a := range n.Args {
		g.emitExpr(ch, a)
	}

	idx := ch.AddName(n.Name)
	op := opcode.OpCallBuiltin
	if _, ok := g.Prog.Functions[n.Name]; ok {
		op = opcode.OpCallUser
	}

	dst := g.allocReg() // discard result
	operand := (int32(len(n.Args)) << 24) | idx
	ch.Emit(op, dst, 0, argStart, operand, n.Line)
	g.nextReg = g.baseReg
}
