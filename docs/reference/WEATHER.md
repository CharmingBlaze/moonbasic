# Weather Commands

**`WEATHER.*`**, **`FOG.*`**, and **`WIND.*`**: combined **weather state** (type, coverage), **distance fog** parameters tracked for the module, and **global wind**. **CGO** required for full behavior; fog may not call every Raylib entry point if a symbol is unavailable — state is still tracked.

Page shape: [DOC_STYLE_GUIDE.md](../DOC_STYLE_GUIDE.md) (**WAVE pattern**).

## Core Workflow

Create a **`WEATHER`** handle, drive **`WEATHER.UPDATE`** / **`WEATHER.DRAW`** each frame, and tune presets with **`WEATHER.SETTYPE`**. Use **`FOG.*`** for distance fog state and **`WIND.*`** for a horizontal wind vector other systems can read. This namespace is **weather-scoped** state; it does not replace your full render pipeline — combine with **`SKY.*`**, **`PARTICLE.*`**, and scene fog as documented in runtime.

---

### `WEATHER.CREATE()` 

Creates a weather controller handle.

---

### `WEATHER.FREE(weather)` 

Frees the controller.

---

### `WEATHER.UPDATE(weather, dt)` 

Per-frame update.

---

### `WEATHER.DRAW(weather)` 

Draw pass (particles / effects as implemented).

---

### `WEATHER.SETTYPE(weather, name)` 

Sets a named preset (for example `"clear"`, `"rain"` — see runtime).

---

### `WEATHER.GETCOVERAGE(weather)` 

Returns normalized coverage **0–1**.

---

### `WEATHER.GETTYPE(weather)` 

Returns the current type string.

---

### `FOG.ENABLE(enabled)` 

Turns fog application on/off where supported.

---

### `FOG.SETNEAR(near)` / `FOG.SETFAR(far)` 

Distance fog start and end.

---

### `FOG.SETCOLOR(r, g, b, a)` 

Fog color components **0–255**.

**Common mistake:** Expecting **`FOG.*`** to duplicate all **`RENDER.*`** fog APIs — this namespace is **weather-scoped** state; combine with your render pipeline as documented in runtime.

---

### `WIND.SET(strength, dx, dz)` 

Sets wind **strength** and a horizontal direction on the XZ plane (components need not be normalized).

---

### `WIND.GETSTRENGTH()` 

Reads current strength.

---

## Full Example

Minimal loop sketch (paths and presets depend on your project):

```basic
weather = WEATHER.CREATE()
WEATHER.SETTYPE(weather, "rain")
FOG.ENABLE(TRUE)
FOG.SETNEAR(8.0)
FOG.SETFAR(120.0)
FOG.SETCOLOR(180, 190, 200, 255)
WIND.SET(2.5, 1.0, 0.2)

WHILE NOT WINDOW.SHOULDCLOSE()
    dt = TIME.DELTA()
    WEATHER.UPDATE(weather, dt)
    RENDER.CLEAR(20, 24, 32)
    ; ... sky / scene ...
    WEATHER.DRAW(weather)
    RENDER.FRAME()
WEND

WEATHER.FREE(weather)
```

---

## Extended Command Reference

| Command | Description |
|--------|-------------|
| `WEATHER.MAKE(...)` | Deprecated alias of `WEATHER.CREATE`. |

---

## See also

- [PARTICLES.md](PARTICLES.md)
- [SKY.md](SKY.md)
