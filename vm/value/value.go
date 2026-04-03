// Package value defines the tagged union for moonBASIC VM runtime values.
// Values are passed by copy (fixed 24 bytes, IR v2). KindString stores an index
// into the program string table, not a Go string.
package value

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"unsafe"
)

func init() {
	if unsafe.Sizeof(Value{}) != 24 {
		panic("value: Value must be exactly 24 bytes (IR v2)")
	}
}

// Kind identifies the runtime type of a Value.
type Kind byte

const (
	KindNil Kind = iota
	KindInt
	KindFloat
	KindString
	KindBool
	KindHandle
)

// Value is the core tagged union (24 bytes on amd64).
// KindString: IVal holds the int32 index (sign-extended in IVal).
type Value struct {
	Kind Kind
	_    [7]byte
	IVal int64
	FVal float64
}

// Global Nil value.
var Nil = Value{Kind: KindNil}

// FromStringIndex builds a string value from a program string-pool index.
func FromStringIndex(idx int32) Value {
	return Value{Kind: KindString, IVal: int64(idx)}
}

func FromInt(v int64) Value    { return Value{Kind: KindInt, IVal: v} }
func FromFloat(v float64) Value { return Value{Kind: KindFloat, FVal: v} }
func FromBool(v bool) Value {
	if v {
		return Value{Kind: KindBool, IVal: 1}
	}
	return Value{Kind: KindBool, IVal: 0}
}
func FromHandle(h int32) Value { return Value{Kind: KindHandle, IVal: int64(h)} }

func Int(v int64) Value       { return FromInt(v) }
func Float(v float64) Value   { return FromFloat(v) }
func Bool(v bool) Value       { return FromBool(v) }
func Handle(h int32) Value    { return FromHandle(h) }

// StringIndex returns the pool index for KindString, or 0 otherwise.
func (v Value) StringIndex() int32 {
	if v.Kind != KindString {
		return 0
	}
	return int32(v.IVal)
}

// StringAt resolves KindString using the program pool; other kinds use String().
func StringAt(v Value, pool []string) string {
	if v.Kind != KindString {
		return v.String()
	}
	i := int32(v.IVal)
	if i < 0 || int(i) >= len(pool) {
		return ""
	}
	return pool[i]
}

// Truthy returns the BASIC boolean truth value. strPool is required for KindString.
func Truthy(v Value, strPool []string) bool {
	switch v.Kind {
	case KindNil:
		return false
	case KindInt, KindHandle:
		return v.IVal != 0
	case KindFloat:
		return v.FVal != 0
	case KindString:
		return StringAt(v, strPool) != ""
	case KindBool:
		return v.IVal != 0
	default:
		return false
	}
}

// String returns a printable representation for tracing and debugging.
// KindString renders as {idx} when no pool is available.
func (v Value) String() string {
	switch v.Kind {
	case KindNil:
		return "NULL"
	case KindInt:
		return strconv.FormatInt(v.IVal, 10)
	case KindFloat:
		s := strconv.FormatFloat(v.FVal, 'g', -1, 64)
		if !strings.ContainsAny(s, ".eE") {
			s += ".0"
		}
		return s
	case KindString:
		return fmt.Sprintf("{%d}", int32(v.IVal))
	case KindBool:
		if v.IVal != 0 {
			return "TRUE"
		}
		return "FALSE"
	case KindHandle:
		return fmt.Sprintf("<handle %d>", v.IVal)
	default:
		return "?"
	}
}

// TypeName returns a mnemonic name for error messages.
func (v Value) TypeName() string {
	switch v.Kind {
	case KindNil:
		return "NULL"
	case KindInt:
		return "INT"
	case KindFloat:
		return "FLOAT"
	case KindString:
		return "STRING"
	case KindBool:
		return "BOOL"
	case KindHandle:
		return "HANDLE"
	default:
		return "UNKNOWN"
	}
}

// ToFloat coerces the value to float64. returns (float, success).
func (v Value) ToFloat() (float64, bool) {
	switch v.Kind {
	case KindFloat:
		return v.FVal, true
	case KindInt:
		return float64(v.IVal), true
	default:
		return 0.0, false
	}
}

// ToInt coerces the value to int64. returns (int, success).
func (v Value) ToInt() (int64, bool) {
	switch v.Kind {
	case KindInt:
		return v.IVal, true
	case KindFloat:
		return int64(v.FVal), true
	case KindHandle:
		return v.IVal, true
	default:
		return 0, false
	}
}

