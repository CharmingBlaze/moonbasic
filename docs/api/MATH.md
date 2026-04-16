# Math, Easing, Noise & Utility Commands

Commands for mathematical operations, easing functions, noise generation, interpolation, and game-math helpers. These are essential for smooth animations, procedural content, camera movement, and gameplay logic.

## Core Math

moonBASIC supports standard math operators and built-in functions:

| Function | Description |
|----------|-------------|
| `SIN(deg)` | Sine (degrees) |
| `COS(deg)` | Cosine (degrees) |
| `TAN(deg)` | Tangent (degrees) |
| `ASIN(v)` | Arc sine → degrees |
| `ACOS(v)` | Arc cosine → degrees |
| `ATAN(v)` | Arc tangent → degrees |
| `ATAN2(y, x)` | Two-argument arc tangent → degrees |
| `SQRT(v)` | Square root |
| `ABS(v)` | Absolute value |
| `SGN(v)` | Sign (-1, 0, or 1) |
| `MIN(a, b)` | Minimum |
| `MAX(a, b)` | Maximum |
| `CLAMP(v, min, max)` | Clamp value to range |
| `FLOOR(v)` | Round down |
| `CEIL(v)` | Round up |
| `INT(v)` | Truncate to integer |
| `FLOAT(v)` | Convert to float |
| `MOD(a, b)` | Modulo |
| `POW(base, exp)` | Power |
| `LOG(v)` | Natural logarithm |
| `LOG10(v)` | Base-10 logarithm |
| `RND(min, max)` | Random float in range |
| `RAND(min, max)` | Random integer in range |
| `SEEDRND(seed)` | Seed the random generator |
| `PI` | 3.14159... constant |

---

## Game Math Helpers

### `CurveValue(destination, current, speed)`

Smoothly interpolates `current` toward `destination` using an exponential curve. Higher `speed` = faster convergence.

- `destination` (float) — Target value.
- `current` (float) — Current value.
- `speed` (float) — Curve speed (1.0–100.0 typical).

**Returns:** `float`

**How it works:** Computes `current + (destination - current) / speed` each frame. When called every frame, this produces smooth exponential decay toward the target.

```basic
; Smooth camera follow
camX = CurveValue(playerX, camX, 10)
camY = CurveValue(playerY, camY, 10)
```

---

### `CurveAngle(destination, current, speed)`

Same as `CurveValue` but handles angle wrapping (0–360 degrees). Smoothly interpolates between angles taking the shortest path.

- `destination` (float) — Target angle in degrees.
- `current` (float) — Current angle in degrees.
- `speed` (float) — Curve speed.

**Returns:** `float`

```basic
; Smooth rotation toward target
facing = CurveAngle(targetAngle, facing, 8)
```

---

### `WrapValue(value, min, max)`

Wraps a value into a range. If it exceeds `max`, it wraps to `min` and vice versa.

- `value` (float) — Input.
- `min`, `max` (float) — Range.

**Returns:** `float`

```basic
angle = WrapValue(angle + rotSpeed * dt, 0, 360)
```

---

### `Approach(current, target, maxDelta)`

Moves `current` toward `target` by at most `maxDelta`. Linear interpolation with clamping.

- `current` (float) — Current value.
- `target` (float) — Target value.
- `maxDelta` (float) — Maximum step size.

**Returns:** `float`

```basic
health = Approach(health, maxHealth, healRate * dt)
```

---

### `Oscillate(speed, min, max)`

Returns a value that oscillates between `min` and `max` over time using a sine wave.

- `speed` (float) — Oscillation speed.
- `min`, `max` (float) — Value range.

**Returns:** `float`

**How it works:** Uses wall-clock elapsed time and sine: `min + (max - min) * (0.5 + 0.5 * sin(time * speed))`.

```basic
; Pulsing glow
alpha = Oscillate(3.0, 100, 255)
```

---

### `PointDir2D(x1, y1, x2, y2)`

Returns the angle in degrees from point 1 to point 2.

**Returns:** `float`

```basic
angleToEnemy = PointDir2D(playerX, playerY, enemyX, enemyY)
```

---

### `PointDir3D(x1, y1, z1, x2, y2, z2, axis)`

Returns the angle between two 3D points along a specific axis (`"x"`, `"y"`, or `"z"`).

**Returns:** `float`

---

### `NewXValue(x, angle, distance)` / `NewYValue(y, angle, distance)` / `NewZValue(z, angleX, angleY, distance)`

Computes a new position by moving from a point at an angle for a distance. Useful for projectile trajectories and direction-based movement.

```basic
; Move forward at facing angle
bulletX = NewXValue(x, facing, speed * dt)
bulletY = NewYValue(y, facing, speed * dt)
```

---

## Easing Functions

Easing functions take a `t` value from 0.0 to 1.0 and return a curved result. Use with timer fractions for smooth animations.

### Quadratic

| Function | Curve |
|----------|-------|
| `EaseIn(t)` | Slow start, fast end |
| `EaseOut(t)` | Fast start, slow end |
| `EaseInOut(t)` | Slow start and end |

### Cubic

| Function | Curve |
|----------|-------|
| `EaseIn3(t)` | Stronger cubic ease in |
| `EaseOut3(t)` | Stronger cubic ease out |
| `EaseInOut3(t)` | Stronger cubic ease in-out |

### Sine

| Function | Curve |
|----------|-------|
| `EaseInSine(t)` | Gentle sine ease in |
| `EaseOutSine(t)` | Gentle sine ease out |
| `EaseInOutSine(t)` | Gentle sine ease in-out |

### Special

