# moonBASIC Architecture (v1.3.0)
## Mandatory architecture reference

This document defines the **Ground Truth** for the moonBASIC compiler and runtime. Any changes must adhere to the modular structure and stable APIs defined here. **DO NOT REVERT TO OLD MONOLITHIC VERSIONS.**

### Notes for maintainers

- Treat this file as **authoritative** over informal notes or older write-ups. If the repo already matches this document, **do not** “restore” or replace `main.go` / `pipeline.go` with older orchestration patterns.
- **`compiler/pipeline`** intentionally imports **`runtime`** here so **`RunProgram`** is a one-call embedder entrypoint. Do not delete `EncodeMOON` / `DecodeMOON` or stub out `--compile` / `--run` in `main.go` if this section lists them as shipped.
- **`CallStmtNode`** must delegate to **`emitCallStmt`** in **`codegen_stmts.go`**; an empty `case` is a **bug** (symptoms: `PRINT` produces no bytecode).
- **`commands.json`** supports **multiple rows per `key`** (overload arities). Semantic analysis uses **`LookupArity`**; do not assume **`Table.Commands[key]`** is a single struct.
- Logical **`AND` / `OR` / `XOR`** must compile to **`OpAnd` / `OpOr` / `OpXor`** (**§5**, **§7.1**). Handle calls **`recv.METHOD(...)`** require **`handleCallBuiltin`** mapping (**§8.2**); **`Camera3D.Begin`** is **not** a registry key.

---

### 1. The First Law: Case Agnosticism
- **Rule**: Every keyword and identifier in moonBASIC is case-agnostic in source.
- **Implementation**: The lexer records **canonical lowercase** `Lit` text for identifiers and keywords (keyword matching still uses an uppercase view of the scanned letters). The **symbol table** uppercases names for variable/function lookup. Built-in **`NAMESPACE.NAME`** resolution uses **uppercase** manifest/registry keys (`NormalizeCommand`, etc.). String literal contents are not case-normalized.
- **Validation**: Do not assume every compiler stage uses the same casing; follow each subsystem (AST names lowercase; symtable keys uppercase; registry keys uppercase).

---

### 2. The Pipeline API (`compiler/pipeline`)
The orchestration logic must live in `compiler/pipeline/pipeline.go`. The CLI driver (`main.go`) must be a thin wrapper around these functions.

- **`Options`**:
  - `Debug bool` (Bytecode disassembly; when true, the CGO window module also logs throttled **`GetFPS`** lines to **`Out`** during **`RENDER.FRAME`**)
  - `Trace bool` (VM state trace)
  - `Out io.Writer` (Output stream)
- **Functions**:
  - `CompileSource(name, src string) (*opcode.Program, error)`
  - `CompileFile(path string) (*opcode.Program, error)`
  - `CheckFile(path string) error`
  - `RunProgram(prog *opcode.Program, opts Options) error`
  - `EncodeMOON(prog *opcode.Program) ([]byte, error)` — MOON container serialization (`.mbc`)
  - `DecodeMOON(data []byte) (*opcode.Program, error)` — validates **MOON** magic + version, then decodes payload

**CLI (`main.go`)**: `--compile` writes `<stem>.mbc` next to the source; `--run <file.mbc>` or a **positional** `moonbasic game.mbc` decodes and runs via `RunProgram`. Default: compile from source and run in memory.

#### INCLUDE (compile-time)

- **`INCLUDE "path.mb"`** is expanded in **`compiler/include`** immediately after parse (`ExpandWithArena`), before semantic analysis.
- Paths resolve relative to the **including** file; if missing, **`TryOpenInclude`** searches **`MOONBASIC_PATH`** and package roots (**`SyncPackageIncludeRoots`** in **`compiler/pipeline/pkgpath.go`**).
- **Duplicate includes** of the same resolved absolute path are **skipped** (include guard): one parse, one merge — no duplicate top-level statements.
- **`CompileSource`** may use a single **`arena.Arena`** for the merged parse; the arena is reset after bytecode is produced — no ongoing retention of source text in the VM.

---

### 3. MOON bytecode (`.mbc`)

- **Package**: `vm/moon` — binary schema (not `encoding/gob` for shipping).
- **Header** (16 bytes): magic `MOON`, big-endian version (**`0x00030000`** for IR v3; IR v2 supported until v1.4), reserved flags, entry offset.
- **IR v3 payload**: program-level string table, then chunks; 8-byte register-based instructions; see **`ENGINE_IR_V3.md`**.
- **Loader** validates header before building VM state so launchers can reject wrong engines quickly.

---

### 4. The Registry Manifest (`commands.json`)
`compiler/builtinmanifest/commands.json` is the **Single Source of Truth** for all built-in commands.

