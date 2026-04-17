# Instance Commands

GPU instancing: draw many copies of one model with per-instance transforms and colors.

Page shape follows [DOC_STYLE_GUIDE.md](../DOC_STYLE_GUIDE.md) (**WAVE pattern**).

## Core Workflow

1. Create with `MODEL.CREATEINSTANCED(path, count)` or `INSTANCE.CREATE(model, count)`.
2. Set per-instance transforms with `INSTANCE.SETPOS`, `INSTANCE.SETROT`, `INSTANCE.SETSCALE`.
3. Call `INSTANCE.UPDATEBUFFER` after changes.
4. Draw with `INSTANCE.DRAW` or `MODEL.DRAW`.
5. Free with `INSTANCE.FREE`.

---

## Creating an instanced batch

| Command | Arguments | Returns | Notes |
|---------|-----------|---------|--------|
| **`MODEL.CREATEINSTANCED`** (canonical) / **`MODEL.MAKEINSTANCED`** (deprecated); **`INSTANCE.CREATEINSTANCED`** (canonical) / **`INSTANCE.MAKEINSTANCED`** (deprecated) | `path`, `instanceCount` | handle | `LoadModel(path)` then allocate `instanceCount` slots (1…200000). |
| **`INSTANCE.CREATE`** (canonical) / **`INSTANCE.MAKE`** (deprecated) | `model`, `instanceCount` | handle | Reloads from the **same path** as **`MODEL.LOAD`** (source model must have been loaded from disk). Fails for procedural **`MODEL.CREATE(mesh)`** / **`MODEL.MAKE(mesh)`**-only models — use **`CREATEINSTANCED`** / **`MAKEINSTANCED`** with an asset path instead. |

---

## Per-instance state

Indices are **`0 .. INSTANCE.COUNT(inst)-1`**.

| Command | Arguments | Notes |
|---------|-----------|--------|
| **`INSTANCE.SETPOS`** / deprecated **`INSTANCE.SETPOSITION`**; **`INSTANCE.SETINSTANCEPOS`** / **`MODEL.SETINSTANCEPOS`** | `inst`, `index`, `x`, `y`, `z` | World position. Clears a manual **`SETMATRIX`** for that index. |
| **`INSTANCE.SETROT`** | `inst`, `index`, `rx`, `ry`, `rz` | Euler rotation **radians** (same order as **`MatrixRotateXYZ`**). Clears manual matrix for that index. |
| **`INSTANCE.SETSCALE`** / **`INSTANCE.SETINSTANCESCALE`** / **`MODEL.SETINSTANCESCALE`** | `inst`, `index`, `sx`, `sy`, `sz` | Clears manual matrix for that index. |
| **`INSTANCE.SETMATRIX`** | `inst`, `index`, `mat` | **`Matrix4`** handle — full row-major transform for that instance. **Manual** mode: **`UPDATEBUFFER`** / **`UPDATEINSTANCES`** skips that index until **`SETPOS`**, **`SETROT`**, or **`SETSCALE`** clears it. |
| **`INSTANCE.SETCOLOR`** | `inst`, `index`, `r`, `g`, `b`, `a` | Tint **0–255** (stored as float). **Uniform** color across all instances uses **`DrawMeshInstanced`**. **Different** colors per instance use a **per-draw `DrawMesh`** loop (slower). |

---

## Buffer update and draw

| Command | Arguments | Notes |
|---------|-----------|--------|
| **`INSTANCE.UPDATEBUFFER`** / **`INSTANCE.UPDATEINSTANCES`** / **`MODEL.UPDATEINSTANCES`** | `inst` | Rebuilds **`Matrix`** from **T × R × S** for non-manual instances. Call after **`Set*`** changes. |
| **`INSTANCE.DRAW`** | `inst` | Instanced draw only (errors if not **`InstancedModel`**). |
| **`MODEL.DRAW`** | `inst` | **Model**, **LODModel**, or **InstancedModel** (shared entry). |
| **`INSTANCE.DRAWLOD`** | `inst`, `lodMesh`, `dist` | If camera distance from the **centroid** of instance positions exceeds **`dist`**, draws using **`lodMesh`** (same transforms/material as the primary mesh). Otherwise uses the default mesh. |
| **`INSTANCE.SETCULLDISTANCE`** | `inst`, `dist` | If **`dist > 0`**, skips **draw** (and shadow) when the camera is farther than **`dist`** from that centroid. **`0`** disables. |

---

## Lifecycle and queries

| Command | Arguments | Returns |
|---------|-----------|---------|
| **`INSTANCE.COUNT`** | `inst` | `int` |
| **`INSTANCE.FREE`** / **`MODEL.FREE`** | handle | Frees **`InstancedModel`** or **`Model`** (same heap **`Free`**). |

---

## Handle methods (`InstancedModel`)

`inst.Method(...)` maps to **`INSTANCE.*`** keys (prefix **`INSTANCE.`**). Examples:

- **`inst.pos()`** / **`inst.rot()`** / **`inst.scale()`** with **no arguments** → **`INSTANCE.GETPOS`** / **`GETROT`** / **`GETSCALE`** (transform of **instance index 0**; use **`INSTANCE.SETINSTANCEPOS`** / **`INSTANCE.SETROT`** / **`SETSCALE`** with an index for other slots)
- **`inst.SetPos`** / **`Instance.SetPos`** → **`INSTANCE.SETPOS`**
- **`inst.Draw`** → **`MODEL.DRAW`** (shared draw path)
- **`inst.DrawLOD mesh, dist`** → **`INSTANCE.DRAWLOD`**
- **`inst.Free`** → **`INSTANCE.FREE`**

See [UNIVERSAL_HANDLE_METHODS.md](UNIVERSAL_HANDLE_METHODS.md).

---

## Full Example

```basic
WINDOW.OPEN(1280, 720, "Instancing Demo")
WINDOW.SETFPS(60)

cam = CAMERA.CREATE()
CAMERA.SETPOS(cam, 0, 10, -20)
CAMERA.SETTARGET(cam, 0, 0, 0)

inst = MODEL.CREATEINSTANCED("tree.glb", 100)
FOR i = 0 TO 99
    INSTANCE.SETPOS(inst, i, RND(20) - 10, 0, RND(20) - 10)
    INSTANCE.SETSCALE(inst, i, 0.5 + RNDF(0, 1), 0.5 + RNDF(0, 1), 0.5 + RNDF(0, 1))
NEXT
INSTANCE.UPDATEBUFFER(inst)

WHILE NOT WINDOW.SHOULDCLOSE()
    RENDER.CLEAR(40, 60, 80)
    RENDER.BEGIN3D(cam)
        INSTANCE.DRAW(inst)
    RENDER.END3D()
    RENDER.FRAME()
WEND

INSTANCE.FREE(inst)
CAMERA.FREE(cam)
WINDOW.CLOSE()
```

---

## See also

- [MODEL.md](MODEL.md) — loading, materials, **`MODEL.INSTANCE`** (scene graph clone; not GPU instancing)
- [MESH.md](MESH.md) — mesh primitives
- [CAMERA.md](CAMERA.md) — **`CAMERA.BEGIN` / `CAMERA.END`**
