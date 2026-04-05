# Sprite — `Sprite.*`, atlas, and `ANIM.*`

Sprites are **GPU textures** (or regions of an atlas) plus **frame layout** and optional **animation** state. They are drawn with **`Sprite.Draw`** (typically inside **`Render.BeginMode2D`** / **`EndMode2D`**).

**Requires CGO** (same as `Texture.*`, `Draw.*`).

**Related:** [ATLAS.md](ATLAS.md) for **`Atlas.Load`**, **`Atlas.GetSprite`**, **`Atlas.Free`**.

---

## `Sprite.Load(path$)`

Loads an image file with Raylib **`LoadTexture`** (e.g. **`.png`**, **`.jpg`**). Returns a **sprite handle** (one full texture, one logical frame until you call **`Sprite.DefAnim`**).

This is **not** an Aseprite/JSON parser; use a **horizontal strip** (all frames same size, laid out left-to-right) or an **atlas** (see below).

---

## Horizontal strip animation

### `Sprite.DefAnim(sprite, frameCount$)`

`frameCount$` is a **decimal string** (e.g. `"4"`) = number of **equal-width** frames in a **single row** inside the texture (from the sprite’s source region). Frame width = available width ÷ frame count.

### `Sprite.PlayAnim(sprite, name$)`

Starts strip playback from frame 0. The **`name$`** argument is accepted for API symmetry but **is not used** by the strip player (reserved for future use).

### `Sprite.UpdateAnim(sprite, deltaTime#)`

Advances the strip using the sprite’s internal FPS (default **8**). Pass **`Time.Delta()`** each frame. If the sprite is using the **`ANIM.*`** state machine (**`ANIM.UPDATE`**), **`Sprite.UpdateAnim`** does not advance strip frames (use one system or the other).

### `Sprite.Draw(sprite, x, y)`

Draws the **current frame** at **screen** `(x, y)`, plus any offset from **`Sprite.SetPos`**.

---

## Position offset & hit test

### `Sprite.SetPos(sprite, x#, y#)` / `Sprite.SetPosition(...)`

Alias pair. Adds a **float** offset to the draw position (useful for sub-pixel movement). **`Sprite.HIT`** uses these offsets with **`frameW` / `frameH`**.

### `Sprite.Hit(spriteA, spriteB)`

Returns **`TRUE`** if the two sprites’ **axis-aligned boxes** overlap, using each sprite’s **`SetPos`** and its **frame width/height** (current frame size).

### `SpriteCollide(spriteA, spriteB)` (registry: `SPRITECOLLIDE`)

Alias for **`Sprite.Hit`** — same AABB test, DarkBASIC-style name.

### `Sprite.PointHit(sprite, x#, y#)`

Returns **`TRUE`** if **screen** point **`(x, y)`** lies inside the sprite’s current-frame axis-aligned rectangle (**`SetPos`** + frame size).

---

## Texture atlas workflow

Use **`Atlas.Load`**, **`Atlas.GetSprite`**, and **`Atlas.Free`** as documented in **[ATLAS.md](ATLAS.md)**. Each **`Atlas.GetSprite`** handle shares the atlas texture; **`Sprite.Draw`** / **`DefAnim`** work the same as for **`Sprite.Load`**.

---

## `ANIM.*` — frame-range state machine (optional)

For **named clips** (frame index ranges) and **parameter-driven transitions**, use **`ANIM.*`** on the same sprite handle. This is separate from **`Sprite.DefAnim`** strip mode; if **`ANIM.UPDATE`** is in use, **`Sprite.UpdateAnim`** is a no-op for that sprite.

| Command | Arguments | Notes |
|--------|-----------|--------|
| `Anim.Define(sprite, name$, first, last, fps#, looping?)` | | Inclusive frame indices; **`looping?`** bool |
| `Anim.AddTransition(sprite, from$, to$, condition$)` | | See conditions below |
| `Anim.Update(sprite, dt#)` | | Advance frames + evaluate transitions |
| `Anim.SetParam(sprite, name$, value)` | | **`value`** numeric or bool; names are compared **case-insensitive** |

### Transition conditions

- **Comparison:** `param >= 1`, `speed == 0`, etc. Left side is a **parameter name** (lowercased internally); right side is a number.
- **Single name:** uses **bool** param, or numeric truthy float param.
- **Literals:** `true` / `false`.

---

## Full example (strip)