- **Overloads**: The same canonical **`"key"`** (e.g. **`RENDER.CLEAR`**) may appear **multiple times** with different **`"args"`** arrays. The JSON loader (`compiler/builtinmanifest/manifest_json.go`) builds **`Table.Commands`** as **`map[string][]Command`**: one map entry per dotted key, value = ordered overload list.
- **Semantic pass** (`compiler/semantic/analyze.go`): Resolves **`NamespaceCallStmt`** and **`NamespaceCallExpr`** with **`LookupArity(ns, method, len(args))`**. If the namespace method exists but no overload matches the argument count → **Compile Error** with **`ArityHint`** (lists valid arities). Unknown **`NS.METHOD`** → **Compile Error** with did-you-mean (**`compiler/semantic/cmdhint.go`**).
- **Runtime dispatch** (`runtime/runtime.go`): **`RegisterFromManifest`** walks every overload but registers **at most one stub per canonical `Command.Key`** — the first-seen **`Namespace`** wins for the stub closure. Natives that support multiple arities implement **branching on `len(args)`** in one **`BuiltinFn`** (e.g. **`RENDER.CLEAR`** in **`runtime/window/raylib_cgo.go`**).
- **LSP** (`lsp/server.go`): Hover for a dotted builtin uses **`FirstOverload(key)`** when arity cannot be inferred from the line; multi-arity commands may show only the first signature unless tooling is extended.

**Inventory tooling**: From the repo root, **`python tools/gen_master_audit.py`** regenerates **`MASTER_AUDIT.txt`** (manifest keys vs **`Register("KEY"`** in **`runtime/**/*.go`**, excluding **`*_test.go`**) plus **`MASTER_AUDIT_REGISTERED.txt`**, **`MASTER_AUDIT_MANIFEST.txt`**, **`MASTER_AUDIT_DUPLICATES.txt`**, and **`docs/audit/REFERENCE_KEY_COVERAGE.txt`** (which manifest keys appear verbatim in **`docs/reference/*.md`** and **`compiler/errors/MoonBasic.md`**).

---

### 4.1 The Easy Mode Layer
MoonBASIC provides a high-level "Easy Mode" API alongside the namespaced engine core.
- **Rule**: Every namespaced command (e.g. `ENTITY.CREATECUBE`) should have a flat, Blitz3D-style global alias (e.g. `CreateCube`) if it is commonly used in game loops.
- **Implementation**:
    - Aliases are registered in the same runtime modules as the core commands.
    - Global aliases DO NOT have a namespace prefix in `commands.json`.
    - Shorthands (like `ENTX`, `ENTY`) are used to simplify access to handle-based fields.
- **Dual-Layer Registration**: When registering a module, provide both the namespaced (canonical) key and the flat (Easy Mode) key to the `Registrar`.

---

### 5. Modular Code Generation
Codegen is split into specialized sub-handlers. It uses a **Register-Based** model (IR v3):
- `codegen.go`: Structural base, orchestration, and register management (`baseReg`, `nextReg`, `allocReg`).
- `codegen_expr.go`: Translates expressions into register-based IR. Every `emitExpr` call returns a register index.
- `codegen_stmts.go`: Handles statements. Resets `nextReg = baseReg` at statement boundaries to minimize register pressure.
- `codegen_calls.go`: Handles built-in and user-function call resolution using contiguous register blocks for arguments.

---

### 6. Register Allocation
- **Strategy**: Static-Base, Dynamic-Temp. 
- **Locals**: Local variables and parameters are assigned fixed registers at the start of the frame (R0, R1, ...).
- **Temporaries**: The compiler maintains a pool of temporary registers starting after the local slots.
- **Reset-on-Statement**: To prevent register exhaustion, the compiler resets the temporary register counter (`nextReg = baseReg`) at the start of most statement nodes.

---

### 7. Expression operator precedence (language rule)

This order is implemented by **`compiler/parser/parser_expr.go`** (`parseExpr` → `parseOr` → `parseXor` → `parseAnd` → `parseNot` → …). The long-form spec **`compiler/errors/MoonBasic.md`** must list the same levels. **Tighter binding = higher in the table** (evaluated before looser operators).

| Precedence (tight ↑) | Operators / forms |
|--------------------|-------------------|
| (loosest) | **`OR`** |
| | **`XOR`** |
| | **`AND`** |
| | **`NOT`** (unary prefix, chains right) |
| | **`=` `<>` `<` `>` `<=` `>=`** |
| | **`+` `-`** (binary) |
| | **`*` `/` `MOD`** |
| | **`^`** (right-associative) |
| | **Unary `-`** |
| | **Postfix** — calls `()`, namespace calls, indexing |
| | **Primary** — literals, identifiers, **`(expr)`** |

**`NOT` vs `OR` / `AND` (permanent rule):** **`NOT` binds tighter than `OR` and `AND`**, matching typical BASIC. Therefore:

- **`NOT a OR b`** parses as **`(NOT a) OR b`**, not **`NOT (a OR b)`**.
- To exit a loop when the key is down **or** the OS requests close, write **`WHILE NOT (Input.KeyDown(KEY_ESCAPE) OR Window.ShouldClose())`** — parentheses are **required**. Without them, **`WHILE NOT Input.KeyDown(KEY_ESCAPE) OR Window.ShouldClose()`** keeps the loop true after a window close (broken).

#### 7.1 Logical operators in bytecode
The parser builds **`OR` / `XOR` / `AND`** as **`ast.BinopNode`**. codegen emits **`OpOr` / `OpXor` / `OpAnd`** as three-register instructions: `Op Dst, SrcA, SrcB`. The VM implements them in the main dispatch loop using `value.Truthy` checks.

---

### 8. VM tracing, `OpCallHandle`, and handle method dispatch

