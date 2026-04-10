package heap

import (
	"fmt"
	"sort"
	"strconv"
)

// ArrayKind selects element storage for a heap array.
type ArrayKind byte

const (
	ArrayKindFloat ArrayKind = iota
	ArrayKindString
	ArrayKindBool
	ArrayKindHandle  // elements are heap handles (e.g. TYPE instances)
	ArrayKindFloat32 // shared float32 buffer (e.g. physics matrices)
)

// MaxArrayCells caps total elements (flat count) for one allocation.
const MaxArrayCells int64 = 64_000_000

// Array is a dense row-major array stored on the heap as a handle.
// Language indices are 1-based per dimension; storage is 0-based flat.
// Float and bool elements use Floats; bool is stored as 0/1.
// String elements use Strings (program string-pool indices).
// Handle elements use Handles (raw heap IDs; 0 = null).
type Array struct {
	Dims    []int64
	Floats   []float64
	Floats32 []float32
	Strings  []int32
	Handles  []int32
	Kind     ArrayKind
	empty    bool
	// VarName is the source variable name (uppercase), for runtime errors; may be empty.
	VarName string
}

func productDims(dims []int64) (int64, error) {
	if len(dims) == 0 {
		return 0, fmt.Errorf("array: need at least one dimension")
	}
	n := int64(1)
	for _, d := range dims {
		if d < 1 {
			return 0, fmt.Errorf("array: dimension size must be >= 1, got %d", d)
		}
		if n > MaxArrayCells/d {
			return 0, fmt.Errorf("array: would require more than %d cells", MaxArrayCells)
		}
		n *= d
	}
	return n, nil
}

// NewArray allocates a zero-filled array with the given dimension sizes (each >= 1).
func NewArray(dims []int64) (*Array, error) {
	return NewArrayOfKind(dims, ArrayKindFloat, 0)
}

// NewArrayOfKind allocates an array; emptyStrIdx is the pool index for "" (string kind only).
func NewArrayOfKind(dims []int64, kind ArrayKind, emptyStrIdx int32) (*Array, error) {
	n, err := productDims(dims)
	if err != nil {
		return nil, err
	}
	a := &Array{Dims: append([]int64(nil), dims...), Kind: kind}
	switch kind {
	case ArrayKindFloat, ArrayKindBool:
		a.Floats = make([]float64, n)
	case ArrayKindString:
		a.Strings = make([]int32, n)
		for i := range a.Strings {
			a.Strings[i] = emptyStrIdx
		}
	case ArrayKindHandle:
		a.Handles = make([]int32, n)
	case ArrayKindFloat32:
		a.Floats32 = make([]float32, n)
	default:
		return nil, fmt.Errorf("array: unknown kind")
	}
	return a, nil
}

// NewSharedArrayF32 creates an array that wraps an existing float32 slice.
func NewSharedArrayF32(buf []float32) (*Array, error) {
	if buf == nil {
		return nil, fmt.Errorf("array: nil shared buffer")
	}
	a := &Array{
		Dims:     []int64{int64(len(buf))},
		Floats32: buf,
		Kind:     ArrayKindFloat32,
	}
	return a, nil
}

// TypeName implements HeapObject.
func (a *Array) TypeName() string { return "Array" }

// TypeTag implements HeapObject.
func (a *Array) TypeTag() uint16 { return TagArray }

// Free implements HeapObject.
func (a *Array) Free() {
	a.empty = true
	a.Floats = nil
	a.Floats32 = nil
	a.Strings = nil
	a.Handles = nil
}

// GetHandle returns the heap handle id at indices (ArrayKindHandle only).
func (a *Array) GetHandle(indices []int64) (int32, error) {
	if a.Kind != ArrayKindHandle {
		return 0, fmt.Errorf("array: not a handle array")
	}
	li, err := a.linearIndex(indices)
	if err != nil {
		return 0, err
	}
	return a.Handles[li], nil
}

