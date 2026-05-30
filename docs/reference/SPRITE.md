# Sprite Commands

2D sprites with frame animation, groups, layers, batching, UI sprites, and 2D particles.

For **camera-facing 3D billboards** (world-space props, texture-sheet animation on entities), see [**`SPRITE3D.md`**](SPRITE3D.md).

Page shape follows [DOC_STYLE_GUIDE.md](../DOC_STYLE_GUIDE.md) (**WAVE pattern**).

## Core Workflow

1. Load with `SPRITE.LOAD` or create from a texture with `SPRITE.CREATE`.
2. Configure frame layout with `SPRITE.SETGRID`, animate with `ANIM.PLAY`.
3. Draw with `SPRITE.DRAW` each frame.
4. Test overlap with `SPRITE.HIT` / `SPRITE.POINTHIT`.
5. Free with `SPRITE.FREE`.

See also [ATLAS.md](ATLAS.md), [TEXTURE.md](TEXTURE.md), [IMAGE.md](IMAGE.md).

---

### `SPRITE.LOAD(path)`
Loads an image and returns a **sprite handle**.

- **Arguments**:
    - `path`: (String) File path relative to working directory.
- **Returns**: (Handle) The new sprite handle.
- **Example**:
    ```basic
    hero = SPRITE.LOAD("hero.png")
    ```

---

### `SPRITE.FREE(handle)`
Unloads the sprite and frees memory.

---

### `SPRITE.DRAW(handle, x, y)`
Draws the current frame at pixel coordinates.

- **Arguments**:
    - `handle`: (Handle) The sprite to draw.
    - `x, y`: (Float) Screen position.
- **Returns**: (Handle) The sprite handle (for chaining).

---

### `SPRITE.SETPOS(handle, x, y)`
Sets a floating-point draw offset.

- **Returns**: (Handle) The sprite handle (for chaining).

---

### `SPRITE.DEFANIM(handle, framesCount)`
Defines a grid animation.

- **Arguments**:
    - `handle`: (Handle) The sprite to animate.
    - `framesCount`: (Integer) Total frames in the sprite sheet.
- **Returns**: (Handle) The sprite handle (for chaining).

---

### `SPRITE.UPDATEANIM(handle, dt)`
Advances animation frame by time.

- **Returns**: (Handle) The sprite handle (for chaining).

---

### `SPRITE.HIT(a, b)`
Returns **`TRUE`** if the two sprites overlap.

- **Arguments**:
    - `a, b`: (Handle) The sprites to test.
- **Returns**: (Boolean)

---

### `SPRITE.POINTHIT(handle, x, y)`
Returns **`TRUE`** if **`(x, y)`** lies inside the sprite.

- **Returns**: (Boolean)

---

## `SPRITEGROUP.*`

### `SPRITEGROUP.CREATE()` 
Creates a new empty sprite group. Returns a handle. **`SPRITEGROUP.MAKE`** is a **deprecated** alias of **`SPRITEGROUP.CREATE`**.

---

### `SPRITEGROUP.ADD(group, sprite)` 
Adds a sprite to the group.

---

### `SPRITEGROUP.DRAW(group, x, y)` 
Draws all sprites in the group relative to a base position.

---

### `SPRITEGROUP.FREE(group)` 
Frees the group object (members remain).

---

## `SPRITELAYER.*`

| Command | Signature | Notes |
|--------|-----------|--------|
| **`SPRITELAYER.CREATE`** / deprecated **`SPRITELAYER.MAKE`** | `(z)` → handle | **`z`** stored for your sorting; draw order is under your control. |
| **`SPRITELAYER.ADD`** | `(layer, spr)` | |
| **`SPRITELAYER.CLEAR`** | `(layer)` | Remove all members. |
| **`SPRITELAYER.SETZ`** | `(layer, z)` | Update stored **z**. |
| **`SPRITELAYER.DRAW`** | `(layer, x, y)` | Same base position semantics as group draw. |
| **`SPRITELAYER.FREE`** | `(layer)` | Frees layer only. |

---

## `SPRITEBATCH.*`

Records **multiple** **`(sprite, x, y)`** draws; **`SPRITEBATCH.DRAW`** executes them in order.

| Command | Notes |
|--------|--------|
| **`SPRITEBATCH.CREATE`** / deprecated **`SPRITEBATCH.MAKE`** | `()` → handle |
| **`SPRITEBATCH.ADD`** | `(batch, spr, x, y)` — **int** positions |
| **`SPRITEBATCH.CLEAR`** | `(batch)` |
| **`SPRITEBATCH.DRAW`** | `(batch)` |
| **`SPRITEBATCH.FREE`** | `(batch)` |

---

## `SPRITEUI.*`

Anchored placement using **fractions of screen size** (e.g. **`0.5, 0.5`** = center).

```basic
ui = SPRITEUI.CREATE(spr, anchorX, anchorY)
SPRITEUI.DRAW(ui, SCREENW(), SCREENH())
SPRITEUI.FREE(ui)
```

