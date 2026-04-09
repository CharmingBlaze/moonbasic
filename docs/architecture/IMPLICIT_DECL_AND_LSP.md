# Implicit declaration, type inference, and LSP integration

This document specifies how moonBASIC’s **modern syntax** (optional explicit `VAR`) interacts with **`DIM`**, **`TYPE`**, and tooling, and how the **symbol table** is exposed for LSP.

## Current implementation (compiler)

| Component | Path | Role |
|-----------|------|------|
| Pipeline options | [compiler/pipeline/compile.go](../../compiler/pipeline/compile.go) `CompileOptions` | `ImplicitDeclaration` and `TypeInference` (defaults **on** in `CompileSource`). |
| Two-pass symbol builder | [compiler/symtable/builder.go](../../compiler/symtable/builder.go) | Pass 1: predeclare `FUNCTION` / `TYPE`. Pass 2: global assignments, `DIM`, `CONST`, `LOCAL`, loop vars; Pass 3: function-local symbols. |
| Symbol table | [compiler/symtable/symtable.go](../../compiler/symtable/symtable.go) | `Table`, `Symbol` (`Kind`, `Type`, `Slot`). **`ExportJSON()`** for tooling. |
| Semantic analysis | [compiler/semantic/analyze.go](../../compiler/semantic/analyze.go) | Validation after symbol build. |
| Codegen | [compiler/codegen/](../../compiler/codegen/) | Uses `NewWithSymbols` when implicit table is present. |

## First-assignment rule (globals)

- At program scope, an **`AssignNode`** to an unknown identifier **declares** a global variable (`DefineGlobalVar`), with **type inference** from suffix (`#`, `$`, `?`) or expression (`inferType`).
- **`DIM`** / **`DIM AS`** still declare arrays and typed arrays explicitly.
- **`CONST`** and **`TYPE`** use dedicated nodes and take precedence over implicit rules.

## Interaction with `DIM` and arrays

- **`DIM name(size)`** declares an array; the builder treats it as a global with `types.Array`.
- Implicit assignment does **not** create an array without **`DIM`** — assigning to `a(i)` requires prior **`DIM a(...)`** (or REDIM where supported) per language rules in semantic analysis.

## Interaction with `FUNCTION` / locals

- Locals are collected per function in **`collectFunctionLocals`**; parameters and **`LOCAL`** statements populate slots.

## Float64 and FFI

- VM **`value.Value`** stores floats as **`float64`**. Runtime builtins often cast to **`float32`** at the Raylib/Jolt boundary. A future “float64-only script” policy would require auditing every **`argFloat`** / cast in `runtime/`.

## LSP integration plan

1. **JSON export:** [symtable.Table.ExportJSON](../../compiler/symtable/symtable.go) exposes globals, funcs, types.
2. **Pipeline helpers:** [compiler/pipeline/symbols_export.go](../../compiler/pipeline/symbols_export.go) — **`BuildSymbolTable`**, **`ExportSymbolTableJSON`**, **`DocumentSymbols`** (parse + INCLUDE + symbol builder, no codegen).
3. **LSP implementation:** [lsp/server.go](../../lsp/server.go) advertises **`documentSymbolProvider`** and **`definitionProvider`**; [lsp/symbols.go](../../lsp/symbols.go) maps symbols to **best-effort** line ranges via source scan (`FUNCTION` / `TYPE` / `CONST` / `DIM` / first assignment).

## Future work

- Store **source line + column** on each `Symbol` for precise **go-to-definition** and **references**.
- Dual-pass **semantic** symbol JSON merged with builder output for diagnostics parity.