#### 8.1 Golden trace
The virtual machine implements a deterministic state dumper.
- Output format: `[trace] <chunk> L<line> IP=<ip> <opcode> R<dst> R<srcA> R<srcB> | regs=...`
- It provides a snapshot of the register file after each instruction.

#### 8.2 `OpCallHandle` (heap `recv.METHOD(args)`)
Statements and expressions may use **handle calls**: load a handle then call a method. Codegen emits an `OpCallHandle` instruction: `OpCallHandle Dst, RecvReg, ArgStartReg, (ArgCount << 24 | MethodIdx)`. 

**Dispatch rule** (**`vm/handlecall.go`** + **`vm/vm_dispatch.go` `doCallHandle`**): The registry keys for natives are **`CAMERA.BEGIN`**, **`MESH.DRAW`**, etc. Heap objects report a **`TypeName()`** (e.g. **`Camera3D`**, **`Mesh`**, **`Matrix4`**). The VM **does not** call **`Registry.Call("Camera3D.Begin", …)`** — that key is not registered. Instead **`handleCallBuiltin(typeName, method)`** maps to a **registered** key and whether to **prepend the receiver handle** as the first argument:

| Heap `TypeName` | Example method | Registry key | Prepend receiver |
|-----------------|----------------|--------------|------------------|
| **`Camera3D`** | **`BEGIN`**, **`SETPOS`**, **`SETTARGET`**, **`SETFOV`**, **`MOVE`**, **`GETRAY`**, **`GETMATRIX`** | **`CAMERA.<METHOD>`** | yes |
| **`Camera3D`** | **`END`** | **`CAMERA.END`** | no (Raylib **`EndMode3D`** is global) |
| **`Mesh`** | **`DRAW`**, **`DRAWROTATED`** | **`MESH.<METHOD>`** | yes |
| **`Matrix4`** | **`SETROTATION`** | **`TRANSFORM.SETROTATION`** | yes |

Unmapped **`TypeName.METHOD`** combinations should fail at runtime with an unknown-command error until explicitly wired. Method names are matched **case-insensitively** after **`strings.ToUpper`**.

#### 8.3 `ENTITY` spatial macros (`OpEntityPropGet` / `OpEntityPropSet`) and validation

- **Lowering**: **`ENTITY.X(id)`** (and **`Y`**, **`Z`**, **`P`**, **`W`/`YAW`**, **`R`**) in expressions and assignments compile to **`OpEntityPropGet`** / **`OpEntityPropSet`** when arity matches the macro fast path (**`codegen_expr.go`**, **`codegen_stmts.go`**).
- **Shared memory**: When **`Registry.Spatial`** points at the host SoA (**`runtime.SpatialBuffer`**), the VM reads/writes **`float32`** columns directly for props **0–5**.
- **Validation layer**:
  - **Compile time**: known **literal** indices outside **`[0, runtime.MaxEntitySpatialIndex)`** are rejected in **semantic analysis** (so **`--check`** catches them) and again in **codegen** (**`compiler/entityspatial/validate.go`**, **`compiler/semantic/analyze.go`**, **`compiler/codegen/entity_macro_validate.go`**).
  - **Run time**: every access runs **`validateEntityMacroID`** (negative or **≥ `MaxEntitySpatialIndex`** → **`ENTITY:`** runtime error). If **`Registry.EntityIDActive`** is set (**mbentity**), in-bounds SoA slots for **inactive** ids are rejected (avoids silent reads/writes on uninitialized rows). See **`docs/COMPILER_SPEC.md`**.

---

### 9. Phase B: Graphics stack (HAL + Raylib)

- **Canonical doc**: **[`docs/architecture/HAL_AND_RENDERING.md`](docs/architecture/HAL_AND_RENDERING.md)** — HAL package, null vs Raylib drivers, injection points, Windows purego deferred DLL load, and how this differs from static linking.
- **Tests**: `go test ./runtime/window/ -v` with **`CGO_ENABLED=0`** exercises registration paths on Windows without requiring **`raylib.dll`** in typical configurations (vendored purego **`init()`** defers DLL load for **`*.test`** binaries; see **`MOONBASIC_SKIP_RAYLIB_DLL`** in the HAL doc). Full Raylib behavior is still verified with **`CGO_ENABLED=1`** builds (linked Raylib on Linux/macOS; Windows per **`internal/driver`** and CI matrix).
- **Hardware Abstraction Layer (HAL)**: Windowing and the **HAL-backed** draw slice use **`hal.VideoDevice`**, **`hal.InputDevice`**, and **`hal.SystemDevice`** via **`Registry.Driver`**. The **`compiler/pipeline`** package injects either the **Raylib** driver (**`fullruntime`**) or the **null** driver (**`!fullruntime`**) through **`DefaultDriver()`**; tests may use **`runtime.NewRegistryHeadless`**.

#### 9.1 The Driver Architecture (`drivers/` + `hal/`)
- **`hal/`**: Pure Go interfaces and math types — **no** Raylib import.
- **`drivers/video/raylib`**: Implements HAL by delegating to **`raylib-go`** (CGO **or** purego sidecar **`raylib.dll`** on Windows, depending on build flags — not “CGO-only”).
- **`drivers/video/null`**: No-op implementation for headless CLI, **`ListBuiltins`**, and VM/compiler tests that should not touch the GPU.
- **Registry injection**: **`runtime.NewRegistry(h, hal.Driver{...})`** requires an explicit **`hal.Driver`**; **`runtime.NewRegistryHeadless`** supplies the null driver for all three facets.

