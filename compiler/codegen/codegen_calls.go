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
	// 1. Push arguments onto the stack
	for _, a := range n.Args {
		g.emitExpr(ch, a)
	}

	// 2. Resolve name and determine OpCode
	idx := ch.AddName(n.Name)
	op := opcode.OpCallBuiltin
	if _, ok := g.Prog.Functions[n.Name]; ok {
		op = opcode.OpCallUser
	}

	// 3. Emit Call
	ac := g.argCountFlags(len(n.Args), n.Line, n.Col)
	if g.err != nil {
		return
	}
	ch.Emit(op, idx, ac, n.Line)

	// 4. Pop the return value (even if it's NULL)
	// MoonBASIC statements never leave anything on the stack.
	ch.Emit(opcode.OpPop, 0, 0, n.Line)
}
