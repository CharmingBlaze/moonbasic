package codegen

import (
	"moonbasic/compiler/ast"
	"moonbasic/compiler/entityspatial"
)

func (g *CodeGen) validateEntityMacroConstArg(e ast.Expr, line, col int) {
	id, ok := entityspatial.ConstEntitySlotID(e)
	if !ok {
		return
	}
	if err := entityspatial.ValidateLiteralSlot(id); err != nil {
		g.codegenError(line, col, err.Error(), entityspatial.LiteralSlotHint())
	}
}
