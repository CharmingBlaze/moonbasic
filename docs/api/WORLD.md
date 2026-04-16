# World Commands

Commands for managing the game world: gravity, fog, sky color, world center, terrain streaming, vegetation, reflections, and the per-frame world update tick. The World namespace ties together physics, terrain, and atmospheric systems.

## Core Concepts

- **World.Update(dt)** — The per-frame world tick. Advances terrain streaming, vegetation LOD, and world-space systems. Call every frame in your game loop.
- **Gravity** — Sets the physics gravity vector via Jolt Physics (`PHYSICS3D.SETGRAVITY`).
- **Fog** — Distance-based color blending to simulate atmospheric depth (linear or exponential).
- **Streaming** — Large terrains are loaded/unloaded in chunks around a center point for open-world games.

---

## World Update

### `World.Update(dt)`

Advances all world systems by `dt` seconds. Call this every frame.

- `dt` (float) — Delta time (seconds since last frame). Use `Time.Delta()`.

**How it works:** Ticks terrain streaming, vegetation LOD recalculation, and any world-space managers (weather, clouds, etc.) that are attached.

```basic
WHILE NOT Window.ShouldClose()
    dt = Time.Delta()
    World.Update(dt)
    ; ... render ...
WEND
```

---

## Gravity

### `World.Gravity(gx, gy, gz)`

Sets the 3D physics gravity vector. Delegates directly to `PHYSICS3D.SETGRAVITY`.

- `gx`, `gy`, `gz` (float) — Gravity vector components. Typical Earth gravity: `(0, -9.81, 0)`.

```basic
World.Gravity(0, -9.81, 0)

; Low gravity moon level
World.Gravity(0, -1.6, 0)

; Zero-G space level
World.Gravity(0, 0, 0)
```

---

## Fog

### `World.FogMode(mode)`

Sets the fog rendering mode.

- `mode` (int) — 0 = Off, 1 = Linear, 2 = Exponential.

**How it works:** Linear fog blends between a start and end distance. Exponential fog uses a density curve. Fog color is separate from fog mode.

```basic
World.FogMode(1)                    ; Linear fog
World.FogColor(180, 180, 200)       ; Light grey-blue
World.FogDensity(0.02)
```

---

### `World.FogColor(r, g, b)`

Sets the fog color.

- `r`, `g`, `b` (int) — Color (0–255).

---

### `World.FogDensity(density)`

Sets the fog density for exponential fog mode.

- `density` (float) — Fog density (e.g., 0.01–0.1).

---

## Sky

### `SkyColor(r, g, b)`

Alias for `Render.Clear`. Clears the screen to a sky color — used in the Blitz3D rendering pattern where `SkyColor` is called at the start of each frame.

```basic
SkyColor(100, 180, 255)   ; Clear blue sky
```

---

## World Center & Streaming

### `World.SetCenter(x, z)`

Sets the world center position for terrain streaming and LOD calculations. Terrain chunks near the center are loaded at full detail.

- `x`, `z` (float) — World-space XZ position.

```basic
World.SetCenter(playerX, playerZ)
```

---

### `World.SetCenterEntity(entityHandle)`

Sets the world center to an entity's position. Automatically tracks the entity.

- `entityHandle` (handle) — Entity to follow.

```basic
World.SetCenterEntity(player)
```

---

### `World.StreamEnable(enabled)`

Enables or disables terrain chunk streaming. When enabled, terrain chunks load/unload based on distance from the world center.

- `enabled` (bool) — `TRUE` to enable streaming.

---

### `World.Preload(radius)`

Forces preloading of terrain chunks within a radius. Use during loading screens to avoid pop-in.

- `radius` (int) — Chunk radius to preload.

```basic
World.Preload(5)   ; Preload 5 chunks in every direction
```

---

### `World.Status()`

Returns a string describing the current world loading status (e.g., number of chunks loaded, pending).

**Returns:** `string`

---

### `World.IsReady()`

Returns `TRUE` when all queued terrain chunks have finished loading.

**Returns:** `bool`

```basic
World.Preload(3)
WHILE NOT World.IsReady()
    Render.Clear(0, 0, 0)
    Draw.Text("Loading world: " + World.Status(), 10, 10, 20, 255, 255, 255, 255)
    Render.Frame()
WEND
```

---

## Environment

### `World.SetVegetation(density)`

Sets the vegetation density for procedural grass/plant rendering.

- `density` (float) — Density multiplier.

---

### `World.SetReflection(mode)`

Enables or configures environment reflections for water and shiny surfaces.

- `mode` (int) — Reflection mode.

---

## Time Scale

### `World.SetTimeScale(scale)` / `Game.SetTimeScale(scale)`

Sets a global time multiplier that affects `Time.Delta()`. Use for slow-motion effects.

- `scale` (float) — 1.0 = normal, 0.5 = half speed, 2.0 = double speed, 0.0 = paused.

