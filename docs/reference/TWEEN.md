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
Creates a new empty tween.

- **Returns**: (Handle) The new tween handle.

---

### `TWEEN.TO(handle, varName, target, seconds, easing)` / `THEN`
Appends an animation segment to the tween.

- **Arguments**:
    - `handle`: (Handle) The tween.
    - `varName`: (String) Global variable to animate.
    - `target`: (Float) Destination value.
    - `seconds`: (Float) Duration.
    - `easing`: (String) Easing function (e.g., "easeout", "bounce").
- **Returns**: (Handle) The tween handle (for chaining).

---

### `TWEEN.START(handle)` / `UPDATE(dt)` / `STOP()`
Controls the playback and progression of the tween.

- **Example**:
    ```basic
    t = TWEEN.MAKE()
    TWEEN.TO(t, "alpha", 255, 1.0, "linear")
    TWEEN.START(t)
    ```

---

### `TWEEN.LOOP(handle, count)` / `YOYO()`
Sets repeating and oscillation behavior.

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

---

## Extended Command Reference

| Command | Description |
|--------|-------------|
| `TWEEN.ISPLAYING(t)` | Returns `TRUE` if the tween is currently running. |
| `TWEEN.ISFINISHED(t)` | Returns `TRUE` if the tween has completed. |
| `TWEEN.GETLOOP(t)` | Returns the loop count setting (-1 = infinite). |
| `TWEEN.GETYOYO(t)` | Returns `TRUE` if yoyo (ping-pong) mode is enabled. |

## See also

- [TIME.md](TIME.md) — `TIME.DELTA` for manual lerp
- [MATH.md](MATH.md) — `MATH.LERP`, easing helpers
