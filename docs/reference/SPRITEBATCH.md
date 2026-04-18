# SpriteBatch / SpriteGroup / SpriteLayer / SpriteUI Commands

Sprite organisation and batch rendering handles. Group sprites for efficient draw calls, layer them by Z-depth, and build simple UI panels from sprites.

## SpriteBatch

A flat list of sprite + position pairs drawn in one GPU batch.

### `SPRITEBATCH.CREATE()` 

Creates an empty sprite batch. Returns a **batch handle**.

---

### `SPRITEBATCH.ADD(batch, spriteHandle, x, y)` 

Adds a sprite at screen position `(x, y)` to the batch.

---

### `SPRITEBATCH.CLEAR(batch)` 

Removes all sprites from the batch.

---

### `SPRITEBATCH.DRAW(batch)` 

Draws all sprites in the batch in one call.

---

### `SPRITEBATCH.FREE(batch)` 

Frees the batch handle.

---

## SpriteGroup

A named group of sprites drawn together at an offset. Useful for compound characters or UI elements built from multiple sprites.

### `SPRITEGROUP.CREATE()` 

Creates an empty sprite group. Returns a **group handle**.

---

### `SPRITEGROUP.ADD(group, spriteHandle)` 

Adds a sprite to the group.

---

### `SPRITEGROUP.REMOVE(group, spriteHandle)` 

Removes a specific sprite from the group.

---

### `SPRITEGROUP.CLEAR(group)` 

Removes all sprites from the group.

---

### `SPRITEGROUP.DRAW(group, x, y)` 

Draws all sprites in the group at world offset `(x, y)`.

---

### `SPRITEGROUP.FREE(group)` 

Frees the group handle.

---

## SpriteLayer

Organises sprites into depth-sorted layers. Multiple layers can be stacked to achieve parallax or Z-ordering.

### `SPRITELAYER.CREATE(zDepth)` 

Creates a sprite layer at depth `zDepth`. Returns a **layer handle**.

---

### `SPRITELAYER.ADD(layer, spriteHandle)` 

Adds a sprite to this layer.

---

### `SPRITELAYER.CLEAR(layer)` 

Removes all sprites from the layer.

---

### `SPRITELAYER.SETZ(layer, zDepth)` 

Updates the layer's Z depth (affects draw order among layers).

---

### `SPRITELAYER.DRAW(layer, x, y)` 

Draws all sprites in the layer at offset `(x, y)`.

---

### `SPRITELAYER.FREE(layer)` 

Frees the layer handle.

---

## SpriteUI

A sprite pinned to screen space, useful for HUD elements like health bars, icons, and panel backgrounds.

### `SPRITEUI.CREATE(spriteHandle, x, y)` 

Creates a UI sprite element at screen position `(x, y)`. Returns a **ui element handle**.

---

### `SPRITEUI.DRAW(handle, x, y)` 

Draws the UI element at `(x, y)` (overrides creation position for this frame).

---

### `SPRITEUI.FREE(handle)` 

Frees the UI element handle.

---

## Core Workflow

1. Load sprites with `SPRITE.*` or `TEXTURE.LOAD`.
2. Create a `SPRITEBATCH` for bulk rendering or `SPRITELAYER` for Z-sorted scenes.
3. `ADD` sprites each frame (or once if static).
4. `DRAW` inside `CAMERA2D.BEGIN` / `CAMERA2D.END` for world-space or outside for screen-space.
5. `FREE` all handles on exit.

---

## Full Example

A 2D scene with a background layer, entity layer, and HUD.

```basic
WINDOW.OPEN(800, 600, "SpriteBatch Demo")
WINDOW.SETFPS(60)

cam = CAMERA2D.CREATE()
CAMERA2D.SETOFFSET(cam, 400, 300)

bgTex     = TEXTURE.LOAD("assets/bg.png")
charTex   = TEXTURE.LOAD("assets/char.png")
iconTex   = TEXTURE.LOAD("assets/icon.png")

bgSprite   = SPRITE.CREATE(bgTex)
charSprite = SPRITE.CREATE(charTex)

bgLayer   = SPRITELAYER.CREATE(0.0)
charLayer = SPRITELAYER.CREATE(1.0)
SPRITELAYER.ADD(bgLayer, bgSprite)
SPRITELAYER.ADD(charLayer, charSprite)

hud = SPRITEUI.CREATE(SPRITE.CREATE(iconTex), 20, 20)

px = 0.0
py = 0.0

WHILE NOT WINDOW.SHOULDCLOSE()
    dt = TIME.DELTA()
    IF INPUT.KEYDOWN(KEY_RIGHT) THEN px = px + 200 * dt
    IF INPUT.KEYDOWN(KEY_LEFT)  THEN px = px - 200 * dt

    CAMERA2D.SETTARGET(cam, px, py)
    SPRITE.SETPOS(charSprite, px, py)

    RENDER.CLEAR(30, 40, 60)
    CAMERA2D.BEGIN(cam)
        SPRITELAYER.DRAW(bgLayer, 0, 0)
        SPRITELAYER.DRAW(charLayer, 0, 0)
    CAMERA2D.END()
    SPRITEUI.DRAW(hud, 20, 20)
    RENDER.FRAME()
WEND

SPRITELAYER.FREE(bgLayer)
SPRITELAYER.FREE(charLayer)
SPRITEUI.FREE(hud)
CAMERA2D.FREE(cam)
WINDOW.CLOSE()
```

---

## Extended Command Reference

| Command | Description |
|--------|-------------|
| `SPRITEBATCH.MAKE(...)` | Deprecated alias of `SPRITEBATCH.CREATE`. |

---

## See also

- [SPRITE.md](SPRITE.md) — sprite creation, animation, properties
- [CAMERA2D.md](CAMERA2D.md) — 2D scrolling camera
- [ATLAS.md](ATLAS.md) — texture atlas for sprite sheets