// Add performs numeric addition (no string concatenation; see VM for STR + STR).
func Add(a, b Value) (Value, error) {
	if a.Kind == KindFloat || b.Kind == KindFloat {
		af, _ := a.ToFloat()
		bf, _ := b.ToFloat()
		return FromFloat(af + bf), nil
	}
	ai, aok := a.ToInt()
	bi, bok := b.ToInt()
	if aok && bok {
		return FromInt(ai + bi), nil
	}
	return Nil, fmt.Errorf("cannot add %s and %s", a.TypeName(), b.TypeName())
}

func Sub(a, b Value) (Value, error) {
	if a.Kind == KindFloat || b.Kind == KindFloat {
		af, _ := a.ToFloat()
		bf, _ := b.ToFloat()
		return FromFloat(af - bf), nil
	}
	ai, aok := a.ToInt()
	bi, bok := b.ToInt()
	if aok && bok {
		return FromInt(ai - bi), nil
	}
	return Nil, fmt.Errorf("cannot subtract %s and %s", a.TypeName(), b.TypeName())
}

func Mul(a, b Value) (Value, error) {
	if a.Kind == KindFloat || b.Kind == KindFloat {
		af, _ := a.ToFloat()
		bf, _ := b.ToFloat()
		return FromFloat(af * bf), nil
	}
	ai, aok := a.ToInt()
	bi, bok := b.ToInt()
	if aok && bok {
		return FromInt(ai * bi), nil
	}
	return Nil, fmt.Errorf("cannot multiply %s and %s", a.TypeName(), b.TypeName())
}

func Div(a, b Value) (Value, error) {
	if a.Kind == KindFloat || b.Kind == KindFloat {
		af, _ := a.ToFloat()
		bf, _ := b.ToFloat()
		if bf == 0 {
			return Nil, fmt.Errorf("division by zero")
		}
		return FromFloat(af / bf), nil
	}
	ai, aok := a.ToInt()
	bi, bok := b.ToInt()
	if aok && bok {
		if bi == 0 {
			return Nil, fmt.Errorf("division by zero")
		}
		return FromInt(ai / bi), nil
	}
	return Nil, fmt.Errorf("cannot divide %s and %s", a.TypeName(), b.TypeName())
}

func Mod(a, b Value) (Value, error) {
	ai, aok := a.ToInt()
	bi, bok := b.ToInt()
	if aok && bok {
		if bi == 0 {
			return Nil, fmt.Errorf("modulo by zero")
		}
		return FromInt(ai % bi), nil
	}
	return Nil, fmt.Errorf("cannot MOD %s and %s", a.TypeName(), b.TypeName())
}

func Pow(a, b Value) (Value, error) {
	af, aok := a.ToFloat()
	bf, bok := b.ToFloat()
	if aok && bok {
		return FromFloat(math.Pow(af, bf)), nil
	}
	return Nil, fmt.Errorf("cannot raise %s to power of %s", a.TypeName(), b.TypeName())
}

func Neg(v Value) (Value, error) {
	switch v.Kind {
	case KindInt:
		return FromInt(-v.IVal), nil
	case KindFloat:
		return FromFloat(-v.FVal), nil
	default:
		return Nil, fmt.Errorf("cannot negate %s", v.TypeName())
	}
}

// Equal compares two values. KindString uses pool index equality (same intern table).
func Equal(a, b Value) bool {
	if a.Kind == KindFloat || b.Kind == KindFloat {
		af, aok := a.ToFloat()
		bf, bok := b.ToFloat()
		return aok && bok && (af == bf)
	}
	if a.Kind == KindString && b.Kind == KindString {
		return a.IVal == b.IVal
	}
	if a.Kind == KindBool && b.Kind == KindBool {
		return a.IVal == b.IVal
	}
	if (a.Kind == KindInt || a.Kind == KindHandle) && (b.Kind == KindInt || b.Kind == KindHandle) {
		return a.IVal == b.IVal
	}
	return false
}

// Less compares values; strPool is used for KindString lexical order.
func Less(a, b Value, strPool []string) (bool, error) {
	if a.Kind == KindFloat || b.Kind == KindFloat {
		af, aok := a.ToFloat()
		bf, bok := b.ToFloat()
		if aok && bok {
			return af < bf, nil
		}
	}
	ai, aok := a.ToInt()
	bi, bok := b.ToInt()
	if aok && bok {
		return ai < bi, nil
	}
	if a.Kind == KindString && b.Kind == KindString {
		return StringAt(a, strPool) < StringAt(b, strPool), nil
	}
	return false, fmt.Errorf("cannot compare %s and %s", a.TypeName(), b.TypeName())
}
