package mbtime

import (
	"time"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

// registerWallClock registers local wall-clock date/time flats (time.Now).
func registerWallClock(reg runtime.Registrar) {
	reg.Register("YEAR", "time", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		_ = rt
		if len(args) != 0 {
			return value.Nil, errArgs(0, len(args))
		}
		return value.FromInt(int64(time.Now().Year())), nil
	})
	reg.Register("MONTH", "time", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		_ = rt
		if len(args) != 0 {
			return value.Nil, errArgs(0, len(args))
		}
		return value.FromInt(int64(time.Now().Month())), nil
	})
	reg.Register("DAY", "time", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		_ = rt
		if len(args) != 0 {
			return value.Nil, errArgs(0, len(args))
		}
		return value.FromInt(int64(time.Now().Day())), nil
	})
	reg.Register("HOUR", "time", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		_ = rt
		if len(args) != 0 {
			return value.Nil, errArgs(0, len(args))
		}
		return value.FromInt(int64(time.Now().Hour())), nil
	})
	reg.Register("MINUTE", "time", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		_ = rt
		if len(args) != 0 {
			return value.Nil, errArgs(0, len(args))
		}
		return value.FromInt(int64(time.Now().Minute())), nil
	})
	reg.Register("SECOND", "time", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		_ = rt
		if len(args) != 0 {
			return value.Nil, errArgs(0, len(args))
		}
		return value.FromInt(int64(time.Now().Second())), nil
	})
	reg.Register("MILLISECOND", "time", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		_ = rt
		if len(args) != 0 {
			return value.Nil, errArgs(0, len(args))
		}
		ms := time.Now().Nanosecond() / 1_000_000
		return value.FromInt(int64(ms)), nil
	})
	reg.Register("TIMESTAMP", "time", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		_ = rt
		if len(args) != 0 {
			return value.Nil, errArgs(0, len(args))
		}
		return value.FromFloat(float64(time.Now().Unix())), nil
	})
	reg.Register("DATE", "time", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 0 {
			return value.Nil, errArgs(0, len(args))
		}
		return rt.RetString(time.Now().Format("2006-01-02")), nil
	})
	reg.Register("TIME", "time", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 0 {
			return value.Nil, errArgs(0, len(args))
		}
		return rt.RetString(time.Now().Format("15:04:05")), nil
	})
	// Blitz-style names (same wall clock as TIME / DATE).
	reg.Register("CURRENTTIME", "time", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 0 {
			return value.Nil, errArgs(0, len(args))
		}
		return rt.RetString(time.Now().Format("15:04:05")), nil
	})
	reg.Register("CURRENTDATE", "time", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 0 {
			return value.Nil, errArgs(0, len(args))
		}
		return rt.RetString(time.Now().Format("02 Jan 2006")), nil
	})
	reg.Register("DATETIME", "time", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 0 {
			return value.Nil, errArgs(0, len(args))
		}
		return rt.RetString(time.Now().Format("2006-01-02 15:04:05")), nil
	})
}
