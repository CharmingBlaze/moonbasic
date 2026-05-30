# Particle3D Commands

3D world-space particle system with billboard sprites, per-emitter gravity, color lifecycle, burst mode, and texture support.

## Core Workflow

1. `PARTICLE3D.CREATE()` — allocate an emitter handle.
2. Configure: `PARTICLE3D.SETPOS`, `PARTICLE3D.SETLIFETIME`, `PARTICLE3D.SETSPEED`, `PARTICLE3D.SETCOLOR`, etc.
3. `PARTICLE3D.PLAY(emitter)` — start continuous emission.
4. Each frame: `PARTICLE3D.UPDATE(emitter, dt)` → `PARTICLE3D.DRAW(emitter)`.
5. `PARTICLE3D.STOP(emitter)` then `PARTICLE3D.FREE(emitter)` when done.

---

## Creation

### `PARTICLE3D.CREATE()` 

Allocates a 3D particle emitter with default settings. Returns an **emitter handle**.

---

## Position

### `PARTICLE3D.SETPOS(emitter, x, y, z)` 

Sets the world-space emission origin.

- *Handle shortcut*: `emitter.setPos(x, y, z)`

---

### `PARTICLE3D.GETPOS(emitter)` 

Returns the emitter position as a 3-element array.

---

## Emission Rate & Lifetime

### `PARTICLE3D.SETEMITRATE(emitter, rate)` / `PARTICLE3D.SETRATE(emitter, rate)` 

Sets the continuous emission rate in particles per second.

---

### `PARTICLE3D.SETLIFETIME(emitter, minSec, maxSec)` 

Sets the random lifetime range (seconds) for each particle.

---

### `PARTICLE3D.SETBURST(emitter, count)` 

Sets the number of particles to spawn instantly on `PARTICLE3D.PLAY`.

---

## Velocity & Direction

### `PARTICLE3D.SETVELOCITY(emitter, vx, vy, vz, speed)` 

Sets the base emission direction `(vx, vy, vz)` and base speed.

---

### `PARTICLE3D.SETDIRECTION(emitter, dx, dy, dz)` 

Sets the emission direction vector.

---

### `PARTICLE3D.SETSPEED(emitter, minSpeed, maxSpeed)` 

Sets the random speed range for emitted particles.

---

### `PARTICLE3D.SETSPREAD(emitter, angleDeg)` 

Sets the cone spread angle in degrees around the direction vector.

---

### `PARTICLE3D.SETGRAVITY(emitter, gy)` / `PARTICLE3D.SETGRAVITY(emitter, gx, gy, gz)` 

Applies gravity to particles. Single-float form sets Y gravity; triple-float form sets full gravity vector.

---

## Size

### `PARTICLE3D.SETSIZE(emitter, minSize, maxSize)` / `PARTICLE3D.SETSTARTSIZE(emitter, minSize, maxSize)` 

Sets the random starting size range.

---

### `PARTICLE3D.SETENDSIZE(emitter, minSize, maxSize)` 

Sets the random ending size at lifetime end (particles interpolate from start to end size).

---

## Color

### `PARTICLE3D.SETCOLOR(emitter, r, g, b, a)` / `PARTICLE3D.SETSTARTCOLOR(emitter, r, g, b, a)` 

Sets the starting color (0–255 per channel).

---

### `PARTICLE3D.SETCOLOREND(emitter, r, g, b, a)` / `PARTICLE3D.SETENDCOLOR(emitter, r, g, b, a)` 

Sets the ending color. Particles interpolate over lifetime.

---

### `PARTICLE3D.GETCOLOR(emitter)` 

Returns the current start color as a color handle.

---

### `PARTICLE3D.GETALPHA(emitter)` 

Returns the current start alpha as a float.

---

## Texture & Billboard

### `PARTICLE3D.SETTEXTURE(emitter, texHandle)` 

Sets a texture for particle sprites. `texHandle` from `TEXTURE.LOAD`.

---

### `PARTICLE3D.SETBILLBOARD(emitter, enabled)` 

When `TRUE` (default), particles always face the camera. Set `FALSE` for oriented quads.

---

## Playback

### `PARTICLE3D.PLAY(emitter)` 

Starts or restarts emission.

---

### `PARTICLE3D.STOP(emitter)` 

Stops new emission; existing particles continue until their lifetime expires.

---

### `PARTICLE3D.ISALIVE(emitter)` 

Returns `1` if any particles are still alive (useful to wait for all particles to fade before freeing).

---

### `PARTICLE3D.COUNT(emitter)` 

Returns the number of currently live particles.

---

## Update & Draw

### `PARTICLE3D.UPDATE(emitter, dt)` 

Advances all live particles by `dt` seconds.

---

### `PARTICLE3D.DRAW(emitter)` / `PARTICLE3D.DRAW(emitter, cam)` 

Draws all live particles. If `cam` is provided, uses that camera for billboarding.

---

## Lifetime

### `PARTICLE3D.FREE(emitter)` 

Destroys the emitter and all its particles.

---

## Full Example

A campfire effect — upward orange sparks fading to transparent.

```basic
WINDOW.OPEN(960, 540, "Particle3D Demo")
WINDOW.SETFPS(60)

cam = CAMERA.CREATE()
CAMERA.SETPOS(cam, 0, 4, -8)
CAMERA.SETTARGET(cam, 0, 2, 0)

fire = PARTICLE3D.CREATE()
PARTICLE3D.SETPOS(fire, 0, 0, 0)
PARTICLE3D.SETRATE(fire, 40)
PARTICLE3D.SETLIFETIME(fire, 0.6, 1.2)
PARTICLE3D.SETSPEED(fire, 1.5, 3.0)
PARTICLE3D.SETDIRECTION(fire, 0, 1, 0)
PARTICLE3D.SETSPREAD(fire, 20)
PARTICLE3D.SETSTARTSIZE(fire, 0.15, 0.3)
PARTICLE3D.SETENDSIZE(fire, 0.0, 0.05)
PARTICLE3D.SETSTARTCOLOR(fire, 255, 160, 30, 255)
PARTICLE3D.SETENDCOLOR(fire, 200, 50, 10, 0)
PARTICLE3D.SETGRAVITY(fire, 0.2)
PARTICLE3D.PLAY(fire)

WHILE NOT WINDOW.SHOULDCLOSE()
    dt = TIME.DELTA()
    PARTICLE3D.UPDATE(fire, dt)

    RENDER.CLEAR(10, 10, 15)
    RENDER.BEGIN3D(cam)
        PARTICLE3D.DRAW(fire, cam)
        DRAW3D.GRID(10, 1.0)
    RENDER.END3D()
    RENDER.FRAME()
WEND

PARTICLE3D.STOP(fire)
PARTICLE3D.FREE(fire)
CAMERA.FREE(cam)
WINDOW.CLOSE()
```

---

## Extended Command Reference

| Command | Description |
|--------|-------------|
| `PARTICLE3D.MAKE(...)` | Deprecated alias of `PARTICLE3D.CREATE`. |
| `PARTICLE3D.SETPOSITION(p, x,y,z)` | Alias of `PARTICLE3D.SETPOS`. |

---

## See also

- [PARTICLE2D.md](PARTICLE2D.md) — 2D screen-space particles
- [PARTICLE.md](PARTICLE.md) — legacy particle emitters
- [TEXTURE.md](TEXTURE.md) — texture loading for particle sprites
- [EFFECT.md](EFFECT.md) — post-process effects (bloom, etc.)
