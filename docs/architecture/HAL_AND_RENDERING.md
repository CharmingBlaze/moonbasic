# HAL and rendering architecture

This document is the **canonical** description of how moonBASIC separates **graphics/window concerns** from the **compiler and VM**, why that matters for CI and headless tooling, and how **linking** (CGO vs purego vs static) relates to that separation.

**Related:** [ZERO_CGO_RAYLIB.md](ZERO_CGO_RAYLIB.md) (purego/sidecar strategy and symbol inventory), [BUILDING.md](../BUILDING.md) (toolchains), [scripts/build_static.ps1](../../scripts/build_static.ps1) (experimental Zig CC build).

---

## 1. Why decouple at all?

- **Compiler and bytecode** should not need a GPU, a display server, or a **`raylib.dll`** on `PATH` to run **`go test`**. If importing a runtime package panics in **`init()`** because a shared library is missing, tests and small tools inherit that fragility.
- **Headless tooling** (asset checks, **`--check`**, LSP, future packers) should stay **fast and portable**.
- **Future backends** (different GL/Vulkan layers, alternate windowing) are only realistic if game logic and the VM talk to a **small interface**, not to a single third-party API spread across every package.

**Important distinction:** *Decoupling* (package boundaries and interfaces) and *linking* (whether Raylib is in-process via CGO, loaded as a DLL via purego, or archived into a static binary) are **orthogonal**. You want clean architecture **first**; linking strategy is a **build-time** choice layered on top.

---

## 2. Package `hal` (hardware abstraction layer)

Location: [`hal/`](../../hal/).

- **`hal.Driver`** — aggregates three facades:
  - **`Video`** — `VideoDevice` (2D/3D draw, clear, begin/end frame segments).
  - **`Input`** — `InputDevice` (keys, mouse, gamepad).
  - **`System`** — `SystemDevice` (window open/close, FPS, poll, dimensions).
- **Types** — `hal.V2`, `hal.V3`, `hal.RGBA`, `hal.Rect`, `hal.Camera3D`, `hal.Matrix`, etc.  
  **`hal.Matrix`** is documented as **column-major** in the same element order as Raylib’s `rl.Matrix` (see [`hal/types.go`](../../hal/types.go)).

**Rule:** `hal` must **not** import Raylib or CGO. It is pure Go contracts only.

---

## 3. Drivers

### 3.1 Null driver — `drivers/video/null`

Used when no GPU interaction is required: **compiler-only** builds, **`ListBuiltins`**, tests that construct a registry with **`runtime.NewRegistryHeadless`**.

Implements all three `hal` device interfaces as **no-ops** (or trivial constants such as a fixed “screen” size).

### 3.2 Raylib driver — `drivers/video/raylib`

Implements the same interfaces by delegating to **`github.com/gen2brain/raylib-go/raylib`**. Whether that binding uses **CGO** or **purego** is decided by the **Go toolchain and OS**, not by this package’s source.

**`compiler/pipeline`** wires the default HAL driver for **`RunProgram`**:

- **`//go:build fullruntime`** — [`driver_fullruntime.go`](../../compiler/pipeline/driver_fullruntime.go) returns **`hal.Driver`** backed by **`raylib.NewDriver()`** for Video/Input/System.
- **`//go:build !fullruntime`** — [`driver_headless.go`](../../compiler/pipeline/driver_headless.go) returns the **null** driver (compiler CLI, no game runtime).

---

## 4. Where the driver is injected

| Entry point | Driver |
|-------------|--------|
| **`runtime.NewRegistry(h, hal.Driver{...})`** | Caller supplies the full **`hal.Driver`** (used by tests and custom embedders). |
| **`runtime.NewRegistryHeadless(h)`** | **`null`** driver for all three facets ([`registry_headless.go`](../../runtime/registry_headless.go)). |
| **`compiler/pipeline.RunProgram`** (fullruntime) | **`DefaultDriver()`** → Raylib implementation ([`runner.go`](../../compiler/pipeline/runner.go)). |
| **`compiler/pipeline.ListBuiltins`** | Builds registry with **null** HAL, then **`setupRegistry`** so listing keys does not load Raylib ([`registry.go`](../../compiler/pipeline/registry.go)). |

At runtime, natives that need hardware services use **`rt.Driver`** on the active **`runtime.Registry`** (see **`runtime.ActiveRegistry()`** in [`execctx.go`](../../runtime/execctx.go)).

---

## 5. Window backend selection (`internal/driver`)

[`internal/driver`](../../internal/driver/) describes **how** the process talks to Raylib when multiple options exist:

