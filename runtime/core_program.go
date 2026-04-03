package runtime

import (
	"time"

	"moonbasic/vm/value"
)

func registerProgramControl(r Registrar) {
	quitStop := func(rt *Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 0 {
			return value.Nil, Errorf("expects 0 arguments, got %d", len(args))
		}
		if rt.TerminateVM != nil {
			rt.TerminateVM()
		}
		return value.Nil, nil
	}
	r.Register("QUIT", "core", quitStop)
	r.Register("STOP", "core", quitStop)
	r.Register("SLEEP", "core", sleepOrWaitBuiltin)
	r.Register("WAIT", "core", sleepOrWaitBuiltin)
}

func sleepOrWaitBuiltin(rt *Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) != 1 {
		return value.Nil, Errorf("expects 1 argument, got %d", len(args))
	}
	d, err := sleepDuration(args[0])
	if err != nil {
		return value.Nil, err
	}
	time.Sleep(d)
	return value.Nil, nil
}

// sleepDuration maps SLEEP/WAIT argument: int → milliseconds, float → seconds.
func sleepDuration(v value.Value) (time.Duration, error) {
	switch v.Kind {
	case value.KindInt:
		if v.IVal < 0 {
			return 0, nil
		}
		return time.Duration(v.IVal) * time.Millisecond, nil
	case value.KindFloat:
		if v.FVal < 0 {
			return 0, nil
		}
		return time.Duration(v.FVal * float64(time.Second)), nil
	default:
		return 0, Errorf("SLEEP expects numeric argument")
	}
}
