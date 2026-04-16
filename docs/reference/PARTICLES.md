# Particle Commands

Emitter-based **3D** particles: **`PARTICLE.*`** (alias namespace **`PARTICLE3D.*`** — same handlers). **2D** CPU particles: **`PARTICLE2D.*`** — see [PARTICLE.md](PARTICLE.md) and [SPRITE.md](SPRITE.md).

Full command table (all **`SET*`**, **`PARTICLE3D.*`**, handle methods): [PARTICLE.md](PARTICLE.md).

Page shape: [DOC_STYLE_GUIDE.md](../DOC_STYLE_GUIDE.md) (**WAVE pattern**).

## Core Workflow

**`PARTICLE.CREATE()`** returns an emitter handle (no max-count argument — capacity is internal). Configure with **`PARTICLE.SETTEXTURE`**, **`PARTICLE.SETEMITRATE`**, **`PARTICLE.SETLIFETIME`**, velocity/color/size/gravity helpers, then **`PARTICLE.PLAY`**. Each frame: **`PARTICLE.UPDATE(emitter, dt)`** then **`PARTICLE.DRAW(emitter)`** inside your **3D** camera pass. **`PARTICLE.FREE`** when done.

**`PARTICLE.MAKE`** / **`PARTICLE3D.MAKE`** are deprecated aliases of **`CREATE`**.

---

### `PARTICLE.CREATE()`

Returns a new **emitter handle**. Same as **`PARTICLE3D.CREATE()`**.

---

### `PARTICLE.FREE(handle)`

Frees the emitter and all associated resources.

---

### `PARTICLE.UPDATE(handle, dt)`

Advances simulation by **`dt`** seconds.

---

### `PARTICLE.DRAW(handle)` / `PARTICLE.DRAW(handle, cameraHandle)`

Draws billboards (or debug cubes if configured). Uses the active camera, or pass an explicit **3D** camera handle when overloaded — see [PARTICLE.md](PARTICLE.md).

---

### `PARTICLE.PLAY(handle)` / `PARTICLE.STOP(handle)`

Starts or stops continuous emission.

---

### `PARTICLE.SETBURST(handle, count)`

Immediately spawns up to **`count`** particles (capped).

---

## Configuration (selected)

The manifest lists every **`PARTICLE.SET*`** overload — see [PARTICLE.md](PARTICLE.md). Common knobs:

- **`PARTICLE.SETTEXTURE(handle, textureHandle)`**
- **`PARTICLE.SETEMITRATE`** / **`PARTICLE.SETRATE`**
- **`PARTICLE.SETPOS`**, **`PARTICLE.SETLIFETIME`**, **`PARTICLE.SETVELOCITY`**, **`PARTICLE.SETSPREAD`**, **`PARTICLE.SETSPEED`**
- **`PARTICLE.SETSTARTSIZE`** / **`PARTICLE.SETENDSIZE`** (or **`PARTICLE.SETSIZE`** shorthand)
- **`PARTICLE.SETCOLOR`** / **`PARTICLE.SETCOLOREND`**, **`PARTICLE.SETGRAVITY`**, **`PARTICLE.SETBILLBOARD`**

---

## Full Example (3D)

```basic
WINDOW.OPEN(960, 540, "Particles")
WINDOW.SETFPS(60)

myTex = TEXTURE.LOAD("assets/fx/spark.png")
cam = CAMERA.CREATE()
CAMERA.SETPOS(cam, 0, 2, 10)
CAMERA.SETTARGET(cam, 0, 0, 0)

p = PARTICLE.CREATE()
PARTICLE.SETTEXTURE(p, myTex)
PARTICLE.SETEMITRATE(p, 40)
PARTICLE.SETLIFETIME(p, 0.5, 1.5)
PARTICLE.SETVELOCITY(p, 0, 1, 0, 0.4)
PARTICLE.SETSPREAD(p, 0.4)
PARTICLE.SETSPEED(p, 0.8, 1.2)
PARTICLE.SETSTARTSIZE(p, 0.15, 0.35)
PARTICLE.SETENDSIZE(p, 0.05, 0.1)
PARTICLE.SETCOLOR(p, 255, 200, 100, 255)
PARTICLE.SETCOLOREND(p, 255, 50, 0, 0)
PARTICLE.SETGRAVITY(p, 0, -2, 0)
PARTICLE.SETPOS(p, 0, 0, 0)
PARTICLE.PLAY(p)

WHILE NOT WINDOW.SHOULDCLOSE()
    dt = TIME.DELTA()
    PARTICLE.UPDATE(p, dt)
    RENDER.CLEAR(20, 24, 32)
    RENDER.BEGIN3D(cam)
        PARTICLE.DRAW(p)
    RENDER.END3D()
    RENDER.FRAME()
WEND

PARTICLE.FREE(p)
TEXTURE.FREE(myTex)
CAMERA.FREE(cam)
WINDOW.CLOSE()
```

---

## Tips

- **Performance:** keep emit rates reasonable; use **`PARTICLE.COUNT`** for debugging.
- **Alias:** **`PARTICLE3D.*`** matches **`PARTICLE.*`** line-for-line.
