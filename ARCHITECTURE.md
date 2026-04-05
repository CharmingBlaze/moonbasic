# moonBASIC Architecture (v1.2.1)
## Mandatory Document for AI Assistants

This document defines the **Ground Truth** for the moonBASIC compiler and runtime. Any changes must adhere to the modular structure and stable APIs defined here. **DO NOT REVERT TO OLD MONOLITHIC VERSIONS.**

### Note for AI assistants (Cursor, etc.)

- Treat this file as **authoritative** over chat memory or older summaries. If the repo already matches this document, **do not** “restore” or replace `main.go` / `pipeline.go` with older orchestration patterns.
- **`compiler/pipeline`** intentionally imports **`runtime`** here so **`RunProgram`** is a one-call embedder entrypoint. Do not delete `EncodeMOON` / `DecodeMOON` or stub out `--compile` / `--run` in `main.go` if this section lists them as shipped.
- **`CallStmtNode`** must delegate to **`emitCallStmt`** in **`codegen_stmts.go`**; an empty `case` is a **bug** (symptoms: `PRINT` produces no bytecode).
- **`commands.json`** supports **multiple rows per `key`** (overload arities). Semantic analysis uses **`LookupArity`**; do not assume **`Table.Commands[key]`** is a single struct.
- Logical **`AND` / `OR` / `XOR`** must compile to **`OpAnd` / `OpOr` / `OpXor`** (**§5**, **§7.1**). Handle calls **`recv.METHOD(...)`** require **`handleCallBuiltin`** mapping (**§8.2**); **`Camera3D.Begin`** is **not** a registry key.

---

### 1. The First Law: Case Agnosticism
- **Rule**: Every keyword and identifier in moonBASIC is case-agnostic.
- **Implementation**: The Lexer and Symbol Table **MUST** unconditionally normalize all tokens to **UPPERCASE**. 
- **Validation**: Any `strings.ToLower` in the pipeline is a bug unless it is inside a string literal.

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

---

### 3. MOON bytecode (`.mbc`)

- **Package**: `vm/moon` — binary schema (not `encoding/gob` for shipping).
- **Header** (16 bytes): magic `MOON`, big-endian version (**`0x00020000`** for IR v2; MOON `0x00010000` is rejected), reserved flags, entry offset.
- **IR v2 payload**: program-level string table, then chunks; 8-byte instructions; see **`ENGINE_IR_V2.md`**.
- **Loader** validates header before building VM state so launchers can reject wrong engines quickly.

---

### 4. The Registry Manifest (`commands.json`)
`compiler/builtinmanifest/commands.json` is the **Single Source of Truth** for all built-in commands.

- **Overloads**: The same canonical **`"key"`** (e.g. **`RENDER.CLEAR`**) may appear **multiple times** with different **`"args"`** arrays. The JSON loader (`compiler/builtinmanifest/manifest_json.go`) builds **`Table.Commands`** as **`map[string][]Command`**: one map entry per dotted key, value = ordered overload list.
- **Semantic pass** (`compiler/semantic/analyze.go`): Resolves **`NamespaceCallStmt`** and **`NamespaceCallExpr`** with **`LookupArity(ns, method, len(args))`**. If the namespace method exists but no overload matches the argument count → **Compile Error** with **`ArityHint`** (lists valid arities). Unknown **`NS.METHOD`** → **Compile Error** with did-you-mean (**`compiler/semantic/cmdhint.go`**).
- **Runtime dispatch** (`runtime/runtime.go`): **`RegisterFromManifest`** walks every overload but registers **at most one stub per canonical `Command.Key`** — the first-seen **`Namespace`** wins for the stub closure. Natives that support multiple arities implement **branching on `len(args)`** in one **`BuiltinFn`** (e.g. **`RENDER.CLEAR`** in **`runtime/window/raylib_cgo.go`**).
- **LSP** (`lsp/server.go`): Hover for a dotted builtin uses **`FirstOverload(key)`** when arity cannot be inferred from the line; multi-arity commands may show only the first signature unless tooling is extended.

