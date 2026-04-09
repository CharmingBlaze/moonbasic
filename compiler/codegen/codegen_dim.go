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
		// OpArrayRedim: Dst=PreserveFlag, SrcA=HandleReg, SrcB=ArgStart, Operand=DimCount
		ch.Emit(opcode.OpArrayRedim, preserve, hReg, argStart, int32(len(n.Dims)), n.Line)
		g.nextReg = g.baseReg
		return
	}
	
	if n.ElemType != "" {
		argStart := g.nextReg
		for _, d := range n.Dims {
			g.emitExpr(ch, d)
		}
		tn := strings.ToUpper(n.ElemType)
		tidx := ch.AddName(tn)
		
		dst := g.allocReg()
		// OpArrayMakeTyped: Dst=handle, SrcA=ArgStart, SrcB=DimCount, Operand=TypeIdx
		ch.Emit(opcode.OpArrayMakeTyped, dst, argStart, uint8(len(n.Dims)), tidx, n.Line)
		
		g.emitStoreNamed(ch, n.Name, n.Line, dst)
		g.nextReg = g.baseReg
		return
	}
	
	argStart := g.nextReg
	for _, d := range n.Dims {
		g.emitExpr(ch, d)
	}
	dst := g.allocReg()
	// OpArrayMake: Dst=handle, SrcA=Kind, SrcB=ArgStart, Operand=DimCount
	ch.Emit(opcode.OpArrayMake, dst, flags, argStart, int32(len(n.Dims)), n.Line)
	
	g.emitStoreNamed(ch, n.Name, n.Line, dst)
	g.nextReg = g.baseReg
}