```basic
IF NOT Window.Open(800, 600, "Sprite strip") THEN END
ENDIF
Window.SetFPS(60)

; sheet.png = one row of 4 equal frames
hero = Sprite.Load("sheet.png")
Sprite.DefAnim(hero, "4")
Sprite.PlayAnim(hero, "walk")

x# = 300
y# = 250

WHILE NOT Window.ShouldClose()
    Sprite.SetPos(hero, x#, y#)
    Sprite.UpdateAnim(hero, Time.Delta())

    Render.Clear(30, 40, 50)
    Render.BeginMode2D()
        Sprite.Draw(hero, 0, 0)
    Render.EndMode2D()
    Render.Frame()
WEND

Window.Close()
```

---

## `SpriteGroup.*` — ordered draws

| Command | Arguments | Notes |
|--------|-----------|--------|
| `SpriteGroup.Make` | (none) | Returns a **group** handle (does not own member sprites). |
| `SpriteGroup.Add` | `(group, sprite)` | Append a **sprite** handle; order is draw order. |
| `SpriteGroup.Clear` | `(group)` | Remove all members. |
| `SpriteGroup.Draw` | `(group, x, y)` | Draws each member at the same base **integer** `(x,y)` (same rules as **`Sprite.Draw`**, including **`SetPos`** offsets). |
| `SpriteGroup.Free` | `(group)` | Frees only the group container. |

---

## `SpriteLayer.*` — layer + z metadata

| Command | Arguments | Notes |
|--------|-----------|--------|
| `SpriteLayer.Make` | `(z#)` | Returns a **layer** handle; **`z`** is stored for your own sorting (draw multiple layers in the order you choose). |
| `SpriteLayer.Add` | `(layer, sprite)` | |
| `SpriteLayer.SetZ` | `(layer, z#)` | Updates stored **z**. |
| `SpriteLayer.Draw` | `(layer, x, y)` | Same drawing semantics as **`SpriteGroup.Draw`**. |
| `SpriteLayer.Free` | `(layer)` | Frees only the layer object. |

---

## `SpriteBatch.*` — repeated placements

| Command | Arguments | Notes |
|--------|-----------|--------|
| `SpriteBatch.Make` | (none) | |
| `SpriteBatch.Add` | `(batch, sprite, x, y)` | Records one draw at **integer** screen `(x,y)` (plus each sprite’s **`SetPos`** offset). |
| `SpriteBatch.Clear` | `(batch)` | Clears recorded entries. |
| `SpriteBatch.Draw` | `(batch)` | Executes all recorded draws. |
| `SpriteBatch.Free` | `(batch)` | |

---

## `SpriteUI.*` — anchor placement

| Command | Arguments | Notes |
|--------|-----------|--------|
| `SpriteUI.Make` | `(sprite, anchorX#, anchorY#)` | **Anchors** in **0–1** range (e.g. `0.5, 0.5` = center). |
| `SpriteUI.Draw` | `(ui, screenW, screenH)` | Positions from **`SCREENW`/`SCREENH`**-style dimensions (pass **`SCREENW()`**, **`SCREENH()`** in game code). |
| `SpriteUI.Free` | `(ui)` | Frees only the UI wrapper (not the sprite). |

---

## `Particle2D.*` — simple CPU circles

Colored circles (no texture). **`Emit`** appends particles until the pool is full; **`Update`** integrates motion and subtracts **`dt`** from **`life`**.

| Command | Arguments |
|--------|-----------|
| `Particle2D.Make` | `(max, r, g, b, a)` |
| `Particle2D.Emit` | `(p, x, y, vx, vy, life#)` |
| `Particle2D.Update` | `(p, dt#)` |
| `Particle2D.Draw` | `(p)` |
| `Particle2D.Free` | `(p)` |

---

## Command checklist (implemented)

**Sprite:** `Load`, `Draw`, `SetPos`, `SetPosition`, `DefAnim`, `PlayAnim`, `UpdateAnim`, `Hit`.

**SpriteGroup / SpriteLayer / SpriteBatch / SpriteUI / Particle2D:** commands listed in the sections above.

**Atlas:** `Load`, `Free`, `GetSprite` — [ATLAS.md](ATLAS.md).

**Anim:** `Define`, `AddTransition`, `Update`, `SetParam` — this page § **`ANIM.*`**.

---

## See also

- [ATLAS.md](ATLAS.md) — packed sheets + JSON
- [TEXTURE.md](TEXTURE.md) — raw GPU textures
- [IMAGE.md](IMAGE.md) — CPU images before upload
