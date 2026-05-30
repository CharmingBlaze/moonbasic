# Level Commands

Load glTF levels, resolve markers/spawns, toggle layers, preload assets, and set up physics.

Page shape follows [DOC_STYLE_GUIDE.md](../DOC_STYLE_GUIDE.md) (**WAVE pattern**).

## Core Workflow

1. Set asset root with `LEVEL.SETROOT`.
2. Load a glTF file with `LEVEL.LOAD`.
3. Query spawn points with `LEVEL.GETSPAWN`, find entities with `LEVEL.FINDENTITY`.
4. Toggle visibility with `LEVEL.SHOWLAYER` / `LEVEL.HIDELAYER`.
5. Optionally set up physics with `LEVEL.SETUP` and static collision with `LEVEL.STATIC` / `LEVEL.AUTOCOLLIDE`.

For game-scene switching see [SCENE.md](SCENE.md).

---

### `LEVEL.SETROOT(path)`
Sets the base directory for relative asset paths.

- **Returns**: (None)

---

### `LEVEL.LOAD(path)`
Loads a glTF level and resolves all markers and spawns.

- **Arguments**:
    - `path`: (String) File path relative to the root.
- **Returns**: (Integer) The ID of the root level entity.
- **Example**:
    ```basic
    LEVEL.SETROOT("assets/")
    lvl = LEVEL.LOAD("level1.glb")
    ```

---

### `LEVEL.STATIC(entity)` / `AUTOCOLLIDE`
Bakes high-performance static collision geometry for levels.

- **Returns**: (None)

---

### `LEVEL.GETSPAWN(name)`
Returns the world transform of a named node.

- **Arguments**:
    - `name`: (String) The node name from the glTF file.
- **Returns**: (Handle) A matrix handle (`MAT4`).

**Limits (current):**

- **One** combined Raylib model per file — not per-node mesh instancing. Multi-mesh scenes still get correct **named transforms** for empties and nodes; see [SCENE_ENGINE_BRIEF.md](SCENE_ENGINE_BRIEF.md) for the full roadmap (texture dedupe, **`ENTITY.INSTANCE`**, Jolt buffer prealloc).
- **`Col_*`** mesh nodes are treated as collision-oriented: if the chosen visual node is **`Col_`**, the entity is **hidden**. **`Col_*`** transforms are also appended to an internal collider list for future **`LEVEL.APPLYPHYSICS`**.
- **`Lgt_*`** and **`KHR_lights_punctual`** are not converted to **`LIGHT.*`** yet.

---

## Names, layers, metadata

