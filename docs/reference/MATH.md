# Math Commands

Commands for mathematical and numerical operations. All functions are available
both as plain names (`SIN`, `CLAMP`) and as `MATH.`-prefixed names (`MATH.SIN`,
`MATH.CLAMP`). The prefixed forms are useful when reading code and avoiding name
collisions.

---

## Trigonometry

Angles are in **radians** for all trig functions. Use `DEG2RAD` / `RAD2DEG` to
convert.

| Command | Returns |
|---|---|
| `SIN(angle)` | Sine of `angle`. |
| `COS(angle)` | Cosine of `angle`. |
| `TAN(angle)` | Tangent of `angle`. |
| `ATN(angle)` | Arctangent of `angle`. Alias: `ATAN`. |
| `ASIN(value)` | Arcsine — returns an angle in radians. |
| `ACOS(value)` | Arccosine — returns an angle in radians. |
| `ATAN2(y, x)` | Two-argument arctangent — returns the angle of the vector `(x, y)`. |

```basic
; Circular motion example
x = COS(angle) * radius
y = SIN(angle) * radius
```

---

## Powers, Roots & Logarithms

| Command | Returns |
|---|---|
| `SQRT(value)` | Square root. Alias: `SQR`. |
| `POW(base, exp)` | `base` raised to the power of `exp`. |
| `EXP(value)` | e raised to the power of `value`. |
| `LOG(value)` | Natural logarithm (base e). |
| `LOG2(value)` | Base-2 logarithm. |
| `LOG10(value)` | Base-10 logarithm. |

---

## Rounding & Truncation

| Command | Returns |
|---|---|
| `FLOOR(value)` | Largest integer ≤ `value`. |
| `CEIL(value)` | Smallest integer ≥ `value`. |
| `ROUND(value, [decimals])` | Nearest integer, or rounded to `decimals` places. |
| `INT(value)` | Truncates toward zero (same as `FIX`). |
| `FIX(value)` | Truncates toward zero — `FIX(-3.7)` = `-3`. |

---

## Arithmetic Helpers

| Command | Returns |
|---|---|
| `ABS(value)` | Absolute value. |
| `SGN(value)` | Sign: returns `-1`, `0`, or `1`. |
| `MIN(a, b)` | The smaller of two values. |
| `MAX(a, b)` | The larger of two values. |
| `CLAMP(value, min, max)` | Constrains `value` to the range [min, max]. |

`CLAMP` is available as both `CLAMP(...)` and `MATH.CLAMP(...)`.

---

## Camera-relative movement (`MOVEX` / `MOVEZ`)

Helpers for **third-person** or **orbit-camera** movement on the **XZ plane**. Pass the camera **yaw** in **radians**, plus **forward** and **strafe** inputs (typically `-1` … `1` from **`Input.Axis`**).

| Command | Returns |
|---|---|
| `MOVEX(yaw, forward, strafe)` | X world component: `-SIN(yaw)*forward + COS(yaw)*strafe`. |
| `MOVEZ(yaw, forward, strafe)` | Z world component: `-COS(yaw)*forward - SIN(yaw)*strafe`. |
| `MOVESTEPX(yaw, forward, strafe, speed, dt)` | Same as **`MOVEX(...)*speed*dt`** — world **X** delta for one frame. |
| `MOVESTEPZ(yaw, forward, strafe, speed, dt)` | Same as **`MOVEZ(...)*speed*dt`** — world **Z** delta for one frame. |

```basic
px = px + MOVEX(camYaw, f, s) * speed * dt
pz = pz + MOVEZ(camYaw, f, s) * speed * dt

; or one call per axis (bundles × speed × dt):
px = px + MOVESTEPX(camYaw, f, s, speed, dt)
pz = pz + MOVESTEPZ(camYaw, f, s, speed, dt)
```

Together they match **forward = −Z** and **right = +X** when **`yaw = 0`**. See [INPUT.md](INPUT.md) and [CAMERA.md](CAMERA.md).

---

## Inline conditional (`IIF`)

