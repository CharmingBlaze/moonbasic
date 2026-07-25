# Codebase map

Quick answers for contributors navigating moonbasic-compiler.

## Compiler

| Question | Location |
|----------|----------|
| Where is the lexer? | [`compiler/lexer/`](../../compiler/lexer/) |
| Where are tokens? | [`compiler/token/`](../../compiler/token/) |
| Where is the parser? | [`compiler/parser/`](../../compiler/parser/) |
| Where is the AST? | [`compiler/ast/`](../../compiler/ast/) |
| Where is semantic analysis? | [`compiler/semantic/`](../../compiler/semantic/) |
| Where is the symbol table? | [`compiler/symtable/`](../../compiler/symtable/) |
| Where is type checking helpers? | [`compiler/types/`](../../compiler/types/) |
| Where is code generation? | [`compiler/codegen/`](../../compiler/codegen/) |
| Where is the optimizer? | [`compiler/opt/`](../../compiler/opt/) |
| Where is INCLUDE expansion? | [`compiler/include/`](../../compiler/include/) |
| Where is the compile pipeline? | [`compiler/pipeline/`](../../compiler/pipeline/) |
| Where are structured diagnostics? | [`compiler/errors/`](../../compiler/errors/) |
| Where are built-in commands declared? | [`compiler/builtinmanifest/commands.json`](../../compiler/builtinmanifest/commands.json) |

## VM / bytecode

| Question | Location |
|----------|----------|
| Where is bytecode encoded (opcodes)? | [`vm/opcode/`](../../vm/opcode/) |
| Where is bytecode executed? | [`vm/moon/`](../../vm/moon/) (and related packages under [`vm/`](../../vm/)) |
| Where is handle method dispatch (shared)? | [`handlecall/`](../../handlecall/) + [`vm/handlecall.go`](../../vm/handlecall.go) |

## Runtime

| Question | Location |
|----------|----------|
| Where are built-in commands implemented? | [`runtime/`](../../runtime/) packages (`draw`, `window`, `mbentity`, `physics3d`, …) |
| Where are runtime modules registered? | [`compiler/pipeline/registry.go`](../../compiler/pipeline/registry.go) (`fullruntime` tag) |
| Where is the HAL? | [`hal/`](../../hal/), drivers in [`drivers/`](../../drivers/) |
| Where is Jolt / physics sync documented? | [`docs/PHYSICS.md`](../PHYSICS.md) |

## Tools and docs generation

| Question | Location |
|----------|----------|
| Where are command docs generated? | [`tools/apidoc/`](../../tools/apidoc/) |
| Where is the IDE docs bundle synced? | [`tools/docsexport/`](../../tools/docsexport/) → `ide/bundled-docs/` |
| Where is IDE language data exported? | [`tools/ideexport/`](../../tools/ideexport/) |
| Where is manifest vs runtime audited? | [`tools/diff_manifest_runtime.py`](../../tools/diff_manifest_runtime.py) |
| Where is the doc site builder? | [`cmd/moondoc/`](../../cmd/moondoc/) |

## Products / UX

| Question | Location |
|----------|----------|
| Where is the desktop IDE? | [`ide/`](../../ide/) (Wails; formerly `moonbasic ide/`) |
| Where is the VS Code extension? | [`editors/vscode-moonbasic/`](../../editors/vscode-moonbasic/) |
| Where is the LSP? | [`lsp/`](../../lsp/) |
| Where is the debugger (DAP)? | [`dap/`](../../dap/) |
| Compiler CLI (default `go run .`)? | Root [`main.go`](../../main.go) + [`cmd/moonbasic/`](../../cmd/moonbasic/) |
| Game runtime player? | [`cmd/moonrun/`](../../cmd/moonrun/) |

## Release and packaging

| Question | Location |
|----------|----------|
| Where are release archives built? | [`.github/workflows/release.yml`](../../.github/workflows/release.yml), [`scripts/release/`](../../scripts/release/), [`scripts/packaging/`](../../scripts/packaging/) |
| Where are IDE zip extras? | [`packaging/`](../../packaging/) (`START-IDE.*`, `samples/`, READMEs) |
| Where are platform libraries stored? | [`third_party/`](../../third_party/) (Jolt, ENet, raylib sources) |

## Examples and tests

| Question | Location |
|----------|----------|
| Where are examples? | [`examples/`](../../examples/) |
| Where are compiler/runtime fixtures? | [`testdata/`](../../testdata/) |
| Where are benchmarks? | [`benchmarks/`](../../benchmarks/) |

## Documentation layout

| Area | Path |
|------|------|
| Getting started | [`docs/GETTING_STARTED.md`](../GETTING_STARTED.md), [`docs/BEGIN_HERE.md`](../BEGIN_HERE.md) |
| Language | [`docs/LANGUAGE.md`](../LANGUAGE.md), [`docs/PROGRAMMING.md`](../PROGRAMMING.md), root [`STYLE_GUIDE.md`](../../STYLE_GUIDE.md) |
| API reference | [`docs/reference/`](../reference/), [`docs/api/`](../api/) |
| Systems guides | [`docs/systems/`](../systems/) |
| Architecture | [`docs/architecture/`](../architecture/), root [`ARCHITECTURE.md`](../../ARCHITECTURE.md) |
| Development / cleanup | [`docs/development/`](./) |
| Audits | [`docs/audit/`](../audit/) |

## Build modes (reminder)

| Mode | Entry |
|------|--------|
| Compiler-only (default) | `go build .` / `go build ./cmd/moonbasic` |
| Full runtime | `go build -tags fullruntime ./cmd/moonrun` |
| Both paths before push | `scripts/build/check_builds.ps1` or `.sh` |

See also [CONTRIBUTING.md](../../CONTRIBUTING.md) and [DEVELOPER.md](../DEVELOPER.md).
