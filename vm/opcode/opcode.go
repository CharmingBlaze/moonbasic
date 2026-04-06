// Package opcode defines the bytecode instruction set and program metadata for the moonBASIC VM.
// Instructions are fixed-width IR v2: OpCode + Flags + pad + Operand (8 bytes total).
// Chunk and Program are defined here so both codegen and vm can import them
// without creating a circular dependency.
package opcode

import (
	"fmt"
	"unsafe"
)

func init() {
	if unsafe.Sizeof(Instruction{}) != 8 {
		panic("opcode: Instruction must be exactly 8 bytes (IR v2)")
	}
}

// OpCode is a single bytecode operation type.
type OpCode byte

//go:generate stringer -type=OpCode
const (
	// Stack manipulation
	OpPushInt    OpCode = iota // Operand: index into Chunk.IntConsts
	OpPushFloat                // Operand: index into Chunk.FloatConsts
	OpPushString               // Operand: index into Program.StringTable
	OpPushBool                 // Operand: 0 = false, 1 = true
	OpPushNull                 // No operands
	OpPop                      // Pop top of stack

	// Variable access (globals use name table, locals use frame slots)
	OpLoadGlobal  // Operand: index into Chunk.Names (interned uppercase string)
	OpStoreGlobal // Operand: index into Chunk.Names
	OpLoadLocal   // Operand: stack frame slot index (local or param)
	OpStoreLocal  // Operand: stack frame slot index

	// Arithmetic
	OpAdd
	OpSub
	OpMul
	OpDiv
	OpMod
	OpPow
	OpNeg

	// Comparison — result is KindBool on stack
	OpEq
	OpNeq
	OpLt
	OpGt
	OpLte
	OpGte

	// Logic
	OpAnd
	OpOr
	OpNot
	OpXor

	// String
	OpConcat

	// Control flow
	OpJump        // Operand: absolute instruction index
	OpJumpIfFalse // Operand: absolute instruction index
	OpJumpIfTrue  // Operand: absolute instruction index

	// Functions & Calls
	OpCallBuiltin // Operand: name index in Chunk.Names; Flags: arg count (0–255)
	OpCallUser    // Operand: name index in Chunk.Names; Flags: arg count
	OpCallHandle  // Operand: method name index in Chunk.Names; Flags: arg count
	OpReturn      // Operand: 1 if returning a value, 0 if void
	OpReturnVoid  // Shorthand for RETURN (0)

	// Arrays
	OpArrayMake // Operand: dimension count
	OpArrayGet  // Operand: dimension count
	OpArraySet  // Operand: dimension count

	// User-defined types
	OpNew      // Operand: type name index in Chunk.Names
	OpDelete   // No operands (pops handle)
	OpFieldGet // Operand: field name index in Chunk.Names
	OpFieldSet // Operand: field name index in Chunk.Names

	// Process control
	OpHalt // Terminate program execution

	// Extended ops (appended to keep earlier opcode values stable)
	OpSwap       // Swap top two operand-stack values
	OpArrayRedim // Operand: dimension count; Flags: 1=preserve contents; stack: handle, dim0..dimN-1 (last dim pushed last)
	OpArrayMakeTyped // Operand: type name index in Chunk.Names; Flags: dimension count; stack: dim sizes (same as ARRAY_MAKE)
	OpNewFilled      // Operand: type name index; Flags: argument count; stack: field values in declaration order
	OpEraseAll       // ERASE ALL — frees entire heap, nulls handle values in globals and operand stack
)

// Instruction is a fixed-width VM decoded unit (8 bytes, IR v2).
// Flags holds the argument count for OpCallBuiltin, OpCallUser, OpCallHandle; else 0.
type Instruction struct {
	Op      OpCode
	Flags   uint8
	_       [2]byte // padding
	Operand int32
}

// String returns a disassembly-style representation of an instruction.
func (i Instruction) String() string {
	return fmt.Sprintf("%-16s %8d %8d", i.Op.String(), i.Operand, int32(i.Flags))
}

