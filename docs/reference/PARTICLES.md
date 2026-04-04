# Particle System Commands

Commands for creating and simulating 2D and 3D particle effects.

The particle system in moonBASIC is emitter-based. You create an emitter, configure its
behaviour (rate, lifetime, velocity, color, size), and then update and draw it each frame.

## Core Workflow

1. **Create**: `Particle.Make(maxCount)` — reserves space for up to `maxCount` live particles.
2. **Configure**: Set rate, lifetime, velocity, color, etc.
3. **Start**: `Particle.Play(handle)` — begins emitting.
4. **Loop**:
   - `Particle.SetPos(handle, x#, y#, z#)` — move the emitter each frame if needed.
   - `Particle.Update(handle, dt#)` — advance the simulation.
   - `Particle.Draw(handle)` — render the particles.
5. **Free**: `Particle.Free(handle)` when you are done.

---

## Creating & Freeing

### `Particle.Make(maxCount)`

Creates a new particle emitter. Returns a handle.

- `maxCount`: The maximum number of simultaneous live particles. Keep this
  reasonable — 500–2000 is typical for on-screen effects; up to 16000 is the
  system maximum.

```basic
fire = Particle.Make(500)
```

---

### `Particle.Free(handle)`

Frees the emitter and all its particle data.

---

## Configuration

Configure the emitter **before** calling `Particle.Play()`. Properties can also
be changed at runtime to create evolving effects.

---

### `Particle.SetTexture(handle, texHandle)`

Sets the sprite texture drawn for each particle. If not set, particles render
as solid white squares.

- `texHandle`: A texture handle from `Texture.Load()`.

```basic
spark_tex = Texture.Load("assets/spark.png")
Particle.SetTexture(fire, spark_tex)
```

---

### `Particle.SetEmitRate(handle, rate#)`

Sets how many particles are spawned per second.

- `rate#`: Particles per second (e.g., `60.0`).

```basic
Particle.SetEmitRate(fire, 80.0)
```

---

### `Particle.SetLifetime(handle, min#, max#)`

Sets the range for how long each particle lives, in seconds. Each particle gets
a random lifetime in this range.

```basic
Particle.SetLifetime(fire, 0.5, 1.2)
```

---

### `Particle.SetVelocity(handle, vx#, vy#, vz#, spread#)`

Sets the base emission direction and the angular spread (in degrees). A spread
of `0` creates a tight beam; `360` creates a full sphere.

```basic
; Emit mostly upward with some spread
Particle.SetVelocity(fire, 0, -120, 0, 30.0)
```

---

### `Particle.SetColor(handle, r, g, b, a)`

Sets the starting color of newly spawned particles (0–255 per channel).

```basic
Particle.SetColor(fire, 255, 140, 40, 255)  ; bright orange
```

---

### `Particle.SetColorEnd(handle, r, g, b, a)`

Sets the color a particle fades toward by the end of its lifetime.

```basic
Particle.SetColorEnd(fire, 100, 0, 0, 0)  ; fade to dark red, transparent
```

---

### `Particle.SetSize(handle, startSize#, endSize#)`

Sets the particle size in pixels (or world units for 3D) at birth and at death.

```basic
Particle.SetSize(fire, 12.0, 2.0)  ; shrinks as it ages
```

---

### `Particle.SetGravity(handle, gx#, gy#, gz#)`

Sets the gravitational acceleration applied to this emitter's particles
independently of global physics.

```basic
; Float upward, ignore real gravity
Particle.SetGravity(fire, 0, -20, 0)
```

---

## Positioning

### `Particle.SetPos(handle, x#, y#, z#)`

Sets the world-space position where new particles spawn. Update this each frame
if the emitter is attached to a moving object.

```basic
; Attach fire to the player's position
Particle.SetPos(fire, player_x#, player_y#, 0)
```

---

## Playback

### `Particle.Play(handle)`

Starts the emitter. Particles begin spawning immediately.

### `Particle.Update(handle, dt#)`

Advances the particle simulation by `dt#` seconds. Call this once per frame.

- `dt#`: Delta time — use `Time.Delta()`.

### `Particle.Draw(handle)`

Renders all currently live particles. Place this call between `Render.Clear()`
and `Render.Frame()`.

---

## Full Example: Campfire Effect

```basic
Window.Open(960, 540, "Campfire Particles")
Window.SetFPS(60)

Audio.Init()

; --- SETUP ---
fire_tex = Texture.Load("assets/spark.png")

fire = Particle.Make(400)
Particle.SetTexture(fire, fire_tex)
Particle.SetEmitRate(fire, 60.0)
Particle.SetLifetime(fire, 0.4, 1.0)
Particle.SetVelocity(fire, 0, -80, 0, 25.0)
Particle.SetColor(fire, 255, 160, 40, 255)
Particle.SetColorEnd(fire, 80, 0, 0, 0)
Particle.SetSize(fire, 10.0, 1.0)
Particle.SetGravity(fire, 0, -30, 0)
Particle.SetPos(fire, 480, 400, 0)

ember = Particle.Make(150)
Particle.SetEmitRate(ember, 20.0)
Particle.SetLifetime(ember, 0.8, 2.0)
Particle.SetVelocity(ember, 0, -60, 0, 60.0)
Particle.SetColor(ember, 255, 200, 100, 255)
Particle.SetColorEnd(ember, 200, 100, 0, 0)
Particle.SetSize(ember, 4.0, 1.0)
Particle.SetGravity(ember, 0, -15, 0)
Particle.SetPos(ember, 480, 400, 0)

Particle.Play(fire)
Particle.Play(ember)

; --- MAIN LOOP ---
WHILE NOT Window.ShouldClose()
    dt# = Time.Delta()

    Particle.Update(fire, dt#)
    Particle.Update(ember, dt#)

    Render.Clear(10, 10, 15)
    Draw.Rectangle(0, 420, 960, 120, 30, 20, 10, 255)  ; ground
    Particle.Draw(fire)
    Particle.Draw(ember)
    Render.Frame()
WEND

; --- CLEANUP ---
Particle.Free(fire)
Particle.Free(ember)
Texture.Free(fire_tex)
Audio.Close()
Window.Close()
```

---

## Tips

- **Layer order matters**: Draw background particles before foreground sprites.
- **Burst effects**: Set a high `EmitRate`, call `Play`, then on the next frame
  set `EmitRate` to `0` — existing particles continue until their lifetime ends.
- **Performance**: Keep `maxCount` tight. 500 particles per emitter is plenty for
  most 2D effects; 3D effects may need more.
- **3D particles**: `SetPos` uses world-space XYZ. Draw the emitter inside
  `cam.Begin()` / `cam.End()` for 3D scenes.
