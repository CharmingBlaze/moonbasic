package codegen

import (
	"moonbasic/compiler/ast"
	"moonbasic/compiler/symtable"
	"moonbasic/vm/opcode"
)

func (g *CodeGen) emitDim(ch *opcode.Chunk, n *ast.DimNode) {
	flags := arrayKindFlags(n.Name)
	if n.IsRedim {
		g.emitExpr(ch, &ast.IdentNode{Name: n.Name, Line: n.Line, Col: n.Col})
		for _, d := range n.Dims {
			g.emitExpr(ch, d)
		}
		preserve := uint8(0)
		if n.Preserve {
			preserve = 1
		}
		ch.Emit(opcode.OpArrayRedim, int32(len(n.Dims)), preserve, n.Line)
		return
	}
	for _, d := range n.Dims {
		g.emitExpr(ch, d)
	}
	ch.Emit(opcode.OpArrayMake, int32(len(n.Dims)), flags, n.Line)
	sym := g.resolveOrDefineAssignTarget(n.Name)
	if sym != nil && (sym.Kind == symtable.Local || sym.Kind == symtable.Param) {
		ch.Emit(opcode.OpStoreLocal, int32(sym.Slot), 0, n.Line)
	} else if sym != nil && sym.Kind == symtable.Static {
		k := ch.AddName(sym.StaticKey)
		ch.Emit(opcode.OpStoreGlobal, k, 0, n.Line)
	} else {
		idx := ch.AddName(n.Name)
		ch.Emit(opcode.OpStoreGlobal, idx, 0, n.Line)
	}
	ch.Emit(opcode.OpPop, 0, 0, n.Line)
}