#### `WINDOW.OPEN`, `WINDOW.CLOSE`, and heap lifecycle
- **`WINDOW.OPEN`**: Manifest has **no return type** (void). After **`InitWindow`**, the native checks **`rl.IsWindowReady()`**; on failure it prints a one-line message to **stderr** and **`os.Exit(1)`**. On success it returns **`value.Nil`**. Scripts that must **probe** without terminating use **`WINDOW.CANOPEN(w,h,title)`** → **`bool`** before calling **`WINDOW.OPEN`**. Use bare **`END`** (**`EndProgramStmt` → `OpHalt`**) to stop cleanly when **`QUIT`** is unavailable.
- **`WINDOW.CLOSE`**: Ends any active frame (**`EndDrawing`**), runs audio close hook, **`CloseWindow`**, then **`rt.Heap.FreeAll()`** so GPU-backed heap objects (meshes, materials, matrices, etc.) are released without requiring explicit **`*.FREE`** in short examples. **`Registry.Shutdown()`** still calls **`Heap.FreeAll()`** at process teardown — **double `FreeAll` is safe** on an already-cleared store.
- **`RENDER.CLEAR` overloads** (single native, arity dispatch; manifest lists multiple rows — **§4**): **`()`** → clear **black** `(0,0,0,255)`; **`(r,g,b)`** → opaque RGB; **`(r,g,b,a)`** → RGBA; **`(colorHandle)`** → resolve via **`mbmatrix.HeapColorRGBA`** (**`runtime/mbmatrix/color_heap.go`**) for heap **`Color`** objects. **`RENDER.CLEAR`** begins a frame (**`BeginDrawing`**) on first use after **`OPEN`** or **`FRAME`**, same as before.
- **High-DPI Support**: **`WINDOW.OPEN`** unconditionally sets **`rl.FlagWindowHighdpi`**. Metrics like **`WINDOW.WIDTH`** / **`HEIGHT`** and **`MOUSEX`** / **`Y`** return **logical** coordinates as **`float64`**. **`RENDER.WIDTH`** / **`HEIGHT`** return the physical **framebuffer** dimensions, and **`WINDOW.DPISCALE`** returns the pixel density ratio.

#### Acceptance test (behavioral “Phase A” on current IR)

Canonical program: [`testdata/pretty_window.mb`](testdata/pretty_window.mb) — opens a 1280×720 window, **`WINDOW.SETFPS(60)`**, clears to RGB **(20, 20, 30)** each frame, **`RENDER.FRAME`**, then **`WINDOW.CLOSE`**. The loop condition must be **`WHILE NOT (Input.KeyDown(KEY_ESCAPE) OR Window.ShouldClose())`**. See **§7** — parentheses are required because **`NOT` binds tighter than `OR`**.

- **Platforms**: verify on **Windows x64** and **Linux x64** with **`CGO_ENABLED=1`**, a C toolchain, and Raylib available to the linker. For a **full windowed run** from the module root use **`go run -tags fullruntime ./cmd/moonrun testdata/pretty_window.mb`** (not `go run ./...`, which is for packages). Plain **`go run . testdata/pretty_window.mb`** only **compiles** to `.mbc` (compiler-only entrypoint). CI: **`.github/workflows/ci.yml`** runs **`go run . --check testdata/pretty_window.mb`** (no window), then **`go test ./...`** with **`CGO_ENABLED=1`** — Linux under **Xvfb** + GL/X11/Wayland packages; Windows via **`msys2/setup-msys2`** (**MINGW64** **`gcc`**) with **`CC`** set to that toolchain (install path comes from the action output, not a fixed drive letter). Integration **`--check`** also includes **`examples/spin_cube/main.mb`** (3D sample; semantic-only on CI).
- **Exit semantics**: ESC or the window **X** ends the loop when the condition above is false; the script then calls **`WINDOW.CLOSE`**. The implementation ends any open drawing frame (**`EndDrawing`**) before **`CloseWindow`**, and avoids double **`CloseWindow`** (undefined behavior in Raylib on some platforms). Process exit status is the CLI’s (0 on normal completion), not a dedicated “ESC code.”
- **“Stable 60fps”**: means **`SetTargetFPS`** (via **`WINDOW.SETFPS`**) sets Raylib’s **target** frame cap, not a hard real-time guarantee. Vsync, GPU load, and the OS scheduler affect measured FPS. **`--info`** prints a one-line **runtime banner** (same libraries as **`--version`**) plus bytecode disassembly before run; with CGO, the window module also prints a throttled **`GetFPS`** line to **`Options.Out`** (~once per second) during **`RENDER.FRAME`** for coarse verification.
- **Resource lifecycle**: **`pipeline.RunProgram`** uses **`defer Registry.Shutdown()`** so Raylib and **`Heap.FreeAll`** run after normal completion or a VM error. **`WINDOW.CLOSE`** additionally frees the heap mid-script when games exit without terminating the process. Treat **valgrind** / **`GODEBUG=gccheckmark`** on Go+CGo binaries as optional signals, not release gates.

