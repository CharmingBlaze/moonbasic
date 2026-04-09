# High-DPI rendering, native resolution, and post-processing API plan

MoonBASIC targets **native window framebuffer** sizes (4K, Retina) with **float** screen and world coordinates at the script level; the runtime must stay consistent with **GLFW** / **Raylib** framebuffer vs window sizes.

## Coordinate system

- Script and **`DRAW.*`** / **`Camera`** paths already use **float** parameters; internal **`float32`** at the Raylib boundary is normal.
- **No virtual “retro” resolution** is required for new work: treat **`GetScreenWidth` / `GetScreenHeight`** (or engine wrappers) as the drawable pixel size unless a **`RENDER.SET*`** scaler is explicitly added.

## High-DPI / framebuffer

| Concern | Plan |
|---------|------|
| **Blurred UI on Windows** | Ensure **DPI-aware** window creation where GLFW exposes it; use **framebuffer** dimensions for **`rl.BeginDrawing`** / render targets, not only logical window size. |
| **Retina (macOS)** | **Framebuffer** size often **2×** logical size; mouse coordinates may need **`GetMousePosition` / scale** — audit [runtime/window/](../../runtime/window/) and [runtime/input/](../../runtime/input/) for consistency. |
| **Documentation** | When adding APIs, document whether arguments are **pixels**, **logical units**, or **normalized**. |

## Post-processing (MSAA, bloom, SSAO)

Raylib 5.x exposes pieces of this via **shaders**, **render textures**, and **RLGL**. Planned **high-level** surface (names illustrative; final keys follow **`commands.json`**):

| User-facing intent | Implementation sketch |
|--------------------|------------------------|
| **MSAA** | `rl.EnableSmoothMultisampling()` / MSAA render target + blit to window (platform-dependent). |
| **Bloom** | Post-process pass: bright pass → blur → combine (custom shader + full-screen quad). |
| **SSAO** | Depth/normal G-buffer pass (requires 3D depth access); heavier; optional **quality** preset. |

**Rule:** Feature-detect; if unsupported, **no-op** with a logged warning in **`--debug`** mode.

## Zero-allocation draw path

- Post FX should **reuse** render textures and shaders for the session, not allocate per frame.
- See [docs/audit/IR_V3_VM_AUDIT.md](../audit/IR_V3_VM_AUDIT.md) and profiling notes in [vm/](../../vm/).

## File touchpoints

- Window / GL context: [runtime/window/raylib_cgo.go](../../runtime/window/raylib_cgo.go), [post_process_cgo.go](../../runtime/window/post_process_cgo.go).
- Future MSAA/bloom: extend the same module family with **`Register`** entries and **stub** parity.
