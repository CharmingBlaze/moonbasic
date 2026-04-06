# Sprite — `Sprite.*`, groups, layers, batch, UI, `Particle2D`, `ANIM.*`, atlas

Sprites are **GPU textures** (Raylib `Texture2D`) plus **source rectangle**, **frame layout**, and optional **`ANIM.*`** state. Drawing uses **`SPRITE.DRAW`**. When **`RENDER.BEGINMODE2D`** / **`RENDER.ENDMODE2D`** are implemented in your build, wrap draws for **Camera2D**-style views; otherwise draw directly after **`Render.Clear`** (as in **`testdata/sprite_complete_test.mb`**).

**Requires CGO** (same as `Texture.*`, `Draw.*`).

Registry keys use **dots and uppercase** (e.g. `SPRITE.LOAD`). This document uses **PascalCase** names aligned with specs where helpful.

**Related:** [ATLAS.md](ATLAS.md) (`ATLAS.LOAD`, `ATLAS.GETSPRITE`, `ATLAS.FREE`), [TEXTURE.md](TEXTURE.md), [IMAGE.md](IMAGE.md).

---

### Sprite.Load

```basic
spr = SPRITE.LOAD(path$)
```

Loads an image file from disk (**PNG**, **JPG**, etc.). Returns a **sprite handle** covering the full texture as a single frame until you **`SPRITE.DEFANIM`** or use an **atlas** sub-rectangle.

**Parameters**

| Name | Type | Description |
|---|---|---|
| path$ | string | Path to image file. |

**Returns** — handle.

> **Common mistake:** Expecting **Aseprite JSON** animation import — use a **horizontal strip** with **`SPRITE.DEFANIM`** or **`ATLAS.LOAD`** + **`ATLAS.GETSPRITE`**.

**See also:** `ATLAS.LOAD`, `SPRITE.DEFANIM`

---

### Sprite.Free

```basic
SPRITE.FREE(spr)
```

Unloads the sprite’s **texture** (unless the sprite came from an **atlas** shared texture) and releases the heap slot. Call when the sprite is no longer needed.

**Parameters**

| Name | Type | Description |
|---|---|---|
| spr | handle | Sprite from **`SPRITE.LOAD`** or **`ATLAS.GETSPRITE`**. |

> **Common mistake:** Freeing a sprite that is still listed in a **`SPRITEGROUP`** / **`SPRITELAYER`** / **`SPRITEBATCH`** — remove it first or clear the container, or you may hold **stale handles**.

**See also:** `SPRITEGROUP.REMOVE`, `ATLAS.FREE`

---

### Sprite.Draw

```basic
SPRITE.DRAW(spr, x, y)
```

Draws the **current frame** at **integer** screen **`(x, y)`**, plus **`SPRITE.SETPOS`** offsets (float).

**Parameters**

| Name | Type | Description |
|---|---|---|
| spr | handle | Sprite handle. |
| x, y | int | Screen position (pixels). |

**Notes** — Use **`CAMERA2D.BEGIN`** / **`CAMERA2D.END`** when your game uses a 2D camera, or **`RENDER.BEGINMODE2D`** when that command is implemented.

---

### Sprite.SetPos / Sprite.SetPosition

```basic
SPRITE.SETPOS(spr, x#, y#)
SPRITE.SETPOSITION(spr, x#, y#)
```

**Alias pair.** Adds a **float** offset applied on top of **`SPRITE.DRAW`** / group draw positions. Used for smooth movement and hit tests.

---

### Sprite.DefAnim / Sprite.PlayAnim / Sprite.UpdateAnim

```basic
SPRITE.DEFANIM(spr, frameCount$)
SPRITE.PLAYANIM(spr, name$)
SPRITE.UPDATEANIM(spr, dt#)
```

**Strip mode:** **`frameCount$`** is a **decimal string** (e.g. `"4"`) = equal-width frames in **one row** inside the texture. Frame width = region width ÷ frame count.

**`SPRITE.PLAYANIM`** accepts **`name$`** for API symmetry; the strip player does not use distinct names yet.

