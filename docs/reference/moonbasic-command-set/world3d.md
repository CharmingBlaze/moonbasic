# 3D world

| Designed | Implementation | Memory / notes |
|----------|----------------|----------------|
| **CreateWorld ()** | Implicit: **`WINDOW.OPEN`** + **`CAMERA.MAKE`** + loop; or **`ENTITY.CLEARSCENE`** for entity layer | No single “world object” handle in core API. |
| **UpdateWorld** *(Blitz name; not a moonBASIC builtin)* | **`ENTITY.UPDATE(dt)`** with **`dt = TIME.DELTA()`** | Same coordinator Blitz called **UpdateWorld** — explicit API only. |
| **RenderWorld** *(Blitz name; not a moonBASIC builtin)* | **`RENDER.CLEAR`**, **`CAMERA.BEGIN`/`END`** or **`RENDER.Begin3D`/`End3D`**, **`ENTITY.DRAWALL`**, **`RENDER.FRAME`** | Use **`CAMERA2D.Begin`/`End`** for HUD overlays. |
| **ClearWorld ()** | **`ENTITY.CLEARSCENE`**, **`SCENE.CLEARSCENE`** | Frees **all** entities’ native resources — [MEMORY.md](../../MEMORY.md). |
| **SetAmbient (r,g,b)** | **`RENDER.SETAMBIENT`**, **`FOG.SETCOLOR`** context, 2D ambient **`RENDER.SET2DAMBIENT`** | Global render state. |
| **SetFog (r,g,b, near, far)** | **`FOG.SETCOLOR`**, **`FOG.SETRANGE`** / **`SETNEAR`+`SETFAR`** | |
| **SetWireframe (mode)** | **`RENDER.SETWIREFRAME`** | |

See [WORLD.md](../WORLD.md), [BLITZ2025.md](../BLITZ2025.md).
