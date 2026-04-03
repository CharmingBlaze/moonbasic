package runtime

import (
	"moonbasic/vm/value"
)

// ArgString extracts a string from the argument list.
func (rt *Runtime) ArgString(args []value.Value, index int) (string, error) {
	if index >= len(args) {
		return "", Errorf("missing argument %d", index+1)
	}
	if args[index].Kind != value.KindString {
		return "", Errorf("argument %d must be a string", index+1)
	}
	idx := int32(args[index].IVal)
	s, ok := rt.Heap.GetString(idx)
	if ok {
		return s, nil
	}
	// Fallback to program's string table for compile-time strings
	if rt.Prog != nil && idx >= 0 && int(idx) < len(rt.Prog.StringTable) {
		return rt.Prog.StringTable[idx], nil
	}
	return "", Errorf("invalid string reference for argument %d", index+1)
}

// ArgInt extracts an integer from the argument list.
func (rt *Runtime) ArgInt(args []value.Value, index int) (int64, error) {
	if index >= len(args) {
		return 0, Errorf("missing argument %d", index+1)
	}
	val, ok := args[index].ToInt()
	if !ok {
		return 0, Errorf("argument %d must be an integer", index+1)
	}
	return val, nil
}

// ArgFloat extracts a float from the argument list.
func (rt *Runtime) ArgFloat(args []value.Value, index int) (float64, error) {
	if index >= len(args) {
		return 0, Errorf("missing argument %d", index+1)
	}
	val, ok := args[index].ToFloat()
	if !ok {
		return 0, Errorf("argument %d must be a float", index+1)
	}
	return val, nil
}

// ArgHandle extracts a heap handle from the argument list.
func (rt *Runtime) ArgHandle(args []value.Value, index int) (int32, error) {
	if index >= len(args) {
		return 0, Errorf("missing argument %d", index+1)
	}
	if args[index].Kind != value.KindHandle {
		return 0, Errorf("argument %d must be a handle", index+1)
	}
	return int32(args[index].IVal), nil
}

// ArgBool extracts a boolean from the argument list.
func (rt *Runtime) ArgBool(args []value.Value, index int) (bool, error) {
	if index >= len(args) {
		return false, Errorf("missing argument %d", index+1)
	}
	var pool []string
	if rt.Prog != nil {
		pool = rt.Prog.StringTable
	}
	return value.Truthy(args[index], pool), nil
}

// RetString creates a string return value.
func (rt *Runtime) RetString(s string) value.Value {
	return value.FromStringIndex(rt.Heap.Intern(s))
}

// RetInt creates an integer return value.
func (rt *Runtime) RetInt(i int64) value.Value {
	return value.FromInt(i)
}

// RetFloat creates a float return value.
func (rt *Runtime) RetFloat(f float64) value.Value {
	return value.FromFloat(f)
}

// RetBool creates a boolean return value.
func (rt *Runtime) RetBool(b bool) value.Value {
	return value.FromBool(b)
}

// RetHandle creates a handle return value.
func (rt *Runtime) RetHandle(h int32) value.Value {
	return value.FromHandle(h)
}
