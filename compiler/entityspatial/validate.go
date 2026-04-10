// Package entityspatial holds shared rules for ENTITY.* spatial macro indices (SoA fast path).
package entityspatial

import (
	"fmt"
	"math"
	"strings"

	"moonbasic/compiler/ast"
	"moonbasic/runtime"
)

// ConstEntitySlotID returns (id, true) for integer or whole-float literals.
func ConstEntitySlotID(e ast.Expr) (int64, bool) {
	switch n := e.(type) {
	case *ast.IntLitNode:
		return n.Value, true
	case *ast.FloatLitNode:
		v := n.Value
		if math.IsNaN(v) || math.IsInf(v, 0) {
			return 0, false
		}
		iv := int64(v)
		if float64(iv) != v {
			return 0, false
		}
		return iv, true
	default:
		return 0, false
	}
}

// SpatialPropID returns the SoA column index for ENTITY.X/Y/Z/P/W/YAW/R macros.
func SpatialPropID(method string) (int, bool) {
	switch strings.ToUpper(strings.TrimSpace(method)) {
	case "X":
		return 0, true
	case "Y":
		return 1, true
	case "Z":
		return 2, true
	case "P":
		return 3, true
	case "W", "YAW":
		return 4, true
	case "R":
		return 5, true
	default:
		return -1, false
	}
}

// ValidateLiteralSlot returns an error if a compile-time-known index is out of range.
func ValidateLiteralSlot(id int64) error {
	if id < 0 {
		return fmt.Errorf("entity id %d is invalid (negative)", id)
	}
	if id >= runtime.MaxEntitySpatialIndex {
		return fmt.Errorf("entity id %d exceeds compile-time limit %d", id, runtime.MaxEntitySpatialIndex-1)
	}
	return nil
}

// LiteralSlotHint is appended as a semantic/codegen hint for out-of-range literals.
func LiteralSlotHint() string {
	return fmt.Sprintf("Use an entity id in [0, %d], or a variable (bounds-checked at runtime).", runtime.MaxEntitySpatialIndex-1)
}