**Inventory tooling**: From the repo root, **`python tools/gen_master_audit.py`** regenerates **`MASTER_AUDIT.txt`** (manifest keys vs **`Register("KEY"`** in **`runtime/**/*.go`**, excluding **`*_test.go`**) plus **`MASTER_AUDIT_REGISTERED.txt`**, **`MASTER_AUDIT_MANIFEST.txt`**, and **`MASTER_AUDIT_DUPLICATES.txt`**.

---

### 5. Modular Code Generation
Codegen is split into specialized sub-handlers to maintain the ~400 line limit:
- `codegen.go`: Structural base and orchestration.
- `codegen_expr.go`: Handles literal and expression emission, including binary **`AND` / `OR` / `XOR`** on **`ast.BinopNode`** → **`OpAnd` / `OpOr` / `OpXor`** (see §7.1).
- `codegen_stmts.go`: Handles statements, control flow (`IF`, `WHILE`, `FOR`, `REPEAT`), and `OpPop` stack hygiene.
- `codegen_calls.go`: Handles built-in and user-function call resolution.

---

### 6. Stack Hygiene
- **Rule**: moonBASIC statements **MUST NOT** leave values on the operand stack.
- **Implementation**: `AssignNode` and `CallStmtNode` (Expression statements) must emit an `OpPop` after the store/call.

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
The parser builds **`OR` / `XOR` / `AND`** as **`ast.BinopNode`** (**`compiler/parser/parser_expr.go`**). **`codegen_expr.go`** must emit **`OpOr` / `OpXor` / `OpAnd`** (not only comparisons and arithmetic). The VM implements them in **`doLogic`** (**`vm/vm_dispatch.go`**) using **`value.Truthy`** on both operands. **`NOT`** remains **`OpNot`**. Regression: **`compiler/codegen`** tests should assert **`OR` / `AND` / `XOR`** appear in disassembly for minimal programs.

---

### 8. VM tracing, `OpCallHandle`, and handle method dispatch

#### 8.1 Golden trace
The virtual machine implements a deterministic state dumper triggered by the **`--trace`** flag.
- Output format: `[trace] <chunk> L<line> IP=<ip> <opcode> | depth=N stack=...` (see **`vm/vm.go`**)
- Use for compiler/VM regression tests where a fixed trace is required.

#### 8.2 `OpCallHandle` (heap `recv.METHOD(args)`)
Statements and expressions may use **handle calls**: load a handle (**`cam`**, **`mesh`**, matrix handle), then call a method (**`cam.Begin()`**, **`mesh.Draw(mat, rot)`**). The parser produces **`HandleCallStmt` / `HandleCallExpr`** (**`compiler/parser`**); codegen emits **`OpLoadLocal` / `OpLoadGlobal`**, pushes args, then **`OpCallHandle`** with the method name in **`Chunk.Names`** (**`compiler/codegen`**).

**Dispatch rule** (**`vm/handlecall.go`** + **`vm/vm_dispatch.go` `doCallHandle`**): The registry keys for natives are **`CAMERA.BEGIN`**, **`MESH.DRAW`**, etc. Heap objects report a **`TypeName()`** (e.g. **`Camera3D`**, **`Mesh`**, **`Matrix4`**). The VM **does not** call **`Registry.Call("Camera3D.Begin", …)`** — that key is not registered. Instead **`handleCallBuiltin(typeName, method)`** maps to a **registered** key and whether to **prepend the receiver handle** as the first argument:

| Heap `TypeName` | Example method | Registry key | Prepend receiver |
|-----------------|----------------|--------------|------------------|
| **`Camera3D`** | **`BEGIN`**, **`SETPOS`**, **`SETTARGET`**, **`SETFOV`**, **`MOVE`**, **`GETRAY`**, **`GETMATRIX`** | **`CAMERA.<METHOD>`** | yes |
| **`Camera3D`** | **`END`** | **`CAMERA.END`** | no (Raylib **`EndMode3D`** is global) |
| **`Mesh`** | **`DRAW`**, **`DRAWROTATED`** | **`MESH.<METHOD>`** | yes |
| **`Matrix4`** | **`SETROTATION`** | **`TRANSFORM.SETROTATION`** | yes |

Unmapped **`TypeName.METHOD`** combinations should fail at runtime with an unknown-command error until explicitly wired. Method names are matched **case-insensitively** after **`strings.ToUpper`**.

---

### 9. Phase B: Raylib window (CGO)

