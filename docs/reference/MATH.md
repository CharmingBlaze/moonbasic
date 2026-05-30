# Math Commands

Mathematical and numerical operations: trig, powers, clamping, lerp, random, and constants.

Page shape follows [DOC_STYLE_GUIDE.md](../DOC_STYLE_GUIDE.md) (**WAVE pattern**).

## Core Workflow

Use global aliases (`SIN`, `COS`, `SQRT`, `ABS`, `CLAMP`, `LERP`, `RND`, `PI`, etc.) or the `MATH.*` namespace equivalents. Angles are in **radians**.

---

## Trigonometry

Angles are in **radians**. Use `Math.Deg2Rad()` / `Math.Rad2Deg()` to convert.

### `MATH.SIN(angle)` / `COS` / `TAN`
Returns the trigonometric result of an angle in **radians**.

- **Arguments**:
    - `angle`: (Float) Angle in radians.
- **Returns**: (Float) The result.
- **Example**:
    ```basic
    y = SIN(angle) * radius
    ```

---

### `MATH.SQRT(value)`
Returns the square root of a non-negative value.

- **Returns**: (Float)

---

### `MATH.ABS(value)`
Returns the absolute value of a number.

- **Returns**: (Float/Integer)

---

### `MATH.POW(base, exp)`
Returns base raised to the power of exp.

- **Returns**: (Float)

---

### `MATH.CLAMP(value, min, max)`
Constrains a value to the range `[min, max]`.

- **Arguments**:
    - `value`: (Float) The number to clamp.
    - `min, max`: (Float) The range boundaries.
- **Returns**: (Float) The clamped value.

---

### `MATH.LERP(a, b, t)`
Linearly interpolates between `a` and `b` by factor `t`.

- **Arguments**:
    - `a, b`: (Float) Start and end values.
    - `t`: (Float) Interpolation factor (0.0 to 1.0).
- **Returns**: (Float) The interpolated value.

---

### `Rnd(limit)`
Returns a random integer from 0 up to `limit-1`.

- **Returns**: (Integer)

---

### `RndF(min, max)`
Returns a random float between `min` and `max`.

- **Returns**: (Float)

---

### `MATH.PI()` / `TAU` / `E`
Returns mathematical constants.

- **Returns**: (Float)

---

## Full Example: Circular Motion

```basic
WINDOW.OPEN(800, 600, "Math Example: Circular Motion")
WINDOW.SETFPS(60)

angle = 0.0
radius = 150.0
center_x = 400
center_y = 300

WHILE NOT WINDOW.SHOULDCLOSE()
    angle = angle + 2.0 * TIME.DELTA()

    x = center_x + INT(MATH.COS(angle) * radius)
    y = center_y + INT(MATH.SIN(angle) * radius)

    RENDER.CLEAR(20, 20, 20)
    DRAW.RECTANGLE(x - 15, y - 15, 30, 30, 200, 50, 150, 255)
    RENDER.FRAME()
WEND

WINDOW.CLOSE()
```

---

## Full Example: Randomized Star Field

```basic
WINDOW.OPEN(800, 600, "Star Field")
WINDOW.SETFPS(60)

RANDOMIZE

CONST STAR_COUNT = 200
DIM sx(STAR_COUNT)
DIM sy(STAR_COUNT)
DIM ss(STAR_COUNT)

FOR i = 1 TO STAR_COUNT
    sx(i) = MATH.RND(800)
    sy(i) = MATH.RND(600)
    ss(i) = MATH.RND(3) + 1
NEXT

WHILE NOT WINDOW.SHOULDCLOSE()
    RENDER.CLEAR(0, 0, 10)
    FOR i = 1 TO STAR_COUNT
        bright = 150 + MATH.RND(105)
        DRAW.RECTANGLE(sx(i), sy(i), ss(i), ss(i), bright, bright, bright, 255)
    NEXT
    RENDER.FRAME()
WEND

WINDOW.CLOSE()
```

---

## Extended Command Reference

### Trig (degree variants)

