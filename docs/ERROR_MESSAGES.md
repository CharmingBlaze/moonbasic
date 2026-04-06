# moonBASIC errors — what, where, how to fix

This document describes how errors are produced today and the quality bar for new code.

## Compile-time (semantic / type checker)

When a call does not match the manifest (`compiler/builtinmanifest/commands.json`):

- **Unknown command** — message names the dotted key (e.g. `Unknown command 'CAMERA.SETPOSITON'`).
- **Did-you-mean** — if another command in the same namespace is within edit distance ≤ 3 of the method name, the compiler suggests it (`compiler/semantic/cmdhint.go`).
- **Wrong arity** — `NS.METHOD: no overload matches N argument(s)` plus an arity hint from the manifest.
- **Wrong argument types** — which argument index, expected vs got, and a short fix hint.
- **Record types (`TYPE` … `ENDTYPE`)** — unknown field name, wrong **`TypeName(...)`** field count vs the type definition, or field type mismatch (messages name the type and field).

Source location always includes **file, line, column** and a **source line** excerpt when available (`compiler/errors`).

## Runtime (VM + native builtins)

When a builtin returns an error, the VM wraps it with **host source path** and **line** for the **call site** (the bytecode instruction that invoked the native), when line information exists:

```text
[moonBASIC] Error in game.mb line 14:
  handle 3 is Mesh, but this operation requires a different resource type
  Hint: Pass the handle returned by the matching MAKE/LOAD ...
```

Implementation: `vm/vm.go` (`runtimeError`), using `Program.SourcePath` / chunk name and `Chunk.SourceLines`.

### Heap handles (`heap.Cast`, invalid / stale / wrong type)

`vm/heap/heap.go` `Cast` errors are written to answer:

1. **What** — null handle, stale handle, or wrong resource type.
2. **Where** — handle id and actual `TypeName` when relevant.
3. **How to fix** — short `Hint:` lines (avoid 0; avoid use-after-free; pass the matching handle class).

### Unknown runtime command

Dynamic dispatch misses use `runtime.FormatUnknownRegistryCommand` (`runtime/suggest_cmd.go`) with **did-you-mean** against registered keys (distance ≤ 3).

### Stubs and unimplemented commands

Commands registered from the manifest without a native body may return:

`[moonBASIC] Runtime Error: command X is not yet implemented`

## Conventions for new native code

- Prefer **`runtime.Errorf(...)`** for user-facing messages so the `[moonBASIC]` prefix is consistent where used.
- After **`heap.Cast`**, return the error unchanged so hints stay intact; the VM adds **file/line** when it surfaces the error.
- State the **operation** in the first line (e.g. `CAMERA.SETPOS: non-numeric position`).
- For argument count mistakes, say **expected** shape explicitly (names + types).

## Related docs

- [API_CONSISTENCY.md](./API_CONSISTENCY.md) — naming conventions and full command listing (generated).
- Roadmap and phases: [ROADMAP.md](./ROADMAP.md).
