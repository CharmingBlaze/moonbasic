# Particle system commands

Emitter-based **3D** particles (**`PARTICLE.*`** and alias **`PARTICLE3D.*`**) and **2D** particles (**`PARTICLE2D.*`**) — see [PARTICLE.md](PARTICLE.md) and [SPRITE.md](SPRITE.md).

---

## 3D workflow

1. **`p = PARTICLE.MAKE`** (or **`PARTICLE3D.MAKE`**) — returns a **`Particle`** handle (no arguments; up to **16000** live particles per emitter).
2. Configure: **`SETTEXTURE`**, **`SETEMITRATE`** / **`SETRATE`**, **`SETLIFETIME`**, direction (**`SETDIRECTION`** + **`SETSPREAD`**, or **`SETVELOCITY`** for combined), **`SETSPEED`**, sizes (**`SETSTARTSIZE`** / **`SETENDSIZE`** or legacy **`SETSIZE`**), colors (**`SETCOLOR`** / **`SETCOLOREND`**), **`SETGRAVITY`**, **`SETPOS`**.
3. **`PLAY`** to emit; **`STOP`** stops new spawns (existing particles age out).
4. Each frame: **`UPDATE(p, dt)`** with **`Time.Delta()`**.
5. Inside **`CAMERA.BEGIN` … `CAMERA.END`**: **`DRAW(p)`**, or **`DRAW(p, camHandle)`** to use a specific camera.
6. Query: **`COUNT(p)`**, **`ISALIVE(p)`** (returns **1** while playing or particles remain).
7. **`FREE(p)`** when done.

**Burst:** **`SETBURST(p, n)`** spawns **n** particles immediately.

**Billboard:** default **`DrawBillboard`**. **`SETBILLBOARD(p, FALSE)`** draws small **cubes** instead (non–camera-facing).

---

## Example (minimal 3D)

```basic
cam = Camera.Make()
cam.SetPos(0, 2, 10)
cam.SetTarget(0, 0, 0)

p = PARTICLE.MAKE
PARTICLE.SETTEXTURE p, myTex
PARTICLE.SETEMITRATE p, 40
PARTICLE.SETLIFETIME p, 0.5, 1.5
PARTICLE.SETDIRECTION p, 0, 1, 0
PARTICLE.SETSPREAD p, 0.4
PARTICLE.SETSPEED p, 0.8, 1.2
PARTICLE.SETSTARTSIZE p, 0.15, 0.35
PARTICLE.SETENDSIZE p, 0.05, 0.1
PARTICLE.SETCOLOR p, 255, 200, 100, 255
PARTICLE.SETCOLOREND p, 255, 50, 0, 0
PARTICLE.SETGRAVITY p, 0, -2, 0
PARTICLE.SETPOS p, 0, 0, 0
PARTICLE.PLAY p

WHILE NOT Window.ShouldClose()
    dt# = Time.Delta()
    PARTICLE.UPDATE p, dt#
    Render.Clear(20, 24, 32)
    cam.Begin()
        PARTICLE.DRAW p
    cam.End()
    Render.Frame()
WEND

PARTICLE.FREE p
Camera.Free cam
```

---

## Tips

- **Performance:** keep emit rates reasonable; **`COUNT`** for debugging.
- **Gravity:** use **`SETGRAVITY p, gx, gy, gz`** for world-space acceleration.
- **Alias:** **`PARTICLE3D.*`** matches **`PARTICLE.*`** line-for-line.