**`SPRITEUI.CREATE`** — **`SPRITEUI.MAKE`** is a **deprecated** alias. **`SPRITEUI.FREE`** releases only the **UI wrapper**; the sprite remains.

---

## `PARTICLE2D.*` (simple filled circles)

CPU-side **circles** (no texture). **`PARTICLE2D.CREATE(max, r, g, b, a)`** sets pool size and colour; **`EMIT`** adds particles; **`UPDATE`** integrates velocity and **`life`**; **`DRAW`** renders. **`PARTICLE2D.MAKE`** is a **deprecated** alias of **`PARTICLE2D.CREATE`**.

| Command | Arguments |
|--------|-----------|
| **`PARTICLE2D.CREATE`** | `(max, r, g, b, a)` |
| **`PARTICLE2D.EMIT`** | `(p, x, y, vx, vy, life)` |
| **`PARTICLE2D.UPDATE`** | `(p, dt)` |
| **`PARTICLE2D.DRAW`** | `(p)` |
| **`PARTICLE2D.FREE`** | `(p)` |

---

## `ANIM.*` (optional state machine)

| Command | Purpose |
|--------|---------|
| **`ANIM.DEFINE`** | Named clip: first/last frame, fps, looping |
| **`ANIM.ADDTRANSITION`** | Conditional clip change |
| **`ANIM.UPDATE`** | Advance + evaluate transitions |
| **`ANIM.SETPARAM`** | Parameters for transition conditions |

See inline tables in earlier revisions of this file for **transition condition** syntax. Do not mix **`ANIM.UPDATE`** with **`SPRITE.UPDATEANIM`** strip advancement on the **same** sprite without understanding the interaction.

---

## Atlas

See **[ATLAS.md](ATLAS.md)** for **`ATLAS.LOAD`**, **`ATLAS.GETSPRITE`**, **`ATLAS.FREE`**.

---

## Full Example (strip + Mode2D)

```basic
WINDOW.OPEN(800, 600, "Sprite strip")
WINDOW.SETFPS(60)

hero = SPRITE.LOAD("sheet.png")
SPRITE.DEFANIM(hero, "4")
SPRITE.PLAYANIM(hero, "walk")

x = 300
y = 250

WHILE NOT WINDOW.SHOULDCLOSE()
    SPRITE.SETPOS(hero, x, y)
    SPRITE.UPDATEANIM(hero, TIME.DELTA())

    RENDER.CLEAR(30, 40, 50)
    SPRITE.DRAW(hero, 0, 0)
    RENDER.FRAME()
WEND

SPRITE.FREE(hero)
WINDOW.CLOSE()
```

---

## Common mistakes

- **Skipping `RENDER.BEGINMODE2D` / `RENDER.ENDMODE2D`** when using cameras or scaled views — align with your 2D camera setup.
- **Leaking sprites** — pair **`SPRITE.LOAD`** / **`ATLAS.GETSPRITE`** with **`SPRITE.FREE`** when done (and **`ATLAS.FREE`** for the atlas).
- **Atlas sprites** — **`SPRITE.FREE`** on an atlas sub-sprite does **not** unload the shared atlas texture (`fromAtlas` path).

---

## Extended Command Reference

### Transform

| Command | Description |
|--------|-------------|
| `SPRITE.SETPOSITION(spr, x, y)` | Alias of `SPRITE.SETPOS`. |
| `SPRITE.GETPOS(spr)` | Returns `[x, y]` position. |
| `SPRITE.ROT(spr, angle)` | Set rotation angle in degrees. |
| `SPRITE.SETROT(spr, angle)` | Alias of `SPRITE.ROT`. |
| `SPRITE.GETROT(spr)` | Returns current rotation angle. |
| `SPRITE.SETSCALE(spr, sx, sy)` | Set XY scale. |
| `SPRITE.GETSCALE(spr)` | Returns `[sx, sy]` scale. |
| `SPRITE.SETORIGIN(spr, ox, oy)` | Set pivot/origin offset in pixels. |

### Color & alpha

| Command | Description |
|--------|-------------|
| `SPRITE.ALPHA(spr, a)` | Set alpha 0.0–1.0. |
| `SPRITE.SETALPHA(spr, a)` | Alias of `SPRITE.ALPHA`. |
| `SPRITE.GETALPHA(spr)` | Returns current alpha. |
| `SPRITE.COLOR(spr, r,g,b,a)` | Set RGBA tint. |
| `SPRITE.SETCOLOR(spr, r,g,b,a)` | Alias of `SPRITE.COLOR`. |
| `SPRITE.GETCOLOR(spr)` | Returns `[r,g,b,a]` tint. |

### Animation

| Command | Description |
|--------|-------------|
| `SPRITE.SETFRAME(spr, frame)` | Set current animation frame index. |

---

## See also

- [ATLAS.md](ATLAS.md) — packed sheets + JSON
- [TEXTURE.md](TEXTURE.md) — raw GPU textures
- [IMAGE.md](IMAGE.md) — CPU images before upload
- [DRAW2D.md](DRAW2D.md) — screen drawing helpers
