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

## VM heap tags (physics / net)

These use the same **`FREE`** / **`ERASE`** rules as other **`KindHandle`** objects:

| Handle | Typical free |
|--------|----------------|
| **`JOINT2D.*`** | **`JOINT2D.FREE`** |
| **`PACKET.CREATE`** (or after **`PEER.SENDPACKET`**, ownership transfers—see [network-enet.md](reference/moonbasic-command-set/network-enet.md)) | **`PACKET.FREE`** if not sent |

---

## Game orbit helpers (`ORBITYAWDELTA` / `ORBITPITCHDELTA` / `ORBITDISTDELTA`)

These **`GAME`** builtins return **numeric floats only** (radians or distance delta). They do **not** allocate VM heap objects — **no `ERASE`**. Pair them with your own **`camYaw#` / `camPitch#` / `camDist#`** variables and **`Camera.SetOrbit`** (see [GAMEHELPERS.md](reference/GAMEHELPERS.md)).

---

## Runtime heap: three layers and the handle table

This section complements the **ENTITY** tables above. VM **`KindHandle`** values index **`heap.Store`** (`vm/heap`).

| Layer | What | Who frees |
|-------|------|-----------|
| **Go heap** | Slices, maps, strings, channels | GC when unreachable |
| **CGo (raylib, Jolt, Box2D, ENet, OS)** | Textures, models, bodies, hosts | Explicit **`Unload` / `Destroy` / `Release`** in **`HeapObject.Free()`** or subsystem teardown |
| **Handle table** | Opaque **`int32`** handles (generation + slot) | **`Store.Free`**, **`Store.FreeAll`**, **`ERASE ALL` / `FREE.ALL`**, pipeline shutdown |

The GC does **not** see C allocations: setting a Go field to `nil` does not free GPU or physics memory.

### Idempotent `Free()` — `ReleaseOnce` vs `freed bool`

`HeapObject.Free()` must be safe to call more than once. Implementations typically use:

- **`heap.ReleaseOnce`** — wraps the native cleanup so only the first call runs (see `vm/heap/release_once.go`), or
- an explicit **`if o.freed { return }`** with **`o.freed = true`** at the end.

Either pattern satisfies the contract in `heap.go`.

### Owned vs shared (textures)

**`TextureObject`** (`runtime/texture/heap_objects_cgo.go`) uses **`Borrowed`**: when true (e.g. view from a render target’s texture), **`Free()`** does **not** unload the GPU texture — the owning **`RenderTargetObject`** does.

### Destruction order (native)

| Resource | Must be freed before / ordering notes |
|----------|--------------------------------------|
| **Jolt constraint / joint** | Remove from world before destroying bodies that still reference it (script/API contract) |
| **Jolt body** | Remove from world → destroy body → **`shape.Release()`** after the body is gone |
| **Model animations** | **`UnloadModelAnimation`** (each) before **`UnloadModel`** |
| **Terrain chunk** | Remove physics body from world → destroy body → release shape → **`UnloadMesh`** (see terrain unload path) |
| **ENet** | **`enet_packet_destroy`** after copying data out; host **`Destroy`** disconnects peers |

### `Store.FreeAll` guarantee

On shutdown, **`RunProgram`** tears down the registry and calls **`Store.FreeAll`**, which walks **every** slot and invokes **`HeapObject.Free()`** on live objects. Long-running games should still **`FREE`** / **`ERASE`** handles they no longer need to cap memory during play; shutdown is the safety net for leaks of VM handles.

### Audit baselines (repo root)

Optional regression artifacts (see **`ARCHITECTURE.md`** for optional Valgrind/gccheckmark notes):

- **`RACE_BASELINE.txt`** — `go test -race ./...`
- **`GCCHECK_BASELINE.txt`** — `GODEBUG=gccheckmark=1 go test ./...`
- **`ESCAPE_ANALYSIS.txt`** — `go build -gcflags="-m -m"` (main package; CLI-heavy escapes are expected)
- **`CGO_ALLOCS.txt`** — grep of raylib/Jolt/ENet allocation sites to cross-check **`Free()`**
- **`HEAP_AUDIT.txt`** — HeapObject contract checklist
- **`VALGRIND_BASELINE.txt`**, **`DRMEM_BASELINE.txt`** — placeholders on Windows; run on Linux or with Dr. Memory locally

### Hot-path allocations

The VM execute loop and per-frame draw paths should avoid unnecessary allocations. Profile with **`go test -bench`**, **`runtime/pprof`**, or **`GODEBUG=gctrace=1`**; there is no moonBASIC **`DEBUG.FRAMEALLOC`** builtin yet — use Go tooling for engine-side regression.

### Lifecycle policy: explicit `Free()` vs `runtime.SetFinalizer`

The engine’s **supported** model for C/C++ resources (Raylib, Jolt, Box2D, ENet) is:

1. **VM heap handles** — every native-backed object implements **`HeapObject.Free()`** idempotently (**`ReleaseOnce`** or equivalent).
2. **Deterministic teardown** — scripts or **`Store.FreeAll`** on shutdown must release C memory; the Go GC **does not** see C heap usage.
3. **`runtime.SetFinalizer`** is **not** used as the primary release path: finalizers run **nondeterministically**, may run **late** or **never** under memory pressure, and are **easy to misuse** with finalizer cycles.

**Optional future addition:** a **debug-only** or **`moonbasic_leakcheck`** build tag that attaches a **finalizer sentinel** logging “handle GC’d without Free” — useful for finding leaks in CI, **without** replacing explicit **`FREE`** / **`ERASE`** in shipped games.

This policy matches the contract in **[vm/heap/heap.go](../vm/heap/heap.go)** and supersedes any directive that suggests GC-driven **`Unload`** as the default.