- **`KindNativeCGO`** — binary built with CGO and Raylib **linked** at build time.
- **`KindPuregoDLL`** — **`CGO_ENABLED=0`** on Windows: **purego** loads **`raylib.dll`** (or similar) at runtime.
- **`MOONBASIC_DRIVER`** — override: `auto`, `cgo`, `purego` / `dll` / `sidecar` ([`driver.go`](../../internal/driver/driver.go)).

The **window** module records this via **`Module.BindDriverSelection`** so implementations can choose code paths (e.g. full raygui vs minimal GUI) consistently with the probe result.

---

## 6. Windows purego: deferring `raylib.dll` load for tests

Vendored **`raylib-go`** purego initialization lives under [`third_party/raylib-go-raylib/`](../../third_party/raylib-go-raylib/). **`init()`** can **skip** loading the Raylib shared library when:

1. The executable name looks like a **Go test binary** (basename contains **`.test`**), or  
2. **`MOONBASIC_SKIP_RAYLIB_DLL`** is set to **`1`**, **`true`**, or **`yes`** (case-insensitive).

When load is deferred, **`raylibDll`** stays zero; any code path that still calls DLL-backed functions without guarding first is a bug. Mitigations in-tree include:

- **`FrustumCullDistances`** — on Windows purego without a loaded DLL, a small stub returns safe default near/far values ([`frustum_cull_purego_windows.go`](../../third_party/raylib-go-raylib/frustum_cull_purego_windows.go)); CGO builds use the real rlgl path ([`frustum_cull_cgo.go`](../../third_party/raylib-go-raylib/frustum_cull_cgo.go)).

Production game binaries **must** load Raylib normally (they are not `*.test` and should not set the skip env var in normal use).

---

## 7. Migration status (honest snapshot)

**Direction:** New and refactored code should use **`hal`** types and **`rt.Driver`** for window/video/input surfaces that are already behind the HAL.

**Reality:** Many **`runtime/*`** packages still **`import` Raylib** directly for meshes, textures, audio, shaders, and legacy paths. Full migration is incremental; [`ZERO_CGO_RAYLIB.md`](ZERO_CGO_RAYLIB.md) tracks the size of the remaining `rl.` surface.

Do not treat “HAL complete” as a binary fact in reviews—treat it as a **gradual constraint** on new work.

---

## 8. Static linking and “single executable” builds

- **Purego + sidecar DLL** — Easy cross-compilation of the **Go** binary; users ship **`raylib.dll`** next to it (Windows).
- **CGO + static or semi-static Raylib** — Can produce a **single** user-facing executable, but requires a **C toolchain**, correct **LDFLAGS**, and platform-specific libraries. Cross-compiling static CGO is harder than pure Go.

The repo includes an **experimental** script [`scripts/build_static.ps1`](../../scripts/build_static.ps1) that sets **`CGO_ENABLED=1`** and **`CC`** to **`zig cc`** for a Windows GNU target; see [`BUILDING.md`](../BUILDING.md) for prerequisites (Raylib and friends on the compiler search path). This does not remove the need for careful platform setup; it only centralizes one invocation pattern.

Using **Zig** as **`CC`** is optional; it is one way to get a predictable C linker for Go CGO without relying on a full MSVC setup for every contributor.

---

## 9. CI and local testing expectations

- **Compiler / VM packages** should remain testable with **`CGO_ENABLED=0`** where possible.
- **Full graphics** integration is validated on CI with **Linux + CGO + Xvfb** (see [`.github/workflows/ci.yml`](../../.github/workflows/ci.yml)); Windows jobs include **purego** builds per workflow configuration.
- Prefer **`NewRegistryHeadless`** (or **`ListBuiltins`’s** null-driver path) for tests that only need registration and dispatch, not pixels.

---

## 10. Summary table

| Concern | Mechanism |
|--------|-----------|
| No GPU for tests / listing builtins | **`drivers/video/null`**, **`NewRegistryHeadless`**, **`ListBuiltins`** |
| Interactive **`RunProgram`** (fullruntime) | **`DefaultDriver()`** → **`drivers/video/raylib`** |
| CGO vs DLL on Windows | **`internal/driver`**, build tags, **`MOONBASIC_DRIVER`** |
| Avoid DLL load during **`go test`** (Windows purego) | Deferred **`init()`** + stubs for select calls; **`MOONBASIC_SKIP_RAYLIB_DLL`** |
| Static-ish Windows experiment | **`scripts/build_static.ps1`** + **`BUILDING.md`** |