// OpCode.String provides human-readable opcode names.
func (op OpCode) String() string {
	names := [...]string{
		"PUSH_INT", "PUSH_FLOAT", "PUSH_STRING", "PUSH_BOOL", "PUSH_NULL", "POP",
		"LOAD_GLOBAL", "STORE_GLOBAL", "LOAD_LOCAL", "STORE_LOCAL",
		"ADD", "SUB", "MUL", "DIV", "MOD", "POW", "NEG",
		"EQ", "NEQ", "LT", "GT", "LTE", "GTE",
		"AND", "OR", "NOT", "XOR", "CONCAT",
		"JUMP", "JUMP_IF_FALSE", "JUMP_IF_TRUE",
		"CALL_BUILTIN", "CALL_USER", "CALL_HANDLE", "RETURN", "RETURN_VOID",
		"ARRAY_MAKE", "ARRAY_GET", "ARRAY_SET",
		"NEW", "DELETE", "FIELD_GET", "FIELD_SET", "HALT",
		"SWAP", "ARRAY_REDIM", "ARRAY_MAKE_TYPED", "NEW_FILLED", "ERASE_ALL",
	}
	if int(op) < 0 || int(op) >= len(names) {
		return fmt.Sprintf("OP_%d", int(op))
	}
	return names[op]
}

// Chunk is a compiled unit of bytecode. A moonBASIC program is composed
// of multiple chunks (main, functions, handlers). It acts as a "compilation unit."
type Chunk struct {
	Name         string
	Instructions []Instruction
	IntConsts    []int64
	FloatConsts  []float64
	Names        []string // Interned uppercase identifiers (variables, functions, fields)
	SourceLines  []int32  // Parallel to Instructions; used for stack traces
}

// AddInt interns an integer constant and returns its index.
func (c *Chunk) AddInt(v int64) int32 {
	for i, x := range c.IntConsts {
		if x == v {
			return int32(i)
		}
	}
	c.IntConsts = append(c.IntConsts, v)
	return int32(len(c.IntConsts) - 1)
}

// AddFloat interns a float constant and returns its index.
func (c *Chunk) AddFloat(v float64) int32 {
	for i, x := range c.FloatConsts {
		if x == v {
			return int32(i)
		}
	}
	c.FloatConsts = append(c.FloatConsts, v)
	return int32(len(c.FloatConsts) - 1)
}

// AddName interns a symbol name (uppercase identifier) and returns its index.
func (c *Chunk) AddName(name string) int32 {
	for i, x := range c.Names {
		if x == name {
			return int32(i)
		}
	}
	c.Names = append(c.Names, name)
	return int32(len(c.Names) - 1)
}

// Emit appends an instruction and its source line, returning the absolute instruction index.
// flags is the arg count for call opcodes; use 0 otherwise.
func (c *Chunk) Emit(op OpCode, operand int32, flags uint8, line int) int {
	c.Instructions = append(c.Instructions, Instruction{Op: op, Flags: flags, Operand: operand})
	c.SourceLines = append(c.SourceLines, int32(line))
	return len(c.Instructions) - 1
}

// Disassemble returns a human-readable listing of the chunk for debugging.
func (c *Chunk) Disassemble() string {
	out := fmt.Sprintf("=== %s ===\n", c.Name)
	for i, instr := range c.Instructions {
		line := int32(-1)
		if i < len(c.SourceLines) {
			line = c.SourceLines[i]
		}
		out += fmt.Sprintf("%04d [L%3d] %s\n", i, line, instr)
	}
	return out
}

// TypeDef describes a user-defined moonBASIC TYPE.
type TypeDef struct {
	Name   string   // uppercase name
	Fields []string // uppercase field names
}

// Program is the complete output of the moonBASIC compiler.
type Program struct {
	StringTable []string // Interned string literals (shared by all chunks); IR v2 SSOT
	// SourcePath is the primary source file path passed to the compiler (used in runtime errors).
	SourcePath string
	Main       *Chunk
	Functions  map[string]*Chunk   // uppercase function name → chunk
	Types      map[string]*TypeDef // uppercase type name → typedef
}

// InternString adds or returns the index of s in the program string pool.
func (p *Program) InternString(s string) int32 {
	for i, x := range p.StringTable {
		if x == s {
			return int32(i)
		}
	}
	p.StringTable = append(p.StringTable, s)
	return int32(len(p.StringTable) - 1)
}

// NewChunk allocates a new empty chunk.
func NewChunk(name string) *Chunk {
	return &Chunk{
		Name: name,
	}
}

// NewProgram allocates a new program container with a main chunk.
func NewProgram() *Program {
	return &Program{
		Main:      NewChunk("<MAIN>"),
		Functions: make(map[string]*Chunk),
		Types:     make(map[string]*TypeDef),
	}
}