| Command | Returns |
|---|---|
| `IIF(condition, trueVal, falseVal)` | `trueVal` if `condition` is true, else `falseVal`. **Both branches are evaluated** (not short-circuit). |
| `IIF(condition, true, false)` | Same for strings. |

Use for compact HUD colours or labels; for side effects, prefer **`IF` / `ENDIF`**.

---

## Interpolation & Easing

| Command | Returns |
|---|---|
| `LERP(a, b, t)` | Linear interpolation between `a` and `b` by `t` (0.0–1.0). |
| `REMAP(value, inMin, inMax, outMin, outMax)` | Maps `value` linearly from the input range to the output range. If `inMin` = `inMax`, returns `outMin` (avoids divide-by-zero). |
| `INVERSE_LERP(a, b, x)` | Returns `(x - a) / (b - a)` — where `x` sits between endpoints (not clamped). If `a` = `b`, returns `0`. Use with `LERP` to convert between ranges. |
| `SATURATE(x)` | Clamps `x` to **[0, 1]** (same idea as GPU “saturate”). |
| `SMOOTHSTEP(lo, hi, x)` | Smooth S-curve interpolation — useful for easing animations. |
| `PINGPONG(t, length)` | Bounces `t` back and forth between 0 and `length`. |
| `WRAP(value, min, max)` | Wraps `value` within [min, max] (like a modulo that respects bounds). |

```basic
; Smooth follow camera
cam_x = LERP(cam_x, target_x, 5.0 * Time.Delta())

; UI slider 0..800 -> -1..1
t = INVERSE_LERP(0.0, 800.0, mouseX)
t = SATURATE(t)
pan = REMAP(t, 0.0, 1.0, -1.0, 1.0)
```

More gameplay-oriented shortcuts (`MATH.CIRCLEPOINT`, `MATH.APPROACH`, `MATH.LERPANGLE`, …) are listed in [LESS_MATH.md](LESS_MATH.md).

### Short game-logic helpers (engine-style names)

| Name | moonBASIC | Notes |
|------|-----------|--------|
| **Approach**(current, target, step) | **`MATH.APPROACH`** (also **`MOVE.TOWARD`**, **`APPROACH`**) | Move **current** toward **target** by at most **\|step\|** without overshooting. |
| **Lerp**(a, b, t) | **`MATH.LERP`** (also **`LERP`**, **`MOVE.LERP`**) | Linear interpolation; **t** usually 0…1 (Raylib-style clamp in implementation). |
| **Curve**(current, target, divisor) | **`MATH.CURVE`** (also **`CURVE`**) | Smooth easing: **`current + (target-current)/divisor`**; if divisor is below 1 it is treated as 1. Same idea as **`CURVEVALUE`** (argument order differs — see manifest). |
| **NewX**(currentX, angle, distance) | **`MATH.NEWX`** | **`currentX + MOVEX(angle, 1, 0) * distance`** — **angle** is **yaw in radians** (same as **`MOVEX`/`MOVEZ`**). |
| **NewZ**(currentZ, angle, distance) | **`MATH.NEWZ`** | **`currentZ + MOVEZ(angle, 1, 0) * distance`** — **radians**. |
| **AngleDiff**(a, b) | **`MATH.ANGLEDIFF`** (also **`ANGLEDIFF`**, **`ANGLE.DIFFERENCE`**) | Shortest signed delta **a → b** in **degrees** (−180…180). For radians use **`MATH.ANGLEDIFFRAD`**. |

2D **degree**-based steps (**cos/sin**) are still available as **`NEWXVALUE`** / **`NEWYVALUE`** in the game module if you need that convention instead of XZ + radians.

---

## Angles

| Command | Returns |
|---|---|
| `DEG2RAD(degrees)` | Converts degrees to radians. |
| `RAD2DEG(radians)` | Converts radians to degrees. |
| `WRAPANGLE(angle)` | Normalizes an angle to the range [0, 360). |
| `WRAPANGLE180(angle)` | Normalizes an angle to the range [-180, 180). |
| `ANGLEDIFF(from, to)` | Returns the shortest signed angle difference in degrees. |

