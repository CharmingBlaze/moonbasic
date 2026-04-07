package blitzengine

import (
	"moonbasic/runtime"
	"moonbasic/vm/value"
)

// call invokes a registry command by name with a variadic argument list.
func call(rt *runtime.Runtime, name string, args ...value.Value) (value.Value, error) {
	return rt.Call(name, args)
}
