# Transition Commands

Full-screen fade and wipe overlays for scene changes and cinematic effects.

Page shape follows [DOC_STYLE_GUIDE.md](../DOC_STYLE_GUIDE.md) (**WAVE pattern**).

## Core Workflow

1. Optionally set the overlay color with `TRANSITION.SETCOLOR`.
2. Start a transition with `TRANSITION.FADEOUT`, `TRANSITION.FADEIN`, or `TRANSITION.WIPE`.
3. Poll `TRANSITION.ISDONE` each frame to know when to proceed.

For automatic scene transitions see [SCENE.md](SCENE.md).

---

### `TRANSITION.FADEOUT(seconds)` / `FADEIN`
Starts a full-screen color fade.

- **Arguments**:
    - `seconds`: (Float) Duration of the effect.
- **Returns**: (None)

---

### `TRANSITION.WIPE(direction, seconds)`
Starts a directional screen wipe.

- **Arguments**:
    - `direction`: (String) "left", "right", "up", "down".
    - `seconds`: (Float) Duration.
- **Returns**: (None)

---

### `TRANSITION.ISDONE()`
Returns `TRUE` if the current transition effect has completed.

- **Returns**: (Boolean)

---

### `TRANSITION.SETCOLOR(r, g, b [, a])`
Sets the color used for transitions (default black).

- **Returns**: (None)

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
