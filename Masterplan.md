# moonBASIC Masterplan — Implementation Roadmap (v1.1.0)
## Current State: Phase B (Engine Core & Native Modules)

moonBASIC is a high-performance, Go-powered game scripting language designed for 2D and 3D game development. This masterplan reflects the consolidated, modular architecture established in Milestones 1–3.

---

### Phase A: Compiler & VM Reconstruction (COMPLETED)
The following foundational components have been successfully refactored and verified:

1.  **The First Law (Case Agnosticism)**: Lexer and SymTable normalize all identifiers and keywords to UPPERCASE.
2.  **Modular Pipeline Library**: All orchestration logic has been extracted from `main.go` into `compiler/pipeline`. This allows the engine to be embedded without CLI dependencies.
3.  **Modular Code Generation**: The generator is split into specialized sub-handlers:
    - `codegen.go` (Base)
    - `codegen_expr.go` (Expressions & Literals)
    - `codegen_stmts.go` (Statements & Control Flow)
    - `codegen_calls.go` (Call Resolution)
4.  **Golden Trace Debugger**: The VM implements a deterministic state-dumping tracer for regression testing.
5.  **Unified Command Manifest**: `commands.json` serves as the Single Source of Truth for both the Semantic Pass and the Runtime Dispatch Registry. Automated registration handles stubbing.
6.  **Stack Hygiene (OpPop)**: Statements no longer leave orphan values on the operand stack.
7.  **MOON bytecode (`.mbc`)**: `vm/moon` implements versioned binary serialization; `pipeline.EncodeMOON` / `DecodeMOON` and CLI `--compile` / `--run` ship bytecode for standalone distribution.

---

### Phase B: Engine Core & Native Modules (ACTIVE)
The goal for Phase B is to implement the actual implementation of engine commands (Window, Render, Physics, Audio).

0.  **Raylib bootstrap (done — first slice)**:
    - [`runtime/window`](runtime/window): `WINDOW.OPEN` / `CLOSE` / `SHOULDCLOSE`, `RENDER.CLEAR` / `FRAME` with `//go:build cgo` vs stub when `CGO_ENABLED=0`.
    - Requires **CGO + C toolchain** for real graphics; see **`ARCHITECTURE.md` §9**.
    - **Acceptance test** (window + input + render loop on the **current** IR, not engine v2): run [`testdata/pretty_window.mb`](testdata/pretty_window.mb) on **Windows x64** and **Linux x64** with `CGO_ENABLED=1` (`go run . testdata/pretty_window.mb`). The sample exits on **ESC** or **window close** (`Window.ShouldClose`). Operator **`NOT` vs `OR`** is specified in **`ARCHITECTURE.md` §7**; full checklist, FPS meaning, **`--info`**, and CI are under **`ARCHITECTURE.md` §9** (“Acceptance test”).
    - Additional manual smoke: [`testdata/rayloop.mb`](testdata/rayloop.mb).

1.  **Rendering Engine (Raylib-Go)**:
    - Extend beyond the first slice: `Render.DrawRectangle`, textures, 3D, etc.
    - Wire handle-based `TEXTURE.LOAD` and `MODEL.LOAD`.
2.  **3D Physics (Jolt)**:
    - Implement the `JOLT.*` suite for 3D rigid bodies and constraints.
3.  **2D Physics (Box2D)**:
    - Implement the `BOX2D.*` suite for 2D gameplay.
4.  **Networking (ENet)**:
    - Implement high-speed packet transfer via the `ENET.*` suite.

---

### Phase C: Optimization & Advanced Analysis
1.  **Static Analysis (CallGraph)**: Use the `Analyzer.CallGraph` for Dead-Code Elimination (DCE).
2.  **Optimization Pass**: Peephole optimization beyond existing constant folding.
3.  **MOON format revisions**: New payload sections or version bumps only when the IR changes (see `ARCHITECTURE.md`).

---

### Phase D: Tooling & Documentation
1.  **The moonBASIC Language Reference**: Detailed documentation of every manifested command.
2.  **Bytecode Disassembler**: Integrated with the CLI for debugging.
3.  **VS Code / Cursor Extensions**: Syntax highlighting and autocomplete using the `commands.json` manifest.

---

## Verification Plan

### Automated Tests
- `go build ./...`: THE FINAL TEST of build integrity.
- `go test ./...`: Full test suite coverage.

### Manual Verification
- `moonbasic --trace script.mbc`: Verify the Golden Trace is clean (depth=0 at halt).
- `moonbasic --check source.mb`: Verify the manifest-driven semantic pass correctly catches typos.
- **Raylib acceptance**: `CGO_ENABLED=1 go run . testdata/pretty_window.mb` on Windows and Linux (see **`ARCHITECTURE.md` §9**). Optional: `moonbasic --info testdata/pretty_window.mb` for runtime banner + bytecode listing; with CGO, throttled **GetFPS** lines on stderr during the loop. **`moonbasic --version`** prints version and the same runtime library line.