| Command | Description |
|--------|-------------|
| `MATH.SIND(deg)` | Sine of degrees. |
| `MATH.COSD(deg)` | Cosine of degrees. |
| `MATH.TAND(deg)` | Tangent of degrees. |
| `MATH.ASIN(v)` | Arc-sine in radians. |
| `MATH.ACOS(v)` | Arc-cosine in radians. |
| `MATH.ATAN(v)` | Arc-tangent in radians. |
| `MATH.ATAN2(y, x)` | Two-argument arc-tangent. |
| `MATH.ATN(v)` | Alias of `MATH.ATAN`. |

### Rounding & signs

| Command | Description |
|--------|-------------|
| `MATH.FLOOR(v)` | Round down to nearest integer. |
| `MATH.CEIL(v)` | Round up to nearest integer. |
| `MATH.ROUND(v)` | Round to nearest integer. |
| `MATH.FIX(v)` | Truncate toward zero. |
| `MATH.SIGN(v)` / `MATH.SGN(v)` | Returns -1, 0, or 1. |
| `MATH.SATURATE(v)` | Clamp to 0.0–1.0. |

### Logarithms & misc

| Command | Description |
|--------|-------------|
| `MATH.LOG(v)` | Natural logarithm. |
| `MATH.LOG10(v)` | Base-10 logarithm. |
| `MATH.LOG2(v)` | Base-2 logarithm. |
| `MATH.CHANCE(p)` | Returns `TRUE` with probability `p` (0.0–1.0). |
| `MATH.RNDSEED(seed)` | Seed the global RNG. |

### Interpolation helpers

| Command | Description |
|--------|-------------|
| `MATH.APPROACH(cur, target, step)` | Move `cur` toward `target` by at most `step`. |
| `MATH.CURVE(cur, target, factor)` | Exponential approach: `cur + (target-cur) * factor`. |
| `MATH.INVERSE_LERP(a, b, v)` | Returns `t` such that `LERP(a,b,t) == v`. |
| `MATH.LERPANGLE(a, b, t)` | Lerp between two angles (shortest path). |
| `MATH.PINGPONG(t, length)` | Bounce `t` back and forth 0..`length`. |
| `MATH.REMAP(v, s0, s1, d0, d1)` | Remap value from source range to destination range. |
| `MATH.SATURATE(v)` | Clamp 0–1 (also listed above). |
| `MATH.SMOOTHSTEP(edge0, edge1, x)` | Hermite smooth step. |
| `MATH.WRAP(v, min, max)` | Wrap value within `[min, max)`. |
| `MATH.WRAPANGLE(a)` | Wrap angle to 0..2π. |
| `MATH.WRAPANGLE180(a)` | Wrap angle to -π..π. |

### Heading helpers

| Command | Description |
|--------|-------------|
| `MATH.NEWX(yaw, dist)` | X offset for heading `yaw` (radians) and `dist`. |
| `MATH.NEWZ(yaw, dist)` | Z offset for heading `yaw` (radians) and `dist`. |
| `MATH.ANGLETO(x1, z1, x2, z2)` | Yaw angle (radians) from point 1 to point 2. |
| `MATH.DEGPERSEC(deg, dt)` | Converts degrees/second to this-frame radians: `deg * dt * π/180`. |
| `MATH.CIRCLEPOINT(cx, cz, r, i, n)` | Returns `[x, z]` of point `i` on a circle of `n` equally-spaced points. |

---

## See also (gameplay-oriented)

- [GAME_MATH_HELPERS.md](GAME_MATH_HELPERS.md) — **`HDIST` / `HDISTSQ`** (XZ distance), **`DIST2D` / `DISTSQ2D`**, **`YAWFROMXZ`**, **`ANGLEDIFFRAD`**, **`SMOOTHERSTEP`**
- [GAME_ENGINE_PATTERNS.md](GAME_ENGINE_PATTERNS.md) — **`MOVE.TOWARD`**, **`MOVE.LERP`**, **`ANGLE.DIFFERENCE`**, rays, lights, sprites, **`RES.*`**
- [LESS_MATH.md](LESS_MATH.md) — camera-relative movement, terrain snap, vector helpers
