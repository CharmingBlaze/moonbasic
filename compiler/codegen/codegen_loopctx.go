package codegen

import (
	"strings"

	"moonbasic/compiler/ast"
	"moonbasic/vm/opcode"
)

type loopFrame struct {
	kind            string
	continueTarget  int32 // -1 until set (FOR); absolute instruction index
	breakPatches    []int
	continuePatches []int
}

func (g *CodeGen) beginLoop(kind string, continueTarget int32) {
	g.loopStack = append(g.loopStack, loopFrame{
		kind:           kind,
		continueTarget: continueTarget,
	})
}

func (g *CodeGen) setLoopContinueTarget(target int32) {
	n := len(g.loopStack)
	if n == 0 {
		return
	}
	g.loopStack[n-1].continueTarget = target
}

// endLoop patches EXIT (break) and CONTINUE jumps for the innermost active loop, then pops it.
func (g *CodeGen) endLoop(ch *opcode.Chunk, breakTarget int) {
	n := len(g.loopStack)
	if n == 0 {
		g.codegenError(1, 1, "internal: endLoop with empty loop stack", "")
		return
	}
	f := &g.loopStack[n-1]
	bt := int32(breakTarget)
	for _, pc := range f.breakPatches {
		ch.Instructions[pc].Operand = bt
	}
	ct := f.continueTarget
	if ct < 0 {
		g.codegenError(1, 1, "internal: CONTINUE target not set for loop", "")
		return
	}
	for _, pc := range f.continuePatches {
		ch.Instructions[pc].Operand = ct
	}
	g.loopStack = g.loopStack[:n-1]
}

func loopKindFromAST(target string) string {
	return strings.ToLower(target)
}

func (g *CodeGen) emitExitStmt(ch *opcode.Chunk, n *ast.ExitStmt) {
	if n.Target == "FUNCTION" {
		if g.fnDepth == 0 {
			g.codegenError(n.Line, n.Col, "EXIT FUNCTION only allowed inside a FUNCTION", "")
			return
		}
		ch.Emit(opcode.OpReturnVoid, 0, 0, 0, 0, n.Line)
		return
	}
	k := loopKindFromAST(n.Target)
	for i := len(g.loopStack) - 1; i >= 0; i-- {
		if g.loopStack[i].kind == k {
			// OpJump: dst=0, srcA=0, srcB=0, operand=target
			j := ch.Emit(opcode.OpJump, 0, 0, 0, 0, n.Line)
			g.loopStack[i].breakPatches = append(g.loopStack[i].breakPatches, j)
			return
		}
	}
	g.codegenError(n.Line, n.Col, "EXIT "+n.Target+": no matching active loop", "")
}

func (g *CodeGen) emitContinueStmt(ch *opcode.Chunk, n *ast.ContinueStmt) {
	k := loopKindFromAST(n.Target)
	for i := len(g.loopStack) - 1; i >= 0; i-- {
		if g.loopStack[i].kind == k {
			j := ch.Emit(opcode.OpJump, 0, 0, 0, 0, n.Line)
			g.loopStack[i].continuePatches = append(g.loopStack[i].continuePatches, j)
			return
		}
	}
	g.codegenError(n.Line, n.Col, "CONTINUE "+n.Target+": no matching active loop", "")
}
