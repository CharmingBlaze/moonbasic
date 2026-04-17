# Transition Commands

Full-screen fade and wipe overlays for scene changes and cinematic effects.

Page shape follows [DOC_STYLE_GUIDE.md](../DOC_STYLE_GUIDE.md) (**WAVE pattern**).

## Core Workflow

1. Optionally set the overlay color with `TRANSITION.SETCOLOR`.
2. Start a transition with `TRANSITION.FADEOUT`, `TRANSITION.FADEIN`, or `TRANSITION.WIPE`.
3. Poll `TRANSITION.ISDONE` each frame to know when to proceed.

For automatic scene transitions see [SCENE.md](SCENE.md).

---

### `TRANSITION.FADEOUT(seconds)` 

Fades the overlay in, blocking the view over `seconds`.

---

### `TRANSITION.FADEIN(seconds)` 

Fades the overlay out, revealing the scene over `seconds`.

---

### `TRANSITION.ISDONE()` 

Returns `TRUE` when the current transition has finished.

---

### `TRANSITION.WIPE(direction, seconds)` 

Wipe overlay using the current color. `direction` (case-insensitive):

- `left` — bar grows from the left edge.
- `right` — from the right.
- `up` / `top` — from the top.
- `down` / `bottom` — from the bottom.

---

### `TRANSITION.SETCOLOR(r, g, b, a)` 

Sets the overlay color (0–255 per channel). Default is opaque black.

---

## Full Example

This example fades out, loads a new scene, then fades back in.

```basic
TRANSITION.SETCOLOR(0, 0, 0, 255)
TRANSITION.FADEOUT(0.5)

; Wait for fade to finish
WHILE NOT TRANSITION.ISDONE()
    RENDER.BEGINFRAME()
    RENDER.ENDFRAME()
WEND

; Load new scene here
SCENE.LOAD("level2")

TRANSITION.FADEIN(0.5)
WHILE NOT TRANSITION.ISDONE()
    RENDER.BEGINFRAME()
    RENDER.ENDFRAME()
WEND
```