// SetHandle sets the handle at indices (ArrayKindHandle only).
func (a *Array) SetHandle(indices []int64, hid int32) error {
	if a.Kind != ArrayKindHandle {
		return fmt.Errorf("array: not a handle array")
	}
	li, err := a.linearIndex(indices)
	if err != nil {
		return err
	}
	a.Handles[li] = hid
	return nil
}

func (a *Array) totalLen() int {
	n, err := productDims(a.Dims)
	if err != nil {
		return 0
	}
	return int(n)
}

func (a *Array) errPrefix() string {
	if a.VarName != "" {
		return a.VarName
	}
	return "array"
}

// linearIndex converts 1-based language indices to a flat storage index.
func (a *Array) linearIndex(indices []int64) (int, error) {
	if a.empty {
		return 0, fmt.Errorf("%s: array handle has been freed", a.errPrefix())
	}
	prefix := a.errPrefix()
	if len(indices) != len(a.Dims) {
		return 0, fmt.Errorf("%s: array is %d-dimensional; got %d indices", prefix, len(a.Dims), len(indices))
	}
	var off int64
	for i, langIdx := range indices {
		dimSize := a.Dims[i]
		if langIdx < 1 || langIdx > dimSize {
			return 0, fmt.Errorf("%s: index out of bounds (dimension %d: index %d, valid 1..%d)", prefix, i+1, langIdx, dimSize)
		}
		internal := langIdx - 1
		var stride int64 = 1
		for j := i + 1; j < len(a.Dims); j++ {
			if stride > MaxArrayCells/a.Dims[j] {
				return 0, fmt.Errorf("%s: array index calculation overflow", prefix)
			}
			stride *= a.Dims[j]
		}
		if internal > 0 && stride > (MaxArrayCells-off)/internal {
			return 0, fmt.Errorf("%s: array index calculation overflow", prefix)
		}
		off += internal * stride
	}
	if off < 0 || int64(int(off)) != off {
		return 0, fmt.Errorf("%s: array index calculation overflow", prefix)
	}
	return int(off), nil
}

// GetFloat returns the numeric element at indices (float or bool storage).
func (a *Array) GetFloat(indices []int64) (float64, error) {
	if a.Kind == ArrayKindString || a.Kind == ArrayKindHandle {
		return 0, fmt.Errorf("array: not a numeric array")
	}
	li, err := a.linearIndex(indices)
	if err != nil {
		return 0, err
	}
	return a.Floats[li], nil
}

// GetStringIndex returns the string pool index at indices.
func (a *Array) GetStringIndex(indices []int64) (int32, error) {
	if a.Kind != ArrayKindString {
		return 0, fmt.Errorf("array: not a string array")
	}
	li, err := a.linearIndex(indices)
	if err != nil {
		return 0, err
	}
	return a.Strings[li], nil
}

// Get is an alias for GetFloat (numeric/bool arrays).
func (a *Array) Get(indices []int64) (float64, error) {
	return a.GetFloat(indices)
}

// Set is an alias for SetFloat (numeric/bool arrays).
func (a *Array) Set(indices []int64, v float64) error {
	return a.SetFloat(indices, v)
}

// SetFloat sets a numeric element.
func (a *Array) SetFloat(indices []int64, v float64) error {
	if a.Kind == ArrayKindString {
		return fmt.Errorf("array: not a numeric array")
	}
	li, err := a.linearIndex(indices)
	if err != nil {
		return err
	}
	if a.Kind == ArrayKindBool {
		if v != 0 {
			a.Floats[li] = 1
		} else {
			a.Floats[li] = 0
		}
		return nil
	}
	a.Floats[li] = v
	return nil
}

// SetStringIndex sets a string element by pool index.
func (a *Array) SetStringIndex(indices []int64, idx int32) error {
	if a.Kind != ArrayKindString {
		return fmt.Errorf("array: not a string array")
	}
	li, err := a.linearIndex(indices)
	if err != nil {
		return err
	}
	a.Strings[li] = idx
	return nil
}