#### 3D slice (Phase D precursors, same CGO stack)
Procedural meshes and camera/matrix helpers ship under **`runtime/mbmodel3d`**, **`runtime/mbcamera`**, **`runtime/mbmatrix`** (see **§11** for the full Phase D vision). Current contracts:
- **`MESH.CREATE*`** (canonical) / deprecated **`MESH.MAKE*`** / **`MESH.CUBE`** / **`MESH.SPHERE`** / **`MESH.PLANE`**: After allocating a **`meshObj`**, **`allocMesh`** calls **`rl.UploadMesh(&mesh, false)`** so scripts do not need **`MESH.UPLOAD`** for static procedural geometry. **`MESH.UPLOAD`** remains for **dynamic** meshes or re-upload after edits.
- **`MESH.DRAWROTATED(mesh, mat, rx, ry, rz)`**: Builds an Euler rotation matrix and **`DrawMesh`** (convenience vs **`Transform` + `MESH.DRAW`**).
- **`CAMERA.CREATE`** (canonical; deprecated **`CAMERA.MAKE`**): Initializes a sensible default **3D** camera (position **(0, 2, 8)**, target origin, up **+Y**, **45°** FOV, perspective) in **`runtime/camera/raylib_cgo.go`** — scripts may still call **`SetPos` / `SetTarget` / `SetFOV`** for clarity.
- **`TRANSFORM.ROTATION`** (legacy **`MAT4.ROTATION`** / **`MAT4.FROMROTATION`**) allocates a new matrix. **`TRANSFORM.SETROTATION(handle, rx, ry, rz)`** overwrites the **`Matrix4`** heap object in place (avoids per-frame alloc/free in loops).
- **Canonical small sample**: [`examples/spin_cube/main.mb`](examples/spin_cube/main.mb) — handle-style **`cam.Begin()` / `cam.End()`**, **`Transform.Identity`** + **`Transform.SetRotation`**, no **`Mesh.Upload`**. Larger sample: [`examples/fps/main.mb`](examples/fps/main.mb).

### 10. Phase C runtime modules

New runtime modules follow the same pattern as `runtime/window` and `runtime/input`:

- **Package path:** `runtime/{name}` (e.g. `runtime/mathmod`, `runtime/draw`, `runtime/file`, `runtime/audio`)
- Each package exports **`NewModule() *Module`**
- **`module.go`** contains `Module` struct, `NewModule`, `Register`, `Shutdown` only — thin orchestrator
- CGO packages split into **`raylib_cgo.go`** and **`stub.go`**
- Pure packages need no build-tag split
- **File split convention:** one file per concern, soft limit **~400** lines, split before **~500**
- **HAL compliance (direction of travel)**: Prefer `hal` types (`V3`, `V2`, `RGBA`, `Matrix`) and **`rt.Driver`** for window/video/input surfaces that already sit behind the HAL. Many existing packages still import Raylib directly for meshes, textures, audio, and legacy paths — migrate incrementally; see **[`docs/architecture/HAL_AND_RENDERING.md`](docs/architecture/HAL_AND_RENDERING.md)** §7.
- **Registration order** in `compiler/pipeline/pipeline.go`: all **`RegisterModule`** calls **before** **`RegisterFromManifest`**
- **`runtime/mbgame`**: instant-game / QOL utilities (shortcuts such as **`SCREENW`**, **`DT`**, collision and movement math, easing, noise, **`CONFIG.*`**, wall-clock and sim timers, **`GAME.*` audio helpers, screen flash, etc.). Do **not** register the same dotted keys again from another “QOL” package — the registry maps **one** handler per uppercase key (**§4**).
- **Data modules** (**`runtime/jsonmod`**, **`runtime/csvmod`**, **`runtime/dbmod`**, **`runtime/tablemod`**): add every new dotted name to **`compiler/builtinmanifest/commands.json`** first (overload rows where arity differs). **`jsonmod`** / **`csvmod`** / **`tablemod`** are pure Go; **`dbmod`** is **SQLite via CGO** (`mattn/go-sqlite3`) with **`//go:build !cgo`** stubs that return clear errors when CGO is off. Register **`jsonmod` → `csvmod` → `dbmod` → `tablemod`** before **`RegisterFromManifest`** so integration commands and bridges resolve in order.

### 11. Phase D — 3D engine extension (planned)

Phase D extends the runtime with models, lighting, environment (skybox / IBL / fog), terrain, custom shaders, animation, bones, immediate 3D drawing, shadows, and render-to-texture / post-processing. It is **not** part of the Phase C closure; implement it in ordered milestones (models and debug draw before lighting and shadows).

**Incremental delivery (already in tree):** procedural **`MESH.CREATE*`** / **`MESH.MAKE*`** with automatic **`UploadMesh`**, **`MESH.DRAW` / `MESH.DRAWROTATED`**, **`CAMERA.*`** 3D mode with **`CAMERA.CREATE`** defaults, **`TRANSFORM.*`** (and legacy **`MAT4.*`**) including **`SETROTATION`** and **`ROTATION`**, handle dispatch for **`cam.Begin()`** / **`mesh.Draw(...)`** (**§8.2**), and samples **`examples/spin_cube`**, **`examples/fps`**. These satisfy **early** items in the milestone list below; the acceptance program remains aspirational until lighting, terrain, and shadows are in scope.

