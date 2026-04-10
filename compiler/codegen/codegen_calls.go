package codegen

import (
	"moonbasic/compiler/ast"
	"moonbasic/vm/opcode"
)

// emitArgsStable evaluates args left-to-right and packs their final values into
// contiguous registers starting at argStart. This prevents nested expressions
// from clobbering earlier argument registers.
func (g *CodeGen) emitArgsStable(ch *opcode.Chunk, args []ast.Expr, line int) uint8 {
	argStart := g.nextReg
	nextArg := argStart
	for _, a := range args {
		r := g.emitExpr(ch, a)
		if r != nextArg {
			ch.Emit(opcode.OpMove, nextArg, r, 0, 0, line)
		}
		nextArg++
		g.nextReg = nextArg
	}
	return argStart
}

// emitCallStmt translates a call statement (builtin or user function) into bytecode.
func (g *CodeGen) emitCallStmt(ch *opcode.Chunk, n *ast.CallStmtNode) {
	if g.err != nil {
		return
	}
	g.nextReg = g.baseReg
	argStart := g.emitArgsStable(ch, n.Args, n.Line)

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