// Describe implements fmt.Stringer for errors.
func (a *Array) Describe() string {
	var b string
	for i, d := range a.Dims {
		if i > 0 {
			b += "x"
		}
		b += strconv.FormatInt(d, 10)
	}
	return b
}

// TotalElements returns the number of elements (product of dimensions).
func (a *Array) TotalElements() int { return a.totalLen() }

// DimSize returns the size of the 1-based dimension dim1 (1 = slowest / first index).
func (a *Array) DimSize(dim1 int) (int64, error) {
	if dim1 < 1 || dim1 > len(a.Dims) {
		return 0, fmt.Errorf("array: invalid dimension %d (have %d)", dim1, len(a.Dims))
	}
	return a.Dims[dim1-1], nil
}

// Redim resizes the array in place, optionally preserving leading elements in row-major order.
func (a *Array) Redim(newDims []int64, preserve bool) error {
	if a.empty {
		return fmt.Errorf("array: use after free")
	}
	newCount, err := productDims(newDims)
	if err != nil {
		return err
	}
	oldKind := a.Kind
	oldFloats := a.Floats
	oldStrings := a.Strings
	oldLen := a.totalLen()

	a.Dims = append([]int64(nil), newDims...)

	switch oldKind {
	case ArrayKindFloat, ArrayKindBool:
		a.Floats = make([]float64, newCount)
		if preserve && oldFloats != nil {
			n := oldLen
			if int(newCount) < n {
				n = int(newCount)
			}
			copy(a.Floats, oldFloats[:n])
		}
		a.Strings = nil
	case ArrayKindString:
		a.Strings = make([]int32, newCount)
		for i := range a.Strings {
			a.Strings[i] = 0
		}
		if preserve && oldStrings != nil {
			n := oldLen
			if int(newCount) < n {
				n = int(newCount)
			}
			copy(a.Strings, oldStrings[:n])
		}
		a.Floats = nil
	default:
		return fmt.Errorf("array: unknown kind")
	}
	return nil
}

// FillScalar sets every numeric/bool element to v (ignored for string arrays).
func (a *Array) FillScalar(v float64) error {
	if a.Kind == ArrayKindString {
		return fmt.Errorf("array: use ARRAYFILL with string arrays via builtin")
	}
	for i := range a.Floats {
		if a.Kind == ArrayKindBool {
			if v != 0 {
				a.Floats[i] = 1
			} else {
				a.Floats[i] = 0
			}
		} else {
			a.Floats[i] = v
		}
	}
	return nil
}

// FillStringIndex sets every string element to the same pool index.
func (a *Array) FillStringIndex(idx int32) error {
	if a.Kind != ArrayKindString {
		return fmt.Errorf("array: not a string array")
	}
	for i := range a.Strings {
		a.Strings[i] = idx
	}
	return nil
}

// CopyFrom copies min(len(src), len(dst)) flat elements from src to a (same kind only).
func (a *Array) CopyFrom(src *Array) error {
	if src == nil || a.empty || src.empty {
		return fmt.Errorf("array: invalid copy")
	}
	if a.Kind != src.Kind {
		return fmt.Errorf("array: kind mismatch")
	}
	n := a.totalLen()
	sn := src.totalLen()
	if sn < n {
		n = sn
	}
	switch a.Kind {
	case ArrayKindFloat, ArrayKindBool:
		copy(a.Floats, src.Floats[:n])
	case ArrayKindString:
		copy(a.Strings, src.Strings[:n])
	}
	return nil
}

// Sort1D sorts the flattened array ascending or descending (any rank, sorts linear storage).
func (a *Array) Sort1D(desc bool, strLess func(i, j int32) bool) error {
	if a.Kind == ArrayKindString {
		if strLess == nil {
			return fmt.Errorf("array: string sort requires comparator")
		}
		sort.Slice(a.Strings, func(i, j int) bool {
			if desc {
				return strLess(a.Strings[j], a.Strings[i])
			}
			return strLess(a.Strings[i], a.Strings[j])
		})
		return nil
	}
	sort.Float64s(a.Floats)
	if desc {
		for i, j := 0, len(a.Floats)-1; i < j; i, j = i+1, j-1 {
			a.Floats[i], a.Floats[j] = a.Floats[j], a.Floats[i]
		}
	}
	return nil
}

