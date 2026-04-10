package codegen

import (
	"strings"

	"moonbasic/compiler/ast"
	"moonbasic/vm/opcode"
)

func (g *CodeGen) emitDim(ch *opcode.Chunk, n *ast.DimNode) {
	g.nextReg = g.baseReg
	flags := arrayKindFlags(n.Name)
	if n.IsRedim {
		hReg := g.emitExpr(ch, &ast.IdentNode{Name: n.Name, Line: n.Line, Col: n.Col})
		argStart := g.nextReg
		for _, d := range n.Dims {
			g.emitExpr(ch, d)
		}
		preserve := uint8(0)
		if n.Preserve {
			preserve = 1
		}
		ch.Emit(opcode.OpArrayRedim, preserve, hReg, argStart, int32(len(n.Dims)), n.Line)
		ch.SetLastArrayDebugName(ch.AddName(strings.ToUpper(n.Name)))
		g.nextReg = g.baseReg
		return
	}

	if n.TypeName != "" {
		tn := strings.ToUpper(strings.TrimSpace(n.TypeName))
		argStart := g.nextReg
		for _, d := range n.Dims {
			g.emitExpr(ch, d)
		}
		dst := g.allocReg()

		switch tn {
		case "HANDLE":
			ch.Emit(opcode.OpArrayMake, dst, argStart, 3, int32(len(n.Dims)), n.Line)
			ch.SetLastArrayDebugName(ch.AddName(strings.ToUpper(n.Name)))
		case "STRING":
			ch.Emit(opcode.OpArrayMake, dst, argStart, 1, int32(len(n.Dims)), n.Line)
			ch.SetLastArrayDebugName(ch.AddName(strings.ToUpper(n.Name)))
		case "INTEGER", "FLOAT":
			ch.Emit(opcode.OpArrayMake, dst, argStart, 0, int32(len(n.Dims)), n.Line)
			ch.SetLastArrayDebugName(ch.AddName(strings.ToUpper(n.Name)))
		default:
			// User TYPE array — heap instance per cell
			tidx := ch.AddName(tn)
			ch.Emit(opcode.OpArrayMakeTyped, dst, argStart, uint8(len(n.Dims)), tidx, n.Line)
			ch.SetLastArrayDebugName(ch.AddName(strings.ToUpper(n.Name)))
		}

		g.emitStoreNamed(ch, n.Name, n.Line, dst)
		g.nextReg = g.baseReg
		return
	}

	argStart := g.nextReg
	for _, d := range n.Dims {
		g.emitExpr(ch, d)
	}
	dst := g.allocReg()
	// OpArrayMake: SrcA = dim value registers start, SrcB = kind flags (see vm doArrayMake)
	ch.Emit(opcode.OpArrayMake, dst, argStart, flags, int32(len(n.Dims)), n.Line)
	ch.SetLastArrayDebugName(ch.AddName(strings.ToUpper(n.Name)))

	g.emitStoreNamed(ch, n.Name, n.Line, dst)
	g.nextReg = g.baseReg
}
