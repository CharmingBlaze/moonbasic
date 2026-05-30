# DrawTex Commands

Retained-mode texture draw objects: configure source/destination rectangles, rotation, origin, and color once, then call `DRAW` each frame. Three variants for different use cases.

| Handle type | Best for |
|---|---|
| `DRAWTEX2` | Simple position + color, no source rect |
| `DRAWTEXREC` | Source rect crop from a texture atlas |
| `DRAWTEXPRO` | Full pro control: src rect, dst rect, origin, rotation |

## Core Workflow

1. Create via the appropriate variant and set its texture with `SETTEXTURE`.
2. Configure properties each frame as needed.
3. Call `DRAW(handle)` to render.
4. `FREE(handle)` when done.

---

## DRAWTEX2 — Simple Texture Draw

### `DRAWTEX2.SETTEXTURE(handle, texHandle)` 

Binds a texture to this draw object.

---

### `DRAWTEX2.POS(handle, x, y)` 

Sets the draw position in screen pixels.

---

### `DRAWTEX2.COLOR(handle, r, g, b, a)` / `DRAWTEX2.COL(handle, r, g, b, a)` 

Sets the color tint (0–255).

---

### `DRAWTEX2.DRAW(handle)` 

Draws the texture at the configured position and tint.

---

### `DRAWTEX2.FREE(handle)` 

Frees the draw object.

---

## DRAWTEXREC — Source Rect Crop

### `DRAWTEXREC.SETTEXTURE(handle, texHandle)` 

Binds a texture atlas.

---

### `DRAWTEXREC.SRC(handle, x, y, width, height)` 

Sets the source rectangle within the texture (pixel coordinates).

---

### `DRAWTEXREC.POS(handle, x, y)` 

Sets the screen draw position.

---

### `DRAWTEXREC.COLOR(handle, r, g, b, a)` / `DRAWTEXREC.COL(handle, r, g, b, a)` 

Sets the color tint.

---

### `DRAWTEXREC.DRAW(handle)` 

Draws the cropped texture region.

---

### `DRAWTEXREC.FREE(handle)` 

Frees the draw object.

---

## DRAWTEXPRO — Full Pro Control

### `DRAWTEXPRO.SETTEXTURE(handle, texHandle)` 

Binds a texture.

---

### `DRAWTEXPRO.SRC(handle, x, y, width, height)` 

Source rectangle in the texture.

---

### `DRAWTEXPRO.DST(handle, x, y, width, height)` 

Destination rectangle on screen (controls scale and position).

---

### `DRAWTEXPRO.ORIGIN(handle, ox, oy)` 

Sets the rotation/scale origin point within the destination rect. `(0, 0)` = top-left, `(width/2, height/2)` = centre.

---

### `DRAWTEXPRO.ROT(handle, degrees)` 

Sets the rotation angle in degrees around the origin.

---

### `DRAWTEXPRO.COLOR(handle, r, g, b, a)` / `DRAWTEXPRO.COL(handle, r, g, b, a)` 

Sets the color tint.

---

### `DRAWTEXPRO.DRAW(handle)` 

Draws the texture with all pro parameters applied.

---

### `DRAWTEXPRO.FREE(handle)` 

Frees the draw object.

---

## Full Example

An atlas sprite spinning with DRAWTEXPRO.

```basic
WINDOW.OPEN(800, 600, "DrawTex Demo")
WINDOW.SETFPS(60)

atlas = TEXTURE.LOAD("assets/sprites.png")

; crop frame 0 from a 64x64 grid
d = DRAWTEXPRO.SRC(0, 0, 0, 64, 64)    ; placeholder — configure after
DRAWTEXPRO.SETTEXTURE(d, atlas)
DRAWTEXPRO.SRC(d, 0, 0, 64, 64)        ; frame 0 top-left of atlas
DRAWTEXPRO.DST(d, 400, 300, 96, 96)    ; drawn 96x96 at screen centre
DRAWTEXPRO.ORIGIN(d, 48, 48)           ; rotate around centre
DRAWTEXPRO.COLOR(d, 255, 255, 255, 255)

angle = 0.0
WHILE NOT WINDOW.SHOULDCLOSE()
    angle = angle + 60 * TIME.DELTA()
    DRAWTEXPRO.ROT(d, angle)

    RENDER.CLEAR(20, 20, 40)
    DRAWTEXPRO.DRAW(d)
    RENDER.FRAME()
WEND

DRAWTEXPRO.FREE(d)
TEXTURE.UNLOAD(atlas)
WINDOW.CLOSE()
```

---

## See also

- [SPRITE.md](SPRITE.md) — higher-level sprite system
- [TEXTURE.md](TEXTURE.md) — texture loading
- [ATLAS.md](ATLAS.md) — texture atlas management
- [DRAW2D.md](DRAW2D.md) — immediate-mode texture drawing