// Reverse reverses flat storage in place.
func (a *Array) Reverse() {
	switch a.Kind {
	case ArrayKindFloat, ArrayKindBool:
		for i, j := 0, len(a.Floats)-1; i < j; i, j = i+1, j-1 {
			a.Floats[i], a.Floats[j] = a.Floats[j], a.Floats[i]
		}
	case ArrayKindString:
		for i, j := 0, len(a.Strings)-1; i < j; i, j = i+1, j-1 {
			a.Strings[i], a.Strings[j] = a.Strings[j], a.Strings[i]
		}
	}
}

// FindFlat returns the flat index of v or -1 (numeric/bool).
func (a *Array) FindFlat(want float64) int {
	if a.Kind == ArrayKindString {
		return -1
	}
	for i, x := range a.Floats {
		if a.Kind == ArrayKindBool {
			wb := want != 0
			xb := x != 0
			if wb == xb {
				return i
			}
			continue
		}
		if x == want {
			return i
		}
	}
	return -1
}

// FindStringIndex returns flat index or -1.
func (a *Array) FindStringIndex(want int32) int {
	if a.Kind != ArrayKindString {
		return -1
	}
	for i, x := range a.Strings {
		if x == want {
			return i
		}
	}
	return -1
}

// ContainsFlat reports whether a numeric/bool value exists.
func (a *Array) ContainsFlat(want float64) bool {
	return a.FindFlat(want) >= 0
}

// ContainsStringIndex reports whether a pool index exists.
func (a *Array) ContainsStringIndex(want int32) bool {
	return a.FindStringIndex(want) >= 0
}

// Push1D appends one element (1D dynamic arrays only).
func (a *Array) Push1D(f float64, strIdx int32, isBool bool) error {
	if len(a.Dims) != 1 {
		return fmt.Errorf("array: PUSH requires 1D array")
	}
	switch a.Kind {
	case ArrayKindFloat:
		a.Dims[0]++
		a.Floats = append(a.Floats, f)
	case ArrayKindBool:
		a.Dims[0]++
		v := 0.0
		if f != 0 {
			v = 1
		}
		a.Floats = append(a.Floats, v)
	case ArrayKindString:
		a.Dims[0]++
		a.Strings = append(a.Strings, strIdx)
	default:
		return fmt.Errorf("array: unknown kind")
	}
	return nil
}

// Pop1D removes and returns the last element (1D only). ok is false if empty.
func (a *Array) Pop1D() (f float64, strIdx int32, ok bool) {
	if len(a.Dims) != 1 || a.Dims[0] < 1 {
		return 0, 0, false
	}
	switch a.Kind {
	case ArrayKindFloat, ArrayKindBool:
		if len(a.Floats) == 0 {
			return 0, 0, false
		}
		i := len(a.Floats) - 1
		f = a.Floats[i]
		a.Floats = a.Floats[:i]
	case ArrayKindString:
		if len(a.Strings) == 0 {
			return 0, 0, false
		}
		i := len(a.Strings) - 1
		strIdx = a.Strings[i]
		a.Strings = a.Strings[:i]
	default:
		return 0, 0, false
	}
	a.Dims[0]--
	return f, strIdx, true
}