- **Packages** [`runtime/window`](c:\Users\rain\Documents\GO\moonbasic\runtime\window) (**`WINDOW.*`**, **`RENDER.CLEAR` / `FRAME`**) and [`runtime/input`](c:\Users\rain\Documents\GO\moonbasic\runtime\input) (**`INPUT.KEYDOWN` / `KEYPRESSED` / `KEYUP`** when CGO is on). Stubs when **`CGO_ENABLED=0`** (input keys never down).
- **`runtime.SeedInputKeyGlobals`**: preloads **`KEY_ESCAPE`**, **`KEY_W`**, **`KEY_A`**, **`KEY_S`**, **`KEY_D`**, **`KEY_SPACE`** on the VM global map so scripts can use **`Input.KeyDown(KEY_ESCAPE)`** (values match raylib `KeyboardKey`).
- **Build tags**: `raylib_cgo.go` uses `//go:build cgo`; `stub.go` uses `//go:build !cgo`. With **`CGO_ENABLED=0`**, builtins return a clear error directing users to enable CGO and install a C compiler (e.g. MinGW on Windows, gcc on Linux).
- **`pipeline.RunProgram`** calls **`runtime.LockOSThread()`** before init (Raylib requires the main OS thread), then **`RegisterModule(window.NewModule())`** **before** **`RegisterFromManifest`** so real handlers win over manifest stubs.
- **Dependency**: [`github.com/gen2brain/raylib-go/raylib`](https://github.com/gen2brain/raylib-go) in `go.mod` — pinned to the **raylib 5.5** line (`v0.56.0-dev` family; see exact pseudo-version in `go.mod`).
- **Tests**: `go test ./runtime/window/ -v` with **`CGO_ENABLED=0`** runs stub registration tests; full Raylib link is verified with **`CGO_ENABLED=1 go build`** on a machine with Raylib dev libraries available to the C linker.

#### `WINDOW.OPEN`, `WINDOW.CLOSE`, and heap lifecycle
- **`WINDOW.OPEN`**: Manifest **`returns: "bool"`**. After **`InitWindow`**, the native checks **`rl.IsWindowReady()`**; on failure it closes the window handle, leaves **`opened`** false, and returns **`FALSE`**. On success returns **`TRUE`**. Scripts can branch with **`IF NOT Window.Open(…)`**; use bare **`END`** (**`EndProgramStmt` → `OpHalt`**) to stop cleanly when **`QUIT`** is unavailable.
- **`WINDOW.CLOSE`**: Ends any active frame (**`EndDrawing`**), runs audio close hook, **`CloseWindow`**, then **`rt.Heap.FreeAll()`** so GPU-backed heap objects (meshes, materials, matrices, etc.) are released without requiring explicit **`*.FREE`** in short examples. **`Registry.Shutdown()`** still calls **`Heap.FreeAll()`** at process teardown — **double `FreeAll` is safe** on an already-cleared store.
- **`RENDER.CLEAR` overloads** (single native, arity dispatch; manifest lists multiple rows — **§4**): **`()`** → clear **black** `(0,0,0,255)`; **`(r,g,b)`** → opaque RGB; **`(r,g,b,a)`** → RGBA; **`(colorHandle)`** → resolve via **`mbmatrix.HeapColorRGBA`** (**`runtime/mbmatrix/color_heap.go`**) for heap **`Color`** objects. **`RENDER.CLEAR`** begins a frame (**`BeginDrawing`**) on first use after **`OPEN`** or **`FRAME`**, same as before.

#### Acceptance test (behavioral “Phase A” on current IR)

Canonical program: [`testdata/pretty_window.mb`](testdata/pretty_window.mb) — opens a 1280×720 window, **`WINDOW.SETFPS(60)`**, clears to RGB **(20, 20, 30)** each frame, **`RENDER.FRAME`**, then **`WINDOW.CLOSE`**. The loop condition must be **`WHILE NOT (Input.KeyDown(KEY_ESCAPE) OR Window.ShouldClose())`**. See **§7** — parentheses are required because **`NOT` binds tighter than `OR`**.

- **Platforms**: verify on **Windows x64** and **Linux x64** with **`CGO_ENABLED=1`**, a C toolchain, and Raylib available to the linker. From the module root use **`go run . testdata/pretty_window.mb`** (not `go run ./...`, which is for packages). CI: **`.github/workflows/ci.yml`** runs **`go run . --check testdata/pretty_window.mb`** (no window), then **`go test ./...`** with **`CGO_ENABLED=1`** — Linux under **Xvfb** + GL/X11/Wayland packages; Windows via **`msys2/setup-msys2`** (**MINGW64** **`gcc`**) with **`CC`** set to that toolchain (install path comes from the action output, not a fixed drive letter). Integration **`--check`** also includes **`examples/spin_cube/main.mb`** (3D sample; semantic-only on CI).
- **Exit semantics**: ESC or the window **X** ends the loop when the condition above is false; the script then calls **`WINDOW.CLOSE`**. The implementation ends any open drawing frame (**`EndDrawing`**) before **`CloseWindow`**, and avoids double **`CloseWindow`** (undefined behavior in Raylib on some platforms). Process exit status is the CLI’s (0 on normal completion), not a dedicated “ESC code.”
- **“Stable 60fps”**: means **`SetTargetFPS`** (via **`WINDOW.SETFPS`**) sets Raylib’s **target** frame cap, not a hard real-time guarantee. Vsync, GPU load, and the OS scheduler affect measured FPS. **`--info`** prints a one-line **runtime banner** (same libraries as **`--version`**) plus bytecode disassembly before run; with CGO, the window module also prints a throttled **`GetFPS`** line to **`Options.Out`** (~once per second) during **`RENDER.FRAME`** for coarse verification.
- **Resource lifecycle**: **`pipeline.RunProgram`** uses **`defer Registry.Shutdown()`** so Raylib and **`Heap.FreeAll`** run after normal completion or a VM error. **`WINDOW.CLOSE`** additionally frees the heap mid-script when games exit without terminating the process. Treat **valgrind** / **`GODEBUG=gccheckmark`** on Go+CGo binaries as optional signals, not release gates.

#### 3D slice (Phase D precursors, same CGO stack)
Procedural meshes and camera/matrix helpers ship under **`runtime/mbmodel3d`**, **`runtime/mbcamera`**, **`runtime/mbmatrix`** (see **§11** for the full Phase D vision). Current contracts:
- **`MESH.MAKE*`** / **`MESH.CUBE`** / **`MESH.SPHERE`** / **`MESH.PLANE`**: After allocating a **`meshObj`**, **`allocMesh`** calls **`rl.UploadMesh(&mesh, false)`** so scripts do not need **`MESH.UPLOAD`** for static procedural geometry. **`MESH.UPLOAD`** remains for **dynamic** meshes or re-upload after edits.
- **`MESH.DRAWROTATED(mesh, mat, rx, ry, rz)`**: Builds an Euler rotation matrix and **`DrawMesh`** (convenience vs **`Transform` + `MESH.DRAW`**).
- **`CAMERA.MAKE`**: Initializes a sensible default **3D** camera (position **(0, 2, 8)**, target origin, up **+Y**, **45°** FOV, perspective) in **`runtime/camera/raylib_cgo.go`** — scripts may still call **`SetPos` / `SetTarget` / `SetFOV`** for clarity.
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
- **Registration order** in `compiler/pipeline/pipeline.go`: all **`RegisterModule`** calls **before** **`RegisterFromManifest`**
- **`runtime/mbgame`**: instant-game / QOL utilities (shortcuts such as **`SCREENW`**, **`DT`**, collision and movement math, easing, noise, **`CONFIG.*`**, wall-clock and sim timers, **`GAME.*` audio helpers, screen flash, etc.). Do **not** register the same dotted keys again from another “QOL” package — the registry maps **one** handler per uppercase key (**§4**).
- **Data modules** (**`runtime/jsonmod`**, **`runtime/csvmod`**, **`runtime/dbmod`**, **`runtime/tablemod`**): add every new dotted name to **`compiler/builtinmanifest/commands.json`** first (overload rows where arity differs). **`jsonmod`** / **`csvmod`** / **`tablemod`** are pure Go; **`dbmod`** is **SQLite via CGO** (`mattn/go-sqlite3`) with **`//go:build !cgo`** stubs that return clear errors when CGO is off. Register **`jsonmod` → `csvmod` → `dbmod` → `tablemod`** before **`RegisterFromManifest`** so integration commands and bridges resolve in order.

### 11. Phase D — 3D engine extension (planned)

Phase D extends the runtime with models, lighting, environment (skybox / IBL / fog), terrain, custom shaders, animation, bones, immediate 3D drawing, shadows, and render-to-texture / post-processing. It is **not** part of the Phase C closure; implement it in ordered milestones (models and debug draw before lighting and shadows).

**Incremental delivery (already in tree):** procedural **`MESH.MAKE*`** with automatic **`UploadMesh`**, **`MESH.DRAW` / `MESH.DRAWROTATED`**, **`CAMERA.*`** 3D mode with **`CAMERA.MAKE`** defaults, **`TRANSFORM.*`** (and legacy **`MAT4.*`**) including **`SETROTATION`** and **`ROTATION`**, handle dispatch for **`cam.Begin()`** / **`mesh.Draw(...)`** (**§8.2**), and samples **`examples/spin_cube`**, **`examples/fps`**. These satisfy **early** items in the milestone list below; the acceptance program remains aspirational until lighting, terrain, and shadows are in scope.

- **Authority**: Same rules as §4 and §10 — add each new command to **`compiler/builtinmanifest/commands.json`** first (reuse **`"key"`**; add **overload rows** when the same dotted name needs multiple arities — **§4**). Then implement natives in **`runtime/{name}`** packages with thin **`module.go`**, **`raylib_cgo.go` / `stub.go`** where CGO is required, and **one file per concern** (soft limit ~400 lines, split before ~500).
- **Registration**: New modules are **`RegisterModule`**’d in **`compiler/pipeline/pipeline.go`** in dependency order, all **before** **`RegisterFromManifest`**.
- **Acceptance**: When Phase D is complete, the canonical behavioral reference shall be a **`testdata/`** program (replace the placeholder [`testdata/phase_d_acceptance.mb`](testdata/phase_d_acceptance.mb)) that exercises a 3D scene: loaded or procedural model, terrain interaction, lighting, skybox or gradient, shadow mapping, and a camera that follows terrain height — analogous in role to §9’s **`testdata/pretty_window.mb`** for the window stack. Until then, **`phase_d_acceptance.mb`** remains a minimal **`--check`**-only stub; CI may run **`go run . --check testdata/phase_d_acceptance.mb`** alongside **`pretty_window.mb`**.
- **Suggested milestone order**: (1) model load/draw and primitives — **partially met** (mesh draw + procedural mesh + **`MODEL.LOAD`** surface exists; keep hardening), (2) immediate **`Draw3D.*`** (or equivalent) for debugging, (3) lighting + shadow maps, (4) sky / environment / fog, (5) shader + render-target pipeline and post, (6) terrain, (7) skeletal animation and bones.

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
- **Dependency**: **[`github.com/bbitechnologies/jolt-go`](https://github.com/bbitechnologies/jolt-go)** (pinned in **`go.mod`**). The binding currently ships CGO for **Linux** (**amd64**, **arm64**) and **Darwin arm64** only. **Windows** builds use **stubs** for **`PHYSICS3D.*`**, **`BODY3D.*`**, and **`CHARCONTROLLER.*`** until a Windows-capable binding exists. **`PHYSICS2D.*` / `BODY2D.*`** are **Box2D stubs** everywhere (clear runtime error).
- **Build tags**: **`runtime/physics3d/jolt_linux.go`** and **`runtime/charcontroller/jolt_linux.go`** use **`//go:build linux && cgo`**; companion **`stub.go`** files use **`//go:build !linux || !cgo`**. Requires **Go 1.25.3+** (jolt-go requirement).
- **Registration**: In **`compiler/pipeline/pipeline.go`**, **`RegisterModule`** for **`charcontroller`**, **`physics2d`**, then **`physics3d`** (all **before** **`RegisterFromManifest`**) so natives override manifest stubs. **`charcontroller`** is registered **before** **`physics3d`** so **`Registry.Shutdown`** tears down **Jolt `CharacterVirtual`** instances **before** the physics world is destroyed. **`physics3d.NewModule().SetUserInvoker(vm.CallUserFunction)`** wires **`PHYSICS3D.PROCESSCOLLISIONS`** to user **`FUNCTION`** callbacks (queued events only; Jolt contact → queue is not fully wired yet).
- **Purity**: **`PHYSICS3D.STEP`**, **`PHYSICS2D.STEP`**, **`BODY3D.*`** mutators, and **`CHARCONTROLLER.MOVE`** are **not** pure; treat them like §9’s render phase for ordering vs **`RENDER.FRAME`**.
- **Heap**: **`BODY3D`** bodies and **`CHARCONTROLLER`** instances are **`HeapObject`** handles; **`BODY3D.FREE`** / **`CHARCONTROLLER.FREE`** (or **`Heap.FreeAll`** on shutdown) release native resources. **`PHYSICS3D.RAYCAST`** returns a **heap numeric array** handle (length 6: hit, normal xyz, fraction, body handle placeholder **0**).
- **Collision callbacks**: **`PHYSICS3D.ONCOLLISION`** registers rules; **`PHYSICS3D.PROCESSCOLLISIONS`** drains the pending queue and invokes the named user function with **`(handle, handle)`**. Callbacks should run **after** **`STEP`** in the script loop. Events are produced when the Jolt listener enqueues them (listener integration is incremental).
- **Acceptance**: Canonical sample: **[`testdata/physics_demo.mb`](testdata/physics_demo.mb)** — static floor, one dynamic box, **`CharController`**, **`Physics3D.Step`** inside the §9 window loop. CI runs **`go run . --check testdata/physics_demo.mb`** for semantics only. A **full** run (native Jolt + Raylib) is optional CI or manual verification on **Linux + CGO**, analogous to the Raylib gate in §9.

### 13. Phase F — Networking (ENet + JSON messages)

- **Authority**: Same as §4 / §10 — canonical keys in **`compiler/builtinmanifest/commands.json`**: **`NET.*`**, **`PEER.*`**, **`EVENT.*`**, **`JSON.*`**. Legacy **`ENET.*`** rows remain manifest stubs with a different shape; prefer **`NET.*`** for new scripts.
- **Packages**: **`runtime/net`** (**`mbnet`**, **`//go:build cgo`** **`enet_cgo.go`** + **`enet_peer_event.go`**, **`stub.go`** for **`!cgo`**); data stacks **`runtime/jsonmod`** (**`mbjson`**), **`runtime/csvmod`** (**`mbcsv`**), **`runtime/tablemod`** (**`mbtable`**) — pure Go; **`runtime/dbmod`** (**`mbdb`**) — SQLite when CGO enabled.
- **Dependency**: **[`github.com/codecat/go-enet`](https://github.com/codecat/go-enet)** (CGO). **Linux**: system **`libenet`** (**`libenet-dev`** on Debian/Ubuntu — see CI). **Windows**: the module vendors ENet sources; use the same **MinGW** toolchain as Raylib when **`CGO_ENABLED=1`**. Call **`NET.START`** before creating hosts; **`NET.STOP`** (or registry **`Shutdown`**) destroys open hosts and **`enet_deinitialize`**.
- **Event model**: **`NET.UPDATE`** pumps **`enet_host_service`** into an internal per-host queue; **`NET.RECEIVE`** pops one queued item into a heap **event** object (handle **`0`** means no event). Types: **1** connect, **2** disconnect, **3** receive (matches ENet **`EventConnect` / `EventDisconnect` / `EventReceive`**). **`EVENT.DATA$`** is UTF-8 text; binary payloads should use **base64** (or another encoding) in the game layer.
- **Peer handles**: Stable heap **`NetPeer`** objects; ENet **`Peer.SetData`** stores the heap id so **`EVENT.PEER`** matches **`NET.CONNECT`** and server connect events. **`NET.SETBANDWIDTH`** and **`NET.SETTIMEOUT`** are reserved no-ops until a lower-level wrapper exposes **`enet_host_bandwidth_limit`** / **`enet_peer_timeout`** (bandwidth is set at **`NewHost`** creation today).
- **Acceptance**: **[`testdata/net_server.mb`](testdata/net_server.mb)** and **[`testdata/net_client.mb`](testdata/net_client.mb)** — JSON **`hello` / `ack`** exchange on port **27777**; CI runs **`--check`** only. Full client/server behavior requires two processes (manual or custom harness).

### 14. Phase H — Developer experience (CLI + LSP + editor)

- **Diagnostics**: Unknown **`NS.METHOD`** engine commands use **`compiler/builtinmanifest`** helpers for **did-you-mean** (edit distance) and **Available:** listings for the namespace (see **`compiler/semantic/cmdhint.go`**).
- **CLI** (**`main.go`**): **`--disasm <file.mbc>`** — human-readable bytecode via **`compiler/pipeline.PrintProgramDisassembly`** (optional same-stem **`.mb`** for source-line annotations). **`--profile <source.mb>`** — per–source-line instruction counts via **`vm.ProfileRecorder`**; prints top 10 after run. **`--watch <source.mb>`** — **`fsnotify`** debounced recompile + **`RunProgram`**. **`--lsp`** — stdio LSP in **`lsp/`** (hover for **`NS.METHOD`** using **`builtinmanifest.FirstOverload`** when arity is not parsed from the line — **§4**; completion after **`.`**, diagnostics from **`pipeline.CheckSource`** with **overload-aware arity** checking).
- **VS Code / Cursor**: Extension under **`editors/vscode-moonbasic`** — TextMate grammar for **`.mb`** and **`.mbc`**, snippets (**WHILE/WEND**, **FOR/NEXT**, **FUNCTION/ENDFUNCTION**, **SELECT/CASE**), Language Client spawning **`moonbasic --lsp`** (override with **`moonbasic.languageServerPath`**). Run **`npm install`** and **`npm run compile`** in that folder before **F5** / packaging (**`vscode-languageclient`** lives in **`node_modules`**; **`out/extension.js`** is emitted by **`tsc`**).
- **gopls / build tags**: The repo includes **`.vscode/settings.json`** setting **`CGO_ENABLED=1`** for **`go`** / **`gopls`** so Raylib and ENet **`*_cgo.go`** files are part of the language-server build (the default gopls env often had **`CGO_ENABLED=0`**, which excluded them and triggered **“no packages found”**). CGO-backed files use **`cgo && !gopls_stub`**; stubs use **`!cgo || gopls_stub`** (Jolt: **`linux && cgo && !gopls_stub`** vs **`!linux || !cgo || gopls_stub`**). Normal **`go build`** is unchanged. **`jolt_linux.go`** still matches only **`GOOS=linux`**; on Windows/macOS, gopls may keep showing that hint for those files—use WSL/Linux, or ignore. If gopls shows **“no packages found”** for a **`stub.go`** while **`CGO_ENABLED=1`**, set **`"gopls": { "build.buildFlags": ["-tags=gopls_stub"] }`** so stubs join the analysis build (CGO sources drop out until you remove the flag).

### 15. Procedural noise (`NOISE.*`, `runtime/procnoise`, `runtime/noisemod`)

- **Authority**: **`compiler/builtinmanifest/commands.json`** lists every **`NOISE.*`** key; implementations live in **`runtime/noisemod`** with sampling core in **`runtime/procnoise`** (pure Go, shared with legacy **`PERLIN` / `SIMPLEX` / `VORONOI` / `FBMNOISE`** in **`runtime/mbgame`** — same algorithms at seed **0** for backward compatibility).
- **Handles**: **`NoiseObject`** uses **`heap.TagNoise`**. **`NOISE.FREE`** or **`Heap.FreeAll`** tear down generator state (no Raylib objects unless **`NOISE.FILLIMAGE`** touches an **`Image`** handle).
- **`NOISE.FILLIMAGE`**: **`//go:build cgo`** — uses **`mbimage.RayImageForTexture`** and **`rl.ImageDrawPixel`**; **`//go:build !cgo`** returns a clear error.
- **Configuration lock**: After the first **`NOISE.GET`** / **`FILL*`**, **`Set*`** returns an error (immutable configuration), matching the intended FastNoise-style workflow.
- **Namespace shadowing**: Identifiers are case-agnostic — a variable named **`noise`** uppercases to **`NOISE`** and **shadows** the **`Noise.*`** namespace; use names like **`ng`**, **`gen`**, **`terrainNoise`** for handles (see **`docs/reference/NOISE.md`**).
- **Docs / samples**: Reference **`docs/reference/NOISE.md`**; **`testdata/noise_test.mb`**, **`testdata/noise_terrain.mb`**.
