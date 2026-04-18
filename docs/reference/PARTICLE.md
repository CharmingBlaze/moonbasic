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
### `PARTICLE.CREATE()`
Creates a new 3D particle emitter. Returns an **emitter handle**.

- **Returns**: (Handle) The new emitter handle.
- **Example**:
    ```basic
    fire = PARTICLE.CREATE()
    ```

---

### `PARTICLE.SETTEXTURE(emitter, textureHandle)`
Binds a texture to the emitter particles.

- **Returns**: (Handle) The emitter handle (for chaining).

---

### `PARTICLE.SETEMITRATE(emitter, rate)`
Sets the number of particles spawned per second.

- **Returns**: (Handle) The emitter handle (for chaining).

---

### `PARTICLE.SETPOS(emitter, x, y, z)`
Sets the emitter world position.

- **Returns**: (Handle) The emitter handle (for chaining).

---

### `PARTICLE.SETLIFETIME(emitter, min, max)`
Sets the lifespan range for new particles.

- **Returns**: (Handle) The emitter handle (for chaining).

---

### `PARTICLE.SETVELOCITY(emitter, vx, vy, vz, spread)`
Sets the initial velocity and random spread.

- **Returns**: (Handle) The emitter handle (for chaining).

---

### `PARTICLE.UPDATE(emitter, dt)`
Advances the particle simulation.

- **Arguments**:
    - `emitter`: (Handle) The emitter to update.
    - `dt`: (Float) Delta time.
- **Returns**: (None)

---

### `PARTICLE.DRAW(emitter [, camera])`
Renders the particles in 3D space.

- **Returns**: (None)

---

### `PARTICLE.FREE(emitter)`
Frees the emitter and all its particles.

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
