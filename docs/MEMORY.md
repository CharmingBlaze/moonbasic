# moonBASIC memory model (entity / scene)

## VM heap: `ERASE ALL` / `FREE.ALL`

**`ERASE ALL`** (statement) and **`FREE.ALL`** (builtin, no arguments) do the same thing:

1. Call **`heap.Store.FreeAll`** — every registered **`HeapObject`** is released (including nested handles inside handle-arrays, with correct ordering).
2. Set every **`KindHandle`** value in **global variables** and on the **operand stack** to **null**, so scripts cannot keep stale integer handles.

Use this when tearing down a program or resetting a scene’s VM-backed resources in one step instead of many **`ERASE`** lines. **Do not** rely on it mid-expression: any handle temporarily on the stack for a pending operation would be cleared.

**Example:** [`examples/mario64/main_orbit_simple.mb`](../examples/mario64/main_orbit_simple.mb) ends with **`ERASE ALL`** then **`Window.Close()`** after the main loop (camera + platform **`DIM`** arrays are VM handles).

**Not covered:** numeric **entity IDs** from **`ENTITY.***` are **not** VM heap handles — use **`ENTITY.CLEARSCENE`** / **`ENTITY.FREE`** as before. **Window**, **input**, and other non–heap-backed state are unchanged.

The identifier **`ALL`** is reserved for this statement form; do not use **`ALL`** as a variable name if you need **`ERASE varname`** for a single array.

---

This document also describes **ownership and cleanup** for the **ENTITY** module paths that use **raylib** (`rl.Model`, `rl.ModelAnimation`, procedural meshes). It is the reference for `ENTITY.LOADMESH`, `ENTITY.CREATEMESH`, `ENTITY.LOADANIMATEDMESH`, `ENTITY.COPY`, `ENTITY.SAVESCENE` / `ENTITY.LOADSCENE`, and `ENTITY.CLEARSCENE`.

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

---

## Game orbit helpers (`ORBITYAWDELTA` / `ORBITPITCHDELTA` / `ORBITDISTDELTA`)

These **`GAME`** builtins return **numeric floats only** (radians or distance delta). They do **not** allocate VM heap objects — **no `ERASE`**. Pair them with your own **`camYaw#` / `camPitch#` / `camDist#`** variables and **`Camera.SetOrbit`** (see [GAMEHELPERS.md](reference/GAMEHELPERS.md)).
