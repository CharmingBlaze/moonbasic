# Instancing — `MODEL.*` and `INSTANCE.*`

GPU **instancing** (many copies of one loaded model with per-instance transforms) is implemented in **`runtime/mbmodel3d`**. There is no separate engine type namespace: heap handles report type **`InstancedModel`**, and commands are exposed under **`MODEL.*`** with **equivalent `INSTANCE.*` aliases** (same handlers).

---

## Registry keys

| Command | Arguments | Returns | Notes |
|---------|-----------|---------|--------|
| **`MODEL.MAKEINSTANCED`** / **`INSTANCE.MAKEINSTANCED`** | `path$`, `instanceCount` | handle | Loads a model from disk (`LoadModel`), allocates transforms for **`instanceCount`** instances (1…200000). |
| **`MODEL.SETINSTANCEPOS`** / **`INSTANCE.SETINSTANCEPOS`** | `inst`, `index`, `x#`, `y#`, `z#` | — | Sets world position for one instance. |
| **`MODEL.SETINSTANCESCALE`** / **`INSTANCE.SETINSTANCESCALE`** | `inst`, `index`, `sx#`, `sy#`, `sz#` | — | Per-instance scale (used when building transform matrices). |
| **`MODEL.UPDATEINSTANCES`** / **`INSTANCE.UPDATEINSTANCES`** | `inst` | — | Rebuilds instance **`Matrix`** array from positions and scales (translate × scale). |
| **`MODEL.DRAW`** | `inst` | — | **`DrawMeshInstanced`** for the instanced handle (also used for regular **`Model`** / **`LODModel`** with different dispatch). |

Use **`CAMERA.BEGIN` / `CAMERA.END`** around drawing (see [CAMERA.md](CAMERA.md)).

---

## Handle methods (`InstancedModel`)

On an **`InstancedModel`** handle, method calls map to **`INSTANCE.*`** (not `MODEL.*`) for **`SetInstancePos`**, **`SetInstanceScale`**, **`UpdateInstances`**; **`Draw`** still maps to **`MODEL.DRAW`** (shared with non-instanced models).

Examples: **`inst.SetInstancePos`**, **`Instance.SetInstancePos`** → **`INSTANCE.SETINSTANCEPOS`**.

---

## See also

- [MODEL.md](MODEL.md) — loading, materials, hierarchy, **`MODEL.INSTANCE`** (scene graph clone)
- [MESH.md](MESH.md) — mesh primitives if you build geometry without a file
