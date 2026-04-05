# Instancing — `MODEL.*` instance APIs

GPU instancing for many copies of a model is exposed as **`MODEL.MAKEINSTANCED`**, **`MODEL.SETINSTANCEPOS`**, **`MODEL.SETINSTANCESCALE`**, **`MODEL.UPDATEINSTANCES`**, and instanced **`MODEL.DRAW`** (see **`runtime/mbmodel3d`**).

There is no separate `Instance.*` namespace in the engine: use the **`MODEL.*`** names above so materials and meshes stay unified.

---

## See also

- [MODEL.md](MODEL.md) — loading, materials, hierarchy
- [MESH.md](MESH.md) — raw mesh generation
