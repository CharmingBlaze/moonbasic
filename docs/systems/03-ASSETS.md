# Asset systems: MESH, MODEL, MATERIAL, TEXTURE, ASSET

> Primitives, image files, materials, loaded models, and manifest-based asset packs.

**All commands:** [COMMAND_REGISTRY.md#assets](COMMAND_REGISTRY.md#assets)

**See also:** [01-CORE](01-CORE.md) · [reference/TEXTURE.md](../reference/TEXTURE.md) · [reference/MODEL.md](../reference/MODEL.md) · [reference/ASSET.md](../reference/ASSET.md)

**Case:** IDs in **`ASSET.*`** packs and command names are **case-insensitive**.

**Deep guides:** [guides/MESHES-MODELS-MATERIALS.md](guides/MESHES-MODELS-MATERIALS.md) · [guides/ASSETS-PIPELINE.md](guides/ASSETS-PIPELINE.md)

---

## Table of contents

- [MESH system](#mesh-system)
- [MODEL system](#model-system)
- [MATERIAL system](#material-system)
- [TEXTURE system](#texture-system)
- [ASSET system](#asset-system)
- [Full example](#full-example)
- [Memory notes](#memory-notes)
- [See also](#see-also)

---

## MESH system

CPU-side mesh data and built-in primitives.

### Core workflow

1. `MESH.CUBE` / `SPHERE` / `PLANE` / `CREATECYLINDER` — primitive handles.
2. Optional: edit vertices, then `MESH.UPLOAD(mesh)`.
3. Assign to entities with `ENTITY.SETMESH` (see [reference/ENTITY.md](../reference/ENTITY.md)).

---

### Primitive commands

| Command | Description |
|---------|-------------|
| `MESH.CUBE(w, h, d)` | Box mesh handle |
| `MESH.SPHERE(radius [, rings, slices])` | Sphere mesh |
| `MESH.PLANE(w, h)` | Flat quad |
| `MESH.CREATECYLINDER(radius, height [, segments])` | Cylinder |

**Returns:** `handle` — mesh

**Example:**

```basic
cubeMesh = MESH.CUBE(2, 2, 2)
```

---

### `MESH.UPLOAD(mesh)`

Uploads mesh vertex/index data to the GPU after editing.

| Argument | Type | Description |
|----------|------|-------------|
| mesh | handle | Mesh to upload |

**Returns:** nothing

**Example:**

```basic
MESH.UPLOAD(customMesh)
```

---

### Editable mesh (advanced)

| Command | Description |
|---------|-------------|
| `MESH.ADDVERTEX(mesh, x, y, z)` | Add vertex |
| `MESH.ADDFACE(mesh, i1, i2, i3)` | Triangle face |
| `MESH.RECALCNORMALS(mesh)` | Recompute normals |

See [reference/MESH.md](../reference/MESH.md) for the full mesh API.

---

## MODEL system

Loaded 3D models (GLB/GLTF, OBJ).

### `MODEL.LOAD(path)`

Loads a model file from disk.

| Argument | Type | Description |
|----------|------|-------------|
| path | string | File path (e.g. `assets/hero.glb`) |

**Returns:** `handle` — model

**Example:**

```basic
hero = MODEL.LOAD("assets/hero.glb")
ENTITY.SETMODEL(player, hero)
```

---

### `MODEL.FREE(model)` / `MODEL.UNLOAD(model)`

Releases GPU and registry resources for a model.

**Returns:** nothing

**Aliases:** checklist `MODEL.UNLOAD`

**Example:**

```basic
MODEL.FREE(hero)
```

---

### `MODEL.ANIMCOUNT(model)` / `MODEL.GETANIMNAME(model, index)`

Query skeletal animation clips on a loaded model.

**Returns:** `int` or `string`

**Example:**

```basic
n = MODEL.ANIMCOUNT(hero)
name = MODEL.GETANIMNAME(hero, 0)
```

Entity playback: `ENTITY.PLAYANIM` / `ENTITY.STOPANIM` — see [07-2D-WORLD](07-2D-WORLD.md#animation-system).

---

## MATERIAL system

Surface appearance (color, texture, PBR-lite values).

### `MATERIAL.CREATE(name)` / `MATERIAL.CREATEDEFAULT()`

Creates a material handle.

| Argument | Type | Description |
|----------|------|-------------|
| name | string | Debug name (optional overloads) |

**Returns:** `handle`

**Example:**

```basic
mat = MATERIAL.CREATEDEFAULT()
```

---

### `MATERIAL.SETCOLOR(mat, r, g, b [, a])`

Base color and optional alpha.

**Returns:** `handle`

**Example:**

```basic
MATERIAL.SETCOLOR(mat, 255, 255, 255, 255)
```

---

### `MATERIAL.SETTEXTURE(mat, texture)`

Assigns a texture handle to the material.

**Returns:** `handle`

**Example:**

```basic
MATERIAL.SETTEXTURE(mat, crateTex)
ENTITY.SETMATERIAL(crate, mat)
```

---

### Other material setters

| Command | Description |
|---------|-------------|
| `MATERIAL.SETALPHA(mat, a)` | Transparency |
| `MATERIAL.SETMETALLIC(mat, v)` | PBR metallic |
| `MATERIAL.SETROUGHNESS(mat, v)` | PBR roughness |

---

## TEXTURE system

Image loading and GPU textures.

### `TEXTURE.LOAD(path)`

Loads an image file (PNG, JPG, …).

| Argument | Type | Description |
|----------|------|-------------|
| path | string | Image path |

**Returns:** `handle`

**Example:**

```basic
tex = TEXTURE.LOAD("assets/crate.png")
```

---

### `TEXTURE.WIDTH(tex)` / `TEXTURE.HEIGHT(tex)`

Returns pixel dimensions.

**Returns:** `int`

**Example:**

```basic
w = TEXTURE.WIDTH(tex)
```

---

### `TEXTURE.FREE(tex)`

Frees the texture handle and GPU memory.

**Aliases:** checklist `TEXTURE.UNLOAD`

**Example:**

```basic
TEXTURE.FREE(tex)
```

---

## ASSET system

Manifest-driven loading — one JSON file lists textures, models, and sounds by id.

### `ASSET.LOADPACK(manifestPath)`

Parses `assets.json` (or similar) and prepares the pack. Paths are relative to the manifest file.

| Argument | Type | Description |
|----------|------|-------------|
| manifestPath | string | Path to manifest JSON |

**Returns:** nothing

**Example:**

```basic
ASSET.LOADPACK("assets/assets.json")
```

Example manifest:

```json
{
  "textures": { "player": "textures/player.png", "crate": "textures/crate.png" },
  "models": { "hero": "models/hero.glb" },
  "sounds": { "jump": "audio/jump.wav" }
}
```

---

### `ASSET.TEXTURE(id)` / `ASSET.MODEL(id)` / `ASSET.SOUND(id)`

Returns a cached handle for a manifest id. **Ids are case-insensitive.**

| Argument | Type | Description |
|----------|------|-------------|
| id | string | Key from the manifest |

**Returns:** `handle`

**Example:**

```basic
playerTex = ASSET.TEXTURE("player")
jumpSnd = ASSET.SOUND("jump")
```

---

### `ASSET.UNLOAD()`

Frees all handles loaded from the current pack (GPU/audio cache).

**Returns:** nothing

**Example:**

```basic
ASSET.UNLOAD()
```

---

## Full example

```basic
APP.OPEN(800, 600, "Assets")
APP.SETFPS(60)

ASSET.LOADPACK("assets/assets.json")
tex = ASSET.TEXTURE("crate")

mat = MATERIAL.CREATEDEFAULT()
MATERIAL.SETTEXTURE(mat, tex)

cam = CAMERA.CREATE()
CAMERA.SETACTIVE(cam)
CAMERA.SETPOS(cam, 0, 2, -6)
CAMERA.LOOKAT(cam, 0, 0, 0)

crate = ENTITY.CREATECUBE(2, 2, 2)
ENTITY.SETMATERIAL(crate, mat)

WHILE NOT APP.SHOULDCLOSE()
    RENDER.CLEAR(30, 32, 40)
    RENDER.BEGIN(cam)
    SCENE.DRAW()
    RENDER.END()
    RENDER.FRAME()
WEND

ASSET.UNLOAD()
APP.CLOSE()
```

---

## Memory notes

- **`TEXTURE.FREE`**, **`MODEL.FREE`**, **`MESH.FREE`** — per-handle cleanup.
- **`ASSET.LOADPACK`** caches handles; call **`ASSET.UNLOAD`** before loading a different pack or on level teardown.
- Reloading the same path without unload may replace cached entries — prefer explicit unload in long sessions.

---

## See also

- [07-2D-WORLD](07-2D-WORLD.md) — sprites and tilemaps use textures too
- [06-AUDIO](06-AUDIO.md) — `ASSET.SOUND` + `AUDIO.PLAY`
- [MEMORY.md](../MEMORY.md) — handle lifetime
