# Instancing — `INSTANCE.*` and `MODEL.*`

GPU **instancing** draws many copies of one loaded model with **per-instance transforms** (`DrawMeshInstanced` in Raylib). Heap handles report type **`InstancedModel`**.

---

## Creating an instanced batch

| Command | Arguments | Returns | Notes |
|---------|-----------|---------|--------|
| **`INSTANCE.MAKEINSTANCED`** / **`MODEL.MAKEINSTANCED`** | `path$`, `instanceCount` | handle | `LoadModel(path)` then allocate `instanceCount` slots (1…200000). |
| **`INSTANCE.MAKE`** | `model`, `instanceCount` | handle | Reloads from the **same path** as **`MODEL.LOAD`** (source model must have been loaded from disk). Fails for procedural **`MODEL.MAKE(mesh)`**-only models — use **`MAKEINSTANCED`** with an asset path instead. |

---

## Per-instance state

Indices are **`0 .. INSTANCE.COUNT(inst)-1`**.

| Command | Arguments | Notes |
|---------|-----------|--------|
| **`INSTANCE.SETPOS`** / **`INSTANCE.SETINSTANCEPOS`** / **`MODEL.SETINSTANCEPOS`** | `inst`, `index`, `x#`, `y#`, `z#` | World position. Clears a manual **`SETMATRIX`** for that index. |
| **`INSTANCE.SETROT`** | `inst`, `index`, `rx#`, `ry#`, `rz#` | Euler rotation **radians** (same order as **`MatrixRotateXYZ`**). Clears manual matrix for that index. |
| **`INSTANCE.SETSCALE`** / **`INSTANCE.SETINSTANCESCALE`** / **`MODEL.SETINSTANCESCALE`** | `inst`, `index`, `sx#`, `sy#`, `sz#` | Clears manual matrix for that index. |
| **`INSTANCE.SETMATRIX`** | `inst`, `index`, `mat` | **`Matrix4`** handle — full row-major transform for that instance. **Manual** mode: **`UPDATEBUFFER`** / **`UPDATEINSTANCES`** skips that index until **`SETPOS`**, **`SETROT`**, or **`SETSCALE`** clears it. |
| **`INSTANCE.SETCOLOR`** | `inst`, `index`, `r`, `g`, `b`, `a` | Tint **0–255** (stored as float). **Uniform** color across all instances uses **`DrawMeshInstanced`**. **Different** colors per instance use a **per-draw `DrawMesh`** loop (slower). |

---

## Buffer update and draw

| Command | Arguments | Notes |
|---------|-----------|--------|
| **`INSTANCE.UPDATEBUFFER`** / **`INSTANCE.UPDATEINSTANCES`** / **`MODEL.UPDATEINSTANCES`** | `inst` | Rebuilds **`Matrix`** from **T × R × S** for non-manual instances. Call after **`Set*`** changes. |
| **`INSTANCE.DRAW`** | `inst` | Instanced draw only (errors if not **`InstancedModel`**). |
| **`MODEL.DRAW`** | `inst` | **Model**, **LODModel**, or **InstancedModel** (shared entry). |
| **`INSTANCE.DRAWLOD`** | `inst`, `lodMesh`, `dist#` | If camera distance from the **centroid** of instance positions exceeds **`dist`**, draws using **`lodMesh`** (same transforms/material as the primary mesh). Otherwise uses the default mesh. |
| **`INSTANCE.SETCULLDISTANCE`** | `inst`, `dist#` | If **`dist > 0`**, skips **draw** (and shadow) when the camera is farther than **`dist`** from that centroid. **`0`** disables. |

---

## Lifecycle and queries

| Command | Arguments | Returns |
|---------|-----------|---------|
| **`INSTANCE.COUNT`** | `inst` | `int` |
| **`INSTANCE.FREE`** / **`MODEL.FREE`** | handle | Frees **`InstancedModel`** or **`Model`** (same heap **`Free`**). |

---

## Handle methods (`InstancedModel`)

`inst.Method(...)` maps to **`INSTANCE.*`** keys (prefix **`INSTANCE.`**). Examples:

- **`inst.SetPos`** / **`Instance.SetPos`** → **`INSTANCE.SETPOS`**
- **`inst.Draw`** → **`MODEL.DRAW`** (shared draw path)
- **`inst.DrawLOD mesh, dist`** → **`INSTANCE.DRAWLOD`**
- **`inst.Free`** → **`INSTANCE.FREE`**

---

## See also

- [MODEL.md](MODEL.md) — loading, materials, **`MODEL.INSTANCE`** (scene graph clone; not GPU instancing)
- [MESH.md](MESH.md) — mesh primitives
- [CAMERA.md](CAMERA.md) — **`CAMERA.BEGIN` / `CAMERA.END`**
