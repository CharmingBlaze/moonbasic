// Package types provides type information and inference for moonBASIC.
// This supports the modern syntax with implicit variable declaration and type inference.
package types

// Tag represents the type of a value or variable.
type Tag byte

const (
	Unknown Tag = iota
	Int
	Float
	String
	Bool
	Array
	Handle
	UserType // For TYPE ... ENDTYPE definitions
)

// String returns the human-readable name of the type.
func (t Tag) String() string {
	names := []string{
		"unknown", "int", "float", "string", "bool", "array", "handle", "usertype",
	}
	if int(t) < len(names) {
		return names[t]
	}
	return "invalid"
}

// FromSuffix returns the type tag based on variable suffix.
// moonBASIC uses: # = float, $ = string, ? = bool, no suffix = int.
func FromSuffix(name string) Tag {
	if len(name) == 0 {
		return Int // Default to int for untyped
	}
	switch name[len(name)-1] {
	case '#':
		return Float
	case '$':
		return String
	case '?':
		return Bool
	default:
		return Int
	}
}

// Info holds detailed type information for a symbol.
type Info struct {
	Tag      Tag
	UserType string // For UserType, stores the type name
	ArrayOf  *Info  // For arrays, stores element type
}

// IsNumeric returns true if the type is Int or Float.
func (t Tag) IsNumeric() bool {
	return t == Int || t == Float
}

// Promote returns the result type of combining two types.
// Float wins over Int. String only combines with String.
func Promote(a, b Tag) Tag {
	if a == Float || b == Float {
		return Float
	}
	if a == Int && b == Int {
		return Int
	}
	if a == String || b == String {
		return String
	}
	if a == Bool && b == Bool {
		return Bool
	}
	return Unknown
}

// IsCompatible checks if a value of type 'from' can be assigned to a variable of type 'to'.
func IsCompatible(from, to Tag) bool {
	if to == Unknown || from == Unknown {
		return true // Allow unknown during inference
	}
	if to == from {
		return true
	}
	// Numeric promotion: int can go to float
	if to == Float && from == Int {
		return true
	}
	return false
}