```basic
; Rotate toward a target, never overshooting
diff = ANGLEDIFF(facing, desired)
facing = facing + CLAMP(diff, -5, 5)
```

---

## Randomization

The RNG is seeded from the system clock at startup. Use `RNDSEED` / `RANDOMIZE`
for reproducible sequences.

| Command | Returns |
|---|---|
| `RND()` | Random float in **[0, 1)**. |
| `RND(limit)` | Random integer **0 .. limit−1** (for `limit` ≥ 1). |
| `RND(lo, hi)` / `RAND(lo, hi)` | **Inclusive** random integer in **[lo, hi]** (Blitz-style **`Rand(min, max)`**). **`MATH.RAND`** is the same. |
| `RNDF(min, max)` | Random float in [min, max]. |
| `RNDSEED(seed)` | Seeds the RNG with a specific integer value. |
| `RANDOMIZE([seed])` | Seeds from `seed`, or from the system clock if omitted. |

**Smooth transitions (Blitz-style):** **`CURVEVALUE(dest, current, speed)`** and **`CURVEANGLE(...)`** are registered in the **`game`** module — see [GAMEHELPERS.md](GAMEHELPERS.md) and [BLITZ_ESSENTIAL_API.md](BLITZ_ESSENTIAL_API.md).

```basic
; Roll a six-sided die
die = RND(6) + 1

; Inclusive range (e.g. 10..20)
n = RAND(10, 20)

; Random float speed
speed = RNDF(80.0, 160.0)

; Reproducible map generation
RNDSEED(12345)
```

---

## Constants

| Command | Returns |
|---|---|
| `PI()` | π ≈ 3.14159265 |
| `TAU()` | τ = 2π ≈ 6.28318530 |
| `E()` | Euler's number ≈ 2.71828182 |

---

## Full Example: Circular Motion

```basic
Window.Open(800, 600, "Math Example: Circular Motion")
Window.SetFPS(60)

angle = 0.0
radius = 150.0
center_x = 400
center_y = 300

WHILE NOT Window.ShouldClose()
    angle = angle + 2.0 * Time.Delta()

    x = center_x + INT(COS(angle) * radius)
    y = center_y + INT(SIN(angle) * radius)

    Render.Clear(20, 20, 20)
    Draw.Rectangle(x - 15, y - 15, 30, 30, 200, 50, 150, 255)
    Render.Frame()
WEND

Window.Close()
```

---

## Full Example: Randomized Star Field

```basic
Window.Open(800, 600, "Star Field")
Window.SetFPS(60)

RANDOMIZE  ; Seed from system clock

CONST STAR_COUNT = 200
DIM sx(STAR_COUNT)
DIM sy(STAR_COUNT)
DIM ss(STAR_COUNT)

FOR i = 1 TO STAR_COUNT
    sx(i) = RND(800)
    sy(i) = RND(600)
    ss(i) = RND(3) + 1
NEXT

WHILE NOT Window.ShouldClose()
    Render.Clear(0, 0, 10)
    FOR i = 1 TO STAR_COUNT
        bright = 150 + RND(105)
        Draw.Rectangle(sx(i), sy(i), ss(i), ss(i), bright, bright, bright, 255)
    NEXT
    Render.Frame()
WEND

Window.Close()
```

---

## See also (gameplay-oriented)

- [GAME_MATH_HELPERS.md](GAME_MATH_HELPERS.md) — **`HDIST` / `HDISTSQ`** (XZ distance), **`DIST2D` / `DISTSQ2D`**, **`YAWFROMXZ`**, **`ANGLEDIFFRAD`**, **`SMOOTHERSTEP`**
- [GAME_ENGINE_PATTERNS.md](GAME_ENGINE_PATTERNS.md) — **`MOVE.TOWARD`**, **`MOVE.LERP`**, **`ANGLE.DIFFERENCE`**, rays, lights, sprites, **`RES.*`**
- [LESS_MATH.md](LESS_MATH.md) — camera-relative movement, terrain snap, vector helpers
