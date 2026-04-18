# 3D Sprite (Billboard) Commands

Camera-facing **billboard** quads in world space: static or **texture-sheet animated**, placed either as **entities** (scene graph, `ENTITY.DRAWALL`) or drawn **immediately** with **`DRAW3D.BILLBOARD`**. This is **not** the same as screen-space [**`SPRITE.*`**](SPRITE.md) (2D).

For **skeletal** animation on loaded meshes, see [**`ANIMATION_3D.md`**](ANIMATION_3D.md) (`ENTITY.LOAD`, bones, clips).

Page shape follows [**`DOC_STYLE_GUIDE.md`**](DOC_STYLE_GUIDE.md) (**WAVE pattern**).

## Core Workflow

### Static billboards

1. **Entity path (recommended for props and levels):** load a texture from disk with **`ENTITY.LOADSPRITE(path)`** or bind a **`TEXTURE.*`** handle with **`ENTITY.CREATESPRITE(tex, w, h [, parent])`**. Position with **`ENTITY.SETPOS`**, scale with **`ENTITY.SETSCALE`**, optional **`SPRITEMODE`** / **`SPRITEVIEWMODE`** for how the quad tracks the camera.
2. **Immediate path (one-off / debug):** inside **`RENDER.BEGIN3D(cam)`** … **`RENDER.END3D()`** (or **`CAMERA.BEGIN`/`CAMERA.END`**), call **`DRAW3D.BILLBOARD`** or **`DRAW3D.BILLBOARDREC`** with a **texture handle** and world-space coordinates.
3. Each frame: **`ENTITY.UPDATE(dt)`** for sprite entities, then draw (**`ENTITY.DRAWALL`** or your immediate **`DRAW3D.*`** calls).

### Animated billboards (texture atlas / grid)

1. Load a **grid-based** sheet with **`TEXTURE.LOADANIM(path, cols, rows)`** (or **`TEXTURE.LOAD`** + **`TEXTURE.SETGRID`** — see [**`TEXTURE.md`**](TEXTURE.md)).
2. Start playback with **`TEXTURE.PLAY(tex, fps, loop)`** (or drive frames manually with **`TEXTURE.SETFRAME`** / **`ENTITY.SETSPRITEFRAME`**).
3. **Once per frame**, call **`TEXTURE.TICKALL`** (optionally with **`dt`**) so playing textures advance.
4. Attach the same **texture handle** to a billboard entity: **`ENTITY.CREATESPRITE(tex, quadW, quadH)`**. The renderer uses the texture’s current **source rectangle** when drawing (`FrameSourceRect` on the billboard path in **`runtime/mbentity/`**).

### Not covered here

- **Skeletal clips** on **`MODEL`/`ENTITY.LOAD`** meshes — [**`ANIMATION_3D.md`**](ANIMATION_3D.md).
- **GPU particle billboards** — [**`PARTICLE3D.md`**](PARTICLE3D.md).

---

### `SPRITEMODE` / `SPRITEVIEWMODE` 

Sets how a **sprite entity** rotates relative to the camera.

- **Arguments**:
  - `sprite`: Entity id (sprite entity from **`LOADSPRITE`** / **`CREATESPRITE`**).
  - `mode` (int): **`1`** = Y-axis billboard, **`2`** = full camera-facing billboard, **`3`** = static quad (see runtime: **`entity_blitz_cgo.go`**, `spriteMode`).
- **Returns**: (none)

---

### `ENTITY.LOADSPRITE(path)` / `ENTITY.LOADSPRITE(path, parentEntity)` 

Loads an image from **path** and returns a **sprite entity id**. Optional **parent** parents the new entity.

- **Returns**: (int) Entity id.

---

### `ENTITY.CREATESPRITE(path)` / `ENTITY.CREATESPRITE(path, parent)` / `ENTITY.CREATESPRITE(tex, w, h)` / `ENTITY.CREATESPRITE(tex, w, h, parent)` 

Creates a billboard entity: either from a **file path** or from an existing **`TEXTURE.*`** handle plus **quad width and height** in world units (atlas / **`TEXTURE.LOADANIM`** friendly).

- **Returns**: (int) Entity id.

---

### `ENTITY.SETSPRITEFRAME(entity, frameIndex)` 

Sets the **atlas frame index** (0-based) for a billboard bound to a **`TEXTURE`** object.

- **Arguments**:
  - `entity` (int): Sprite entity id.
  - `frameIndex` (int): Cell index.

---

### `DRAW3D.BILLBOARD(tex, x, y, z, size, r, g, b, a)` 

Draws a **camera-facing** billboard using the **full texture**. Must run inside **`RENDER.BEGIN3D(cam)`** … **`RENDER.END3D()`** (or **`CAMERA.BEGIN`/`END`**).

- **Arguments**:
  - `tex` (handle): Texture handle.
  - `x, y, z` (float): World position.
  - `size` (float): Uniform scale.
  - `r, g, b, a` (int): Tint 0–255.
- **Returns**: (none)

**Alias:** **`DRAW.BILLBOARD`** (same arguments).

---

### `DRAW3D.BILLBOARDREC(tex, srcX, srcY, srcW, srcH, x, y, z, w, h, r, g, b, a)` 

Draws a billboard using a **source rectangle** from the texture (sub-image / atlas region).

- **Arguments**:
  - `tex` (handle): Texture handle.
  - `srcX, srcY, srcW, srcH` (float): Source rectangle in texels.
  - `x, y, z` (float): World position.
  - `w, h` (float): World size of the quad.
  - `r, g, b, a` (int): Tint 0–255.
- **Returns**: (none)

