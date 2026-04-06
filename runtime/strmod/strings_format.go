package strmod

import (
	"encoding/binary"
	"fmt"
	"math"
	"strconv"
	"strings"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

func registerStringsFormat(r runtime.Registrar) {
	r.Register("BIN$", "core", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 1 {
			return value.Value{}, runtime.Errorf("BIN$ expects 1 argument")
		}
		n, err := rt.ArgInt(args, 0)
		if err != nil {
			return value.Value{}, err
		}
		return rt.RetString(strconv.FormatInt(n, 2)), nil
	})
	r.Register("HEX$", "core", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 1 {
			return value.Value{}, runtime.Errorf("HEX$ expects 1 argument")
		}
		n, err := rt.ArgInt(args, 0)
		if err != nil {
			return value.Value{}, err
		}
		return rt.RetString(strings.ToUpper(strconv.FormatInt(n, 16))), nil
	})
	r.Register("OCT$", "core", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 1 {
			return value.Value{}, runtime.Errorf("OCT$ expects 1 argument")
		}
		n, err := rt.ArgInt(args, 0)
		if err != nil {
			return value.Value{}, err
		}
		return rt.RetString(strconv.FormatInt(n, 8)), nil
	})
	r.Register("FORMAT$", "core", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 2 || args[1].Kind != value.KindString {
			return value.Value{}, runtime.Errorf("FORMAT$ expects (v, pattern$)")
		}
		pat, err := rt.ArgString(args, 1)
		if err != nil {
			return value.Value{}, err
		}
		v := args[0]
		switch v.Kind {
		case value.KindInt:
			return rt.RetString(fmt.Sprintf(pat, v.IVal)), nil
		case value.KindFloat:
			return rt.RetString(fmt.Sprintf(pat, v.FVal)), nil
		case value.KindString:
			s, err := rt.ArgString(args, 0)
			if err != nil {
				return value.Value{}, err
			}
			return rt.RetString(fmt.Sprintf(pat, s)), nil
		case value.KindBool:
			b := v.IVal != 0
			return rt.RetString(fmt.Sprintf(pat, b)), nil
		case value.KindHandle:
			return rt.RetString(fmt.Sprintf(pat, v.IVal)), nil
		default:
			return rt.RetString(fmt.Sprintf(pat, v.String())), nil
		}
	})

	r.Register("MKSHORT$", "core", mkFixedLE(2, func(b []byte, rt *runtime.Runtime, args []value.Value) error {
		n, err := rt.ArgInt(args, 0)
		if err != nil {
			return err
		}
		binary.LittleEndian.PutUint16(b, uint16(int16(n)))
		return nil
	}))
	r.Register("MKINT$", "core", mkFixedLE(4, func(b []byte, rt *runtime.Runtime, args []value.Value) error {
		n, err := rt.ArgInt(args, 0)
		if err != nil {
			return err
		}
		binary.LittleEndian.PutUint32(b, uint32(int32(n)))
		return nil
	}))
	r.Register("MKLONG$", "core", mkFixedLE(8, func(b []byte, rt *runtime.Runtime, args []value.Value) error {
		n, err := rt.ArgInt(args, 0)
		if err != nil {
			return err
		}
		binary.LittleEndian.PutUint64(b, uint64(n))
		return nil
	}))
	r.Register("MKFLOAT$", "core", mkFixedLE(4, func(b []byte, rt *runtime.Runtime, args []value.Value) error {
		f, ok := args[0].ToFloat()
		if !ok {
			if args[0].Kind == value.KindInt {
				f = float64(args[0].IVal)
				ok = true
			}
		}
		if !ok {
			return runtime.Errorf("MKFLOAT$: numeric expected")
		}
		binary.LittleEndian.PutUint32(b, math.Float32bits(float32(f)))
		return nil
	}))
	r.Register("MKDOUBLE$", "core", mkFixedLE(8, func(b []byte, rt *runtime.Runtime, args []value.Value) error {
		f, ok := args[0].ToFloat()
		if !ok {
			if args[0].Kind == value.KindInt {
				f = float64(args[0].IVal)
				ok = true
			}
		}
		if !ok {
			return runtime.Errorf("MKDOUBLE$: numeric expected")
		}
		binary.LittleEndian.PutUint64(b, math.Float64bits(f))
		return nil
	}))

	r.Register("CVSHORT", "core", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 1 {
			return value.Value{}, runtime.Errorf("CVSHORT expects 1 string argument")
		}
		s, err := rt.ArgString(args, 0)
		if err != nil {
			return value.Value{}, err
		}
		b := []byte(s)
		if len(b) < 2 {
			return value.FromInt(0), nil
		}
		return value.FromInt(int64(int16(binary.LittleEndian.Uint16(b[:2])))), nil
	})
	r.Register("CVINT", "core", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 1 {
			return value.Value{}, runtime.Errorf("CVINT expects 1 string argument")
		}
		s, err := rt.ArgString(args, 0)
		if err != nil {
			return value.Value{}, err
		}
		b := []byte(s)
		if len(b) < 4 {
			return value.FromInt(0), nil
		}
		return value.FromInt(int64(int32(binary.LittleEndian.Uint32(b[:4])))), nil
	})
	r.Register("CVLONG", "core", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 1 {
			return value.Value{}, runtime.Errorf("CVLONG expects 1 string argument")
		}
		s, err := rt.ArgString(args, 0)
		if err != nil {
			return value.Value{}, err
		}
		b := []byte(s)
		if len(b) < 8 {
			return value.FromInt(0), nil
		}
		return value.FromInt(int64(binary.LittleEndian.Uint64(b[:8]))), nil
	})
	r.Register("CVFLOAT", "core", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 1 {
			return value.Value{}, runtime.Errorf("CVFLOAT expects 1 string argument")
		}
		s, err := rt.ArgString(args, 0)
		if err != nil {
			return value.Value{}, err
		}
		b := []byte(s)
		if len(b) < 4 {
			return value.FromFloat(0), nil
		}
		u := binary.LittleEndian.Uint32(b[:4])
		return value.FromFloat(float64(math.Float32frombits(u))), nil
	})
	r.Register("CVDOUBLE", "core", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 1 {
			return value.Value{}, runtime.Errorf("CVDOUBLE expects 1 string argument")
		}
		s, err := rt.ArgString(args, 0)
		if err != nil {
			return value.Value{}, err
		}
		b := []byte(s)
		if len(b) < 8 {
			return value.FromFloat(0), nil
		}
		u := binary.LittleEndian.Uint64(b[:8])
		return value.FromFloat(math.Float64frombits(u)), nil
	})
}

func mkFixedLE(n int, fill func([]byte, *runtime.Runtime, []value.Value) error) runtime.BuiltinFn {
	return func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 1 {
			return value.Value{}, runtime.Errorf("expects 1 argument")
		}
		b := make([]byte, n)
		if err := fill(b, rt, args); err != nil {
			return value.Value{}, err
		}
		return rt.RetString(string(b)), nil
	}
}
