# moonBASIC memory model (entity / scene)

This document describes **ownership and cleanup** for the **ENTITY** module paths that use **raylib** (`rl.Model`, `rl.ModelAnimation`, procedural meshes). It is the reference for `ENTITY.LOADMESH`, `ENTITY.CREATEMESH`, `ENTITY.LOADANIMATEDMESH`, `ENTITY.COPY`, `ENTITY.SAVESCENE` / `ENTITY.LOADSCENE`, and `ENTITY.CLEARSCENE`.

## Three layers (short)

1. **Go heap** — `ent` structs, maps, slices: freed when unreferenced (GC).
2. **CGo / raylib** — `rl.LoadModel`, `rl.LoadModelAnimations`, `rl.GenMeshCube` / `rl.LoadModelFromMesh`: **not** visible to the Go GC; must be released with `rl.UnloadModel` / `rl.UnloadModelAnimations` / `rl.UnloadMesh` as documented by raylib.
3. **Entity IDs** — integer handles into `entityStore.ents`; **not** the same as VM heap handles. Lifetime is managed only by **ENTITY.FREE** and **ENTITY.CLEARSCENE** (and failed **ENTITY.LOADSCENE** rollback — see below).

## What each entity owns

| Field / state | Owned? | Released in |
|---------------|--------|-------------|
| `rlModel` when `hasRLModel` | Yes | `entFree`: `UnloadModelAnimations` then `UnloadModel` |
| `modelAnims` | Yes | `entFree`: `UnloadModelAnimations` before model |
| Procedural mesh used only to build the model | Yes (transient) | `CREATEMESH` / scene mesh: `UnloadMesh` immediately after `LoadModelFromMesh` |

**Destruction order** (required by raylib): unload **animations first**, then **model**.

## ENTITY.LOADSCENE and partial failure

`ENTITY.LOADSCENE` clears the scene, then builds entities in a loop. If **any** step fails after Raylib resources were created, the implementation sets a **rollback** flag: on return with error it calls **`ENTITY.CLEARSCENE`**, which walks all entities and calls **`ENTITY.FREE`** for each, so **no `rl.Model` / animation array is left allocated** for entities that were only partially loaded.

Successful loads validate `MeshCount > 0` after `LoadModel` / `LoadModelFromMesh` so empty failed loads do not commit live entities.

## ENTITY.LOADMESH, ENTITY.CREATEMESH, ENTITY.COPY

- **LOADMESH** loads the model **before** allocating a new entity id. If the load is empty (`MeshCount <= 0`), the model is **unloaded** and the function returns an error **without** consuming the next id.
- **CREATEMESH** unloads the mesh after `LoadModelFromMesh`; if the model is invalid, it **unloads the model** and returns an error **without** registering an entity.
- **COPY** reloads from `loadPath` for mesh-backed entities. It **rejects** duplication of procedural meshes without a path (e.g. **CREATEMESH**) **before** allocating a new id. Failed `LoadModel` results in **UnloadModel** and an error **without** bumping `nextID`.

## ENTITY.CLEARSCENE

Clears groups, resets `nextID` to 1, and frees every entity by calling **`entFree`** for each id so Raylib resources are always released.