**Alias:** **`DRAW.BILLBOARDREC`**.

---

### `TEXTURE.LOADANIM(path, cols, rows)` 

Loads a texture and configures a **grid** of frames (see [**`TEXTURE.md`**](TEXTURE.md)). Use with **`TEXTURE.PLAY`** and **`TEXTURE.TICKALL`** for automatic frame cycling.

- **Returns**: (handle) Texture handle.

---

### `TEXTURE.PLAY(tex, fps, loop)` 

Enables automatic frame advance for a texture; call **`TEXTURE.TICKALL`** each frame.

- **Arguments**:
  - `tex` (handle)
  - `fps` (float)
  - `loop` (bool)

---

### `TEXTURE.TICKALL()` / `TEXTURE.TICKALL(dt)` 

Advances all **playing** texture animations. Call **once per frame** from your main loop.

---

### `TEXTURE.SETFRAME(tex, frameIndex)` 

Selects a **grid cell** (0-based) without using **`TEXTURE.PLAY`**.

---

## Atlas and maps

### Two “atlas” meanings

| Workflow | Use when | Typical APIs |
|----------|-----------|--------------|
| **Texture grid on `TEXTURE.*`** | 3D billboard entities, **`DRAW3D.BILLBOARDREC`**, shared sheets | **`TEXTURE.LOADANIM`**, **`TEXTURE.SETFRAME`**, **`TEXTURE.PLAY`**, **`ENTITY.CREATESPRITE(tex, w, h)`** |
| **TexturePacker JSON + `ATLAS.*`** | 2D sprites from named regions | [**`ATLAS.md`**](ATLAS.md) — **`ATLAS.LOAD`**, **`ATLAS.GETSPRITE`** → **`SPRITE.DRAW`** (screen space) |

For **world** billboards, prefer loading the same PNG once as a **`TEXTURE`** (grid or manual UVs) rather than **`ATLAS.GETSPRITE`** handles, which target **`SPRITE.*`** drawing.

### Building 3D maps

- **Terrain:** heightfields and placement helpers — [**`TERRAIN.md`**](TERRAIN.md). Snap or parent billboard **entities** to the ground.
- **Levels / scenes:** packaged worlds and switching — [**`LEVEL.md`**](LEVEL.md), [**`SCENE.md`**](SCENE.md).
- **Props:** billboard **entities** for signs, torches, stylized trees; parent to groups or terrain-aligned empties.
- **Vegetation / scatter:** **`WORLD.SETVEGETATION(terrain, billboard, density)`** reserves a billboard handle for future instanced drawing; see [**`WORLD.md`**](WORLD.md) for current behavior.
- **Visibility / LOD:** frustum and distance tests — [**`CULL.md`**](CULL.md), **`Entity.InView`**. Many billboard entities may need profiling; see performance note in [**`ENTITY.md`**](ENTITY.md).

---

## Full Example

**Static entity billboard** and **animated texture** (grid) in one loop pattern:

```basic
WINDOW.OPEN(960, 540, "3D sprite demo")
WINDOW.SETFPS(60)

cam = CAMERA.CREATE()
CAMERA.SETPOS(cam, 0, 4, 10)
CAMERA.SETTARGET(cam, 0, 1, 0)

; Optional: animated sheet (replace path with your asset)
tex = TEXTURE.LOADANIM("sprites/hero_sheet.png", 4, 4)
TEXTURE.PLAY(tex, 12.0, TRUE)

; Billboard entity using the same texture handle (world units for quad size)
hero = ENTITY.CREATESPRITE(tex, 2.0, 2.0)
ENTITY.SETPOS(hero, 0, 1, 0, TRUE)
SPRITEMODE(hero, 2)

WHILE NOT WINDOW.SHOULDCLOSE()
    dt = TIME.DELTA()
    TEXTURE.TICKALL(dt)
    ENTITY.UPDATE(dt)

    RENDER.CLEAR(20, 25, 35)
    RENDER.BEGIN3D(cam)
        ENTITY.DRAWALL()
        DRAW3D.GRID(20, 1.0)
    RENDER.END3D()
    RENDER.FRAME()
WEND

ENTITY.FREE(hero)
TEXTURE.FREE(tex)
CAMERA.FREE(cam)
WINDOW.CLOSE()
```

**Immediate draw** (no entity) for a single quad:

```basic
; Inside RENDER.BEGIN3D(cam) … RENDER.END3D()
; DRAW3D.BILLBOARD(myTex, x, y, z, 2.0, 255, 255, 255, 255)
```

---

## See also

- [**`SPRITE.md`**](SPRITE.md) — 2D screen sprites
- [**`ENTITY.md`**](ENTITY.md) — entity graph, **`ENTITY.DRAWALL`**, transforms
- [**`DRAW3D.md`**](DRAW3D.md) — primitives and billboard aliases
- [**`TEXTURE.md`**](TEXTURE.md) — grids, **`LOADANIM`**, **`TICKALL`**
- [**`ATLAS.md`**](ATLAS.md) — TexturePacker atlases for **`SPRITE.*`**
- [**`TERRAIN.md`**](TERRAIN.md), [**`LEVEL.md`**](LEVEL.md), [**`SCENE.md`**](SCENE.md) — worlds and maps
- [**`CULL.md`**](CULL.md) — visibility helpers
- [**`ANIMATION_3D.md`**](ANIMATION_3D.md) — skeletal mesh animation (not billboards)
- [**`PARTICLE3D.md`**](PARTICLE3D.md) — 3D particle billboards
- [**`CAMERA.md`**](CAMERA.md) — 3D camera and picking rays
