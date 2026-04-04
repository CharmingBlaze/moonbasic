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
| `SIN(angle#)` | Sine of `angle`. |
| `COS(angle#)` | Cosine of `angle`. |
| `TAN(angle#)` | Tangent of `angle`. |
| `ATN(angle#)` | Arctangent of `angle`. Alias: `ATAN`. |
| `ASIN(value#)` | Arcsine — returns an angle in radians. |
| `ACOS(value#)` | Arccosine — returns an angle in radians. |
| `ATAN2(y#, x#)` | Two-argument arctangent — returns the angle of the vector `(x, y)`. |

```basic
; Circular motion example
x# = COS(angle#) * radius#
y# = SIN(angle#) * radius#
```

---

## Powers, Roots & Logarithms

| Command | Returns |
|---|---|
| `SQRT(value#)` | Square root. Alias: `SQR`. |
| `POW(base#, exp#)` | `base` raised to the power of `exp`. |
| `EXP(value#)` | e raised to the power of `value`. |
| `LOG(value#)` | Natural logarithm (base e). |
| `LOG2(value#)` | Base-2 logarithm. |
| `LOG10(value#)` | Base-10 logarithm. |

---

## Rounding & Truncation

| Command | Returns |
|---|---|
| `FLOOR(value#)` | Largest integer ≤ `value`. |
| `CEIL(value#)` | Smallest integer ≥ `value`. |
| `ROUND(value#, [decimals])` | Nearest integer, or rounded to `decimals` places. |
| `INT(value#)` | Truncates toward zero (same as `FIX`). |
| `FIX(value#)` | Truncates toward zero — `FIX(-3.7)` = `-3`. |

---

## Arithmetic Helpers

| Command | Returns |
|---|---|
| `ABS(value#)` | Absolute value. |
| `SGN(value#)` | Sign: returns `-1`, `0`, or `1`. |
| `MIN(a#, b#)` | The smaller of two values. |
| `MAX(a#, b#)` | The larger of two values. |
| `CLAMP(value#, min#, max#)` | Constrains `value` to the range [min, max]. |

---

## Interpolation & Easing

| Command | Returns |
|---|---|
| `LERP(a#, b#, t#)` | Linear interpolation between `a` and `b` by `t` (0.0–1.0). |
| `SMOOTHSTEP(lo#, hi#, x#)` | Smooth S-curve interpolation — useful for easing animations. |
| `PINGPONG(t#, length#)` | Bounces `t` back and forth between 0 and `length`. |
| `WRAP(value#, min#, max#)` | Wraps `value` within [min, max] (like a modulo that respects bounds). |

```basic
; Smooth follow camera
cam_x# = LERP(cam_x#, target_x#, 5.0 * Time.Delta())
```

---

## Angles

| Command | Returns |
|---|---|
| `DEG2RAD(degrees#)` | Converts degrees to radians. |
| `RAD2DEG(radians#)` | Converts radians to degrees. |
| `WRAPANGLE(angle#)` | Normalizes an angle to the range [0, 360). |
| `WRAPANGLE180(angle#)` | Normalizes an angle to the range [-180, 180). |
| `ANGLEDIFF(from#, to#)` | Returns the shortest signed angle difference in degrees. |

```basic
; Rotate toward a target, never overshooting
diff# = ANGLEDIFF(facing#, desired#)
facing# = facing# + CLAMP(diff#, -5, 5)
```

---

## Randomization

The RNG is seeded from the system clock at startup. Use `RNDSEED` / `RANDOMIZE`
for reproducible sequences.

| Command | Returns |
|---|---|
| `RND([limit])` | Random integer 0 to `limit-1`. With no argument, a float in [0, 1). |
| `RNDF(min#, max#)` | Random float in [min, max]. |
| `RNDSEED(seed)` | Seeds the RNG with a specific integer value. |
| `RANDOMIZE([seed])` | Seeds from `seed`, or from the system clock if omitted. |

```basic
; Roll a six-sided die
die = RND(6) + 1

; Random float speed
speed# = RNDF(80.0, 160.0)

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

angle# = 0.0
radius# = 150.0
center_x = 400
center_y = 300

WHILE NOT Window.ShouldClose()
    angle# = angle# + 2.0 * Time.Delta()

    x = center_x + INT(COS(angle#) * radius#)
    y = center_y + INT(SIN(angle#) * radius#)

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

FOR i = 0 TO STAR_COUNT - 1
    sx(i) = RND(800)
    sy(i) = RND(600)
    ss(i) = RND(3) + 1
NEXT

WHILE NOT Window.ShouldClose()
    Render.Clear(0, 0, 10)
    FOR i = 0 TO STAR_COUNT - 1
        bright = 150 + RND(105)
        Draw.Rectangle(sx(i), sy(i), ss(i), ss(i), bright, bright, bright, 255)
    NEXT
    Render.Frame()
WEND

Window.Close()
```