**`SPRITE.UPDATEANIM`** advances time using internal FPS (default **8**). Pass **`TIME.DELTA()`** each frame.

**Note:** If **`ANIM.UPDATE`** drives the same sprite, **`SPRITE.UPDATEANIM`** does not advance strip frames for that object — use **one** animation system per sprite.

---

### Sprite.Hit / Sprite.PointHit / SPRITECOLLIDE

```basic
hit = SPRITE.HIT(a, b)
hit = SPRITECOLLIDE(a, b)
hit = SPRITE.POINTHIT(spr, x#, y#)
```

**`SPRITE.HIT`** / **`SPRITECOLLIDE`** — axis-aligned box overlap using **`SETPOS`** and current frame size.

**`SPRITE.POINTHIT`** — point-in-rect for **screen** coordinates.

**Returns** — boolean.

---

## SpriteGroup.* (handle-based group)

Groups are **named by handle**, not by string. **`SPRITEGROUP.MAKE`** takes **no arguments** and returns a **group handle**.

| Command | Signature | Notes |
|---|---|---|
| `SPRITEGROUP.MAKE` | `()` → handle | New empty group. |
| `SPRITEGROUP.ADD` | `(group, spr)` | Append sprite. |
| `SPRITEGROUP.REMOVE` | `(group, spr)` | Remove first occurrence of **spr**; no error if absent. |
| `SPRITEGROUP.CLEAR` | `(group)` | Remove all members. |
| `SPRITEGROUP.DRAW` | `(group, x, y)` | Draw each member at base **`(x,y)`** + **`SETPOS`**. |
| `SPRITEGROUP.FREE` | `(group)` | Frees **only** the group object (not member sprites). |

> **Spec note:** Some docs describe **`SpriteGroup.Make(name$)`** — in moonBASIC the group is a **handle**; use your own string table if you need names.

---

## SpriteLayer.*

| Command | Signature | Notes |
|---|---|---|
| `SPRITELAYER.MAKE` | `(z#)` → handle | **`z#`** stored for your sorting; draw order is under your control. |
| `SPRITELAYER.ADD` | `(layer, spr)` | |
| `SPRITELAYER.CLEAR` | `(layer)` | Remove all members. |
| `SPRITELAYER.SETZ` | `(layer, z#)` | Update stored **z**. |
| `SPRITELAYER.DRAW` | `(layer, x, y)` | Same base position semantics as group draw. |
| `SPRITELAYER.FREE` | `(layer)` | Frees layer only. |

---

## SpriteBatch.*

Records **multiple** **`(sprite, x, y)`** draws; **`SPRITEBATCH.DRAW`** executes them in order.

| Command | Notes |
|---|---|
| `SPRITEBATCH.MAKE` | `()` → handle |
| `SPRITEBATCH.ADD` | `(batch, spr, x, y)` — **int** positions |
| `SPRITEBATCH.CLEAR` | `(batch)` |
| `SPRITEBATCH.DRAW` | `(batch)` |
| `SPRITEBATCH.FREE` | `(batch)` |

---

## SpriteUI.*

Anchored placement using **fractions of screen size** (e.g. **`0.5, 0.5`** = center).

```basic
ui = SPRITEUI.MAKE(spr, anchorX#, anchorY#)
SPRITEUI.DRAW(ui, SCREENW(), SCREENH())
SPRITEUI.FREE(ui)
```

**`SPRITEUI.FREE`** releases only the **UI wrapper**; the sprite remains.

---

## Particle2D.* (simple filled circles)

CPU-side **circles** (no texture). **`PARTICLE2D.MAKE(max, r, g, b, a)`** sets pool size and colour; **`EMIT`** adds particles; **`UPDATE`** integrates velocity and **`life`**; **`DRAW`** renders.

| Command | Arguments |
|---|---|
| `PARTICLE2D.MAKE` | `(max, r, g, b, a)` |
| `PARTICLE2D.EMIT` | `(p, x, y, vx, vy, life#)` |
| `PARTICLE2D.UPDATE` | `(p, dt#)` |
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

x# = 300
y# = 250

WHILE NOT Window.ShouldClose()
    SPRITE.SETPOS(hero, x#, y#)
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
