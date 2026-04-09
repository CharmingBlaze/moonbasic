package window

import "moonbasic/vm/value"

func argInt(v value.Value) (int64, bool) {
	if i, ok := v.ToInt(); ok {
		return i, true
	}
	if f, ok := v.ToFloat(); ok {
		return int64(f), true
	}
	return 0, false
}

func clampU8(v int64) uint8 {
	switch {
	case v < 0:
		return 0
	case v > 255:
		return 255
	default:
		return uint8(v)
	}
}