| Function | Curve |
|----------|-------|
| `EaseInBack(t)` | Overshoots slightly before easing in |
| `EaseOutBack(t)` | Overshoots slightly at the end |
| `EaseInBounce(t)` | Bounces at the start |
| `EaseOutBounce(t)` | Bounces at the end |
| `EaseInElastic(t)` | Elastic spring at start |
| `EaseOutElastic(t)` | Elastic spring at end |

### Generic Easing Lerp

### `EaseLerp(a, b, t, easingName)`

Interpolates between `a` and `b` at time `t` using a named easing function.

- `a`, `b` (float) — Start and end values.
- `t` (float) — Progress (0.0–1.0).
- `easingName` (string) — Easing function name (e.g., `"easeInQuad"`, `"easeOutBounce"`).

**Returns:** `float`

```basic
; Smooth menu slide-in
t = Timer.Fraction(slideTimer)
menuX = EaseLerp(-300, 0, t, "easeOutBack")
```

---

## Noise Generation

### `Perlin(x, y)` / `Perlin(x, y, z)` / `Perlin(x, y, octaves, persistence, lacunarity)`

Generates Perlin noise at a position. Returns a value from -1.0 to 1.0.

- `x`, `y` (float) — 2D coordinates.
- `z` (float) — Optional third dimension.
- `octaves` (int) — Number of noise octaves (detail layers).
- `persistence` (float) — How much each octave contributes.
- `lacunarity` (float) — Frequency multiplier per octave.

**Returns:** `float`

```basic
; Terrain height
height = Perlin(x * 0.05, z * 0.05) * 10

; Detailed terrain with octaves
height = Perlin(x * 0.02, z * 0.02, 6, 0.5, 2.0) * 20
```

---

### `Simplex(x, y)` / `Simplex(x, y, z)`

Generates Simplex noise. Similar to Perlin but with fewer directional artifacts.

**Returns:** `float`

---

## Random

### `RND(min, max)`

Returns a random float between `min` and `max`.

```basic
spawnX = RND(-10, 10)
```

---

### `RAND(min, max)`

Returns a random integer between `min` and `max` (inclusive).

```basic
damage = RAND(10, 25)
```

---

### `SeedRnd(seed)`

Seeds the random number generator for reproducible results.

```basic
SeedRnd(12345)   ; Same sequence every time
```

---

## Color Utilities

### `Color.Create(r, g, b, a)` / `Color.Make(r, g, b, a)`

Creates a color handle.

**Returns:** `handle`

---

### `Color.Lerp(colorA, colorB, t)`

Interpolates between two colors.

**Returns:** `handle`

---

### `Color.FromHSV(h, s, v)` / `Color.ToHSV(colorHandle)`

Convert between RGB and HSV color spaces.

---

## Vec2 / Vec3 / Mat4

### `Vec2.Create(x, y)` / `Vec3.Create(x, y, z)`

Creates a vector handle.

### `Vec3.Add(a, b)` / `Vec3.Sub(a, b)` / `Vec3.Scale(v, s)` / `Vec3.Normalize(v)` / `Vec3.Dot(a, b)` / `Vec3.Cross(a, b)` / `Vec3.Length(v)` / `Vec3.Distance(a, b)` / `Vec3.Lerp(a, b, t)`

Vector math operations.

### `Mat4.Identity()` / `Mat4.Translate(x, y, z)` / `Mat4.RotateX(angle)` / `Mat4.RotateY(angle)` / `Mat4.RotateZ(angle)` / `Mat4.Scale(sx, sy, sz)` / `Mat4.Multiply(a, b)` / `Mat4.Invert(m)`

Matrix operations.

---

## Full Example

Smooth camera follow, easing animations, and noise-based terrain.

```basic
Window.Open(800, 600, "Math Demo")
Window.SetFPS(60)

cam = Camera.Create()
cam.fov(60)

; Generate noise terrain
playerX = 0
playerZ = 0
camX = 0
camY = 10
camZ = 20

WHILE NOT Window.ShouldClose()
    dt = Time.Delta()

    ; Move player
    IF Input.KeyDown(KEY_W) THEN playerZ = playerZ - 10 * dt
    IF Input.KeyDown(KEY_S) THEN playerZ = playerZ + 10 * dt
    IF Input.KeyDown(KEY_A) THEN playerX = playerX - 10 * dt
    IF Input.KeyDown(KEY_D) THEN playerX = playerX + 10 * dt

    ; Smooth camera follow (CurveValue)
    camX = CurveValue(playerX, camX, 8)
    camZ = CurveValue(playerZ + 20, camZ, 8)
    cam.pos(camX, camY, camZ)
    cam.look(playerX, 0, playerZ)

    ; Oscillating light
    lightIntensity = Oscillate(2.0, 0.5, 1.0)

    Render.Clear(30, 30, 50)
    Camera.Begin(cam)
        Draw.Grid(40, 2.0)

        ; Player cube
        Draw.Cube(playerX, 1, playerZ, 1, 2, 1, 100, 200, 255, 255)

        ; Noise-based terrain dots
        FOR gx = -10 TO 10
            FOR gz = -10 TO 10
                wx = gx * 2
                wz = gz * 2
                h = Perlin(wx * 0.1, wz * 0.1) * 3
                Draw.Cube(wx, h, wz, 1.8, 0.2, 1.8, 60 + INT(h * 20), 100, 60, 255)
            NEXT
        NEXT
    Camera.End(cam)

    Draw.Text("WASD = Move | Perlin terrain + CurveValue camera", 10, 10, 16, 255, 255, 255, 255)
    Render.Frame()
WEND

Camera.Free(cam)
Window.Close()
```

---

## See Also

- [TIME](TIME.md) — Timer fractions for easing animations
- [CAMERA](CAMERA.md) — Smooth camera follow
- [WORLD](WORLD.md) — Procedural terrain with noise
