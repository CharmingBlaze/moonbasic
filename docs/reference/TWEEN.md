# Tween Commands

Keyframe-style animation chains that read and write global variables by name.

Page shape follows [DOC_STYLE_GUIDE.md](../DOC_STYLE_GUIDE.md) (**WAVE pattern**).

## Core Workflow

1. Create a tween with `TWEEN.MAKE`.
2. Append segments with `TWEEN.TO` / `TWEEN.THEN`.
3. Optionally set `TWEEN.LOOP`, `TWEEN.YOYO`, or `TWEEN.ONCOMPLETE`.
4. Start with `TWEEN.START` and advance with `TWEEN.UPDATE(tween, dt)` each frame.
5. Stop with `TWEEN.STOP` if needed.

---

### `TWEEN.MAKE()` 

Creates an empty tween (default: **one** loop, no yoyo).

---

### `TWEEN.TO(tweenHandle, varName, target, seconds, easing)` 

Appends a segment that animates the global `varName` toward `target` over `seconds`.

- `easing`: `"linear"`, `"easein"`, `"easeout"`, `"easeinout"`, `"bounce"`, `"elastic"`, `"back"`, `"circ"`, `"expo"`, `"sine"` (case-insensitive; unknown defaults to linear).

---

### `TWEEN.THEN(tweenHandle, varName, target, seconds, easing)` 

Alias for `TWEEN.TO` — appends another segment after the previous.

---

### `TWEEN.ONCOMPLETE(tweenHandle, functionName)` 

Registers a user function to call when all loops finish.

---

### `TWEEN.LOOP(tweenHandle, count)` 

Sets loop count. `count <= 0` means infinite loops.

---

### `TWEEN.YOYO(tweenHandle)` 

Enables yoyo mode — the tween plays forward then backward per loop.

---

### `TWEEN.START(tweenHandle)` 

Begins playback from the first segment.

---

### `TWEEN.UPDATE(tweenHandle, dt)` 

Advances the tween by `dt` seconds. Call each frame.

---

### `TWEEN.STOP(tweenHandle)` 

Stops playback without calling OnComplete.

---

## Full Example

This example tweens a global `posX` from 0 to 400 with easing.

```basic
posX = 0.0

t = TWEEN.MAKE()
TWEEN.TO(t, "posX", 400.0, 2.0, "easeout")
TWEEN.THEN(t, "posX", 0.0, 2.0, "easein")
TWEEN.LOOP(t, 0)
TWEEN.YOYO(t)
TWEEN.START(t)

WHILE NOT WINDOW.SHOULDCLOSE()
    dt = DELTATIME()
    TWEEN.UPDATE(t, dt)

    RENDER.BEGINFRAME()
    DRAW.RECT(INT(posX), 300, 20, 20, 255, 100, 50, 255)
    RENDER.ENDFRAME()
WEND

TWEEN.STOP(t)
```
