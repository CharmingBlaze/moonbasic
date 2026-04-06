package strmod

import (
	"strings"

	"moonbasic/runtime"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

func registerStringsSplitJoin(r runtime.Registrar) {
	r.Register("SPLIT$", "core", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if rt.Heap == nil {
			return value.Value{}, runtime.Errorf("SPLIT$: runtime heap not available")
		}
		if len(args) != 2 {
			return value.Value{}, runtime.Errorf("SPLIT$ expects 2 arguments (s$, delim$)")
		}
		s, err := rt.ArgString(args, 0)
		if err != nil {
			return value.Value{}, err
		}
		d, err := rt.ArgString(args, 1)
		if err != nil {
			return value.Value{}, err
		}
		var parts []string
		if d == "" {
			for _, ch := range s {
				parts = append(parts, string(ch))
			}
		} else {
			parts = strings.Split(s, d)
		}
		id, err := rt.Heap.Alloc(&heap.StringList{Items: parts})
		if err != nil {
			return value.Value{}, err
		}
		return value.FromHandle(id), nil
	})
	r.Register("JOIN$", "core", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if rt.Heap == nil {
			return value.Value{}, runtime.Errorf("JOIN$: runtime heap not available")
		}
		if len(args) != 2 {
			return value.Value{}, runtime.Errorf("JOIN$ expects 2 arguments (list, delim$)")
		}
		if args[0].Kind != value.KindHandle {
			return value.Value{}, runtime.Errorf("JOIN$: first argument must be a handle (string list or string array)")
		}
		delim, err := rt.ArgString(args, 1)
		if err != nil {
			return value.Value{}, err
		}
		h := heap.Handle(args[0].IVal)
		if sl, err := heap.Cast[*heap.StringList](rt.Heap, h); err == nil {
			return rt.RetString(strings.Join(sl.Items, delim)), nil
		}
		if arr, err := heap.Cast[*heap.Array](rt.Heap, h); err == nil {
			if arr.Kind != heap.ArrayKindString {
				return value.Value{}, runtime.Errorf("JOIN$: string array expected")
			}
			return rt.RetString(arr.JoinStrings(rt.Prog.StringTable, delim)), nil
		}
		return value.Value{}, runtime.Errorf("JOIN$: expected StringList or string Array handle")
	})
}