| Command | Purpose |
|--------|---------|
| **`LEVEL.FINDENTITY(name)`** | Same as **`ENTITY.FIND`** — looks up **`ENTITY.SETNAME`** / loader-assigned names. |
| **`LEVEL.GETMARKER(name)`** | 3-float array: **translation** from the named node’s world matrix (empties and mesh nodes). |
| **`LEVEL.GETSPAWN(name)`** | **`MAT4`** handle: full **world** matrix for that node name. |
| **`LEVEL.SHOWLAYER(layerName, visible)** | Shows/hides entities registered to a **`layer`** extra on the primary loaded mesh node, or falls back to **`ENTITY.GROUPCREATE`** membership for the same **`layerName`**. |
| **`ENTITY.GETMETADATA(entity, key)`** → **string** | Reads flattened **glTF extras** from the primary mesh node used by **`LEVEL.LOAD`** (nested keys use **`.`**, e.g. **`door.options.label`**). Blender custom properties are typically surfaced here. Empty string if missing. |

Node **`extras`** may include JSON **`{"layer":"MyLayer"}`** (string or number) to associate the root loaded entity with a layer for **`SHOWLAYER`**. A string **`tag`** in extras is stored for **`MATERIAL.BULKASSIGN`** matching.

---

## Global textures & materials

| Command | Purpose |
|--------|---------|
| **`TEXTURE.SETDEFAULTFILTER(mode)`** | Sets the **default min/mag filter** for **new** file loads (e.g. **`FILTER_POINT`** for a PS1 look). Pass **`-1`** to clear and use the normal **`TEXTURE.LOAD`** flag presets again. |
| **`MATERIAL.AUTOFILTER(mode)`** | Alias of **`TEXTURE.SETDEFAULTFILTER`** (same engine hook). |
| **`MATERIAL.BULKASSIGN(pattern, textureHandle [, materialIndex])`** → **count** | For every entity with a model whose **`ENTITY`** name **or** Blender **`tag`** extra matches **`pattern`** (case-insensitive `path` glob: **`*`**, **`?`**), sets the albedo map on **material 0** by default, or **`materialIndex`** when given. Returns how many entities were updated. |
| **`RENDER.CLEARCACHE`** | **`TEXTURE.FREE`** on all handles recorded by **`LEVEL.PRELOAD`** (safe between levels to drop unused preload textures). Other textures are unaffected. |

---

## Atmosphere & tone mapping

| Command | Purpose |
|--------|---------|
| **`LEVEL.LOADSKYBOX(hdrPath)`** → **texture handle** | Loads an HDR (or other image) through **`TEXTURE.LOAD`** rules. Returns a **texture handle** for drawing (e.g. sky sphere) or post workflows. **IBL** and automatic env lighting are **not** wired to PBR yet. |
| **`RENDER.SETTONEMAPPING(mode)`** | Alias of **`POST.SETTONEMAP`**: **0** none, **1** Reinhard, **2** Filmic, **3** ACES (requires post stack). |

---

## Script binding (data-driven dispatch)

The VM does not auto-call BASIC functions on collision yet. **`LEVEL.BINDSCRIPT`** registers **glob patterns** → **function names**; you resolve them at runtime:

| Command | Purpose |
|--------|---------|
| **`LEVEL.BINDSCRIPT(pattern, functionName)`** | Records a binding (e.g. **`GoldCoin*`** → **`CollectCoin`**). |
| **`LEVEL.MATCHSCRIPTBIND(objectName)`** → **string** | Returns the **first** matching **`functionName`**, or **empty** if none. Use with **`EntityName`**, ray hits, or physics contact to **`SELECT`** / branch in BASIC. |

---

## Physics triggers & optimization

| Command | Status |
|--------|--------|
| **`TRIGGER.CREATEFROMENTITY(entity)`** | Not implemented — Jolt sensor-from-mesh still blocked on bindings; use **`ENTITY.SETTRIGGER`** when exposed. |
| **`LEVEL.OPTIMIZE(entity)`** | Not implemented — static mesh merging / batching is future work; use **`MODEL.CREATEINSTANCED`** (canonical) / **`MODEL.MAKEINSTANCED`** for GPU instancing today. |
| **`WORLD.SETREFLECTION(entity)`** | Not implemented — reflection probe capture / env map path not wired. |

---

## Repetition grid

| Command | Purpose |
|--------|---------|
| **`ENTITY.INSTANCEGRID(entity, countX, countZ, spacing)`** → **total** | Places **`countX * countZ`** copies on the **XZ** plane: the original entity moves to the first cell; additional cells use **`ENTITY.COPY`** (separate VRAM per copy). For **true** hardware instancing with one draw path, prefer **`MODEL.CREATEINSTANCED`**. |

---

## Stubs (errors explain next steps)

| Command | Status |
|--------|--------|
| **`LEVEL.APPLYPHYSICS(entity)`** | Not implemented — use **`BODY3D.*`** + **`PHYSICS3D.*`** manually ([PHYSICS3D.md](PHYSICS3D.md)). |
| **`LEVEL.SYNCLIGHTS(toggle, optional)`** | Not implemented — **`KHR_lights_punctual`** → **`LIGHT.*`** is future work. |
| **`PHYSICS.AUTOCREATE(entity)`** | Not implemented — use **`ENTITY.GETBOUNDS`** + **`BODY3D.ADDBOX`** / **`ADDMESH`**. |
| **`ENTITY.SETSTATIC(entity, toggle)`** | Marks an entity as static (for **`LEVEL.AUTOCOLLIDE`** or internal culling). |
| **`ENTITY.SETTRIGGER(entity)`** | Not implemented — sensors pending Jolt exposure. |
| **`ENTITY.INSTANCE`** | Not implemented — **`MODEL.CREATEINSTANCED`** / **`MODEL.MAKEINSTANCED`** or **`ENTITY.COPY`** / **`ENTITY.INSTANCEGRID`** (VRAM tradeoff). |

---

## Hot reload

| Command | Purpose |
|--------|---------|
| **`TEXTURE.RELOAD(texHandle)`** | Reloads GPU data from **`SourcePath`** (textures created with **`TEXTURE.LOAD`** from a file). Runs on the main thread. |

---

## Engine roadmap (compiler / host)

Resource **deduplication** on **`LEVEL.LOAD`**, optional **texture arrays / atlases**, **WASM shared memory** for scene-graph reads, and **automatic** script callbacks are described in [SCENE_ENGINE_BRIEF.md](SCENE_ENGINE_BRIEF.md).

---

## Full Example

```basic
WINDOW.OPEN(1280, 720, "Level Demo")
WINDOW.SETFPS(60)

cam = CAMERA.CREATE()
CAMERA.SETPOS(cam, 0, 10, -15)
CAMERA.SETTARGET(cam, 0, 0, 0)

LEVEL.SETROOT("assets/levels/")
lvl = LEVEL.LOAD("town.glb")

spawn = LEVEL.GETSPAWN("PlayerStart")
px = spawn(0) : py = spawn(1) : pz = spawn(2)

WHILE NOT WINDOW.SHOULDCLOSE()
    RENDER.CLEAR(80, 120, 180)
    RENDER.BEGIN3D(cam)
        ENTITY.DRAWALL()
    RENDER.END3D()
    RENDER.FRAME()
WEND

ENTITY.FREE(lvl)
CAMERA.FREE(cam)
WINDOW.CLOSE()
```

---

## See also

- [SCENE_ENGINE_BRIEF.md](SCENE_ENGINE_BRIEF.md) — architecture roadmap and WASM/Jolt notes  
- [SPRITE3D.md](SPRITE3D.md) — 3D billboards and animated sheets in levels  
- [ENTITY.md](ENTITY.md) — entity ids, groups, drawing  
- [ANIMATION_3D.md](ANIMATION_3D.md) — skinned models  
- [WORLD.md](WORLD.md) — global setup and streaming
- [VEHICLE.md](VEHICLE.md) — cars and aircraft
