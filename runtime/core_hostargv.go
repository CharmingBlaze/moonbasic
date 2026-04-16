package runtime

import (
	"os"
	"strings"

	"moonbasic/vm/value"
)

func registerHostArgv(r Registrar) {
	r.Register("ARGC", "core", func(rt *Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 0 {
			return value.Nil, Errorf("ARGC expects 0 arguments, got %d", len(args))
		}
		argv := argvFor(rt)
		return value.FromInt(int64(len(argv))), nil
	})
	commandFn := func(rt *Runtime, args ...value.Value) (value.Value, error) {
		argv := argvFor(rt)
		switch len(args) {
		case 0:
			return rt.RetString(strings.Join(argv, " ")), nil
		case 1:
			idx, ok := args[0].ToInt()
			if !ok {
				return value.Nil, Errorf("COMMAND expects int index")
			}
			if idx < 0 || int(idx) >= len(argv) {
				return rt.RetString(""), nil
			}
			return rt.RetString(argv[idx]), nil
		default:
			return value.Nil, Errorf("COMMAND expects 0 or 1 arguments, got %d", len(args))
		}
	}
	r.Register("COMMAND", "core", commandFn)
	r.Register("COMMAND$", "core", commandFn)
}

func argvFor(rt *Runtime) []string {
	if rt == nil || rt.HostArgs == nil {
		return os.Args
	}
	return rt.HostArgs
}
