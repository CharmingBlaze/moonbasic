# Particle Commands

Commands for creating and controlling 3D particle emitters. Particles are small, GPU-rendered billboards that simulate fire, smoke, sparks, explosions, rain, dust, magic effects, and other visual phenomena.

## Core Concepts

- **Emitter** — A point in 3D space that spawns particles over time. Each emitter has its own configuration (color, speed, life, gravity, etc.).
- **Particle** — A single billboard quad spawned by an emitter. Lives for a set duration, then dies.
- **Burst** — Spawns a batch of particles instantly (for explosions, impacts).
- **Continuous** — Spawns particles at a steady rate (for fire, smoke, trails).
- Emitters are **handles** and must be freed when no longer needed.

---

## Creation

### `Particle.Create(x, y, z)` / `Particle.Make(x, y, z)`

Creates a particle emitter at a world position.

- `x`, `y`, `z` (float) — World position.

**Returns:** `handle`

```basic
fireEmitter = Particle.Create(0, 0, 0)
```

---

### `Particle.Free(emitterHandle)`

Frees an emitter and all its particles.

- `emitterHandle` (handle) — Emitter to free.

---

## Position

### `Particle.SetPos(emitterHandle, x, y, z)` / `emitter.pos(x, y, z)`

Moves the emitter to a new position. New particles spawn from this point.

- `x`, `y`, `z` (float) — World position.

```basic
Particle.SetPos(fireEmitter, playerX, playerY + 1, playerZ)
```

---

## Particle Properties

### `Particle.SetColor(emitterHandle, r, g, b)` / `Particle.SetCol(emitterHandle, r, g, b, a)`

Sets the color of new particles.

- `r`, `g`, `b` (int) — Color (0–255).
- `a` (int) — Alpha (optional, 0–255).

```basic
Particle.SetColor(fireEmitter, 255, 150, 50)  ; Orange fire
```

---

### `Particle.SetLife(emitterHandle, seconds)`

Sets how long each particle lives before fading out.

- `seconds` (float) — Lifetime in seconds.

```basic
Particle.SetLife(fireEmitter, 1.5)
```

---

### `Particle.SetSpeed(emitterHandle, speed)`

Sets the initial speed of particles when spawned.

- `speed` (float) — Speed in world units per second.

```basic
Particle.SetSpeed(fireEmitter, 3.0)
```

---

### `Particle.SetGravity(emitterHandle, gravity)`

Sets gravity applied to particles. Positive = falls down, negative = floats up.

- `gravity` (float) — Gravity force.

```basic
Particle.SetGravity(fireEmitter, -1.0)   ; Fire floats upward
Particle.SetGravity(sparkEmitter, 5.0)   ; Sparks fall fast
```

---

### `Particle.SetSize(emitterHandle, size)`

Sets the size of particle billboards.

- `size` (float) — Size in world units.

---

### `Particle.SetRate(emitterHandle, particlesPerSecond)`

Sets the continuous emission rate.

- `particlesPerSecond` (float) — How many particles spawn per second.

```basic
Particle.SetRate(smokeEmitter, 20)
```

---

### `Particle.SetSpread(emitterHandle, angle)`

Sets the cone angle for particle emission direction.

- `angle` (float) — Spread angle in degrees. 0 = straight up, 180 = all directions.

---

### `Particle.SetTexture(emitterHandle, textureHandle)`

Sets a custom texture for particles (instead of the default white quad).

- `textureHandle` (handle) — Texture handle.

---

### `Particle.SetBlend(emitterHandle, mode)`

Sets the blend mode for particles. Additive blending (mode 1) creates a glow effect.

- `mode` (int) — 0 = Alpha, 1 = Additive.

```basic
Particle.SetBlend(fireEmitter, 1)   ; Additive = glowing fire
```

---

## Emission

### `Particle.Emit(emitterHandle)` / `emitter.emit()`

Triggers continuous emission. Particles spawn at the configured rate.

---

### `Particle.Burst(emitterHandle, count)` / `emitter.burst(count)`

Spawns `count` particles instantly.

- `count` (int) — Number of particles to spawn.

```basic
; Explosion burst
Particle.Burst(explosionEmitter, 50)
```

---

### `Particle.Stop(emitterHandle)`

Stops continuous emission. Existing particles continue their lifecycle.

---

## Rendering

### `Particle.Draw(emitterHandle)` / `Particle.DrawAll()`

Draws particles. Must be inside a camera block.

```basic
Camera.Begin(cam)
    Entity.DrawAll()
    Particle.DrawAll()
Camera.End(cam)
```

---

## Full Example

A fire and explosion particle system.

```basic
Window.Open(1280, 720, "Particle Demo")
Window.SetFPS(60)

cam = Camera.Create()
cam.pos(0, 5, 10)
cam.look(0, 2, 0)
cam.fov(60)

; Fire emitter
fire = Particle.Create(0, 0.5, 0)
Particle.SetColor(fire, 255, 150, 50)
Particle.SetLife(fire, 1.0)
Particle.SetSpeed(fire, 2.0)
Particle.SetGravity(fire, -1.5)
Particle.SetSize(fire, 0.3)
Particle.SetRate(fire, 30)
Particle.SetSpread(fire, 15)
Particle.SetBlend(fire, 1)
Particle.Emit(fire)

; Explosion emitter (burst only)
explosion = Particle.Create(0, 0, 0)
Particle.SetColor(explosion, 255, 200, 100)
Particle.SetLife(explosion, 0.8)
Particle.SetSpeed(explosion, 8.0)
Particle.SetGravity(explosion, 3.0)
Particle.SetSize(explosion, 0.5)
Particle.SetSpread(explosion, 180)
Particle.SetBlend(explosion, 1)

WHILE NOT Window.ShouldClose()
    dt = Time.Delta()

    ; Trigger explosion on click
    IF Input.MousePressed(0) THEN
        Particle.SetPos(explosion, RND(-3, 3), 1, RND(-3, 3))
        Particle.Burst(explosion, 40)
        Camera.Shake(cam, 2.0, 0.2)
    ENDIF

    Render.Clear(20, 20, 30)
    Camera.Begin(cam)
        Draw.Grid(10, 1.0)
        Particle.DrawAll()
    Camera.End(cam)

    Draw.Text("Click to explode!", 10, 10, 20, 255, 255, 255, 255)
    Render.Frame()
WEND

Particle.Free(fire)
Particle.Free(explosion)
Camera.Free(cam)
Window.Close()
```

---

## See Also

- [DRAW](DRAW.md) — Immediate-mode drawing
- [ENTITY](ENTITY.md) — Attach emitters to entities
- [CAMERA](CAMERA.md) — Camera for rendering particles
- [RENDER](RENDER.md) — Blend modes for particle effects