- **Authority**: Same rules as §4 and §10 — add each new command to **`compiler/builtinmanifest/commands.json`** first (reuse **`"key"`**; add **overload rows** when the same dotted name needs multiple arities — **§4**). Then implement natives in **`runtime/{name}`** packages with thin **`module.go`**, **`raylib_cgo.go` / `stub.go`** where CGO is required, and **one file per concern** (soft limit ~400 lines, split before ~500).
- **Registration**: New modules are **`RegisterModule`**’d in **`compiler/pipeline/pipeline.go`** in dependency order, all **before** **`RegisterFromManifest`**.
- **Acceptance**: When Phase D is complete, the canonical behavioral reference shall be a **`testdata/`** program (replace the placeholder [`testdata/phase_d_acceptance.mb`](testdata/phase_d_acceptance.mb)) that exercises a 3D scene: loaded or procedural model, terrain interaction, lighting, skybox or gradient, shadow mapping, and a camera that follows terrain height — analogous in role to §9’s **`testdata/pretty_window.mb`** for the window stack. Until then, **`phase_d_acceptance.mb`** remains a minimal **`--check`**-only stub; CI may run **`go run . --check testdata/phase_d_acceptance.mb`** alongside **`pretty_window.mb`**.
- **Suggested milestone order**: (1) model load/draw and primitives — **partially met** (mesh draw + procedural mesh + **`MODEL.LOAD`** surface exists; keep hardening), (2) immediate **`Draw3D.*`** (or equivalent) for debugging, (3) lighting + shadow maps, (4) sky / environment / fog, (5) shader + render-target pipeline and post — **partially met** (**Modern Post-Processing Architecture**), (6) terrain, (7) skeletal animation and bones.

- **Physical Precision**: The post-pipeline operates at the physical resolution returned by **`RENDER.WIDTH/HEIGHT`**, ensuring effects scale correctly on High-DPI displays.

#### Zero-Allocation Render Contract
- **No Per-Frame Allocations**: Native handlers called in the main loop (e.g., `RENDER.FRAME`, `Draw.*`, `SHADER.SET*`) **MUST NOT** allocate on the Go heap.
- **Scratch Buffers**: All CGO calls taking slices (like `SetShaderValue`) must use pre-allocated scratch buffers stored in the module state or global variables.
- **Hook Optimization**: The window module uses `frameHookScratch` to avoid slice allocation when iterating over frame draw hooks.

#### Open-world runtime (Phase D extension — shipped incrementally)

- **Packages** (all **`RegisterModule`** before **`RegisterFromManifest`**, after data modules): [`runtime/terrain`](../runtime/terrain) (**`TERRAIN.*`**, **`CHUNK.*`**) heightfield + chunked **`GenMeshHeightmap`** meshes; [`runtime/worldmgr`](../runtime/worldmgr) (**`WORLD.*`**) streaming center / preload / status; [`runtime/water`](../runtime/water) (**`WATER.*`**) water plane; [`runtime/sky`](../runtime/sky) (**`SKY.*`**) day/night tint sphere; [`runtime/cloudmod`](../runtime/cloudmod) (**`CLOUD.*`**) coverage state (draw hook reserved); [`runtime/weathermod`](../runtime/weathermod) (**`WEATHER.*`**, **`FOG.*`**, **`WIND.*`**); [`runtime/scatter`](../runtime/scatter) (**`SCATTER.*`**, **`PROP.*`**) instanced markers; [`runtime/biome`](../runtime/biome) (**`BIOME.*`**).
- **Threading**: Raylib calls stay on the main OS thread (**§9**); terrain mesh rebuild runs synchronously in **`TERRAIN.DRAW`** / chunk paths on the main thread.
- **Navigation**: terrain-adjacent pathfinding continues to use existing **`NAV.*` / `PATH.*` / `NAVAGENT.*`** in [`runtime/mbnav`](../runtime/mbnav) — add geometry with **`NAV.ADDTERRAIN`** / grid as documented in [docs/reference/NAV_AI.md](docs/reference/NAV_AI.md) and the terrain integration notes in [docs/reference/NAVMESH.md](docs/reference/NAVMESH.md).

**Conceptual overview (how open-world fits together):**

- **Terrain vs world:** **`runtime/terrain`** owns the **heightfield** (CPU height samples) and **chunk meshes** (GPU, built from those samples). **`runtime/worldmgr`** does not store terrain data; it updates a **streaming center** (`WORLD.SETCENTER`) and each frame runs **`WORLD.UPDATE`** so the terrain module loads, unloads, or rebuilds **chunks** near that center. In short: **terrain = data + drawing**, **world = which chunks should exist given player/camera position**.
- **Chunk distances:** **`CHUNK.SETRANGE`** sets **load** vs **unload** distances in world units so nearby chunks stay resident and distant ones can be dropped; two radii avoid **thrashing** (load/unload every frame at a boundary).
- **Typical draw order (inside `CAMERA.BEGIN` / `END`):** **sky** → **opaque terrain** (and opaque props/scatter) → **water** (transparent) → **weather/clouds/particles** last. Exact blending depends on shader and pass setup; subsystem reference pages spell out each API.
- **Handles:** Subsystems return **heap handles**; scripts must call the matching **`*.FREE`** (or rely on shutdown **`Heap.FreeAll`**). Wrong-handle **`CAST`** errors are preferred over crashes (**§10** heap tags).
- **Authority vs roadmap:** Features described in external design docs (async mesh worker pools, per-chunk Jolt heightfields, **`REGION.*` files, lightning callbacks) may be **partial** or **absent**. **`compiler/builtinmanifest/commands.json`** plus **`r.Register("KEY", …)`** in each package are the **source of truth** for what ships today.

