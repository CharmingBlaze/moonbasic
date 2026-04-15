# Sprite — `Sprite.*`, groups, layers, batch, UI, `Particle2D`, `ANIM.*`, atlas

Sprites are **GPU textures** (Raylib `Texture2D`) plus **source rectangle**, **frame layout**, and optional **`ANIM.*`** state. Drawing uses **`SPRITE.DRAW`**. When **`RENDER.BEGINMODE2D`** / **`RENDER.ENDMODE2D`** are implemented in your build, wrap draws for **Camera2D**-style views; otherwise draw directly after **`Render.Clear`** (as in **`testdata/sprite_complete_test.mb`**).

**Requires CGO** (same as `Texture.*`, `Draw.*`).

Registry keys use **dots and uppercase** (e.g. `SPRITE.LOAD`). This document uses **PascalCase** names aligned with specs where helpful.

**Related:** [ATLAS.md](ATLAS.md) (`ATLAS.LOAD`, `ATLAS.GETSPRITE`, `ATLAS.FREE`), [TEXTURE.md](TEXTURE.md), [IMAGE.md](IMAGE.md).

**Blitz-style “sprite collide”:** there is no separate **`Sprite.Collide`** name — use **`SPRITE.HIT`** / **`SPRITE.POINTHIT`**. Overlap uses the same **scaled** destination quad, **origin**, and **rotation** as **`SPRITE.DRAW`** (raylib **`DrawTexturePro`**), not a separate axis-aligned box. For pixel-perfect work, use **`IMAGE.*`** CPU pixels or custom overlap — see table below and [BLITZ_ESSENTIAL_API.md](BLITZ_ESSENTIAL_API.md).

---

### `Sprite.Load(path)`
Loads an image and returns a **sprite handle**.

### `Sprite.Free(handle)`
Unloads the sprite and frees memory.

---

### `Sprite.Draw(handle, x, y)`
Draws the current frame at pixel coordinates.

### `Sprite.SetPos(handle, x, y)`
Sets a floating-point draw offset.

---

### `Sprite.DefAnim(handle, count)`
Defines a grid animation (count is a string).

### `Sprite.UpdateAnim(handle, dt)`
Advances animation frame by time.

---

### `Sprite.Hit(a, b)`
Returns `TRUE` if the two drawn quads overlap (SAT on the four corners; matches **`DrawTexturePro`** geometry).

### `Sprite.PointHit(handle, x, y)`
Returns `TRUE` if **`(x, y)`** lies inside that quad in the same coordinate space as **`Sprite.Draw`**’s **`x, y`** plus **`SetPos`** offsets (inverse rotation into local frame size).

---

## Sprite Groups

### `SpriteGroup.Make()`

Creates a new empty sprite group. Returns a handle.

### `SpriteGroup.Add(group, sprite)`

Adds a sprite to the group.

### `SpriteGroup.Draw(group, x, y)`

Draws all sprites in the group relative to a base position.

### `SpriteGroup.Free(group)`

Frees the group object (members remain).

---

## SpriteLayer.*

| Command | Signature | Notes |
|---|---|---|
| **`SPRITELAYER.CREATE`** / deprecated `SPRITELAYER.MAKE` | `(z)` → handle | **`z`** stored for your sorting; draw order is under your control. |
| `SPRITELAYER.ADD` | `(layer, spr)` | |
| `SPRITELAYER.CLEAR` | `(layer)` | Remove all members. |
| `SPRITELAYER.SETZ` | `(layer, z)` | Update stored **z**. |
| `SPRITELAYER.DRAW` | `(layer, x, y)` | Same base position semantics as group draw. |
| `SPRITELAYER.FREE` | `(layer)` | Frees layer only. |

---

## SpriteBatch.*

Records **multiple** **`(sprite, x, y)`** draws; **`SPRITEBATCH.DRAW`** executes them in order.

| Command | Notes |
|---|---|
| **`SPRITEBATCH.CREATE`** / deprecated `SPRITEBATCH.MAKE` | `()` → handle |
| `SPRITEBATCH.ADD` | `(batch, spr, x, y)` — **int** positions |
| `SPRITEBATCH.CLEAR` | `(batch)` |
| `SPRITEBATCH.DRAW` | `(batch)` |
| `SPRITEBATCH.FREE` | `(batch)` |

---

## SpriteUI.*

Anchored placement using **fractions of screen size** (e.g. **`0.5, 0.5`** = center).

```basic
ui = SPRITEUI.CREATE(spr, anchorX, anchorY)
SPRITEUI.DRAW(ui, SCREENW(), SCREENH())
SPRITEUI.FREE(ui)
```

**`SPRITEUI.FREE`** releases only the **UI wrapper**; the sprite remains.

---

## Particle2D.* (simple filled circles)

CPU-side **circles** (no texture). **`PARTICLE2D.CREATE(max, r, g, b, a)`** sets pool size and colour; **`EMIT`** adds particles; **`UPDATE`** integrates velocity and **`life`**; **`DRAW`** renders.

| Command | Arguments |
|---|---|
| `PARTICLE2D.MAKE` | `(max, r, g, b, a)` |
| `PARTICLE2D.EMIT` | `(p, x, y, vx, vy, life)` |
| `PARTICLE2D.UPDATE` | `(p, dt)` |
| `PARTICLE2D.DRAW` | `(p)` |
| `PARTICLE2D.FREE` | `(p)` |

---

## ANIM.* (optional state machine)

| Command | Purpose |
|---|---|
| `ANIM.DEFINE` | Named clip: first/last frame, fps, looping |
| `ANIM.ADDTRANSITION` | Conditional clip change |
| `ANIM.UPDATE` | Advance + evaluate transitions |
| `ANIM.SETPARAM` | Parameters for transition conditions |

See inline tables in earlier revisions of this file for **transition condition** syntax. Do not mix **`ANIM.UPDATE`** with **`SPRITE.UPDATEANIM`** strip advancement on the **same** sprite without understanding the interaction.

---

## Atlas

See **[ATLAS.md](ATLAS.md)** for **`ATLAS.LOAD`**, **`ATLAS.GETSPRITE`**, **`ATLAS.FREE`**.

---

## Example (strip + Mode2D)

```basic
Window.Open(800, 600, "Sprite strip")
Window.SetFPS(60)

hero = SPRITE.LOAD("sheet.png")
SPRITE.DEFANIM(hero, "4")
SPRITE.PLAYANIM(hero, "walk")

x = 300
y = 250

WHILE NOT Window.ShouldClose()
    SPRITE.SETPOS(hero, x, y)
    SPRITE.UPDATEANIM(hero, TIME.DELTA())

    Render.Clear(30, 40, 50)
    SPRITE.DRAW(hero, 0, 0)
    Render.Frame()
WEND

SPRITE.FREE(hero)
Window.Close()
```

---

## Common mistakes

- **Skipping `BeginMode2D`** when using cameras or scaled views — align with your **`Camera2D`** setup.
- **Leaking sprites** — pair **`SPRITE.LOAD`** / **`ATLAS.GETSPRITE`** with **`SPRITE.FREE`** when done (and **`ATLAS.FREE`** for the atlas).
- **Atlas sprites** — **`SPRITE.FREE`** on an atlas sub-sprite does **not** unload the shared atlas texture (`fromAtlas` path).

---

## See also

- [ATLAS.md](ATLAS.md) — packed sheets + JSON
- [TEXTURE.md](TEXTURE.md) — raw GPU textures
- [IMAGE.md](IMAGE.md) — CPU images before upload
- [DRAW2D.md](DRAW2D.md) — screen drawing helpers
