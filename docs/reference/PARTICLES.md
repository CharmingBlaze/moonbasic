# Particle system commands

Emitter-based **3D** particles (**`PARTICLE.*`** and alias **`PARTICLE3D.*`**) and **2D** particles (**`PARTICLE2D.*`**) — see [PARTICLE.md](PARTICLE.md) and [SPRITE.md](SPRITE.md).

---

## 3D workflow

### `Particles.Make(maxCount)`
Creates a new particle system with a fixed maximum number of particles. Returns an **emitter handle**.

### `Particles.Free(handle)`
Frees the particle system and all associated resources from memory.

### `Particles.Emit(handle, x, y, z, count)`
Emits a burst of particles at the specified world position.

### `Particles.Update(handle, dt)`
Updates the positions and lifetimes of all active particles based on elapsed time.

### `Particles.Draw(handle)`
Renders all active particles. This must be called within a **`Camera.Begin()`** / **`Camera.End()`** block.

### `Particle.SetEmitRate(handle, rate)`
Sets the number of particles emitted per second.

### `Particle.SetLifetime(handle, min, max)`
Sets the minimum and maximum lifetime (in seconds) for newly emitted particles.

### `Particle.SetVelocity(handle, vx, vy, vz, spread)`
Sets the initial emission velocity vector and randomness spread.

### `Particle.SetColor(handle, r, g, b, a)`
Sets the starting color for newly emitted particles.

### `Particle.SetSize(handle, start, end)`
Sets the starting and ending sizes for particles over their lifetime.

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
PARTICLE.SETVELOCITY p, 0, 1, 0, 0.4
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
    dt = Time.Delta()
    PARTICLE.UPDATE p, dt
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