### 12. Phase E — Physics (Jolt 3D, Box2D 2D stub, character)

- **Authority**: Same as §4 / §10 — commands are defined in **`compiler/builtinmanifest/commands.json`**; implementations live under **`runtime/physics3d`**, **`runtime/physics2d`**, and **`runtime/charcontroller`** with thin **`module.go`** and a **CGO / stub** split.
- **Dependency**: **[`github.com/bbitechnologies/jolt-go`](https://github.com/bbitechnologies/jolt-go)** (pinned in **`go.mod`**). The binding supports CGO for **Linux**, **Windows**, and **Darwin**. **Windows** builds link against prebuilt static libraries in `third_party/jolt-go/jolt/lib/windows_amd64`. **`PHYSICS2D.*` / `BODY2D.*`** are **Box2D stubs** everywhere (clear runtime error).
- **Build tags**: Jolt-specific files use **`//go:build (linux || windows) && cgo`**; companion **`stub.go`** files use **`//go:build (!linux && !windows) || !cgo`**. Requires **Go 1.25.3+**.
- **Registration**: In **`compiler/pipeline/pipeline.go`**, **`RegisterModule`** for **`charcontroller`**, **`physics2d`**, then **`physics3d`** (all **before** **`RegisterFromManifest`**) so natives override manifest stubs. **`charcontroller`** is registered **before** **`physics3d`** so **`Registry.Shutdown`** tears down **Jolt `CharacterVirtual`** instances **before** the physics world is destroyed. **`physics3d.NewModule().SetUserInvoker(vm.CallUserFunction)`** wires **`PHYSICS3D.PROCESSCOLLISIONS`** to user **`FUNCTION`** callbacks (queued events only; Jolt contact → queue is not fully wired yet).
- **Purity**: **`PHYSICS3D.STEP`**, **`PHYSICS2D.STEP`**, **`BODY3D.*`** mutators, and **`CHARCONTROLLER.MOVE`** are **not** pure; treat them like §9’s render phase for ordering vs **`RENDER.FRAME`**.
- **Heap**: **`BODY3D`** bodies and **`CHARCONTROLLER`** instances are **`HeapObject`** handles; **`BODY3D.FREE`** / **`CHARCONTROLLER.FREE`** (or **`Heap.FreeAll`** on shutdown) release native resources. **`PHYSICS3D.RAYCAST`** returns a **heap numeric array** handle (length 6: hit, normal xyz, fraction, body handle placeholder **0**).
- **Collision callbacks**: **`PHYSICS3D.ONCOLLISION`** registers rules; **`PHYSICS3D.PROCESSCOLLISIONS`** drains the pending queue and invokes the named user function with **`(handle, handle)`**. Callbacks should run **after** **`STEP`** in the script loop. Events are produced when the Jolt listener enqueues them (listener integration is incremental).
- **Acceptance**: Canonical sample: **[`testdata/physics_demo.mb`](testdata/physics_demo.mb)** — static floor, one dynamic box, **`CharController`**, **`Physics3D.Step`** inside the §9 window loop. CI runs **`go run . --check testdata/physics_demo.mb`** for semantics only. A **full** run (native Jolt + Raylib) is optional CI or manual verification on **Linux + CGO**, analogous to the Raylib gate in §9.

### 13. Phase F — Networking (ENet + JSON messages)

- **Authority**: Same as §4 / §10 — canonical keys in **`compiler/builtinmanifest/commands.json`**: **`NET.*`**, **`PEER.*`**, **`EVENT.*`**, **`JSON.*`**. Legacy **`ENET.*`** rows remain manifest stubs with a different shape; prefer **`NET.*`** for new scripts.
- **Packages**: **`runtime/net`** (**`mbnet`**, **`//go:build cgo`** **`enet_cgo.go`** + **`enet_peer_event.go`**, **`stub.go`** for **`!cgo`**); data stacks **`runtime/jsonmod`** (**`mbjson`**), **`runtime/csvmod`** (**`mbcsv`**), **`runtime/tablemod`** (**`mbtable`**) — pure Go; **`runtime/dbmod`** (**`mbdb`**) — SQLite when CGO enabled.
- **Dependency**: **[`github.com/codecat/go-enet`](https://github.com/codecat/go-enet)** (CGO). **Linux**: system **`libenet`** (**`libenet-dev`** on Debian/Ubuntu — see CI). **Windows**: the module vendors ENet sources; use the same **MinGW** toolchain as Raylib when **`CGO_ENABLED=1`**. Call **`NET.START`** before creating hosts; **`NET.STOP`** (or registry **`Shutdown`**) destroys open hosts and **`enet_deinitialize`**.
- **Event model**: **`NET.UPDATE`** pumps **`enet_host_service`** into an internal per-host queue; **`NET.RECEIVE`** pops one queued item into a heap **event** object (handle **`0`** means no event). Types: **1** connect, **2** disconnect, **3** receive (matches ENet **`EventConnect` / `EventDisconnect` / `EventReceive`**). **`EVENT.DATA`** is UTF-8 text; binary payloads should use **base64** (or another encoding) in the game layer.
- **Peer handles**: Stable heap **`NetPeer`** objects; ENet **`Peer.SetData`** stores the heap id so **`EVENT.PEER`** matches **`NET.CONNECT`** and server connect events. **`NET.SETBANDWIDTH`** and **`NET.SETTIMEOUT`** are reserved no-ops until a lower-level wrapper exposes **`enet_host_bandwidth_limit`** / **`enet_peer_timeout`** (bandwidth is set at **`NewHost`** creation today).
- **Acceptance**: **[`testdata/net_server.mb`](testdata/net_server.mb)** and **[`testdata/net_client.mb`](testdata/net_client.mb)** — JSON **`hello` / `ack`** exchange on port **27777**; CI runs **`--check`** only. Full client/server behavior requires two processes (manual or custom harness).

### 14. Phase H — Developer experience (CLI + LSP + editor)

- **Diagnostics**: Unknown **`NS.METHOD`** engine commands use **`compiler/builtinmanifest`** helpers for **did-you-mean** (edit distance) and **Available:** listings for the namespace (see **`compiler/semantic/cmdhint.go`**).
- **CLI** (**`main.go`**): **`--disasm <file.mbc>`** — human-readable bytecode via **`compiler/pipeline.PrintProgramDisassembly`** (optional same-stem **`.mb`** for source-line annotations). **`--profile <source.mb>`** — per–source-line instruction counts via **`vm.ProfileRecorder`**; prints top 10 after run. **`--watch <source.mb>`** — **`fsnotify`** debounced recompile + **`RunProgram`**. **`--lsp`** — stdio LSP in **`lsp/`** (hover for **`NS.METHOD`** using **`builtinmanifest.FirstOverload`** when arity is not parsed from the line — **§4**; completion after **`.`**, diagnostics from **`pipeline.CheckSource`** with **overload-aware arity** checking).
- **VS Code**: Extension under **`editors/vscode-moonbasic`** — TextMate grammar for **`.mb`** and **`.mbc`**, snippets (**WHILE/WEND**, **FOR/NEXT**, **FUNCTION/ENDFUNCTION**, **SELECT/CASE**), Language Client spawning **`moonbasic --lsp`** (override with **`moonbasic.languageServerPath`**). Run **`npm install`** and **`npm run compile`** in that folder before **F5** / packaging (**`vscode-languageclient`** lives in **`node_modules`**; **`out/extension.js`** is emitted by **`tsc`**).
- **gopls / build tags**: **`.vscode/settings.json`** sets **`fullruntime,gopls_stub`** plus **`CGO_ENABLED=1`** so Raylib AND Jolt **`*_cgo.go`** files analyze on Windows. **`gopls_stub`** switches **`runtime/terrain`** and Jolt components to their **stub** sources for the language server when required (see **`//go:build`** on **`heap_objects_stub.go`** vs **`heap_objects_raylib.go`**); **`go build -tags fullruntime`** does **not** use **`gopls_stub`**. **`main.go`** / **`cmd/moonbasic`** use **`!fullruntime`** — you may see **“no packages found”** for those until you temporarily adjust gopls tags (**[`docs/DEVELOPER.md`](docs/DEVELOPER.md)**).

### 15. Procedural noise (`NOISE.*`, `runtime/procnoise`, `runtime/noisemod`)

- **Authority**: **`compiler/builtinmanifest/commands.json`** lists every **`NOISE.*`** key; implementations live in **`runtime/noisemod`** with sampling core in **`runtime/procnoise`** (pure Go, shared with legacy **`PERLIN` / `SIMPLEX` / `VORONOI` / `FBMNOISE`** in **`runtime/mbgame`** — same algorithms at seed **0** for backward compatibility).
- **Handles**: **`NoiseObject`** uses **`heap.TagNoise`**. **`NOISE.FREE`** or **`Heap.FreeAll`** tear down generator state (no Raylib objects unless **`NOISE.FILLIMAGE`** touches an **`Image`** handle).
- **`NOISE.FILLIMAGE`**: **`//go:build cgo`** — uses **`mbimage.RayImageForTexture`** and **`rl.ImageDrawPixel`**; **`//go:build !cgo`** returns a clear error.
- **Configuration lock**: After the first **`NOISE.GET`** / **`FILL*`**, **`Set*`** returns an error (immutable configuration), matching the intended FastNoise-style workflow.
- **Namespace shadowing**: Identifiers are case-agnostic — a variable named **`noise`** uppercases to **`NOISE`** and **shadows** the **`Noise.*`** namespace; use names like **`ng`**, **`gen`**, **`terrainNoise`** for handles (see **`docs/reference/NOISE.md`**).
- **Docs / samples**: Reference **`docs/reference/NOISE.md`**; **`testdata/noise_test.mb`**, **`testdata/noise_terrain.mb`**.
