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
Window.Open(800, 600, "Math Example: Circular Motion")
Window.SetFPS(60)

angle = 0.0
radius = 150.0
center_x = 400
center_y = 300

WHILE NOT Window.ShouldClose()
    angle = angle + 2.0 * Time.Delta()

    x = center_x + INT(Math.Cos(angle) * radius)
    y = center_y + INT(Math.Sin(angle) * radius)

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

Math.Randomize  ; Seed from system clock

CONST STAR_COUNT = 200
DIM sx(STAR_COUNT)
DIM sy(STAR_COUNT)
DIM ss(STAR_COUNT)

FOR i = 1 TO STAR_COUNT
    sx(i) = Math.Rnd(800)
    sy(i) = Math.Rnd(600)
    ss(i) = Math.Rnd(3) + 1
NEXT

WHILE NOT Window.ShouldClose()
    Render.Clear(0, 0, 10)
    FOR i = 1 TO STAR_COUNT
        bright = 150 + Math.Rnd(105)
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
