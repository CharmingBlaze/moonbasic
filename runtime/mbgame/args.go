package mbgame

import "moonbasic/vm/value"

func argF(v value.Value) (float64, bool) {
	if f, ok := v.ToFloat(); ok {
		return f, true
	}
	if i, ok := v.ToInt(); ok {
		return float64(i), true
	}
	return 0, false
}

func argI(v value.Value) (int64, bool) {
	if i, ok := v.ToInt(); ok {
		return i, true
	}
	if f, ok := v.ToFloat(); ok {
		return int64(f), true
	}
	return 0, false
}
