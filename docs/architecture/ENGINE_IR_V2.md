# Engine IR v2 / MOON v2 (shipped)

This document describes the **current** VM/compiler IR after Phase G. Older **MOON 1.x** bytecode (`0x00010000`) is **rejected** at load time with a message to recompile from `.mb`.

## MOON container (`vm/moon`)

- **Version**: `0x00020000` (`moon.Version`).
- **Header**: unchanged 16-byte layout (magic `MOON`, big-endian version, flags, entry offset).
- **Payload** (v2): program **string table** (count + UTF-8 strings) first, then encoded chunks. Each chunk no longer embeds its own string constant pool; `OpPushString` operands index **`Program.StringTable`**.

## Instruction (`vm/opcode`, 8 bytes)

- Layout: `Op uint8`, `Flags uint8`, padding `[2]byte`, `Operand int32`.
- **Call opcodes** (`CALL_BUILTIN`, `CALL_USER`, `CALL_HANDLE`): `Operand` = name index in `Chunk.Names`; `Flags` = argument count (0–255).
- MOON file stores per instruction: op, flags, 2-byte zero pad, operand (BE), source line (BE).

## Value (`vm/value`, 24 bytes)

- Layout: `Kind uint8`, `[7]byte` pad, `IVal int64`, `FVal float64`.
- **`KindString`**: `IVal` holds the **`int32` index** into `Program.StringTable`. Runtime builtins resolve text via `runtime.ArgString` / `RetString` (active `Registry.Prog` during `Call`).

## Compiler support

- **`compiler/intern`**: lexer interns identifier/keyword spellings per compile.
- **`compiler/arena`**: parser allocates AST nodes from an arena; **`pipeline.CompileSource` / `CheckSource`** call **`arena.Reset()`** after the pipeline — do not retain AST pointers past compile.
- **`compiler/opt`**: bytecode peephole (e.g. push/pop) and jump threading; invoked at end of codegen.

## Heap (`vm/heap`)

- Handles are opaque **`int32`**: high 16 bits = **generation** (incremented on `Free`), low 16 bits = **slot**. **`Get`** / **`Cast`** fail on generation mismatch (stale handle).

## Reference architecture

- **`unsafe.Sizeof`** is asserted in **`vm/opcode`** (8) and **`vm/value`** (24) on init; primary target is **amd64**.

See **`ARCHITECTURE.md`** for the full pipeline and manifest rules.
