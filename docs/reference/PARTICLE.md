# Particle Commands

3D billboard and 2D CPU particle emitters for effects like fire, smoke, and sparks.

Page shape follows [DOC_STYLE_GUIDE.md](../DOC_STYLE_GUIDE.md) (**WAVE pattern**).

## Core Workflow

1. Create an emitter with `PARTICLE.CREATE`.
2. Configure with `PARTICLE.SETTEXTURE`, `PARTICLE.SETEMITRATE`, `PARTICLE.SETLIFETIME`, `PARTICLE.SETPOS`, etc.
3. Call `PARTICLE.UPDATE(dt)` each frame.
4. Draw with `PARTICLE.DRAW` inside the 3D pass.
5. Free with `PARTICLE.FREE`.

`PARTICLE3D.*` is a full alias set. For 2D particles see `PARTICLE2D.*` in [SPRITE.md](SPRITE.md).

---

## `PARTICLE.*` / `PARTICLE3D.*` (3D)

| Command | Notes |
|---------|--------|
| **`PARTICLE.CREATE`** / **`PARTICLE3D.CREATE`** (canonical) / deprecated **`PARTICLE.MAKE`** / **`PARTICLE3D.MAKE`** | No args → emitter **handle**. |
| **`PARTICLE.FREE`** / **`PARTICLE3D.FREE`** | Free emitter. |
| **`PARTICLE.SETTEXTURE`** | `(emitter, textureHandle)` |
| **`PARTICLE.SETEMITRATE`** / **`PARTICLE.SETRATE`** | `(emitter, per_sec)` — **`SETRATE`** is an alias. **`PARTICLE3D.SETRATE`** same. |
| **`PARTICLE.SETPOS`** | `(emitter, x, y, z)` |
| **`PARTICLE.SETLIFETIME`** | `(emitter, min, max)` seconds |
| **`PARTICLE.SETVELOCITY`** | `(emitter, vx, vy, vz, spread)` — **spread** is random component noise (same units as velocity components). |
| **`PARTICLE.SETDIRECTION`** | `(emitter, vx, vy, vz)` — base direction; combine with **`SETSPREAD`**. |
| **`PARTICLE.SETSPREAD`** | `(emitter, angle)` — random jitter added to each velocity component at spawn. |
| **`PARTICLE.SETSPEED`** | `(emitter, min, max)` — per-particle scalar applied after direction+spread. Default `1..1`. |
| **`PARTICLE.SETSTARTSIZE`** | `(emitter, min, max)` — spawn size range. |
| **`PARTICLE.SETENDSIZE`** | `(emitter, min, max)` — end-of-life size range. |
| **`PARTICLE.SETSIZE`** | `(emitter, start, end)` — sets **both** start and end to **single** values (legacy shorthand). |
| **`PARTICLE.SETCOLOR`** / **`PARTICLE.SETSTARTCOLOR`** | `(emitter, r, g, b, a)` 0–255 |
| **`PARTICLE.SETCOLOREND`** / **`PARTICLE.SETENDCOLOR`** | End color |
| **`PARTICLE.SETGRAVITY`** | `(emitter, g)` **legacy** → `(0, g, 0)`, or `(emitter, gx, gy, gz)` |
| **`PARTICLE.SETBURST`** | `(emitter, count)` — spawn **count** particles immediately (capped). |
| **`PARTICLE.SETBILLBOARD`** | `(emitter, TRUE/FALSE)` — **`TRUE`**: **`DrawBillboard`**. **`FALSE`**: draw **cubes** at particle positions (debug / non-camera-facing). |
| **`PARTICLE.PLAY`** / **`PARTICLE.STOP`** | Start/stop continuous emission. |
| **`PARTICLE.UPDATE`** | `(emitter, dt)` |
| **`PARTICLE.DRAW`** | `(emitter)` uses **`CAMERA.BEGIN` … `CAMERA.END`** active camera, or **`(emitter, cameraHandle)`** for an explicit **`Camera3D`** handle. |
| **`PARTICLE.ISALIVE`** | `→ int` (`1` = still playing **or** live particles remain). |
| **`PARTICLE.COUNT`** | `→ int` live particles |

Every row above exists under **`PARTICLE3D.*`** as well (e.g. **`PARTICLE3D.SETTEXTURE`**, **`PARTICLE3D.DRAW`**, …).

---

## Handle methods

On a **`Particle`** handle, method calls map to **`PARTICLE.*`** keys (see **`vm/handlecall.go`**), e.g. **`emitter.SetPos`**, **`emitter.Play`**, **`emitter.SetStartColor`**.

---

## Full Example

```basic
WINDOW.OPEN(800, 600, "Particle Demo")
WINDOW.SETFPS(60)

cam = CAMERA.CREATE()
CAMERA.SETPOS(cam, 0, 3, -8)
CAMERA.SETTARGET(cam, 0, 1, 0)

emitter = PARTICLE.CREATE()
PARTICLE.SETPOS(emitter, 0, 0, 0)
PARTICLE.SETEMITRATE(emitter, 50)
PARTICLE.SETLIFETIME(emitter, 0.5, 2.0)
PARTICLE.SETSTARTCOLOR(emitter, 255, 200, 50, 255)
PARTICLE.SETENDCOLOR(emitter, 255, 50, 10, 0)

WHILE NOT WINDOW.SHOULDCLOSE()
    PARTICLE.UPDATE(emitter, TIME.DELTA())

    RENDER.CLEAR(20, 20, 30)
    RENDER.BEGIN3D(cam)
        PARTICLE.DRAW(emitter)
    RENDER.END3D()
    RENDER.FRAME()
WEND

PARTICLE.FREE(emitter)
CAMERA.FREE(cam)
WINDOW.CLOSE()
```

---

## See also

- [PARTICLES.md](PARTICLES.md) — longer examples and workflow
- [CAMERA.md](CAMERA.md) — 3D pass for **`PARTICLE.DRAW`**
- [SPRITE.md](SPRITE.md) — **`PARTICLE2D.*`**
