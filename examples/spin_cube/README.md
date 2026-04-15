# Spinning cube (3D)

Minimal 3D demo: `CreateCamera()` (Easy Mode → `CAMERA.CREATE`; defaults plus explicit framing), **`Transform.Identity`** + **`Transform.SetRotation`** for a single reusable matrix (no per-frame alloc/free), **`Mesh.MakeCube`** (GPU upload is automatic), and **`Mesh.Draw`**. Uses handle-style **`cam.Begin()`** / **`cam.End()`** so the active camera is obvious.

Alternatives you can try in code:

- **`Mesh.DrawRotated(cubeMesh, cubeMat, rx#, ry#, rz#)`** — same draw without a matrix handle.
- **`Transform.Rotation(rx#, ry#, rz#)`** — allocates a new rotation matrix when you cannot reuse one handle with **`SetRotation`**.

**`Window.Open`** returns **`TRUE`** when Raylib reports the window is ready (useful for headless or init failures). **`Window.Close`** releases the heap (Raylib mesh/material cleanup runs via **`Heap.FreeAll`**), so explicit **`Mesh.Free`** / **`Transform.Free`** are optional if you exit right after close.

Run from the repo root with **CGO** enabled (same toolchain as `examples/fps`).

**Run the demo** (opens a window):

```bash
CGO_ENABLED=1 go run -tags fullruntime ./cmd/moonrun examples/spin_cube/main.mb
```

**Compile only** (writes `main.mbc`, no window): `CGO_ENABLED=1 go run . examples/spin_cube/main.mb`

Controls: **ESC** or close the window to exit.

## What to explore next

- **`examples/fps`** — top-down arena (WASD + mouse aim); good next step for game loop patterns.
- **[ARCHITECTURE.md](../../ARCHITECTURE.md) §11** — Phase D 3D roadmap (lighting, terrain, shaders, etc.).
