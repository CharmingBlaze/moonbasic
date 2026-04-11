package mathmod

import (
	"fmt"
	"math"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

func (m *Module) registerLogic(r runtime.Registrar) {
	r.Register("IIF", "math", m.builtinIIF)
	r.Register("CHOOSE", "math", m.builtinChoose)
	r.Register("SWITCH", "math", m.builtinSwitch)
}

// valuesEqual compares two values for equality in logic functions.
func valuesEqual(rt *runtime.Runtime, a, b value.Value) bool {
	if a.Kind != b.Kind {
		af, aok := a.ToFloat()
		bf, bok := b.ToFloat()
		if aok && bok {
			return af == bf
		}
		ai, aiok := a.ToInt()
		bi, biok := b.ToInt()
		if aiok && biok {
			return ai == bi
		}
		return false
	}
	switch a.Kind {
	case value.KindNil:
		return true
	case value.KindInt, value.KindBool, value.KindHandle:
		return a.IVal == b.IVal
	case value.KindFloat:
		return a.FVal == b.FVal
	case value.KindString:
		var pool []string
		if rt.Prog != nil {
			pool = rt.Prog.StringTable
		}
		return value.EqualStringValue(a, b, pool, rt.Heap)
	default:
		return false
	}
}

func (m *Module) builtinIIF(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("IIF expects 3 arguments (cond, a, b)")
	}
	if value.Truthy(args[0], rt.Prog.StringTable, rt.Heap) {
		return args[1], nil
	}
	return args[2], nil
}

func (m *Module) builtinChoose(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) < 2 {
		return value.Nil, fmt.Errorf("CHOOSE expects at least 2 arguments (index, v1, ...)")
	}
	idxF, _ := args[0].ToFloat()
	idx := int(math.Floor(idxF))
	if idx < 1 || idx >= len(args) {
		return value.Nil, nil
	}
	return args[idx], nil
}

func (m *Module) builtinSwitch(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) < 4 || (len(args)-2)%2 != 0 {
		return value.Nil, fmt.Errorf("SWITCH expects (expr, case1, val1, ..., caseN, valN, default)")
	}
	expr := args[0]
	defaultVal := args[len(args)-1]
	for i := 1; i < len(args)-1; i += 2 {
		if valuesEqual(rt, expr, args[i]) {
			return args[i+1], nil
		}
	}
	return defaultVal, nil
}