**How it works:** Modifies `rt.TimeScale` on the runtime. All calls to `Time.Delta()` are scaled by this value. Does not affect real-time timers or physics step rate — only the `dt` value returned to scripts.

```basic
; Slow-motion hit effect
World.SetTimeScale(0.25)
SLEEP 500
World.SetTimeScale(1.0)
```

---

### `Game.GetTimeScale()`

Returns the current time scale.

**Returns:** `float`

---

## Visual Effects

### `World.Flash(r, g, b, a, duration)`

Flashes the entire screen with a color that fades out over the given duration. Used for damage feedback, explosions, or transitions.

- `r`, `g`, `b`, `a` (int) — Flash color.
- `duration` (float) — Duration in seconds.

**How it works:** Stores the flash color and duration. During `Render.Frame`, a full-screen rectangle is drawn with decreasing alpha each frame until the duration expires.

```basic
; Red damage flash
World.Flash(255, 0, 0, 180, 0.3)

; White explosion flash
World.Flash(255, 255, 255, 255, 0.15)
```

---

### `World.HitStop(duration)`

Freezes game time for a brief moment (hit-stop / impact frame). During the hit-stop, `Time.Delta()` returns 0, creating a dramatic pause.

- `duration` (float) — Duration in seconds (typically 0.05–0.2).

**How it works:** Sets `rt.HitStopEndAt` to the current wall clock plus the duration. `DeltaSeconds` returns 0 until the hit-stop expires.

```basic
; Impact frame on a heavy attack
World.HitStop(0.08)
Camera.Shake(cam, 3.0, 0.2)
```

---

## Easy Mode Shortcuts

| Shortcut | Maps To |
|----------|---------|
| `SKYCOLOR(r, g, b)` | `Render.Clear(r, g, b)` |
| `FOGMODE(m)` | `World.FogMode(m)` |
| `FOGCOLOR(r, g, b)` | `World.FogColor(r, g, b)` |
| `FOGDENSITY(d)` | `World.FogDensity(d)` |
| `UPDATEPHYSICS` | `World.Update + Physics step` |
| `DELTATIME` | `Time.Delta()` |

---

## Full Example

An open-world scene with terrain streaming, fog, and time-scale effects.

```basic
Window.Open(1280, 720, "World Demo")
Window.SetFPS(60)

cam = Camera.Create()
cam.pos(0, 20, 40)
cam.look(0, 0, 0)
cam.fov(60)

; Setup world atmosphere
World.Gravity(0, -9.81, 0)
World.FogMode(1)
World.FogColor(180, 190, 210)
World.FogDensity(0.015)

; Enable terrain streaming
World.StreamEnable(TRUE)
World.Preload(3)

; Wait for terrain to load
WHILE NOT World.IsReady()
    Render.Clear(0, 0, 0)
    Draw.Text("Loading: " + World.Status(), 10, 10, 24, 255, 255, 255, 255)
    Render.Frame()
WEND

playerX = 0
playerZ = 0

WHILE NOT Window.ShouldClose()
    dt = Time.Delta()

    ; Move player
    IF Input.KeyDown(KEY_W) THEN playerZ = playerZ - 20 * dt
    IF Input.KeyDown(KEY_S) THEN playerZ = playerZ + 20 * dt
    IF Input.KeyDown(KEY_A) THEN playerX = playerX - 20 * dt
    IF Input.KeyDown(KEY_D) THEN playerX = playerX + 20 * dt

    ; Update world center for streaming
    World.SetCenter(playerX, playerZ)
    World.Update(dt)

    ; Slow-motion toggle
    IF Input.KeyPressed(KEY_TAB) THEN
        IF Game.GetTimeScale() < 1.0 THEN
            World.SetTimeScale(1.0)
        ELSE
            World.SetTimeScale(0.3)
        ENDIF
    ENDIF

    ; Render
    Render.Clear(100, 160, 220)
    Camera.Begin(cam)
        Draw.Grid(40, 2.0)
        Draw.Cube(playerX, 1, playerZ, 2, 2, 2, 200, 80, 80, 255)
    Camera.End(cam)

    Draw.Text("WASD = Move | TAB = Slow-mo", 10, 10, 18, 255, 255, 255, 255)
    Draw.Text("Time Scale: " + STR(Game.GetTimeScale()), 10, 32, 18, 200, 200, 200, 255)
    Render.Frame()
WEND

Camera.Free(cam)
Window.Close()
```

---

## See Also

- [RENDER](RENDER.md) — Frame lifecycle and fog settings
- [CAMERA](CAMERA.md) — Camera used for world rendering
- [PHYSICS](PHYSICS.md) — Gravity and physics step
- [TERRAIN](TERRAIN.md) — Terrain generation and streaming
- [WEATHER](WEATHER.md) — Weather and atmospheric effects
