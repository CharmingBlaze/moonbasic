# Meshes, models, materials, and textures

> Build visible 3D objects — primitives, image textures, materials, and imported models — before or alongside asset packs.

**Namespaces:** `MESH` · `MODEL` · `MATERIAL` · `TEXTURE` · **Status:** Shipped

**Commands:** [COMMAND_REGISTRY.md#assets](../COMMAND_REGISTRY.md#assets) · **Overview:** [03-ASSETS.md](../03-ASSETS.md)

**Related:** Packaged shipping → [ASSETS-PIPELINE.md](ASSETS-PIPELINE.md) ( `ASSET.LOADPACK` )

---

## Table of contents

- [At a glance](#at-a-glance)
- [When to use this system](#when-to-use-this-system)
- [Choose the right tool](#choose-the-right-tool)
- [Core workflow](#core-workflow)
- [TEXTURE and MATERIAL](#texture-and-material)
- [MESH primitives](#mesh-primitives)
- [MODEL loading](#model-loading)
- [Full example](#full-example)
- [Common mistakes](#common-mistakes)
- [Memory notes](#memory-notes)
- [See also](#see-also)

---

## At a glance

| Idea | Detail |
|------|--------|
| **You get** | GPU textures, material slots, primitive meshes, animated model files |
| **You need first** | Entities to display them ([ENTITY-SYSTEM.md](ENTITY-SYSTEM.md)) |
| **Typical games** | 3D props, characters, styled cubes |
| **Not for** | Manifest-only shipping — use `ASSET.*` when ids are stable |

**Why four namespaces:** **Texture** = image on GPU. **Material** = how surface reacts to light. **Mesh** = geometry. **Model** = file with meshes + optional skeleton.

---

## When to use this system

**Use when:**

- Prototyping with `MESH.CUBE` before art exists.
- Loading `character.glb` with `MODEL.LOAD`.
- Custom colors and textures per entity (`MATERIAL.SETCOLOR`, `SETTEXTURE`).
- Editing procedural geometry then `MESH.UPLOAD`.

**Skip when:**

- All art is in `assets.json` with stable ids — prefer [ASSETS-PIPELINE.md](ASSETS-PIPELINE.md).
- Pure 2D — `SPRITE.LOAD` ([SPRITES-TILEMAPS-2D.md](SPRITES-TILEMAPS-2D.md)).

---

## Choose the right tool

| I want to… | Use | Not |
|------------|-----|-----|
| Gray box level | `ENTITY.CREATECUBE` | `MESH.CUBE` + manual entity (unless custom mesh) |
| Custom mesh shape | `MESH.*` + `ENTITY.SETMESH` | Many entities for one merged mesh |
| Character from file | `MODEL.LOAD` + `ENTITY.SETMODEL` | `TEXTURE.LOAD` alone |
| Wood / metal look | `MATERIAL.CREATE` + `SETTEXTURE` | Tint only on untextured material |
| Ship game with ids | `ASSET.LOADPACK` | Hard-coded paths everywhere |
| Play skeletal clip | `ENTITY.PLAYANIM` | [ANIMATION.md](ANIMATION.md) |

---

## Core workflow

1. **Texture** (optional) — `TEXTURE.LOAD("wall.png")`.  
   **Why:** Image data for material slots.

2. **Material** — `MATERIAL.CREATE()` → `SETCOLOR` / `SETTEXTURE`.  
   **Why:** Binds shading + texture for draw calls.

3. **Geometry** — `MESH.CUBE` or `MODEL.LOAD("hero.glb")`.  
   **Why:** Defines vertices; model includes animation data.

4. **Entity** — `ENTITY.CREATE` → `SETMESH` or `SETMODEL` → assign material.  
   **Why:** Scene graph position + `SCENE.DRAW`.

5. **Free** — `TEXTURE.FREE`, `MODEL.FREE` on level unload.

---

## TEXTURE and MATERIAL

```basic
tex = TEXTURE.LOAD("assets/brick.png")
mat = MATERIAL.CREATE()
MATERIAL.SETTEXTURE(mat, tex)
MATERIAL.SETCOLOR(mat, 200, 200, 200)

wall = ENTITY.CREATECUBE(4, 3, 0.2)
ENTITY.SETMATERIAL(wall, mat)
```

| Command | Why |
|---------|-----|
| `TEXTURE.LOAD(path)` | Load image to GPU |
| `TEXTURE.WIDTH` / `HEIGHT` | UI scaling, terrain splat sizing |
| `TEXTURE.FREE(tex)` | Free VRAM |
| `MATERIAL.CREATE()` | New material slot |
| `MATERIAL.SETCOLOR(mat, r, g, b)` | Flat tint multiplier |
| `MATERIAL.SETTEXTURE(mat, tex)` | Albedo / diffuse map |

---

## MESH primitives

| Command | Why |
|---------|-----|
| `MESH.CUBE(w, h, d)` | Boxes, walls |
| `MESH.SPHERE(radius)` | Balls, planets |
| `MESH.PLANE(w, h)` | Floors, banners |
| `MESH.CREATECYLINDER(r, h)` | Pipes, columns |
| `MESH.UPLOAD(mesh)` | Push edited vertices to GPU |

Assign: `ENTITY.SETMESH(entity, mesh)` after `ENTITY.CREATE`.

---

## MODEL loading

```basic
heroModel = MODEL.LOAD("assets/hero.glb")
hero = ENTITY.CREATE("Hero")
ENTITY.SETMODEL(hero, heroModel)
hero.pos(0, 0, 5)
```

| Command | Why |
|---------|-----|
| `MODEL.LOAD(path)` | Import GLB/OBJ/etc. |
| `MODEL.ANIMCOUNT` / `GETANIMNAME` | List clips for [ANIMATION.md](ANIMATION.md) |
| `MODEL.FREE(model)` | Release mesh + skeleton VRAM |

---

## Full example

```basic
APP.OPEN(800, 600, "Materials")
APP.SETFPS(60)

cam = CAMERA.CREATE()
CAMERA.SETACTIVE(cam)
CAMERA.SETPOS(cam, 0, 2, -6)
CAMERA.LOOKAT(cam, 0, 0, 0)

sun = LIGHT.CREATEDIRECTIONAL()
LIGHT.SETDIR(sun, -0.3, -1, -0.2)

tex = TEXTURE.LOAD("assets/checker.png")
mat = MATERIAL.CREATE()
MATERIAL.SETTEXTURE(mat, tex)

box = ENTITY.CREATECUBE(2, 2, 2)
ENTITY.SETMATERIAL(box, mat)
box.pos(0, 0, 3)

WHILE NOT APP.SHOULDCLOSE()
    box.turn(0, 45 * APP.DELTA(), 0)
    RENDER.CLEAR(25, 28, 35)
    RENDER.BEGIN(cam)
    SCENE.DRAW()
    RENDER.END()
    RENDER.FRAME()
WEND

TEXTURE.FREE(tex)
APP.CLOSE()
```

---

## Common mistakes

| Mistake | Fix |
|---------|-----|
| Pink / missing texture | Path wrong; load pack first or use `FILE.EXISTS` |
| Flat gray mesh | Add `LIGHT.*` ([LIGHTING.md](LIGHTING.md)) |
| Model not visible | `SETMODEL` on entity; camera `LOOKAT` entity |
| VRAM growth | `FREE` textures/models on level change |
| Duplicate loads | Cache handle; or `ASSET.MODEL("id")` |

---

## Memory notes

- `TEXTURE.FREE`, `MODEL.FREE` — per-handle GPU release.
- `ASSET.UNLOAD` — bulk release from pack ([ASSETS-PIPELINE.md](ASSETS-PIPELINE.md)).

---

## See also

- [ASSETS-PIPELINE.md](ASSETS-PIPELINE.md) — `assets.json` manifest
- [LIGHTING.md](LIGHTING.md) — materials need light
- [ANIMATION.md](ANIMATION.md) — clips on loaded models
- [ENTITY-SYSTEM.md](ENTITY-SYSTEM.md) — attach mesh/model to entity
