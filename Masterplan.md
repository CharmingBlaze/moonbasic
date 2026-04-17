# moonBASIC Masterplan — Implementation Roadmap (v1.1.0)
## Current State: Phase B (Engine Core & Native Modules)

moonBASIC is a high-performance, Go-powered game scripting language designed for 2D and 3D game development. This masterplan reflects the consolidated, modular architecture established in Milestones 1–3.

---

### Phase A: Compiler & VM Reconstruction (COMPLETED)
The following foundational components have been successfully refactored and verified:

1.  **The First Law (Case Agnosticism)**: Source is case-insensitive; the lexer keeps **lowercase** canonical identifier/keyword spellings, SymTable uppercases **symbol** keys for lookup, and the manifest/registry use **uppercase** dotted command names.
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

#### Handle method chaining (API consistency)
- **Rule:** Builtins that mutate an object through a leading **handle** argument should return that handle on success (`args[0]`) so scripts can chain: `cam = CAMERA.CREATE().pos(0, 10, 20).look(0, 0, 0).fov(60)`.
- **Status (Part 1 — 2026-04):** `runtime/camera/raylib_cgo.go` core `CAMERA.SET*` paths already returned the handle; extended fixes landed in `camera_blitz_cgo.go`, `camera_more_cgo.go`, `orbit_follow_cgo.go`, `cam2d_cgo.go`, `blitz_extra_cgo.go` (`CAMERA.LOOKATENTITY`), plus `BODY3D.ADDMESH` / no-op physics paths in `runtime/physics3d/jolt_body3d_cgo.go`. `runtime/sprite/raylib_cgo.go`, `runtime/mblight/builtins.go`, `runtime/mbmodel3d/model_transform_cgo.go`, and `runtime/physics2d/box2d.go` body setters were already handle-chaining–friendly.
- **Terrain:** `ModifyTerrain`, `TERRAIN.APPLYMAP`, deferred `CHUNK.GENERATE` out-of-range, and `TerrainShading` stub return the terrain handle where appropriate for chaining. `handleCallBuiltin` lists common `TERRAIN.*` operations on `TagTerrain`; `normalizeHandleMethod` maps `DETAIL`→`SETDETAIL` (zero-arg `.detail()` still resolves to `TERRAIN.GETDETAIL` via `handleCallDispatch`).
- **Water:** `WATER.DRAW` returns the water handle (`args[0]`) for chaining with other `WATER.SET*` helpers. `vm/handlecall.go` maps `TagWater` methods (`SETPOS`, `DRAW`, `SETWAVE` / `.wave()`, colors, queries); `WATER.UPDATE(dt)` stays global (no per-handle receiver).
- **Model / material (continued):** `MODEL.SET*` helpers in `model_render_hierarchy_cgo.go`, `model_texture_stages_cgo.go`, `model_material_cgo.go`, `material_cmds_cgo.go` (MATERIAL setters), `model_lod_cgo.go`, deferred draw paths in `model_inst_draw_cgo.go`, and early exits in `model_complete_cgo.go` (`DRAWAT`/`DRAWEX` hidden, `UPDATEANIM` no-op) now return the leading handle on success where applicable. `vm/handlecall.go` maps extended script methods on `TagModel` / `TagLODModel` (e.g. `SETCULL`, texture stages, `SETMATERIAL`, `MOVE`, `ATTACHTO`, `SETLODDISTANCES`) and `TagMaterial` (`MATERIAL.SETSHADER`, `SETTEXTURE`, …); `normalizeHandleMethod` adds short aliases (`cull`→`SETCULL`, `attach`→`ATTACHTO`, …). `vm/handlecall_chaining_test.go` covers dispatch + aliases.
- **Intentionally not chained:** `CAMERA.END`, `CAMERA2D.END`, `MATRIX.FREE`, `BODY3D.FREE`, etc. (no receiver / destructive free).
- **Parts 2–8** (MODEL/LIGHT getters, manifest, tests, docs): see task backlog; `vm/handlecall.go` already maps many `TagModel`, `TagLight`, `TagParticle`, `TagTerrain`, `TagInstancedModel` method names to registry keys.

0.  **Raylib bootstrap (done — first slice)**:
    - [`runtime/window`](runtime/window): `WINDOW.OPEN` / `CLOSE` / `SHOULDCLOSE`, `RENDER.CLEAR` / `FRAME` with `//go:build cgo` vs stub when `CGO_ENABLED=0`.
    - Requires **CGO + C toolchain** for real graphics; see **`ARCHITECTURE.md` §9**.
    - **Acceptance test** (window + input + render loop on the **current** IR, not engine v2): run [`testdata/pretty_window.mb`](testdata/pretty_window.mb) on **Windows x64** and **Linux x64** with `CGO_ENABLED=1` using the **full runtime**, e.g. `go run -tags fullruntime ./cmd/moonrun testdata/pretty_window.mb` (plain `go run .` only compiles to `.mbc`). The sample exits on **ESC** or **window close** (`Window.ShouldClose`). Operator **`NOT` vs `OR`** is specified in **`ARCHITECTURE.md` §7**; full checklist, FPS meaning, **`--info`**, and CI are under **`ARCHITECTURE.md` §9** (“Acceptance test”).
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
3.  **VS Code extension**: Syntax highlighting and autocomplete using the `commands.json` manifest.

---

## Verification Plan

### Automated Tests
- `go build ./...`: THE FINAL TEST of build integrity.
- `go test ./...`: Full test suite coverage.

### Manual Verification
- `moonrun --trace script.mbc`: Verify the Golden Trace is clean (depth=0 at halt) — requires a **fullruntime**-built `moonrun` (or `go run -tags fullruntime ./cmd/moonrun --trace …`).
- `moonbasic --check source.mb`: Verify the manifest-driven semantic pass correctly catches typos (compiler-only binary).
- **Raylib acceptance**: `CGO_ENABLED=1 go run -tags fullruntime ./cmd/moonrun testdata/pretty_window.mb` on Windows and Linux (see **`ARCHITECTURE.md` §9**). Optional: fullruntime **`moonbasic --info`** on a script for runtime banner + bytecode listing before run. **`moonbasic --version`** prints compiler version (use **`moonrun --version`** for the game runtime line when using release binaries).
