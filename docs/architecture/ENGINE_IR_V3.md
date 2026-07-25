# Engine IR v3 / MOON v3 (shipped)

This document describes the **register-based** VM and compiler IR. 
Older **MOON v2** bytecode (`0x00020000`) is supported for loading but internally converted or deprecated. IR v3 is the current performance standard.

## Platform portability (tier-1: Windows + Linux)

- **Bytecode and `.mbc` files are OS-neutral**: the MOON container, fixed-width opcodes, and value layout do not depend on `GOOS` or host pointer width (tier-1 targets are **Windows x64** and **Linux x64**; see `compiler/errors/MoonBasic.md`).
- **Source line endings**: the parser normalizes `\r\n` / `\r` to `\n` before lexing, so the same source compiles to the same IR on Windows and Linux.
- **Regression tests**: `vm/moon/ir_v3_roundtrip_test.go` exercises compile → MOON encode → decode without importing the full runtime registry (runs even when Raylib is not installed).

## MOON container (`vm/moon`)

- **Version**: `0x00030000` (`moon.Version`).
- **Header**: unchanged 16-byte layout (magic `MOON`, big-endian version, flags, entry offset).
- **Payload** (v3): Identical to v2 (string table then chunks), but optimized for 8-byte fixed-width register instructions.

## Instruction (`vm/opcode`, 8 bytes)

- **Layout**: `Op uint8`, `Dst uint8`, `SrcA uint8`, `SrcB uint8`, `Operand int32`.
- **Register File**: 256 virtual registers (R000–R255).
- **Control Flow**: `Operand` typically holds the absolute instruction index for jumps.
- **Calls**: `Dst` holds the return register index; `SrcA` holds the base of the argument register block; `Operand` encodes argument count and name index.

## Value (`vm/value`, 24 bytes)

- **Layout**: Unchanged from v2 (`Kind uint8`, `[7]byte` pad, `IVal int64`, `FVal float64`).
- **Memory Safety**: Register-based execution eliminates stack-depth errors but requires strict tracking of register lifetimes in the compiler.

## Compiler support

- **`compiler/codegen`**: Moved to a "reset-on-statement" register allocation strategy. 
  - `baseReg` marks the end of local/param storage.
  - `nextReg` tracks temporaries, returning to `baseReg` after each statement to keep the register file shallow.
- **`compiler/opt`**: Optimized for register-based PEEPHOLE passes (e.g., removing redundant `OpMove` instructions).

## Heap (`vm/heap`)

- **Handles**: Handle-based memory remains the core safety mechanism. Opaque `int32` handles with generation checks prevent use-after-free and raw pointer hazards.

## Reference architecture

- **Performance**: Register-based execution reduces total instructions per line of code by ~30%, enabling 4K/144Hz performance on modern CPUs (Goal: Milestone 2).
- **`unsafe.Sizeof`** remains asserted for `opcode.Instruction` (8) and `value.Value` (24).

See **`ARCHITECTURE.md`** for the full pipeline and modern syntax rules.

**See also:** [docs/audit/IR_V3_VM_AUDIT.md](docs/audit/IR_V3_VM_AUDIT.md) — implementation audit (legacy opcode names vs register behavior).
