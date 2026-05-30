# Texture Draw Wrapper Commands

Persistent draw objects for textured quads, source rectangles, and pro (rotation + origin) texture draws, plus font-based text objects.

Page shape follows [DOC_STYLE_GUIDE.md](../DOC_STYLE_GUIDE.md) (**WAVE pattern**).

## Core Workflow

1. Create a wrapper with a factory call (`DRAWTEX2(tex)`, `DRAWTEXREC(tex)`, `DRAWTEXPRO(tex)`, or `TEXTOBJEX(font, text)`).
2. Configure position, color, source rect, rotation, etc.
3. Call `.DRAW(handle)` each frame.
4. Free with `.FREE(handle)` when done.

These are retained-mode wrappers — set properties once, draw many frames without reconfiguring. For immediate-mode draws see [DRAW2D.md](DRAW2D.md) and [DRAW_WRAPPERS.md](DRAW_WRAPPERS.md).

---

## DRAWTEX2 — Simple Texture Draw

### `DRAWTEX2(textureHandle)` 

Creates a simple texture draw object. Returns a handle. Factory; also listed in the manifest as a global.

---

### `DRAWTEX2.POS(handle, x, y)` 

Sets the screen position (integer pixels).

---

### `DRAWTEX2.COLOR(handle, r, g, b, a)` 

Sets the tint color (0–255 per channel).

---

### `DRAWTEX2.COL(handle, r, g, b, a)` 

Alias for `DRAWTEX2.COLOR`.

---

### `DRAWTEX2.SETTEXTURE(handle, textureHandle)` 

Swaps the texture on an existing draw object.

---

### `DRAWTEX2.DRAW(handle)` 

Draws the texture at the configured position and tint.

---

### `DRAWTEX2.FREE(handle)` 

Frees the draw object.

---

## DRAWTEXREC — Source-Rectangle Texture Draw

### `DRAWTEXREC(textureHandle)` 

Creates a texture draw object with a configurable source rectangle. Returns a handle.

---

### `DRAWTEXREC.SRC(handle, x, y, w, h)` 

Sets the source rectangle on the texture (in texels).

---

### `DRAWTEXREC.POS(handle, x, y)` 

Sets the destination position (float).

---

### `DRAWTEXREC.COLOR(handle, r, g, b, a)` 

Sets the tint color.

---

### `DRAWTEXREC.COL(handle, r, g, b, a)` 

Alias for `DRAWTEXREC.COLOR`.

---

### `DRAWTEXREC.SETTEXTURE(handle, textureHandle)` 

Swaps the texture.

---

### `DRAWTEXREC.DRAW(handle)` 

Draws the source rectangle at the configured position.

---

### `DRAWTEXREC.FREE(handle)` 

Frees the draw object.

---

## DRAWTEXPRO — Pro Texture Draw (Rotation + Origin)

### `DRAWTEXPRO(textureHandle)` 

Creates a pro texture draw object with source rect, destination rect, origin, and rotation. Returns a handle.

---

### `DRAWTEXPRO.SRC(handle, x, y, w, h)` 

Sets the source rectangle on the texture.

---

### `DRAWTEXPRO.DST(handle, x, y, w, h)` 

Sets the destination rectangle on screen.

---

### `DRAWTEXPRO.ORIGIN(handle, ox, oy)` 

Sets the rotation origin (relative to destination rect).

---

### `DRAWTEXPRO.ROT(handle, angle)` 

Sets the rotation angle in degrees.

---

### `DRAWTEXPRO.COLOR(handle, r, g, b, a)` 

Sets the tint color.

---

### `DRAWTEXPRO.COL(handle, r, g, b, a)` 

Alias for `DRAWTEXPRO.COLOR`.

---

### `DRAWTEXPRO.SETTEXTURE(handle, textureHandle)` 

Swaps the texture.

---

### `DRAWTEXPRO.DRAW(handle)` 

Draws the texture with source rect, destination rect, origin, and rotation applied.

---

### `DRAWTEXPRO.FREE(handle)` 

Frees the draw object.

---

## TEXTEXOBJ — Font-Based Text Object

### `TEXTOBJEX(fontHandle, text)` 

Creates a font-based text draw object. Returns a handle.

---

### `TEXTEXOBJ.POS(handle, x, y)` 

Sets the draw position (float).

---

### `TEXTEXOBJ.SIZE(handle, fontSize)` 

Sets the font size.

---

### `TEXTEXOBJ.SPACING(handle, charSpacing)` 

Sets the character spacing.

---

### `TEXTEXOBJ.COLOR(handle, r, g, b, a)` 

Sets the text color.

---

### `TEXTEXOBJ.SETTEXT(handle, newText)` 

Changes the displayed string.

---

### `TEXTEXOBJ.DRAW(handle)` 

Draws the text at the configured position, size, and color.

---

### `TEXTEXOBJ.FREE(handle)` 

Frees the text object.

---

## Full Example

This example creates several texture draw wrappers and a text object, renders them, then cleans up.

```basic
tex = TEXTURE.LOAD("spritesheet.png")
font = LOADFONT("mono.ttf", 24)

; Simple texture draw
simple = DRAWTEX2(tex)
DRAWTEX2.POS(simple, 10, 10)
DRAWTEX2.COLOR(simple, 255, 255, 255, 255)

; Source-rect draw (show a 32x32 tile from the sheet)
tiled = DRAWTEXREC(tex)
DRAWTEXREC.SRC(tiled, 0, 0, 32, 32)
DRAWTEXREC.POS(tiled, 100, 10)

; Pro draw with rotation
pro = DRAWTEXPRO(tex)
DRAWTEXPRO.SRC(pro, 0, 0, 64, 64)
DRAWTEXPRO.DST(pro, 200, 50, 128, 128)
DRAWTEXPRO.ORIGIN(pro, 64, 64)
DRAWTEXPRO.ROT(pro, 45.0)

; Font text object
label = TEXTOBJEX(font, "Score: 0")
TEXTEXOBJ.POS(label, 10, 200)
TEXTEXOBJ.SIZE(label, 20)
TEXTEXOBJ.COLOR(label, 255, 255, 0, 255)

WHILE NOT WINDOW.SHOULDCLOSE()
    RENDER.BEGINFRAME()
    DRAWTEX2.DRAW(simple)
    DRAWTEXREC.DRAW(tiled)
    DRAWTEXPRO.DRAW(pro)
    TEXTEXOBJ.DRAW(label)
    RENDER.ENDFRAME()
WEND

; Cleanup
DRAWTEX2.FREE(simple)
DRAWTEXREC.FREE(tiled)
DRAWTEXPRO.FREE(pro)
TEXTEXOBJ.FREE(label)
```
