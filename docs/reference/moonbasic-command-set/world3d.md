# 3D world

| Designed | Implementation | Memory / notes |
|----------|----------------|----------------|
| **CreateWorld ()** | Implicit: **`WINDOW.OPEN`** + **`CAMERA.MAKE`** + loop; or **`ENTITY.CLEARSCENE`** for entity layer | No single “world object” handle in core API. |
| **UpdateWorld (speed#)** | **`TIME.DELTA`**, **`ENTITY.UPDATE`**, your systems | **Entity ids** updated in place — no extra free per frame. |
| **RenderWorld ()** | **`RENDER.CLEAR`**, **`CAMERA.BEGIN`/`END`**, **`ENTITY.DRAWALL`**, **`RENDER.FRAME`** | Order matters for Z-buffer. |
| **ClearWorld ()** | **`ENTITY.CLEARSCENE`**, **`SCENE.CLEARSCENE`** | Frees **all** entities’ native resources — [MEMORY.md](../../MEMORY.md). |
| **SetAmbient (r,g,b)** | **`RENDER.SETAMBIENT`**, **`FOG.SETCOLOR`** context, 2D ambient **`RENDER.SET2DAMBIENT`** | Global render state. |
| **SetFog (r,g,b, near, far)** | **`FOG.SETCOLOR`**, **`FOG.SETRANGE`** / **`SETNEAR`+`SETFAR`** | |
| **SetWireframe (mode)** | **`RENDER.SETWIREFRAME`** | |

See [WORLD.md](../WORLD.md), [BLITZ2025.md](../BLITZ2025.md).
