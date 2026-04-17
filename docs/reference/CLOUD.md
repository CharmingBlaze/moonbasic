# Cloud Commands

**Coverage** and timing state for volumetric-style clouds. The current implementation stores parameters and reserves a draw hook; visual detail may be minimal until shaders are extended. **CGO** required.

Page shape: [DOC_STYLE_GUIDE.md](../DOC_STYLE_GUIDE.md) (**WAVE pattern**).

## Core Workflow

Create a layer with **`CLOUD.CREATE()`**, tune **`CLOUD.SETCOVERAGE`**, then drive **`CLOUD.UPDATE(cloud, dt)`** and **`CLOUD.DRAW(cloud)`** each frame. Draw order with **`SKY.*`** matters — see [SKY.md](SKY.md).

---

### `CLOUD.CREATE()` 

Creates a cloud layer object. **`CLOUD.MAKE()`** is a deprecated alias.

---

### `CLOUD.FREE(cloud)` 

Frees the cloud handle.

---

### `CLOUD.UPDATE(cloud, dt)` 

Advances simulation time.

---

### `CLOUD.DRAW(cloud)` 

Draw pass (may be a no-op depending on build / shader support).

---

### `CLOUD.SETCOVERAGE(cloud, coverage)` 

**`coverage`** in **0–1** (clamped), affecting density/opacity where implemented.

---

## Full Example

Minimal sketch paired with a sky (camera setup omitted):

```basic
sky = SKY.CREATE()
clouds = CLOUD.CREATE()
CLOUD.SETCOVERAGE(clouds, 0.45)

WHILE NOT WINDOW.SHOULDCLOSE()
    dt = TIME.DELTA()
    SKY.UPDATE(sky, dt)
    CLOUD.UPDATE(clouds, dt)
    RENDER.CLEAR(8, 10, 14)
    SKY.DRAW(sky)
    CLOUD.DRAW(clouds)
    ; ... rest of scene ...
    RENDER.FRAME()
WEND

CLOUD.FREE(clouds)
SKY.FREE(sky)
```

---

## See also

- [SKY.md](SKY.md) — draw order
- [WEATHER.md](WEATHER.md) — precipitation coverage
