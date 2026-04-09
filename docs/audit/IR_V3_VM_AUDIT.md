# IR v3 VM audit (gap analysis vs ENGINE_IR_V3.md)

This document confirms alignment between [ENGINE_IR_V3.md](../../ENGINE_IR_V3.md), the compiler, and [vm/vm.go](../../vm/vm.go), and lists legacy naming or paths to retire or document.

## Confirmed alignment

Tier-1 platforms for shipping and CI are **Windows x64** and **Linux x64**; IR v3 / MOON are defined to be **portable** across them (see [ENGINE_IR_V3.md](../../ENGINE_IR_V3.md) § Platform portability).

| Item | Location | Notes |
|------|----------|--------|
| Instruction width | [vm/opcode/opcode.go](../../vm/opcode/opcode.go) `Instruction` | 8 bytes; `init()` asserts `unsafe.Sizeof(Instruction{}) == 8`. |
| Register file | VM frame in [vm/callstack](../../vm/callstack/) | Per-frame registers; `Dst`/`SrcA`/`SrcB` in each opcode. |
| MOON v3 | [vm/moon/](../../vm/moon/) | Header version `0x00030000` per ARCHITECTURE. |
| Codegen | [compiler/codegen/](../../compiler/codegen/) | Register allocation; `Emit(op, dst, srcA, srcB, operand, line)`. |
| Value model | [vm/value/value.go](../../vm/value/value.go) | 24-byte tagged union; `float64` for float storage. |

## Opcodes with “stack-era” names (implementation is register-based)

Several opcode names retain **stack** vocabulary from IR v2 naming, but the VM implements them as **register** operations:

| Opcode | VM behavior ([vm/vm.go](../../vm/vm.go)) |
|--------|------------------------------------------|
| `OpPushInt` / `OpPushFloat` / `OpPushString` / `OpPushBool` / `OpPushNull` | Load immediate or pool value into **`Dst`** register. |
| `OpPop` | **No-op** (reserved; codegen should not rely on stack depth). |
| `OpLoadLocal` / `OpStoreLocal` | **Error** at runtime: “use register moves instead” — indicates stale bytecode if hit. |
| `OpSwap` | Swaps **`SrcA`** and **`SrcB`** register contents. |

**Action:** No retirement required for correctness; optional future rename to `OpLoadImm*` / `OpMoveImm` in a breaking IR v4. Document for compiler engineers.

## Complex dispatch

Arithmetic, compares, calls, arrays, and control flow are handled in [vm/vm_arith.go](../../vm/vm_arith.go), [vm/vm_control.go](../../vm/vm_control.go), and `dispatchComplex` — all register-oriented.

## Peephole / performance follow-ups (non-blocking)

- [compiler/opt/](../../compiler/opt/) — redundant `OpMove` removal per ENGINE_IR_V3.
- Reduce `OpCallBuiltin` overhead via manifest inlining for tiny pure helpers (future).

## Conclusion

The execution model is **register-based IR v3** as documented. Remaining “stack” labels are **historical names** on load-style opcodes, not a separate stack machine interpreter.
