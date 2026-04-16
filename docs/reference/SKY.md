# Sky Commands

Day/night **tinted sky dome** (drawn as a large sphere) with time-of-day in **hours** and configurable **day length**. **CGO + Raylib** required.

Page shape: [DOC_STYLE_GUIDE.md](../DOC_STYLE_GUIDE.md) (**WAVE pattern**).

## Core Workflow

Create a sky with **`SKY.CREATE()`**, configure **`SKY.SETTIME`** / **`SKY.SETDAYLENGTH`**, then each frame call **`SKY.UPDATE(sky, dt)`** and **`SKY.DRAW(sky)`**. Draw the sky **early** in the 3D pass (usually **before** opaque terrain) so depth behaves correctly.

---

### `SKY.CREATE()`

Creates a sky object with default time and day length. **`SKY.MAKE()`** is a deprecated alias.

---

### `SKY.FREE(sky)`

Frees the sky handle.

---

### `SKY.UPDATE(sky, dt)`

Advances internal time using **`dt`** and day length.

---

### `SKY.DRAW(sky)`

Draws the sky dome. Call order is user-defined, but the sky should usually be drawn **first** inside the camera block.

---

### `SKY.SETTIME(sky, hours)` / `SKY.SETDAYLENGTH(sky, seconds)`

**`SETTIME`**: **0–24** style hours. **`SETDAYLENGTH`**: real-time **seconds** for a full day/night cycle.

---

### `SKY.GETTIMEHOURS(sky)`

Returns the current simulated hour (**float**).

---

### `SKY.ISNIGHT(sky)`

Returns **`TRUE`** when the sun is below the horizon (implementation threshold).

---

## Full Example

Minimal frame sketch (camera setup omitted):

```basic
sky = SKY.CREATE()
SKY.SETDAYLENGTH(sky, 600.0)

WHILE NOT WINDOW.SHOULDCLOSE()
    dt = TIME.DELTA()
    SKY.UPDATE(sky, dt)
    RENDER.CLEAR(10, 12, 18)
    ; Begin 3D / camera, then draw sky before terrain:
    SKY.DRAW(sky)
    ; ... terrain, entities ...
    RENDER.FRAME()
WEND

SKY.FREE(sky)
```

---

## See also

- [CLOUD.md](CLOUD.md)
- [WEATHER.md](WEATHER.md)
