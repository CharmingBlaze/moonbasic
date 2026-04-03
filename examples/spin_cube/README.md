# Spinning cube (3D)

Minimal 3D demo: `Camera.Make` (defaults plus explicit framing), **`Mat4.Identity`** + **`Mat4.SetRotation`** for a single reusable matrix (no per-frame alloc/free), **`Mesh.MakeCube`** (GPU upload is automatic), and **`Mesh.Draw`**. Uses handle-style **`cam.Begin()`** / **`cam.End()`** so the active camera is obvious.

Alternatives you can try in code:

- **`Mesh.DrawRotated(cubeMesh, cubeMat, rx#, ry#, rz#)`** — same draw without a matrix handle.
- **`Mat4.Rotation(rx#, ry#, rz#)`** — alias for **`Mat4.FromRotation`** (allocates a new matrix).

**`Window.Open`** returns **`TRUE`** when Raylib reports the window is ready (useful for headless or init failures). **`Window.Close`** releases the heap (Raylib mesh/material cleanup runs via **`Heap.FreeAll`**), so explicit **`Mesh.Free`** / **`Mat4.Free`** are optional if you exit right after close.

Run from the repo root with CGO enabled (same toolchain as `examples/fps`):

```bash
go run . examples/spin_cube/main.mb
```

Controls: **ESC** or close the window to exit.

## What to explore next

- **`examples/fps`** — first-person 3D movement and mouselook.
- **[ARCHITECTURE.md](../../ARCHITECTURE.md) §11** — Phase D 3D roadmap (lighting, terrain, shaders, etc.).
