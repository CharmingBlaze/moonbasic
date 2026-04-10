# Scene engine brief (data-driven levels & performance)

This document captures a **target architecture** for Blender → moonBASIC workflows: glTF-aware loading, shared GPU resources, physics wired to the **WASM/Jolt shared matrix buffer**, and culling—without confusing it with APIs that already exist under similar names.

---

## 1. Three different “scene” concepts today

| Name | Package / role | What it is **not** |
|------|------------------|---------------------|
| **`SCENE.*` (`mbscene`)** | Game **scene state machine**: `SCENE.REGISTER`, `SCENE.LOAD`, `SCENE.UPDATE`, `SCENE.DRAW` | A glTF file loader |
| **`ENTITY.LOADSCENE` / `SCENE.LOADSCENE`** | **Serialize / deserialize** entity ids to JSON ([MEMORY.md](../MEMORY.md)) | Parsing a 3D level format |
| **Raylib `LoadModel` via `ENTITY.LOAD` / `MODEL.LOAD`** | Single-file mesh + materials (glTF/glB, OBJ, … per build) | Full node graph, KHR_lights_punctual, or Blender custom props |

A future **“smart level loader”** should use a **distinct prefix** (e.g. **`LEVEL.LOAD`**, **`GLTF.LOAD`**, or **`SCENEGRAPH.LOAD`**) so docs and tooling do not overload **`SCENE.Load`** from `mbscene`.

---

## 2. Target: scoped scene graph (brief §1)

**Goal:** Move from ad hoc entity creation to a **loader-owned** structure:

- **Parse glTF** in Go (nodes, meshes, skins, materials, extras).
- **Deduplicate** embedded images by hash before `LoadTexture` / material assignment.
- **Map PBR** (`pbrMetallicRoughness`) to the same Raylib PBR path already used in [`DrawEntityModel`](../../runtime/mbmodel3d/model_entity_draw_cgo.go) / material conversion helpers.
- **Tag pass** on `node.Name` (convention):
  - **`Col_*`** → collision-only: no draw (or debug draw); register shape into Jolt / shared buffer.
  - **`Lgt_*`** → spawn punctual lights (requires **`KHR_lights_punctual`** parsing + existing [LIGHT](LIGHT.md) / shader path).

**Status:** Not implemented as a single pipeline; today, artists rely on **`ENTITY.LOAD`** + manual **`BODY3D.*`** / **`ENTITY.LINKPHYSBUFFER`**.

---

## 3. Memory & performance (brief §2)

| Idea | Direction | Current hooks |
|------|-----------|----------------|
| **Zero-copy Jolt sync** | Pre-allocate rows in **`PHYSICS3D.GETMATRIXBUFFER`** / shared float buffer; avoid per-frame body lookup | **`ENTITY.LINKPHYSBUFFER`**, **`PHYSICS3D.STEP`**, **`PHYSICS3D.SYNCWASMTOPHYSREGS`** (WASM paths) — see [PHYSICS3D.md](PHYSICS3D.md) |
| **Instancing** | Second entity references same GPU mesh/materials; only transform + body id differ | **`MODEL.MAKEINSTANCED`** / instanced draw path in `mbmodel3d`; entity-side **`ENTITY.COPY`** duplicates logic, not necessarily GPU sharing — **needs explicit “instance” semantics** if we mirror the brief |
| **Frustum culling** | **`Scene.Draw`**-style batch skips off-screen draws | **`ENTITY.INVIEW`**, cull modes on entities; deferred renderer already has some culling hooks — **no** dedicated “level draw” with hierarchical culling yet |

---

## 4. Command wishlist → implementation owner

The **compiler** parses and emits bytecode; **runtime** owns loaders, heaps, Raylib, and Jolt. Most items below are **runtime + tooling**, not compiler optimizations unless we add intrinsics or static fusion later.

| Brief command | moonBASIC surface | Status / notes |
|---------------|---------------------|----------------|
| **`Scene.Load(path)`** (glTF) | **`LEVEL.LOAD(path)`** | **Implemented** — [LEVEL.md](LEVEL.md); not **`SCENE.LOAD`** (`mbscene`). Uses [qmuntal/gltf](https://github.com/qmuntal/gltf) + one Raylib **`ENTITY.LOAD`** per file. |
| **`Scene.SetRoot(path)`** | **`LEVEL.SETROOT`** | **Implemented** — base path for relative loads. |
| **`Scene.FindEntity(name)`** | **`LEVEL.FINDENTITY`** | **Implemented** — delegates to **`ENTITY.FIND`**. |
| **`Scene.GetMarker(name)`** | **`LEVEL.GETMARKER`** | **Implemented** — world translation (**3 floats**) for named nodes. |
| **`Scene.GetSpawn` / matrix** | **`LEVEL.GETSPAWN`** | **Implemented** — **`MAT4`** world matrix handle. |
| **`Scene.ShowLayer(name)`** | **`LEVEL.SHOWLAYER`** | **Implemented** — **`extras.layer`** or entity group name. |
| **`Scene.ApplyPhysics(id)`** | **`LEVEL.APPLYPHYSICS`** | **Stub** — use **`BODY3D.*`** manually ([PHYSICS3D.md](PHYSICS3D.md)). |
| **`Physics.AutoCreate(id)`** | **`PHYSICS.AUTOCREATE`** | **Stub** — use **`ENTITY.GETBOUNDS`** + **`BODY3D.ADD*`** |
| **`Entity.SetStatic` / `SetTrigger`** | **`ENTITY.SETSTATIC`** / **`SETTRIGGER`** | **Stub** — use **`BODY3D.MAKE("STATIC")`**; sensors N/A in binding yet. |
| **`Material.AutoFilter`** | **`MATERIAL.AUTOFILTER`** | **Stub** — use **`TEXTURE.SETFILTER`**. |
| **`Scene.SyncLights`** | **`LEVEL.SYNCLIGHTS`** | **Stub** — **`KHR_lights_punctual`** not wired. |
| **`Texture.Reload`** | **`TEXTURE.RELOAD`** | **Implemented** — reloads from **`SourcePath`** (main thread). |

---

## 5. WASM shared buffer (brief closing paragraph)

The intended workflow—**bodies created at load time with stable indices into the shared matrix buffer**—aligns with:

- Pre-step sync: **`PHYSICS3D.SYNCWASMTOPHYSREGS`**
- Entity follow: **`ENTITY.LINKPHYSBUFFER(entity, bufferIndex)`** after **`BODY3D.BUFFERINDEX`**

A glTF loader should **reserve buffer slots** when it creates colliders, and document the mapping `{ entityId → buffer row }` so the VM loop does not search for bodies by name.

---

## 6. Suggested phased roadmap

1. **Loader MVP:** glTF → entity ids + `node.Name` map + single PBR material path; no physics automation.
2. **Resources:** Texture hash dedupe + root path for external images.
3. **Tags:** `Col_*` / `Lgt_*` conventions + optional **`KHR_lights_punctual`**.
4. **Physics:** Custom properties → **`BODY3D`** templates + buffer slot reservation.
5. **Draw path:** Batched draw with frustum test per node (or per group).

---

## See also

- [SCENE.md](SCENE.md) — **`SCENE.*`** game scenes (`mbscene`)
- [ENTITY.md](ENTITY.md) — entity ids, groups, **`LOADSCENE`**
- [ANIMATION_3D.md](ANIMATION_3D.md) — **`ENTITY.*`** animation & drawing
- [PHYSICS3D.md](PHYSICS3D.md) — Jolt / WASM buffer
- [ARCHITECTURE.md](../../ARCHITECTURE.md) — overall runtime layout
