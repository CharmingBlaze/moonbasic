//go:build cgo

package input

/*
extern const char *GetKeyName(int key);
*/
import "C"

import (
	"fmt"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

func (m *Module) inGetKeyName(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("INPUT.GETKEYNAME expects (key)")
	}
	k, err := keyCodeArg(args[0])
	if err != nil {
		return value.Nil, fmt.Errorf("INPUT.GETKEYNAME: %w", err)
	}
	p := C.GetKeyName(C.int(k))
	if p == nil {
		return rt.RetString(""), nil
	}
	return rt.RetString(C.GoString(p)), nil
}
