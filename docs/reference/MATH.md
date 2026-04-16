# Math Commands

Commands for mathematical and numerical operations.

---

## Trigonometry

Angles are in **radians**. Use `Math.Deg2Rad()` / `Math.Rad2Deg()` to convert.

### `Math.Sin(angle)`
Returns the sine of an angle in **radians**. Alias: `SIN()`.

### `Math.Cos(angle)`
Returns the cosine of an angle in **radians**. Alias: `COS()`.

### `Math.Tan(angle)`
Returns the tangent of an angle in **radians**. Alias: `TAN()`.

---

## Powers, Roots & Logarithms

### `Math.Sqrt(value)`
Returns the square root of a non-negative value. Alias: `SQRT()`, `SQR()`.

### `Math.Abs(value)`
Returns the absolute value of a number. Alias: `ABS()`.

### `Math.Pow(base, exp)`
Returns base raised to the power of exp. Alias: `POW()`.

---

## Arithmetic Helpers

### `Math.Clamp(value, min, max)`
Constrains a value to the range `[min, max]`. Alias: `CLAMP()`.

### `Math.Lerp(a, b, t)`
Linearly interpolates between `a` and `b` by factor `t` (0.0–1.0). Alias: `LERP()`.

---

## Randomization

### `Rnd(limit)`
Returns a random integer from 0 up to `limit-1`.

### `RndF(min, max)`
Returns a random float between `min` and `max`.

---

## Constants

### `Math.Pi()`
Returns **π** (3.14159...). Alias: `PI()`.

### `Math.Tau()`
Returns **τ** (6.28318...). Alias: `TAU()`.

### `Math.E()`
Returns Euler's number (2.71828...). Alias: `E()`.

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

## See also (gameplay-oriented)

- [GAME_MATH_HELPERS.md](GAME_MATH_HELPERS.md) — **`HDIST` / `HDISTSQ`** (XZ distance), **`DIST2D` / `DISTSQ2D`**, **`YAWFROMXZ`**, **`ANGLEDIFFRAD`**, **`SMOOTHERSTEP`**
- [GAME_ENGINE_PATTERNS.md](GAME_ENGINE_PATTERNS.md) — **`MOVE.TOWARD`**, **`MOVE.LERP`**, **`ANGLE.DIFFERENCE`**, rays, lights, sprites, **`RES.*`**
- [LESS_MATH.md](LESS_MATH.md) — camera-relative movement, terrain snap, vector helpers