// Shift1D removes and returns the first element (1D only).
func (a *Array) Shift1D() (f float64, strIdx int32, ok bool) {
	if len(a.Dims) != 1 || a.Dims[0] < 1 {
		return 0, 0, false
	}
	switch a.Kind {
	case ArrayKindFloat, ArrayKindBool:
		if len(a.Floats) == 0 {
			return 0, 0, false
		}
		f = a.Floats[0]
		a.Floats = append([]float64(nil), a.Floats[1:]...)
	case ArrayKindString:
		if len(a.Strings) == 0 {
			return 0, 0, false
		}
		strIdx = a.Strings[0]
		a.Strings = append([]int32(nil), a.Strings[1:]...)
	default:
		return 0, 0, false
	}
	a.Dims[0]--
	return f, strIdx, true
}

// Unshift1D prepends one element (1D only).
func (a *Array) Unshift1D(f float64, strIdx int32) error {
	if len(a.Dims) != 1 {
		return fmt.Errorf("array: UNSHIFT requires 1D array")
	}
	a.Dims[0]++
	switch a.Kind {
	case ArrayKindFloat:
		a.Floats = append([]float64{f}, a.Floats...)
	case ArrayKindBool:
		v := 0.0
		if f != 0 {
			v = 1
		}
		a.Floats = append([]float64{v}, a.Floats...)
	case ArrayKindString:
		a.Strings = append([]int32{strIdx}, a.Strings...)
	default:
		return fmt.Errorf("array: unknown kind")
	}
	return nil
}

// Splice1D removes count elements starting at pos (0-based); returns error if out of range.
func (a *Array) Splice1D(pos, count int64) error {
	if len(a.Dims) != 1 {
		return fmt.Errorf("array: SPLICE requires 1D array")
	}
	if count < 0 || pos < 0 {
		return fmt.Errorf("array: invalid splice range")
	}
	n := int64(a.totalLen())
	if pos > n || pos+count > n {
		return fmt.Errorf("array: splice out of range")
	}
	pi := int(pos)
	c := int(count)
	switch a.Kind {
	case ArrayKindFloat, ArrayKindBool:
		a.Floats = append(a.Floats[:pi], a.Floats[pi+c:]...)
	case ArrayKindString:
		a.Strings = append(a.Strings[:pi], a.Strings[pi+c:]...)
	}
	a.Dims[0] -= count
	return nil
}

// Slice1D returns a new 1D array with elements [start:end) (end exclusive).
func (a *Array) Slice1D(start, end int64, emptyStrIdx int32) (*Array, error) {
	if len(a.Dims) != 1 {
		return nil, fmt.Errorf("array: SLICE requires 1D array")
	}
	n := int64(a.totalLen())
	if start < 0 || end < start || end > n {
		return nil, fmt.Errorf("array: slice range out of bounds")
	}
	length := end - start
	out, err := NewArrayOfKind([]int64{length}, a.Kind, emptyStrIdx)
	if err != nil {
		return nil, err
	}
	switch a.Kind {
	case ArrayKindFloat, ArrayKindBool:
		copy(out.Floats, a.Floats[start:end])
	case ArrayKindString:
		copy(out.Strings, a.Strings[start:end])
	}
	return out, nil
}

// FlatToMulti converts a linear index to multi-dimensional indices.
func (a *Array) FlatToMulti(flat int) ([]int64, error) {
	if flat < 0 || flat >= a.totalLen() {
		return nil, fmt.Errorf("array: flat index out of range")
	}
	indices := make([]int64, len(a.Dims))
	rem := flat
	for i := 0; i < len(a.Dims); i++ {
		stride := 1
		for j := i + 1; j < len(a.Dims); j++ {
			stride *= int(a.Dims[j])
		}
		indices[i] = int64(rem / stride)
		rem = rem % stride
	}
	return indices, nil
}

// JoinStrings builds a delimiter-separated string (pool indices resolved via pool).
func (a *Array) JoinStrings(pool []string, delim string) string {
	if a.Kind != ArrayKindString {
		return ""
	}
	var b []byte
	first := true
	for _, ix := range a.Strings {
		if !first {
			b = append(b, delim...)
		}
		first = false
		s := ""
		if int(ix) >= 0 && int(ix) < len(pool) {
			s = pool[ix]
		}
		b = append(b, s...)
	}
	return string(b)
}
